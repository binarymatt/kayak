package store

import (
	"errors"
	"strconv"
	"strings"
	"time"

	"github.com/dgraph-io/badger/v4"
	"github.com/oklog/ulid/v2"
	"google.golang.org/protobuf/proto"

	kayakv1 "github.com/binarymatt/kayak/gen/kayak/v1"
)

type Store interface {
	PutStream(stream *kayakv1.Stream) error
	GetStream(name string) (*kayakv1.Stream, error)
	PutRecords(streamName string, records ...*kayakv1.Record) error
	GetRecords(streamName string, partition int64, startPosition string, limit int) ([]*kayakv1.Record, error)

	ExtendLease(worker *kayakv1.Worker, expires time.Duration) error
	RemoveLease(worker *kayakv1.Worker) error
	GetWorkerPosition(worker *kayakv1.Worker) (string, error)
	CommitWorkerPosition(worker *kayakv1.Worker) error
	PartitionAssignment(stream, group string, partition int64) (string, error)
	GetPartitionAssignments(stream, group string) (map[int64]string, error)
}

var _ Store = (*store)(nil)

type store struct {
	db *badger.DB
}

func (s *store) PutStream(stream *kayakv1.Stream) error {
	return s.db.Update(func(txn *badger.Txn) error {
		key := streamsKey(stream.Name)
		raw, err := proto.Marshal(stream)
		if err != nil {
			return err
		}
		entry := badger.NewEntry(key, raw)
		return txn.SetEntry(entry)
	})
}
func (s *store) GetStream(name string) (*kayakv1.Stream, error) {
	txn := s.db.NewTransaction(false)
	defer txn.Discard()
	item, err := txn.Get(streamsKey(name))
	if err != nil {
		return nil, err
	}
	rawData, err := item.ValueCopy(nil)
	if err != nil {
		return nil, err
	}

	var stream *kayakv1.Stream
	if err := proto.Unmarshal(rawData, stream); err != nil {
		return nil, err
	}
	return stream, nil
}

func (s *store) PutRecords(streamName string, records ...*kayakv1.Record) error {
	stream, err := s.GetStream(streamName)
	if err != nil {
		return err
	}
	txn := s.db.NewTransaction(true)
	defer txn.Discard()

	// batch := s.db.NewBatchWithSize(len(records))
	for _, record := range records {
		marshalledData, err := proto.Marshal(record)
		if err != nil {
			return err
		}
		id, err := ulid.Parse(string(record.InternalId))
		if err != nil {
			return err
		}
		// <stream_name>:<partition_num>:<id>
		key := recordKey(record.StreamName, record.Partition, id.String())
		entry := badger.NewEntry(key, marshalledData)
		if stream.Ttl > 0 {
			entry = entry.WithTTL(time.Second * time.Duration(stream.Ttl))
		}
		if err := txn.SetEntry(entry); err != nil {
			return err
		}
	}
	return txn.Commit()
}

func (s *store) GetRecords(streamName string, partition int64, startPosition string, limit int) ([]*kayakv1.Record, error) {
	txn := s.db.NewTransaction(false)
	defer txn.Discard()

	prefix := recordPrefixKey(streamName, partition)
	start := recordKey(streamName, partition, startPosition)
	if startPosition == "" {
		start = nil
	}
	options := badger.IteratorOptions{
		PrefetchValues: true,
		PrefetchSize:   limit,
		Reverse:        false,
		AllVersions:    false,
	}
	it := txn.NewIterator(options)
	defer it.Close()
	records := []*kayakv1.Record{}
	for it.Seek(start); it.ValidForPrefix(prefix); it.Next() {

		item := it.Item()
		val, err := item.ValueCopy(nil)
		if err != nil {
			return nil, err
		}
		r := &kayakv1.Record{}
		if err := proto.Unmarshal(val, r); err != nil {
			return nil, err
		}
		records = append(records, r)
	}
	return records, nil
}

func (s *store) GetWorkerPosition(worker *kayakv1.Worker) (string, error) {
	txn := s.db.NewTransaction(false)
	defer txn.Discard()
	positionKey := workerPositionKey(worker.StreamName, worker.GroupName, worker.Id)
	item, err := txn.Get(positionKey)
	if err != nil {
		return "", err
	}
	payload, err := item.ValueCopy(nil)
	if err != nil {
		return "", err
	}
	return string(payload), nil
}

func (s *store) CommitWorkerPosition(worker *kayakv1.Worker) error {
	return s.db.Update(func(tx *badger.Txn) error {
		positionKey := workerPositionKey(worker.StreamName, worker.GroupName, worker.Id)
		position := worker.Position

		return tx.Set(positionKey, []byte(position))
	})
}

func (s *store) PartitionAssignment(stream, group string, partition int64) (string, error) {
	assignedId := ""
	err := s.db.View(func(tx *badger.Txn) error {
		item, err := tx.Get(partitionAssignmentKey(stream, group, partition))
		if err != nil {
			if errors.Is(err, badger.ErrKeyNotFound) {
				return nil
			}
			return err
		}
		val, err := item.ValueCopy(nil)
		if err != nil {
			return err
		}
		assignedId = string(val)
		return nil
	})
	return assignedId, err
}

var ErrAlreadyRegistered = errors.New("partition is alreay registered")

func (s *store) ExtendLease(worker *kayakv1.Worker, expires time.Duration) error {
	assignedId, err := s.PartitionAssignment(worker.StreamName, worker.GroupName, worker.PartitionAssignment)
	if err != nil {
		return err
	}
	if assignedId != "" && assignedId != worker.Id {
		return ErrAlreadyRegistered
	}

	return s.db.Update(func(tx *badger.Txn) error {
		key := partitionAssignmentKey(worker.StreamName, worker.GroupName, worker.PartitionAssignment)
		entry := badger.NewEntry(key, nil).WithTTL(expires)
		return tx.SetEntry(entry)
	})
}
func (s *store) RemoveLease(worker *kayakv1.Worker) error {
	assignedId, err := s.PartitionAssignment(worker.StreamName, worker.GroupName, worker.PartitionAssignment)
	if err != nil {
		return err
	}
	if assignedId != worker.Id {
		return nil
	}

	return s.db.Update(func(tx *badger.Txn) error {
		key := partitionAssignmentKey(worker.StreamName, worker.GroupName, worker.PartitionAssignment)
		return tx.Delete(key)
	})
}

func (s *store) GetPartitionAssignments(stream, group string) (map[int64]string, error) {
	txn := s.db.NewTransaction(false)
	defer txn.Discard()

	prefix := partitionAssignmentPrefix(stream, group)

	it := txn.NewIterator(badger.DefaultIteratorOptions)
	defer it.Close()
	assignments := map[int64]string{}
	for it.Seek(prefix); it.ValidForPrefix(prefix); it.Next() {
		item := it.Item()
		k := item.Key()
		kvals := strings.Split(string(k), ":")
		partitionStr := kvals[3]
		n, err := strconv.ParseInt(partitionStr, 10, 64)
		if err != nil {
			return nil, err
		}
		if err := item.Value(func(v []byte) error {
			assignments[n] = string(v)
			return nil
		}); err != nil {
			return nil, err
		}
	}

	return assignments, nil
}

func New(db *badger.DB) *store {
	return &store{db}
}
