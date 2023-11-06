package store

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"sync"
	"time"

	badger "github.com/dgraph-io/badger/v4"
	"google.golang.org/protobuf/types/known/timestamppb"
	"log/slog"

	kayakv1 "github.com/binarymatt/kayak/gen/kayak/v1"
)

type badgerStore struct {
	db        *badger.DB
	metaMutex sync.Mutex
	timeFunc  func() time.Time
}

func NewBadger(db *badger.DB) *badgerStore {
	slog.Debug("created new db", "db", db)
	return &badgerStore{
		db:       db,
		timeFunc: time.Now,
	}
}

func intToBytes(i int64) []byte {
	return strconv.AppendInt(nil, i, 10)
}

func (b *badgerStore) initTopicMeta(ctx context.Context, tx *badger.Txn, topic string) error {
	prefix := fmt.Sprintf("%s#", topic)
	recordCount := intToBytes(0)
	createdAt := intToBytes(b.timeFunc().UTC().Unix())
	if err := tx.Set(key(prefix+"created_at"), createdAt); err != nil {
		return ErrCreatingTopic
	}
	if err := tx.Set(key("topics#"+topic), createdAt); err != nil {
		slog.Error("could not add to topics list", "error", err)
		return ErrCreatingTopic
	}
	if err := tx.Set(key(prefix+"record_count"), recordCount); err != nil {
		slog.Error("could not create topic record count", "error", err)
		return ErrCreatingTopic
	}
	return nil
}
func (b *badgerStore) CreateTopic(ctx context.Context, name string) error {
	return b.db.Update(func(tx *badger.Txn) error {

		err := b.isTopicArchived(tx, name)
		if err != nil {
			return err
		}

		// create metadata entry
		// add to topics set
		return b.initTopicMeta(ctx, tx, name)
	})
}

func (b *badgerStore) topicExists(tx *badger.Txn, topic string) error {
	_, err := tx.Get(key("topics#" + topic))
	if err != nil {
		return ErrInvalidTopic
	}
	return err
}

func (b *badgerStore) isTopicArchived(tx *badger.Txn, topic string) error {
	k := fmt.Sprintf("%s#archived", topic)
	item, err := tx.Get(key(k))
	if errors.Is(err, badger.ErrKeyNotFound) {
		return nil
	}
	if err != nil {
		return err
	}
	if item.ValueSize() > 0 {
		return ErrTopicArchived
	}
	return nil
}

func (b *badgerStore) topicRecordCount(tx *badger.Txn, topic string) (int64, error) {
	recordCountItem, err := tx.Get(key(topic + "#record_count"))
	if err != nil {
		return -1, err
	}
	rcBytes, err := recordCountItem.ValueCopy(nil)
	if err != nil {
		return -1, err
	}
	return strconv.ParseInt(string(rcBytes), 10, 64)
}
func (b *badgerStore) topicCreatedAt(tx *badger.Txn, topic string) (time.Time, error) {
	ts := time.Time{}
	createdItem, err := tx.Get(key(topic + "#created_at"))
	if err != nil {
		return ts, err
	}
	createdBytes, err := createdItem.ValueCopy(nil)
	if err != nil {
		return ts, err
	}
	raw, err := strconv.ParseInt(string(createdBytes), 10, 64)
	if err != nil {
		return ts, err
	}
	return time.Unix(raw, 0), nil
}

func (b *badgerStore) retrieveTopicMeta(tx *badger.Txn, topic string) (*kayakv1.TopicMetadata, error) {
	// metaKey := key(fmt.Sprintf("topics#%s", topic))
	var meta kayakv1.TopicMetadata
	return &meta, nil
}
func (b *badgerStore) ArchiveTopic(ctx context.Context, topic string) (err error) {
	err = b.db.DropPrefix(key(fmt.Sprintf("%s#", topic)))
	if err != nil {
		return
	}
	err = b.db.Update(func(tx *badger.Txn) error {
		b.metaMutex.Lock()
		defer b.metaMutex.Unlock()
		data := strconv.FormatBool(true)
		key := key(fmt.Sprintf("%s#archived", topic))
		if err := tx.Set(key, []byte(data)); err != nil {
			return err
		}
		return nil
	})
	return err
}

func (b *badgerStore) PurgeTopic(ctx context.Context, topic string) (err error) {
	err = b.db.DropPrefix(key(fmt.Sprintf("%s#", topic)))
	if err != nil {
		return
	}
	return b.db.DropPrefix(key(fmt.Sprintf("topics#%s", topic)))
}

