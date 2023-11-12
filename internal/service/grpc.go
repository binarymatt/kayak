package service

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"

	"connectrpc.com/connect"
	"github.com/bufbuild/protovalidate-go"
	"github.com/hashicorp/raft"
	"github.com/spaolacci/murmur3"
	"go.opentelemetry.io/otel"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/reflect/protoreflect"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/structpb"
	"log/slog"

	kayakv1 "github.com/binarymatt/kayak/gen/kayak/v1"
	"github.com/binarymatt/kayak/gen/kayak/v1/kayakv1connect"
	"github.com/binarymatt/kayak/internal/fsm"
)

var (
	ErrInvalidConsumer = errors.New("invalid consumer")
)

func validate(msg protoreflect.ProtoMessage) error {
	validator, err := protovalidate.New()
	if err != nil {
		return err
	}
	return validator.Validate(msg)
}
func (s *service) PutRecords(ctx context.Context, req *connect.Request[kayakv1.PutRecordsRequest]) (*connect.Response[emptypb.Empty], error) {
	fmt.Println("putting records")
	slog.Info("put records request")
	if err := validate(req.Msg); err != nil {
		slog.Error("invalid put request", "error", err)
		return nil, connect.NewError(connect.CodeInvalidArgument, err)
	}

	for _, record := range req.Msg.Records {
		id := s.idGenerator()
		record.Id = id.String()
		record.Topic = req.Msg.Topic
	}
	command := &kayakv1.Command{
		Payload: &kayakv1.Command_PutRecordsRequest{
			PutRecordsRequest: req.Msg,
		},
	}
	_, err := s.applyCommand(ctx, command)
	return connect.NewResponse(&emptypb.Empty{}), err
}

func (s *service) GetRecords(ctx context.Context, req *connect.Request[kayakv1.GetRecordsRequest]) (*connect.Response[kayakv1.GetRecordsResponse], error) {
	if err := validate(req.Msg); err != nil {
		slog.Error("invalid put request", "error", err)
		return nil, connect.NewError(connect.CodeInvalidArgument, err)
	}

	limit := int(req.Msg.Limit)
	def := false
	if limit == 0 {
		def = true
		limit = 99
	}
	slog.Info("getting records from store", "default", def, "limit", limit)
	records, err := s.store.GetRecords(ctx, req.Msg.Topic, req.Msg.Start, limit)
	if err != nil {
		slog.Error("error getting records from store", "error", err)
		return nil, connect.NewError(connect.CodeInvalidArgument, err)
	}
	return connect.NewResponse(&kayakv1.GetRecordsResponse{
		Records: records,
	}), nil
}
func findConsumer(consumers []*kayakv1.TopicConsumer, id string) *kayakv1.TopicConsumer {
	for _, consumer := range consumers {
		if consumer.Id == id {
			return consumer
		}
	}
	return nil
}

// return partition index that this message should live in
func balancer(key string, partitionCount int64) int64 {
	partition := murmur3.Sum64([]byte(key)) % uint64(partitionCount)
	return int64(partition)
}
func (s *service) FetchRecord(ctx context.Context, req *connect.Request[kayakv1.FetchRecordRequest]) (*connect.Response[kayakv1.FetchRecordsResponse], error) {
	logger := slog.With("consumer", req.Msg.ConsumerId, "topic", req.Msg.Topic)
	if err := validate(req.Msg); err != nil {
		slog.Error("invalid request", "error", err)
		return nil, connect.NewError(connect.CodeInvalidArgument, err)
	}
	meta, err := s.store.LoadMeta(ctx, req.Msg.Topic)
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}
	group, ok := meta.GroupMetadata[req.Msg.ConsumerGroup]
	if !ok {
		return nil, connect.NewError(connect.CodeInvalidArgument, ErrInvalidConsumer)
	}
	consumers := group.Consumers
	c := findConsumer(consumers, req.Msg.ConsumerId)
	if c == nil {
		return nil, connect.NewError(connect.CodeInvalidArgument, ErrInvalidConsumer)
	}
	logger.Info("fetching records via grpc")

	logger.Debug("checking consumer group details")

	slog.Info("retrieved consumer group position from local store", "fetched_position", c.Position)
	position := c.Position
	logger.Info("fetching record")

	records, err := s.store.GetRecords(ctx, req.Msg.Topic, position, 10)
	if err != nil {
		logger.Error("error getting records from store", "error", err)
		return nil, connect.NewError(connect.CodeInvalidArgument, err)
	}
	var record *kayakv1.Record
	for _, record := range records {
		fmt.Println(record)
		partition := balancer(record.Id, group.Partitions)
		if partition == c.Partition {
			return connect.NewResponse(&kayakv1.FetchRecordsResponse{
				Record: record,
			}), nil
		}
	}

	return connect.NewResponse(&kayakv1.FetchRecordsResponse{
		Record: record,
	}), nil
}

