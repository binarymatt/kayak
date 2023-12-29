package store

import (
	"context"
	"fmt"
	"io"
	"testing"
	"time"

	"github.com/oklog/ulid/v2"
	"github.com/stretchr/testify/suite"
	"gorm.io/datatypes"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"log/slog"

	kayakv1 "github.com/binarymatt/kayak/gen/kayak/v1"
	"github.com/binarymatt/kayak/internal/store/models"
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
	// s.db.Logger = logger.Default.LogMode(logger.Info)
}

func (s *SqlTestSuite) TestCreateTopic_Simple() {
	r := s.Require()
	t := &models.Topic{
		ID:             "test",
		PartitionCount: 1,
	}

	err := s.store.CreateTopic(s.ctx, t)
	r.NoError(err)

	var topic models.Topic
	result := s.db.First(&topic)
	r.Equal(int64(1), result.RowsAffected)
	r.NoError(result.Error)
	s.Equal(models.Topic{
		ID:             "test",
		PartitionCount: 1,
		Archived:       false,
		CreatedAt:      s.ts.Unix(),
		UpdatedAt:      s.ts.Unix(),
		DefaultHash:    kayakv1.Hash_HASH_MURMUR3.String(),
	}, topic)
}
func (s *SqlTestSuite) TestCreateTopic_Existing() {
	r := s.Require()
	t := &models.Topic{
		ID:             "test",
		PartitionCount: 1,
	}

	err := s.store.CreateTopic(s.ctx, t)
	r.NoError(err)

	err = s.store.CreateTopic(s.ctx, t)
	r.Error(err)

}

func (s *SqlTestSuite) TestDeleteTopic_Force() {
	r := s.Require()
	s.createTopic("test", 1)

	err := s.store.DeleteTopic(s.ctx, &models.Topic{ID: "test", Archived: false})
	r.NoError(err)

	var topic models.Topic
	result := s.db.First(&topic)
	r.Error(result.Error)
}

func (s *SqlTestSuite) TestDeleteTopic_Archive() {
	r := s.Require()

	s.createTopic("test", 1)

	err := s.store.DeleteTopic(s.ctx, &models.Topic{ID: "test", Archived: true})
	r.NoError(err)

	var topic models.Topic
	result := s.db.First(&topic)
	r.NoError(result.Error)
	r.Equal(int64(1), result.RowsAffected)
	r.Equal(models.Topic{
		ID:             "test",
		Archived:       true,
		PartitionCount: 1,
		CreatedAt:      s.ts.Unix(),
		UpdatedAt:      s.ts.Unix(),
		DefaultHash:    kayakv1.Hash_HASH_MURMUR3.String(),
		TTL:            30,
	}, topic)
}

func (s *SqlTestSuite) createTopic(name string, count int) {

	t := &models.Topic{
		ID:             name,
		PartitionCount: count,
		TTL:            30,
	}

	err := s.store.CreateTopic(s.ctx, t)
	s.Require().NoError(err)
}

func (s *SqlTestSuite) TestAddRecords() {
	r := s.Require()
	s.createTopic("test", 1)

	var records []models.Record
	err := s.store.AddRecords(s.ctx, "test", &models.Record{
		TopicID: "test",
		ID:      "testID",
		Headers: datatypes.JSONMap(map[string]interface{}{"header": "test"}),
		Payload: []byte("test record"),
	})
	r.NoError(err)

	result := s.db.Find(&records)
	r.NoError(result.Error)
	r.Equal(int64(1), result.RowsAffected)
	r.Equal("test", records[0].TopicID)
	r.Equal(datatypes.JSONMap(map[string]interface{}{"header": "test"}), records[0].Headers)
	r.Equal("testID", records[0].ID)
}

