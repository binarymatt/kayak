package models

import (
	"gorm.io/datatypes"
)

type Topic struct {
	ID        string `gorm:"primaryKey"`
	Name      string `gorm:"unique"`
	Archived  bool
	UpdatedAt int64 `gorm:"autoUpdateTime"`
	CreatedAt int64 `gorm:"autoCreateTime"`
}
type Record struct {
	ID      string `gorm:"primaryKey"`
	TopicID string
	// Headers map[string]string
	Headers   datatypes.JSONMap
	Payload   []byte
	UpdatedAt int64 `gorm:"autoUpdateTime"`
	CreatedAt int64 `gorm:"autoCreateTime"`
}

func (r *Record) AddHeaders(headers map[string]string) {
	items := make(map[string]interface{}, len(headers))
	for k, v := range headers {
		items[k] = v
	}
	r.Headers = items
}

type ConsumerGroup struct {
	ID             uint
	Name           string `gorm:"uniqueIndex:topic_name"`
	TopicID        string `gorm:"uniqueIndex:topic_name"`
	PartitionCount int64
	Hash           string
	UpdatedAt      int64 `gorm:"autoUpdateTime"`
	CreatedAt      int64 `gorm:"autoCreateTime"`
}

type Consumer struct {
	ID        string `gorm:"primaryKey"`
	GroupID   uint
	Partition int64
	Position  string
	CreatedAt int64 `gorm:"autoCreateTime"`
	UpdatedAt int64 `gorm:"autoUpdateTime"`
}
