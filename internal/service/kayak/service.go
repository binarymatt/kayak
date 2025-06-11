package kayak

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"buf.build/go/protovalidate"
	"connectrpc.com/connect"
	"github.com/coder/quartz"
	"github.com/dgraph-io/badger/v4"
	"github.com/hashicorp/raft"
	"github.com/oklog/ulid/v2"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/timestamppb"

	v1 "github.com/binarymatt/kayak/gen/kayak/v1"
	"github.com/binarymatt/kayak/gen/kayak/v1/kayakv1connect"
	internal_raft "github.com/binarymatt/kayak/internal/raft"
	"github.com/binarymatt/kayak/internal/store"
)

var (
	_            kayakv1connect.KayakServiceHandler = (*service)(nil)
	ErrNotLeader                                    = errors.New("node is not the leader")
)

type service struct {
	idGenerator      func() ulid.ULID
	store            store.Store
	workerExpiration time.Duration
	logger           *slog.Logger
	raft             internal_raft.RaftInterface
	testLeaderClient kayakv1connect.KayakServiceClient
	clock            quartz.Clock
}

func (s *service) applyCommand(ctx context.Context, cmd *v1.RaftCommand) error {
	if s.raft.State() != raft.Leader {
		slog.Info("sending apply to leader")
		client := s.getLeaderClient()
		if _, err := client.Apply(ctx, connect.NewRequest(&v1.ApplyRequest{Command: cmd})); err != nil {
			return err
		}
		return nil
	}

	data, err := proto.Marshal(cmd)
	if err != nil {
		return err
	}
	slog.Debug("apply to raft members")
	applyFuture := s.raft.Apply(data, 10*time.Millisecond)
	if err := applyFuture.Error(); err != nil {
		s.logger.Error("could not apply command to raft", "error", err)
		//	return nil, connect.NewError(connect.CodeInternal, err)
		return err
	}
	resp := applyFuture.Response()
	if resp != nil {
		as := resp.(*store.ApplyResponse)
		return as.Error
	}
	return errors.New("not sure why we are here")

}
func (s *service) PutRecords(ctx context.Context, req *connect.Request[v1.PutRecordsRequest]) (*connect.Response[emptypb.Empty], error) {
	s.logger.Info("put records request", "count", len(req.Msg.Records), "raft", s.raft)

	stream, err := s.store.GetStream(req.Msg.StreamName)
	if err != nil {
		slog.Error("could not get stream info", "error", err)
		return nil, err
	}
	for _, r := range req.Msg.Records {
		id := s.idGenerator()
		r.InternalId = id.String()
		if r.Id == "" {
			r.Id = id.String()
		}
		r.StreamName = stream.Name
		r.Partition = balancer(r.Id, stream.PartitionCount)
		r.AcceptTimestamp = timestamppb.New(s.clock.Now())
	}

	//err = s.store.PutRecords(req.Msg.StreamName, req.Msg.GetRecords()...)
	cmd := &v1.RaftCommand{
		Payload: &v1.RaftCommand_PutRecords{
			PutRecords: &v1.PutRecords{
				Records:    req.Msg.GetRecords(),
				StreamName: req.Msg.StreamName,
			},
		},
	}

	err = s.applyCommand(ctx, cmd)
	return connect.NewResponse(&emptypb.Empty{}), err
}

func (s *service) CommitRecord(ctx context.Context, req *connect.Request[v1.CommitRecordRequest]) (*connect.Response[emptypb.Empty], error) {
	slog.Info("commiting record", "msg", req.Msg)
	// check worker lease
	if err := s.store.HasLease(req.Msg.Worker); err != nil {
		return nil, connect.NewError(connect.CodeInvalidArgument, err)
	}
	err := s.applyCommand(ctx, &v1.RaftCommand{
		Payload: &v1.RaftCommand_ExtendLease{
			ExtendLease: &v1.ExtendLease{
				Worker:    req.Msg.Worker,
				ExpiresMs: s.clock.Now().Add(s.workerExpiration).UnixMilli(),
			},
		},
	})
	if err != nil {
		// if err := s.store.ExtendLease(req.Msg.Worker, s.workerExpiration); err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}
	worker := req.Msg.GetWorker()
	record := req.Msg.GetRecord()
	id, err := ulid.Parse(record.InternalId)
	if err != nil {
		slog.Error("could not parse internal ID", "error", err, "record")
		return nil, connect.NewError(connect.CodeInternal, err)
	}
	position := id.String()
	err = s.applyCommand(ctx, &v1.RaftCommand{
		Payload: &v1.RaftCommand_CommitGroupPosition{
			CommitGroupPosition: &v1.CommitGroupPosition{
				StreamName: worker.StreamName,
				GroupName:  worker.GroupName,
				Partition:  worker.PartitionAssignment,
				Position:   position,
			},
		},
	})
	if err != nil {
		//if err := s.store.CommitGroupPosition(worker.StreamName, worker.GroupName, worker.PartitionAssignment, position); err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}
	return connect.NewResponse(&emptypb.Empty{}), nil
}

func (s *service) FetchRecords(ctx context.Context, req *connect.Request[v1.FetchRecordsRequest]) (*connect.Response[v1.FetchRecordsResponse], error) {
	worker := req.Msg.Worker
	// is worker registered?
	if err := s.store.HasLease(worker); err != nil {
		return nil, connect.NewError(connect.CodeInvalidArgument, err)
	}

	// TODO: cache position in worker client and  use if not too old
	position, err := s.store.GetGroupPosition(req.Msg.StreamName, worker.GroupName, worker.PartitionAssignment)
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}
	records, err := s.store.GetRecords(req.Msg.StreamName, worker.PartitionAssignment, position, int(req.Msg.Limit))

	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}
	// Get records based on worker info
	return connect.NewResponse(&v1.FetchRecordsResponse{Records: records}), nil
}

