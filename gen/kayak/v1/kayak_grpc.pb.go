// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.2.0
// - protoc             (unknown)
// source: proto/kayak/v1/kayak.proto

package kayakv1

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
	emptypb "google.golang.org/protobuf/types/known/emptypb"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.32.0 or later.
const _ = grpc.SupportPackageIsVersion7

// KayakServiceClient is the client API for KayakService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type KayakServiceClient interface {
	// Apply applies a command on the primary node.
	Apply(ctx context.Context, in *Command, opts ...grpc.CallOption) (*ApplyResponse, error)
	// PutRecords adds records to the stream
	PutRecords(ctx context.Context, in *PutRecordsRequest, opts ...grpc.CallOption) (*emptypb.Empty, error)
	// Commits Consumer position for a topic/group
	CommitRecord(ctx context.Context, in *CommitRecordRequest, opts ...grpc.CallOption) (*emptypb.Empty, error)
	// Creates new topic on server
	CreateTopic(ctx context.Context, in *CreateTopicRequest, opts ...grpc.CallOption) (*emptypb.Empty, error)
	// Deletes Topic across server - permantly or via archive
	DeleteTopic(ctx context.Context, in *DeleteTopicRequest, opts ...grpc.CallOption) (*emptypb.Empty, error)
	// rpc CreateConsumerGroup(CreateConsumerGroupRequest) returns (google.protobuf.Empty) {}
	RegisterConsumer(ctx context.Context, in *RegisterConsumerRequest, opts ...grpc.CallOption) (*emptypb.Empty, error)
	// Read Procedures
	GetRecords(ctx context.Context, in *GetRecordsRequest, opts ...grpc.CallOption) (*GetRecordsResponse, error)
	FetchRecord(ctx context.Context, in *FetchRecordRequest, opts ...grpc.CallOption) (*FetchRecordsResponse, error)
	StreamRecords(ctx context.Context, in *StreamRecordsRequest, opts ...grpc.CallOption) (KayakService_StreamRecordsClient, error)
	ListTopics(ctx context.Context, in *ListTopicsRequest, opts ...grpc.CallOption) (*ListTopicsResponse, error)
	Stats(ctx context.Context, in *emptypb.Empty, opts ...grpc.CallOption) (*StatsResponse, error)
	GetNodeDetails(ctx context.Context, in *emptypb.Empty, opts ...grpc.CallOption) (*GetNodeDetailsResponse, error)
}

type kayakServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewKayakServiceClient(cc grpc.ClientConnInterface) KayakServiceClient {
	return &kayakServiceClient{cc}
}

