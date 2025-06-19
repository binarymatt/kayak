package client

import (
	"time"

	"connectrpc.com/connect"
)

func WithHttpClient(client connect.HTTPClient) ConfigOpt {
	return func(c *config) {
		c.client = client
	}
}
func WithWorkerId(id string) ConfigOpt {
	return func(c *config) {
		c.id = id
	}
}
func WithStream(stream string) ConfigOpt {
	return func(c *config) {
		c.streamName = stream
	}
}
func WithGroup(group string) ConfigOpt {
	return func(c *config) {
		c.group = group
	}
}
func WithTimer(d time.Duration) ConfigOpt {
	return func(c *config) {
		c.ticker = d
	}
}
