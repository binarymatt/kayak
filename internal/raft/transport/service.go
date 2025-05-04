package transport

import (
	"context"
	"errors"
	"io"
	"log/slog"

	"connectrpc.com/connect"
	"github.com/hashicorp/raft"

	v1 "github.com/binarymatt/kayak/gen/transport/v1"
	"github.com/binarymatt/kayak/gen/transport/v1/transportv1connect"
)

var _ transportv1connect.RaftTransportHandler = (*transportService)(nil)

type transportService struct {
	rpcChan      chan raft.RPC
	shutdownChan chan struct{}
}

func (ts *transportService) handleRPC(command interface{}, data io.Reader) (interface{}, error) {
	ch := make(chan raft.RPCResponse, 1)
	rpc := raft.RPC{
		Command:  command,
		RespChan: ch,
		Reader:   data,
	}
	select {
	case ts.rpcChan <- rpc:
	case <-ts.shutdownChan:
		return nil, raft.ErrTransportShutdown
	}
	select {
	case resp := <-ch:
		if resp.Error != nil {
			return nil, resp.Error
		}
		return resp.Response, nil
	case <-ts.shutdownChan:
		return nil, raft.ErrTransportShutdown
	}
}
func (ts *transportService) AppendEntries(ctx context.Context, req *connect.Request[v1.AppendEntriesRequest]) (*connect.Response[v1.AppendEntriesResponse], error) {
	res, err := ts.handleRPC(decodeAppendEntriesRequest(req.Msg), nil)
	if err != nil {
		slog.Error("could not handle rpc", "error", err)
		return nil, err
	}
	resp := encodeAppendEntriesResponse(res.(*raft.AppendEntriesResponse))
	return connect.NewResponse(resp), nil
}

func (ts *transportService) AppendEntriesPipeline(ctx context.Context, req *connect.BidiStream[v1.AppendEntriesRequest, v1.AppendEntriesResponse]) error {
	for {
		msg, err := req.Receive()
		if err != nil {
			if errors.Is(err, io.EOF) {
				return nil
			}
			slog.Error("error from receive in transportService", "error", err)
			return err
		}
		resp, err := ts.handleRPC(decodeAppendEntriesRequest(msg), nil)
		if err != nil {
			slog.Error("could not handleRPC", "error", err)
			return err
		}
		if err := req.Send(encodeAppendEntriesResponse(resp.(*raft.AppendEntriesResponse))); err != nil {
			slog.Error("could not send response", "error", err)
			return err
		}

	}
}

func (ts *transportService) InstallSnapshot(ctx context.Context, req *connect.ClientStream[v1.InstallSnapshotRequest]) (*connect.Response[v1.InstallSnapshotResponse], error) {
	ok := req.Receive()
	if !ok {
		return nil, errors.New("oops")
	}
	isr := req.Msg()

	resp, err := ts.handleRPC(decodeInstallSnapshotRequest(isr), &snapshotReader{req, isr.GetData()})
	if err != nil {
		return nil, err
	}
	return connect.NewResponse(encodeInstallSnapshotResponse(resp.(*raft.InstallSnapshotResponse))), nil
}

type snapshotReader struct {
	req *connect.ClientStream[v1.InstallSnapshotRequest]
	buf []byte
}

func (s *snapshotReader) Read(b []byte) (int, error) {
	if len(s.buf) > 0 {
		n := copy(b, s.buf)
		s.buf = s.buf[n:]
		return n, nil
	}
	if ok := s.req.Receive(); !ok {
		return 0, errors.New("issue with receive")
	}

	n := copy(b, s.req.Msg().GetData())
	if n < len(s.req.Msg().GetData()) {
		s.buf = s.req.Msg().GetData()[n:]
	}
	return n, nil
}

func (ts *transportService) RequestPreVote(ctx context.Context, req *connect.Request[v1.RequestPreVoteRequest]) (*connect.Response[v1.RequestPreVoteResponse], error) {
	resp, err := ts.handleRPC(decodeRequestPreVoteRequest(req.Msg), nil)
	if err != nil {
		return nil, err
	}
	res := encodeRequestPreVoteResponse(resp.(*raft.RequestPreVoteResponse))
	return connect.NewResponse(res), nil
}
func (ts *transportService) RequestVote(ctx context.Context, req *connect.Request[v1.RequestVoteRequest]) (*connect.Response[v1.RequestVoteResponse], error) {
	resp, err := ts.handleRPC(decodeRequestVoteRequest(req.Msg), nil)
	if err != nil {
		return nil, err
	}
	res := encodeRequestVoteResponse(resp.(*raft.RequestVoteResponse))
	return connect.NewResponse(res), nil
}
func (ts *transportService) TimeoutNow(ctx context.Context, req *connect.Request[v1.TimeoutNowRequest]) (*connect.Response[v1.TimeoutNowResponse], error) {
	resp, err := ts.handleRPC(decodeTimeoutNowRequest(req.Msg), nil)
	if err != nil {
		return nil, err
	}
	res := encodeTimeoutNowResponse(resp.(*raft.TimeoutNowResponse))
	return connect.NewResponse(res), nil
}
