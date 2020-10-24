package logger

import "context"

// NewNopManager should be used for testing only.
func NewNopManager() Manager {
	return NewLoggerManager(NewNopLogger())
}

type nopLogger struct{}

// NewNopLogger returns a logger that doesn't do anything.
// Should be used for testing only
func NewNopLogger() Logger { return nopLogger{} }

func (logger nopLogger) Log(ctx context.Context, loggerLevel Level, message string, params ...interface{}) error {
	return nil
}