func (s *SqlTestSuite) TestAddRecords_Balance() {
	r := s.Require()
	t := &models.Topic{
		ID:             "test",
		PartitionCount: 2,
	}

	err := s.store.CreateTopic(s.ctx, t)
	r.NoError(err)

	_, err = s.store.RegisterConsumer(s.ctx, &models.Consumer{
		ID:         "testConsumer",
		TopicID:    "test",
		Partitions: datatypes.NewJSONSlice([]int64{1}),
	})
	r.NoError(err)
	records := generateRecords("test", 10)
	err = s.store.AddRecords(s.ctx, "test", records...)
	r.NoError(err)

	var items []models.Record
	s.db.Find(&items)
	r.Len(items, 10)
	partitions := map[string]int64{}
	for i, item := range items {
		partitions[fmt.Sprintf("%d", item.Partition)] = item.Partition
		r.Equal(models.Record{
			ID:        records[i].ID,
			TopicID:   "test",
			Partition: item.Partition,
			Headers:   records[i].Headers,
			Payload:   records[i].Payload,
			CreatedAt: s.ts.Unix(),
			UpdatedAt: s.ts.Unix(),
		}, item)
	}
	r.NotEmpty(partitions)
	r.Len(partitions, 2)

}

func (s *SqlTestSuite) TestAddRecords_InvalidTopic() {
	r := s.Require()
	err := s.store.AddRecords(s.ctx, "test", &models.Record{})
	r.ErrorIs(err, ErrInvalidTopic)
}

func (s *SqlTestSuite) TestAddRecords_ArchivedTopic() {
	r := s.Require()
	s.createTopic("test", 1)

	result := s.db.Model(&models.Topic{}).Where("id = ?", "test").Update("archived", true)
	r.NoError(result.Error)
	err := s.store.AddRecords(s.ctx, "test", &models.Record{})
	r.ErrorIs(err, ErrArchivedTopic)
}

func (s *SqlTestSuite) TestGetRecords_Simple() {
	r := s.Require()
	s.createTopic("testget", 1)

	records := generateRecords("testget", 3)

	err := s.store.AddRecords(s.ctx, "testget", records...)
	r.NoError(err)

	items, err := s.store.GetRecords(s.ctx, "testget", "", 10)
	r.NoError(err)
	r.Len(items, 3)
}
func (s *SqlTestSuite) TestGetRecords_Empty() {
	r := s.Require()

	err := s.store.CreateTopic(s.ctx, &models.Topic{ID: "testget"})
	r.NoError(err)

	items, err := s.store.GetRecords(s.ctx, "testget", "", 10)
	r.NoError(err)
	r.Len(items, 0)
}
func (s *SqlTestSuite) TestGetRecords_StartKey() {
	r := s.Require()

	err := s.store.CreateTopic(s.ctx, &models.Topic{
		ID:             "testget",
		PartitionCount: 1,
	})
	r.NoError(err)

	records := generateRecords("testget", 3)

	err = s.store.AddRecords(s.ctx, "testget", records...)
	r.NoError(err)

	items, err := s.store.GetRecords(s.ctx, "testget", records[0].ID, 10)
	r.NoError(err)
	r.Len(items, 2)
	r.Equal(records[1].ID, items[0].ID)
	r.Equal(records[2].ID, items[1].ID)
}
func (s *SqlTestSuite) TestGetRecords_Limit() {

	r := s.Require()
	err := s.store.CreateTopic(s.ctx, &models.Topic{ID: "testget", PartitionCount: 1})
	r.NoError(err)

	records := generateRecords("testget", 3)

	err = s.store.AddRecords(s.ctx, "testget", records...)
	r.NoError(err)

	items, err := s.store.GetRecords(s.ctx, "testget", "", 2)
	r.NoError(err)
	r.Len(items, 2)
	r.Equal(records[0].ID, items[0].ID)
	r.Equal(records[1].ID, items[1].ID)
}

func (s *SqlTestSuite) TestGetRecords_InvalidTopic() {
	r := s.Require()
	items, err := s.store.GetRecords(s.ctx, "testget", "", 2)
	r.ErrorIs(err, ErrInvalidTopic)
	r.Len(items, 0)
}

