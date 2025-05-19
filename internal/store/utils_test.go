package store

import (
	"testing"

	"github.com/shoenig/test/must"
)

func TestPartitionAssignmentPrefix(t *testing.T) {
	res := partitionAssignmentPrefix("test", "group1")
	must.Eq(t, "registrations:test:group1", string(res))
}

func TestPartitionAssignmentKey(t *testing.T) {
	res := partitionAssignmentKey("test", "group1", 0)
	must.Eq(t, "registrations:test:group1:0", string(res))
}

func TestStreamsKey(t *testing.T) {
	res := streamsKey("test")
	must.Eq(t, "streams:test", string(res))
}

func TestRecordKey(t *testing.T) {
	res := recordKey("test", 0, "1234")
	must.Eq(t, "test:0:1234", string(res))
}

func TestRecordPrefixKey(t *testing.T) {
	res := recordPrefixKey("test", 0)
	must.Eq(t, "test:0", string(res))
}

func TestGroupPositionKey(t *testing.T) {
	res := groupPositionKey("testStream", "testGroup", 1)
	must.Eq(t, "groups:testStream:testGroup:1", string(res))
}
