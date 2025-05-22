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
	"strings"
	"syscall"
	"time"

	"connectrpc.com/connect"
	connectcors "connectrpc.com/cors"
	"connectrpc.com/otelconnect"
	transport "github.com/Jille/raft-grpc-transport"
	"github.com/dgraph-io/badger/v4"
	"github.com/hashicorp/raft"
	raftboltdb "github.com/hashicorp/raft-boltdb/v2"
	"github.com/lmittmann/tint"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/rs/cors"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/exporters/prometheus"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/resource"
	"go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.26.0"
	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"
	"golang.org/x/sync/errgroup"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	kayakv1 "github.com/binarymatt/kayak/gen/kayak/v1"
	"github.com/binarymatt/kayak/gen/kayak/v1/kayakv1connect"
	"github.com/binarymatt/kayak/internal/log"
	"github.com/binarymatt/kayak/internal/service/admin"
	"github.com/binarymatt/kayak/internal/service/kayak"
	"github.com/binarymatt/kayak/internal/store"
)

func withCORS(h http.Handler) http.Handler {
	middleware := cors.New(cors.Options{
		AllowedOrigins: []string{"*"},
		AllowedMethods: connectcors.AllowedMethods(),
		AllowedHeaders: connectcors.AllowedHeaders(),
		ExposedHeaders: connectcors.ExposedHeaders(),
	})
	return middleware.Handler(h)
}

func setupTracing(ctx context.Context) (func(), error) {
	client := otlptracegrpc.NewClient()
	var err error
	var exp trace.SpanExporter
	exp, err = otlptrace.New(ctx, client)
	if err != nil {
		slog.Error("Error initializing trace", "error", err)
		return func() {}, err
	}
	r, err := resource.Merge(
		resource.Default(),
		resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceName("kayak"),
		),
	)
	if err != nil {
		slog.Error("Error creating trace resource", "error", err)
		return func() {}, err
	}

	tp := trace.NewTracerProvider(
		trace.WithBatcher(exp),
		trace.WithResource(r),
	)
	closer := func() {
		if err := tp.Shutdown(ctx); err != nil {
			slog.Error("error shutting down trace provider", "error", err)
		}
		if err := exp.Shutdown(ctx); err != nil {
			slog.Error("error shutting down trace exporter", "error", err)
		}
	}
	otel.SetTracerProvider(tp)
	otel.SetTextMapPropagator(
		propagation.NewCompositeTextMapPropagator(
			propagation.TraceContext{},
			propagation.Baggage{},
		),
	)
	return closer, nil

}

func setupMetrics() error {
	exporter, err := prometheus.New()
	if err != nil {
		return err
	}
	provider := metric.NewMeterProvider(metric.WithReader(exporter))
	otel.SetMeterProvider(provider)
	return nil

}

type config struct {
	NodeId         string
	ListenAddress  string
	GRPCAddress    string
	RaftDataDir    string
	KayakDataDir   string
	Bootstrap      bool
	InMemory       bool
	JoinAddr       string
	ConsoleLogging bool
	LogLevel       string
}

func parseConfig() (*config, error) {
	pflag.String("node_id", "", "")
	pflag.String("raft_address", "127.0.0.1:1200", "ip:port to use for raft communication")
	pflag.String("listen_address", "0.0.0.0:8080", "ip:port to use for kayak service")
	pflag.String("grpc_address", "0.0.0.0:28080", "ip:port to use for grpc")
	pflag.String("raft_data_dir", "./raft_data", "")
	pflag.String("data_dir", "./data", "")
	pflag.Bool("in_memory", true, "")
	pflag.String("join_addr", "", "")
	pflag.Bool("console", true, "log to console")
	pflag.String("log_level", "info", "log level")
	pflag.Parse()
	if err := viper.BindPFlags(pflag.CommandLine); err != nil {
		return nil, err
	}
	nodeId := viper.GetString("node_id")
	if nodeId == "" {
		nodeId = viper.GetString("raft_address")
	}
	advertise := viper.GetString("advertise_address")
	if advertise == "" {
		advertise = viper.GetString("listen_address")
	}
	return &config{
		NodeId:         nodeId,
		ListenAddress:  viper.GetString("listen_address"),
		GRPCAddress:    viper.GetString("grpc_address"),
		RaftDataDir:    viper.GetString("raft_data_dir"),
		KayakDataDir:   viper.GetString("data_dir"),
		Bootstrap:      viper.GetString("join_addr") == "",
		InMemory:       viper.GetBool("in_memory"),
		JoinAddr:       viper.GetString("join_addr"),
		ConsoleLogging: viper.GetBool("console"),
		LogLevel:       viper.GetString("log_level"),
	}, nil
}