func (b *badgerStore) DeleteTopic(ctx context.Context, topic string, archive bool) (err error) {
	if archive {
		return b.ArchiveTopic(ctx, topic)
	}
	return b.PurgeTopic(ctx, topic)
}

func (b *badgerStore) AddRecords(ctx context.Context, topic string, records ...*kayakv1.Record) error {
	return b.db.Update(func(tx *badger.Txn) error {
		b.metaMutex.Lock()
		defer b.metaMutex.Unlock()
		err := b.isTopicArchived(tx, topic)
		if err != nil {
			return err
		}

		current, err := b.topicRecordCount(tx, topic)
		if err != nil {
			return err
		}

		recordCount := int64(len(records)) + current

		for _, record := range records {
			record.Topic = topic
			recordKey := fmt.Sprintf("%s#records#%s", topic, record.Id)
			data, err := encode(record)
			if err != nil {
				return err
			}
			if err := tx.Set(key(recordKey), data); err != nil {
				return err
			}
		}
		cnt := fmt.Sprintf("%d", recordCount)
		return tx.Set(key(topic+"#record_count"), []byte(cnt))
	})
}

func (b *badgerStore) GetRecords(ctx context.Context, topic string, start string, limit int) (records []*kayakv1.Record, err error) {
	err = b.db.View(func(tx *badger.Txn) error {
		it := tx.NewIterator(badger.DefaultIteratorOptions)
		defer it.Close()

		prefix := []byte(fmt.Sprintf("%s#", topic))
		startKey := []byte(fmt.Sprintf("%s#records#%s", topic, start))

		counter := 0
		for it.Seek(startKey); it.ValidForPrefix(prefix); it.Next() {
			if counter >= limit {
				slog.Debug("past limit, exit")
				break
			}
			item := it.Item()
			data, err := item.ValueCopy(nil)
			if err != nil {
				return err
			}
			record, err := decode(nil, data)
			if err != nil {
				return err
			}
			if record.Id == start {
				slog.Debug("skipping start value", "start", start)
				continue
			}
			records = append(records, record)
			counter++

		}

		return nil
	})
	slog.InfoContext(ctx, "returning records from badger", "count", len(records))
	return
}

func (b *badgerStore) getId(item *badger.Item) string {
	k := string(item.Key())
	key := strings.Split(k, "#")[1]
	return key
}

func (b *badgerStore) ListTopics(ctx context.Context) (topics []string, err error) {
	err = b.db.View(func(tx *badger.Txn) error {
		opts := badger.DefaultIteratorOptions
		it := tx.NewIterator(opts)
		defer it.Close()
		prefix := []byte("topics#")
		for it.Seek(prefix); it.ValidForPrefix(prefix); it.Next() {
			item := it.Item()
			key := b.getId(item)
			topics = append(topics, key)
		}

		return nil
	})
	return
}

func (b *badgerStore) GetConsumerPosition(ctx context.Context, consumer *kayakv1.TopicConsumer) (position string, err error) {

	err = b.db.View(func(tx *badger.Txn) error {
		position, err = b.getConsumerPosition(ctx, tx, consumer)
		return nil
	})
	return
}

func (b *badgerStore) CommitConsumerPosition(ctx context.Context, consumer *kayakv1.TopicConsumer) (err error) {
	err = b.db.Update(func(tx *badger.Txn) error {
		b.metaMutex.Lock()
		defer b.metaMutex.Unlock()
		err := b.topicExists(tx, consumer.Topic)
		if err != nil {
			return err
		}

		err = b.isTopicArchived(tx, consumer.Topic)
		if err != nil {
			return err
		}

		positionKey := fmt.Sprintf("%s#groups#%s#consumer_position#%s", consumer.Topic, consumer.Group, consumer.Id)
		return tx.Set(key(positionKey), bytes(consumer.Position))
	})
	return
}

func (b *badgerStore) consumerGoupExists(ctx context.Context, tx *badger.Txn, topic, group string) (bool, error) {

	partitionCountKey := fmt.Sprintf("%s#groups#%s#partition_count", topic, group)

	existing, err := tx.Get(key(partitionCountKey))

	if err != nil && !errors.Is(err, badger.ErrKeyNotFound) {
		return false, err
	}
	if existing != nil && existing.ValueSize() > 0 {
		return true, nil
	}
	return false, nil
}

