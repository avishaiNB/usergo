package logger_test

import (
	"context"
	"testing"

	"github.com/thelotter-enterprise/usergo/core/logger"
)

func TestInregrationLogger(t *testing.T) {
	var isIntegrationTest bool
	isIntegrationTest = false
	if isIntegrationTest {
		loggerConfig := logger.Config{
			Env:         "Dev",
			LevelName:   logger.Info,
			LoggerName:  "FileLogger",
			ProcessName: "UserGo",
		}
		ctx := context.Background()
		logUser := LogUser{
			Name: "David",
			Age:  3,
		}
		var loggers []logger.Logger
		fileLogger := logger.NewFileLogger(loggerConfig)
		loggerConfig.LoggerName = "StdOutLogger"
		stdlogger := logger.NewStdOutLogger(loggerConfig)
		loggers = append(loggers, stdlogger)
		loggers = append(loggers, fileLogger)
		loggerManager := logger.NewLoggerManager(loggers)
		loggerManager.Info(ctx, "Text", "Log user", logUser)
		loggerManager.Info(ctx, "Text", "Log user1", logUser)
	}
}
