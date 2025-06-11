package client

import (
	"context"
	"fmt"
	"log/slog"

	"connectrpc.com/connect"

	kayakv1 "github.com/binarymatt/kayak/gen/kayak/v1"
)

func (kc *KayakClient) timer(ctx context.Context) error {
	ticker := kc.clock.NewTicker(kc.cfg.ticker)
	for {
		select {
		case <-ctx.Done():
			return nil
		case <-ticker.C:
			if err := kc.extendLease(ctx); err != nil {
				return err
			}
		}
	}
}
func (kc *KayakClient) Init(ctx context.Context) error {
	ctx, cancel := context.WithCancel(ctx)
	if err := kc.Register(ctx); err != nil {
		cancel()
		return err
	}
	go kc.timer(ctx)
	kc.closeFns = append(kc.closeFns, cancel)
	return nil
}

func (kc *KayakClient) extendLease(ctx context.Context) error {
	kc.mu.Lock()
	defer kc.mu.Unlock()
	req := &kayakv1.RenewRegistrationRequest{
		Worker: kc.worker,
	}
	_, err := kc.client.RenewRegistration(ctx, connect.NewRequest(req))
	return err
}

func (kc *KayakClient) Register(ctx context.Context) error {
	req := &kayakv1.RegisterWorkerRequest{
		StreamName: kc.cfg.streamName,
		Group:      kc.cfg.group,
		Id:         kc.cfg.id,
	}
	resp, err := kc.client.RegisterWorker(ctx, connect.NewRequest(req))
	if err != nil {
		fmt.Println(err)
		return err
	}
	kc.worker = resp.Msg.GetWorker()
	slog.Info("returned config", "worker", kc.worker)
	return nil
}
func (kc *KayakClient) Close() {
	for _, fn := range kc.closeFns {
		fn()
	}
}

func (kc *KayakClient) FetchRecord(ctx context.Context) (*kayakv1.Record, error) {
	req := &kayakv1.FetchRecordsRequest{
		Worker: kc.worker,
		Limit:  1,
	}
	resp, err := kc.client.FetchRecords(ctx, connect.NewRequest(req))
	slog.Debug("fetched records", "resp", resp, "error", err)
	if len(resp.Msg.GetRecords()) > 0 {
		return resp.Msg.GetRecords()[0], nil
	}
	return nil, nil
}

func (kc *KayakClient) CommitRecord(ctx context.Context, record *kayakv1.Record) error {
	kc.mu.Lock()
	defer kc.mu.Unlock()
	req := &kayakv1.CommitRecordRequest{
		Worker: kc.worker,
		Record: record,
	}
	_, err := kc.client.CommitRecord(ctx, connect.NewRequest(req))
	return err
}
