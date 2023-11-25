package store

import (
	"context"
	"io"
	"testing"
	"time"

	"github.com/oklog/ulid/v2"
	"github.com/stretchr/testify/suite"
	"golang.org/x/exp/slog"
	"gorm.io/datatypes"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	kayakv1 "github.com/binarymatt/kayak/gen/kayak/v1"
	"github.com/binarymatt/kayak/internal/store/models"
	"github.com/binarymatt/kayak/internal/test"
)

type SqlTestSuite struct {
	suite.Suite
	store *sqlStore
	db    *gorm.DB
	ctx   context.Context
	ts    time.Time
	id    ulid.ULID
}

func TestSqlTestSuite(t *testing.T) {
	suite.Run(t, new(SqlTestSuite))
}
func (b *SqlTestSuite) SetupSuite() {
	l := slog.New(slog.NewTextHandler(io.Discard, nil))
	slog.SetDefault(l)
}
func (s *SqlTestSuite) SetupTest() {
	s.ts = time.Now().UTC()
	db := sqlite.Open("file::memory:")
	//db := sqlite.Open("temp.db")
	s.db, _ = gorm.Open(db, &gorm.Config{
		NowFunc: func() time.Time {
			return s.ts
		},
		Logger: logger.Discard,
	})
	s.store = NewSqlStore(s.db)
	s.id = ulid.Make()
	s.store.idFunc = func() ulid.ULID {
		return s.id
	}
	err := s.store.RunMigrations()
	s.Require().NoError(err)
}
func (s *SqlTestSuite) TestCreateTopic_Simple() {
	r := s.Require()

	err := s.store.CreateTopic(s.ctx, "test")
	r.NoError(err)

	var topic models.Topic
	result := s.db.First(&topic)
	r.Equal(int64(1), result.RowsAffected)
	r.NoError(result.Error)
	s.Equal(models.Topic{
		ID:        s.id.String(),
		Name:      "test",
		Archived:  false,
		CreatedAt: s.ts.Unix(),
		UpdatedAt: s.ts.Unix(),
	}, topic)
}
func (s *SqlTestSuite) TestCreateTopic_Existing() {
	r := s.Require()

	err := s.store.CreateTopic(s.ctx, "test")
	r.NoError(err)

	err = s.store.CreateTopic(s.ctx, "test")
	r.Error(err)

}

func (s *SqlTestSuite) TestDeleteTopic_Force() {
	r := s.Require()

	err := s.store.CreateTopic(s.ctx, "test")
	r.NoError(err)

	err = s.store.DeleteTopic(s.ctx, "test", true)
	r.NoError(err)
	var topic models.Topic
	result := s.db.First(&topic)
	r.Error(result.Error)
}

func (s *SqlTestSuite) TestDeleteTopic_Archive() {
	r := s.Require()

	err := s.store.CreateTopic(s.ctx, "test")
	r.NoError(err)

	err = s.store.DeleteTopic(s.ctx, "test", false)
	r.NoError(err)

	var topic models.Topic
	result := s.db.First(&topic)
	r.NoError(result.Error)
	r.Equal(int64(1), result.RowsAffected)
	r.Equal(models.Topic{
		ID:        s.id.String(),
		Name:      "test",
		Archived:  true,
		CreatedAt: s.ts.Unix(),
		UpdatedAt: s.ts.Unix(),
	}, topic)
}

func (s *SqlTestSuite) TestAddRecords() {
	r := s.Require()

	err := s.store.CreateTopic(s.ctx, "test")
	r.NoError(err)
	var records []models.Record
	err = s.store.AddRecords(s.ctx, "test", &kayakv1.Record{
		Topic:   "test",
		Id:      "testID",
		Headers: map[string]string{"header": "test"},
		Payload: []byte("test record"),
	})
	r.NoError(err)

	result := s.db.Find(&records)
	r.NoError(result.Error)
	r.Equal(int64(1), result.RowsAffected)
	r.Equal(s.id.String(), records[0].TopicID)
	r.Equal(datatypes.JSONMap(map[string]interface{}{"header": "test"}), records[0].Headers)
	r.Equal("testID", records[0].ID)
}

func (s *SqlTestSuite) TestAddRecords_InvalidTopic() {
	r := s.Require()
	err := s.store.AddRecords(s.ctx, "test", &kayakv1.Record{})
	r.ErrorIs(err, ErrInvalidTopic)
}

func (s *SqlTestSuite) TestAddRecords_ArchivedTopic() {
	r := s.Require()
	err := s.store.CreateTopic(s.ctx, "test")
	r.NoError(err)
	result := s.db.Model(&models.Topic{}).Where("name = ?", "test").Update("archived", true)
	r.NoError(result.Error)
	err = s.store.AddRecords(s.ctx, "test", &kayakv1.Record{})
	r.ErrorIs(err, ErrTopicArchived)
}

