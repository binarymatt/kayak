package client

import (
	"context"
	"testing"
	"time"

	"connectrpc.com/connect"
	"github.com/coder/quartz"
	"github.com/shoenig/test/must"
	"github.com/stretchr/testify/mock"
	"google.golang.org/protobuf/testing/protocmp"

	kayakv1 "github.com/binarymatt/kayak/gen/kayak/v1"
	"github.com/binarymatt/kayak/gen/kayak/v1/kayakv1connect"
)

func protoEq[V any](t *testing.T, expected, actual V) {
	t.Helper()
	must.Eq(t, expected, actual, must.Cmp(protocmp.Transform()))
}

type testContainer struct {
	kc         *KayakClient
	mockClient *kayakv1connect.MockKayakServiceClient
	mockClock  *quartz.Mock
}

func setupTest(t *testing.T) *testContainer {
	t.Helper()
	mockServiceClient := kayakv1connect.NewMockKayakServiceClient(t)
	mockClock := quartz.NewMock(t)
	client := &KayakClient{
		clock:  mockClock,
		client: mockServiceClient,
		cfg: &config{
			ticker:     1 * time.Second,
			streamName: "stream",
			group:      "group",
			id:         "one",
		},
		worker: &kayakv1.Worker{
			StreamName: "test",
		},
	}
	return &testContainer{
		kc:         client,
		mockClient: mockServiceClient,
		mockClock:  mockClock,
	}
}
func TestInit(t *testing.T) {
	tc := setupTest(t)
	ctx, cancel := context.WithCancel(context.Background())
	worker := &kayakv1.Worker{
		StreamName:          "stream",
		GroupName:           "group",
		Id:                  "one",
		PartitionAssignment: 1,
	}
	// ctx, cancel := context.WithCancel(context.Background())
	registerRequest := &kayakv1.RegisterWorkerRequest{
		StreamName: "stream",
		Group:      "group",
		Id:         "one",
	}
	registerResponse := &kayakv1.RegisterWorkerResponse{
		Worker: worker,
	}
	renewRequest := &kayakv1.RenewRegistrationRequest{
		Worker: worker,
	}
	tc.mockClient.
		EXPECT().
		RegisterWorker(mock.AnythingOfType("*context.cancelCtx"), connect.NewRequest(registerRequest)).
		Return(connect.NewResponse(registerResponse), nil).
		Once()
	tc.mockClient.
		EXPECT().
		RenewRegistration(mock.AnythingOfType("*context.cancelCtx"), connect.NewRequest(renewRequest)).
		Return(nil, nil).Once()

	err := tc.kc.Init(ctx) //nolint:errcheck
	must.NoError(t, err)
	w := tc.mockClock.Advance(1 * time.Second)
	err = w.Wait(ctx)
	must.NoError(t, err)
	cancel()

}

func TestFetchRecord(t *testing.T) {
	tc := setupTest(t)
	ctx := context.Background()
	expectedRecord := &kayakv1.Record{
		StreamName: "test",
		Partition:  0,
		Id:         "id",
	}
	req := &kayakv1.FetchRecordsRequest{
		StreamName: "test",
		Worker:     tc.kc.worker,
		Limit:      1,
	}
	resp := &kayakv1.FetchRecordsResponse{
		Records: []*kayakv1.Record{expectedRecord},
	}

	tc.mockClient.EXPECT().
		FetchRecords(ctx, connect.NewRequest(req)).
		Return(connect.NewResponse(resp), nil).Once()

	record, err := tc.kc.FetchRecord(ctx)
	must.NoError(t, err)
	protoEq(t, expectedRecord, record)
}

func TestCommitRecord(t *testing.T) {
	tc := setupTest(t)
	ctx := context.Background()
	record := &kayakv1.Record{
		StreamName: "test",
		Partition:  0,
		Id:         "test",
	}
	req := &kayakv1.CommitRecordRequest{
		Record: record,
		Worker: tc.kc.worker,
	}
	tc.mockClient.EXPECT().CommitRecord(ctx, connect.NewRequest(req)).Return(nil, nil).Once()
	err := tc.kc.CommitRecord(ctx, record)
	must.NoError(t, err)
}
