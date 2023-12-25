package service

import (
	"context"
	"testing"

	"connectrpc.com/connect"
	"github.com/stretchr/testify/require"
	"gorm.io/datatypes"

	kayakv1 "github.com/binarymatt/kayak/gen/kayak/v1"
	"github.com/binarymatt/kayak/internal/store/models"
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

func (s *ServiceTestSuite) TestRegisterConsumer() {
	s.mockStore.EXPECT().
		RegisterConsumer(context.TODO(), &models.Consumer{TopicID: "topic", Group: "group", ID: "id1", Partitions: datatypes.JSONSlice[int64]{1}}).
		Return(&models.Consumer{TopicID: "topic", Group: "group", ID: "id1", Position: "", Partitions: datatypes.JSONSlice[int64]{1}}, nil).Once()
	_, err := s.service.RegisterConsumer(context.Background(),
		connect.NewRequest(&kayakv1.RegisterConsumerRequest{
			Consumer: &kayakv1.TopicConsumer{
				Topic:      "topic",
				Group:      "group",
				Id:         "id1",
				Partitions: []int64{1},
			},
		}))
	s.NoError(err)
}
