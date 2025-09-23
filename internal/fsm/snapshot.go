package fsm

import (
	"fmt"
	"io"
	"log/slog"

	"github.com/dgraph-io/badger/v4"
	"github.com/hashicorp/raft"
)

var _ raft.FSMSnapshot = (*badgerFsmSnapshot)(nil)

type SnapshotDB interface {
	Close() error
	Backup(w io.Writer, since uint64) (uint64, error)
}
type badgerFsmSnapshot struct {
	db SnapshotDB
}

func (fs *badgerFsmSnapshot) Persist(sink raft.SnapshotSink) error {
	defer sink.Close() //nolint: errcheck
	_, err := fs.db.Backup(sink, 0)
	if err != nil {
		return fmt.Errorf("failed to persist: %w", err)
	}
	return nil
}
func (fs *badgerFsmSnapshot) Release() {
	slog.Info("releasing fsm badger snapshot")
}

func NewFSMSnapshot(db *badger.DB) *badgerFsmSnapshot {
	return &badgerFsmSnapshot{db: db}
}
