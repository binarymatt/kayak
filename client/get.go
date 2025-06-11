package client

import (
	"context"

	"buf.build/go/protovalidate"
	"connectrpc.com/connect"

	kayakv1 "github.com/binarymatt/kayak/gen/kayak/v1"
)

func (kc *KayakClient) GetRecords(ctx context.Context, streamName string, partition int64, start string, limit int64) ([]*kayakv1.Record, error) {
	req := &kayakv1.GetRecordsRequest{
		StreamName: streamName,
		Partition:  partition,
		StartId:    start,
		Limit:      limit,
	}
	if err := protovalidate.Validate(req); err != nil {
		return nil, err
	}
	resp, err := kc.client.GetRecords(ctx, connect.NewRequest(req))
	if err != nil {
		return nil, err
	}
	return resp.Msg.GetRecords(), nil
}