func (s *service) GetRecords(ctx context.Context, req *connect.Request[v1.GetRecordsRequest]) (*connect.Response[v1.GetRecordsResponse], error) {
	slog.Info("getting records", "partitio", req.Msg.Partition, "start", req.Msg.StartId, "stream", req.Msg.StreamName, "limit", req.Msg.Limit)
	records, err := s.store.GetRecords(req.Msg.StreamName, req.Msg.Partition, req.Msg.StartId, int(req.Msg.Limit))
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}
	return connect.NewResponse(&v1.GetRecordsResponse{
		Records: records,
	}), nil
}

func (s *service) CreateStream(ctx context.Context, req *connect.Request[v1.CreateStreamRequest]) (*connect.Response[emptypb.Empty], error) {
	s.logger.Info("creating stream", "name", req.Msg.Name, "partitions", req.Msg.PartitionCount)
	// TODO: validate request
	stream := &v1.Stream{
		Name:           req.Msg.Name,
		PartitionCount: req.Msg.PartitionCount,
		Ttl:            req.Msg.Ttl,
	}
	cmd := &v1.RaftCommand{
		Payload: &v1.RaftCommand_PutStream{
			PutStream: &v1.PutStream{
				Stream: stream,
			},
		},
	}
	if err := s.applyCommand(ctx, cmd); err != nil {
		slog.Error("could not create stream", "error", err)
		return nil, connect.NewError(connect.CodeInternal, err)
	}
	return connect.NewResponse(&emptypb.Empty{}), nil
}

func (s *service) GetStream(ctx context.Context, req *connect.Request[v1.GetStreamRequest]) (*connect.Response[v1.GetStreamResponse], error) {
	stream, err := s.store.GetStream(req.Msg.Name)
	if err != nil {
		if errors.Is(err, badger.ErrKeyNotFound) {
			return nil, connect.NewError(connect.CodeInvalidArgument, err)
		}
		return nil, connect.NewError(connect.CodeInternal, err)
	}
	stats, err := s.store.GetStreamStats(stream.Name)
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}
	stream.Stats = stats
	return connect.NewResponse(&v1.GetStreamResponse{
		Stream: stream,
	}), nil
}

func (s *service) GetStreams(ctx context.Context, req *connect.Request[v1.GetStreamsRequest]) (*connect.Response[v1.GetStreamsResponse], error) {
	streams, err := s.store.GetStreams()
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}
	resp := &v1.GetStreamsResponse{
		Streams: streams,
	}
	return connect.NewResponse(resp), nil
}

func (s *service) RegisterWorker(ctx context.Context, req *connect.Request[v1.RegisterWorkerRequest]) (*connect.Response[v1.RegisterWorkerResponse], error) {
	worker := &v1.Worker{
		StreamName: req.Msg.StreamName,
		GroupName:  req.Msg.Group,
		Id:         req.Msg.Id,
	}
	stream, err := s.store.GetStream(req.Msg.StreamName)
	if err != nil {
		if errors.Is(err, badger.ErrKeyNotFound) {
			return nil, connect.NewError(connect.CodeInvalidArgument, err)
		}
		return nil, connect.NewError(connect.CodeInternal, err)
	}
	assignments, err := s.store.GetPartitionAssignments(req.Msg.StreamName, req.Msg.Group)
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}
	if len(assignments) < int(stream.PartitionCount) {
		// get first available partition
		for i := range stream.PartitionCount {
			_, ok := assignments[i]
			if !ok {
				worker.PartitionAssignment = i
				break
			}
		}
	} else {
		return nil, connect.NewError(connect.CodeOutOfRange, errors.New("no partitions available"))
	}
	n := s.clock.Now().Add(s.workerExpiration)
	worker.LeaseExpires = n.UnixMilli()
	cmd := &v1.RaftCommand{
		Payload: &v1.RaftCommand_ExtendLease{
			ExtendLease: &v1.ExtendLease{
				Worker:    worker,
				ExpiresMs: worker.LeaseExpires,
			},
		},
	}
	if err := s.applyCommand(ctx, cmd); err != nil {
		//if err := s.store.ExtendLease(worker, s.workerExpiration); err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}
	resp := &v1.RegisterWorkerResponse{
		Worker: worker,
	}
	return connect.NewResponse(resp), nil
}