func (c *kayakServiceClient) Apply(ctx context.Context, in *Command, opts ...grpc.CallOption) (*ApplyResponse, error) {
	out := new(ApplyResponse)
	err := c.cc.Invoke(ctx, "/kayak.v1.KayakService/Apply", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *kayakServiceClient) PutRecords(ctx context.Context, in *PutRecordsRequest, opts ...grpc.CallOption) (*emptypb.Empty, error) {
	out := new(emptypb.Empty)
	err := c.cc.Invoke(ctx, "/kayak.v1.KayakService/PutRecords", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *kayakServiceClient) CommitRecord(ctx context.Context, in *CommitRecordRequest, opts ...grpc.CallOption) (*emptypb.Empty, error) {
	out := new(emptypb.Empty)
	err := c.cc.Invoke(ctx, "/kayak.v1.KayakService/CommitRecord", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *kayakServiceClient) CreateTopic(ctx context.Context, in *CreateTopicRequest, opts ...grpc.CallOption) (*emptypb.Empty, error) {
	out := new(emptypb.Empty)
	err := c.cc.Invoke(ctx, "/kayak.v1.KayakService/CreateTopic", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *kayakServiceClient) DeleteTopic(ctx context.Context, in *DeleteTopicRequest, opts ...grpc.CallOption) (*emptypb.Empty, error) {
	out := new(emptypb.Empty)
	err := c.cc.Invoke(ctx, "/kayak.v1.KayakService/DeleteTopic", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *kayakServiceClient) RegisterConsumer(ctx context.Context, in *RegisterConsumerRequest, opts ...grpc.CallOption) (*emptypb.Empty, error) {
	out := new(emptypb.Empty)
	err := c.cc.Invoke(ctx, "/kayak.v1.KayakService/RegisterConsumer", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *kayakServiceClient) GetRecords(ctx context.Context, in *GetRecordsRequest, opts ...grpc.CallOption) (*GetRecordsResponse, error) {
	out := new(GetRecordsResponse)
	err := c.cc.Invoke(ctx, "/kayak.v1.KayakService/GetRecords", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *kayakServiceClient) FetchRecord(ctx context.Context, in *FetchRecordRequest, opts ...grpc.CallOption) (*FetchRecordsResponse, error) {
	out := new(FetchRecordsResponse)
	err := c.cc.Invoke(ctx, "/kayak.v1.KayakService/FetchRecord", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *kayakServiceClient) StreamRecords(ctx context.Context, in *StreamRecordsRequest, opts ...grpc.CallOption) (KayakService_StreamRecordsClient, error) {
	stream, err := c.cc.NewStream(ctx, &KayakService_ServiceDesc.Streams[0], "/kayak.v1.KayakService/StreamRecords", opts...)
	if err != nil {
		return nil, err
	}
	x := &kayakServiceStreamRecordsClient{stream}
	if err := x.ClientStream.SendMsg(in); err != nil {
		return nil, err
	}
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	return x, nil
}

type KayakService_StreamRecordsClient interface {
	Recv() (*Record, error)
	grpc.ClientStream
}

type kayakServiceStreamRecordsClient struct {
	grpc.ClientStream
}

func (x *kayakServiceStreamRecordsClient) Recv() (*Record, error) {
	m := new(Record)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func (c *kayakServiceClient) ListTopics(ctx context.Context, in *ListTopicsRequest, opts ...grpc.CallOption) (*ListTopicsResponse, error) {
	out := new(ListTopicsResponse)
	err := c.cc.Invoke(ctx, "/kayak.v1.KayakService/ListTopics", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *kayakServiceClient) Stats(ctx context.Context, in *emptypb.Empty, opts ...grpc.CallOption) (*StatsResponse, error) {
	out := new(StatsResponse)
	err := c.cc.Invoke(ctx, "/kayak.v1.KayakService/Stats", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *kayakServiceClient) GetNodeDetails(ctx context.Context, in *emptypb.Empty, opts ...grpc.CallOption) (*GetNodeDetailsResponse, error) {
	out := new(GetNodeDetailsResponse)
	err := c.cc.Invoke(ctx, "/kayak.v1.KayakService/GetNodeDetails", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// KayakServiceServer is the server API for KayakService service.
// All implementations should embed UnimplementedKayakServiceServer
// for forward compatibility
type KayakServiceServer interface {
	// Apply applies a command on the primary node.
	Apply(context.Context, *Command) (*ApplyResponse, error)
	// PutRecords adds records to the stream
	PutRecords(context.Context, *PutRecordsRequest) (*emptypb.Empty, error)
	// Commits Consumer position for a topic/group
	CommitRecord(context.Context, *CommitRecordRequest) (*emptypb.Empty, error)
	// Creates new topic on server
	CreateTopic(context.Context, *CreateTopicRequest) (*emptypb.Empty, error)
	// Deletes Topic across server - permantly or via archive
	DeleteTopic(context.Context, *DeleteTopicRequest) (*emptypb.Empty, error)
	// rpc CreateConsumerGroup(CreateConsumerGroupRequest) returns (google.protobuf.Empty) {}
	RegisterConsumer(context.Context, *RegisterConsumerRequest) (*emptypb.Empty, error)
	// Read Procedures
	GetRecords(context.Context, *GetRecordsRequest) (*GetRecordsResponse, error)
	FetchRecord(context.Context, *FetchRecordRequest) (*FetchRecordsResponse, error)
	StreamRecords(*StreamRecordsRequest, KayakService_StreamRecordsServer) error
	ListTopics(context.Context, *ListTopicsRequest) (*ListTopicsResponse, error)
	Stats(context.Context, *emptypb.Empty) (*StatsResponse, error)
	GetNodeDetails(context.Context, *emptypb.Empty) (*GetNodeDetailsResponse, error)
}

// UnimplementedKayakServiceServer should be embedded to have forward compatible implementations.
type UnimplementedKayakServiceServer struct {
}

func (UnimplementedKayakServiceServer) Apply(context.Context, *Command) (*ApplyResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Apply not implemented")
}
func (UnimplementedKayakServiceServer) PutRecords(context.Context, *PutRecordsRequest) (*emptypb.Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method PutRecords not implemented")
}
func (UnimplementedKayakServiceServer) CommitRecord(context.Context, *CommitRecordRequest) (*emptypb.Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CommitRecord not implemented")
}
func (UnimplementedKayakServiceServer) CreateTopic(context.Context, *CreateTopicRequest) (*emptypb.Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CreateTopic not implemented")
}
func (UnimplementedKayakServiceServer) DeleteTopic(context.Context, *DeleteTopicRequest) (*emptypb.Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method DeleteTopic not implemented")
}
func (UnimplementedKayakServiceServer) RegisterConsumer(context.Context, *RegisterConsumerRequest) (*emptypb.Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method RegisterConsumer not implemented")
}
func (UnimplementedKayakServiceServer) GetRecords(context.Context, *GetRecordsRequest) (*GetRecordsResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetRecords not implemented")
}
func (UnimplementedKayakServiceServer) FetchRecord(context.Context, *FetchRecordRequest) (*FetchRecordsResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method FetchRecord not implemented")
}
func (UnimplementedKayakServiceServer) StreamRecords(*StreamRecordsRequest, KayakService_StreamRecordsServer) error {
	return status.Errorf(codes.Unimplemented, "method StreamRecords not implemented")
}
func (UnimplementedKayakServiceServer) ListTopics(context.Context, *ListTopicsRequest) (*ListTopicsResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ListTopics not implemented")
}
func (UnimplementedKayakServiceServer) Stats(context.Context, *emptypb.Empty) (*StatsResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Stats not implemented")
}
func (UnimplementedKayakServiceServer) GetNodeDetails(context.Context, *emptypb.Empty) (*GetNodeDetailsResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetNodeDetails not implemented")
}

// UnsafeKayakServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to KayakServiceServer will
// result in compilation errors.
type UnsafeKayakServiceServer interface {
	mustEmbedUnimplementedKayakServiceServer()
}

func RegisterKayakServiceServer(s grpc.ServiceRegistrar, srv KayakServiceServer) {
	s.RegisterService(&KayakService_ServiceDesc, srv)
}

func _KayakService_Apply_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Command)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(KayakServiceServer).Apply(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/kayak.v1.KayakService/Apply",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(KayakServiceServer).Apply(ctx, req.(*Command))
	}
	return interceptor(ctx, in, info, handler)
}

func _KayakService_PutRecords_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(PutRecordsRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(KayakServiceServer).PutRecords(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/kayak.v1.KayakService/PutRecords",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(KayakServiceServer).PutRecords(ctx, req.(*PutRecordsRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _KayakService_CommitRecord_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CommitRecordRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(KayakServiceServer).CommitRecord(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/kayak.v1.KayakService/CommitRecord",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(KayakServiceServer).CommitRecord(ctx, req.(*CommitRecordRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _KayakService_CreateTopic_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CreateTopicRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(KayakServiceServer).CreateTopic(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/kayak.v1.KayakService/CreateTopic",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(KayakServiceServer).CreateTopic(ctx, req.(*CreateTopicRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _KayakService_DeleteTopic_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(DeleteTopicRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(KayakServiceServer).DeleteTopic(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/kayak.v1.KayakService/DeleteTopic",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(KayakServiceServer).DeleteTopic(ctx, req.(*DeleteTopicRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _KayakService_RegisterConsumer_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(RegisterConsumerRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(KayakServiceServer).RegisterConsumer(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/kayak.v1.KayakService/RegisterConsumer",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(KayakServiceServer).RegisterConsumer(ctx, req.(*RegisterConsumerRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _KayakService_GetRecords_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetRecordsRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(KayakServiceServer).GetRecords(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/kayak.v1.KayakService/GetRecords",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(KayakServiceServer).GetRecords(ctx, req.(*GetRecordsRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _KayakService_FetchRecord_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(FetchRecordRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(KayakServiceServer).FetchRecord(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/kayak.v1.KayakService/FetchRecord",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(KayakServiceServer).FetchRecord(ctx, req.(*FetchRecordRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _KayakService_StreamRecords_Handler(srv interface{}, stream grpc.ServerStream) error {
	m := new(StreamRecordsRequest)
	if err := stream.RecvMsg(m); err != nil {
		return err
	}
	return srv.(KayakServiceServer).StreamRecords(m, &kayakServiceStreamRecordsServer{stream})
}

type KayakService_StreamRecordsServer interface {
	Send(*Record) error
	grpc.ServerStream
}

type kayakServiceStreamRecordsServer struct {
	grpc.ServerStream
}

func (x *kayakServiceStreamRecordsServer) Send(m *Record) error {
	return x.ServerStream.SendMsg(m)
}

func _KayakService_ListTopics_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ListTopicsRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(KayakServiceServer).ListTopics(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/kayak.v1.KayakService/ListTopics",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(KayakServiceServer).ListTopics(ctx, req.(*ListTopicsRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _KayakService_Stats_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(emptypb.Empty)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(KayakServiceServer).Stats(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/kayak.v1.KayakService/Stats",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(KayakServiceServer).Stats(ctx, req.(*emptypb.Empty))
	}
	return interceptor(ctx, in, info, handler)
}

func _KayakService_GetNodeDetails_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(emptypb.Empty)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(KayakServiceServer).GetNodeDetails(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/kayak.v1.KayakService/GetNodeDetails",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(KayakServiceServer).GetNodeDetails(ctx, req.(*emptypb.Empty))
	}
	return interceptor(ctx, in, info, handler)
}

// KayakService_ServiceDesc is the grpc.ServiceDesc for KayakService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var KayakService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "kayak.v1.KayakService",
	HandlerType: (*KayakServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Apply",
			Handler:    _KayakService_Apply_Handler,
		},
		{
			MethodName: "PutRecords",
			Handler:    _KayakService_PutRecords_Handler,
		},
		{
			MethodName: "CommitRecord",
			Handler:    _KayakService_CommitRecord_Handler,
		},
		{
			MethodName: "CreateTopic",
			Handler:    _KayakService_CreateTopic_Handler,
		},
		{
			MethodName: "DeleteTopic",
			Handler:    _KayakService_DeleteTopic_Handler,
		},
		{
			MethodName: "RegisterConsumer",
			Handler:    _KayakService_RegisterConsumer_Handler,
		},
		{
			MethodName: "GetRecords",
			Handler:    _KayakService_GetRecords_Handler,
		},
		{
			MethodName: "FetchRecord",
			Handler:    _KayakService_FetchRecord_Handler,
		},
		{
			MethodName: "ListTopics",
			Handler:    _KayakService_ListTopics_Handler,
		},
		{
			MethodName: "Stats",
			Handler:    _KayakService_Stats_Handler,
		},
		{
			MethodName: "GetNodeDetails",
			Handler:    _KayakService_GetNodeDetails_Handler,
		},
	},
	Streams: []grpc.StreamDesc{
		{
			StreamName:    "StreamRecords",
			Handler:       _KayakService_StreamRecords_Handler,
			ServerStreams: true,
		},
	},
	Metadata: "proto/kayak/v1/kayak.proto",
}
