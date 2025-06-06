package store

import "fmt"

func partitionAssignmentPrefix(stream, group string) []byte {
	return fmt.Appendf(nil, "registrations:%s:%s", stream, group)
}

// partitionAssignmentKey returns key in following pattern registrations:<stream>:<group>:<partition>
func partitionAssignmentKey(stream, group string, partition int64) []byte {
	return fmt.Appendf(nil, "registrations:%s:%s:%d", stream, group, partition)
}

// streams key returns key for stream information: streams:<name>
func streamsKey(name string) []byte {
	return fmt.Appendf(nil, "streams:%s", name)
}

// <stream_name>:<partition_num>:<id>
func recordKey(stream string, partition int64, id string) []byte {
	prefix := recordPrefixKey(stream, partition)
	key := fmt.Appendf(nil, "%s:%s", string(prefix), id)
	return key
}

// <stream_name>:<partition>
func recordPrefixKey(stream string, partition int64) []byte {
	return fmt.Appendf(nil, "%s:%d", stream, partition)
}

// groups:stream_name:group_name:partition
func groupPositionKey(stream, group string, partition int64) []byte {
	prefix := groupPrefix(stream)
	key := fmt.Sprintf("%s:%s:%d", prefix, group, partition)
	return []byte(key)
}

func groupPrefix(streamName string) string {
	return fmt.Sprintf("groups:%s", streamName)
}

// TODO: decide if this will be used
func streamMetadataKey(stream string) []byte { //nolint:unused
	key := fmt.Sprintf("meta:%s", stream)
	return []byte(key)
}
