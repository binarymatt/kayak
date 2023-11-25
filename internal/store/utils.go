package store

import (
	"errors"
	"strings"

	"github.com/oklog/ulid/v2"
	"google.golang.org/protobuf/proto"

	kayakv1 "github.com/binarymatt/kayak/gen/kayak/v1"
)

var (
	ErrIterFinished              = errors.New("ERR iteration finished successfully")
	ErrInvalidTopic              = errors.New("invalid topic")
	ErrTopicAlreadyExists        = errors.New("topic already exists")
	ErrTopicArchived             = errors.New("topic has been archived")
	ErrTopicRecreate             = errors.New("cannot recreate archived topic")
	ErrMissingBucket             = errors.New("missing topics bucket")
	ErrCreatingTopic             = errors.New("could not create topic")
	ErrGroupFull                 = errors.New("no more slots left in consumer group")
	ErrConsumerGroupInvalid      = errors.New("invalid consumer group")
	ErrInvalidConsumer           = errors.New("invalid consumer for group")
	ErrConsumerGroupExists       = errors.New("consumer group already exists")
	ErrConsumerAlreadyRegistered = errors.New("consumer is already registered")
)

func parseID(data []byte) (ulid.ULID, error) {
	id := ulid.ULID{}
	err := id.UnmarshalBinary(data)
	return id, err
}

func encode(record *kayakv1.Record) ([]byte, error) {
	return proto.Marshal(record)
}

func decode(key, value []byte) (*kayakv1.Record, error) {
	var record kayakv1.Record
	err := proto.Unmarshal(value, &record)
	if err != nil {
		return nil, err
	}
	if key == nil {
		return &record, nil
	}
	parsedKey := strings.Split(string(key), "#")[2]
	id := ulid.ULID{}
	err = id.UnmarshalBinary([]byte(parsedKey))
	if err != nil {
		return nil, err
	}
	record.Id = id.String()
	return &record, nil
}
func key(key string) []byte {
	return []byte(key)
}
func bytes(item string) []byte {
	return []byte(item)
}

type KVItem struct {
	Bucket []byte
	Key    []byte
	Value  []byte
	Err    error
}

func (i *KVItem) IsFinished() bool {
	return i.Err == ErrIterFinished
}
