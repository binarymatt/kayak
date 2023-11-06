package store

import (
	"context"

	kayakv1 "github.com/binarymatt/kayak/gen/kayak/v1"
)

type Store interface {
	CreateTopic(ctx context.Context, name string) error
	DeleteTopic(ctx context.Context, topic string, force bool) error
	AddRecords(ctx context.Context, topic string, records ...*kayakv1.Record) error
	GetRecords(ctx context.Context, topic string, start string, limit int) ([]*kayakv1.Record, error)
	ListTopics(ctx context.Context) ([]string, error)
	GetConsumerPartitions(ctx context.Context, topic, group string) ([]*kayakv1.TopicConsumer, error)
	RegisterConsumerGroup(ctx context.Context, group *kayakv1.ConsumerGroup) error
	RegisterConsumer(ctx context.Context, consumer *kayakv1.TopicConsumer) (*kayakv1.TopicConsumer, error)
	GetConsumerPosition(ctx context.Context, consumer *kayakv1.TopicConsumer) (string, error)
	CommitConsumerPosition(ctx context.Context, consumer *kayakv1.TopicConsumer) error
	Stats() map[string]*kayakv1.TopicMetadata
	Impl() any
	Close()
	SnapshotItems() <-chan DataItem
}

type DataItem interface{}
