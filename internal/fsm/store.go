package fsm

import (
	"context"
	"encoding/binary"
	"errors"
	"io"

	"github.com/boltdb/bolt"
	"github.com/dgraph-io/badger/v4"
	"github.com/hashicorp/raft"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
	"log/slog"

	kayakv1 "github.com/binarymatt/kayak/gen/kayak/v1"
	"github.com/binarymatt/kayak/internal/store"
	"github.com/binarymatt/kayak/internal/store/models"
)

type storeFSM struct {
	store store.Store
	path  string
}

var (
	ErrMarhsallRecord    = errors.New("could not marshal record")
	ErrTransactionCommit = errors.New("could not commit transaction")
)

type ApplyResponse struct {
	Error error
	Data  interface{}
}

func (s storeFSM) Apply(log *raft.Log) interface{} {
	switch log.Type {
	case raft.LogCommand:
		var command kayakv1.Command
		if err := protojson.Unmarshal(log.Data, &command); err != nil {
			slog.Error("error unmarshal payload for apply", "error", err)
			return nil
		}
		if command.GetPutRecordsRequest() != nil {
			topic := command.GetPutRecordsRequest().Topic
			records := make([]*models.Record, len(command.GetPutRecordsRequest().GetRecords()))
			for i, pr := range command.GetPutRecordsRequest().GetRecords() {
				records[i] = models.RecordFromProto(pr)
				records[i].TopicID = topic
			}
			err := s.store.AddRecords(context.TODO(), topic, records...)
			if err != nil {
				slog.Error("error while adding records in fsm", "error", err)
			}
			return &ApplyResponse{
				Error: err,
			}
		}
		if command.GetCreateTopicRequest() != nil {
			topic := command.GetCreateTopicRequest().Topic
			name := topic.Name
			slog.Info("creating topic in state machine", "topic", name)
			err := s.store.CreateTopic(context.TODO(), models.TopicFromProto(topic))
			if err != nil {
				slog.Error("Error creating topic in state machine", "topic", name, "error", err)
			}
			if errors.Is(err, bolt.ErrBucketExists) {
				slog.Info("bucket exists already, no error")
				err = nil
			}

			return &ApplyResponse{
				Error: err,
			}
		}
		if command.GetCommitRecordRequest() != nil {
			consumer := models.ConsumerFromProto(command.GetCommitRecordRequest().GetConsumer())
			err := s.store.CommitConsumerPosition(context.TODO(), consumer)
			if err != nil {
				slog.Error("error committing group position", "error", err)
			}

			slog.Info("Commited record in fsm", "topic", consumer.TopicID, "consumer", consumer.ID, "position", consumer.Position)
			return &ApplyResponse{
				Error: err,
			}
		}
		if command.GetDeleteTopicRequest() != nil {
			req := command.GetDeleteTopicRequest()
			slog.Info("deleting topic from state machine", "topic", req.Topic.Name, "archived", req.Topic.Archived)
			topic := models.TopicFromProto(req.Topic)

			err := s.store.DeleteTopic(context.TODO(), topic)
			if err != nil {
				slog.Error("Error deleting topic in state machine", "topic", topic, "error", err)
			}
			if errors.Is(err, bolt.ErrBucketNotFound) {
				slog.Info("ignoring not found bucket for delete")
				err = nil
			}
			return &ApplyResponse{
				Error: err,
			}
		}
		if req := command.GetCreateConsumerGroupRequest(); req != nil {
			return &ApplyResponse{
				Error: errors.New("not available."),
			}
		}
		if req := command.GetRegisterConsumerRequest(); req != nil {
			_, err := s.store.RegisterConsumer(context.TODO(), models.ConsumerFromProto(req.Consumer))
			return &ApplyResponse{
				Error: err,
			}
		}

	}

	slog.Warn("not raft log type")
	return nil
}
func (s storeFSM) SnapshotBolt() (raft.FSMSnapshot, error) {
	slog.Warn("generating fsm snapshot")
	return &fsmSnapshot{
		store: s.store,
	}, nil
	//return newSnapshotNoop()
}
func (s storeFSM) SnapshotBadger() (raft.FSMSnapshot, error) {

	db, ok := s.store.Impl().(*badger.DB)
	if !ok {
		return nil, errors.New("not a badger db")
	}
	return &badgerFsmSnapshot{db: db}, nil
}
func (s storeFSM) Snapshot() (raft.FSMSnapshot, error) {
	//_, ok := s.store.Impl().(*bolt.DB)
	//if ok {
	//	return s.SnapshotBolt()
	//}
	//return s.SnapshotBadger()
	return &sqliteFSMSnapshot{dbPath: s.path}, nil
}

// TODO: figure out how to restore while db is open
func (s storeFSM) Restore(rClose io.ReadCloser) error {
	return nil
}
func (s storeFSM) RestoreBader(rClose io.ReadCloser) error {
	db, ok := s.store.Impl().(*badger.DB)
	if !ok {
		slog.Error("could not case store", "method", "Restore")
		return errors.New("invalid store")
	}
	return db.Load(rClose, 1)
}
func (s storeFSM) RestoreBolt(rClose io.ReadCloser) error {
	defer func() {
		if err := rClose.Close(); err != nil {
			slog.Error("RESTORE: close error", "error", err)
		}
		slog.Info("closing reader done")
	}()
	slog.Info("restoring snapshot")
	db, ok := s.store.Impl().(*bolt.DB)
	if !ok {
		slog.Error("could not case store", "method", "Restore")
		return errors.New("invalid store")
	}
	topics := map[string]int{}
	err := db.Update(func(tx *bolt.Tx) error {
		recordCount := 0
		for {
			msg, err := s.Read(rClose)
			if err != nil {
				break
			}
			var item kayakv1.KVItem
			if err := proto.Unmarshal(msg, &item); err != nil {
				slog.Error("could not unmarshal", "error", err)
				return err
			}
			b, err := tx.CreateBucketIfNotExists(item.GetBucket())
			if err != nil {
				return err
			}
			if err := b.Put(item.GetKey(), item.GetValue()); err != nil {
				return err
			}
			topics[string(item.Bucket)] = topics[string(item.Bucket)] + 1
			recordCount = recordCount + 1
		}
		slog.Warn("done with restore", "stats", topics, "record_count", recordCount)
		return nil
	})
	return err
}

func (s storeFSM) Read(r io.Reader) ([]byte, error) {
	buf := make([]byte, 4)
	if _, err := io.ReadFull(r, buf); err != nil {
		return nil, err
	}

	size := binary.LittleEndian.Uint32(buf)

	msg := make([]byte, size)
	if _, err := io.ReadFull(r, msg); err != nil {
		return nil, err
	}

	return msg, nil
}

func NewStore(s store.Store, path string) raft.FSM {
	return &storeFSM{
		store: s,
		path:  path,
	}
}
