package service

import (
	"context"
	"time"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric"
)

var (
	meter            = otel.GetMeterProvider().Meter("kayak.service")
	latencyHistogram metric.Int64Histogram
)

func init() {
	var err error
	latencyHistogram, err = meter.Int64Histogram("latency",
		metric.WithDescription("The number of rolls by roll value"),
		metric.WithUnit("{roll}"),
	)
	if err != nil {
		panic(err)
	}
}
func RecordLatency(ctx context.Context, endpoint string, start time.Time) {
	opts := metric.WithAttributes(
		attribute.Key("endpoint").String(endpoint),
	)
	dur := time.Since(start)
	latencyHistogram.Record(ctx, dur.Microseconds(), opts)
}
