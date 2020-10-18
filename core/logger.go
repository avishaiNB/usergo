package core

import (
	"context"
	"os"

	"github.com/go-kit/kit/log"
)

// Log will create a new instance of the Log with ready to use loggers
// Logger should:
// output: log to stdout, stderr, file
// data: each log should append information: correlation id, env name, process name, host name, duration, deadline, logger name, level
// levels: need to check how the levels play here and how we can control them in order to write the log
// read performance related concerns for using file appenders
// TBD: funnel logger
type Log struct {
	Logger log.Logger
	Level  int
}

const (
	// LogLevelCritical ..
	LogLevelCritical int = 5

	// LogLevelError ..
	LogLevelError int = 4

	// LogLevelWarn ..
	LogLevelWarn int = 3

	// LogLevelInfo ..
	LogLevelInfo int = 2

	// LogLevelDebug ..
	LogLevelDebug int = 1
)

// NewLogWithDefaults ..
func NewLogWithDefaults() Log {
	var logger log.Logger
	{
		logger = log.NewLogfmtLogger(os.Stdout)

		logger = log.With(logger, "ts", log.DefaultTimestampUTC)
		logger = log.With(logger, "caller", log.DefaultCaller)
		//log.With(logger, "component", "HTTP")
	}

	logger.Log()
	return NewLog(logger, LogLevelError)
}

// NewLog ...
func NewLog(logger log.Logger, level int) Log {
	return Log{
		Logger: logger,
		Level:  level,
	}
}

// Error will log an error
func (log Log) Error(ctx context.Context, message string, err error, logger log.Logger) bool {
	wasLogged := false

	if log.ShouldLog(LogLevelError) {
		err := logger.Log(
			"level", "error",
			"message", message,
			"error", err,
			// additional important information
		)

		if err == nil {
			wasLogged = true
		}
	}

	return wasLogged
}

// ShouldLog will return a bool indicating if the log message should be logged based on the log level
func (log Log) ShouldLog(level int) bool {
	return level >= log.Level
}
