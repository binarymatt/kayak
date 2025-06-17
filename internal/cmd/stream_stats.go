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
		stream, err := client.Stats(ctx, cmd.String("stream"))
		if err != nil {
			return err
		}
		stats := stream.GetStats()
		slog.Info("stream statistics", "name", stream.Name, "partitions", stream.PartitionCount, "ttl", stream.Ttl)
		slog.Info("records", "count", stats.RecordCount)
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
