package service

import (
	"context"
	"time"

	"connectrpc.com/connect"
	"go.opentelemetry.io/otel/trace"
	"google.golang.org/protobuf/proto"
	"log/slog"

	"github.com/binarymatt/kayak/internal/log"
)

func NewLogInterceptor() connect.UnaryInterceptorFunc {
	interceptor := func(next connect.UnaryFunc) connect.UnaryFunc {
		return connect.UnaryFunc(func(
			ctx context.Context,
			req connect.AnyRequest,
		) (connect.AnyResponse, error) {
			httpMethod := req.HTTPMethod()
			method := req.Spec().Procedure
			scontext := trace.SpanContextFromContext(ctx)
			logger := slog.Default()
			if scontext.HasTraceID() {
				logger = logger.With("trace_id", scontext.TraceID())
			}
			log.WithContext(ctx, logger)
			logger.DebugContext(ctx, "grpc endpoint called", slog.Group("grpc", slog.String("http_method", httpMethod), slog.String("method", method)))
			//slog.SetDefault(slog.Default().With(slog.Group("grpc", slog.String("http_method", httpMethod), slog.String("method", method))))
			return next(ctx, req)
		})
	}
	return connect.UnaryInterceptorFunc(interceptor)
}

func MetricsInterceptor() connect.UnaryInterceptorFunc {
	interceptor := func(next connect.UnaryFunc) connect.UnaryFunc {
		return connect.UnaryFunc(func(
			ctx context.Context,
			req connect.AnyRequest,
		) (connect.AnyResponse, error) {
			start := time.Now()
			endpoint := req.Spec().Procedure

			var reqSize int
			if req != nil {
				if msg, ok := req.Any().(proto.Message); ok {
					reqSize = proto.Size(msg)
				}
			}
			resp, err := next(ctx, req)
			var respSize int
			statusCode := "OK"
			if err == nil {
				if msg, ok := resp.Any().(proto.Message); ok {
					respSize = proto.Size(msg)
				}
			} else {
				statusCode = connect.CodeOf(err).String()
			}
			duration := time.Since(start).Milliseconds()

			requestSize.WithLabelValues(endpoint).Observe(float64(reqSize))
			grpcStarted.WithLabelValues(endpoint).Inc()

			responseSize.WithLabelValues(endpoint, statusCode).Observe(float64(respSize))
			serverLatency.WithLabelValues(endpoint, statusCode).Observe(float64(duration))
			grpcHandled.WithLabelValues(endpoint, statusCode).Inc()

			return resp, err
		})
	}
	return connect.UnaryInterceptorFunc(interceptor)
}
