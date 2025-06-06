package raft

import (
	"time"

	"github.com/hashicorp/raft"
)

type RaftInterface interface {
	Apply(cmd []byte, timeout time.Duration) raft.ApplyFuture
	State() raft.RaftState
	Leader() raft.ServerAddress
	AddVoter(id raft.ServerID, address raft.ServerAddress, prevIndex uint64, timeout time.Duration) raft.IndexFuture
	GetConfiguration() raft.ConfigurationFuture
	LeaderWithID() (raft.ServerAddress, raft.ServerID)
	LastContact() time.Time
	Stats() map[string]string
}

type TestFuture struct {
	Err       error
	ResponseI any
}

func (t *TestFuture) Error() error {
	return t.Err
}
func (t *TestFuture) Index() uint64 {
	return 0
}
func (t *TestFuture) Response() any {
	return t.ResponseI
}
