package main

import (
	"os"

	"github.com/urfave/cli/v2"
	"log/slog"

	internal "github.com/binarymatt/kayak/internal/cli"
	"github.com/binarymatt/kayak/internal/log"
)

func main() {
	app := &cli.App{
		Before: func(ctx *cli.Context) error {
			consoleLogging := ctx.Bool("console")
			var logger *slog.Logger
			if consoleLogging {
				logger = slog.New(log.NewTextHandler(os.Stderr, nil))
			} else {
				logger = slog.New(slog.NewJSONHandler(os.Stdout, nil))
			}
			slog.SetDefault(logger)
			return nil
		},
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    "host",
				Value:   "http://localhost:8081",
				Usage:   "kayak command host",
				EnvVars: []string{"KAYAK_HOST"},
			},
			&cli.BoolFlag{
				Name:    "console",
				Value:   true,
				Aliases: []string{"c"},
			},
		},
		Commands: []*cli.Command{
			{
				Name: "admin",
				Subcommands: []*cli.Command{
					{
						Name:   "setup_cluster",
						Action: internal.SetupCluster,
						Flags: []cli.Flag{
							&cli.StringFlag{
								Name:     "nomad_address",
								Required: true,
							},
						},
					},
					{
						Name:   "leader",
						Action: internal.AdminLeader,
					},
					{
						Name: "add_voter",
						Flags: []cli.Flag{
							&cli.StringFlag{
								Name:     "id",
								Required: true,
							},
							&cli.StringFlag{
								Name:     "address",
								Required: true,
							},
						},
						Action: internal.AddVoter,
					},
					{
						Name:   "get_configuration",
						Action: internal.GetConfiguration,
					},
					{
						Name:   "remove_server",
						Action: internal.RemoveServer,
						Flags: []cli.Flag{
							&cli.StringFlag{
								Name:     "id",
								Required: true,
							},
						},
					},
				},
			},
			{
				Name: "records",
				Subcommands: []*cli.Command{
					{
						Name:   "create",
						Action: internal.CreateRecords,
						Flags: []cli.Flag{
							&cli.StringFlag{
								Name:     "data",
								Required: true,
							},
						},
					},
					{
						Name:   "list",
						Action: internal.ListRecords,
					},
				},
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:     "topic",
						Required: true,
					},
				},
			},
			{
				Name:   "stats",
				Action: internal.Stats,
			},
			{
				Name: "topics",
				Subcommands: []*cli.Command{
					{
						Name: "create",
						Flags: []cli.Flag{
							&cli.StringFlag{
								Name:     "name",
								Required: true,
							},
						},
						Action: internal.CreateTopic,
					},
					{
						Name:   "list",
						Action: internal.ListTopics,
					},
					{
						Name:   "delete",
						Action: internal.DeleteTopic,
						Flags: []cli.Flag{
							&cli.BoolFlag{
								Name:  "archive",
								Value: false,
							},
							&cli.StringFlag{
								Name:     "name",
								Required: true,
							},
						},
					},
				},
			},
		},
	}
	if err := app.Run(os.Args); err != nil {
		slog.Error("could not run cli", "error", err)
	}
}
