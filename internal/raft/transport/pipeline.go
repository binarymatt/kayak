package transport

import (
	"errors"
	"io"
	"log/slog"
	"sync"
	"time"

	"connectrpc.com/connect"
	"github.com/hashicorp/raft"

	transportv1 "github.com/binarymatt/kayak/gen/transport/v1"
)

var _ raft.AppendPipeline = (*appendPipeline)(nil)

type appendPipeline struct {
	stream        *connect.BidiStreamForClient[transportv1.AppendEntriesRequest, transportv1.AppendEntriesResponse]
	inflightChMtx sync.Mutex
	inflightCh    chan *appendFuture
	cancel        func()
	doneCh        chan raft.AppendFuture
}

func (ap *appendPipeline) AppendEntries(req *raft.AppendEntriesRequest, res *raft.AppendEntriesResponse) (raft.AppendFuture, error) {
	af := &appendFuture{
		start:   time.Now(),
		request: req,
		done:    make(chan struct{}),
	}
	if err := ap.stream.Send(encodeAppendEntriesRequest(req)); err != nil {
		slog.Debug("appendPipeline: AppendEntries failure", "error", err)
		if errors.Is(err, io.EOF) {
			_, err := ap.stream.Receive()
			slog.Error("failure during send", "error", err)
			return nil, err
		}
		return nil, err
	}
	ap.inflightChMtx.Lock()
	ap.inflightCh <- af
	ap.inflightChMtx.Unlock()
	return af, nil
}

func (ap *appendPipeline) Close() error {
	ap.cancel()
	ap.inflightChMtx.Lock()
	close(ap.inflightCh)
	ap.inflightChMtx.Unlock()
	return nil
}

func (ap *appendPipeline) Consumer() <-chan raft.AppendFuture {
	return ap.doneCh
}

func (ap *appendPipeline) receiver() {
	slog.Info("goroutine recieiver running")
	for af := range ap.inflightCh {
		msg, err := ap.stream.Receive()
		if err != nil {
			slog.Error("error during receive", "error", err)
			if !errors.Is(err, io.EOF) {
				af.err = err
			}
		} else {
			af.response = *decodeAppendEntriesResponse(msg)
		}
		close(af.done)
		ap.doneCh <- af
	}
	slog.Warn("done with append pipeline goroutine")
}

var _ raft.AppendFuture = (*appendFuture)(nil)

type appendFuture struct {
	start    time.Time
	request  *raft.AppendEntriesRequest
	response raft.AppendEntriesResponse
	err      error
	done     chan struct{}
}

func (af *appendFuture) Start() time.Time {
	return af.start
}
func (af *appendFuture) Error() error {
	<-af.done
	return af.err
}
func (af *appendFuture) Request() *raft.AppendEntriesRequest {
	return af.request
}
func (af *appendFuture) Response() *raft.AppendEntriesResponse {
	return &af.response
}
