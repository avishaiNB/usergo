package core

import (
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
}

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
	return NewLog(logger)
}

// NewLog ...
func NewLog(logger log.Logger) Log {
	return Log{
		Logger: logger,
	}
}
