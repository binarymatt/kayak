package store

import (
	"errors"
	"fmt"
	"io"
	"log/slog"
	"strconv"
	"strings"
	"time"

	"github.com/dgraph-io/badger/v4"
	"github.com/hashicorp/raft"
	"github.com/oklog/ulid/v2"
	"google.golang.org/protobuf/proto"

	kayakv1 "github.com/binarymatt/kayak/gen/kayak/v1"
)

type Store interface {
	PutStream(stream *kayakv1.Stream) error
	GetStream(name string) (*kayakv1.Stream, error)
	GetStreams() ([]*kayakv1.Stream, error)
	DeleteStream(name string) error

	PutRecords(streamName string, records ...*kayakv1.Record) error
	GetRecords(streamName string, partition int64, startPosition string, limit int) ([]*kayakv1.Record, error)

	// GetLease()
	ExtendLease(worker *kayakv1.Worker, expires time.Duration) error
	RemoveLease(worker *kayakv1.Worker) error
	HasLease(worker *kayakv1.Worker) error
	GetGroupPosition(stream, group string, partition int64) (string, error)
	CommitGroupPosition(stream, group string, parition int64, position string) error
	GetPartitionAssignment(stream, group string, partition int64) (string, error)
	GetPartitionAssignments(stream, group string) (map[int64]*kayakv1.PartitionAssignment, error)
	GetGroupInformation(streamName, groupName string) (*kayakv1.Group, error)
	GetStreamStats(name string) (*kayakv1.StreamStats, error)

	// raft FSM
	Apply(l *raft.Log) any
	Snapshot() (raft.FSMSnapshot, error)
	Restore(snapshot io.ReadCloser) error
}

var (
	_ Store = (*store)(nil)
)

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

	var stream kayakv1.Stream
	if err := proto.Unmarshal(rawData, &stream); err != nil {
		return nil, err
	}
	return &stream, nil
}

func (s *store) GetStreams() ([]*kayakv1.Stream, error) {
	txn := s.db.NewTransaction(false)
	defer txn.Discard()
	it := txn.NewIterator(badger.DefaultIteratorOptions)
	defer it.Close()
	streams := []*kayakv1.Stream{}
	prefix := []byte("streams:")
	for it.Seek(prefix); it.ValidForPrefix(prefix); it.Next() {
		item := it.Item()
		val, err := item.ValueCopy(nil)
		if err != nil {
			return nil, err
		}

		var stream kayakv1.Stream
		if err := proto.Unmarshal(val, &stream); err != nil {
			return nil, err
		}
		streams = append(streams, &stream)
	}
	return streams, nil
}