func main() {
	ctx := context.TODO()
	cfg, err := parseConfig()
	if err != nil {
		os.Exit(-1)
		return
	}
	setupLogging(cfg.ConsoleLogging, cfg.LogLevel)
	if err := setupMetrics(); err != nil {
		slog.Error("failure to setup metrics", "error", err)
		os.Exit(-1)
		return
	}
	/*
		closer, err := setupTracing(ctx)
		if err != nil {
			slog.Error("failure to setup tracing", "error", err)
			os.Exit(-1)
			return
		}
		defer closer()
	*/
	otelInterceptor, err := otelconnect.NewInterceptor()
	if err != nil {
		slog.Error("failure to initialize otel connect interceptor", "error", err)
		os.Exit(-1)
	}

	ctx, cancel := signal.NotifyContext(
		ctx,
		syscall.SIGINT,
		syscall.SIGTERM,
		syscall.SIGQUIT,
		syscall.SIGKILL,
	)

	defer cancel()

	g := new(errgroup.Group)
	mux := http.NewServeMux()
	mux.Handle("/metrics", promhttp.Handler())
	opts := badger.DefaultOptions(cfg.KayakDataDir).WithLogger(nil)
	db, err := badger.Open(opts)
	if err != nil {
		slog.Error("could not open db", "error", err)
		return
	}
	store := store.New(db)

	_, port, err := net.SplitHostPort(cfg.GRPCAddress)
	if err != nil {
		slog.Error("failed to parse local address", "error", err, "address", cfg.GRPCAddress)
		os.Exit(1)
		return
	}
	sock, err := net.Listen("tcp", fmt.Sprintf(":%s", port))
	if err != nil {
		slog.Error("failed to listen", "error", err)
		os.Exit(1)
	}
	tm := transport.New(raft.ServerAddress(cfg.GRPCAddress), []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())})
	s := grpc.NewServer()
	tm.Register(s)
	ra, err := setupRaft(cfg.NodeId, cfg.RaftDataDir, store, tm.Transport(), cfg.InMemory, cfg.Bootstrap)
	if err != nil {
		slog.Error("could not setup raft", "error", err)
		return
	}
	kayakService := kayak.New(store, ra)
	mux.Handle(kayakv1connect.NewKayakServiceHandler(
		kayakService,
		connect.WithInterceptors(
			log.NewLogInterceptor(),
			otelInterceptor,
		),
	))
	adminService := admin.New(ra)
	mux.Handle(kayakv1connect.NewAdminServiceHandler(
		adminService,
		connect.WithInterceptors(
			log.NewLogInterceptor(),
			otelInterceptor,
		),
	))

	handler := withCORS(mux)
	server := &http.Server{
		Addr:              cfg.ListenAddress,
		ReadHeaderTimeout: 3 * time.Second,
		// Use h2c so we can serve HTTP/2 without TLS.
		Handler: h2c.NewHandler(handler, &http2.Server{}),
	}
	g.Go(func() error {
		slog.Info("listening", "id", cfg.NodeId, "address", cfg.ListenAddress, "join", cfg.JoinAddr != "", "join_addr", cfg.JoinAddr)
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
		if err := s.Serve(sock); err != nil {
			slog.Warn("grpc server error", "error", err)
			return err
		}
		return nil
	})
	g.Go(func() error {
		<-ctx.Done()
		slog.Debug("ctx is done")
		s.GracefulStop()
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
			slog.Info("trying to join", "address", cfg.JoinAddr, "grpc_address", cfg.GRPCAddress, "id", cfg.NodeId)
			leader := fmt.Sprintf("http://%s", cfg.JoinAddr)
			client := kayakv1connect.NewAdminServiceClient(http.DefaultClient, leader)
			_, err := client.AddVoter(ctx, connect.NewRequest(&kayakv1.AddVoterRequest{
				Id:      cfg.NodeId,
				Address: cfg.GRPCAddress,
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

func setupRaft(nodeId, raftDir string, st store.Store, transport raft.Transport, inMemory, bootstrap bool) (*raft.Raft, error) {
	config := raft.DefaultConfig()
	//slogger := slog.Default().With("component", "raft")
	//logger := shclog.New(slogger)
	//config.Logger = logger
	config.LocalID = raft.ServerID(nodeId)

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

func setupLogging(console bool, logLevel string) {
	var logger *slog.Logger
	var level slog.Leveler
	switch strings.ToLower(logLevel) {
	case "info":
		level = slog.LevelInfo
	case "debug":
		level = slog.LevelDebug
	case "warn":
		level = slog.LevelWarn
	case "error":
		level = slog.LevelError
	default:
		level = slog.LevelInfo
	}

	opts := &slog.HandlerOptions{Level: level, AddSource: true}
	if console {
		logger = slog.New(tint.NewHandler(os.Stderr, &tint.Options{
			Level:     level,
			AddSource: true,
		}))
	} else {
		logger = slog.New(slog.NewJSONHandler(os.Stderr, opts))
	}

	slog.SetDefault(logger)
}
