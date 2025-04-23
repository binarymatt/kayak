package admin

import (
	"context"
	"errors"

	"connectrpc.com/connect"
	"github.com/hashicorp/raft"

	kayakv1 "github.com/binarymatt/kayak/gen/kayak/v1"
	"github.com/binarymatt/kayak/gen/kayak/v1/kayakv1connect"
)

var _ kayakv1connect.AdminServiceHandler = (*adminService)(nil)

type adminService struct {
	raft *raft.Raft
}

func (a *adminService) AddVoter(ctx context.Context, req *connect.Request[kayakv1.AddVoterRequest]) (*connect.Response[kayakv1.AddVoterResponse], error) {
	if a.raft.State() != raft.Leader {
		return nil, connect.NewError(connect.CodeInvalidArgument, errors.New("not leader"))
	}
	f := a.raft.AddVoter(raft.ServerID(req.Msg.Id), raft.ServerAddress(req.Msg.Address), 0, 0)
	if f.Error() != nil {
		return nil, connect.NewError(connect.CodeInternal, f.Error())
	}
	return connect.NewResponse(&kayakv1.AddVoterResponse{}), nil
}

func (a *adminService) Stats(ctx context.Context, req *connect.Request[kayakv1.StatsRequest]) (*connect.Response[kayakv1.StatsResponse], error) {
	stats := map[string]string{}
	for k, v := range a.raft.Stats() {
		stats[k] = v
	}
	return connect.NewResponse(&kayakv1.StatsResponse{Stats: stats}), nil
}
func (a *adminService) Leader(ctx context.Context, req *connect.Request[kayakv1.LeaderRequest]) (*connect.Response[kayakv1.LeaderResponse], error) {
	address, id := a.raft.LeaderWithID()
	return connect.NewResponse(&kayakv1.LeaderResponse{Id: string(id), Address: string(address)}), nil
}
func New(ra *raft.Raft) *adminService {
	return &adminService{raft: ra}
}
