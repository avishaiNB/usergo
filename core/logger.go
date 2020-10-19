package core

import (
	"context"

	"github.com/go-kit/kit/log"
)

// Log will create a new instance of the Log with ready to use loggers
// Logger should:
// output: log to stdout, stderr, file
// data: each log should append information: correlation id, env name, process name, host name, duration, deadline, logger name, level
// levels: need to check how the levels play here and how we can control them in order to write the log
// read performance related concerns for using file appenders
// TBD: funnel logger

// LoggerLevel represent logger level
type LoggerLevel int8

const (
	// DebugLoggerLevel contains thelotter interpretation value of debug level
	DebugLoggerLevel LoggerLevel = 1
	// InfoLoggerLevel contains thelotter interpretation value of info level
	InfoLoggerLevel LoggerLevel = 2
	// WarnLoggerLevel contains thelotter interpretation value of warn level
	WarnLoggerLevel LoggerLevel = 3
	// ErrorLoggerLevel contains thelotter interpretation value of error level
	ErrorLoggerLevel LoggerLevel = 4
	// PanicLoggerLevel contains thelotter interpretation value of panic level
	PanicLoggerLevel LoggerLevel = 5
)

// Log is main class og logger
// Logger - represents log.Logger
// LoggerManager - represents LoggerManager
type Log struct {
	Logger        log.Logger
	LoggerManager LoggerManager
}

type logger struct {
	LoggerManager LoggerManager
}

// Logger represent convention of all possible loggers (file logger , stdout logger , humio logger etc.)
type Logger interface {
	Log(context.Context, LoggerLevel, string, ...interface{}) error
}

//NewGoKitLogger create new logger which represents go-kit Logger
func NewGoKitLogger(loggerManager LoggerManager) log.Logger {
	return &logger{
		LoggerManager: loggerManager,
	}
}

// NewLogWithDefaults create Log object with
func NewLogWithDefaults() Log {
	return Log{
		Logger:        logger{},
		LoggerManager: loggerManager{},
	}
}

// NewLog create new Log object
// logger - represent struct which "implement" log.Logger contract
// loggerManager - represent struct which "implement" LoggerManager contract
func NewLog(logger log.Logger, loggerManager LoggerManager) Log {
	return Log{
		Logger:        logger,
		LoggerManager: loggerManager,
	}
}

// Log func in part of go-kit logger contract
// kvs argument reuire next fields:
// "level" as LoggerLevel - level of log (info , warn etc.)
// "context" as context.Context
// "message" as string
func (logger logger) Log(kvs ...interface{}) error {
	args := make(map[string]interface{})
	for i := 0; i < len(kvs); i += 2 {
		key := kvs[i].(string)
		args[key] = kvs[i+1]
	}

	logLevel, levelOk := args["level"].(LoggerLevel)
	if levelOk {
		delete(args, "level")
	} else {
		logLevel = InfoLoggerLevel
	}

	ctx, ctxOk := args["context"].(context.Context)
	if ctxOk {
		delete(args, "context")
	} else {
		ctx = context.Background()
	}

	message, messageOk := args["message"].(string)
	if messageOk {
		delete(args, "message")
	} else {
		message = ""
	}

	switch logLevel {
	case DebugLoggerLevel:
		return logger.LoggerManager.Debug(ctx, message, args)
	case InfoLoggerLevel:
		return logger.LoggerManager.Info(ctx, message, args)
	case WarnLoggerLevel:
		return logger.LoggerManager.Warn(ctx, message, args)
	case ErrorLoggerLevel:
		return logger.LoggerManager.Error(ctx, message, args)
	case PanicLoggerLevel:
		return logger.LoggerManager.Panic(ctx, message, args)
	default:
		return logger.LoggerManager.Debug(ctx, message, args)
	}
}
