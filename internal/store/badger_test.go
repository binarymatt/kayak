package store

import (
	"context"
	"fmt"
	"io"
	"strconv"
	"testing"
	"time"

	"github.com/dgraph-io/badger/v4"
	"github.com/google/go-cmp/cmp"
	"github.com/oklog/ulid/v2"
	"github.com/stretchr/testify/suite"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/testing/protocmp"
	"google.golang.org/protobuf/types/known/timestamppb"
	"log/slog"

	kayakv1 "github.com/binarymatt/kayak/gen/kayak/v1"
)

const (
	testTopicName = "test"
)

func createRecords() []*kayakv1.Record {
	return []*kayakv1.Record{
		{Topic: "test", Id: ulid.Make().String(), Payload: []byte("sample message 1")},
		{Topic: "test", Id: ulid.Make().String(), Payload: []byte("sample message 2")},
		{Topic: "test", Id: ulid.Make().String(), Payload: []byte("sample message 3")},
	}
}

type BadgerTestSuite struct {
	suite.Suite
	store *badgerStore
	db    *badger.DB
	ctx   context.Context
	ts    time.Time
}

func (b *BadgerTestSuite) SetupSuite() {
	l := slog.New(slog.NewTextHandler(io.Discard, nil))
	slog.SetDefault(l)
}
func (b *BadgerTestSuite) SetupTest() {
	opt := badger.DefaultOptions("").WithInMemory(true).WithLoggingLevel(badger.ERROR)
	db, err := badger.Open(opt)
	b.NoError(err)
	b.db = db
	b.ctx = context.Background()
	now := time.Now()
	b.ts = now
	b.store = &badgerStore{
		db:       db,
		timeFunc: func() time.Time { return now },
	}
}

func (b *BadgerTestSuite) TestCreateTopic() {
	err := b.store.CreateTopic(b.ctx, testTopicName)
	b.NoError(err)

	err = b.db.View(func(tx *badger.Txn) error {
		_, err := tx.Get(key("topics#test"))
		b.NoError(err)

		return nil
	})
	b.NoError(err)
}

func (b *BadgerTestSuite) TestAddRecords_HappyPath() {
	err := b.store.CreateTopic(b.ctx, testTopicName)
	b.NoError(err)

	records := createRecords()

	err = b.store.AddRecords(context.Background(), testTopicName, records...)
	b.NoError(err)

	err = b.db.View(func(tx *badger.Txn) error {
		// Validate metadata
		meta, err := b.store.loadMeta(b.ctx, tx, "test")
		b.NoError(err)
		if err != nil {
			return err
		}
		b.Equal(int64(3), meta.RecordCount)
		// Validate records exist
		it := tx.NewIterator(badger.DefaultIteratorOptions)
		defer it.Close()
		prefix := []byte("test#records")
		i := 0
		for it.Seek(prefix); it.ValidForPrefix(prefix); it.Next() {
			item := it.Item()
			data, err := item.ValueCopy(nil)
			b.NoError(err)
			if err != nil {
				return err
			}

			fmt.Println("testing decode")
			r, err := decode(nil, data)
			b.NoError(err)
			if err != nil {
				return err
			}
			b.NotNil(r)

			fmt.Println("checking id", r)
			b.NotEmpty(r.Id)
			b.Equal(records[i].Topic, r.Topic)
			b.Equal(records[i].Payload, r.Payload)
			i++
		}

		return nil
	})
	b.NoError(err)

}
func (b *BadgerTestSuite) TestGetRecords_AllRecords() {
	ctx := context.Background()
	_ = b.store.CreateTopic(ctx, "test")
	records := createRecords()
	err := b.store.AddRecords(ctx, "test", records...)
	b.NoError(err)
	items, err := b.store.GetRecords(ctx, "test", "", 100)
	b.NoError(err)
	b.Len(items, 3)
}

func (b *BadgerTestSuite) TestGetRecords_TwoRecords() {
	ctx := context.Background()
	_ = b.store.CreateTopic(ctx, "test")
	records := createRecords()
	err := b.store.AddRecords(ctx, "test", records...)
	b.NoError(err)
	items, err := b.store.GetRecords(ctx, "test", "", 2)
	b.NoError(err)
	b.Len(items, 2)
	b.Empty(cmp.Diff(records[0], items[0], protocmp.Transform()))
	b.Empty(cmp.Diff(records[1], items[1], protocmp.Transform()))
}

