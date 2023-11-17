package service

import (
	"context"
	"time"

	"connectrpc.com/connect"
	"go.opentelemetry.io/otel/trace"
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
			logger.InfoContext(ctx, "grpc endpoint called", slog.Group("grpc", slog.String("http_method", httpMethod), slog.String("method", method)))
			//slog.SetDefault(slog.Default().With(slog.Group("grpc", slog.String("http_method", httpMethod), slog.String("method", method))))
			return next(ctx, req)
		})
	}
	return connect.UnaryInterceptorFunc(interceptor)
}

func LatencyInterceptor() connect.UnaryInterceptorFunc {
	interceptor := func(next connect.UnaryFunc) connect.UnaryFunc {
		return connect.UnaryFunc(func(
			ctx context.Context,
			req connect.AnyRequest,
		) (connect.AnyResponse, error) {
			method := req.Spec().Procedure
			start := time.Now()
			defer RecordLatency(ctx, method, start)
			return next(ctx, req)
		})
	}
	return connect.UnaryInterceptorFunc(interceptor)
}
