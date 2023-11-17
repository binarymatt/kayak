package service

import (
	"context"

	"connectrpc.com/connect"

	kayakv1 "github.com/binarymatt/kayak/gen/kayak/v1"
	"github.com/binarymatt/kayak/internal/test"
)

func (s *ServiceTestSuite) TestFetchRecord_SinglePartition() {
	r := s.Require()
	request := connect.NewRequest(&kayakv1.FetchRecordRequest{
		Topic:         "testTopic",
		ConsumerGroup: "testGroup",
		ConsumerId:    "1",
	})
	records := []*kayakv1.Record{
		{
			Topic:     "testTopic",
			Partition: 0,
			Payload:   []byte("first"),
		},
		{
			Topic:     "testTopic",
			Partition: 0,
			Payload:   []byte("second"),
		},
	}
	ctx := context.Background()
	s.mockStore.EXPECT().
		LoadMeta(ctx, "testTopic").
		Return(&kayakv1.TopicMetadata{
			GroupMetadata: map[string]*kayakv1.GroupPartitions{
				"testGroup": &kayakv1.GroupPartitions{
					Name:       "testGroup",
					Partitions: 1,
					Consumers: []*kayakv1.TopicConsumer{
						{
							Id:    "1",
							Topic: "testTopic",
						},
					},
				},
			},
		}, nil).
		Once()
	s.mockStore.EXPECT().
		GetRecords(ctx, "testTopic", "", 10).
		Return(records, nil).
		Once()
	resp, err := s.service.FetchRecord(ctx, request)
	r.NoError(err)
	test.ProtoEqual(s.T(), records[0], resp.Msg.Record)
}

func (s *ServiceTestSuite) TestFetchRecord_MultiplePartitions() {
	r := s.Require()
	request := connect.NewRequest(&kayakv1.FetchRecordRequest{
		Topic:         "testTopic",
		ConsumerGroup: "testGroup",
		ConsumerId:    "1",
	})
	records := []*kayakv1.Record{
		{
			Id:      "2Y5Pw5vJP4ATrUdjNENthrlucjb",
			Topic:   "testTopic",
			Payload: []byte("first"),
		},
		{
			Id:      "2Y5PvAX8aY98M5u2oQMrobCahtL",
			Topic:   "testTopic",
			Payload: []byte("second"),
		},
	}
	ctx := context.Background()
	s.mockStore.EXPECT().
		LoadMeta(ctx, "testTopic").
		Return(&kayakv1.TopicMetadata{
			GroupMetadata: map[string]*kayakv1.GroupPartitions{
				"testGroup": &kayakv1.GroupPartitions{
					Name:       "testGroup",
					Partitions: 2,
					Consumers: []*kayakv1.TopicConsumer{
						{
							Id:        "1",
							Topic:     "testTopic",
							Partition: 1,
						},
					},
				},
			},
		}, nil).
		Once()
	s.mockStore.EXPECT().
		GetRecords(ctx, "testTopic", "", 10).
		Return(records, nil).
		Once()
	resp, err := s.service.FetchRecord(ctx, request)
	r.NoError(err)
	test.ProtoEqual(s.T(), records[0], resp.Msg.Record)
}
