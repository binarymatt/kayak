package main

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"
	"time"

	"connectrpc.com/connect"
	"github.com/dgraph-io/badger/v4"
	"github.com/hashicorp/raft"
	raftboltdb "github.com/hashicorp/raft-boltdb/v2"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"
	"golang.org/x/sync/errgroup"

	kayakv1 "github.com/binarymatt/kayak/gen/kayak/v1"
	"github.com/binarymatt/kayak/gen/kayak/v1/kayakv1connect"
	"github.com/binarymatt/kayak/internal/service/admin"
	"github.com/binarymatt/kayak/internal/service/kayak"
	"github.com/binarymatt/kayak/internal/store"
)

func NewLogInterceptor() connect.UnaryInterceptorFunc {
	interceptor := func(next connect.UnaryFunc) connect.UnaryFunc {
		return connect.UnaryFunc(func(
			ctx context.Context,
			req connect.AnyRequest,
		) (connect.AnyResponse, error) {
			httpMethod := req.HTTPMethod()
			method := req.Spec().Procedure
			//scontext := trace.SpanContextFromContext(ctx)
			logger := slog.Default()
			//if scontext.HasTraceID() {
			//	logger = logger.With("trace_id", scontext.TraceID())
			// }
			// log.WithContext(ctx, logger)
			logger.InfoContext(ctx, "grpc endpoint called", slog.Group("grpc", slog.String("http_method", httpMethod), slog.String("method", method)))
			//slog.SetDefault(slog.Default().With(slog.Group("grpc", slog.String("http_method", httpMethod), slog.String("method", method))))
			return next(ctx, req)
		})
	}
	return connect.UnaryInterceptorFunc(interceptor)
}

type config struct {
	NodeId        string
	ListenAddress string
	RaftAddress   string
	RaftDataDir   string
	KayakDataDir  string
	Bootstrap     bool
	InMemory      bool
	JoinAddr      string
}

