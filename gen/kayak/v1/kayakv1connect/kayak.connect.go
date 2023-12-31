// Code generated by protoc-gen-connect-go. DO NOT EDIT.
//
// Source: proto/kayak/v1/kayak.proto

package kayakv1connect

import (
	connect "connectrpc.com/connect"
	context "context"
	errors "errors"
	v1 "github.com/binarymatt/kayak/gen/kayak/v1"
	emptypb "google.golang.org/protobuf/types/known/emptypb"
	http "net/http"
	strings "strings"
)

// This is a compile-time assertion to ensure that this generated file and the connect package are
// compatible. If you get a compiler error that this constant is not defined, this code was
// generated with a version of connect newer than the one compiled into your binary. You can fix the
// problem by either regenerating this code with an older version of connect or updating the connect
// version compiled into your binary.
const _ = connect.IsAtLeastVersion0_1_0

const (
	// KayakServiceName is the fully-qualified name of the KayakService service.
	KayakServiceName = "kayak.v1.KayakService"
)

// These constants are the fully-qualified names of the RPCs defined in this package. They're
// exposed at runtime as Spec.Procedure and as the final two segments of the HTTP route.
//
// Note that these are different from the fully-qualified method names used by
// google.golang.org/protobuf/reflect/protoreflect. To convert from these constants to
// reflection-formatted method names, remove the leading slash and convert the remaining slash to a
// period.
const (
	// KayakServiceApplyProcedure is the fully-qualified name of the KayakService's Apply RPC.
	KayakServiceApplyProcedure = "/kayak.v1.KayakService/Apply"
	// KayakServicePutRecordsProcedure is the fully-qualified name of the KayakService's PutRecords RPC.
	KayakServicePutRecordsProcedure = "/kayak.v1.KayakService/PutRecords"
	// KayakServiceCommitRecordProcedure is the fully-qualified name of the KayakService's CommitRecord
	// RPC.
	KayakServiceCommitRecordProcedure = "/kayak.v1.KayakService/CommitRecord"
	// KayakServiceCreateTopicProcedure is the fully-qualified name of the KayakService's CreateTopic
	// RPC.
	KayakServiceCreateTopicProcedure = "/kayak.v1.KayakService/CreateTopic"
	// KayakServiceDeleteTopicProcedure is the fully-qualified name of the KayakService's DeleteTopic
	// RPC.
	KayakServiceDeleteTopicProcedure = "/kayak.v1.KayakService/DeleteTopic"
	// KayakServiceRegisterConsumerProcedure is the fully-qualified name of the KayakService's
	// RegisterConsumer RPC.
	KayakServiceRegisterConsumerProcedure = "/kayak.v1.KayakService/RegisterConsumer"
	// KayakServiceGetRecordsProcedure is the fully-qualified name of the KayakService's GetRecords RPC.
	KayakServiceGetRecordsProcedure = "/kayak.v1.KayakService/GetRecords"
	// KayakServiceFetchRecordProcedure is the fully-qualified name of the KayakService's FetchRecord
	// RPC.
	KayakServiceFetchRecordProcedure = "/kayak.v1.KayakService/FetchRecord"
	// KayakServiceStreamRecordsProcedure is the fully-qualified name of the KayakService's
	// StreamRecords RPC.
	KayakServiceStreamRecordsProcedure = "/kayak.v1.KayakService/StreamRecords"
	// KayakServiceListTopicsProcedure is the fully-qualified name of the KayakService's ListTopics RPC.
	KayakServiceListTopicsProcedure = "/kayak.v1.KayakService/ListTopics"
	// KayakServiceStatsProcedure is the fully-qualified name of the KayakService's Stats RPC.
	KayakServiceStatsProcedure = "/kayak.v1.KayakService/Stats"
	// KayakServiceGetNodeDetailsProcedure is the fully-qualified name of the KayakService's
	// GetNodeDetails RPC.
	KayakServiceGetNodeDetailsProcedure = "/kayak.v1.KayakService/GetNodeDetails"
)

