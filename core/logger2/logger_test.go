package logger2_test

import (
	"testing"

	tlectx "github.com/thelotter-enterprise/usergo/core/context"
	logger "github.com/thelotter-enterprise/usergo/core/logger2"
)

func TestLoggerWithContextReturnNil(t *testing.T) {
	ctx := tlectx.New()
	nl := logger.NewNopLogger()

	logRes := logger.ErrorWithContext(ctx.Context, nl, "this is an error with context")
	if logRes != nil {
		t.Errorf("TestLoggerWithContextReturnNil return error %v", logRes)
	}
}