func (b *BadgerTestSuite) TestGetRecords_Single() {
	ctx := context.Background()
	_ = b.store.CreateTopic(ctx, "test")
	records := createRecords()
	err := b.store.AddRecords(ctx, "test", records...)
	b.NoError(err)

	record, err := b.store.GetRecords(ctx, "test", "", 1)
	b.NoError(err)
	b.Len(record, 1)
	b.Empty(cmp.Diff(records[0], record[0], protocmp.Transform()))

	record2, err := b.store.GetRecords(ctx, "test", records[0].Id, 1)
	b.NoError(err)
	b.Empty(cmp.Diff(records[1], record2[0], protocmp.Transform()))

	record3, err := b.store.GetRecords(ctx, "test", records[1].Id, 1)
	b.NoError(err)
	b.Empty(cmp.Diff(records[2], record3[0], protocmp.Transform()))
}

func (b *BadgerTestSuite) TestGetRecords_WithStart() {
	ctx := context.Background()
	_ = b.store.CreateTopic(ctx, "test")
	records := createRecords()
	err := b.store.AddRecords(ctx, "test", records...)
	b.NoError(err)
	items, err := b.store.GetRecords(ctx, "test", records[0].Id, 100)
	b.NoError(err)
	b.Len(items, 2)
	b.Empty(cmp.Diff(records[1], items[0], protocmp.Transform()))
	b.Empty(cmp.Diff(records[2], items[1], protocmp.Transform()))
}

func (b *BadgerTestSuite) TestCommitConsumerPosition() {
	ctx := context.Background()
	err := b.store.CreateTopic(ctx, testTopicName)
	b.NoError(err)

	err = b.store.CommitConsumerPosition(ctx, &kayakv1.TopicConsumer{Topic: testTopicName, Group: "testGroup", Id: "1", Position: "recordID"})
	b.NoError(err)
	err = b.db.View(func(tx *badger.Txn) error {
		item, err := tx.Get(key("topics#test"))
		b.NoError(err)

		data, err := item.ValueCopy(nil)
		b.NoError(err)
		var meta kayakv1.TopicMetadata
		err = proto.Unmarshal(data, &meta)
		b.NoError(err)

		return nil
	})
	b.NoError(err)
}

func (b *BadgerTestSuite) TestCommitConsumerPosition_ArchivedTopic() {
	ctx := context.Background()
	err := b.store.CreateTopic(ctx, testTopicName)
	b.NoError(err)

	err = b.store.DeleteTopic(ctx, testTopicName, true)
	b.NoError(err)

	err = b.store.CommitConsumerPosition(ctx, &kayakv1.TopicConsumer{Topic: testTopicName, Group: "group", Id: "testConsumer", Position: "1"})
	b.ErrorIs(err, ErrTopicArchived)
}

func (b *BadgerTestSuite) TestCommitConsumerPosition_NonExistentTopic() {
	err := b.store.CommitConsumerPosition(context.Background(), &kayakv1.TopicConsumer{Topic: testTopicName, Group: "group", Id: "testConsumer", Position: "1"})
	b.ErrorIs(err, ErrInvalidTopic)
}

func (b *BadgerTestSuite) TestGetConsumerPosition_NonExistentTopic() {
}

func (b *BadgerTestSuite) TestGetConsumerPosition_NonExistentConsumer() {
}

func (b *BadgerTestSuite) TestGetConsumerPosition_Happy() {

}

func (b *BadgerTestSuite) TestListTopics() {
	topics, err := b.store.ListTopics(context.Background())
	b.NoError(err)
	b.Empty(topics)
	err = b.store.CreateTopic(context.Background(), testTopicName)
	b.NoError(err)
	topics, err = b.store.ListTopics(context.Background())
	b.NoError(err)
	b.Equal([]string{testTopicName}, topics)

}

