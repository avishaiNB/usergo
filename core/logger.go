package core

import (
	"context"

	"github.com/go-kit/kit/log"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// Log will create a new instance of the Log with ready to use loggers
// Logger should:
// output: log to stdout, stderr, file
// data: each log should append information: correlation id, env name, process name, host name, duration, deadline, logger name, level
// levels: need to check how the levels play here and how we can control them in order to write the log
// read performance related concerns for using file appenders
// TBD: funnel logger

// LoggerLevel ...
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

type Log struct {
	Logger log.Logger
}

type logger struct {
	LoggerManager LoggerManager
}

type Logger interface {
	Log(ctx context.Context, message string, params ...interface{})
}

func NewLog(logger log.Logger) Log {
	return Log{
		Logger: logger,
	}
}

// func (logger logger) Log(ctx Ctx, level zapcore.Level, message string, params ...interface{}) {
// 	correlationID := ctx.GetCorrelationFromContext(ctx.Context)
// 	duration, timeout := ctx.GetTimeoutFromContext(ctx.Context)
// 	gokitLogger := gokitZap.NewZapSugarLogger(logger.zapLogger, level)
// 	gokitLogger.Log(message,
// 		"correlationID", correlationID,
// 		"duration", duration,
// 		"timeout", timeout,
// 		"params", params)
// }

func getAtomicLevel(atomicLevel interface{}) zap.AtomicLevel {
	atom := zap.NewAtomicLevel()
	if atomicLevel == nil {
		atom.SetLevel(zapcore.InfoLevel)
	} else {
		switch al := atomicLevel.(string); al {
		case "Debug":
			atom.SetLevel(zapcore.DebugLevel)
		case "Info":
			atom.SetLevel(zapcore.InfoLevel)
		case "Warn":
			atom.SetLevel(zapcore.WarnLevel)
		case "Error":
			atom.SetLevel(zapcore.ErrorLevel)
		default:
			atom.SetLevel(zapcore.InfoLevel)
		}
	}
	return atom
}

func (logger logger) Log(kvs ...interface{}) {
	var args map[string]interface{}
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
		logger.LoggerManager.Debug(ctx, message, args)
	case InfoLoggerLevel:
		logger.LoggerManager.Info(ctx, message, args)
	case WarnLoggerLevel:
		logger.LoggerManager.Warn(ctx, message, args)
	case ErrorLoggerLevel:
		logger.LoggerManager.Error(ctx, message, args)
	case PanicLoggerLevel:
		logger.LoggerManager.Panic(ctx, message, args)
	default:
		logger.LoggerManager.Debug(ctx, message, args)
	}
}