func (s *SqlTestSuite) TestGetRecords_Simple() {
	r := s.Require()
	err := s.store.CreateTopic(s.ctx, "testget")
	r.NoError(err)

	records := []*kayakv1.Record{
		{
			Topic:   "testget",
			Id:      ulid.Make().String(),
			Headers: map[string]string{"header": "first"},
			Payload: []byte("first record"),
		},
		{
			Topic:   "testget",
			Id:      ulid.Make().String(),
			Headers: map[string]string{"header": "second"},
			Payload: []byte("second record"),
		},
		{
			Topic:   "testget",
			Id:      ulid.Make().String(),
			Headers: map[string]string{"header": "third"},
			Payload: []byte("third record"),
		},
	}

	err = s.store.AddRecords(s.ctx, "testget", records...)
	r.NoError(err)

	items, err := s.store.GetRecords(s.ctx, "testget", "", 10)
	r.NoError(err)
	r.Len(items, 3)
}
func (s *SqlTestSuite) TestGetRecords_Empty() {
	r := s.Require()

	err := s.store.CreateTopic(s.ctx, "testget")
	r.NoError(err)

	items, err := s.store.GetRecords(s.ctx, "testget", "", 10)
	r.NoError(err)
	r.Len(items, 0)
}
func (s *SqlTestSuite) TestGetRecords_StartKey() {
	r := s.Require()

	err := s.store.CreateTopic(s.ctx, "testget")
	r.NoError(err)

	records := []*kayakv1.Record{
		{
			Id:      ulid.Make().String(),
			Topic:   "testget",
			Headers: map[string]string{"header": "first"},
			Payload: []byte("first record"),
		},
		{
			Id:      ulid.Make().String(),
			Topic:   "testget",
			Headers: map[string]string{"header": "second"},
			Payload: []byte("second record"),
		},
		{
			Id:      ulid.Make().String(),
			Topic:   "testget",
			Headers: map[string]string{"header": "third"},
			Payload: []byte("third record"),
		},
	}

	err = s.store.AddRecords(s.ctx, "testget", records...)
	r.NoError(err)

	items, err := s.store.GetRecords(s.ctx, "testget", records[0].Id, 10)
	r.NoError(err)
	r.Len(items, 2)
	r.Equal(records[1].Id, items[0].Id)
	r.Equal(records[2].Id, items[1].Id)
}
func (s *SqlTestSuite) TestGetRecords_Limit() {

	r := s.Require()
	err := s.store.CreateTopic(s.ctx, "testget")
	r.NoError(err)

	records := []*kayakv1.Record{
		{
			Id:      ulid.Make().String(),
			Topic:   "testget",
			Headers: map[string]string{"header": "first"},
			Payload: []byte("first record"),
		},
		{
			Id:      ulid.Make().String(),
			Topic:   "testget",
			Headers: map[string]string{"header": "second"},
			Payload: []byte("second record"),
		},
		{
			Id:      ulid.Make().String(),
			Topic:   "testget",
			Headers: map[string]string{"header": "third"},
			Payload: []byte("third record"),
		},
	}

	err = s.store.AddRecords(s.ctx, "testget", records...)
	r.NoError(err)

	items, err := s.store.GetRecords(s.ctx, "testget", "", 2)
	r.NoError(err)
	r.Len(items, 2)
	r.Equal(records[0].Id, items[0].Id)
	r.Equal(records[1].Id, items[1].Id)
}

func (s *SqlTestSuite) TestGetRecords_InvalidTopic() {
	r := s.Require()
	items, err := s.store.GetRecords(s.ctx, "testget", "", 2)
	r.ErrorIs(err, ErrInvalidTopic)
	r.Len(items, 0)
}

func (s *SqlTestSuite) TestGetRecords_ArchivedTopic() {
	r := s.Require()

	err := s.store.CreateTopic(s.ctx, "testget")
	r.NoError(err)

	err = s.store.DeleteTopic(s.ctx, "testget", false)
	r.NoError(err)

	items, err := s.store.GetRecords(s.ctx, "testget", "", 2)
	r.ErrorIs(err, ErrTopicArchived)
	r.Len(items, 0)
}

func (s *SqlTestSuite) TestListTopics() {
	r := s.Require()

	topics, err := s.store.ListTopics(s.ctx)
	r.NoError(err)
	r.Empty(topics)

	err = s.store.CreateTopic(s.ctx, "test")
	r.NoError(err)
	topics, err = s.store.ListTopics(s.ctx)
	r.NoError(err)
	r.ElementsMatch([]string{"test"}, topics)
}

func (s *SqlTestSuite) TestRegisterConsumerGroup() {
	r := s.Require()

	err := s.store.CreateTopic(s.ctx, "test")
	r.NoError(err)

	err = s.store.RegisterConsumerGroup(s.ctx, &kayakv1.ConsumerGroup{
		Name:           "testGroup",
		Topic:          "test",
		PartitionCount: 1,
		Hash:           kayakv1.Hash_HASH_MURMUR3,
	})
	r.NoError(err)

	var cg models.ConsumerGroup
	result := s.db.First(&cg)
	r.NoError(result.Error)
	r.Equal(int64(1), result.RowsAffected)
	r.Equal(models.ConsumerGroup{
		ID:             1,
		Name:           "testGroup",
		TopicID:        s.id.String(),
		PartitionCount: 1,
		Hash:           "HASH_MURMUR3",
		CreatedAt:      s.ts.Unix(),
		UpdatedAt:      s.ts.Unix(),
	}, cg)
}
func (s *SqlTestSuite) TestRegisterConsumerGroup_InvalidTopic() {
	err := s.store.RegisterConsumerGroup(s.ctx, &kayakv1.ConsumerGroup{
		Name:           "testGroup",
		Topic:          "test",
		PartitionCount: 1,
		Hash:           kayakv1.Hash_HASH_MURMUR3,
	})
	s.Require().ErrorIs(err, ErrInvalidTopic)
}

