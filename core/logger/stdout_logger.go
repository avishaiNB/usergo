package logger

import (
	"context"
	"fmt"

	gokitZap "github.com/go-kit/kit/log/zap"
	tlecontext "github.com/thelotter-enterprise/usergo/core/context"
	"github.com/thelotter-enterprise/usergo/core/utils"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

const (
	//Message - log message
	Message = "message"
	//CorrelationID - correlation id of current stack trace
	CorrelationID = "correlationID"
	//Duration - duration of current stack trace
	Duration = "duration"
	//Timeout - timeout of current stack trace
	Timeout = "timeout"
)

type stdoutLogger struct {
	zapLogger *zap.Logger
}

// NewStdOutLogger create new stdoutLogger of type Logger
// "params map[string]interface{}" should contains next field:
// fileAtomicLevel - minimal log level (Debug , Info , Warn , Error or Panic)
// env - name of env
// processName - name of the current process
func NewStdOutLogger(loggerConfig Config) Logger {

	config := zap.NewProductionConfig()
	config.OutputPaths = []string{"stdout"}
	config.Level = getAtomicLevel(loggerConfig.LevelName)
	config.EncoderConfig.LevelKey = "level"
	config.EncoderConfig.TimeKey = "timestamp"
	config.EncoderConfig.CallerKey = "caller"
	processName := utils.ProcessName()

	config.InitialFields = make(map[string]interface{})
	config.InitialFields["env"] = loggerConfig.Env
	config.InitialFields["loggerName"] = loggerConfig.LoggerName
	config.InitialFields["processName"] = processName
	logger, err := config.Build()
	if err != nil {
		fmt.Println("Cannot init stdout logger", err)
	}
	return &stdoutLogger{
		zapLogger: logger,
	}
}

func (stdoutLogger stdoutLogger) Log(ctx context.Context, loggerLevel Level, message string, params ...interface{}) error {
	logLevel := stdoutLogger.castLoggerLevel(loggerLevel)
	correlationID := tlecontext.GetCorrelationFromContext(ctx)
	duration, timeout := tlecontext.GetTimeoutFromContext(ctx)
	gokitLogger := gokitZap.NewZapSugarLogger(stdoutLogger.zapLogger, logLevel)
	params = addParamsToLog(CorrelationID, correlationID, params)
	params = addParamsToLog(Message, message, params)
	params = addParamsToLog(Duration, duration, params)
	params = addParamsToLog(Timeout, timeout, params)
	return gokitLogger.Log(params...)
}

func addParamsToLog(key string, value interface{}, params []interface{}) []interface{} {
	params = append(params, key)
	params = append(params, value)
	return params
}

func (stdoutLogger stdoutLogger) castLoggerLevel(loggerLevel Level) zapcore.Level {
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
