package kayak

import (
	"context"
	"log/slog"
	"testing"
	"time"

	"connectrpc.com/connect"
	"github.com/coder/quartz"
	"github.com/hashicorp/raft"
	"github.com/oklog/ulid/v2"
	"github.com/shoenig/test/must"
	"google.golang.org/protobuf/proto"

	kayakv1 "github.com/binarymatt/kayak/gen/kayak/v1"
	"github.com/binarymatt/kayak/gen/kayak/v1/kayakv1connect"
	"github.com/binarymatt/kayak/internal/store"
)

type TestFuture struct {
	err      error
	response any
}

func (t *TestFuture) Error() error {
	return t.err
}
func (t *TestFuture) Index() uint64 {
	return 0
}
func (t *TestFuture) Response() any {
	return t.response
}

type testServiceSuite struct {
	service        *service
	mockStore      *store.MockStore
	mockRaft       *MockRaftInterface
	mockTestClient *kayakv1connect.MockKayakServiceClient
	id             ulid.ULID
	clock          *quartz.Mock
}

func setupTest(t *testing.T) *testServiceSuite {
	slog.Warn("setting up test")

	s := store.NewMockStore(t)
	r := NewMockRaftInterface(t)
	client := kayakv1connect.NewMockKayakServiceClient(t)
	c := quartz.NewMock(t)
	ts := &testServiceSuite{
		mockStore:      s,
		mockRaft:       r,
		mockTestClient: client,
		id:             ulid.Make(),
		clock:          c,
	}

	ts.service = &service{
		testLeaderClient: client,
		raft:             r,
		logger:           slog.Default(),
		store:            s,
		idGenerator: func() ulid.ULID {
			return ts.id
		},
		clock:            c,
		workerExpiration: 10 * time.Second,
	}
	return ts
}

func TestPutRecords_Leader(t *testing.T) {
	ts := setupTest(t)
	ctx := context.Background()

	records := &kayakv1.Record{
		Payload: []byte("test"),
		Id:      []byte("test"),
	}

	ts.mockStore.EXPECT().GetStream("test_stream").
		Return(&kayakv1.Stream{Name: "test_stream", PartitionCount: 1}, nil).
		Once()

	ts.mockRaft.EXPECT().State().Return(raft.Leader).Once()
	cmd := &kayakv1.RaftCommand{
		Payload: &kayakv1.RaftCommand_PutRecords{
			PutRecords: &kayakv1.PutRecords{
				Records: []*kayakv1.Record{
					{
						Id:         []byte("test"),
						InternalId: ts.id.String(),
						Payload:    []byte("test"),
						Partition:  0,
						StreamName: "test_stream",
					},
				},
				StreamName: "test_stream",
			},
		},
	}
	raw, _ := proto.Marshal(cmd)
	ts.mockRaft.EXPECT().Apply(raw, 10*time.Millisecond).Return(&TestFuture{
		response: &store.ApplyResponse{},
	}).Once()

	req := connect.NewRequest(&kayakv1.PutRecordsRequest{
		StreamName: "test_stream",
		Records:    []*kayakv1.Record{records},
	})
	_, err := ts.service.PutRecords(ctx, req)
	must.NoError(t, err)
}

func TestApplyCommand(t *testing.T) {
	ctx := context.Background()
	cases := []struct {
		name   string
		cmd    *kayakv1.RaftCommand
		leader bool
		err    error
	}{
		{
			name: "put records leader",
			cmd: &kayakv1.RaftCommand{
				Payload: &kayakv1.RaftCommand_PutRecords{
					PutRecords: &kayakv1.PutRecords{
						StreamName: "test",
						Records:    []*kayakv1.Record{},
					},
				},
			},
			leader: true,
		},
		{
			name: "put records follower",
			cmd: &kayakv1.RaftCommand{
				Payload: &kayakv1.RaftCommand_PutRecords{
					PutRecords: &kayakv1.PutRecords{
						StreamName: "test",
						Records:    []*kayakv1.Record{},
					},
				},
			},
			leader: false,
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			ts := setupTest(t)
			if tc.leader {
				ts.mockRaft.EXPECT().State().Return(raft.Leader).Once()
				raw, _ := proto.Marshal(tc.cmd)
				future := &TestFuture{
					response: &store.ApplyResponse{
						Error: tc.err,
					},
				}
				ts.mockRaft.EXPECT().Apply(raw, 10*time.Millisecond).Return(future).Once()
			} else {
				ts.mockRaft.EXPECT().State().Return(raft.Follower).Once()
				ts.mockTestClient.EXPECT().Apply(ctx, connect.NewRequest(&kayakv1.ApplyRequest{
					Command: tc.cmd,
				})).Return(connect.NewResponse(&kayakv1.ApplyResponse{}), tc.err)
			}
			err := ts.service.applyCommand(ctx, tc.cmd)
			must.ErrorIs(t, err, tc.err)
		})
	}
}

func TestRenewRegistration_HappyPath(t *testing.T) {
	ts := setupTest(t)
	ctx := context.Background()
	worker := &kayakv1.Worker{
		StreamName:          "stream",
		GroupName:           "group",
		PartitionAssignment: 1,
		Id:                  "worker1",
	}
	leaseExpires := ts.clock.Now().Add(10 * time.Second)
	cmd, _ := proto.Marshal(&kayakv1.RaftCommand{
		Payload: &kayakv1.RaftCommand_ExtendLease{
			ExtendLease: &kayakv1.ExtendLease{
				Worker: &kayakv1.Worker{
					StreamName:          "stream",
					GroupName:           "group",
					PartitionAssignment: 1,
					Id:                  "worker1",
					LeaseExpires:        leaseExpires.UnixMilli(),
				},
				ExpiresMs: leaseExpires.UnixMilli(),
			},
		},
	})

	ts.mockRaft.EXPECT().State().Return(raft.Leader).Once()
	ts.mockRaft.EXPECT().Apply(cmd, 10*time.Millisecond).Return(&TestFuture{
		response: &store.ApplyResponse{},
	}).Once()
	ts.mockStore.EXPECT().GetPartitionAssignment("stream", "group", int64(1)).
		Return("worker1", nil).Once()

	_, err := ts.service.RenewRegistration(ctx, connect.NewRequest(&kayakv1.RenewRegistrationRequest{Worker: worker}))
	must.NoError(t, err)

}

func TestRenewRegistration_InvalidInput(t *testing.T) {
	ts := setupTest(t)

	_, err := ts.service.RenewRegistration(context.Background(), connect.NewRequest(&kayakv1.RenewRegistrationRequest{}))
	must.EqError(t, err, "invalid_argument: validation error:\n - worker: value is required [required]")
}

func TestRenewRegistration_MisMatchingWorkers(t *testing.T) {
	ts := setupTest(t)
	ctx := context.Background()
	worker := &kayakv1.Worker{
		StreamName:          "stream",
		GroupName:           "group",
		PartitionAssignment: 1,
		Id:                  "worker1",
	}
	ts.mockStore.EXPECT().GetPartitionAssignment("stream", "group", int64(1)).
		Return("worker2", nil).Once()

	_, err := ts.service.RenewRegistration(ctx, connect.NewRequest(&kayakv1.RenewRegistrationRequest{Worker: worker}))
	must.ErrorIs(t, err, ErrNoAssignmentMatch)

}
func TestGetRecords(t *testing.T) {}
