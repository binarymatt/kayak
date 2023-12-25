package service

import (
	"context"

	"connectrpc.com/connect"

	kayakv1 "github.com/binarymatt/kayak/gen/kayak/v1"
	"github.com/binarymatt/kayak/internal/store/models"
	"github.com/binarymatt/kayak/internal/test"
)

func (s *ServiceTestSuite) TestFetchRecord() {
	r := s.Require()
	records := createRecords()
	s.mockStore.EXPECT().FetchRecord(context.Background(), &models.Consumer{TopicID: "test", ID: "testConsumer"}).Return(records[0], nil)
	resp, err := s.service.FetchRecord(context.Background(), connect.NewRequest(&kayakv1.FetchRecordRequest{
		Consumer: &kayakv1.TopicConsumer{
			Id:    "testConsumer",
			Topic: "test",
		},
	}))
	r.NoError(err)
	test.ProtoEqual(s.T(), &kayakv1.Record{
		Id:        records[0].ID,
		Topic:     "test",
		Partition: records[0].Partition,
		Headers:   records[0].GetHeaders(),
		Payload:   records[0].Payload,
	}, resp.Msg.Record)
}
