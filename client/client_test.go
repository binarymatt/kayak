package client

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNewConfig(t *testing.T) {
	r := require.New(t)
	cfg := NewConfig("test")
	r.Equal("", cfg.Address)
	r.Equal("test", cfg.ID)
	r.IsType(&http.Client{}, cfg.HTTPClient)
	r.Equal("", cfg.ConsumerID)
	r.Equal("", cfg.Topic)
}
func TestWithAddress(t *testing.T) {
	cfg := &Config{}
	WithAddress("127.0.0.1")(cfg)
	require.Equal(t, "127.0.0.1", cfg.Address)
}
func TestWithTopic(t *testing.T) {
	cfg := &Config{}
	WithTopic("test-topic")(cfg)
	require.Equal(t, "test-topic", cfg.Topic)
}
func TestWithConsumerID(t *testing.T) {
	cfg := &Config{}
	WithConsumerID("testID")(cfg)
	require.Equal(t, "testID", cfg.ConsumerID)
}
func TestWithConsumerGroup(t *testing.T) {
	cfg := &Config{}
	WithConsumerGroup("testGroup")(cfg)
	require.Equal(t, "testGroup", cfg.ConsumerGroup)
}

func TestWithHTTPClient(t *testing.T) {
	cfg := &Config{}
	client := &http.Client{}
	WithHTTPClient(client)(cfg)
	require.Equal(t, client, cfg.HTTPClient)
}
