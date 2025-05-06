package store

import (
	"os"
	"testing"
	"time"

	"github.com/dgraph-io/badger/v4"
	"github.com/oklog/ulid/v2"
	"github.com/shoenig/test/must"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/testing/protocmp"

	kayakv1 "github.com/binarymatt/kayak/gen/kayak/v1"
)

type testContainer struct {
	db    *badger.DB
	store *store
}

func setupTest(t *testing.T) *testContainer {
	t.Helper()
	opts := badger.DefaultOptions("").WithInMemory(true).WithLogger(nil)
	db, err := badger.Open(opts)
	must.NoError(t, err)
	t.Cleanup(func() {
		db.Close() //nolint: errcheck
	})
	store := New(db)
	return &testContainer{
		db:    db,
		store: store,
	}
}

func TestPutStream(t *testing.T) {
	ts := setupTest(t)
	expected := &kayakv1.Stream{
		Name:           "test",
		PartitionCount: 1,
	}
	err := ts.store.PutStream(expected)
	must.NoError(t, err)
	txn := ts.db.NewTransaction(false)
	defer txn.Discard()
	item, err := txn.Get([]byte("streams:test"))
	must.NoError(t, err)
	raw, err := item.ValueCopy(nil)
	must.NoError(t, err)
	var s kayakv1.Stream
	err = proto.Unmarshal(raw, &s)
	must.NoError(t, err)
	protoEq(t, expected, &s)
}

func TestPutStream_Overwrite(t *testing.T) {
	ts := setupTest(t)
	expected := &kayakv1.Stream{
		Name:           "test",
		PartitionCount: 2,
	}
	err := ts.store.PutStream(&kayakv1.Stream{
		Name:           "test",
		PartitionCount: 1,
	})
	must.NoError(t, err)

	ts.db.View(func(txn *badger.Txn) error { //nolint:errcheck
		item, err := txn.Get([]byte("streams:test"))
		must.NoError(t, err)
		raw, err := item.ValueCopy(nil)
		must.NoError(t, err)
		var s kayakv1.Stream
		err = proto.Unmarshal(raw, &s)
		must.NoError(t, err)
		must.Eq(t, 1, s.PartitionCount)
		return nil
	})

	err = ts.store.PutStream(expected)
	must.NoError(t, err)
	ts.db.View(func(txn *badger.Txn) error { //nolint:errcheck
		item, err := txn.Get([]byte("streams:test"))
		must.NoError(t, err)
		raw, err := item.ValueCopy(nil)
		must.NoError(t, err)
		var s kayakv1.Stream
		err = proto.Unmarshal(raw, &s)
		must.NoError(t, err)
		protoEq(t, expected, &s)
		return nil
	})

}

func TestGetStream(t *testing.T) {
	ts := setupTest(t)
	expected := &kayakv1.Stream{
		Name:           "test",
		Ttl:            1,
		PartitionCount: 1,
	}
	err := ts.store.PutStream(expected)
	must.NoError(t, err)

	stream, err := ts.store.GetStream("test")
	must.NoError(t, err)
	protoEq(t, expected, stream)
}

func TestGetStream_Missing(t *testing.T) {
	ts := setupTest(t)
	stream, err := ts.store.GetStream("test")
	must.ErrorIs(t, err, badger.ErrKeyNotFound)
	must.Nil(t, stream)
}

func TestPutRecords(t *testing.T) {
	ts := setupTest(t)
	stream := &kayakv1.Stream{
		Name:           "test",
		PartitionCount: 1,
		Ttl:            0,
	}
	record := &kayakv1.Record{
		StreamName: "test",
		InternalId: ulid.Make().String(),
		Id:         []byte("test"),
		Partition:  0,
	}
	err := ts.store.PutStream(stream)
	must.NoError(t, err)
	err = ts.store.PutRecords("test", record)
	must.NoError(t, err)
	ts.db.View(func(tx *badger.Txn) error { //nolint:errcheck
		item, err := tx.Get(recordKey("test", 0, record.InternalId))
		must.NoError(t, err)
		var actual kayakv1.Record
		raw, err := item.ValueCopy(nil)
		must.NoError(t, err)
		proto.Unmarshal(raw, &actual) //nolint:errcheck
		protoEq(t, record, &actual)
		return nil
	})
}