// KayakServiceClient is a client for the kayak.v1.KayakService service.
type KayakServiceClient interface {
	// Apply applies a command on the primary node.
	Apply(context.Context, *connect.Request[v1.Command]) (*connect.Response[v1.ApplyResponse], error)
	// PutRecords adds records to the stream
	PutRecords(context.Context, *connect.Request[v1.PutRecordsRequest]) (*connect.Response[emptypb.Empty], error)
	// Commits Consumer position for a topic/group
	CommitRecord(context.Context, *connect.Request[v1.CommitRecordRequest]) (*connect.Response[emptypb.Empty], error)
	// Creates new topic on server
	CreateTopic(context.Context, *connect.Request[v1.CreateTopicRequest]) (*connect.Response[emptypb.Empty], error)
	// Deletes Topic across server - permantly or via archive
	DeleteTopic(context.Context, *connect.Request[v1.DeleteTopicRequest]) (*connect.Response[emptypb.Empty], error)
	// rpc CreateConsumerGroup(CreateConsumerGroupRequest) returns (google.protobuf.Empty) {}
	RegisterConsumer(context.Context, *connect.Request[v1.RegisterConsumerRequest]) (*connect.Response[emptypb.Empty], error)
	// Read Procedures
	GetRecords(context.Context, *connect.Request[v1.GetRecordsRequest]) (*connect.Response[v1.GetRecordsResponse], error)
	FetchRecord(context.Context, *connect.Request[v1.FetchRecordRequest]) (*connect.Response[v1.FetchRecordsResponse], error)
	StreamRecords(context.Context, *connect.Request[v1.StreamRecordsRequest]) (*connect.ServerStreamForClient[v1.Record], error)
	ListTopics(context.Context, *connect.Request[v1.ListTopicsRequest]) (*connect.Response[v1.ListTopicsResponse], error)
	Stats(context.Context, *connect.Request[emptypb.Empty]) (*connect.Response[v1.StatsResponse], error)
	GetNodeDetails(context.Context, *connect.Request[emptypb.Empty]) (*connect.Response[v1.GetNodeDetailsResponse], error)
}

// NewKayakServiceClient constructs a client for the kayak.v1.KayakService service. By default, it
// uses the Connect protocol with the binary Protobuf Codec, asks for gzipped responses, and sends
// uncompressed requests. To use the gRPC or gRPC-Web protocols, supply the connect.WithGRPC() or
// connect.WithGRPCWeb() options.
//
// The URL supplied here should be the base URL for the Connect or gRPC server (for example,
// http://api.acme.com or https://acme.com/grpc).
func NewKayakServiceClient(httpClient connect.HTTPClient, baseURL string, opts ...connect.ClientOption) KayakServiceClient {
	baseURL = strings.TrimRight(baseURL, "/")
	return &kayakServiceClient{
		apply: connect.NewClient[v1.Command, v1.ApplyResponse](
			httpClient,
			baseURL+KayakServiceApplyProcedure,
			opts...,
		),
		putRecords: connect.NewClient[v1.PutRecordsRequest, emptypb.Empty](
			httpClient,
			baseURL+KayakServicePutRecordsProcedure,
			opts...,
		),
		commitRecord: connect.NewClient[v1.CommitRecordRequest, emptypb.Empty](
			httpClient,
			baseURL+KayakServiceCommitRecordProcedure,
			opts...,
		),
		createTopic: connect.NewClient[v1.CreateTopicRequest, emptypb.Empty](
			httpClient,
			baseURL+KayakServiceCreateTopicProcedure,
			opts...,
		),
		deleteTopic: connect.NewClient[v1.DeleteTopicRequest, emptypb.Empty](
			httpClient,
			baseURL+KayakServiceDeleteTopicProcedure,
			opts...,
		),
		registerConsumer: connect.NewClient[v1.RegisterConsumerRequest, emptypb.Empty](
			httpClient,
			baseURL+KayakServiceRegisterConsumerProcedure,
			opts...,
		),
		getRecords: connect.NewClient[v1.GetRecordsRequest, v1.GetRecordsResponse](
			httpClient,
			baseURL+KayakServiceGetRecordsProcedure,
			opts...,
		),
		fetchRecord: connect.NewClient[v1.FetchRecordRequest, v1.FetchRecordsResponse](
			httpClient,
			baseURL+KayakServiceFetchRecordProcedure,
			opts...,
		),
		streamRecords: connect.NewClient[v1.StreamRecordsRequest, v1.Record](
			httpClient,
			baseURL+KayakServiceStreamRecordsProcedure,
			opts...,
		),
		listTopics: connect.NewClient[v1.ListTopicsRequest, v1.ListTopicsResponse](
			httpClient,
			baseURL+KayakServiceListTopicsProcedure,
			opts...,
		),
		stats: connect.NewClient[emptypb.Empty, v1.StatsResponse](
			httpClient,
			baseURL+KayakServiceStatsProcedure,
			opts...,
		),
		getNodeDetails: connect.NewClient[emptypb.Empty, v1.GetNodeDetailsResponse](
			httpClient,
			baseURL+KayakServiceGetNodeDetailsProcedure,
			opts...,
		),
	}
}

