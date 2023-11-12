package test

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/testing/protocmp"
)

func ProtoEqual(t *testing.T, expected, actual proto.Message, msgAndArgs ...interface{}) {
	r := require.New(t)
	r.Empty(cmp.Diff(expected, actual, protocmp.Transform()), msgAndArgs...)
}
