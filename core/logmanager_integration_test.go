package core_test

import (
	"context"
	"testing"

	"github.com/thelotter-enterprise/usergo/core"
)

func TestInregrationLogger(t *testing.T) {
	var isIntegrationTest bool
	isIntegrationTest = false
	if isIntegrationTest {
		loggerConfig := core.LoggerConfig{
			Env:         "Dev",
			LevelName:   core.Info,
			LoggerName:  "FileLogger",
			ProcessName: "UserGo",
		}
		ctx := context.Background()
		logUser := LogUser{
			Name: "David",
			Age:  3,
		}
		var loggers []core.Logger
		fileLogger := core.NewFileLogger(loggerConfig)
		loggerConfig.LoggerName = "StdOutLogger"
		stdlogger := core.NewStdOutLogger(loggerConfig)
		loggers = append(loggers, stdlogger)
		loggers = append(loggers, fileLogger)
		loggerManager := core.NewLoggerManager(loggers)
		loggerManager.Info(ctx, "Text", "Log user", logUser)
		loggerManager.Info(ctx, "Text", "Log user1", logUser)
	}
}