func parseConfig() (*config, error) {
	pflag.String("node_id", "", "")
	pflag.String("raft_address", "127.0.0.1:1200", "ip:port to use for raft communication")
	pflag.String("listen_address", "0.0.0.0:8080", "ip:port to use for kayak service")
	pflag.String("raft_data_dir", "./raft_data", "")
	pflag.String("data_dir", "./data", "")
	pflag.Bool("in_memory", true, "")
	pflag.String("join_addr", "", "")
	pflag.Parse()
	if err := viper.BindPFlags(pflag.CommandLine); err != nil {
		return nil, err
	}
	nodeId := viper.GetString("node_id")
	if nodeId == "" {
		nodeId = viper.GetString("raft_address")
	}
	return &config{
		NodeId:        nodeId,
		ListenAddress: viper.GetString("listen_address"),
		RaftAddress:   viper.GetString("raft_address"),
		RaftDataDir:   viper.GetString("raft_data_dir"),
		KayakDataDir:  viper.GetString("data_dir"),
		Bootstrap:     viper.GetString("join_addr") == "",
		InMemory:      viper.GetBool("in_memory"),
		JoinAddr:      viper.GetString("join_addr"),
	}, nil
}
func main() {
	cfg, err := parseConfig()
	if err != nil {
		os.Exit(-1)
		return
	}
	ctx, cancel := signal.NotifyContext(
		context.Background(),
		syscall.SIGINT,
		syscall.SIGTERM,
		syscall.SIGQUIT,
		syscall.SIGKILL,
	)

	defer cancel()

	g := new(errgroup.Group)
	mux := http.NewServeMux()
	opts := badger.DefaultOptions(cfg.KayakDataDir).WithLogger(nil)
	db, err := badger.Open(opts)
	if err != nil {
		slog.Error("could not open db", "error", err)
		return
	}
	store := store.New(db)

	ra, err := setupRaft(cfg.NodeId, cfg.RaftAddress, cfg.RaftDataDir, store, cfg.InMemory, cfg.Bootstrap)
	if err != nil {
		slog.Error("could not setup raft", "error", err)
		return
	}
	kayakService := kayak.New(store, ra)
	mux.Handle(kayakv1connect.NewKayakServiceHandler(
		kayakService,
		connect.WithInterceptors(
			NewLogInterceptor(),
		),
	))
	adminService := admin.New(ra)
	mux.Handle(kayakv1connect.NewAdminServiceHandler(
		adminService,
		connect.WithInterceptors(
			NewLogInterceptor(),
		),
	))
	server := &http.Server{
		Addr:              cfg.ListenAddress,
		ReadHeaderTimeout: 3 * time.Second,
		// Use h2c so we can serve HTTP/2 without TLS.
		Handler: h2c.NewHandler(mux, &http2.Server{}),
	}
	g.Go(func() error {
		slog.Info("listening", "cfg", cfg)
		if err := server.ListenAndServe(); err != nil {
			if errors.Is(err, http.ErrServerClosed) {
				slog.Warn("server is closed, exiting")
				return nil
			}
			slog.Error("error running application", "error", err)
			return err
		}
		slog.Info("done serving")
		return nil
	})
	g.Go(func() error {
		<-ctx.Done()
		slog.Debug("ctx is done")
		ct, cancel := context.WithTimeout(context.Background(), time.Second*10)
		defer cancel()
		if err := server.Shutdown(ct); err != nil {
			slog.Error("http server shutdown issue", "error", err)
			return err
		}

		return nil
	})
	if cfg.JoinAddr != "" {
		g.Go(func() error {
			// create client
			time.Sleep(5 * time.Second)
			slog.Info("trying to join", "address", cfg.JoinAddr)
			leader := fmt.Sprintf("http://%s", cfg.JoinAddr)
			client := kayakv1connect.NewAdminServiceClient(http.DefaultClient, leader)
			_, err := client.AddVoter(ctx, connect.NewRequest(&kayakv1.AddVoterRequest{
				Id:      cfg.NodeId,
				Address: cfg.RaftAddress,
			}))
			if err != nil {
				slog.Error("could not join cluster", "error", err)
				return err
			}
			return nil
		})
		// call join endpoint
	}
	if err := g.Wait(); err != nil {
		slog.Error("error from group", "error", err)
	}
	slog.Info("kayak is done")
}

func setupRaft(nodeId, listenAddr, raftDir string, st store.Store, inMemory, bootstrap bool) (*raft.Raft, error) {
	config := raft.DefaultConfig()
	config.LocalID = raft.ServerID(nodeId)

	addr, err := net.ResolveTCPAddr("tcp", listenAddr)
	if err != nil {
		return nil, err
	}
	transport, err := raft.NewTCPTransport(listenAddr, addr, 3, 10*time.Second, os.Stderr)
	if err != nil {
		return nil, err
	}

	// Create the snapshot store. This allows the Raft to truncate the log.
	snapshots, err := raft.NewFileSnapshotStore(raftDir, 1, os.Stderr)
	if err != nil {
		return nil, fmt.Errorf("file snapshot store: %w", err)
	}
	// Create the log store and stable store.
	var logStore raft.LogStore
	var stableStore raft.StableStore
	if inMemory {
		logStore = raft.NewInmemStore()
		stableStore = raft.NewInmemStore()
	} else {
		boltDB, err := raftboltdb.New(raftboltdb.Options{
			Path: filepath.Join(raftDir, "raft.db"),
		})
		if err != nil {
			return nil, fmt.Errorf("new bbolt store: %w", err)
		}
		logStore = boltDB
		stableStore = boltDB
	}
	ra, err := raft.NewRaft(config, st, logStore, stableStore, snapshots, transport)
	if err != nil {
		return nil, err
	}
	if bootstrap {
		configuration := raft.Configuration{
			Servers: []raft.Server{
				{
					ID:      config.LocalID,
					Address: transport.LocalAddr(),
				},
			},
		}
		ra.BootstrapCluster(configuration)
	}
	return ra, nil
}
