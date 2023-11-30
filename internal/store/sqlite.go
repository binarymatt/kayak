package store

import (
	"context"
	"errors"
	"strings"
	"time"

	"github.com/oklog/ulid/v2"
	"google.golang.org/protobuf/types/known/timestamppb"
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
		&models.ConsumerGroup{},
		&models.Consumer{},
	)
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

func (s *sqlStore) getTopic(topicName string) (*models.Topic, error) {
	var topic models.Topic
	result := s.db.First(&topic, "name = ?", topicName)
	return &topic, result.Error
}

func (s *sqlStore) getConsumerGroup(topicName, groupName string) (*models.ConsumerGroup, error) {
	topic, err := s.getTopic(topicName)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrInvalidTopic
		}
		return nil, err
	}
	var cg models.ConsumerGroup
	result := s.db.First(&cg, "topic_id = ? AND name = ?", topic.ID, groupName)
	if result.RowsAffected < 1 {
		return nil, ErrConsumerGroupInvalid
	}
	return &cg, nil
}

func (s *sqlStore) RegisterConsumerGroup(ctx context.Context, group *kayakv1.ConsumerGroup) error {
	topic, err := s.getTopic(group.Topic)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrInvalidTopic
		}
		return err
	}
	if topic.Archived {
		return ErrTopicArchived
	}
	cg := &models.ConsumerGroup{
		Name:           group.Name,
		TopicID:        topic.ID,
		PartitionCount: group.PartitionCount,
		Hash:           group.Hash.String(),
	}
	result := s.db.Create(cg)
	if result.Error != nil {
		if strings.Contains(result.Error.Error(), "UNIQUE constraint failed: consumer_groups.name") {
			return ErrConsumerGroupExists
		}
	}
	return result.Error
}

func (s *sqlStore) RegisterConsumer(ctx context.Context, consumer *kayakv1.TopicConsumer) (*kayakv1.TopicConsumer, error) {
	group, err := s.getConsumerGroup(consumer.Topic, consumer.Group)
	if err != nil {
		return nil, err
	}
	result := s.db.Find(&models.Consumer{}, "group_id = ?", group.ID)
	if result.Error != nil {
		return nil, result.Error
	}
	if result.RowsAffected >= group.PartitionCount {
		return nil, ErrGroupFull
	}
	c := models.Consumer{
		ID:        consumer.Id,
		GroupID:   group.ID,
		Partition: result.RowsAffected,
	}
	result = s.db.Create(&c)
	if result.Error != nil {
		return nil, result.Error
	}
	return &kayakv1.TopicConsumer{
		Id:        consumer.Id,
		Topic:     consumer.Topic,
		Group:     consumer.Group,
		Partition: c.Partition,
	}, nil
}

func (s *sqlStore) GetConsumerPartitions(ctx context.Context, topic, group string) ([]*kayakv1.TopicConsumer, error) {
	// 1. get group id
	consumerGroup, err := s.getConsumerGroup(topic, group)
	if err != nil {
		return nil, err
	}
	var consumers []models.Consumer
	result := s.db.Find(&consumers, "group_id = ?", consumerGroup.ID)
	if result.Error != nil {
		return nil, result.Error
	}
	topicConsumers := make([]*kayakv1.TopicConsumer, len(consumers))
	// 2. get all consumer information for group
	for index, consumer := range consumers {
		topicConsumers[index] = &kayakv1.TopicConsumer{
			Id:        consumer.ID,
			Group:     group,
			Topic:     topic,
			Partition: consumer.Partition,
			Position:  consumer.Position,
		}
	}
	return topicConsumers, nil
}

func (s *sqlStore) getConsumer(id string) (*models.Consumer, error) {
	var consumer models.Consumer
	result := s.db.First(&consumer, "id = ?", id)
	if result.Error != nil {
		return nil, result.Error
	}
	return &consumer, nil
}
func (s *sqlStore) GetConsumerPosition(ctx context.Context, consumer *kayakv1.TopicConsumer) (string, error) {
	c, err := s.getConsumer(consumer.Id)
	if err != nil {
		return "", err
	}
	return c.Position, nil
}

