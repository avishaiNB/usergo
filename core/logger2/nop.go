package logger2

import (
	"github.com/go-kit/kit/log"
)

type nopLogger struct{}

// NewNopLogger returns a log.Logger that doesn't do anything.
// Should be used `for testing only
func NewNopLogger() log.Logger {
	return nopLogger{}
}

func (logger nopLogger) Log(keyvals ...interface{}) error {
	return nil
}
