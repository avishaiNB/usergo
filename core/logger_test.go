package core_test

import (
	"context"
	"errors"
	"testing"

	"github.com/thelotter-enterprise/usergo/core"
)

func TestLoggerManagerReturnNoError(t *testing.T) {
	ctx := context.Background()
	logUser := LogUser{
		Name: "David",
		Age:  3,
	}
	var loggers []core.Logger
	stdlogger := &fileLoggerMock{}
	fileLogger := &fileLoggerMock{}
	loggers = append(loggers, stdlogger)
	loggers = append(loggers, fileLogger)
	loggerManager := core.NewLoggerManager(loggers)
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
	var loggers []core.Logger
	stdlogger := &fileLoggerMock{}
	fileLogger := &fileLoggerMock{}
	loggers = append(loggers, stdlogger)
	loggers = append(loggers, fileLogger)
	loggerManager := core.NewLoggerManager(loggers)
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
	params = append(params, core.InfoLoggerLevel)
	params = append(params, "context")
	params = append(params, ctx)
	params = append(params, "LogUser")
	params = append(params, logUser)
	var loggers []core.Logger
	stdlogger := &fileLoggerMock{}
	fileLogger := &fileLoggerMock{}
	loggers = append(loggers, stdlogger)
	loggers = append(loggers, fileLogger)
	loggerManager := core.NewLoggerManager(loggers)
	goKitLogger := core.NewLogger(loggerManager)
	log := core.SetLog(goKitLogger, loggerManager)
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
	params = append(params, core.ErrorLoggerLevel)
	params = append(params, "context")
	params = append(params, ctx)
	params = append(params, "LogUser")
	params = append(params, logUser)
	var loggers []core.Logger
	stdlogger := &fileLoggerMock{}
	fileLogger := &fileLoggerMock{}
	loggers = append(loggers, stdlogger)
	loggers = append(loggers, fileLogger)
	loggerManager := core.NewLoggerManager(loggers)
	goKitLogger := core.NewLogger(loggerManager)
	log := core.SetLog(goKitLogger, loggerManager)
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
	kvs = append(kvs, core.WarnLoggerLevel)
	kvs = append(kvs, "context")
	kvs = append(kvs, ctx)
	kvs = append(kvs, "message")
	kvs = append(kvs, message)
	loggerBuild := core.BuildLogData(kvs...)
	if loggerBuild.Level != core.WarnLoggerLevel {
		t.Errorf("loggerBuild return wrong log level %v ; want %v", loggerBuild.Level, core.WarnLoggerLevel)
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
	loggerBuild := core.BuildLogData(kvs...)
	if loggerBuild.Level != core.InfoLoggerLevel {
		t.Errorf("loggerBuild return wrong log level %v ; want %v", loggerBuild.Level, core.InfoLoggerLevel)
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
	loggerBuild := core.BuildLogData(kvs...)
	if loggerBuild.Level != core.InfoLoggerLevel {
		t.Errorf("loggerBuild return wrong log level %v ; want %v", loggerBuild.Level, core.InfoLoggerLevel)
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

func (fileLogger *fileLoggerMock) Log(ctx context.Context, loggerLevel core.LoggerLevel, message string, params ...interface{}) error {
	if loggerLevel == core.ErrorLoggerLevel {
		return errors.New("Custom error")
	}
	return nil
}