func TestPutRecords_NoStream(t *testing.T) {
	ts := setupTest(t)
	record := &kayakv1.Record{
		StreamName: "test",
		InternalId: ulid.Make().String(),
		Id:         []byte("test"),
		Partition:  0,
	}
	err := ts.store.PutRecords("test", record)
	must.ErrorIs(t, err, badger.ErrKeyNotFound)
	must.Eq(t, "missing stream test: Key not found", err.Error())
}

func TestPutRecords_TTL(t *testing.T) {
	ts := setupTest(t)
	stream := &kayakv1.Stream{
		Name:           "test",
		PartitionCount: 1,
		Ttl:            1,
	}
	record := &kayakv1.Record{
		StreamName: "test",
		InternalId: ulid.Make().String(),
		Id:         []byte("test"),
		Partition:  0,
	}
	now := time.Now()
	err := ts.store.PutStream(stream)
	must.NoError(t, err)
	err = ts.store.PutRecords("test", record)
	must.NoError(t, err)
	ts.db.View(func(tx *badger.Txn) error { //nolint:errcheck
		item, err := tx.Get(recordKey("test", 0, record.InternalId))
		must.NoError(t, err)
		must.Eq(t, uint64(now.Unix())+1, item.ExpiresAt())
		return nil
	})
}

func TestPutRecords_Expired(t *testing.T) {
	_, ok := os.LookupEnv("RUN_LONG_TESTS")
	if !ok {
		t.SkipNow()
	}
	ts := setupTest(t)
	stream := &kayakv1.Stream{
		Name:           "test",
		PartitionCount: 1,
		Ttl:            1,
	}
	record := &kayakv1.Record{
		StreamName: "test",
		InternalId: ulid.Make().String(),
		Id:         []byte("test"),
		Partition:  0,
	}
	err := ts.store.PutStream(stream)
	must.NoError(t, err)
	err = ts.store.PutRecords("test", record)
	must.NoError(t, err)
	ts.db.View(func(tx *badger.Txn) error { //nolint:errcheck
		time.Sleep(1 * time.Second)
		_, err := tx.Get(recordKey("test", 0, record.InternalId))
		must.ErrorIs(t, err, badger.ErrKeyNotFound)
		return nil
	})
}

func TestGetRecords(t *testing.T) {
	firstId := ulid.Make().String()
	time.Sleep(2 * time.Millisecond)
	secondId := ulid.Make().String()

	cases := []struct {
		name      string
		records   []*kayakv1.Record
		limit     int
		expected  []*kayakv1.Record
		postition string
	}{
		{
			name:  "all records",
			limit: 5,
			records: []*kayakv1.Record{
				{
					Partition:  0,
					InternalId: firstId,
					Id:         []byte(firstId),
					Payload:    []byte("test1"),
				},
				{
					Partition:  0,
					InternalId: secondId,
					Id:         []byte(secondId),
					Payload:    []byte("test2"),
				},
			},
			expected: []*kayakv1.Record{
				{
					Partition:  0,
					InternalId: firstId,
					Id:         []byte(firstId),
					Payload:    []byte("test1"),
					StreamName: "test",
				},
				{
					Partition:  0,
					InternalId: secondId,
					Id:         []byte(secondId),
					Payload:    []byte("test2"),
					StreamName: "test",
				},
			},
		},
		{
			name:  "one record",
			limit: 1,
			records: []*kayakv1.Record{
				{
					Partition:  0,
					InternalId: firstId,
					Id:         []byte(firstId),
					Payload:    []byte("test1"),
				},
				{
					Partition:  0,
					InternalId: secondId,
					Id:         []byte(secondId),
					Payload:    []byte("test2"),
				},
			},
			expected: []*kayakv1.Record{
				{
					Partition:  0,
					InternalId: firstId,
					Id:         []byte(firstId),
					Payload:    []byte("test1"),
					StreamName: "test",
				},
			},
		},
		{
			name:      "use position",
			limit:     5,
			postition: firstId,
			records: []*kayakv1.Record{
				{
					Partition:  0,
					InternalId: firstId,
					Id:         []byte(firstId),
					Payload:    []byte("test1"),
				},
				{
					Partition:  0,
					InternalId: secondId,
					Id:         []byte(secondId),
					Payload:    []byte("test2"),
				},
			},
			expected: []*kayakv1.Record{
				{
					Partition:  0,
					InternalId: secondId,
					Id:         []byte(secondId),
					Payload:    []byte("test2"),
					StreamName: "test",
				},
			},
		},
		{
			name:  "only correct partition",
			limit: 5,
			records: []*kayakv1.Record{
				{
					Partition:  1,
					InternalId: firstId,
					Id:         []byte(firstId),
					Payload:    []byte("test1"),
				},
				{
					Partition:  0,
					InternalId: secondId,
					Id:         []byte(secondId),
					Payload:    []byte("test2"),
				},
			},
			expected: []*kayakv1.Record{
				{
					Partition:  0,
					InternalId: secondId,
					Id:         []byte(secondId),
					Payload:    []byte("test2"),
					StreamName: "test",
				},
			},
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			ts := setupTest(t)
			ts.store.PutStream(&kayakv1.Stream{ //nolint:errcheck
				Name:           "test",
				PartitionCount: 1,
			})
			ts.store.PutRecords("test", tc.records...) //nolint:errcheck

			actual, err := ts.store.GetRecords("test", 0, tc.postition, tc.limit)
			must.NoError(t, err)
			protoEq(t, tc.expected, actual)
		})
	}
}