func (b *BadgerTestSuite) TestRegisterConsumerGroup_HappyPath() {
	err := b.store.RegisterConsumerGroup(
		context.Background(),
		&kayakv1.ConsumerGroup{
			Name:           "testName",
			Topic:          "testTopic",
			PartitionCount: 1,
		},
	)
	b.NoError(err)

	// Check for correct records in db
	err = b.db.View(func(tx *badger.Txn) error {
		partitionCountKey := fmt.Sprintf("%s#groups#%s#partition_count", "testTopic", "testName")
		item, err := tx.Get(key(partitionCountKey))
		if err != nil {
			return err
		}
		data, err := item.ValueCopy(nil)
		if err != nil {
			return err
		}
		cnt, err := strconv.ParseInt(string(data), 10, 64)
		if err != nil {
			return err
		}
		b.Equal(int64(1), cnt)
		return nil
	})
	b.NoError(err)

}
func (b *BadgerTestSuite) TestRegisterConsumerGroup_ExistingGroup() {
	group := &kayakv1.ConsumerGroup{
		Name:           "groupOne",
		Topic:          "testTopic",
		PartitionCount: 10,
	}
	err := b.store.RegisterConsumerGroup(
		context.Background(),
		group,
	)
	b.NoError(err)
	err = b.store.RegisterConsumerGroup(
		context.Background(),
		group,
	)
	b.ErrorIs(err, ErrConsumerGroupExists)
}
func (b *BadgerTestSuite) TestRegisterConsumer_EmptyGroup() {
	err := b.store.RegisterConsumerGroup(
		context.Background(),
		&kayakv1.ConsumerGroup{
			Name:           "group",
			Topic:          "topic",
			PartitionCount: 10,
		},
	)
	b.NoError(err)

	consumer, err := b.store.RegisterConsumer(context.Background(), &kayakv1.TopicConsumer{
		Topic: "topic",
		Group: "group",
		Id:    "id",
	})
	b.NoError(err)
	b.Equal(int64(0), consumer.Partition)

	err = b.db.Update(func(tx *badger.Txn) error {
		// validate parition
		partitionKey, err := tx.Get(key("topic#groups#group#consumers#id"))
		b.NoError(err)
		data, err := partitionKey.ValueCopy(nil)
		if err != nil {
			return err
		}
		b.Equal("0", string(data))
		// validate position
		positionKey, err := tx.Get(key("topic#groups#group#consumer_position#id"))
		b.NoError(err)
		if err != nil {
			return err
		}
		posData, err := positionKey.ValueCopy(nil)
		b.NoError(err)
		if err != nil {
			return err
		}
		b.Equal("", string(posData))
		return nil
	})
	b.NoError(err)
}

func (b *BadgerTestSuite) TestRegisterConsumer_NonEmptyGroup() {
	err := b.store.RegisterConsumerGroup(
		context.Background(),
		&kayakv1.ConsumerGroup{
			Name:           "group",
			Topic:          "topic",
			PartitionCount: 2,
		},
	)
	b.NoError(err)

	consumer, err := b.store.RegisterConsumer(context.Background(), &kayakv1.TopicConsumer{
		Topic: "topic",
		Group: "group",
		Id:    "one",
	})
	b.NoError(err)
	b.Equal(int64(0), consumer.Partition)

	consumer, err = b.store.RegisterConsumer(context.Background(), &kayakv1.TopicConsumer{
		Topic: "topic",
		Group: "group",
		Id:    "two",
	})
	b.NoError(err)
	b.Equal(int64(1), consumer.Partition)
}
func (b *BadgerTestSuite) TestRegisterConsumer_AlreadyRegistered() {
	err := b.store.RegisterConsumerGroup(
		context.Background(),
		&kayakv1.ConsumerGroup{
			Name:           "group",
			Topic:          "topic",
			PartitionCount: 2,
		},
	)
	b.NoError(err)

	consumer, err := b.store.RegisterConsumer(context.Background(), &kayakv1.TopicConsumer{
		Topic: "topic",
		Group: "group",
		Id:    "one",
	})
	b.NoError(err)
	b.Equal(int64(0), consumer.Partition)

	consumer, err = b.store.RegisterConsumer(context.Background(), &kayakv1.TopicConsumer{
		Topic: "topic",
		Group: "group",
		Id:    "one",
	})
	b.ErrorIs(err, ErrConsumerAlreadyRegistered)
	b.Equal(int64(0), consumer.Partition)
}