func (s *SqlTestSuite) TestGetRecords_ArchivedTopic() {
	r := s.Require()

	err := s.store.CreateTopic(s.ctx, &models.Topic{ID: "testget"})
	r.NoError(err)

	err = s.store.DeleteTopic(s.ctx, &models.Topic{ID: "testget", Archived: true})
	r.NoError(err)

	items, err := s.store.GetRecords(s.ctx, "testget", "", 2)
	r.ErrorIs(err, ErrArchivedTopic)
	r.Len(items, 0)
}

func (s *SqlTestSuite) TestListTopics() {
	r := s.Require()

	topics, err := s.store.ListTopics(s.ctx)
	r.NoError(err)
	r.Empty(topics)

	s.createTopic("test", 1)
	topics, err = s.store.ListTopics(s.ctx)
	r.NoError(err)
	r.ElementsMatch([]string{"test"}, topics)
}

func (s *SqlTestSuite) TestRegisterConsumer_InvalidTopic() {
	r := s.Require()

	_, err := s.store.RegisterConsumer(s.ctx, &models.Consumer{})
	r.ErrorIs(err, ErrInvalidTopic)
}

func (s *SqlTestSuite) TestRegisterConsumer_HappyPath() {
	r := s.Require()
	s.createTopic("test", 1)

	topicConsumer, err := s.store.RegisterConsumer(s.ctx, &models.Consumer{
		TopicID:    "test",
		Group:      "testGroup",
		ID:         "testID",
		Partitions: datatypes.NewJSONSlice([]int64{1, 2, 3}),
	})
	r.NoError(err)
	r.Equal(&models.Consumer{
		ID:         "testID",
		TopicID:    "test",
		Group:      "testGroup",
		Partitions: datatypes.JSONSlice[int64]{1, 2, 3},
		CreatedAt:  s.ts.Unix(),
		UpdatedAt:  s.ts.Unix(),
	}, topicConsumer)

	var consumer models.Consumer
	result := s.db.First(&consumer, "id = ?", "testID")
	r.NoError(result.Error)
	r.Equal(models.Consumer{
		ID:         "testID",
		TopicID:    "test",
		Group:      "testGroup",
		CreatedAt:  s.ts.Unix(),
		UpdatedAt:  s.ts.Unix(),
		Partitions: datatypes.JSONSlice[int64]{1, 2, 3},
	}, consumer)
}

func (s *SqlTestSuite) TestGetConsumer_Unknown() {
	r := s.Require()

	consumer, err := s.store.getConsumer("1234")
	r.ErrorIs(err, gorm.ErrRecordNotFound)
	r.Nil(consumer)
}
func (s *SqlTestSuite) TestGetConsumer() {
	r := s.Require()
	res := s.db.Create(&models.Consumer{
		ID:       "123",
		Position: "",
	})
	r.NoError(res.Error)

	consumer, err := s.store.getConsumer("123")
	r.NoError(err)
	r.Equal(&models.Consumer{
		ID:        "123",
		Position:  "",
		CreatedAt: s.ts.Unix(),
		UpdatedAt: s.ts.Unix(),
	}, consumer)

}

func (s *SqlTestSuite) TestGetConsumerPosition() {

	r := s.Require()
	res := s.db.Create(&models.Consumer{
		ID:       "123",
		Position: "start",
	})
	r.NoError(res.Error)

	position, err := s.store.GetConsumerPosition(s.ctx, &models.Consumer{ID: "123"})
	r.NoError(err)
	r.Equal("start", position)
}

func (s *SqlTestSuite) TestCommitConsumerPosition_Happy() {
	r := s.Require()

	res := s.db.Create(&models.Consumer{
		ID:       "123",
		Position: "",
	})
	r.NoError(res.Error)

	err := s.store.CommitConsumerPosition(s.ctx, &models.Consumer{
		ID:       "123",
		Position: "item1",
	})
	r.NoError(err)
}

