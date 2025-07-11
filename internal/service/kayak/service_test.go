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
	"google.golang.org/protobuf/types/known/timestamppb"

	kayakv1 "github.com/binarymatt/kayak/gen/kayak/v1"
	"github.com/binarymatt/kayak/gen/kayak/v1/kayakv1connect"
	iraft "github.com/binarymatt/kayak/internal/raft"
	"github.com/binarymatt/kayak/internal/store"
)

type testServiceSuite struct {
	service        *service
	mockStore      *store.MockStore
	mockRaft       *iraft.MockRaftInterface
	mockTestClient *kayakv1connect.MockKayakServiceClient
	id             ulid.ULID
	clock          *quartz.Mock
}

func setupTest(t *testing.T) *testServiceSuite {
	slog.Warn("setting up test")

	s := store.NewMockStore(t)
	r := iraft.NewMockRaftInterface(t)
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
		Id:      "test",
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
						Id:              "test",
						InternalId:      ts.id.String(),
						Payload:         []byte("test"),
						Partition:       0,
						StreamName:      "test_stream",
						AcceptTimestamp: timestamppb.New(ts.clock.Now()),
					},
				},
				StreamName: "test_stream",
			},
		},
	}
	raw, _ := proto.Marshal(cmd)
	ts.mockRaft.EXPECT().Apply(raw, 10*time.Millisecond).Return(&iraft.TestFuture{
		ResponseI: &store.ApplyResponse{},
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
				future := &iraft.TestFuture{
					ResponseI: &store.ApplyResponse{
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
	ts.mockRaft.EXPECT().Apply(cmd, 10*time.Millisecond).Return(&iraft.TestFuture{
		ResponseI: &store.ApplyResponse{},
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

func TestGetStreamStatistis(t *testing.T) {
}

func TestRegisterWorker_NoUnassignedPartitions(t *testing.T) {

	ts := setupTest(t)
	ctx := context.Background()
	req := &kayakv1.RegisterWorkerRequest{
		StreamName: "test",
		Group:      "groupTest",
		Id:         "1",
	}

	ts.mockStore.EXPECT().GetStream("test").
		Return(&kayakv1.Stream{Name: "test", PartitionCount: 1}, nil).
		Once()

	// map[int64]*kayakv1.PartitionAssignment, error
	ts.mockStore.EXPECT().GetPartitionAssignments("test", "groupTest").Return(map[int64]*kayakv1.PartitionAssignment{
		1: {StreamName: "test", GroupName: "groupTest", WorkerId: "2"},
	}, nil).Once()
	_, err := ts.service.RegisterWorker(ctx, connect.NewRequest(req))
	must.Eq(t, connect.CodeOutOfRange, connect.CodeOf(err))
}
