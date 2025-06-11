package client

import (
	"context"

	"connectrpc.com/connect"

	kayakv1 "github.com/binarymatt/kayak/gen/kayak/v1"
)

func (kc *KayakClient) Stats(ctx context.Context, stream string) (*kayakv1.StreamStats, error) {
	resp, err := kc.client.GetStreamStatistics(ctx, connect.NewRequest(&kayakv1.GetStreamStatisticsRequest{Name: stream}))
	return resp.Msg.GetStreamStats(), err
}
