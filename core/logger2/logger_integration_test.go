package logger2_test

import (
	"testing"

	tlectx "github.com/thelotter-enterprise/usergo/core/context"
	logger "github.com/thelotter-enterprise/usergo/core/logger2"
)

func TestLogger(t *testing.T) {
	var isIntegrationTest bool
	isIntegrationTest = false
	if isIntegrationTest {
		ctx := tlectx.New()
		l, _ := logger.NewLogger("test", "file logger", logger.InfoLogLevel)

		logger.ErrorWithContext(ctx.Context, l, "this is an error with context", "arg1", "value1")
		logger.WarnWithContext(ctx.Context, l, "this is an warn with context", "arg1", "value1")
		logger.InfoWithContext(ctx.Context, l, "this is an info with context", "arg1", "value1")
		logger.DebugWithContext(ctx.Context, l, "this is a debug with context", "arg1", "value1")

		logger.Error(l, "this is an error", "arg1", "value1")
		logger.Warn(l, "this is an warn", "arg1", "value1")
		logger.Info(l, "this is an info", "arg1", "value1")
		logger.Debug(l, "this is a debug", "arg1", "value1")

		t.Fail()
	}
}
