package kayak

import (
	"context"
	"fmt"
	"log/slog"
	"testing"
	"time"

	"connectrpc.com/connect"
	"github.com/google/go-cmp/cmp"
	"github.com/hashicorp/go-hclog"
	"github.com/hashicorp/raft"
	"github.com/oklog/ulid/v2"
	"github.com/shoenig/test/must"
	"github.com/stretchr/testify/mock"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/testing/protocmp"

	kayakv1 "github.com/binarymatt/kayak/gen/kayak/v1"
	"github.com/binarymatt/kayak/internal/store"
)

type testServiceSuite struct {
	service   *service
	mockStore *store.MockStore
	id        ulid.ULID
}

func setupTest(t *testing.T) *testServiceSuite {
	slog.Warn("setting up test")
	conf := raft.DefaultConfig()
	conf.HeartbeatTimeout = 50 * time.Millisecond
	conf.ElectionTimeout = 50 * time.Millisecond
	conf.LeaderLeaseTimeout = 5 * time.Millisecond
	conf.CommitTimeout = 5 * time.Millisecond
	//conf.LogOutput = io.Discard
	//conf.LogLevel = "ERROR"
	conf.Logger = hclog.NewNullLogger()
	conf.Logger.SetLevel(hclog.Off)

	s := store.NewMockStore(t)
	ts := &testServiceSuite{}
	ts.mockStore = s
	c := raft.MakeClusterCustom(t, &raft.MakeClusterOpts{
		Peers:     1,
		Bootstrap: true,
		Conf:      conf,
		MakeFSMFunc: func() raft.FSM {
			return s
		},
	})
	r := c.Leader()

	ts.id = ulid.Make()
	ts.service = &service{
		raft:   r,
		logger: slog.Default(),
		store:  s,
		idGenerator: func() ulid.ULID {
			return ts.id
		},
	}
	return ts
}
func TestPutRecords(t *testing.T) {
	ts := setupTest(t)
	// t := ts.T()
	ctx := context.Background()
	records := &kayakv1.Record{
		Payload: []byte("test"),
		Id:      []byte("test"),
	}
	expected := []*kayakv1.Record{
		{
			Id:         []byte("test"),
			InternalId: ts.id.String(),
			Payload:    []byte("test"),
			Partition:  0,
		},
	}

	ts.mockStore.EXPECT().GetStream("test_stream").
		Return(&kayakv1.Stream{Name: "test_stream", PartitionCount: 1}, nil).
		Once()
	// ts.mockStore.EXPECT().PutRecords("test_stream", tc.expected).Return(nil).Once()
	cmd := &kayakv1.RaftCommand{
		Payload: &kayakv1.RaftCommand_PutRecords{
			PutRecords: &kayakv1.PutRecords{
				Records:    expected,
				StreamName: "test_stream",
			},
		},
	}
	matcher := func(l *raft.Log) bool {
		var c kayakv1.RaftCommand
		proto.Unmarshal(l.Data, &c) //nolint:errcheck
		diff := cmp.Diff(cmd, &c, protocmp.Transform())
		fmt.Println(diff)
		return diff == ""
	}
	ts.mockStore.EXPECT().Apply(mock.MatchedBy(matcher)).Return(&store.ApplyResponse{})

	req := connect.NewRequest(&kayakv1.PutRecordsRequest{
		StreamName: "test_stream",
		Records:    []*kayakv1.Record{records},
	})
	_, err := ts.service.PutRecords(ctx, req)
	must.NoError(t, err)

}
func TestGetRecords(t *testing.T) {}
