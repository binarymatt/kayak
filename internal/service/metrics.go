package service

import (
	"context"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"gorm.io/gorm"
	"log/slog"

	"github.com/binarymatt/kayak/internal/store"
	"github.com/binarymatt/kayak/internal/store/models"
)

/*
0
5
10
25
50
75
100
250
500
750
1000
2500
5000
7500
10000
*/
var (
	DefaultBytesDistribution        = []float64{1024, 2048, 4096, 16384, 65536, 262144, 1048576, 4194304, 16777216, 67108864, 268435456, 1073741824, 4294967296}
	DefaultMillisecondsDistribution = []float64{0.01, 0.05, 0.1, 0.3, 0.6, 0.8, 1, 2, 3, 4, 5, 6, 8, 10, 13, 16, 20, 25, 30, 40, 50, 65, 80, 100, 130, 160, 200, 250, 300, 400, 500, 650, 800, 1000, 2000, 5000, 10000, 20000, 50000, 100000}
	DefaultMicrosecondsDistribution = []float64{0, 5, 10, 25, 50, 75, 100, 250, 500, 750, 1000, 2500, 5000, 7500, 10000}
	DefaultMessageCountDistribution = []float64{1, 2, 4, 8, 16, 32, 64, 128, 256, 512, 1024, 2048, 4096, 8192, 16384, 32768, 65536}

	serverLatency = promauto.NewHistogramVec(prometheus.HistogramOpts{
		Name:    "endpoint_latency",
		Help:    "endpoint latency measured in milliseconds",
		Buckets: DefaultMillisecondsDistribution,
	}, []string{
		"endpoint",
		"status_code",
	})
	requestSize = promauto.NewHistogramVec(prometheus.HistogramOpts{
		Name:    "request_size",
		Help:    "request size bytes",
		Buckets: DefaultBytesDistribution,
	}, []string{
		"endpoint",
	})
	responseSize = promauto.NewHistogramVec(prometheus.HistogramOpts{
		Name:    "response_size",
		Help:    "response size bytes",
		Buckets: DefaultBytesDistribution,
	}, []string{
		"endpoint",
		"status_code",
	})

	grpcStarted = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "grpc_started",
	}, []string{
		"endpoint",
	})
	grpcHandled = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "grpc_handled",
	}, []string{
		"endpoint",
		"status_code",
	})

	consumerLag = promauto.NewGaugeVec(prometheus.GaugeOpts{
		Name: "consumer_lag",
	}, []string{
		"topic",
		"group",
		"consumer",
	})
	visible_records = promauto.NewGaugeVec(prometheus.GaugeOpts{
		Name: "visible_records",
	}, []string{
		"topic",
	})
)

func recordVisibleCount(topic string, number float64) {
	visible_records.WithLabelValues(topic).Set(number)
}
func recordRequestSize(endpoint string, requestSize int64) {}
func recordDuration(endpoint string, start time.Time) {
	dur := time.Since(start).Microseconds()
	serverLatency.With(prometheus.Labels{
		"endpoint": endpoint,
	}).Observe(float64(dur))
}
func recordResponseSize(endpoint string, responseSize int64) {}
func recordRPCRequest(endpoint string)                       {}
func recordRPCResponse(endpoint string)                      {}

func recordLag(ctx context.Context, topic, group, consumer string, lag int64) {
	slog.Info("calculating consumer lag", "topic", topic, "consumer", consumer)
	consumerLag.WithLabelValues(topic, group, consumer).Set(float64(lag))
}

func CalculateVisibleRecords(ctx context.Context, store store.Store) {
	db, ok := store.Impl().(*gorm.DB)
	if !ok {
		return
	}
	topics, err := store.ListTopics(ctx)
	if err != nil {
		return
	}
	for _, topic := range topics {
		var count int64
		db.Model(&models.Record{}).Where("topic_id = ?", topic).Count(&count)
		recordVisibleCount(topic, float64(count))
	}
}
func CalculateLag(ctx context.Context, store store.Store) {
	topics, err := store.ListTopics(ctx)
	if err != nil {
		return
	}
	for _, topic := range topics {
		meta, err := store.LoadMeta(ctx, topic)
		if err != nil {
			slog.Error("could not get meta information for topic", "topic", topic)
			continue
		}
		slog.Debug("topic info", "topic", meta)
		/*
			for groupName, partitions := range meta.GroupMetadata {
				for _, consumer := range partitions.Consumers {
					lag, err := store.GetConsumerLag(ctx, consumer)
					if err != nil {
						continue
					}
					recordLag(ctx, topic, groupName, consumer.Id, lag)
				}
			}
		*/
	}
}
