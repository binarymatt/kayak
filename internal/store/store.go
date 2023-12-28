package store

import (
	"context"

	"github.com/binarymatt/kayak/internal/store/models"
)

type Store interface {
	CreateTopic(ctx context.Context, topic *models.Topic) error
	DeleteTopic(ctx context.Context, topic *models.Topic) error
	AddRecords(ctx context.Context, topic string, records ...*models.Record) error
	GetRecords(ctx context.Context, topic string, start string, limit int) ([]*models.Record, error)
	FetchRecord(ctx context.Context, consumer *models.Consumer) (*models.Record, error)
	ListTopics(ctx context.Context) ([]string, error)
	RegisterConsumer(ctx context.Context, consumer *models.Consumer) (*models.Consumer, error)
	GetConsumerPosition(ctx context.Context, consumer *models.Consumer) (string, error)
	CommitConsumerPosition(ctx context.Context, consumer *models.Consumer) error
	GetConsumerLag(ctx context.Context, consumer *models.Consumer) (int64, error)
	LoadMeta(ctx context.Context, topic string) (*models.Topic, error)
	Stats() map[string]*models.Topic
	Impl() any
	Close()
	SnapshotItems() <-chan DataItem
	PruneOldRecords(ctx context.Context) error
}

type DataItem interface{}