func (b *badgerStore) RegisterConsumerGroup(ctx context.Context, group *kayakv1.ConsumerGroup) (err error) {
	err = b.db.Update(func(tx *badger.Txn) error {
		return b.registerConsumerGroup(ctx, tx, group)
	})
	return
}

func (b *badgerStore) registerConsumerGroup(ctx context.Context, tx *badger.Txn, group *kayakv1.ConsumerGroup) (err error) {
	cnt := intToBytes(group.PartitionCount)
	partitionCountKey := fmt.Sprintf("%s#groups#%s#partition_count", group.Topic, group.Name)
	exists, err := b.consumerGoupExists(ctx, tx, group.Topic, group.Name)
	if err != nil {
		return err
	}
	if exists {
		return ErrConsumerGroupExists
	}
	if err := tx.Set(key(partitionCountKey), cnt); err != nil {
		return err
	}
	return nil
}

func (b *badgerStore) Stats() (results map[string]*kayakv1.TopicMetadata) {
	results = map[string]*kayakv1.TopicMetadata{}
	return
}

func (b *badgerStore) Impl() any {
	return b.db
}

func (b *badgerStore) Close() {
	_ = b.db.Close()
}

func (b *badgerStore) SnapshotItems() <-chan DataItem {
	// loop through topics
	return nil
}

func (b *badgerStore) isConsumerRegistered(ctx context.Context, tx *badger.Txn, consumer *kayakv1.TopicConsumer) error {
	prefix := fmt.Sprintf("%s#groups#%s", consumer.Topic, consumer.Group)
	_, err := tx.Get(key(fmt.Sprintf("%s#consumers#%s", prefix, consumer.Id)))
	if errors.Is(err, badger.ErrKeyNotFound) {
		return nil
	}
	if err != nil {
		return err
	}
	return ErrConsumerAlreadyRegistered

}
func (b *badgerStore) RegisterConsumer(ctx context.Context, consumer *kayakv1.TopicConsumer) (*kayakv1.TopicConsumer, error) {
	registeredConsumer := &kayakv1.TopicConsumer{
		Id:        consumer.Id,
		Topic:     consumer.Topic,
		Group:     consumer.Group,
		Partition: 0,
		Position:  "",
	}
	partition := int64(0)
	err := b.db.Update(func(tx *badger.Txn) error {
		prefix := fmt.Sprintf("%s#groups#%s", consumer.Topic, consumer.Group)
		partitionCountKey := key(fmt.Sprintf("%s#partition_count", prefix))
		groupExists, err := b.consumerGoupExists(ctx, tx, consumer.Topic, consumer.Group)
		if err != nil {
			return err
		}
		if !groupExists {
			if err := b.registerConsumerGroup(ctx, tx, &kayakv1.ConsumerGroup{
				Name:           consumer.Group,
				Topic:          consumer.Topic,
				PartitionCount: int64(1),
			}); err != nil {
				return err
			}
		}
		err = b.isConsumerRegistered(ctx, tx, consumer)
		if err != nil {
			return err
		}

		item, err := tx.Get(partitionCountKey)
		if err != nil {
			return err
		}
		raw, err := item.ValueCopy(nil)
		if err != nil {
			return err
		}
		cnt, err := strconv.Atoi(string(raw))
		if err != nil {
			return err
		}
		items, err := b.getConsumerPartitions(ctx, tx, consumer.Group, consumer.Topic)
		if err != nil {
			return err
		}
		nextPartition := len(items)
		if nextPartition >= cnt {
			return ErrGroupFull
		}
		partition = int64(nextPartition)
		// setting consumer position at empty (start at begining)
		if err := tx.Set(key(fmt.Sprintf("%s#consumers#%s", prefix, consumer.Id)), key(fmt.Sprintf("%d", nextPartition))); err != nil {
			return err
		}
		if err := tx.Set(key(fmt.Sprintf("%s#consumer_position#%s", prefix, consumer.Id)), key("")); err != nil {
			return err
		}
		return nil
	})
	registeredConsumer.Partition = partition
	return registeredConsumer, err
}
func (b *badgerStore) getConsumerGroupNames(tx *badger.Txn, topic string) (names []string, err error) {
	s := map[string]bool{}
	it := tx.NewIterator(badger.DefaultIteratorOptions)
	defer it.Close()
	prefix := key(fmt.Sprintf("%s#groups#", topic))
	for it.Seek(prefix); it.ValidForPrefix(prefix); it.Next() {
		key := string(it.Item().Key())
		items := strings.Split(key, "#")
		s[items[2]] = true
	}
	for k := range s {
		names = append(names, k)
	}
	return
}

