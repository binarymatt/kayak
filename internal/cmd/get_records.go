package cmd

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/urfave/cli/v3"

	"github.com/binarymatt/kayak/client"
)

var getRecords = &cli.Command{
	Name:     "get",
	Category: "records",
	Action: func(ctx context.Context, cmd *cli.Command) error {
		client, ok := ctx.Value(clientCtxKey).(*client.KayakClient)
		if !ok {
			return fmt.Errorf("could not get client")
		}
		stream := cmd.String("stream")
		partition := cmd.Int64("partition")

		records, err := client.GetRecords(ctx, stream, partition, "", 10)
		if err != nil {
			return err
		}
		for _, record := range records {
			slog.Info("",
				"accept_timestamp", record.GetAcceptTimestamp().AsTime(),
				"id", record.GetId(),
				"internal_id", record.GetInternalId(),
				"stream", record.GetStreamName(),
				"partition", record.GetPartition(),
				"payload", string(record.GetPayload()),
			)
		}
		return nil
	},
}
