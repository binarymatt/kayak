package service

import (
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/hashicorp/go-hclog"
	"github.com/hashicorp/raft"
	"github.com/oklog/ulid/v2"
	"github.com/stretchr/testify/suite"
	"log/slog"

	"github.com/binarymatt/kayak/internal/config"
	"github.com/binarymatt/kayak/internal/fsm"
	"github.com/binarymatt/kayak/mocks"
)

type ServiceTestSuite struct {
	suite.Suite
	mockStore *mocks.Store
	service   *service
	cluster   Cluster
	testID    ulid.ULID
}
type Cluster interface {
	Close()
}

func (s *ServiceTestSuite) SetupSuite() {
	fmt.Println("setting up suite")
	l := slog.New(slog.NewTextHandler(io.Discard, nil))
	slog.SetDefault(l)
}

func (s *ServiceTestSuite) TearDownTest() {
	s.cluster.Close()
}
func (s *ServiceTestSuite) SetupTest() {
	fmt.Println("setting up test")
	store := mocks.NewStore(s.T())
	hclog.SetDefault(hclog.NewNullLogger())
	conf := raft.DefaultConfig()
	conf.HeartbeatTimeout = 50 * time.Millisecond
	conf.ElectionTimeout = 50 * time.Millisecond
	conf.LeaderLeaseTimeout = 50 * time.Millisecond
	conf.CommitTimeout = 5 * time.Millisecond
	conf.LogOutput = io.Discard
	conf.LogLevel = "ERROR"
	conf.Logger = hclog.NewNullLogger()
	conf.Logger.SetLevel(hclog.Off)

	// c := raft.MakeCluster(1, s.T(), conf)
	finiteStateMachine := fsm.NewStore(store)
	c := raft.MakeClusterCustom(s.T(), &raft.MakeClusterOpts{
		Peers:     1,
		Bootstrap: true,
		Conf:      conf,
		MakeFSMFunc: func() raft.FSM {
			return finiteStateMachine
		},
	})
	r := c.Leader()

	//setup service
	svc, err := New(store, &config.Config{})
	svc.idGenerator = func() ulid.ULID {
		return s.testID
	}
	s.NoError(err)
	svc.raft = r
	s.service = svc
	s.cluster = c
	s.mockStore = store
}

func TestServiceTestSuite(t *testing.T) {
	suite.Run(t, new(ServiceTestSuite))
}

func (s *ServiceTestSuite) TestReady() {
	r := s.Require()
	req := httptest.NewRequest("GET", "http://example.com/ready", nil)
	w := httptest.NewRecorder()
	s.service.Ready(w, req)
	resp := w.Result()
	r.Equal(http.StatusServiceUnavailable, resp.StatusCode)

	s.service.SetReady()
	w2 := httptest.NewRecorder()
	s.service.Ready(w2, req)
	resp2 := w2.Result()
	r.Equal(http.StatusOK, resp2.StatusCode)
}