func (s *sqlStore) CommitConsumerPosition(ctx context.Context, consumer *kayakv1.TopicConsumer) error {
	result := s.db.First(&models.Consumer{}, "id = ?", consumer.Id)
	if result.RowsAffected == 0 {
		return ErrInvalidConsumer
	}
	result = s.db.Model(&models.Consumer{}).Where("id = ?", consumer.Id).Update("position", consumer.Position)
	return result.Error
}

func (s *sqlStore) GetConsumerGroupInformation(ctx context.Context, group *kayakv1.ConsumerGroup) (*kayakv1.ConsumerGroup, error) {
	topic, err := s.getTopic(group.Topic)
	if err != nil {
		return nil, err
	}
	consumerGroup, err := s.getConsumerGroup(group.Topic, group.Name)
	if err != nil {
		return nil, err
	}
	return &kayakv1.ConsumerGroup{
		Name:           consumerGroup.Name,
		Topic:          topic.Name,
		PartitionCount: consumerGroup.PartitionCount,
		Hash:           kayakv1.Hash_HASH_MURMUR3,
	}, nil
}

func (s *sqlStore) LoadMeta(ctx context.Context, topicName string) (*kayakv1.TopicMetadata, error) {
	slog.Info("loading topic meta from store", "topic", topicName)
	groupMeta := map[string]*kayakv1.GroupPartitions{}
	topic, err := s.getTopic(topicName)
	if err != nil {
		return nil, err
	}
	var groups []models.ConsumerGroup
	result := s.db.Find(&groups, "topic_id = ?", topic.ID)
	if result.Error != nil {
		return nil, result.Error
	}
	for _, group := range groups {
		consumers, err := s.getGroupConsumers(topicName, group.Name, group.ID)
		if err != nil {
			return nil, err
		}
		groupMeta[group.Name] = &kayakv1.GroupPartitions{
			Name:       group.Name,
			Partitions: group.PartitionCount,
			Consumers:  consumers,
		}
	}

	var count int64
	result = s.db.Model(&models.Record{}).Where("topic_id = ?", topic.ID).Count(&count)
	if result.Error != nil {
		return nil, result.Error
	}
	ts := time.Unix(topic.CreatedAt, 0)
	meta := &kayakv1.TopicMetadata{
		Name:          topicName,
		RecordCount:   count,
		CreatedAt:     timestamppb.New(ts),
		Archived:      topic.Archived,
		GroupMetadata: groupMeta,
	}
	return meta, nil
}
func (s *sqlStore) getGroupConsumers(topic, group string, groupID uint) ([]*kayakv1.TopicConsumer, error) {
	var consumers []models.Consumer
	result := s.db.Find(&consumers, "group_id = ?", groupID)
	if result.Error != nil {
		return nil, result.Error
	}
	topicConsumers := make([]*kayakv1.TopicConsumer, len(consumers))
	for index, consumer := range consumers {
		topicConsumers[index] = &kayakv1.TopicConsumer{
			Id:        consumer.ID,
			Topic:     topic,
			Group:     group,
			Partition: consumer.Partition,
			Position:  consumer.Position,
		}
	}
	return topicConsumers, nil
}
func (s *sqlStore) Stats() map[string]*kayakv1.TopicMetadata {
	var topics []models.Topic
	result := s.db.Find(&topics)
	if result.Error != nil {
		slog.Error("error with topics find query", "error", result.Error)
		return nil
	}
	meta := map[string]*kayakv1.TopicMetadata{}
	for _, topic := range topics {
		topicMeta, err := s.LoadMeta(context.Background(), topic.Name)
		if err != nil || topicMeta == nil {
			slog.Error("could not retrieve topic metadata", "error", err)
			return nil
		}
		meta[topicMeta.Name] = topicMeta
	}

	return meta
}

func (s *sqlStore) Impl() any {
	return nil
}

func (s *sqlStore) Close() {}

func (s *sqlStore) SnapshotItems() <-chan DataItem {
	return nil
}

func (s *sqlStore) GetConsumerLag(ctx context.Context, consumer *kayakv1.TopicConsumer) (int64, error) {
	topic, err := s.getTopic(consumer.Topic)
	if err != nil {
		return 0, err
	}
	var count int64
	result := s.db.Model(&models.Record{}).Where("topic_id = ? and id > ?", topic.ID, consumer.Position).Count(&count)
	return count, result.Error
}