// kayakServiceClient implements KayakServiceClient.
type kayakServiceClient struct {
	apply            *connect.Client[v1.Command, v1.ApplyResponse]
	putRecords       *connect.Client[v1.PutRecordsRequest, emptypb.Empty]
	commitRecord     *connect.Client[v1.CommitRecordRequest, emptypb.Empty]
	createTopic      *connect.Client[v1.CreateTopicRequest, emptypb.Empty]
	deleteTopic      *connect.Client[v1.DeleteTopicRequest, emptypb.Empty]
	registerConsumer *connect.Client[v1.RegisterConsumerRequest, emptypb.Empty]
	getRecords       *connect.Client[v1.GetRecordsRequest, v1.GetRecordsResponse]
	fetchRecord      *connect.Client[v1.FetchRecordRequest, v1.FetchRecordsResponse]
	streamRecords    *connect.Client[v1.StreamRecordsRequest, v1.Record]
	listTopics       *connect.Client[v1.ListTopicsRequest, v1.ListTopicsResponse]
	stats            *connect.Client[emptypb.Empty, v1.StatsResponse]
	getNodeDetails   *connect.Client[emptypb.Empty, v1.GetNodeDetailsResponse]
}

// Apply calls kayak.v1.KayakService.Apply.
func (c *kayakServiceClient) Apply(ctx context.Context, req *connect.Request[v1.Command]) (*connect.Response[v1.ApplyResponse], error) {
	return c.apply.CallUnary(ctx, req)
}

// PutRecords calls kayak.v1.KayakService.PutRecords.
func (c *kayakServiceClient) PutRecords(ctx context.Context, req *connect.Request[v1.PutRecordsRequest]) (*connect.Response[emptypb.Empty], error) {
	return c.putRecords.CallUnary(ctx, req)
}

// CommitRecord calls kayak.v1.KayakService.CommitRecord.
func (c *kayakServiceClient) CommitRecord(ctx context.Context, req *connect.Request[v1.CommitRecordRequest]) (*connect.Response[emptypb.Empty], error) {
	return c.commitRecord.CallUnary(ctx, req)
}

// CreateTopic calls kayak.v1.KayakService.CreateTopic.
func (c *kayakServiceClient) CreateTopic(ctx context.Context, req *connect.Request[v1.CreateTopicRequest]) (*connect.Response[emptypb.Empty], error) {
	return c.createTopic.CallUnary(ctx, req)
}

// DeleteTopic calls kayak.v1.KayakService.DeleteTopic.
func (c *kayakServiceClient) DeleteTopic(ctx context.Context, req *connect.Request[v1.DeleteTopicRequest]) (*connect.Response[emptypb.Empty], error) {
	return c.deleteTopic.CallUnary(ctx, req)
}

