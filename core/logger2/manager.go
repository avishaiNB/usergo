package logger2

import (
	"context"

	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	tlecontext "github.com/thelotter-enterprise/usergo/core/context"
)

// ErrorWithContext ...
func ErrorWithContext(ctx context.Context, logger log.Logger, message string, args ...interface{}) {
	logWithContext(ctx, level.Error(logger), message, args)
}

// WarnWithContext ...
func WarnWithContext(ctx context.Context, logger log.Logger, message string, args ...interface{}) {
	logWithContext(ctx, level.Warn(logger), message, args)
}

// InfoWithContext ...
func InfoWithContext(ctx context.Context, logger log.Logger, message string, args ...interface{}) {
	logWithContext(ctx, level.Info(logger), message, args)
}

// DebugWithContext ...
func DebugWithContext(ctx context.Context, logger log.Logger, message string, args ...interface{}) {
	logWithContext(ctx, level.Debug(logger), message, args)
}

// Error ...
func Error(logger log.Logger, message string, args ...interface{}) {
	l(level.Error(logger), message, args)
}

// Warn ...
func Warn(logger log.Logger, message string, args ...interface{}) {
	l(level.Warn(logger), message, args)
}

// Info ...
func Info(logger log.Logger, message string, args ...interface{}) {
	l(level.Info(logger), message, args)
}

// Debug ...
func Debug(logger log.Logger, message string, args ...interface{}) {
	l(level.Debug(logger), message, args)
}

func logWithContext(ctx context.Context, logger log.Logger, message string, args ...interface{}) {
	correlationID := tlecontext.GetCorrelationFromContext(ctx)
	duration, deadline := tlecontext.GetTimeoutFromContext(ctx)
	logger.Log(
		"message", message,
		"correaltionId", correlationID,
		"duration", duration,
		"deadline", deadline,
		"args", args)
}

func l(logger log.Logger, message string, args ...interface{}) {
	logger.Log(
		"message", message,
		"args", args)
}
