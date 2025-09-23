package fsm

import (
	"errors"
	"testing"

	"github.com/dgraph-io/badger/v4"
	"github.com/hashicorp/raft"
	"github.com/shoenig/test/must"
)

func TestPersist(t *testing.T) {
	mockedDB := NewMockSnapshotDB(t)
	bs := &badgerFsmSnapshot{db: mockedDB}
	sink := &raft.DiscardSnapshotSink{}

	mockedDB.EXPECT().Backup(sink, uint64(0)).Return(0, nil).Once()
	err := bs.Persist(sink)
	must.NoError(t, err)
}

func TestPersist_Error(t *testing.T) {

	mockedDB := NewMockSnapshotDB(t)
	bs := &badgerFsmSnapshot{db: mockedDB}
	sink := &raft.DiscardSnapshotSink{}
	expectedErr := errors.New("oops")

	mockedDB.EXPECT().Backup(sink, uint64(0)).Return(0, expectedErr).Once()
	err := bs.Persist(sink)
	must.ErrorIs(t, err, expectedErr)
}

func TestRelease(t *testing.T) {
	mockedDB := NewMockSnapshotDB(t)
	bs := &badgerFsmSnapshot{db: mockedDB}
	bs.Release()
}

func TestNewFSMSnapshot(t *testing.T) {
	db, err := badger.Open(badger.DefaultOptions("").WithInMemory(true))
	must.NoError(t, err)
	snp := NewFSMSnapshot(db)
	casted, ok := snp.db.(*badger.DB)
	must.True(t, ok)
	must.Eq(t, db, casted)
	defer db.Close()
}
