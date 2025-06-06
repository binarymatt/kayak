package admin

import (
	"context"
	"testing"
	"time"

	"buf.build/go/protovalidate"
	"connectrpc.com/connect"
	"github.com/hashicorp/raft"
	"github.com/shoenig/test/must"

	kayakv1 "github.com/binarymatt/kayak/gen/kayak/v1"
	iraft "github.com/binarymatt/kayak/internal/raft"
)

type testContainer struct {
	mockRaft *iraft.MockRaftInterface
	service  *adminService
}

func setupTest(t *testing.T) *testContainer {
	t.Helper()
	mockRaft := iraft.NewMockRaftInterface(t)
	s := &adminService{
		raft: mockRaft,
	}
	return &testContainer{
		service:  s,
		mockRaft: mockRaft,
	}
}

func TestAddVoter_NotLeader(t *testing.T) {
	tc := setupTest(t)
	ctx := context.Background()
	req := &kayakv1.AddVoterRequest{}

	tc.mockRaft.EXPECT().State().Return(raft.Follower).Once()

	resp, err := tc.service.AddVoter(ctx, connect.NewRequest(req))
	must.ErrorIs(t, err, ErrNotLeader)
	must.Nil(t, resp)
}

func TestAddVoter_MissingArguments(t *testing.T) {
	tc := setupTest(t)
	ctx := context.Background()
	req := &kayakv1.AddVoterRequest{}

	tc.mockRaft.EXPECT().State().Return(raft.Leader).Once()

	resp, err := tc.service.AddVoter(ctx, connect.NewRequest(req))
	ve := &protovalidate.ValidationError{}
	must.ErrorAs(t, err, &ve)
	must.Nil(t, resp)
	must.Len(t, 2, ve.Violations)
	must.Eq(t, ve.Violations[0].FieldDescriptor.TextName(), "id")
	must.Eq(t, ve.Violations[0].RuleDescriptor.TextName(), "required")
	must.True(t, ve.Violations[0].RuleValue.Bool())
	must.Eq(t, ve.Violations[1].FieldDescriptor.TextName(), "address")
	must.Eq(t, ve.Violations[1].RuleDescriptor.TextName(), "required")
	must.True(t, ve.Violations[1].RuleValue.Bool())
}
func TestAddVoter_HappyPath(t *testing.T) {
	tc := setupTest(t)
	ctx := context.Background()
	req := &kayakv1.AddVoterRequest{
		Id:      "worker2",
		Address: "127.0.0.1",
	}

	tc.mockRaft.EXPECT().State().Return(raft.Leader).Once()
	tc.mockRaft.EXPECT().
		AddVoter(raft.ServerID("worker2"), raft.ServerAddress("127.0.0.1"), uint64(0), 0*time.Second).
		Return(&iraft.TestFuture{})

	_, err := tc.service.AddVoter(ctx, connect.NewRequest(req))
	must.NoError(t, err)

}
