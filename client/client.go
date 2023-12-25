package client

import (
	"context"
	"errors"
	"net/http"
	"strings"

	"connectrpc.com/connect"
	kayakv1 "github.com/binarymatt/kayak/gen/kayak/v1"
	"github.com/binarymatt/kayak/gen/kayak/v1/kayakv1connect"
	"github.com/binarymatt/kayak/internal/service"
	"github.com/binarymatt/kayak/internal/store"
	"github.com/hashicorp/go-retryablehttp"
	"github.com/oklog/ulid/v2"
	"log/slog"
)

var (
	ErrConsumerAlreadyRegistered = store.ErrConsumerAlreadyRegistered
	ErrTopicEmpty                = errors.New("topic cannot be empty")
	ErrCfgConsumerIDEmpty        = errors.New("consumer id must be set when a consumer group exists")
	ErrCfgConsumerGroupEmpty     = errors.New("consumer group must be set when a consumer id exists")
)

type Config struct {
	Address       string
	ConsumerID    string
	ConsumerGroup string
	HTTPClient    *http.Client
	ID            string
	Topic         string
	Partitions    []int64
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
func WithPartitions(partitions []int64) CfgOption {
	return func(c *Config) {
		c.Partitions = partitions
	}
}

type Client struct {
	cfg    *Config
	client kayakv1connect.KayakServiceClient
}

func (c *Client) validateConfig() error {
	if c.cfg.Topic == "" {
		return ErrTopicEmpty
	}
	if c.cfg.ConsumerGroup != "" || c.cfg.ConsumerID != "" {
		// consumer id must be set
		if c.cfg.ConsumerID == "" {
			return ErrCfgConsumerIDEmpty
		}
		// consumer group must be set
		if c.cfg.ConsumerGroup == "" {
			return ErrCfgConsumerGroupEmpty
		}
	}
	return nil
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
	if err := c.validateConfig(); err != nil {
		return err
	}
	_, err := c.client.RegisterConsumer(ctx, connect.NewRequest(&kayakv1.RegisterConsumerRequest{
		Consumer: &kayakv1.TopicConsumer{
			Id:         c.cfg.ConsumerID,
			Topic:      c.cfg.Topic,
			Group:      c.cfg.ConsumerGroup,
			Partitions: c.cfg.Partitions,
		},
	}))
	if err != nil && errors.Is(err, ErrConsumerAlreadyRegistered) {
		return nil
	}
	return err
}

func (c *Client) FetchRecord(ctx context.Context) (*kayakv1.Record, error) {
	if err := c.validateConfig(); err != nil {
		return nil, err
	}
	if err := c.RegisterConsumer(ctx); err != nil {
		if errors.Is(err, ErrConsumerAlreadyRegistered) {
			slog.Warn("consumer is already registered")
		} else {
			return nil, err
		}

	}
	req := connect.NewRequest(&kayakv1.FetchRecordRequest{
		Consumer: &kayakv1.TopicConsumer{
			Id:    c.cfg.ConsumerID,
			Topic: c.cfg.Topic,
		},
	})
	slog.Info("kayak client fetch", "topic", c.cfg.Topic, "consumer", c.cfg.ConsumerID)
	resp, err := c.client.FetchRecord(ctx, req)
	if err != nil {
		slog.Info("record not found", "error", err)
		if errors.Is(err, service.ErrRecordNotFound) || strings.Contains(err.Error(), "no record found") {
			return nil, nil
		}
		return nil, err
	}
	return resp.Msg.GetRecord(), err
}

func (c *Client) CommitRecord(ctx context.Context, record *kayakv1.Record) error {
	if err := c.validateConfig(); err != nil {
		return err
	}
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
	if topic == "" {
		return nil, ErrTopicEmpty
	}
	req := connect.NewRequest(&kayakv1.GetRecordsRequest{
		Topic: topic,
		Start: start,
		Limit: limit,
	})
	resp, err := c.client.GetRecords(ctx, req)
	return resp.Msg.GetRecords(), err
}

func (c *Client) CreateTopic(ctx context.Context, topicName string, partitions int64) error {
	topic := &kayakv1.Topic{
		Name:       topicName,
		Partitions: partitions,
	}
	_, err := c.client.CreateTopic(ctx, connect.NewRequest(&kayakv1.CreateTopicRequest{
		Topic: topic,
	}))
	return err
}

func (c *Client) DeleteTopic(ctx context.Context, topic string) error {

	req := connect.NewRequest(&kayakv1.DeleteTopicRequest{
		Topic: &kayakv1.Topic{
			Name:       topic,
			Archived:   false,
			Partitions: 1,
		},
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
