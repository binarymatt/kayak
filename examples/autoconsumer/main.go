package main

import (
	"context"
	"log"
	"log/slog"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/lmittmann/tint"
	"github.com/urfave/cli/v3"

	"github.com/binarymatt/kayak/client/auto"
	kayakv1 "github.com/binarymatt/kayak/gen/kayak/v1"
)

func main() {
	cmd := &cli.Command{
		Name:  "greet",
		Usage: "fight the loneliness!",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    "log_level",
				Value:   "info",
				Usage:   "debug, info, warn, error",
				Aliases: []string{"l"},
			},
			&cli.StringFlag{
				Name:    "log_format",
				Value:   "console",
				Usage:   "valid values are console, json",
				Aliases: []string{"f"},
			},
		},
		Action: func(ctx context.Context, cmd *cli.Command) error {
			ctx, cancel := signal.NotifyContext(
				ctx,
				syscall.SIGINT,
				syscall.SIGTERM,
				syscall.SIGQUIT,
				syscall.SIGKILL,
			)
			defer cancel()
			configureLogging(cmd.String("log_level"), cmd.String("log_format"))
			client := auto.New(auto.Config{
				BaseUrl:      "http://localhost:8080",
				Stream:       "local",
				Group:        "consumer",
				WorkerPrefix: "worker",
				MaxWorkers:   5,
			}, processRecord)

			return client.Run(ctx)

		},
	}

	if err := cmd.Run(context.Background(), os.Args); err != nil {
		log.Fatal(err)
	}
}
func configureLogging(levelStr, formatStr string) {

	var logger *slog.Logger
	var level slog.Leveler
	switch strings.ToLower(levelStr) {
	case "info":
		level = slog.LevelInfo
	case "debug":
		level = slog.LevelDebug
	case "warn":
		level = slog.LevelWarn
	case "error":
		level = slog.LevelError
	default:
		level = slog.LevelInfo
	}

	opts := &slog.HandlerOptions{Level: level, AddSource: true}
	if formatStr == "console" {
		logger = slog.New(tint.NewHandler(os.Stderr, &tint.Options{
			Level:     level,
			AddSource: true,
		}))
	} else {
		logger = slog.New(slog.NewJSONHandler(os.Stderr, opts))
	}

	slog.SetDefault(logger)
}
func processRecord(ctx context.Context, record *kayakv1.Record) error {
	logger := ctx.Value(auto.ContextKey).(*slog.Logger)
	logger.Debug("processing record", "id", record.GetId(), "ts", record.AcceptTimestamp.AsTime())
	return nil
}