func (b *BadgerTestSuite) TestRegisterConsumer_FullGroup() {
	err := b.store.RegisterConsumerGroup(
		context.Background(),
		&kayakv1.ConsumerGroup{
			Name:           "group",
			Topic:          "topic",
			PartitionCount: 1,
		},
	)
	b.NoError(err)

	consumer, err := b.store.RegisterConsumer(context.Background(), &kayakv1.TopicConsumer{
		Topic: "topic",
		Group: "group",
		Id:    "id",
	})
	b.NoError(err)
	b.Equal(int64(0), consumer.Partition)
	_, err = b.store.RegisterConsumer(context.Background(), &kayakv1.TopicConsumer{
		Topic: "topic",
		Group: "group",
		Id:    "id2",
	})
	b.ErrorIs(err, ErrGroupFull)
}
func (b *BadgerTestSuite) TestRegisterConsumer_NonExistantGroup() {
	err := b.store.CreateTopic(b.ctx, "topic")
	b.NoError(err)
	consumer, err := b.store.RegisterConsumer(context.Background(), &kayakv1.TopicConsumer{
		Topic: "topic",
		Group: "group",
		Id:    "id",
	})
	b.NoError(err)
	b.Equal(int64(0), consumer.Partition)
}
func (b *BadgerTestSuite) TestRegisterConsumer_NoSlotsLeft() {

	err := b.store.CreateTopic(b.ctx, "topic")
	b.NoError(err)
	err = b.store.RegisterConsumerGroup(b.ctx, &kayakv1.ConsumerGroup{
		Name:           "testGroup",
		Topic:          "topic",
		PartitionCount: 1,
	})
	b.NoError(err)

	consumer, err := b.store.RegisterConsumer(b.ctx, &kayakv1.TopicConsumer{
		Topic: "topic",
		Group: "testGroup",
		Id:    "consumerID",
	})
	b.NoError(err)
	b.Equal(int64(0), consumer.Partition)

	_, err = b.store.RegisterConsumer(b.ctx, &kayakv1.TopicConsumer{
		Topic: "topic",
		Group: "testGroup",
		Id:    "second",
	})
	b.ErrorIs(err, ErrGroupFull)

}
func (b *BadgerTestSuite) TestRegisterConsumer_MultiplePartitions() {

}

func (b *BadgerTestSuite) TestGetConsumerGroupNames() {
	err := b.db.Update(func(tx *badger.Txn) error {
		if err := tx.Set(key("test#groups#groupOne#partition_count"), key("2")); err != nil {
			return err
		}
		return tx.Set(key("test#groups#groupTwo#partition_count"), key("3"))
	})
	b.NoError(err)
	err = b.db.View(func(tx *badger.Txn) error {
		names, err := b.store.getConsumerGroupNames(tx, "test")
		b.NoError(err)
		b.Equal([]string{"groupOne", "groupTwo"}, names)
		return nil
	})
	b.NoError(err)
}
func (b *BadgerTestSuite) TestLoadMeta() {
	err := b.db.Update(func(tx *badger.Txn) error {

		if err := tx.Set(key("test#groups#consumerOne#partition_count"), key("2")); err != nil {
			return err
		}
		if err := tx.Set(key("test#groups#consumerTwo#partition_count"), key("3")); err != nil {
			return err
		}
		return b.store.initTopicMeta(context.Background(), tx, "test")
	})
	b.NoError(err)
	err = b.db.View(func(tx *badger.Txn) error {
		meta, err := b.store.loadMeta(context.Background(), tx, "test")
		b.NoError(err)
		ts := time.Unix(b.ts.UTC().Unix(), 0)

		b.protoEqual(&kayakv1.TopicMetadata{
			Name:        "test",
			RecordCount: 0,
			CreatedAt:   timestamppb.New(ts),
			Archived:    false,
			GroupMetadata: map[string]*kayakv1.GroupPartitions{
				"consumerOne": &kayakv1.GroupPartitions{
					Name: "consumerOne",
				},
				"consumerTwo": &kayakv1.GroupPartitions{
					Name: "consumerTwo",
				},
			},
		}, meta)
		return nil
	})
	b.NoError(err)
}

func (s *BadgerTestSuite) protoEqual(expected, actual proto.Message) {
	s.Empty(cmp.Diff(expected, actual, protocmp.Transform()))
}
func (b *BadgerTestSuite) TestGetConsumerPartitions() {
	err := b.db.Update(func(tx *badger.Txn) error {
		if err := tx.Set(key("test#groups#groupOne#consumers#one"), key("1")); err != nil {
			return err
		}
		return tx.Set(key("test#groups#groupOne#consumers#two"), key("0"))
	})
	b.NoError(err)
	positions, err := b.store.GetConsumerPartitions(context.Background(), "test", "groupOne")
	b.NoError(err)
	expected := []*kayakv1.ConsumerGroupPartition{
		{
			Position:        "",
			PartitionNumber: 1,
			ConsumerId:      "one",
		},
	}
	b.Equal(expected, positions)
}
func TestBadgerTestSuite(t *testing.T) {
	suite.Run(t, new(BadgerTestSuite))
}
