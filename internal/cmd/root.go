package cmd

import (
	"context"
	"log/slog"
	"os"

	"github.com/urfave/cli/v3"

	"github.com/binarymatt/kayak/client"
)

var clientCtxKey = &ctxKey{}

type ctxKey struct{}

var root = &cli.Command{
	Name: "kyk",
	Action: func(context.Context, *cli.Command) error {
		slog.Info("root, nothing to see here")
		return nil
	},
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name:  "host",
			Value: "http://localhost:8080",
		},
		&cli.StringFlag{
			Name:  "stream",
			Value: "local",
		},
		&cli.Int64Flag{
			Name:    "partition",
			Aliases: []string{"p"},
			Value:   0,
		},
	},
	Before: func(ctx context.Context, cmd *cli.Command) (context.Context, error) {
		client := client.New(cmd.String("host"), client.WithStream(cmd.String("stream")))
		ctx = context.WithValue(ctx, clientCtxKey, client)
		return ctx, nil
	},
	Commands: []*cli.Command{
		records,
		streams,
	},
}

func Execute() {
	if err := root.Run(context.Background(), os.Args); err != nil {
		slog.Error("failed during run", "error", err, "args", os.Args)
		os.Exit(1)
	}
}
