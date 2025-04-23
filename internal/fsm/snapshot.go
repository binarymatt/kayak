package fsm

import (
	"log/slog"

	"github.com/dgraph-io/badger/v4"
	"github.com/hashicorp/raft"
)

type badgerFsmSnapshot struct {
	db *badger.DB
}

func (fs *badgerFsmSnapshot) Persist(sink raft.SnapshotSink) error {
	defer sink.Close()
	_, err := fs.db.Backup(sink, 0)
	return err
}
func (fs *badgerFsmSnapshot) Release() {
	slog.Info("releasing fsm badger snapshot")
}

func NewFSMSnapshot(db *badger.DB) *badgerFsmSnapshot {
	return &badgerFsmSnapshot{db: db}
}