func TestGetPartitionAssignments(t *testing.T) {
	ts := setupTest(t)
	ts.db.Update(func(tx *badger.Txn) error { //nolint:errcheck
		tx.SetEntry(badger.NewEntry([]byte("registrations:test:group1:0"), []byte("worker1")).WithTTL(1 * time.Second))  //nolint:errcheck
		tx.SetEntry(badger.NewEntry([]byte("registrations:test:group1:1"), []byte("worker2")).WithTTL(-1 * time.Second)) //nolint:errcheck
		tx.SetEntry(badger.NewEntry([]byte("registrations:test:group1:2"), []byte("worker3")).WithTTL(1 * time.Second))  //nolint:errcheck
		tx.SetEntry(badger.NewEntry([]byte("registrations:test:group2:3"), []byte("worker4")).WithTTL(1 * time.Second))  //nolint:errcheck
		return nil
	})
	expires := time.Now().Add(1 * time.Second).Unix()
	expected := map[int64]*kayakv1.PartitionAssignment{
		0: {
			StreamName: "test",
			GroupName:  "group1",
			Partition:  0,
			WorkerId:   "worker1",
			ExpiresAt:  expires,
		},
		2: {
			StreamName: "test",
			GroupName:  "group1",
			Partition:  2,
			WorkerId:   "worker3",
			ExpiresAt:  expires,
		},
	}
	assignments, _ := ts.store.GetPartitionAssignments("test", "group1")

	protoEq(t, expected, assignments)
}
func TestGetPartitionAssignment(t *testing.T) {
	ts := setupTest(t)
	// no assignment
	id, err := ts.store.GetPartitionAssignment("test", "group1", 0)
	must.NoError(t, err)
	must.Eq(t, "", id)

	ts.db.Update(func(tx *badger.Txn) error { //nolint:errcheck
		tx.Set(partitionAssignmentKey("test", "group1", 0), []byte("worker1")) //nolint:errcheck
		return nil
	})

	id, err = ts.store.GetPartitionAssignment("test", "group1", 0)
	must.NoError(t, err)
	must.Eq(t, "worker1", id)
}
func TestExtendLease(t *testing.T) {
	worker := &kayakv1.Worker{
		StreamName:          "test",
		GroupName:           "group",
		Id:                  "worker1",
		PartitionAssignment: 0,
		LeaseExpires:        time.Now().Add(10 * time.Second).Unix(),
	}
	cases := []struct {
		name           string
		existingWorker *kayakv1.Worker
		err            error
	}{
		{
			name: "no lease",
		},
		{
			name: "existing lease different worker",
			existingWorker: &kayakv1.Worker{
				StreamName:          "test",
				GroupName:           "group",
				Id:                  "worker2",
				PartitionAssignment: 0,
				LeaseExpires:        time.Now().Add(10 * time.Second).Unix(),
			},
			err: ErrAlreadyRegistered,
		},
		{
			name: "existing lease same worker",
			existingWorker: &kayakv1.Worker{
				StreamName:          "test",
				GroupName:           "group",
				Id:                  "worker1",
				PartitionAssignment: 0,
				LeaseExpires:        time.Now().Add(2 * time.Second).Unix(),
			},
		},
		{
			name: "expired lease",
			existingWorker: &kayakv1.Worker{
				StreamName:          "test",
				GroupName:           "group",
				Id:                  "worker2",
				PartitionAssignment: 0,
				LeaseExpires:        time.Now().Add(-2 * time.Second).Unix(),
			},
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			ts := setupTest(t)
			if tc.existingWorker != nil {
				ts.db.Update(func(tx *badger.Txn) error { //nolint:errcheck
					wts := time.Unix(tc.existingWorker.LeaseExpires, 0)
					ttl := time.Until(wts)
					tx.SetEntry( //nolint:errcheck
						badger.NewEntry(
							partitionAssignmentKey(tc.existingWorker.StreamName, tc.existingWorker.GroupName, tc.existingWorker.PartitionAssignment),
							[]byte(tc.existingWorker.Id),
						).WithTTL(ttl),
					)
					return nil
				})
			}
			expires := time.Until(time.Unix(worker.LeaseExpires, 0))
			err := ts.store.ExtendLease(worker, expires)
			must.ErrorIs(t, err, tc.err)
			expected := worker
			if err != nil {
				expected = tc.existingWorker
			}
			ts.db.View(func(tx *badger.Txn) error { //nolint:errcheck
				item, err := tx.Get(partitionAssignmentKey(expected.StreamName, expected.GroupName, expected.PartitionAssignment))
				must.NoError(t, err)
				must.Eq(t, uint64(expected.LeaseExpires), item.ExpiresAt())
				raw, _ := item.ValueCopy(nil)
				must.Eq(t, expected.Id, string(raw))
				return nil
			})
		})
	}
}

