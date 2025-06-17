package main

import (
	"context"
	"encoding/json"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"connectrpc.com/connect"
	"github.com/oklog/ulid/v2"

	kayakv1 "github.com/binarymatt/kayak/gen/kayak/v1"
	"github.com/binarymatt/kayak/gen/kayak/v1/kayakv1connect"
)

type samplePayload struct {
	Id          string `json:"id"`
	Description string `json:"description"`
	Index       int
}

func generateRecords() []*kayakv1.Record {
	records := []*kayakv1.Record{}
	for i := range 100 {
		id := ulid.Make().String()
		payload := samplePayload{
			Id:          id,
			Description: "sample payload",
			Index:       i,
		}
		raw, _ := json.Marshal(payload)
		records = append(records, &kayakv1.Record{
			Id:      id,
			Payload: raw,
		})

	}
	return records
}
func main() {
	ctx, cancel := signal.NotifyContext(
		context.TODO(),
		syscall.SIGINT,
		syscall.SIGTERM,
		syscall.SIGQUIT,
		syscall.SIGKILL,
	)
	defer cancel()
	client := kayakv1connect.NewKayakServiceClient(http.DefaultClient,
		"http://localhost:8080",
	)
	for range 10000 {
		records := generateRecords()
		_, err := client.PutRecords(ctx, connect.NewRequest(&kayakv1.PutRecordsRequest{
			StreamName: "local",
			Records:    records,
		}))
		if err != nil {
			slog.Error("error adding records", "error", err)
			os.Exit(-1)
		}
		time.Sleep(1 * time.Millisecond)
	}
}
