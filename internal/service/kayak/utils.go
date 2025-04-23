package kayak

import "github.com/spaolacci/murmur3"

// return partition index that this message should live in
func balancer(key []byte, partitionCount int64) int64 {
	partition := murmur3.Sum64(key) % uint64(partitionCount)
	return int64(partition)
}
