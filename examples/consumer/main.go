package main

import (
	"context"
	"fmt"
	"log/slog"
	"os/signal"
	"syscall"
	"time"

	"github.com/binarymatt/kayak/client"
)

func main() {
	ctx := context.TODO()
	client := client.New(
		"http://localhost:8080",
		client.WithGroup("consumer"),
		client.WithStream("local"),
		client.WithWorkerId("worker1"),
	)
	if err := client.Init(ctx); err != nil {
		slog.Error("could not initilize client", "error", err)
	}
	defer client.Close()
	ctx, cancel := signal.NotifyContext(
		ctx,
		syscall.SIGINT,
		syscall.SIGTERM,
		syscall.SIGQUIT,
		syscall.SIGKILL,
	)
	defer cancel()
	ticker := time.NewTicker(5 * time.Second)
	for {
		select {
		case <-ctx.Done():
			slog.Error("context is done")
			return
		case <-ticker.C:
			if err := processRecord(ctx, client); err != nil {
				slog.Error("error consumer records", "error", err)
				return
			}
		}
	}
}
func processRecord(ctx context.Context, client *client.KayakClient) error {
	record, err := client.FetchRecord(ctx)
	if err != nil {
		slog.Error("could not fetch records", "error", err)
		return err
	}
	if record == nil {
		slog.Info("nil record returned, done processing")
		return nil
	}
	fmt.Println(record.GetId(), string(record.GetPayload()))
	if err := client.CommitRecord(ctx, record); err != nil {
		slog.Error("error commiting record", "error", err, "record", record)
		return err
	}
	return nil
}