func (s *SqlTestSuite) TestRegisterConsumerGroup_ArchivedTopic() {
	r := s.Require()
	result := s.db.Create(&models.Topic{
		ID:       s.id.String(),
		Name:     "test",
		Archived: true,
	})
	r.NoError(result.Error)
	err := s.store.RegisterConsumerGroup(s.ctx, &kayakv1.ConsumerGroup{
		Name:           "testGroup",
		Topic:          "test",
		PartitionCount: 1,
		Hash:           kayakv1.Hash_HASH_MURMUR3,
	})
	r.ErrorIs(err, ErrTopicArchived)
}
func (s *SqlTestSuite) TestRegisterConsumerGroup_ExistingGroup() {
	r := s.Require()
	result := s.db.Create(&models.Topic{
		ID:       s.id.String(),
		Name:     "test",
		Archived: false,
	})
	r.NoError(result.Error)
	result = s.db.Create(&models.ConsumerGroup{
		Name:           "testGroup",
		TopicID:        s.id.String(),
		PartitionCount: 1,
		Hash:           kayakv1.Hash_HASH_MURMUR3.String(),
	})
	r.NoError(result.Error)

	err := s.store.RegisterConsumerGroup(s.ctx, &kayakv1.ConsumerGroup{
		Name:           "testGroup",
		Topic:          "test",
		PartitionCount: 2,
	})
	r.ErrorIs(err, ErrConsumerGroupExists)
}

func (s *SqlTestSuite) TestGetConsumerGroup() {
	r := s.Require()

	err := s.store.CreateTopic(s.ctx, "test")
	r.NoError(err)

	group, err := s.store.getConsumerGroup("test", "group")
	r.ErrorIs(err, ErrConsumerGroupInvalid)
	r.Nil(group)

	err = s.store.RegisterConsumerGroup(s.ctx, &kayakv1.ConsumerGroup{
		Name:  "group",
		Topic: "test",
		Hash:  kayakv1.Hash_HASH_MURMUR3,
	})
	r.NoError(err)

	group, err = s.store.getConsumerGroup("test", "group")
	r.NoError(err)
	r.Equal(&models.ConsumerGroup{
		ID:        1,
		Name:      "group",
		TopicID:   s.id.String(),
		Hash:      "HASH_MURMUR3",
		CreatedAt: s.ts.Unix(),
		UpdatedAt: s.ts.Unix(),
	}, group)
}

func (s *SqlTestSuite) TestRegisterConsumer_InvalidTopic() {
	r := s.Require()

	_, err := s.store.RegisterConsumer(s.ctx, &kayakv1.TopicConsumer{})
	r.ErrorIs(err, ErrInvalidTopic)
}

func (s *SqlTestSuite) TestRegisterConsumer_InvalidGroup() {

	r := s.Require()
	err := s.store.CreateTopic(s.ctx, "test")
	r.NoError(err)

	consumer, err := s.store.RegisterConsumer(s.ctx, &kayakv1.TopicConsumer{
		Topic: "test",
		Group: "group",
	})
	r.ErrorIs(err, ErrConsumerGroupInvalid)
	r.Nil(consumer)
}

func (s *SqlTestSuite) TestRegisterConsumer_HappyPath() {
	r := s.Require()
	err := s.store.CreateTopic(s.ctx, "test")
	r.NoError(err)
	err = s.store.RegisterConsumerGroup(s.ctx, &kayakv1.ConsumerGroup{
		Topic:          "test",
		Name:           "testGroup",
		PartitionCount: 1,
		Hash:           kayakv1.Hash_HASH_MURMUR3,
	})
	r.NoError(err)

	topicConsumer, err := s.store.RegisterConsumer(s.ctx, &kayakv1.TopicConsumer{
		Topic: "test",
		Group: "testGroup",
		Id:    "testID",
	})
	r.NoError(err)
	test.ProtoEqual(s.T(), &kayakv1.TopicConsumer{
		Id:        "testID",
		Topic:     "test",
		Group:     "testGroup",
		Partition: 0,
	}, topicConsumer)

	var consumer models.Consumer
	result := s.db.First(&consumer, "id = ?", "testID")
	r.NoError(result.Error)
	r.Equal(models.Consumer{
		ID:        "testID",
		GroupID:   1,
		CreatedAt: s.ts.Unix(),
		UpdatedAt: s.ts.Unix(),
	}, consumer)
}
func (s *SqlTestSuite) TestRegisterConsumer_GroupFull() {}
