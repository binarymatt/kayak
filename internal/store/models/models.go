package models

import (
	"math/rand"

	"github.com/spaolacci/murmur3"
	"golang.org/x/exp/slog"
	"gorm.io/datatypes"

	kayakv1 "github.com/binarymatt/kayak/gen/kayak/v1"
)

type Topic struct {
	ID             string `gorm:"primaryKey"` // name
	PartitionCount int
	Archived       bool
	DefaultHash    string
	UpdatedAt      int64 `gorm:"autoUpdateTime"`
	CreatedAt      int64 `gorm:"autoCreateTime"`
	Consumers      []Consumer
	TTL            int64
}

func (t *Topic) RandomBalancer(record *Record) {
	partitions := t.Partitions()
	partition := partitions[rand.Int()%t.PartitionCount]
	record.Partition = partition
}
func (t *Topic) MurmurBalancer(record *Record) {
	idx := murmur3.Sum64([]byte(record.ID)) % uint64(t.PartitionCount)
	partition := t.Partitions()[idx]
	record.Partition = partition
}
func (t *Topic) Balance(record *Record) {
	switch t.DefaultHash {
	case kayakv1.Hash_HASH_MURMUR3.String():
		t.MurmurBalancer(record)
	case kayakv1.Hash_HASH_UNSPECIFIED.String():
		t.RandomBalancer(record)
	default:
		slog.Warn("no balancer found")
	}
}
func (t *Topic) Partitions() []int64 {
	partitions := []int64{}
	for i := 0; i < int(t.PartitionCount); i++ {
		partitions = append(partitions, int64(i))
	}
	return partitions
}

func TopicFromProto(pt *kayakv1.Topic) *Topic {
	topic := &Topic{
		ID:             pt.Name,
		PartitionCount: int(pt.Partitions),
		DefaultHash:    kayakv1.Hash_HASH_MURMUR3.String(),
	}
	return topic
}

type Record struct {
	ID        string `gorm:"primaryKey"`
	TopicID   string
	Partition int64
	// Headers map[string]string
	Headers   datatypes.JSONMap
	Payload   []byte
	UpdatedAt int64 `gorm:"autoUpdateTime"`
	CreatedAt int64 `gorm:"autoCreateTime"`
}

func (r *Record) ToProto() *kayakv1.Record {
	return &kayakv1.Record{
		Id:        r.ID,
		Topic:     r.TopicID,
		Partition: r.Partition,
		Headers:   r.GetHeaders(),
		Payload:   r.Payload,
	}
}
func RecordFromProto(pr *kayakv1.Record) *Record {
	record := &Record{
		Payload:   pr.Payload,
		Partition: pr.Partition,
		TopicID:   pr.Topic,
	}
	record.AddHeaders(pr.Headers)
	if pr.Id != "" {
		record.ID = pr.Id
	}
	return record
}

func (r *Record) GetHeaders() map[string]string {
	m := map[string]string{}
	for k, v := range r.Headers {
		m[k] = v.(string)
	}
	return m
}
func (r *Record) AddHeaders(headers map[string]string) {
	items := make(map[string]interface{}, len(headers))
	for k, v := range headers {
		items[k] = v
	}
	r.Headers = items
}

type Consumer struct {
	ID         string `gorm:"primaryKey"`
	TopicID    string
	Group      string
	Position   string
	CreatedAt  int64 `gorm:"autoCreateTime"`
	UpdatedAt  int64 `gorm:"autoUpdateTime"`
	Partitions datatypes.JSONSlice[int64]
}

func (c *Consumer) GetPartitions() []int64 {
	partitions := []int64{}
	for _, i := range c.Partitions {
		partitions = append(partitions, i)
	}
	return partitions
}
func ConsumerFromProto(pc *kayakv1.TopicConsumer) *Consumer {
	slog.Warn("translating consumer", "partitions", pc.Partitions)
	c := &Consumer{
		ID:         pc.Id,
		TopicID:    pc.Topic,
		Group:      pc.Group,
		Position:   pc.Position,
		Partitions: pc.Partitions,
	}
	return c
}
