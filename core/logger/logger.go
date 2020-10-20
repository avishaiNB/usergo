package logger

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

// Level represent logger level
type Level int8

// AtomicLevelName represent name of specific log level
type AtomicLevelName string

const (
	// DebugLoggerLevel contains thelotter interpretation value of debug level
	DebugLoggerLevel Level = 1
	// InfoLoggerLevel contains thelotter interpretation value of info level
	InfoLoggerLevel Level = 2
	// WarnLoggerLevel contains thelotter interpretation value of warn level
	WarnLoggerLevel Level = 3
	// ErrorLoggerLevel contains thelotter interpretation value of error level
	ErrorLoggerLevel Level = 4
	// PanicLoggerLevel contains thelotter interpretation value of panic level
	PanicLoggerLevel Level = 5
	// Debug contains name of debug level
	Debug AtomicLevelName = "DEBUG"
	// Info contains name of info level
	Info AtomicLevelName = "INFO"
	// Warn contains name of warn level
	Warn AtomicLevelName = "WARN"
	// Error contains name of error level
	Error AtomicLevelName = "ERROR"
	// Panic contains name of panic level
	Panic AtomicLevelName = "PANIC"
)

// Log is main class og logger
// Logger - represents log.Logger
// LoggerManager - represents LoggerManager
type Log struct {
	Logger        log.Logger
	LoggerManager Manager
}

// Config represents base configurations of logger
// LevelName - minimal log level
// Env- name of current environment
// LoggerName - name of the logger
// ProcessName - name of the current process
type Config struct {
	LevelName   AtomicLevelName
	Env         string
	LoggerName  string
	ProcessName string
}

// LogData represents log data created by BuildLogData function based on kv interface array
type LogData struct {
	Level   Level
	Message string
	Context context.Context
	Data    map[string]interface{}
}

type logger struct {
	LoggerManager Manager
}

// Logger represent convention of all possible loggers (file logger , stdout logger , humio logger etc.)
type Logger interface {
	Log(context.Context, Level, string, ...interface{}) error
}

//NewLogger create new logger which represents go-kit Logger
func NewLogger(loggerManager Manager) log.Logger {
	return &logger{
		LoggerManager: loggerManager,
	}
}

// NewLog create Log object with dafault params and stdout logger
func NewLog() Log {
	loggerConfig := Config{
		LoggerName: "StdoutLogger",
	}

	stdOutLogger := NewStdOutLogger(loggerConfig)
	loggers := []Logger{stdOutLogger}
	logManager := NewLoggerManager(loggers)
	log := NewLogger(logManager)
	return Log{
		Logger:        log,
		LoggerManager: logManager,
	}
}

// SetLog create new Log object
// logger - represent struct which "implement" log.Logger contract
// loggerManager - represent struct which "implement" LoggerManager contract
func SetLog(logger log.Logger, loggerManager Manager) Log {
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
	logData := BuildLogData(kvs...)

	switch logData.Level {
	case DebugLoggerLevel:
		return logger.LoggerManager.Debug(logData.Context, logData.Message, logData.Data)
	case InfoLoggerLevel:
		return logger.LoggerManager.Info(logData.Context, logData.Message, logData.Data)
	case WarnLoggerLevel:
		return logger.LoggerManager.Warn(logData.Context, logData.Message, logData.Data)
	case ErrorLoggerLevel:
		return logger.LoggerManager.Error(logData.Context, logData.Message, logData.Data)
	case PanicLoggerLevel:
		return logger.LoggerManager.Panic(logData.Context, logData.Message, logData.Data)
	default:
		return logger.LoggerManager.Debug(logData.Context, logData.Message, logData.Data)
	}
}

// BuildLogData build LogData based on kvs array , each odd value is key(string) and even value(any type)
func BuildLogData(kvs ...interface{}) LogData {
	args := make(map[string]interface{})
	for i := 0; i < len(kvs); i += 2 {
		key := kvs[i].(string)
		args[key] = kvs[i+1]
	}

	logLevel, levelOk := args["level"].(Level)
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
	return LogData{
		Message: message,
		Context: ctx,
		Level:   logLevel,
		Data:    args,
	}
}

func getAtomicLevel(atomicLevelName AtomicLevelName) zap.AtomicLevel {
	atom := zap.NewAtomicLevel()
	if atomicLevelName == "" {
		atom.SetLevel(zapcore.InfoLevel)
	} else {
		switch atomicLevelName {
		case Debug:
			atom.SetLevel(zapcore.DebugLevel)
		case Info:
			atom.SetLevel(zapcore.InfoLevel)
		case Warn:
			atom.SetLevel(zapcore.WarnLevel)
		case Error:
			atom.SetLevel(zapcore.ErrorLevel)
		case Panic:
			atom.SetLevel(zapcore.PanicLevel)
		default:
			atom.SetLevel(zapcore.InfoLevel)
		}
	}
	return atom
}
