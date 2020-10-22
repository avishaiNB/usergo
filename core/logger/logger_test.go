package logger_test

import (
	"context"
	"errors"
	"testing"

	"github.com/thelotter-enterprise/usergo/core/logger"
)

func TestLoggerManagerReturnNoError(t *testing.T) {
	ctx := context.Background()
	logUser := LogUser{
		Name: "David",
		Age:  3,
	}
	var loggers []logger.Logger
	stdlogger := &fileLoggerMock{}
	fileLogger := &fileLoggerMock{}
	loggers = append(loggers, stdlogger)
	loggers = append(loggers, fileLogger)
	loggerManager := logger.NewLoggerManager(loggers...)
	loggerErr := loggerManager.Info(ctx, "Text", "Log user", logUser)
	if loggerErr != nil {
		t.Error("loggerManager.Info return unexpected error", loggerErr)
	}
}

func TestLoggerManagerReturnError(t *testing.T) {
	ctx := context.Background()
	logUser := LogUser{
		Name: "David",
		Age:  3,
	}
	var loggers []logger.Logger
	stdlogger := &fileLoggerMock{}
	fileLogger := &fileLoggerMock{}
	loggers = append(loggers, stdlogger)
	loggers = append(loggers, fileLogger)
	loggerManager := logger.NewLoggerManager(loggers...)
	loggerErr := loggerManager.Error(ctx, "Text", "Log user", logUser)
	if loggerErr == nil {
		t.Error("loggerManager.Error should return error")
	}
}

func TestLoggerReturnNoError(t *testing.T) {
	ctx := context.Background()
	var params []interface{}
	logUser := LogUser{
		Name: "David",
		Age:  3,
	}
	params = append(params, "level")
	params = append(params, logger.InfoLoggerLevel)
	params = append(params, "context")
	params = append(params, ctx)
	params = append(params, "LogUser")
	params = append(params, logUser)
	var loggers []logger.Logger
	stdlogger := &fileLoggerMock{}
	fileLogger := &fileLoggerMock{}
	loggers = append(loggers, stdlogger)
	loggers = append(loggers, fileLogger)
	loggerManager := logger.NewLoggerManager(loggers...)
	goKitLogger := logger.NewLogger(loggerManager)
	log := logger.SetLog(goKitLogger, loggerManager)
	logErr := log.Logger.Log(params...)
	if logErr != nil {
		t.Error("log.Logger.Log return unexpected error", logErr)
	}
}

func TestLoggerReturnError(t *testing.T) {
	ctx := context.Background()
	var params []interface{}
	logUser := LogUser{
		Name: "David",
		Age:  3,
	}
	params = append(params, "level")
	params = append(params, logger.ErrorLoggerLevel)
	params = append(params, "context")
	params = append(params, ctx)
	params = append(params, "LogUser")
	params = append(params, logUser)
	var loggers []logger.Logger
	stdlogger := &fileLoggerMock{}
	fileLogger := &fileLoggerMock{}
	loggers = append(loggers, stdlogger)
	loggers = append(loggers, fileLogger)
	loggerManager := logger.NewLoggerManager(loggers...)
	goKitLogger := logger.NewLogger(loggerManager)
	log := logger.SetLog(goKitLogger, loggerManager)
	logErr := log.Logger.Log(params...)
	if logErr == nil {
		t.Error("Expected result from log.Logger.Log cannot ne nil")
	}
}

func TestBuildLogDataWithAllData(t *testing.T) {
	var kvs []interface{}
	ctx := context.Background()
	message := "text"
	kvs = append(kvs, "level")
	kvs = append(kvs, logger.WarnLoggerLevel)
	kvs = append(kvs, "context")
	kvs = append(kvs, ctx)
	kvs = append(kvs, "message")
	kvs = append(kvs, message)
	loggerBuild := logger.BuildLogData(kvs...)
	if loggerBuild.Level != logger.WarnLoggerLevel {
		t.Errorf("loggerBuild return wrong log level %v ; want %v", loggerBuild.Level, logger.WarnLoggerLevel)
	}

	if loggerBuild.Message != message {
		t.Errorf("loggerBuild return wrong message %s ; want %s", loggerBuild.Message, message)
	}

	if loggerBuild.Context != ctx {
		t.Errorf("loggerBuild return wrong context %s ; want %s", loggerBuild.Context, ctx)
	}
}

func TestBuildLogDataWithoutLevel(t *testing.T) {
	var kvs []interface{}
	ctx := context.Background()
	message := "text"
	kvs = append(kvs, "context")
	kvs = append(kvs, ctx)
	kvs = append(kvs, "message")
	kvs = append(kvs, message)
	loggerBuild := logger.BuildLogData(kvs...)
	if loggerBuild.Level != logger.InfoLoggerLevel {
		t.Errorf("loggerBuild return wrong log level %v ; want %v", loggerBuild.Level, logger.InfoLoggerLevel)
	}

	if loggerBuild.Message != message {
		t.Errorf("loggerBuild return wrong message %s ; want %s", loggerBuild.Message, message)
	}

	if loggerBuild.Context != ctx {
		t.Errorf("loggerBuild return wrong context %s ; want %s", loggerBuild.Context, ctx)
	}
}

func TestBuildLogDataWithExtraData(t *testing.T) {
	var kvs []interface{}
	ctx := context.Background()
	message := "text"
	customValue := "customValue"
	kvs = append(kvs, "context")
	kvs = append(kvs, ctx)
	kvs = append(kvs, "message")
	kvs = append(kvs, message)
	kvs = append(kvs, customValue)
	kvs = append(kvs, customValue)
	loggerBuild := logger.BuildLogData(kvs...)
	if loggerBuild.Level != logger.InfoLoggerLevel {
		t.Errorf("loggerBuild return wrong log level %v ; want %v", loggerBuild.Level, logger.InfoLoggerLevel)
	}

	if loggerBuild.Message != message {
		t.Errorf("loggerBuild return wrong message %s ; want %s", loggerBuild.Message, message)
	}

	if loggerBuild.Context != ctx {
		t.Errorf("loggerBuild return wrong context %s ; want %s", loggerBuild.Context, ctx)
	}

	loggerBuildCustomResult := loggerBuild.Data["customValue"].(string)
	if loggerBuildCustomResult != customValue {
		t.Errorf("loggerBuild return wrong Data %v ; want %v", loggerBuildCustomResult, customValue)
	}
}

type LogUser struct {
	Name string
	Age  int
}

type fileLoggerMock struct {
}

func (fileLogger *fileLoggerMock) Log(ctx context.Context, loggerLevel logger.Level, message string, params ...interface{}) error {
	if loggerLevel == logger.ErrorLoggerLevel {
		return errors.New("Custom error")
	}
	return nil
}
