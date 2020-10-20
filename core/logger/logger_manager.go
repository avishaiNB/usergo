package logger

import (
	"context"
	"errors"
	"fmt"
)

// Manager represents contract of logger with all log levels (Panic , Error , Warn , Info and Debug)
type Manager interface {
	// Panic represents convention of palic log function
	Panic(context.Context, string, ...interface{}) error
	// Error represents convention of error log function
	Error(context.Context, string, ...interface{}) error
	// Warn represents convention of warn log function
	Warn(context.Context, string, ...interface{}) error
	// Info represents convention of info log function
	Info(context.Context, string, ...interface{}) error
	// Debug represents convention of debug log function
	Debug(context.Context, string, ...interface{}) error
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
func (loggerManager loggerManager) Debug(ctx context.Context, message string, params ...interface{}) error {
	var isError bool
	for _, log := range loggerManager.Loggers {
		logErr := log.Log(ctx, DebugLoggerLevel, message, params...)
		if logErr != nil {
			isError = true
			fmt.Println("Cannot print DebugLoggerLevel log", logErr, log)
		}
	}
	if isError {
		return errors.New("One of the loggers throw exception")
	}
	return nil
}

// Info print logs to all Loggers on info level
func (loggerManager loggerManager) Info(ctx context.Context, message string, params ...interface{}) error {
	var isError bool
	for _, log := range loggerManager.Loggers {
		logErr := log.Log(ctx, InfoLoggerLevel, message, params...)
		if logErr != nil {
			isError = true
			fmt.Println("Cannot print InfoLoggerLevel log", logErr, log)
		}
	}
	if isError {
		return errors.New("One of the loggers throw exception")
	}
	return nil
}

// Warn print logs to all Loggers on warn level
func (loggerManager loggerManager) Warn(ctx context.Context, message string, params ...interface{}) error {
	var isError bool
	for _, log := range loggerManager.Loggers {
		logErr := log.Log(ctx, WarnLoggerLevel, message, params...)
		if logErr != nil {
			isError = true
			fmt.Println("Cannot print WarnLoggerLevel log", logErr, log)
		}
	}
	if isError {
		return errors.New("One of the loggers throw exception")
	}
	return nil
}

// Error print logs to all Loggers on error level
func (loggerManager loggerManager) Error(ctx context.Context, message string, params ...interface{}) error {
	var isError bool
	for _, log := range loggerManager.Loggers {
		logErr := log.Log(ctx, ErrorLoggerLevel, message, params...)
		if logErr != nil {
			isError = true
			fmt.Println("Cannot print ErrorLoggerLevel log", logErr, log)
		}
	}
	if isError {
		return errors.New("One of the loggers throw exception")
	}
	return nil
}

// Panic print logs to all Loggers on panic level
func (loggerManager loggerManager) Panic(ctx context.Context, message string, params ...interface{}) error {
	var isError bool
	for _, log := range loggerManager.Loggers {
		logErr := log.Log(ctx, PanicLoggerLevel, message, params...)
		if logErr != nil {
			isError = true
			fmt.Println("Cannot print PanicLoggerLevel log", logErr, log)
		}
	}
	if isError {
		return errors.New("One of the loggers throw exception")
	}
	return nil
}
