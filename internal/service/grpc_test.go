package service

import (
	"context"

	"connectrpc.com/connect"

	kayakv1 "github.com/binarymatt/kayak/gen/kayak/v1"
)

func (s *ServiceTestSuite) TestCreateConsumerGroup() {

	group := &kayakv1.ConsumerGroup{
		Name:           "group",
		Topic:          "topic",
		PartitionCount: 2,
	}
	s.mockStore.EXPECT().RegisterConsumerGroup(
		context.TODO(),
		group,
	).Return(nil).Once()
	_, err := s.service.CreateConsumerGroup(
		context.Background(),
		connect.NewRequest(&kayakv1.CreateConsumerGroupRequest{
			Group: group,
		}),
	)
	s.NoError(err)
}

func (s *ServiceTestSuite) TestRegisterConsumer() {
	s.mockStore.EXPECT().
		RegisterConsumer(context.TODO(), &kayakv1.TopicConsumer{Topic: "topic", Group: "group", Id: "id1"}).
		Return(&kayakv1.TopicConsumer{Topic: "topic", Group: "group", Id: "id1", Position: "", Partition: 1}, nil).Once()
	_, err := s.service.RegisterConsumer(context.Background(),
		connect.NewRequest(&kayakv1.RegisterConsumerRequest{
			Consumer: &kayakv1.TopicConsumer{
				Topic: "topic",
				Group: "group",
				Id:    "id1",
			},
		}))
	s.NoError(err)
}
