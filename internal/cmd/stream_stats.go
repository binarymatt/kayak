package cmd

import (
	"context"
	"log/slog"

	"github.com/urfave/cli/v3"

	"github.com/binarymatt/kayak/client"
)

var streamStatsCommand = &cli.Command{
	Name: "stats",
	Action: func(ctx context.Context, cmd *cli.Command) error {
		client := ctx.Value(clientCtxKey).(*client.KayakClient)
		stats, err := client.Stats(ctx, cmd.String("stream"))
		if err != nil {
			return err
		}
		slog.Info("stream statistics", "stream_name", cmd.String("stream"))
		slog.Info("", "record_count", stats.RecordCount)
		for k, v := range stats.PartitionCounts {
			slog.Info("partition info", "id", k, "count", v)
		}
		for _, group := range stats.Groups {
			slog.Info("group info", "name", group.Name, "stream", group.StreamName)
			for k, v := range group.PartitionPositions {
				slog.Info("partition position", "parition_id", k, "position", v)
			}
		}

		return nil
	},
}
