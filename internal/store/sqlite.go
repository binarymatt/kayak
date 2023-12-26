package store

import (
	"context"
	"errors"
	"fmt"
	"math/rand"
	"strings"

	"github.com/oklog/ulid/v2"
	"gorm.io/gorm"
	"log/slog"

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
		db:     db,
		idFunc: ulid.Make,
	}
}
func (s *sqlStore) RunMigrations() error {
	return s.db.AutoMigrate(
		&models.Topic{},
		&models.Record{},
		&models.Consumer{},
	)
}

func (s *sqlStore) CreateTopic(ctx context.Context, topic *models.Topic) error {
	if topic.PartitionCount <= 0 {
		topic.PartitionCount = 1
	}
	topic.DefaultHash = kayakv1.Hash_HASH_MURMUR3.String()
	return s.db.Create(topic).Error
}

func (s *sqlStore) DeleteTopic(ctx context.Context, topic *models.Topic) error {
	return s.db.Transaction(func(tx *gorm.DB) error {
		selector := tx.Model(&models.Topic{}).Where("id = ?", topic.ID)
		if topic.Archived {
			if err := selector.Update("archived", true).Error; err != nil {
				fmt.Println("oops, during archive", err)
				return err
			}
		} else {
			if err := selector.Delete(&models.Topic{}).Error; err != nil {
				fmt.Println("oops during delete", err)
				return err
			}
			if err := tx.Delete(&models.Record{}, "topic_id = ?", topic.ID).Error; err != nil {
				return err
			}
		}
		if err := tx.Delete(&models.Consumer{}, "topic_id = ?", topic.ID).Error; err != nil {
			return err
		}
		return nil
	})
}

func (s *sqlStore) AddRecords(ctx context.Context, topicName string, records ...*models.Record) error {
	var topic models.Topic
	result := s.db.Where("id = ?", topicName).First(&topic)
	if result.Error != nil && errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return ErrInvalidTopic
	}
	if topic.Archived {
		return ErrArchivedTopic
	}

	for _, record := range records {
		if record.Partition <= 0 {
			topic.Balance(record)

		}
	}
	result = s.db.Create(records)
	return result.Error
}

func translateHeaders(items map[string]interface{}) map[string]string {
	headers := make(map[string]string, len(items))
	for k, v := range items {
		headers[k] = v.(string)
	}
	return headers
}

func (s *sqlStore) GetRecords(ctx context.Context, topic string, start string, limit int) ([]*models.Record, error) {
	var t models.Topic
	result := s.db.Where("id = ?", topic).First(&t)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, ErrInvalidTopic
		}
		return nil, result.Error
	}
	if t.Archived {
		return nil, ErrArchivedTopic
	}

	// Where topic = ? and id > start LIMIT ?
	// var topic models.Topic
	var records []*models.Record
	result = s.db.Where("topic_id = ? AND id > ? ", t.ID, start).Limit(limit).Find(&records)
	if result.Error != nil {
		return nil, result.Error
	}
	return records, nil
}

func (s *sqlStore) FetchRecord(ctx context.Context, consumer *models.Consumer) (*models.Record, error) {
	topic, err := s.getTopic(consumer.TopicID)
	if err != nil {
		return nil, err
	}
	if topic.Archived {
		return nil, ErrArchivedTopic
	}
	c, err := s.getConsumer(consumer.ID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrInvalidConsumer
		}
		return nil, err
	}
	var record models.Record
	slog.Info("check partitions", "partitions", c.GetPartitions())
	if err := s.db.First(&record, "topic_id = ? and id > ? and partition in (?)", c.TopicID, c.Position, c.GetPartitions()).Error; err != nil {
		return nil, err
	}
	return &record, nil
}

func (s *sqlStore) ListTopics(ctx context.Context) ([]string, error) {
	var topics []models.Topic
	result := s.db.Find(&topics)
	if result.Error != nil {
		return nil, result.Error
	}
	names := make([]string, len(topics))
	for index, topic := range topics {
		names[index] = topic.ID
	}
	return names, nil
}

func (s *sqlStore) getTopic(id string) (*models.Topic, error) {
	var topic models.Topic
	result := s.db.First(&topic, "id = ?", id)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, ErrInvalidTopic
	}
	return &topic, result.Error
}
func (s *sqlStore) getTopicByName(name string) (*models.Topic, error) {
	var topic models.Topic
	result := s.db.First(&topic, "name = ?", name)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, ErrInvalidTopic
	}
	return &topic, result.Error
}