// GetConsumerPartitions returns the current partition assignments for a consumer group
func (b *badgerStore) GetConsumerPartitions(ctx context.Context, topic, group string) (items []*kayakv1.TopicConsumer, err error) {
	err = b.db.View(func(tx *badger.Txn) error {
		var e error
		items, e = b.getConsumerPartitions(ctx, tx, group, topic)
		return e
	})
	return items, err
}
func (b *badgerStore) getConsumerPartitions(ctx context.Context, tx *badger.Txn, group, topic string) ([]*kayakv1.TopicConsumer, error) {
	prefix := fmt.Sprintf("%s#groups#%s", topic, group)
	items := []*kayakv1.TopicConsumer{}
	// partitionCountKey := key(fmt.Sprintf("%s#partition_count", prefix))

	consumersKeyPrefix := key(fmt.Sprintf("%s#consumers", prefix))
	it := tx.NewIterator(badger.DefaultIteratorOptions)
	defer it.Close()

	for it.Seek(consumersKeyPrefix); it.ValidForPrefix(consumersKeyPrefix); it.Next() {
		item := it.Item()
		consumerID := strings.Split(string(item.Key()), "#")[4]
		raw, err := item.ValueCopy(nil)
		if err != nil {
			return nil, err
		}
		partition, err := strconv.ParseInt(string(raw), 10, 64)
		if err != nil {
			return nil, err
		}
		position, err := b.getConsumerPosition(ctx, tx, &kayakv1.TopicConsumer{Topic: topic, Group: group, Id: consumerID})
		if err != nil {
			return nil, err
		}
		items = append(items, &kayakv1.TopicConsumer{
			Id:        consumerID,
			Topic:     topic,
			Group:     group,
			Partition: partition,
			Position:  position,
		})
	}
	return items, nil
}

func (b *badgerStore) getConsumerPosition(ctx context.Context, tx *badger.Txn, consumer *kayakv1.TopicConsumer) (string, error) {
	positionKey := fmt.Sprintf("%s#groups#%s#consumer_position#%s", consumer.Topic, consumer.Group, consumer.Id)
	item, err := tx.Get(key(positionKey))
	if errors.Is(err, badger.ErrKeyNotFound) {
		return "", ErrInvalidConsumer
	}
	if err != nil {
		return "", err
	}
	data, err := item.ValueCopy(nil)
	if err != nil {
		return "", err
	}
	return string(data), nil
}

func (b *badgerStore) loadMeta(ctx context.Context, tx *badger.Txn, topic string) (*kayakv1.TopicMetadata, error) {
	groupMeta := map[string]*kayakv1.GroupPartitions{}
	groups, err := b.getConsumerGroupNames(tx, topic)
	if err != nil {
		return nil, err
	}
	for _, name := range groups {
		partitions, err := b.getConsumerPartitions(ctx, tx, name, topic)
		if err != nil {
			return nil, err
		}
		partitionCountKey := fmt.Sprintf("%s#groups#%s#partition_count", topic, name)

		data, err := tx.Get(key(partitionCountKey))
		if err != nil {
			return nil, err
		}
		raw, err := data.ValueCopy(nil)
		if err != nil {
			return nil, err
		}
		cnt, err := strconv.ParseInt(string(raw), 10, 64)
		if err != nil {
			return nil, err
		}
		groupMeta[name] = &kayakv1.GroupPartitions{
			Name:       name,
			Partitions: cnt,
			Consumers:  partitions,
		}
	}
	// get record count
	recordCount, err := b.topicRecordCount(tx, topic)
	if err != nil {
		return nil, err
	}
	// get created at
	ts, err := b.topicCreatedAt(tx, topic)
	if err != nil {
		return nil, err
	}
	// get archived
	err = b.isTopicArchived(tx, topic)
	if err != nil && !errors.Is(err, ErrTopicArchived) {
		return nil, err
	}

	meta := &kayakv1.TopicMetadata{
		Name:          topic,
		RecordCount:   recordCount,
		CreatedAt:     timestamppb.New(ts),
		Archived:      errors.Is(err, ErrTopicArchived),
		GroupMetadata: groupMeta,
	}
	return meta, nil
}
