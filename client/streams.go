package client

import (
	"context"
	"log/slog"

	"connectrpc.com/connect"
	"google.golang.org/protobuf/proto"

	kayakv1 "github.com/binarymatt/kayak/gen/kayak/v1"
)

func (kc *KayakClient) Stats(ctx context.Context, stream string) (*kayakv1.Stream, error) {
	resp, err := kc.client.GetStream(ctx, connect.NewRequest(&kayakv1.GetStreamRequest{Name: stream, IncludeStats: proto.Bool(true)}))
	if err != nil {
		slog.Error("error durring client stats call", "resp", resp)
		return nil, err
	}

	return resp.Msg.GetStream(), nil
}