func (s *SqlTestSuite) TestCommitConsumerPosition_NonExistent() {
	r := s.Require()
	err := s.store.CommitConsumerPosition(s.ctx, &models.Consumer{
		ID:       "1233",
		Position: "item1",
	})
	r.ErrorIs(err, ErrInvalidConsumer)
}

func (s *SqlTestSuite) TestFetchRecord_SinglePartition() {
	r := s.Require()
	s.createTopic("test", 1)

	_, err := s.store.RegisterConsumer(s.ctx, &models.Consumer{
		ID:         "testConsumer",
		TopicID:    "test",
		Partitions: datatypes.NewJSONSlice([]int64{0}),
	})
	r.NoError(err)
	records := generateRecords("test", 3)
	err = s.store.AddRecords(s.ctx, "test", records...)
	r.NoError(err)

	record, err := s.store.FetchRecord(s.ctx, &models.Consumer{TopicID: "test", ID: "testConsumer"})
	r.NoError(err)
	s.Equal(records[0].ID, record.ID)

}
func (s *SqlTestSuite) TestFetchRecord_AllPartitions() {
	s.createTopic("test", 2)

	_, err := s.store.RegisterConsumer(s.ctx, &models.Consumer{
		ID:         "testConsumer",
		TopicID:    "test",
		Partitions: datatypes.NewJSONSlice([]int64{0, 1}),
	})
	s.NoError(err)
	records := generateRecords("test", 3)
	err = s.store.AddRecords(s.ctx, "test", records...)
	s.NoError(err)

	record, err := s.store.FetchRecord(s.ctx, &models.Consumer{TopicID: "test", ID: "testConsumer"})
	s.NoError(err)
	s.Equal(records[0].ID, record.ID)
}
func (s *SqlTestSuite) TestFetchRecord_TopicArchived() {

	s.createTopic("test", 1)

	_, err := s.store.RegisterConsumer(s.ctx, &models.Consumer{
		ID:         "testConsumer",
		TopicID:    "test",
		Partitions: datatypes.NewJSONSlice([]int64{0, 1}),
	})
	s.NoError(err)
	records := generateRecords("test", 3)
	err = s.store.AddRecords(s.ctx, "test", records...)
	s.NoError(err)
	s.NoError(s.store.DeleteTopic(s.ctx, &models.Topic{ID: "test", Archived: true}))
	record, err := s.store.FetchRecord(s.ctx, &models.Consumer{TopicID: "test", ID: "testConsumer"})
	s.ErrorIs(err, ErrArchivedTopic)
	s.Nil(record)
}
func (s *SqlTestSuite) TestFetchRecord_UnregisteredConsumer() {
	s.createTopic("test", 1)
	record, err := s.store.FetchRecord(s.ctx, &models.Consumer{TopicID: "test", ID: "testConsumer"})
	s.ErrorIs(err, ErrInvalidConsumer)
	s.Nil(record)
}

func (s *SqlTestSuite) TestPruneOldRecords() {
	s.createTopic("test", 1)
	records := generateRecords("test", 5)
	old := time.Now().Add(-(45 * time.Second)).Unix()
	records[0].CreatedAt = old
	err := s.store.AddRecords(s.ctx, "test", records...)
	s.NoError(err)

	err = s.store.PruneOldRecords(s.ctx)
	s.NoError(err)

	var r []models.Record
	err = s.db.Find(&r).Error
	s.NoError(err)

	s.Len(r, 4)
}

func generateRecords(topic string, total int) []*models.Record {
	records := []*models.Record{}
	for i := 0; i < total; i++ {
		records = append(records, &models.Record{
			TopicID: topic,
			ID:      ulid.Make().String(),
			Headers: datatypes.JSONMap(map[string]interface{}{"header": fmt.Sprintf("%d", i)}),
			Payload: []byte(fmt.Sprintf("record %d", i)),
		})
	}
	return records
}
