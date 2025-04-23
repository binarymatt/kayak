package kayak

import (
	"context"
	"errors"
	"time"

	"connectrpc.com/connect"
	"github.com/dgraph-io/badger/v4"
	"github.com/oklog/ulid/v2"
	"google.golang.org/protobuf/types/known/emptypb"

	v1 "github.com/binarymatt/kayak/gen/kayak/v1"
	"github.com/binarymatt/kayak/gen/kayak/v1/kayakv1connect"
	"github.com/binarymatt/kayak/internal/store"
)

var _ kayakv1connect.KayakServiceHandler = (*service)(nil)

type service struct {
	// kayakv1connect.UnimplementedKayakServiceHandler
	idGenerator      func() ulid.ULID
	store            store.Store
	workerExpiration time.Duration
}

func (s *service) PutRecords(ctx context.Context, req *connect.Request[v1.PutRecordsRequest]) (*connect.Response[emptypb.Empty], error) {
	stream, err := s.store.GetStream(req.Msg.StreamName)
	if err != nil {
		return nil, err
	}
	for _, r := range req.Msg.Records {
		id := s.idGenerator()
		r.InternalId = id.Bytes()
		if r.Id == nil {
			r.Id = id.Bytes()
		}
		r.Partition = balancer(string(r.Id), stream.PartitionCount)
	}
	err = s.store.PutRecords(req.Msg.StreamName, req.Msg.GetRecords()...)
	return connect.NewResponse(&emptypb.Empty{}), err
}
func (s *service) CommitRecord(ctx context.Context, req *connect.Request[v1.CommitRecordRequest]) (*connect.Response[emptypb.Empty], error) {
	// extend lease of worker
	return nil, nil
}
func (s *service) FetchRecords(ctx context.Context, req *connect.Request[v1.FetchRecordsRequest]) (*connect.Response[v1.FetchRecordsResponse], error) {
	// is worker registered?
	// if so, extend lease
	// if not, return error
	// Get records based on worker info
	return nil, nil
}

func (s *service) GetRecords(ctx context.Context, req *connect.Request[v1.GetRecordsRequest]) (*connect.Response[v1.GetRecordsResponse], error) {
	records, err := s.store.GetRecords(req.Msg.StreamName, req.Msg.Partition, req.Msg.StartId, int(req.Msg.Limit))
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}
	return connect.NewResponse(&v1.GetRecordsResponse{
		Records: records,
	}), nil
}

func (s *service) CreateStream(ctx context.Context, req *connect.Request[v1.CreateStreamRequest]) (*connect.Response[emptypb.Empty], error) {
	// TODO: validate request
	stream := &v1.Stream{
		Name:           req.Msg.Name,
		PartitionCount: req.Msg.PartitionCount,
		Ttl:            req.Msg.Ttl,
	}
	if err := s.store.PutStream(stream); err != nil {
		return nil, err
	}
	return connect.NewResponse(&emptypb.Empty{}), nil
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
	n := time.Now().Add(s.workerExpiration)
	worker.LeaseExpires = n.UnixMilli()
	if err := s.store.ExtendLease(worker, s.workerExpiration); err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}
	resp := &v1.RegisterWorkerResponse{
		Worker: worker,
	}
	return connect.NewResponse(resp), nil
}

func (s *service) DeregisterWorker(ctx context.Context, req *connect.Request[v1.DeregisterWorkerRequest]) (*connect.Response[emptypb.Empty], error) {
	return nil, nil
}

func New() *service {
	return &service{
		idGenerator:      ulid.Make,
		workerExpiration: 60 * time.Second,
	}
}
