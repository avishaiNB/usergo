package logger

import (
	"context"
	"fmt"
	"strconv"

	tleerrors "github.com/thelotter-enterprise/usergo/core/errors"
)

// Manager represents contract of logger with all log levels (Panic , Error , Warn , Info and Debug)
type Manager interface {
	// Panic represents convention of palic log function
	Panic(context.Context, string, ...interface{}) tleerrors.ApplicationError
	// Error represents convention of error log function
	Error(context.Context, string, ...interface{}) tleerrors.ApplicationError
	// Warn represents convention of warn log function
	Warn(context.Context, string, ...interface{}) tleerrors.ApplicationError
	// Info represents convention of info log function
	Info(context.Context, string, ...interface{}) tleerrors.ApplicationError
	// Debug represents convention of debug log function
	Debug(context.Context, string, ...interface{}) tleerrors.ApplicationError
}

type loggerManager struct {
	Loggers []Logger
}

// NewLoggerManager create loggerManager and give you control on all loggers
func NewLoggerManager(loggers []Logger) Manager {
	return &loggerManager{
		Loggers: loggers,
	}
}

// Debug print logs to all Loggers on debug level
func (loggerManager loggerManager) Debug(ctx context.Context, message string, params ...interface{}) tleerrors.ApplicationError {
	err := tleerrors.ApplicationError{}
	for _, log := range loggerManager.Loggers {
		logErr := log.Log(ctx, DebugLoggerLevel, message, params...)
		if logErr != nil {
			fmt.Println("Cannot print DebugLoggerLevel log", logErr, log)
			loggerManager.addError(&err, logErr)
		}
	}
	return err
}

// Info print logs to all Loggers on info level
func (loggerManager loggerManager) Info(ctx context.Context, message string, params ...interface{}) tleerrors.ApplicationError {
	err := tleerrors.ApplicationError{}
	for _, log := range loggerManager.Loggers {
		logErr := log.Log(ctx, InfoLoggerLevel, message, params...)
		if logErr != nil {
			fmt.Println("Cannot print InfoLoggerLevel log", logErr, log)
			loggerManager.addError(&err, logErr)
		}
	}
	return err
}

// Warn print logs to all Loggers on warn level
func (loggerManager loggerManager) Warn(ctx context.Context, message string, params ...interface{}) tleerrors.ApplicationError {
	err := tleerrors.ApplicationError{}
	for _, log := range loggerManager.Loggers {
		logErr := log.Log(ctx, WarnLoggerLevel, message, params...)
		if logErr != nil {
			fmt.Println("Cannot print WarnLoggerLevel log", logErr, log)
			loggerManager.addError(&err, logErr)
		}
	}
	return err
}

// Error print logs to all Loggers on error level
func (loggerManager loggerManager) Error(ctx context.Context, message string, params ...interface{}) tleerrors.ApplicationError {
	err := tleerrors.ApplicationError{}
	for _, log := range loggerManager.Loggers {
		logErr := log.Log(ctx, ErrorLoggerLevel, message, params...)
		if logErr != nil {
			fmt.Println("Cannot print ErrorLoggerLevel log", logErr, log)
			loggerManager.addError(&err, logErr)
		}
	}
	return err
}

// Panic print logs to all Loggers on panic level
func (loggerManager loggerManager) Panic(ctx context.Context, message string, params ...interface{}) tleerrors.ApplicationError {
	err := tleerrors.ApplicationError{}
	for _, log := range loggerManager.Loggers {
		logErr := log.Log(ctx, PanicLoggerLevel, message, params...)
		if logErr != nil {
			fmt.Println("Cannot print PanicLoggerLevel log", logErr, log)
			loggerManager.addError(&err, logErr)
		}
	}
	return err
}

func (loggerManager loggerManager) addError(applicatioError *tleerrors.ApplicationError, err error) {
	if applicatioError.Err == nil {
		*applicatioError = tleerrors.NewApplicationError("One of the loggers throw exception", make(map[string]interface{}))
	}
	key := strconv.Itoa(len(applicatioError.Params))
	applicatioError.Params[key] = err
}