// RegisterConsumer calls kayak.v1.KayakService.RegisterConsumer.
func (c *kayakServiceClient) RegisterConsumer(ctx context.Context, req *connect.Request[v1.RegisterConsumerRequest]) (*connect.Response[emptypb.Empty], error) {
	return c.registerConsumer.CallUnary(ctx, req)
}

// GetRecords calls kayak.v1.KayakService.GetRecords.
func (c *kayakServiceClient) GetRecords(ctx context.Context, req *connect.Request[v1.GetRecordsRequest]) (*connect.Response[v1.GetRecordsResponse], error) {
	return c.getRecords.CallUnary(ctx, req)
}

// FetchRecord calls kayak.v1.KayakService.FetchRecord.
func (c *kayakServiceClient) FetchRecord(ctx context.Context, req *connect.Request[v1.FetchRecordRequest]) (*connect.Response[v1.FetchRecordsResponse], error) {
	return c.fetchRecord.CallUnary(ctx, req)
}

// StreamRecords calls kayak.v1.KayakService.StreamRecords.
func (c *kayakServiceClient) StreamRecords(ctx context.Context, req *connect.Request[v1.StreamRecordsRequest]) (*connect.ServerStreamForClient[v1.Record], error) {
	return c.streamRecords.CallServerStream(ctx, req)
}

// ListTopics calls kayak.v1.KayakService.ListTopics.
func (c *kayakServiceClient) ListTopics(ctx context.Context, req *connect.Request[v1.ListTopicsRequest]) (*connect.Response[v1.ListTopicsResponse], error) {
	return c.listTopics.CallUnary(ctx, req)
}

// Stats calls kayak.v1.KayakService.Stats.
func (c *kayakServiceClient) Stats(ctx context.Context, req *connect.Request[emptypb.Empty]) (*connect.Response[v1.StatsResponse], error) {
	return c.stats.CallUnary(ctx, req)
}

// GetNodeDetails calls kayak.v1.KayakService.GetNodeDetails.
func (c *kayakServiceClient) GetNodeDetails(ctx context.Context, req *connect.Request[emptypb.Empty]) (*connect.Response[v1.GetNodeDetailsResponse], error) {
	return c.getNodeDetails.CallUnary(ctx, req)
}

// KayakServiceHandler is an implementation of the kayak.v1.KayakService service.
type KayakServiceHandler interface {
	// Apply applies a command on the primary node.
	Apply(context.Context, *connect.Request[v1.Command]) (*connect.Response[v1.ApplyResponse], error)
	// PutRecords adds records to the stream
	PutRecords(context.Context, *connect.Request[v1.PutRecordsRequest]) (*connect.Response[emptypb.Empty], error)
	// Commits Consumer position for a topic/group
	CommitRecord(context.Context, *connect.Request[v1.CommitRecordRequest]) (*connect.Response[emptypb.Empty], error)
	// Creates new topic on server
	CreateTopic(context.Context, *connect.Request[v1.CreateTopicRequest]) (*connect.Response[emptypb.Empty], error)
	// Deletes Topic across server - permantly or via archive
	DeleteTopic(context.Context, *connect.Request[v1.DeleteTopicRequest]) (*connect.Response[emptypb.Empty], error)
	// rpc CreateConsumerGroup(CreateConsumerGroupRequest) returns (google.protobuf.Empty) {}
	RegisterConsumer(context.Context, *connect.Request[v1.RegisterConsumerRequest]) (*connect.Response[emptypb.Empty], error)
	// Read Procedures
	GetRecords(context.Context, *connect.Request[v1.GetRecordsRequest]) (*connect.Response[v1.GetRecordsResponse], error)
	FetchRecord(context.Context, *connect.Request[v1.FetchRecordRequest]) (*connect.Response[v1.FetchRecordsResponse], error)
	StreamRecords(context.Context, *connect.Request[v1.StreamRecordsRequest], *connect.ServerStream[v1.Record]) error
	ListTopics(context.Context, *connect.Request[v1.ListTopicsRequest]) (*connect.Response[v1.ListTopicsResponse], error)
	Stats(context.Context, *connect.Request[emptypb.Empty]) (*connect.Response[v1.StatsResponse], error)
	GetNodeDetails(context.Context, *connect.Request[emptypb.Empty]) (*connect.Response[v1.GetNodeDetailsResponse], error)
}

