package main

import (
	"context"
	"errors"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"
	"time"

	"github.com/glebarez/sqlite"
	"github.com/lmittmann/tint"
	"github.com/urfave/cli/v2"
	"github.com/urfave/cli/v2/altsrc"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/exporters/prometheus"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/resource"
	"go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.21.0"
	"gorm.io/gorm"
	"log/slog"

	"github.com/binarymatt/kayak/internal/config"
	"github.com/binarymatt/kayak/internal/log"
	"github.com/binarymatt/kayak/internal/service"
	"github.com/binarymatt/kayak/internal/store"
)

func main() {
	flags := []cli.Flag{
		altsrc.NewStringFlag(&cli.StringFlag{
			Name:     "id",
			Required: true,
		}),
		altsrc.NewStringFlag(&cli.StringFlag{
			Name:  "host",
			Value: "0.0.0.0",
		}),
		altsrc.NewIntFlag(&cli.IntFlag{
			Name:  "port",
			Value: 8080,
		}),
		altsrc.NewIntFlag(&cli.IntFlag{
			Name:  "serf_port",
			Value: 9000,
		}),
		altsrc.NewBoolFlag(&cli.BoolFlag{
			Name:     "bootstrap",
			Required: false,
			Value:    false,
		}),
		altsrc.NewBoolFlag(&cli.BoolFlag{
			Name:     "console",
			Required: false,
			Value:    false,
		}),
		altsrc.NewStringFlag(&cli.StringFlag{
			Name:     "dir",
			Required: true,
		}),
		altsrc.NewStringSliceFlag(&cli.StringSliceFlag{
			Name: "peers",
		}),
		&cli.StringFlag{
			Name: "config",
		},
		altsrc.NewDurationFlag(&cli.DurationFlag{
			Name:     "stats_timer",
			Required: false,
			Value:    15 * time.Second,
		}),
		altsrc.NewDurationFlag(&cli.DurationFlag{
			Name:     "background_timer",
			Required: false,
			Value:    5 * time.Second,
		}),
		altsrc.NewStringFlag(&cli.StringFlag{
			Name:    "log_level",
			Aliases: []string{"L"},
			Value:   "info",
		}),
	}

	app := &cli.App{
		Action: kayakRun,
		Before: altsrc.InitInputSourceWithContext(flags, altsrc.NewYamlSourceFromFlagFunc("config")),
		Flags:  flags,
	}

	if err := app.Run(os.Args); err != nil {
		slog.Error("error running kayak", "error", err)
	}

}
func setupLogging(console bool, level slog.Leveler) {
	var logger *slog.Logger
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
func initializeDataDir(path string) error {
	if _, err := os.Stat(path); errors.Is(err, os.ErrNotExist) {
		err := os.MkdirAll(path, os.ModePerm)
		if err != nil {
			slog.Error("error ensuring directory", "error", err, "path", path)
			return err
		}
	}
	return nil
}
func initProvider(ctx context.Context) (func(), error) {
	client := otlptracegrpc.NewClient()
	exp, err := otlptrace.New(ctx, client)
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
func kayakRun(cctx *cli.Context) error {
	ctx := cctx.Context
	cfg := config.New(cctx)
	var level slog.Leveler
	switch cctx.String("log_level") {
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

	setupLogging(cctx.Bool("console"), level)
	if err := setupMetrics(); err != nil {
		return err
	}

	closer, err := initProvider(ctx)
	if err != nil {
		return err
	}
	defer closer()

	path := filepath.Join(cfg.DataDir, cfg.ServerID)
	slog.Info("setting up kayak run", "id", cfg.ServerID, "path", path)
	if err := initializeDataDir(path); err != nil {
		return err
	}

	ctx, cancel := signal.NotifyContext(
		ctx,
		syscall.SIGINT,
		syscall.SIGTERM,
		syscall.SIGQUIT,
		syscall.SIGKILL,
	)
	defer cancel()
	slog.Warn("opening db")
	// options := badger.DefaultOptions(filepath.Join(path, "badger.db")).
	// WithLogger(&log.BadgerLogger{})
	// db, err := badger.Open(options)

	sdb := sqlite.Open(cfg.DataPath())
	db, err := gorm.Open(sdb, &gorm.Config{
		Logger: &log.GormLogger{},
	})
	if err != nil {
		slog.Error("could not open db file", "error", err)
		return err
	}
	// s := store.NewBadger(db)
	s := store.NewSqlStore(db)
	defer s.Close()

	err = s.RunMigrations()
	if err != nil {
		slog.Error("error running migrations")
		return err
	}

	slog.Debug("creating service")
	svc, err := service.New(s, cfg)
	if err != nil {
		return err
	}
	if err := svc.Init(); err != nil {
		return err
	}
	defer svc.Stop(ctx)

	return svc.Start(ctx)

}
