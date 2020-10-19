package core_test

import (
	"context"
	"errors"
	"testing"

	"github.com/thelotter-enterprise/usergo/core"
)

func TestInregrationLogger(t *testing.T) {
	var isIntegrationTest bool
	if isIntegrationTest {
		params := make(map[string]interface{})
		ctx := context.Background()
		logUser := LogUser{
			Name: "David",
			Age:  3,
		}
		params["stdOutAtomicLevel"] = "Info"
		var loggers []core.Logger
		stdlogger := core.NewStdOutLogger(params)
		fileLogger := core.NewFileLogger(params)
		loggers = append(loggers, stdlogger)
		loggers = append(loggers, fileLogger)
		loggerManager := core.NewLoggerManager(loggers)
		loggerManager.Info(ctx, "Text", "Log user", logUser)
	}
}

func Test_LoggerManagerReturnNoError(t *testing.T) {
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

func Test_LoggerManagerReturnError(t *testing.T) {
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

func Test_LoggerReturnNoError(t *testing.T) {
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
	goKitLogger := core.NewGoKitLogger(loggerManager)
	log := core.NewLog(goKitLogger, loggerManager)
	logErr := log.Logger.Log(params...)
	if logErr != nil {
		t.Error("log.Logger.Log return unexpected error", logErr)
	}
}

func Test_LoggerReturnError(t *testing.T) {
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
	goKitLogger := core.NewGoKitLogger(loggerManager)
	log := core.NewLog(goKitLogger, loggerManager)
	logErr := log.Logger.Log(params...)
	if logErr == nil {
		t.Error("Expected result from log.Logger.Log cannot ne nil")
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
