package store

import "fmt"

func partitionAssignmentPrefix(stream, group string) []byte {
	return []byte(fmt.Sprintf("registrations:%s:%s", stream, group))
}

// partitionAssignmentKey returns key in following pattern registrations:<stream>:<group>:<partition>
func partitionAssignmentKey(stream, group string, partition int64) []byte {
	return []byte(fmt.Sprintf("registrations:%s:%s:%d", stream, group, partition))
}

// streams key returns key for stream information: streams:<name>
func streamsKey(name string) []byte {
	return []byte(fmt.Sprintf("streams:%s", name))
}

// <stream_name>:<partition_num>:<id>
func recordKey(stream string, partition int64, id string) []byte {
	prefix := recordPrefixKey(stream, partition)
	key := fmt.Sprintf("%s:%s", string(prefix), id)
	return []byte(key)
}

// <stream_name>:<partition>
func recordPrefixKey(stream string, partition int64) []byte {
	key := fmt.Sprintf("%s:%d", stream, partition)
	return []byte(key)
}

// stream_name:group_name:worker_id
func workerPositionKey(stream, group, id string) []byte {
	key := fmt.Sprintf("%s:%s:%s", stream, group, id)
	return []byte(key)
}