func (s *service) StreamRecords(ctx context.Context, req *connect.Request[kayakv1.StreamRecordsRequest], stream *connect.ServerStream[kayakv1.Record]) error {
	validator, err := protovalidate.New()
	if err != nil {
		return err
	}
	if err := validator.Validate(req.Msg); err != nil {
		return connect.NewError(connect.CodeInvalidArgument, err)
	}
	startTime := time.Now()
	startPosition := req.Msg.Position
	limit := req.Msg.BatchSize
	if limit == 0 {
		limit = 1
	}
	d := req.Msg.GetTimeout()
	var duration time.Duration
	if d == nil {
		duration = 1 * time.Second
	} else {
		duration = d.AsDuration()
	}
	for {
		records, err := s.store.GetRecords(ctx, req.Msg.Topic, startPosition, int(limit))
		if err != nil {
			return err
		}
		if len(records) > 0 {
			startTime = time.Now()
		}
		for _, record := range records {
			if err := stream.Send(record); err != nil {
				return err
			}
		}
		if startTime.Add(duration).After(time.Now()) {
			return nil
		}

	}
}

func (s *service) CommitRecord(ctx context.Context, req *connect.Request[kayakv1.CommitRecordRequest]) (*connect.Response[emptypb.Empty], error) {
	if err := req.Msg.Validate(); err != nil {
		return nil, connect.NewError(connect.CodeInvalidArgument, err)
	}
	command := &kayakv1.Command{
		Payload: &kayakv1.Command_CommitRecordRequest{
			CommitRecordRequest: req.Msg,
		},
	}
	_, err := s.applyCommand(ctx, command)
	return connect.NewResponse(&emptypb.Empty{}), err
}
func ToStructValue(item interface{}) (*structpb.Value, error) {
	// to map[string]interface{}
	raw, err := json.Marshal(item)
	if err != nil {
		return nil, err
	}
	var v map[string]interface{}
	if err := json.Unmarshal(raw, &v); err != nil {
		return nil, err
	}
	return structpb.NewValue(v)

}
func (s *service) applyCommand(ctx context.Context, command *kayakv1.Command) (*connect.Response[kayakv1.ApplyResponse], error) {
	if s.raft.State() != raft.Leader {
		fmt.Println("not leader")
		leader := fmt.Sprintf("http://%s", s.raft.Leader())
		client := kayakv1connect.NewKayakServiceClient(http.DefaultClient, leader)
		slog.InfoContext(ctx, "applying to leader", "leader", s.raft.Leader())
		return client.Apply(ctx, connect.NewRequest(command))
	}
	data, err := protojson.Marshal(command)
	if err != nil {
		return nil, connect.NewError(connect.CodeInvalidArgument, err)
	}
	applyFuture := s.raft.Apply(data, 500*time.Millisecond)
	if err := applyFuture.Error(); err != nil {
		slog.ErrorContext(ctx, "could not apply command to raft", "error", err)
		//	return nil, connect.NewError(connect.CodeInternal, err)
		return nil, err
	}
	var val *structpb.Value
	if applyFuture.Response() != nil {
		resp := applyFuture.Response().(*fsm.ApplyResponse)
		val, err = ToStructValue(resp.Data)
		if err != nil {
			return nil, err
		}
		//val, err = structpb.NewValue(resp.Data)

		//anypb.New(resp.Data)
		//if err != nil {
		//	fmt.Println("error getting struct value")
		//	return nil, err
		//}
	}
	return connect.NewResponse(&kayakv1.ApplyResponse{
		Data: val,
	}), applyFuture.Error()
}

func (s *service) DeleteTopic(ctx context.Context, req *connect.Request[kayakv1.DeleteTopicRequest]) (*connect.Response[emptypb.Empty], error) {
	validator, err := protovalidate.New()
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}
	if err := validator.Validate(req.Msg); err != nil {
		return nil, connect.NewError(connect.CodeInvalidArgument, err)
	}
	command := &kayakv1.Command{
		Payload: &kayakv1.Command_DeleteTopicRequest{
			DeleteTopicRequest: req.Msg,
		},
	}
	_, err = s.applyCommand(ctx, command)
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}
	return connect.NewResponse(&emptypb.Empty{}), nil
}
func (s *service) CreateTopic(ctx context.Context, req *connect.Request[kayakv1.CreateTopicRequest]) (*connect.Response[emptypb.Empty], error) {
	if err := req.Msg.Validate(); err != nil {
		return nil, connect.NewError(connect.CodeInvalidArgument, err)
	}
	slog.InfoContext(ctx, "creating topic", "topic", req.Msg.Name)

	command := &kayakv1.Command{
		Payload: &kayakv1.Command_CreateTopicRequest{
			CreateTopicRequest: req.Msg,
		},
	}
	_, err := s.applyCommand(ctx, command)
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, fmt.Errorf("could not apply to raft: %w", err))
	}
	return connect.NewResponse(&emptypb.Empty{}), nil
}

