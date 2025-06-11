package client

import (
	"net/http"
	"testing"
	"time"

	"github.com/shoenig/test/must"
)

func TestWithHttpClient(t *testing.T) {
	config := &config{}
	WithHttpClient(http.DefaultClient)(config)
	must.NotNil(t, config.client)
}

func TestWithWorkerId(t *testing.T) {
	config := &config{}
	WithWorkerId("test")(config)
	must.Eq(t, "test", config.id)
}

func TestWithStream(t *testing.T) {
	config := &config{}
	WithStream("stream")(config)
	must.Eq(t, "stream", config.streamName)
}
func TestWithGroup(t *testing.T) {
	config := &config{}
	WithGroup("group")(config)
	must.Eq(t, "group", config.group)

}
func TestWithTimer(t *testing.T) {
	config := &config{}
	WithTimer(1 * time.Second)(config)
	must.Eq(t, 1*time.Second, config.ticker)
}