// NewKayakServiceHandler builds an HTTP handler from the service implementation. It returns the
// path on which to mount the handler and the handler itself.
//
// By default, handlers support the Connect, gRPC, and gRPC-Web protocols with the binary Protobuf
// and JSON codecs. They also support gzip compression.
func NewKayakServiceHandler(svc KayakServiceHandler, opts ...connect.HandlerOption) (string, http.Handler) {
	kayakServiceApplyHandler := connect.NewUnaryHandler(
		KayakServiceApplyProcedure,
		svc.Apply,
		opts...,
	)
	kayakServicePutRecordsHandler := connect.NewUnaryHandler(
		KayakServicePutRecordsProcedure,
		svc.PutRecords,
		opts...,
	)
	kayakServiceCommitRecordHandler := connect.NewUnaryHandler(
		KayakServiceCommitRecordProcedure,
		svc.CommitRecord,
		opts...,
	)
	kayakServiceCreateTopicHandler := connect.NewUnaryHandler(
		KayakServiceCreateTopicProcedure,
		svc.CreateTopic,
		opts...,
	)
	kayakServiceDeleteTopicHandler := connect.NewUnaryHandler(
		KayakServiceDeleteTopicProcedure,
		svc.DeleteTopic,
		opts...,
	)
	kayakServiceRegisterConsumerHandler := connect.NewUnaryHandler(
		KayakServiceRegisterConsumerProcedure,
		svc.RegisterConsumer,
		opts...,
	)
	kayakServiceGetRecordsHandler := connect.NewUnaryHandler(
		KayakServiceGetRecordsProcedure,
		svc.GetRecords,
		opts...,
	)
	kayakServiceFetchRecordHandler := connect.NewUnaryHandler(
		KayakServiceFetchRecordProcedure,
		svc.FetchRecord,
		opts...,
	)
	kayakServiceStreamRecordsHandler := connect.NewServerStreamHandler(
		KayakServiceStreamRecordsProcedure,
		svc.StreamRecords,
		opts...,
	)
	kayakServiceListTopicsHandler := connect.NewUnaryHandler(
		KayakServiceListTopicsProcedure,
		svc.ListTopics,
		opts...,
	)
	kayakServiceStatsHandler := connect.NewUnaryHandler(
		KayakServiceStatsProcedure,
		svc.Stats,
		opts...,
	)
	kayakServiceGetNodeDetailsHandler := connect.NewUnaryHandler(
		KayakServiceGetNodeDetailsProcedure,
		svc.GetNodeDetails,
		opts...,
	)
	return "/kayak.v1.KayakService/", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case KayakServiceApplyProcedure:
			kayakServiceApplyHandler.ServeHTTP(w, r)
		case KayakServicePutRecordsProcedure:
			kayakServicePutRecordsHandler.ServeHTTP(w, r)
		case KayakServiceCommitRecordProcedure:
			kayakServiceCommitRecordHandler.ServeHTTP(w, r)
		case KayakServiceCreateTopicProcedure:
			kayakServiceCreateTopicHandler.ServeHTTP(w, r)
		case KayakServiceDeleteTopicProcedure:
			kayakServiceDeleteTopicHandler.ServeHTTP(w, r)
		case KayakServiceRegisterConsumerProcedure:
			kayakServiceRegisterConsumerHandler.ServeHTTP(w, r)
		case KayakServiceGetRecordsProcedure:
			kayakServiceGetRecordsHandler.ServeHTTP(w, r)
		case KayakServiceFetchRecordProcedure:
			kayakServiceFetchRecordHandler.ServeHTTP(w, r)
		case KayakServiceStreamRecordsProcedure:
			kayakServiceStreamRecordsHandler.ServeHTTP(w, r)
		case KayakServiceListTopicsProcedure:
			kayakServiceListTopicsHandler.ServeHTTP(w, r)
		case KayakServiceStatsProcedure:
			kayakServiceStatsHandler.ServeHTTP(w, r)
		case KayakServiceGetNodeDetailsProcedure:
			kayakServiceGetNodeDetailsHandler.ServeHTTP(w, r)
		default:
			http.NotFound(w, r)
		}
	})
}

