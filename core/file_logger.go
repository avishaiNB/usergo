package core

import (
	"context"
	"fmt"
	"os"
	"time"

	gokitZap "github.com/go-kit/kit/log/zap"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type fileLogger struct {
	Ctx Ctx

	zapLogger   *zap.Logger
	config      zap.Config
	dateCreated time.Time
}

// NewFileLogger create new fileLogger of type Logger
// "params map[string]interface{}" should contains next field:
// fileAtomicLevel - minimal log level (Debug , Info , Warn , Error or Panic)
// env - name of env
// processName - name of the current process
func NewFileLogger(params map[string]interface{}) Logger {
	os.Mkdir("logs", os.ModePerm)
	ctx := NewCtx()
	dateNow := time.Now()

	config := zap.NewProductionConfig()
	config.OutputPaths = []string{
		buildLogFilePath(dateNow),
	}
	config.Level = getAtomicLevel(params["fileAtomicLevel"])
	config.EncoderConfig.LevelKey = "level"
	config.EncoderConfig.TimeKey = "timestamp"
	config.EncoderConfig.CallerKey = "caller"

	config.InitialFields = make(map[string]interface{})
	config.InitialFields["env"] = params["env"]
	config.InitialFields["loggerName"] = "FileLogger"
	config.InitialFields["processName"] = params["serviceName"]

	logger, err := config.Build()
	if err != nil {
		fmt.Println("Cannot init file logger", err)
		return nil
	}
	return &fileLogger{
		zapLogger:   logger,
		Ctx:         ctx,
		config:      config,
		dateCreated: dateNow,
	}
}

func (fileLogger *fileLogger) Log(ctx context.Context, loggerLevel LoggerLevel, message string, params ...interface{}) error {
	if fileLogger.isLogFileExpired() {
		fileLogger.reload()
	}

	logLevel := fileLogger.castLoggerLevel(loggerLevel)
	correlationID := fileLogger.Ctx.GetCorrelationFromContext(ctx)
	duration, timeout := fileLogger.Ctx.GetTimeoutFromContext(ctx)

	gokitLogger := gokitZap.NewZapSugarLogger(fileLogger.zapLogger, logLevel)
	params = addParamsToLog("correlationID", correlationID, params)
	params = addParamsToLog("Message", message, params)
	params = addParamsToLog("duration", duration, params)
	params = addParamsToLog("timeout", timeout, params)
	return gokitLogger.Log(params...)
}

func (fileLogger fileLogger) castLoggerLevel(loggerLevel LoggerLevel) zapcore.Level {
	switch loggerLevel {
	case DebugLoggerLevel:
		return zapcore.DebugLevel
	case InfoLoggerLevel:
		return zapcore.InfoLevel
	case WarnLoggerLevel:
		return zapcore.WarnLevel
	case ErrorLoggerLevel:
		return zapcore.ErrorLevel
	case PanicLoggerLevel:
		return zapcore.PanicLevel
	default:
		return zapcore.InfoLevel
	}
}

func (fileLogger *fileLogger) reload() {
	dateNow := time.Now()
	fileLogger.config.OutputPaths = []string{
		buildLogFilePath(dateNow),
	}
	fileLogger.dateCreated = dateNow
	zapLogger, err := fileLogger.config.Build()
	if err != nil {
		fmt.Println("Cannot reload file logger", err)
		return
	}
	fileLogger.zapLogger = zapLogger
}

func (fileLogger fileLogger) isLogFileExpired() bool {
	return fileLogger.dateCreated.Day() != time.Now().Day()
}

func buildLogFilePath(date time.Time) string {
	layoutISO := "2006-01-02T02T15:04:05-0700"
	return "logs/log_" + date.Format(layoutISO)
}
