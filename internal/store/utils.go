package store

import (
	"errors"
)

var (
	ErrIterFinished              = errors.New("ERR iteration finished successfully")
	ErrInvalidTopic              = errors.New("invalid topic")
	ErrTopicAlreadyExists        = errors.New("topic already exists")
	ErrArchivedTopic             = errors.New("topic has been archived")
	ErrTopicRecreate             = errors.New("cannot recreate archived topic")
	ErrMissingBucket             = errors.New("missing topics bucket")
	ErrCreatingTopic             = errors.New("could not create topic")
	ErrInvalidConsumer           = errors.New("invalid consumer for group")
	ErrConsumerAlreadyRegistered = errors.New("consumer is already registered")
)

type KVItem struct {
	Bucket []byte
	Key    []byte
	Value  []byte
	Err    error
}

func (i *KVItem) IsFinished() bool {
	return i.Err == ErrIterFinished
}