func (s *store) PutRecords(streamName string, records ...*kayakv1.Record) error {
	stream, err := s.GetStream(streamName)
	if err != nil {
		return fmt.Errorf("missing stream %v: %w", streamName, err)
	}
	txn := s.db.NewTransaction(true)
	defer txn.Discard()

	// batch := s.db.NewBatchWithSize(len(records))
	for _, record := range records {
		record.StreamName = streamName
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
	p := prefix
	if startPosition != "" {
		p = recordKey(streamName, partition, startPosition)
	}
	options := badger.IteratorOptions{
		PrefetchValues: true,
		// PrefetchSize:   limit,
		Reverse:     false,
		AllVersions: false,
	}
	it := txn.NewIterator(options)
	defer it.Close()
	records := []*kayakv1.Record{}
	i := 0
	for it.Seek(p); it.ValidForPrefix(p); it.Next() {
		p = prefix
		item := it.Item()
		val, err := item.ValueCopy(nil)
		if err != nil {
			return nil, err
		}
		var r kayakv1.Record
		if err := proto.Unmarshal(val, &r); err != nil {
			return nil, err
		}
		if r.InternalId == startPosition {
			continue
		}
		records = append(records, &r)
		i++
		if i >= limit {
			break
		}
	}
	return records, nil
}

func (s *store) GetGroupPosition(stream, group string, partition int64) (string, error) {
	txn := s.db.NewTransaction(false)
	defer txn.Discard()
	positionKey := groupPositionKey(stream, group, partition)
	item, err := txn.Get(positionKey)
	if err != nil {
		if errors.Is(err, badger.ErrKeyNotFound) {
			return "", nil
		}
		return "", err
	}
	payload, err := item.ValueCopy(nil)
	if err != nil {
		return "", err
	}
	return string(payload), nil
}

var ErrNewerPositionCommitted = errors.New("a newer position has already been commited")

func (s *store) CommitGroupPosition(stream, group string, partition int64, position string) error {
	// only commit if the current position is lower
	new, err := ulid.Parse(position)
	if err != nil {
		return err
	}
	return s.db.Update(func(tx *badger.Txn) error {
		positionKey := groupPositionKey(stream, group, partition)
		item, err := tx.Get(positionKey)
		if err != nil && !errors.Is(err, badger.ErrKeyNotFound) {
			slog.Error("error getting group positiong during commit", "error", err)
			return err
		}
		if item != nil {
			val, err := item.ValueCopy(nil)
			if err != nil {
				return err
			}
			existing, err := ulid.Parse(string(val))
			if err != nil {
				return err
			}
			if new.Compare(existing) < 1 {
				return ErrNewerPositionCommitted
			}
		}

		return tx.Set(positionKey, []byte(position))
	})
}

func (s *store) GetPartitionAssignment(stream, group string, partition int64) (string, error) {
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
var ErrInvalidLease = errors.New("invalid lease")

func (s *store) HasLease(worker *kayakv1.Worker) error {
	assignedId, err := s.GetPartitionAssignment(worker.StreamName, worker.GroupName, worker.PartitionAssignment)
	if err != nil {
		return err
	}
	if assignedId != worker.Id {
		return ErrInvalidLease
	}
	return nil
}
func (s *store) ExtendLease(worker *kayakv1.Worker, expires time.Duration) error {
	assignedId, err := s.GetPartitionAssignment(worker.StreamName, worker.GroupName, worker.PartitionAssignment)
	if err != nil {
		return err
	}
	if assignedId != "" && assignedId != worker.Id {
		return ErrAlreadyRegistered
	}

	return s.db.Update(func(tx *badger.Txn) error {
		key := partitionAssignmentKey(worker.StreamName, worker.GroupName, worker.PartitionAssignment)
		entry := badger.NewEntry(key, []byte(worker.Id)).WithTTL(expires)
		return tx.SetEntry(entry)
	})
}
func (s *store) RemoveLease(worker *kayakv1.Worker) error {
	assignedId, err := s.GetPartitionAssignment(worker.StreamName, worker.GroupName, worker.PartitionAssignment)
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

func (s *store) GetPartitionAssignments(stream, group string) (map[int64]*kayakv1.PartitionAssignment, error) {
	txn := s.db.NewTransaction(false)
	defer txn.Discard()

	prefix := partitionAssignmentPrefix(stream, group)

	it := txn.NewIterator(badger.DefaultIteratorOptions)
	defer it.Close()
	assignments := map[int64]*kayakv1.PartitionAssignment{}
	for it.Seek(prefix); it.ValidForPrefix(prefix); it.Next() {
		item := it.Item()
		k := item.Key()
		kvals := strings.Split(string(k), ":")
		partitionStr := kvals[3]
		n, err := strconv.ParseInt(partitionStr, 10, 64)
		if err != nil {
			return nil, err
		}
		assignment := &kayakv1.PartitionAssignment{
			StreamName: stream,
			GroupName:  group,
			Partition:  n,
			ExpiresAt:  int64(item.ExpiresAt()),
		}
		if err := item.Value(func(v []byte) error {
			assignment.WorkerId = string(v)
			return nil
		}); err != nil {
			return nil, err
		}
		assignments[n] = assignment
	}

	return assignments, nil
}
func (s *store) getPartitionCounts(stream string) (map[int64]int64, error) {

	txn := s.db.NewTransaction(false)
	defer txn.Discard()
	prefix := []byte(stream)
	opts := badger.DefaultIteratorOptions
	opts.PrefetchValues = false
	it := txn.NewIterator(opts)
	defer it.Close()

	partitionMapping := map[int64]int64{}
	for it.Seek(prefix); it.ValidForPrefix(prefix); it.Next() {
		item := it.Item()
		key := item.Key()
		parts := strings.Split(string(key), ":")
		partitionStr := parts[1]
		partition, err := strconv.ParseInt(partitionStr, 10, 64)
		if err != nil {
			return nil, err
		}
		partitionMapping[partition]++
		// split key by
	}
	return partitionMapping, nil
}
func (s *store) GetStreamStats(name string) (*kayakv1.StreamStats, error) {
	m, err := s.getPartitionCounts(name)
	if err != nil {
		return nil, err
	}
	groups, err := s.getStreamGroups(name)
	if err != nil {
		return nil, err
	}
	stats := &kayakv1.StreamStats{
		PartitionCounts: m,
		Groups:          groups,
	}
	return stats, nil
}
func (s *store) getStreamGroups(streamName string) ([]*kayakv1.Group, error) {
	prefix := []byte(groupPrefix(streamName))

	txn := s.db.NewTransaction(false)
	defer txn.Discard()
	opts := badger.DefaultIteratorOptions
	// opts.PrefetchValues = false
	it := txn.NewIterator(opts)
	defer it.Close()

	container := map[string]*kayakv1.Group{}
	for it.Seek(prefix); it.ValidForPrefix(prefix); it.Next() {
		item := it.Item()
		key := item.Key()
		parts := strings.Split(string(key), ":")
		groupName := parts[2]
		partitionStr := parts[3]
		partitionNumber, err := strconv.ParseInt(partitionStr, 10, 64)
		if err != nil {
			return nil, err
		}
		val, err := item.ValueCopy(nil)
		if err != nil {
			return nil, err
		}
		position := string(val)
		group, ok := container[groupName]
		if !ok {
			group = &kayakv1.Group{
				StreamName:         streamName,
				Name:               groupName,
				PartitionPositions: map[int64]string{},
			}
		}
		group.PartitionPositions[partitionNumber] = position
		container[groupName] = group
	}
	groups := []*kayakv1.Group{}
	for _, v := range container {
		groups = append(groups, v)
	}
	return groups, nil
}
func (s *store) GetGroupInformation(streamName, groupName string) (*kayakv1.Group, error) {
	stream, err := s.GetStream(streamName)
	if err != nil {
		return nil, err
	}
	group := &kayakv1.Group{
		StreamName:         stream.Name,
		Name:               groupName,
		PartitionPositions: map[int64]string{},
	}
	for i := range stream.PartitionCount {
		position, err := s.GetGroupPosition(stream.Name, groupName, i)
		if err != nil {
			return nil, err
		}
		group.PartitionPositions[i] = position
	}
	return group, nil
}

func (s *store) DeleteStream(name string) error {
	return s.db.Update(func(tx *badger.Txn) error {
		// delete stream info
		sKey := streamsKey(name)
		if err := tx.Delete(sKey); err != nil {
			return err
		}
		// delete stream record
		recordPrefix := []byte(fmt.Sprintf("%s:", name))

		opts := badger.DefaultIteratorOptions
		opts.PrefetchValues = false

		recordIt := tx.NewIterator(opts)
		defer recordIt.Close()
		for recordIt.Seek(recordPrefix); recordIt.ValidForPrefix(recordPrefix); recordIt.Next() {
			key := recordIt.Item().Key()
			if err := tx.Delete(key); err != nil {
				return err
			}
		}

		// delete group info
		groupPre := []byte(groupPrefix(name))
		groupIt := tx.NewIterator(opts)
		defer groupIt.Close()
		for groupIt.Seek(groupPre); groupIt.ValidForPrefix(groupPre); groupIt.Next() {
			key := groupIt.Item().Key()

			if err := tx.Delete(key); err != nil {
				return err
			}
		}
		// delete partition assignments
		assignmentPrefix := []byte(fmt.Sprintf("registrations:%s:", name))

		it := tx.NewIterator(opts)
		defer it.Close()
		for it.Seek(assignmentPrefix); it.ValidForPrefix(assignmentPrefix); it.Next() {

			key := groupIt.Item().Key()

			if err := tx.Delete(key); err != nil {
				return err
			}
		}
		return nil
	})

}

func New(db *badger.DB) *store {
	return &store{db: db}
}