func chooseConsumer(topic *models.Topic, consumers []models.Consumer) models.Consumer {

	return consumers[rand.Int()%topic.PartitionCount]
}

// balance partitions
func (s *sqlStore) balancePartitionAssignments(topic *models.Topic, groupID string) error {
	return s.db.Transaction(func(tx *gorm.DB) error {

		var consumers []models.Consumer
		if err := tx.Find(&consumers, "group_id = ?", groupID).Error; err != nil {
			return err
		}
		ids := make([]string, len(consumers))
		for index, consumer := range consumers {
			ids[index] = consumer.ID
		}
		// clear existing assignments
		// assign partitions to consumers
		for _, partition := range topic.Partitions() {
			consumer := chooseConsumer(topic, consumers)
			consumer.Partitions = append(consumer.Partitions, partition)
		}
		return nil
	})
}

func (s *sqlStore) RegisterConsumer(ctx context.Context, consumer *models.Consumer) (*models.Consumer, error) {

	topic, err := s.getTopic(consumer.TopicID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrInvalidTopic
		}
		return nil, err
	}
	if topic.Archived {
		return nil, ErrArchivedTopic
	}

	slog.Warn("registering consumer", "partitions", consumer.Partitions)
	if err := s.db.Create(consumer).Error; err != nil {
		slog.Error("could not create consumer record", "error", err)
		if strings.Contains(err.Error(), "UNIQUE constraint failed: consumers.id") {
			return consumer, ErrConsumerAlreadyRegistered
		}
		return nil, err
	}
	return consumer, nil
}

func (s *sqlStore) getConsumer(id string) (*models.Consumer, error) {
	var consumer models.Consumer
	result := s.db.First(&consumer, "id = ?", id)
	if result.Error != nil {
		return nil, result.Error
	}
	return &consumer, nil
}
func (s *sqlStore) GetConsumerPosition(ctx context.Context, consumer *models.Consumer) (string, error) {
	c, err := s.getConsumer(consumer.ID)
	if err != nil {
		return "", err
	}
	return c.Position, nil
}

func (s *sqlStore) CommitConsumerPosition(ctx context.Context, consumer *models.Consumer) error {
	if _, err := s.getConsumer(consumer.ID); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrInvalidConsumer
		}
		return err
	}
	if err := s.db.Model(&models.Consumer{}).Where("id = ?", consumer.ID).Update("position", consumer.Position).Error; err != nil {
		return err
	}
	return nil
}

func (s *sqlStore) LoadMeta(ctx context.Context, topicName string) (*models.Topic, error) {
	slog.Debug("loading topic info from store", "topic", topicName)
	var topic models.Topic
	if err := s.db.Preload("Consumers").First(&topic, "id = ?", topicName).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrInvalidTopic
		}
		return nil, err
	}
	return &topic, nil
}
func (s *sqlStore) getGroupConsumers(topic, group string) ([]*models.Consumer, error) {
	var consumers []*models.Consumer
	result := s.db.Find(&consumers, "group = ?", group)
	if result.Error != nil {
		return nil, result.Error
	}
	return consumers, nil
}
func (s *sqlStore) Stats() map[string]*models.Topic {
	var topics []models.Topic
	result := s.db.Find(&topics)
	if result.Error != nil {
		slog.Error("error with topics find query", "error", result.Error)
		return nil
	}
	meta := map[string]*models.Topic{}
	for _, topic := range topics {
		topicMeta, err := s.LoadMeta(context.Background(), topic.ID)
		if err != nil || topicMeta == nil {
			slog.Error("could not retrieve topic metadata", "error", err)
			return nil
		}
		meta[topicMeta.ID] = topicMeta
	}

	return meta
}

func (s *sqlStore) Impl() any {
	return s.db
}

func (s *sqlStore) Close() {}

func (s *sqlStore) SnapshotItems() <-chan DataItem {
	return nil
}

func (s *sqlStore) GetConsumerLag(ctx context.Context, consumer *models.Consumer) (int64, error) {
	topic, err := s.getTopic(consumer.TopicID)
	if err != nil {
		return 0, err
	}
	var count int64
	result := s.db.Model(&models.Record{}).Where("topic_id = ? and id > ?", topic.ID, consumer.Position).Count(&count)
	return count, result.Error
}