func TestHasLease(t *testing.T) {
	ts := setupTest(t)
	err := ts.store.HasLease(&kayakv1.Worker{StreamName: "test", GroupName: "group", PartitionAssignment: 1, Id: "worker1"})
	must.ErrorIs(t, err, ErrInvalidLease)
	ts.db.Update(func(tx *badger.Txn) error { //nolint:errcheck
		tx.Set(partitionAssignmentKey("test", "group", 1), []byte("worker1")) //nolint:errcheck
		return nil
	})
	err = ts.store.HasLease(&kayakv1.Worker{StreamName: "test", GroupName: "group", PartitionAssignment: 1, Id: "worker1"})
	must.NoError(t, err)
	err = ts.store.HasLease(&kayakv1.Worker{StreamName: "test", GroupName: "group", PartitionAssignment: 1, Id: "worker2"})
	must.ErrorIs(t, err, ErrInvalidLease)
}

func TestRemoveLease(t *testing.T) {
	ts := setupTest(t)
	worker := &kayakv1.Worker{StreamName: "test", GroupName: "group", PartitionAssignment: 1, Id: "worker1"}
	ts.store.ExtendLease(worker, 1*time.Second) //nolint:errcheck
	must.NoError(t, ts.store.HasLease(worker))

	err := ts.store.RemoveLease(worker)
	must.NoError(t, err)

	must.ErrorIs(t, ts.store.HasLease(worker), ErrInvalidLease)
}

func TestCommitGroupPosition(t *testing.T) {
	ts := setupTest(t)
	ulidFirst := ulid.Make()
	time.Sleep(1 * time.Millisecond)
	ulidSecond := ulid.Make()
	err := ts.store.CommitGroupPosition("stream1", "group1", 0, ulidSecond.String())
	must.NoError(t, err)

	ts.db.View(func(tx *badger.Txn) error { //nolint:errcheck
		item, _ := tx.Get(groupPositionKey("stream1", "group1", 0))
		val, _ := item.ValueCopy(nil)
		must.Eq(t, ulidSecond.String(), string(val))
		return nil
	})
	err = ts.store.CommitGroupPosition("stream1", "group1", 0, ulidFirst.String())
	must.ErrorIs(t, err, ErrNewerPositionCommitted)
}
func TestGetWorkerPosition(t *testing.T) {
	ts := setupTest(t)
	ul := ulid.Make()
	err := ts.store.CommitGroupPosition("stream1", "group1", 0, ul.String())
	must.NoError(t, err)

	value, err := ts.store.GetGroupPosition("stream1", "group1", 0)
	must.NoError(t, err)
	must.Eq(t, ul.String(), value)

	value, err = ts.store.GetGroupPosition("stream1", "group1", 1)
	must.NoError(t, err)
	must.Eq(t, "", value)
}

func protoEq[V any](t *testing.T, expected, actual V) {
	t.Helper()
	must.Eq(t, expected, actual, must.Cmp(protocmp.Transform()))
}
