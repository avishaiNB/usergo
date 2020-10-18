package core

import "context"

type LoggerManager interface {
	Panic(context.Context, string, ...interface{})
	Error(context.Context, string, ...interface{})
	Warn(context.Context, string, ...interface{})
	Info(context.Context, string, ...interface{})
	Debug(context.Context, string, ...interface{})
}

type loggerManager struct {
	Loggers []Logger
}

func NewLoggerClient(loggers []Logger) LoggerManager {
	return &loggerManager{
		Loggers: loggers,
	}
}

func (loggerManager loggerManager) Debug(ctx context.Context, message string, params ...interface{}) {

}

func (loggerManager loggerManager) Info(ctx context.Context, message string, params ...interface{}) {

}

func (loggerManager loggerManager) Warn(ctx context.Context, message string, params ...interface{}) {

}

func (loggerManager loggerManager) Error(ctx context.Context, message string, params ...interface{}) {

}

func (loggerManager loggerManager) Panic(ctx context.Context, message string, params ...interface{}) {

}
