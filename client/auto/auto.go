package auto

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"sync"
	"time"

	"connectrpc.com/connect"
	"github.com/coder/quartz"
	"github.com/oklog/ulid/v2"
	"golang.org/x/sync/errgroup"

	"github.com/binarymatt/kayak/client"
	kayakv1 "github.com/binarymatt/kayak/gen/kayak/v1"
)

type processFn = func(context.Context, *kayakv1.Record) error
type Config struct {
	MaxWorkers    int
	WorkerPrefix  string
	Group         string
	Stream        string
	BaseUrl       string
	SleepDuration time.Duration
}
type AutoClient struct {
	cfg                Config
	fn                 processFn
	clock              quartz.Clock
	eg                 *errgroup.Group
	mu                 sync.Mutex
	currentWorkerCount int
}
type contextKeyType struct{}

// contextKey is the key used for the context to store the logger.
var ContextKey = contextKeyType{}

func New(cfg Config, fn processFn) *AutoClient {
	if cfg.SleepDuration == 0 {
		cfg.SleepDuration = time.Second * 2
	}
	return &AutoClient{
		cfg:   cfg,
		fn:    fn,
		clock: quartz.NewReal(),
	}
}

// work is run in a goroutine with it's own kayak client.
func (ac *AutoClient) work(ctx context.Context, uniqueID string) error {
	defer func() {
		ac.mu.Lock()
		ac.currentWorkerCount--
		ac.mu.Unlock()
	}()
	id := fmt.Sprintf("%s-%s", ac.cfg.WorkerPrefix, uniqueID)
	slog.Warn("starting worker", "id", id)

	kc := client.New(
		ac.cfg.BaseUrl,
		client.WithGroup(ac.cfg.Group),
		client.WithStream(ac.cfg.Stream),
		client.WithWorkerId(id),
	)
	if err := kc.Init(ctx); err != nil {
		// there are no more partition assignments, exit
		if connect.CodeOf(err) == connect.CodeOutOfRange {
			slog.Warn("all partitions have been assigned, discarding worker", "id", id)
			return nil
		}
		slog.Error("error initializing worker", "error", err)
		return err
	}
	defer kc.Close()
	logger := slog.With("id", id, "stream", ac.cfg.Stream, "group", ac.cfg.Group)
	for {
		select {
		case <-ctx.Done():
			return nil
		default:
			record, err := kc.FetchRecord(ctx)
			if err != nil {
				logger.Error("could not fetch record", "error", err)
				return err
			}
			if record == nil {
				logger.Debug("record is nil, sleeping")
				sleeper := ac.clock.NewTimer(ac.cfg.SleepDuration)
				select {
				case <-sleeper.C:
				case <-ctx.Done():
				}
				sleeper.Stop()
				continue
			}
			ctx := context.WithValue(ctx, ContextKey, logger)

			if err := ac.fn(ctx, record); err != nil {
				logger.Error("error processing record", "eror", err, "record", record)
				return err
			}
			if err := kc.CommitRecord(ctx, record); err != nil {
				logger.Error("error committing record", "eror", err, "record", record)
				return err
			}
		}
	}
}

func (ac *AutoClient) spinupWorkers(ctx context.Context) {
	slog.Info("attemping worker spin up", "max_worker_count", ac.cfg.MaxWorkers, "worker_count", ac.currentWorkerCount)

	for range ac.cfg.MaxWorkers {
		if ac.currentWorkerCount < ac.cfg.MaxWorkers {
			ac.mu.Lock()
			ac.eg.Go(func() error {
				return ac.work(ctx, ulid.Make().String())
			})
			ac.currentWorkerCount++
			ac.mu.Unlock()
		}
	}

}
func (ac *AutoClient) Run(ctx context.Context) error {
	g, ctx := errgroup.WithContext(ctx)
	ac.eg = g
	ac.spinupWorkers(ctx)
	ac.eg.Go(func() error {
		ticker := ac.clock.NewTicker(30 * time.Second)
		for {
			select {
			case <-ctx.Done():
				slog.Warn("context has been canceled", "error", ctx.Err())
				if errors.Is(ctx.Err(), context.Canceled) {
					return nil
				}
				return ctx.Err()
			case <-ticker.C:
				ac.spinupWorkers(ctx)
			}
		}
	})
	return g.Wait()
}
