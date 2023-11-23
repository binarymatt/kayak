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
	Headers datatypes.JSONMap
	Payload []byte
}

func (r *Record) AddHeaders(headers map[string]string) {
	items := make(map[string]interface{}, len(headers))
	for k, v := range headers {
		items[k] = v
	}
	r.Headers = items
}
