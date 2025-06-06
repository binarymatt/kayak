package store

import (
	"testing"

	"github.com/dgraph-io/badger/v4"
	"github.com/hashicorp/raft"
	"github.com/shoenig/test/must"
	"google.golang.org/protobuf/proto"

	kayakv1 "github.com/binarymatt/kayak/gen/kayak/v1"
)

func TestApply_UnknownCommand(t *testing.T) {

	ts := setupTest(t)
	resp := ts.store.Apply(&raft.Log{})
	ar, ok := resp.(*ApplyResponse)
	must.True(t, ok)
	must.Error(t, ar.Error)
	must.ErrorContains(t, ar.Error, "unknown command")
}

func TestApply_BadProto(t *testing.T) {

	ts := setupTest(t)
	resp := ts.store.Apply(&raft.Log{Data: []byte("{}")})
	ar, ok := resp.(*ApplyResponse)
	must.True(t, ok)
	must.Error(t, ar.Error)
	must.ErrorContains(t, ar.Error, "cannot parse invalid wire-format data")
}

func TestApply(t *testing.T) {
	cases := []struct {
		name    string
		command *kayakv1.RaftCommand
		check   func(t *testing.T, db *badger.DB)
	}{
		{
			name: "PutStream",
			command: &kayakv1.RaftCommand{
				Payload: &kayakv1.RaftCommand_PutStream{
					PutStream: &kayakv1.PutStream{
						Stream: &kayakv1.Stream{
							Name:           "test",
							PartitionCount: 1,
						},
					},
				},
			},
			check: func(t *testing.T, db *badger.DB) {
			},
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			ts := setupTest(t)
			raw, err := proto.Marshal(tc.command)
			must.NoError(t, err)
			l := &raft.Log{
				Data: raw,
			}
			resp := ts.store.Apply(l)
			ar, ok := resp.(*ApplyResponse)
			must.True(t, ok)
			must.NoError(t, ar.Error)
			tc.check(t, ts.db)
		})
	}
}
