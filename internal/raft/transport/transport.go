package transport

import (
	"context"
	"crypto/tls"
	"fmt"
	"io"
	"log/slog"
	"net"
	"net/http"
	"sync"
	"time"

	"connectrpc.com/connect"
	"github.com/hashicorp/raft"
	"golang.org/x/net/http2"

	transportv1 "github.com/binarymatt/kayak/gen/transport/v1"
	"github.com/binarymatt/kayak/gen/transport/v1/transportv1connect"
)

var _ raft.Transport = (*raftTransport)(nil)
var _ raft.WithClose = (*raftTransport)(nil)
var _ raft.WithPeers = (*raftTransport)(nil)
var _ raft.WithPreVote = (*raftTransport)(nil)

type raftTransport struct {
	clientOpts   []connect.ClientOption
	rpcChan      chan raft.RPC
	localAddress raft.ServerAddress

	heartbeatFunc    func(raft.RPC)
	heartbeatFuncMtx sync.Mutex
	closer           func()

	heartbeatTimeout time.Duration
}

func newClient() *http.Client {
	return &http.Client{
		Timeout: time.Second * 10,
		Transport: &http2.Transport{
			AllowHTTP: true,
			DialTLS: func(network, addr string, _ *tls.Config) (net.Conn, error) {
				// If you're also using this client for non-h2c traffic, you may want
				// to delegate to tls.Dial if the network isn't TCP or the addr isn't
				// in an allowlist.
				return net.Dial(network, addr)
			},
		},
	}
}
func (rt *raftTransport) getPeer(target raft.ServerAddress) transportv1connect.RaftTransportClient {
	baseURL := fmt.Sprintf("http://%s", target)

	client := transportv1connect.NewRaftTransportClient(newClient(), baseURL, rt.clientOpts...)
	return client
}

func (rt *raftTransport) AppendEntries(id raft.ServerID, target raft.ServerAddress, args *raft.AppendEntriesRequest, resp *raft.AppendEntriesResponse) error {
	c := rt.getPeer(target)

	ctx := context.TODO()
	// add check for heartbeat
	if rt.heartbeatTimeout > 0 && isHeartbeat(args) {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, rt.heartbeatTimeout)
		defer cancel()
	}
	ret, err := c.AppendEntries(ctx, connect.NewRequest(encodeAppendEntriesRequest(args)))
	if err != nil {
		return err
	}
	*resp = *decodeAppendEntriesResponse(ret.Msg)
	return nil
}

func (rt *raftTransport) AppendEntriesPipeline(id raft.ServerID, target raft.ServerAddress) (raft.AppendPipeline, error) {
	slog.Info("starting transport append entries pipeline")
	c := rt.getPeer(target)

	ctx, cancel := context.WithCancel(context.TODO())
	stream := c.AppendEntriesPipeline(ctx)
	pipeline := &appendPipeline{
		stream:     stream,
		cancel:     cancel,
		inflightCh: make(chan *appendFuture, 20),
		doneCh:     make(chan raft.AppendFuture, 20),
	}
	go pipeline.receiver()
	return pipeline, nil
}

func (rt *raftTransport) Consumer() <-chan raft.RPC {
	return rt.rpcChan
}

func (rt *raftTransport) DecodePeer(p []byte) raft.ServerAddress {
	return raft.ServerAddress(p)
}
func (rt *raftTransport) EncodePeer(id raft.ServerID, addr raft.ServerAddress) []byte {

	return []byte(addr)
}

func (rt *raftTransport) InstallSnapshot(id raft.ServerID, target raft.ServerAddress, req *raft.InstallSnapshotRequest, resp *raft.InstallSnapshotResponse, data io.Reader) error {
	c := rt.getPeer(target)
	ctx := context.TODO()
	stream := c.InstallSnapshot(ctx)
	if err := stream.Send(encodeInstallSnapshotRequest(req)); err != nil {
		return err
	}
	var buf [16384]byte
	for {
		n, err := data.Read(buf[:])
		if err == io.EOF || (err == nil && n == 0) {
			break
		}
		if err != nil {
			return err
		}
		if err := stream.Send(&transportv1.InstallSnapshotRequest{
			Data: buf[:n],
		}); err != nil {
			return err
		}
	}
	res, err := stream.CloseAndReceive()
	if err != nil {
		return err
	}
	*resp = *decodeInstallSnapshotResponse(res.Msg)
	return nil
}

func (rt *raftTransport) LocalAddr() raft.ServerAddress {
	return rt.localAddress
}

func (rt *raftTransport) RequestVote(id raft.ServerID, target raft.ServerAddress, args *raft.RequestVoteRequest, resp *raft.RequestVoteResponse) error {
	c := rt.getPeer(target)
	ret, err := c.RequestVote(context.TODO(), connect.NewRequest(encodeRequestVoteRequest(args)))
	if err != nil {
		return err
	}
	*resp = *decodeRequestVoteResponse(ret.Msg)
	return nil
}

func (rt *raftTransport) SetHeartbeatHandler(cb func(rpc raft.RPC)) {
	rt.heartbeatFuncMtx.Lock()
	rt.heartbeatFunc = cb
	rt.heartbeatFuncMtx.Unlock()
}

func (rt *raftTransport) TimeoutNow(id raft.ServerID, target raft.ServerAddress, args *raft.TimeoutNowRequest, resp *raft.TimeoutNowResponse) error {
	client := rt.getPeer(target)
	res, err := client.TimeoutNow(context.TODO(), connect.NewRequest(encodeTimeoutNowRequest(args)))
	if err != nil {
		return err
	}
	*resp = *decodeTimeoutNowResponse(res.Msg)
	return nil
}

func (rt *raftTransport) Close() error {
	rt.closer()
	return nil
}

func (rt *raftTransport) Connect(target raft.ServerAddress, t raft.Transport) {
	rt.getPeer(target)
}

func (rt *raftTransport) Disconnect(target raft.ServerAddress) {}
func (rt *raftTransport) DisconnectAll()                       {}

func (rt *raftTransport) RequestPreVote(id raft.ServerID, target raft.ServerAddress, args *raft.RequestPreVoteRequest, resp *raft.RequestPreVoteResponse) error {
	c := rt.getPeer(target)
	ret, err := c.RequestPreVote(context.TODO(), connect.NewRequest(encodeRequestPreVoteRequest(args)))
	if err != nil {
		return err
	}
	*resp = *decodeRequestPreVoteResponse(ret.Msg)
	return nil
}

func isHeartbeat(command interface{}) bool {
	req, ok := command.(*raft.AppendEntriesRequest)
	if !ok {
		return false
	}
	return req.Term != 0 && len(req.Addr) != 0 && req.PrevLogEntry == 0 && req.PrevLogTerm == 0 && len(req.Entries) == 0 && req.LeaderCommitIndex == 0
}

func New(localAddress raft.ServerAddress, clientOptions []connect.ClientOption) (*raftTransport, *transportService) {
	rpcCh := make(chan raft.RPC)
	shutdownCh := make(chan struct{})

	tr := &raftTransport{
		localAddress: localAddress,
		clientOpts:   clientOptions,
		rpcChan:      rpcCh,
		closer: func() {
			close(shutdownCh)
		},
	}
	ts := &transportService{
		rpcChan:      rpcCh,
		shutdownChan: shutdownCh,
	}
	return tr, ts
}
