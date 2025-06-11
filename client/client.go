package client

import (
	"net/http"
	"sync"
	"time"

	"connectrpc.com/connect"
	"github.com/coder/quartz"

	kayakv1 "github.com/binarymatt/kayak/gen/kayak/v1"
	"github.com/binarymatt/kayak/gen/kayak/v1/kayakv1connect"
)

type KayakClient struct {
	baseUrl  string
	client   kayakv1connect.KayakServiceClient
	cfg      *config
	worker   *kayakv1.Worker
	mu       sync.Mutex
	closeFns []func()
	clock    quartz.Clock
}
type ConfigOpt func(*config)
type config struct {
	id         string
	group      string
	client     connect.HTTPClient
	streamName string
	ticker     time.Duration
}

func New(baseUrl string, fns ...ConfigOpt) *KayakClient {
	cfg := &config{
		client: http.DefaultClient,
		ticker: 30 * time.Second,
	}
	for _, fn := range fns {
		fn(cfg)
	}
	client := kayakv1connect.NewKayakServiceClient(cfg.client, baseUrl)
	return &KayakClient{
		baseUrl: baseUrl,
		client:  client,
		cfg:     cfg,
		clock:   quartz.NewReal(),
	}
}
