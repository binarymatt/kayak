package main

import (
	"context"
	"errors"
	"fmt"
	"os"

	"github.com/oklog/ulid/v2"
	"github.com/spaolacci/murmur3"
	"github.com/urfave/cli/v2"
	"golang.org/x/exp/slog"
	"golang.org/x/sync/errgroup"

	"github.com/binarymatt/kayak/client"
	kayakv1 "github.com/binarymatt/kayak/gen/kayak/v1"
)

const (
	topic = "messages"
	index = 1
	limit = 100
)

func createRecords(count int) (records []*kayakv1.Record) {
	for i := 0; i < count; i++ {
		records = append(records, &kayakv1.Record{
			Topic:   topic,
			Id:      ulid.Make().String(),
			Payload: []byte(fmt.Sprint(i)),
		})
	}
	return
}
func mainHash() {
	keys := []string{
		"2Y5PvAX8aY98M5u2oQMrobCahtL",
		"2Y5Pw5vJP4ATrUdjNENthrlucjb",
		"2Y5PwwUdt1Sl8wIZBEqjNgHx4jx",
		"2Y5Pxp8JKQeYI7rJpbwL4s9emxc",
		"test",
		"test2",
	}
	for _, key := range keys {
		fmt.Println(murmur3.Sum64([]byte(key)) % 2)
	}
}
func mainClient() {
	ctx := context.Background()
	c := client.New(
		client.NewConfig(""),
		client.WithAddress("http://127.0.0.1:8081"),
		client.WithTopic("test"),
		client.WithConsumerID("cli"),
		client.WithConsumerGroup("validation"),
	)
	start := "01HF75KXVWV74R8QZ8E0R426QN"
	records, err := c.GetRecords(ctx, "test", start, 10)

	if err != nil {
		slog.Error("could not fetch record", "error", err)
		return
	}
	slog.Info("record returned", "records", records, "len", len(records))
	/*
		record, err := c.FetchRecord(ctx)
		if err != nil {
			slog.Error("could not fetch record", "error", err)
			return
		}
		slog.Info("record returned", "record", record)
		if err := c.CommitRecord(ctx, record); err != nil {
			slog.Error("coudl not commit record", "error", err)
		}
	*/
}

type testing struct {
	c *client.Client
}

func (t *testing) consumer(cctx *cli.Context) error {
	for {
		record, err := t.c.FetchRecord(cctx.Context)
		if err != nil {
			return err
		}
		slog.Info("record retrieved", "record", record)
		if record == nil {
			break
		}
		if err := t.c.CommitRecord(cctx.Context, record); err != nil {
			slog.Error("could not commit record", "record", record)
			return err
		}
		slog.Info("record committed")
	}
	return nil
}
func (t *testing) initTest(cctx *cli.Context) error {
	ctx := cctx.Context
	if err := t.c.CreateTopic(ctx, topic); err != nil {
		return err
	}
	records := []*kayakv1.Record{
		{Headers: map[string]string{"number": "one"}, Payload: []byte(`{"json":"payload"}`)},
		{Payload: []byte(`{"json":"payload2"}`)},
		{Payload: []byte(`{"json":"payload3"}`)},
		{Payload: []byte("text payload")},
	}
	return t.c.PutRecords(cctx.Context, records...)
}
func (t *testing) cleanup(cctx *cli.Context) error {
	if err := t.c.DeleteTopic(cctx.Context, topic); err != nil {
		return err
	}
	return nil
}
func (t *testing) query(cctx *cli.Context) error {
	ctx := cctx.Context
	start := cctx.String("start")
	records, err := t.c.GetRecords(ctx, topic, start, 10)

	if err != nil {
		slog.Error("could not fetch record", "error", err)
		return err
	}
	slog.Info("record returned", "records", records, "len", len(records))
	return nil
}
func (t *testing) fetch(cctx *cli.Context) error {
	ctx := cctx.Context
	record, err := t.c.FetchRecord(ctx)

	if err != nil {
		slog.Error("could not fetch record", "error", err)
		return err
	}
	slog.Info("record returned", "records", record)
	return nil
}
func (t *testing) commit(cctx *cli.Context) error {
	ctx := cctx.Context
	id := cctx.String("id")
	if id == "" {
		return errors.New("empty id for commit")
	}

	return t.c.CommitRecord(ctx, &kayakv1.Record{
		Topic: topic,
		Id:    id,
	})
}
func (t *testing) get(cctx *cli.Context) error {
	ctx := cctx.Context
	records, err := t.c.GetRecords(ctx, topic, "", 10)
	if err != nil {
		slog.Error("could not query records")
		return err
	}
	for _, record := range records {
		slog.Info("returned record", "id", record.Id)
	}
	return nil
}
func worker(ctx context.Context, c *client.Client, id int, recordCount int) error {
	slog.Info("worker started", "id", id, "count", recordCount)
	i := 0
	for i < recordCount {
		if i%4 == id {
			// create record
			record := &kayakv1.Record{
				Headers: map[string]string{"number": fmt.Sprintf("%d", i)},
				Payload: []byte(fmt.Sprintf(`{"json":"payload", id:"%d"}`, id)),
			}
			if err := c.PutRecords(ctx, record); err != nil {
				return err
			}
		}
		i++
	}
	slog.Info("done with worker", "id", id)
	return nil
}
func (t *testing) load(cctx *cli.Context) error {
	number := cctx.Int("count")
	g := new(errgroup.Group)
	g.Go(func() error {
		return worker(cctx.Context, t.c, 0, number)
	})
	g.Go(func() error {
		return worker(cctx.Context, t.c, 1, number)
	})
	g.Go(func() error {
		return worker(cctx.Context, t.c, 2, number)
	})
	g.Go(func() error {
		return worker(cctx.Context, t.c, 3, number)
	})

	slog.Info("waiting on records")
	return g.Wait()
}
func main() {

	c := client.New(
		client.NewConfig(""),
		client.WithAddress("http://127.0.0.1:8081"),
		client.WithTopic(topic),
		client.WithConsumerID("test"),
		client.WithConsumerGroup("validation"),
	)
	t := &testing{c: c}
	app := &cli.App{
		Commands: []*cli.Command{
			{
				Name:   "init",
				Action: t.initTest,
			},
			{
				Name:   "cleanup",
				Action: t.cleanup,
			},
			{
				Name:   "get",
				Action: t.get,
			},
			{
				Name:   "consume",
				Action: t.consumer,
			},
			{
				Name: "query",
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:  "start",
						Value: "01HF75KXVWV74R8QZ8E0R426QN",
					},
				},
				Action: t.query,
			},
			{
				Name: "commit",
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:  "id",
						Value: "",
					},
				},
				Action: t.commit,
			},
			{
				Name:   "fetch",
				Action: t.fetch,
			},
			{
				Name: "load",
				Flags: []cli.Flag{
					&cli.IntFlag{
						Name:  "count",
						Value: 5,
					},
					&cli.IntFlag{
						Name:  "pool",
						Value: 8,
					},
				},
				Action: t.load,
			},
		},
	}
	if err := app.Run(os.Args); err != nil {
		slog.Error("could not run cli", "error", err)
	}
}
