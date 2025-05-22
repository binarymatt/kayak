package store

import (
	"errors"
	"io"
	"time"

	"github.com/hashicorp/raft"
	"google.golang.org/protobuf/proto"

	kayakv1 "github.com/binarymatt/kayak/gen/kayak/v1"
	"github.com/binarymatt/kayak/internal/fsm"
)

var _ raft.FSM = (*store)(nil)

type ApplyResponse struct {
	Error error
}

func (s *store) Apply(l *raft.Log) any {
	var c kayakv1.RaftCommand
	response := &ApplyResponse{}
	if err := proto.Unmarshal(l.Data, &c); err != nil {
		response.Error = err
		return response
	}
	switch p := c.Payload.(type) {
	case *kayakv1.RaftCommand_PutStream:
		response.Error = s.PutStream(p.PutStream.Stream)
	case *kayakv1.RaftCommand_PutRecords:
		response.Error = s.PutRecords(p.PutRecords.StreamName, p.PutRecords.Records...)
	case *kayakv1.RaftCommand_ExtendLease:
		duration := time.Until(time.UnixMilli(p.ExtendLease.ExpiresMs))
		response.Error = s.ExtendLease(p.ExtendLease.Worker, duration)
	case *kayakv1.RaftCommand_RemoveLease:
		response.Error = s.RemoveLease(p.RemoveLease.Worker)
	case *kayakv1.RaftCommand_CommitGroupPosition:
		response.Error = s.CommitGroupPosition(p.CommitGroupPosition.StreamName, p.CommitGroupPosition.GroupName, p.CommitGroupPosition.Partition, p.CommitGroupPosition.Position)
	case *kayakv1.RaftCommand_DeleteStream:
		response.Error = s.DeleteStream(p.DeleteStream.StreamName)
	default:
		response.Error = errors.New("unknown command")
	}
	return response
}

func (s *store) Snapshot() (raft.FSMSnapshot, error) {
	return fsm.NewFSMSnapshot(s.db), nil
}
func (s *store) Restore(snapshot io.ReadCloser) error {
	return s.db.Load(snapshot, 1)
}
