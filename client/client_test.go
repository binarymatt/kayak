package client

import (
	"net/http"
	"testing"

	"github.com/shoenig/test/must"
)

func TestNew(t *testing.T) {
	kc := New("localhost")
	must.Eq(t, &config{client: http.DefaultClient}, kc.cfg)
	must.Eq(t, "localhost", kc.baseUrl)
}

func TestNewWithOptions(t *testing.T) {
	kc := New("localhost", WithStream("test"))
	must.Eq(t, &config{client: http.DefaultClient, streamName: "test"}, kc.cfg)
}