var (
	ErrNoAssignmentMatch = errors.New("assignment does not match")
)

func (s *service) RenewRegistration(ctx context.Context, req *connect.Request[v1.RenewRegistrationRequest]) (*connect.Response[emptypb.Empty], error) {
	if err := protovalidate.Validate(req.Msg); err != nil {
		fmt.Println("invalid request")
		return nil, connect.NewError(connect.CodeInvalidArgument, err)
	}

	worker := req.Msg.GetWorker()
	// check partition assignment
	id, err := s.store.GetPartitionAssignment(worker.GetStreamName(), worker.GetGroupName(), worker.GetPartitionAssignment())
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}
	if id != req.Msg.Worker.Id {
		return nil, connect.NewError(connect.CodeInvalidArgument, fmt.Errorf("partition is assigned to a different id %s: %w", id, ErrNoAssignmentMatch))
	}
	n := s.clock.Now().Add(s.workerExpiration)
	worker.LeaseExpires = n.UnixMilli()
	cmd := &v1.RaftCommand{
		Payload: &v1.RaftCommand_ExtendLease{
			ExtendLease: &v1.ExtendLease{
				Worker:    worker,
				ExpiresMs: worker.LeaseExpires,
			},
		},
	}
	if err := s.applyCommand(ctx, cmd); err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}

	return connect.NewResponse(&emptypb.Empty{}), nil
}
func (s *service) DeregisterWorker(ctx context.Context, req *connect.Request[v1.DeregisterWorkerRequest]) (*connect.Response[emptypb.Empty], error) {
	cmd := &v1.RaftCommand{
		Payload: &v1.RaftCommand_RemoveLease{
			RemoveLease: &v1.RemoveLease{
				Worker: req.Msg.GetWorker(),
			},
		},
	}
	if err := s.applyCommand(ctx, cmd); err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}
	return connect.NewResponse(&emptypb.Empty{}), nil
}

func (s *service) Apply(ctx context.Context, req *connect.Request[v1.ApplyRequest]) (*connect.Response[v1.ApplyResponse], error) {
	if err := s.applyCommand(ctx, req.Msg.GetCommand()); err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}
	return connect.NewResponse(&v1.ApplyResponse{}), nil
}

func (s *service) DeleteStream(ctx context.Context, req *connect.Request[v1.DeleteStreamRequest]) (*connect.Response[emptypb.Empty], error) {

	// TODO: validate request
	cmd := &v1.RaftCommand{
		Payload: &v1.RaftCommand_DeleteStream{
			DeleteStream: &v1.DeleteStream{
				StreamName: req.Msg.Name,
			},
		},
	}
	if err := s.applyCommand(ctx, cmd); err != nil {
		slog.Error("could not create stream", "error", err)
		return nil, connect.NewError(connect.CodeInternal, err)
	}
	return connect.NewResponse(&emptypb.Empty{}), nil
}

func (s *service) GetStreamStatistics(ctx context.Context, req *connect.Request[v1.GetStreamStatisticsRequest]) (*connect.Response[v1.GetStreamStatisticsResponse], error) {
	if err := protovalidate.Validate(req.Msg); err != nil {
		fmt.Println("invalid request")
		return nil, connect.NewError(connect.CodeInvalidArgument, err)
	}
	slog.Info("getting stream statistics", "stream", req.Msg.Name)

	stats, err := s.store.GetStreamStats(req.Msg.Name)
	if err != nil {
		slog.Error("could not retreive stream stats", "error", err)
		return nil, connect.NewError(connect.CodeInternal, err)
	}
	return connect.NewResponse(&v1.GetStreamStatisticsResponse{StreamStats: stats}), nil
}
func (s *service) getLeaderClient() kayakv1connect.KayakServiceClient {
	if s.testLeaderClient != nil {
		return s.testLeaderClient
	}

	leader := fmt.Sprintf("http://%s", s.raft.Leader())
	fmt.Println(leader)
	client := kayakv1connect.NewKayakServiceClient(http.DefaultClient, leader)
	return client
}

func New(st store.Store, ra *raft.Raft) *service {
	return &service{
		raft:             ra,
		idGenerator:      ulid.Make,
		workerExpiration: 60 * time.Second,
		store:            st,
		logger:           slog.Default().With("service", "kayak"),
		clock:            quartz.NewReal(),
	}
}
