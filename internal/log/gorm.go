package log

import (
	"context"
	"errors"
	"fmt"
	"time"

	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"log/slog"
)

type GormLogger struct{}

func (g *GormLogger) LogMode(level logger.LogLevel) logger.Interface {
	return g
}
func (g *GormLogger) Info(ctx context.Context, msg string, data ...interface{}) {
	slog.InfoContext(ctx, fmt.Sprintf(msg, data...))
}
func (g *GormLogger) Warn(ctx context.Context, msg string, data ...interface{}) {
	slog.WarnContext(ctx, fmt.Sprintf(msg, data...))
}
func (g *GormLogger) Error(ctx context.Context, msg string, data ...interface{}) {
	slog.ErrorContext(ctx, fmt.Sprintf(msg, data...))
}
func (g *GormLogger) Trace(ctx context.Context, begin time.Time, fc func() (sql string, rowsAffected int64), err error) {
	sql, rows := fc()
	elapsed := time.Since(begin)

	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		slog.Error("query error", "error", err, "count", rows, "elapsed", float64(elapsed.Nanoseconds())/1e6)
	}
	if elapsed > 500*time.Millisecond {
		slog.Warn("slow query", "query", sql, "elapsed", float64(elapsed.Nanoseconds())/1e6)
	}
}

/*
type Interface interface {
	LogMode(LogLevel) Interface
	Info(context.Context, string, ...interface{})
	Warn(context.Context, string, ...interface{})
	Error(context.Context, string, ...interface{})
	Trace(ctx context.Context, begin time.Time, fc func() (sql string, rowsAffected int64), err error)
}
*/
