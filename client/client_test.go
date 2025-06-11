package client

import (
	"net/http"
	"testing"
	"time"

	"github.com/shoenig/test/must"
)

func TestNew(t *testing.T) {
	kc := New("localhost")
	must.Eq(t, &config{client: http.DefaultClient, ticker: time.Second * 30}, kc.cfg)
	must.Eq(t, "localhost", kc.baseUrl)
}

func TestNewWithOptions(t *testing.T) {
	kc := New("localhost", WithStream("test"))
	must.Eq(t, &config{client: http.DefaultClient, streamName: "test", ticker: time.Second * 30}, kc.cfg)
}
