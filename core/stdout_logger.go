package core

import (
	"context"
	"fmt"

	gokitZap "github.com/go-kit/kit/log/zap"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type stdoutLogger struct {
	zapLogger   *zap.Logger
	Ctx         Ctx
	LoggerName  string
	EnvName     string
	ProcessName string
}

// NewStdOutLogger create new stdoutLogger of type Logger
// "params map[string]interface{}" should contains next field:
// fileAtomicLevel - minimal log level (Debug , Info , Warn , Error or Panic)
// env - name of env
// processName - name of the current process
func NewStdOutLogger(params map[string]interface{}) Logger {
	ctx := NewCtx()

	config := zap.NewProductionConfig()
	config.OutputPaths = []string{"stdout"}
	config.Level = getAtomicLevel(params["stdOutAtomicLevel"])
	config.EncoderConfig.LevelKey = "level"
	config.EncoderConfig.TimeKey = "timestamp"
	config.EncoderConfig.CallerKey = "caller"

	config.InitialFields = make(map[string]interface{})
	config.InitialFields["env"] = params["env"]
	config.InitialFields["loggerName"] = "FileLogger"
	config.InitialFields["processName"] = params["serviceName"]
	logger, err := config.Build()
	if err != nil {
		fmt.Println("Cannot init stdout logger", err)
	}
	return &fileLogger{
		zapLogger: logger,
		Ctx:       ctx,
	}
}

func (stdoutLogger stdoutLogger) Log(ctx context.Context, loggerLevel LoggerLevel, message string, params ...interface{}) error {
	logLevel := stdoutLogger.castLoggerLevel(loggerLevel)
	correlationID := stdoutLogger.Ctx.GetCorrelationFromContext(ctx)
	duration, timeout := stdoutLogger.Ctx.GetTimeoutFromContext(ctx)
	gokitLogger := gokitZap.NewZapSugarLogger(stdoutLogger.zapLogger, logLevel)
	params = addParamsToLog("correlationID", correlationID, params)
	params = addParamsToLog("Message", message, params)
	params = addParamsToLog("duration", duration, params)
	params = addParamsToLog("timeout", timeout, params)
	return gokitLogger.Log(params...)
}

func addParamsToLog(key string, value interface{}, params []interface{}) []interface{} {
	params = append(params, key)
	params = append(params, value)
	return params
}

func (stdoutLogger stdoutLogger) castLoggerLevel(loggerLevel LoggerLevel) zapcore.Level {
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

func getAtomicLevel(atomicLevel interface{}) zap.AtomicLevel {
	atom := zap.NewAtomicLevel()
	if atomicLevel == nil {
		atom.SetLevel(zapcore.InfoLevel)
	} else {
		switch al := atomicLevel.(string); al {
		case "Debug":
			atom.SetLevel(zapcore.DebugLevel)
		case "Info":
			atom.SetLevel(zapcore.InfoLevel)
		case "Warn":
			atom.SetLevel(zapcore.WarnLevel)
		case "Error":
			atom.SetLevel(zapcore.ErrorLevel)
		default:
			atom.SetLevel(zapcore.InfoLevel)
		}
	}
	return atom
}
