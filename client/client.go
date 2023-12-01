package client

import (
	"context"
	"errors"
	"net/http"

	"connectrpc.com/connect"
	kayakv1 "github.com/binarymatt/kayak/gen/kayak/v1"
	"github.com/binarymatt/kayak/gen/kayak/v1/kayakv1connect"
	"github.com/binarymatt/kayak/internal/store"
	"github.com/hashicorp/go-retryablehttp"
	"github.com/oklog/ulid/v2"
	"log/slog"
)

var (
	ErrConsumerAlreadyRegistered = store.ErrConsumerAlreadyRegistered
	ErrConsumerGroupExists       = store.ErrConsumerGroupExists
)

type Config struct {
	Address       string
	ConsumerID    string
	ConsumerGroup string
	HTTPClient    *http.Client
	ID            string
	Topic         string
}
type CfgOption func(*Config)

func NewConfig(id string) *Config {
	if id == "" {
		id = ulid.Make().String()
	}
	client := retryablehttp.NewClient()
	client.RetryMax = 3
	client.Logger = slog.With("client", "")

	return &Config{
		ID:         id,
		HTTPClient: client.StandardClient(),
	}
}

func WithAddress(address string) CfgOption {
	return func(c *Config) {
		c.Address = address
	}
}

func WithTopic(topic string) CfgOption {
	return func(c *Config) {
		c.Topic = topic
	}
}

func WithConsumerID(consumer string) CfgOption {
	return func(c *Config) {
		c.ConsumerID = consumer
	}
}

func WithConsumerGroup(group string) CfgOption {
	return func(c *Config) {
		c.ConsumerGroup = group
	}
}

func WithHTTPClient(client *http.Client) CfgOption {
	return func(c *Config) {
		c.HTTPClient = client
	}
}

type Client struct {
	cfg    *Config
	client kayakv1connect.KayakServiceClient
}

func (c *Client) UpdateConfig(options ...CfgOption) {
	for _, f := range options {
		f(c.cfg)
	}
}

func (c *Client) PutRecords(ctx context.Context, records ...*kayakv1.Record) error {
	topic := c.cfg.Topic
	_, err := c.client.PutRecords(ctx, connect.NewRequest(&kayakv1.PutRecordsRequest{
		Topic:   topic,
		Records: records,
	}))
	return err
}

func (c *Client) RegisterConsumer(ctx context.Context) error {
	_, err := c.client.RegisterConsumer(ctx, connect.NewRequest(&kayakv1.RegisterConsumerRequest{
		Consumer: &kayakv1.TopicConsumer{
			Id:    c.cfg.ConsumerID,
			Topic: c.cfg.Topic,
			Group: c.cfg.ConsumerGroup,
		},
	}))
	return err
}
func (c *Client) CreateConsumerGroup(ctx context.Context, partitionCount int64) error {
	if c.cfg.ConsumerGroup == "" {
		return errors.New("missing group config")
	}
	_, err := c.client.CreateConsumerGroup(ctx, connect.NewRequest(&kayakv1.CreateConsumerGroupRequest{
		Group: &kayakv1.ConsumerGroup{
			Name:           c.cfg.ConsumerGroup,
			Topic:          c.cfg.Topic,
			PartitionCount: partitionCount,
			Hash:           kayakv1.Hash_HASH_MURMUR3,
		},
	}))
	return err
}

func (c *Client) FetchRecord(ctx context.Context) (*kayakv1.Record, error) {
	req := connect.NewRequest(&kayakv1.FetchRecordRequest{
		Topic:         c.cfg.Topic,
		ConsumerGroup: c.cfg.ConsumerGroup,
		ConsumerId:    c.cfg.ConsumerID,
	})
	slog.Info("kayak client fetch", "topic", c.cfg.Topic, "consumer", c.cfg.ConsumerID)
	resp, err := c.client.FetchRecord(ctx, req)
	if err != nil {
		return nil, err
	}
	return resp.Msg.GetRecord(), err
}

func (c *Client) CommitRecord(ctx context.Context, record *kayakv1.Record) error {
	_, err := c.client.CommitRecord(ctx, connect.NewRequest(&kayakv1.CommitRecordRequest{
		Consumer: &kayakv1.TopicConsumer{
			Group:    c.cfg.ConsumerGroup,
			Topic:    c.cfg.Topic,
			Id:       c.cfg.ConsumerID,
			Position: record.GetId(),
		},
	}))
	return err
}

func (c *Client) GetRecords(ctx context.Context, topic string, start string, limit int64) ([]*kayakv1.Record, error) {
	req := connect.NewRequest(&kayakv1.GetRecordsRequest{
		Topic: topic,
		Start: start,
		Limit: limit,
	})
	resp, err := c.client.GetRecords(ctx, req)
	return resp.Msg.GetRecords(), err
}

func (c *Client) CreateTopic(ctx context.Context, topic string) error {
	_, err := c.client.CreateTopic(ctx, connect.NewRequest(&kayakv1.CreateTopicRequest{
		Name: topic,
	}))
	return err
}

func (c *Client) DeleteTopic(ctx context.Context, topic string) error {
	req := connect.NewRequest(&kayakv1.DeleteTopicRequest{
		Topic:   topic,
		Archive: false,
	})
	_, err := c.client.DeleteTopic(ctx, req)
	return err
}

func New(cfg *Config, options ...CfgOption) *Client {
	for _, f := range options {
		f(cfg)
	}
	c := kayakv1connect.NewKayakServiceClient(cfg.HTTPClient, cfg.Address)
	return &Client{
		client: c,
		cfg:    cfg,
	}
}
