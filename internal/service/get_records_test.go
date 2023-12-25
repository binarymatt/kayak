package service

import (
	"context"

	"connectrpc.com/connect"
	"github.com/oklog/ulid/v2"

	kayakv1 "github.com/binarymatt/kayak/gen/kayak/v1"
	"github.com/binarymatt/kayak/internal/store/models"
)

func createProtoRecords() []*kayakv1.Record {
	return []*kayakv1.Record{
		{
			Topic: "test",
			Id:    ulid.Make().String(),
			Headers: map[string]string{
				"name": "one",
			},
			Payload: []byte("first record"),
		},
		{
			Topic: "test",
			Id:    ulid.Make().String(),
			Headers: map[string]string{
				"name": "two",
			},
			Payload: []byte("second record"),
		},
		{
			Topic: "test",
			Id:    ulid.Make().String(),
			Headers: map[string]string{
				"name": "three",
			},
			Payload: []byte("third record"),
		},
	}
}
func createRecords() []*models.Record {
	items := createProtoRecords()
	records := make([]*models.Record, len(items))
	for i, item := range items {
		records[i] = models.RecordFromProto(item)
	}
	return records
}
func (s *ServiceTestSuite) TestGetRecords_SimplePath() {
	records := createRecords()
	request := &kayakv1.GetRecordsRequest{
		Topic: "test",
	}
	s.mockStore.EXPECT().GetRecords(context.Background(), "test", "", 99).Return(records, nil).Once()
	resp, err := s.service.GetRecords(context.Background(), connect.NewRequest(request))
	s.NoError(err)
	s.Len(resp.Msg.Records, 3)
}