// UnimplementedKayakServiceHandler returns CodeUnimplemented from all methods.
type UnimplementedKayakServiceHandler struct{}

func (UnimplementedKayakServiceHandler) Apply(context.Context, *connect.Request[v1.Command]) (*connect.Response[v1.ApplyResponse], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, errors.New("kayak.v1.KayakService.Apply is not implemented"))
}

func (UnimplementedKayakServiceHandler) PutRecords(context.Context, *connect.Request[v1.PutRecordsRequest]) (*connect.Response[emptypb.Empty], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, errors.New("kayak.v1.KayakService.PutRecords is not implemented"))
}

func (UnimplementedKayakServiceHandler) CommitRecord(context.Context, *connect.Request[v1.CommitRecordRequest]) (*connect.Response[emptypb.Empty], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, errors.New("kayak.v1.KayakService.CommitRecord is not implemented"))
}

func (UnimplementedKayakServiceHandler) CreateTopic(context.Context, *connect.Request[v1.CreateTopicRequest]) (*connect.Response[emptypb.Empty], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, errors.New("kayak.v1.KayakService.CreateTopic is not implemented"))
}

func (UnimplementedKayakServiceHandler) DeleteTopic(context.Context, *connect.Request[v1.DeleteTopicRequest]) (*connect.Response[emptypb.Empty], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, errors.New("kayak.v1.KayakService.DeleteTopic is not implemented"))
}

func (UnimplementedKayakServiceHandler) RegisterConsumer(context.Context, *connect.Request[v1.RegisterConsumerRequest]) (*connect.Response[emptypb.Empty], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, errors.New("kayak.v1.KayakService.RegisterConsumer is not implemented"))
}

func (UnimplementedKayakServiceHandler) GetRecords(context.Context, *connect.Request[v1.GetRecordsRequest]) (*connect.Response[v1.GetRecordsResponse], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, errors.New("kayak.v1.KayakService.GetRecords is not implemented"))
}

func (UnimplementedKayakServiceHandler) FetchRecord(context.Context, *connect.Request[v1.FetchRecordRequest]) (*connect.Response[v1.FetchRecordsResponse], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, errors.New("kayak.v1.KayakService.FetchRecord is not implemented"))
}

func (UnimplementedKayakServiceHandler) StreamRecords(context.Context, *connect.Request[v1.StreamRecordsRequest], *connect.ServerStream[v1.Record]) error {
	return connect.NewError(connect.CodeUnimplemented, errors.New("kayak.v1.KayakService.StreamRecords is not implemented"))
}

func (UnimplementedKayakServiceHandler) ListTopics(context.Context, *connect.Request[v1.ListTopicsRequest]) (*connect.Response[v1.ListTopicsResponse], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, errors.New("kayak.v1.KayakService.ListTopics is not implemented"))
}

func (UnimplementedKayakServiceHandler) Stats(context.Context, *connect.Request[emptypb.Empty]) (*connect.Response[v1.StatsResponse], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, errors.New("kayak.v1.KayakService.Stats is not implemented"))
}

func (UnimplementedKayakServiceHandler) GetNodeDetails(context.Context, *connect.Request[emptypb.Empty]) (*connect.Response[v1.GetNodeDetailsResponse], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, errors.New("kayak.v1.KayakService.GetNodeDetails is not implemented"))
}
