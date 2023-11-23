package store

import (
	"context"
	"errors"

	"github.com/oklog/ulid/v2"
	"gorm.io/gorm"

	kayakv1 "github.com/binarymatt/kayak/gen/kayak/v1"
	"github.com/binarymatt/kayak/internal/store/models"
)

type sqlStore struct {
	db     *gorm.DB
	idFunc func() ulid.ULID
	// timeFunc func() time.Time
}

func NewSqlStore(db *gorm.DB) *sqlStore {
	return &sqlStore{
		db: db,
	}
}
func (s *sqlStore) RunMigrations() error {
	return s.db.AutoMigrate(&models.Topic{}, &models.Record{})
}

func (s *sqlStore) CreateTopic(ctx context.Context, name string) error {
	topic := &models.Topic{
		ID:   s.idFunc().String(),
		Name: name,
	}
	result := s.db.Create(topic)
	return result.Error
}
func (s *sqlStore) DeleteTopic(ctx context.Context, topic string, force bool) error {
	tx := s.db.Model(&models.Topic{}).Where("name = ?", topic)
	if force {
		result := tx.Delete(&models.Topic{})
		return result.Error
	}
	result := tx.Update("archived", true)
	return result.Error

}

func (s *sqlStore) AddRecords(ctx context.Context, name string, records ...*kayakv1.Record) error {
	var topic models.Topic
	result := s.db.Where("name = ? ", name).First(&topic)
	if result.Error != nil && errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return ErrInvalidTopic
	}
	if topic.Archived {
		return ErrTopicArchived
	}

	items := make([]*models.Record, len(records))
	for index, record := range records {

		items[index] = &models.Record{
			TopicID: topic.ID,
			ID:      record.Id,
			Payload: record.Payload,
		}
		items[index].AddHeaders(record.Headers)
	}
	result = s.db.Create(items)
	return result.Error
}
func translateHeaders(items map[string]interface{}) map[string]string {
	headers := make(map[string]string, len(items))
	for k, v := range items {
		headers[k] = v.(string)
	}
	return headers
}
func (s *sqlStore) GetRecords(ctx context.Context, topic string, start string, limit int) ([]*kayakv1.Record, error) {
	var t models.Topic
	result := s.db.Where("name = ?", topic).First(&t)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, ErrInvalidTopic
		}
		return nil, result.Error
	}
	if t.Archived {
		return nil, ErrTopicArchived
	}

	// Where topic = ? and id > start LIMIT ?
	// var topic models.Topic
	var items []models.Record
	result = s.db.Where("topic_id = ? AND id > ? ", t.ID, start).Limit(limit).Find(&items)
	if result.Error != nil {
		return nil, result.Error
	}
	records := make([]*kayakv1.Record, result.RowsAffected)
	for index, item := range items {
		records[index] = &kayakv1.Record{
			Topic:   topic,
			Id:      item.ID,
			Headers: translateHeaders(item.Headers),
			Payload: item.Payload,
		}

	}
	return records, nil
}

func (s *sqlStore) ListTopics(ctx context.Context) ([]string, error) {
	var topics []models.Topic
	result := s.db.Find(&topics)
	if result.Error != nil {
		return nil, result.Error
	}
	names := make([]string, len(topics))
	for index, topic := range topics {
		names[index] = topic.Name
	}
	return names, nil
}

func (s *sqlStore) GetConsumerPartitions(ctx context.Context, topic, group string) ([]*kayakv1.TopicConsumer, error) {
	return nil, nil
}

func (s *sqlStore) RegisterConsumerGroup(ctx context.Context, group *kayakv1.ConsumerGroup) error {
	return nil
}

func (s *sqlStore) RegisterConsumer(ctx context.Context, consumer *kayakv1.TopicConsumer) (*kayakv1.TopicConsumer, error) {
	return nil, nil
}

func (s *sqlStore) GetConsumerPosition(ctx context.Context, consumer *kayakv1.TopicConsumer) (string, error) {
	return "", nil
}

func (s *sqlStore) CommitConsumerPosition(ctx context.Context, consumer *kayakv1.TopicConsumer) error {
	return nil
}

func (s *sqlStore) GetConsumerGroupInformation(ctx context.Context, group *kayakv1.ConsumerGroup) (*kayakv1.ConsumerGroup, error) {
	return nil, nil
}

func (s *sqlStore) LoadMeta(ctx context.Context, topic string) (*kayakv1.TopicMetadata, error) {
	return nil, nil
}

func (s *sqlStore) Stats() map[string]*kayakv1.TopicMetadata {
	return nil
}

func (s *sqlStore) Impl() any {
	return nil
}

func (s *sqlStore) Close() {}

func (s *sqlStore) SnapshotItems() <-chan DataItem {
	return nil
}
