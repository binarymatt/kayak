package cli

import (
	"errors"
	"fmt"

	"connectrpc.com/connect"
	nomad "github.com/hashicorp/nomad/api"
	"github.com/urfave/cli/v2"
	"log/slog"

	adminv1 "github.com/binarymatt/kayak/gen/admin/v1"
)

var (
	ErrMultipleLeaders = errors.New("multiple leaders found")
	ErrMissingLeader   = errors.New("no leader found")
)

// SetupNomadCluster promotes first nomad alloc to leader and joins other allocs to leader.
func SetupNomadCluster(cctx *cli.Context) error {
	host := cctx.String("nomad_address")

	//kayakHosts := []khost{}
	var leader string
	nc, err := nomad.NewClient(&nomad.Config{
		Address: host,
	})
	if err != nil {
		return err
	}

	allocs, _, err := nc.Jobs().Allocations("kayak", false, nil)

	if err != nil {
		return err
	}
	for index, alloc := range allocs {
		if alloc.DesiredStatus == "stop" {
			continue
		}
		a, _, _ := nc.Allocations().Info(alloc.ID, nil)
		ports := a.AllocatedResources.Shared.Ports
		serfInfo := ""
		grpcInfo := ""

		for _, port := range ports {
			if port.Label == "serf" {
				serfInfo = fmt.Sprintf("%s:%d", port.HostIP, port.Value)
			}
			if port.Label == "grpc" {
				grpcInfo = fmt.Sprintf("%s:%d", port.HostIP, port.Value)
			}
		}
		slog.Info(a.ID, "serf", serfInfo, "grpc", grpcInfo)
		ac := buildAdminClient(fmt.Sprintf("http://%s", grpcInfo))
		if index == 0 {
			leader = serfInfo
		}
		slog.Info("trying to join cluster", "leader", leader)
		_, err := ac.Join(cctx.Context, connect.NewRequest(&adminv1.JoinRequest{
			Address: leader,
		}))
		if err != nil {
			slog.Error("oops", "error", err)
		}
		resp, err := ac.State(cctx.Context, connect.NewRequest(&adminv1.StateRequest{}))
		if err != nil {
			return err
		}
		slog.Info("server state", "state", resp.Msg.State)
	}
	return nil
}
