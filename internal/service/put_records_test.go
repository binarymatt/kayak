package service

import (
	"context"
	"strconv"

	"connectrpc.com/connect"

	kayakv1 "github.com/binarymatt/kayak/gen/kayak/v1"
	"github.com/binarymatt/kayak/internal/store/models"
)

func (s *ServiceTestSuite) TestPutRecords() {
	records := []*kayakv1.Record{
		{
			Id:      s.testID.String(),
			Topic:   "test",
			Payload: []byte(" "),
		},
	}
	record := models.RecordFromProto(records[0])
	s.mockStore.EXPECT().AddRecords(context.TODO(), "test", record).Return(nil).Once()
	_, err := s.service.PutRecords(context.Background(), connect.NewRequest(&kayakv1.PutRecordsRequest{
		Topic:   "test",
		Records: records,
	}))
	s.NoError(err)
}

func (s *ServiceTestSuite) TestPutRecords_MissingTopic() {
	_, err := s.service.PutRecords(context.Background(), connect.NewRequest(&kayakv1.PutRecordsRequest{
		Records: []*kayakv1.Record{
			{Topic: "test", Payload: []byte(" ")},
		},
	}))
	s.Error(err)
	s.Equal("invalid_argument: validation error:\n - topic: value length must be at least 1 runes [string.min_len]", err.Error())
}

func (s *ServiceTestSuite) TestPutRecords_EmptyRecords() {
	_, err := s.service.PutRecords(context.Background(), connect.NewRequest(&kayakv1.PutRecordsRequest{
		Topic: "test",
	}))
	s.Error(err)
	s.Equal("invalid_argument: validation error:\n - records: value must contain at least 1 item(s) [repeated.min_items]", err.Error())
}

func (s *ServiceTestSuite) TestPutRecords_TooManyRecords() {
	records := []*kayakv1.Record{}
	for i := 0; i < 101; i++ {
		records = append(records, &kayakv1.Record{
			Id:      strconv.Itoa(i),
			Payload: []byte(" "),
		})
	}
	_, err := s.service.PutRecords(context.Background(), connect.NewRequest(&kayakv1.PutRecordsRequest{
		Topic:   "test",
		Records: records,
	}))
	s.Error(err)
	s.Equal("invalid_argument: validation error:\n - records: value must contain no more than 100 item(s) [repeated.max_items]", err.Error())
}