func (s *service) ListTopics(ctx context.Context, req *connect.Request[kayakv1.ListTopicsRequest]) (*connect.Response[kayakv1.ListTopicsResponse], error) {
	names, err := s.store.ListTopics(ctx)
	if err != nil {
		slog.Error("error with store ListTopics", "error", err)
	}
	var topics []*kayakv1.Topic
	for _, name := range names {
		topics = append(topics, &kayakv1.Topic{
			Name: name,
		})
	}
	return connect.NewResponse(&kayakv1.ListTopicsResponse{Topics: topics}), nil
}

func (s *service) Stats(ctx context.Context, _ *connect.Request[emptypb.Empty]) (*connect.Response[kayakv1.StatsResponse], error) {
	_, span := otel.GetTracerProvider().Tracer("").Start(ctx, "stats-span")
	defer span.End()
	return connect.NewResponse(&kayakv1.StatsResponse{
		Raft:  s.raft.Stats(),
		Store: s.store.Stats(),
	}), nil
}

func (s *service) GetNodeDetails(ctx context.Context, _ *connect.Request[emptypb.Empty]) (*connect.Response[kayakv1.GetNodeDetailsResponse], error) {
	return connect.NewResponse(&kayakv1.GetNodeDetailsResponse{
		Id: s.cfg.ServerID,
	}), nil
}

func (s *service) Apply(ctx context.Context, req *connect.Request[kayakv1.Command]) (*connect.Response[kayakv1.ApplyResponse], error) {
	_, span := otel.GetTracerProvider().Tracer("").Start(ctx, "apply-span")
	defer span.End()
	if req.Msg.GetCreateTopicRequest() != nil {
		_, err := s.CreateTopic(ctx, connect.NewRequest(req.Msg.GetCreateTopicRequest()))
		return connect.NewResponse(&kayakv1.ApplyResponse{}), err
	}
	if req.Msg.GetPutRecordsRequest() != nil {
		fmt.Println("apply grpc")
		_, err := s.PutRecords(ctx, connect.NewRequest(req.Msg.GetPutRecordsRequest()))
		return connect.NewResponse(&kayakv1.ApplyResponse{}), err
	}
	if req.Msg.GetCommitRecordRequest() != nil {
		_, err := s.CommitRecord(ctx, connect.NewRequest(req.Msg.GetCommitRecordRequest()))
		return connect.NewResponse(&kayakv1.ApplyResponse{}), err
	}
	if req.Msg.GetDeleteTopicRequest() != nil {
		_, err := s.DeleteTopic(ctx, connect.NewRequest(req.Msg.GetDeleteTopicRequest()))
		return connect.NewResponse(&kayakv1.ApplyResponse{}), err
	}
	// CreateConsumerGroup
	if req.Msg.GetCreateConsumerGroupRequest() != nil {
		_, err := s.CreateConsumerGroup(ctx, connect.NewRequest(req.Msg.GetCreateConsumerGroupRequest()))
		return connect.NewResponse(&kayakv1.ApplyResponse{}), err
	}

	// RegisterConsumer
	if req.Msg.GetRegisterConsumerRequest() != nil {
		_, err := s.RegisterConsumer(ctx, connect.NewRequest(req.Msg.GetRegisterConsumerRequest()))
		return connect.NewResponse(&kayakv1.ApplyResponse{}), err
	}

	return nil, connect.NewError(connect.CodeInvalidArgument, errors.New("not valid"))
}

func (s *service) CreateConsumerGroup(ctx context.Context, req *connect.Request[kayakv1.CreateConsumerGroupRequest]) (*connect.Response[emptypb.Empty], error) {
	if err := validate(req.Msg); err != nil {
		slog.Error("invalid create consumer group request", "error", err)
		return nil, connect.NewError(connect.CodeInvalidArgument, err)
	}
	command := &kayakv1.Command{
		Payload: &kayakv1.Command_CreateConsumerGroupRequest{
			CreateConsumerGroupRequest: req.Msg,
		},
	}
	_, err := s.applyCommand(ctx, command)
	return connect.NewResponse(&emptypb.Empty{}), err

}
func (s *service) RegisterConsumer(ctx context.Context, req *connect.Request[kayakv1.RegisterConsumerRequest]) (*connect.Response[emptypb.Empty], error) {
	if err := validate(req.Msg); err != nil {
		slog.Error("invalid create consumer group request", "error", err)
		return nil, connect.NewError(connect.CodeInvalidArgument, err)
	}
	command := &kayakv1.Command{
		Payload: &kayakv1.Command_RegisterConsumerRequest{
			RegisterConsumerRequest: req.Msg,
		},
	}
	resp, err := s.applyCommand(ctx, command)
	slog.Debug("apply response", "response", resp)
	return connect.NewResponse(&emptypb.Empty{}), err
}
