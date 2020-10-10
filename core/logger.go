package core

import (
	"os"

	"github.com/go-kit/kit/log"
)

// Log ...
type Log struct {
	Logger log.Logger
}

// NewLogWithDefaults ..
func NewLogWithDefaults() Log {
	var logger log.Logger
	{
		logger = log.NewLogfmtLogger(os.Stderr)
		logger = log.With(logger, "ts", log.DefaultTimestampUTC)
		logger = log.With(logger, "caller", log.DefaultCaller)
		//log.With(logger, "component", "HTTP")
	}

	return NewLog(logger)
}

// NewLog ...
func NewLog(logger log.Logger) Log {
	return Log{
		Logger: logger,
	}
}
