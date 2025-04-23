package kayak

import (
	"context"
	"testing"

	"connectrpc.com/connect"
	"github.com/oklog/ulid/v2"
	"github.com/shoenig/test/must"
	"github.com/stretchr/testify/suite"

	kayakv1 "github.com/binarymatt/kayak/gen/kayak/v1"
	"github.com/binarymatt/kayak/internal/store"
)

type TestServiceSuite struct {
	suite.Suite
	service   *service
	mockStore *store.MockStore
	id        ulid.ULID
}

func (ts *TestServiceSuite) SetupTest() {
	s := store.NewMockStore(ts.T())
	ts.mockStore = s
	ts.id = ulid.Make()
	ts.service = &service{
		store: s,
		idGenerator: func() ulid.ULID {
			return ts.id
		},
	}
}
func (ts *TestServiceSuite) TestPutRecords() {
	t := ts.T()
	ctx := context.Background()
	cases := []struct {
		name        string
		record      *kayakv1.Record
		expected    []*kayakv1.Record
		expectedErr error
	}{
		{
			name: "simple case",
			record: &kayakv1.Record{
				Payload: []byte("test"),
				Id:      []byte("test"),
			},
			expected: []*kayakv1.Record{
				{
					Id:         []byte("test"),
					InternalId: ts.id.Bytes(),
					Payload:    []byte("test"),
					Partition:  0,
				},
			},
		},
		{
			name: "existing id",
			record: &kayakv1.Record{
				Payload: []byte("test"),
			},
			expected: []*kayakv1.Record{
				{
					Id:         ts.id.Bytes(),
					InternalId: ts.id.Bytes(),
					Payload:    []byte("test"),
					Partition:  0,
				},
			},
		},
	}
	for _, tc := range cases {

		ts.mockStore.EXPECT().GetStream("test_stream").
			Return(&kayakv1.Stream{Name: "test_stream", Partitions: []*kayakv1.Partition{{}}}, nil).
			Once()
		ts.mockStore.EXPECT().PutRecords("test_stream", tc.expected).Return(nil).Once()

		req := connect.NewRequest(&kayakv1.PutRecordsRequest{
			StreamName: "test_stream",
			Records:    []*kayakv1.Record{tc.record},
		})
		_, err := ts.service.PutRecords(ctx, req)
		must.ErrorIs(t, err, tc.expectedErr)
	}
}
func (ts *TestServiceSuite) TestGetRecords() {}

func TestExampleTestSuite(t *testing.T) {
	suite.Run(t, new(TestServiceSuite))
}
