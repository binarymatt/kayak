package client

import (
	"context"

	"connectrpc.com/connect"

	kayakv1 "github.com/binarymatt/kayak/gen/kayak/v1"
)

func (kc *KayakClient) PutRecords(ctx context.Context, streamName string, records ...*kayakv1.Record) error {
	req := connect.NewRequest(&kayakv1.PutRecordsRequest{
		StreamName: streamName,
		Records:    records,
	})
	_, err := kc.client.PutRecords(ctx, req)
	return err
}
