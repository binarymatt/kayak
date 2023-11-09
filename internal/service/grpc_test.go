package service

import (
	"context"
	"testing"

	"connectrpc.com/connect"
	"github.com/stretchr/testify/require"

	kayakv1 "github.com/binarymatt/kayak/gen/kayak/v1"
)

func TestBalancer(t *testing.T) {
	cases := []struct {
		name       string
		key        string
		partitions int64
		expected   int64
	}{
		{
			name:       "single partition",
			key:        "test",
			partitions: 1,
			expected:   0,
		},
		{
			name:       "two partitions",
			key:        "test",
			partitions: 2,
			expected:   1,
		},
		{
			name:       "three partitions",
			key:        "test",
			partitions: 3,
			expected:   1,
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			r := require.New(t)
			res := balancer(tc.key, tc.partitions)
			r.Equal(tc.expected, res)
		})
	}
}
func (s *ServiceTestSuite) TestCreateConsumerGroup() {

	group := &kayakv1.ConsumerGroup{
		Name:           "group",
		Topic:          "topic",
		PartitionCount: 2,
		Hash:           kayakv1.Hash_HASH_MURMUR3,
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
