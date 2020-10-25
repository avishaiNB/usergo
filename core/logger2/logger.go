package logger2

import (
	"os"
	"time"

	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	kitzap "github.com/go-kit/kit/log/zap"
	"github.com/thelotter-enterprise/usergo/core/utils"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// NewLogger ...
func NewLogger(env string, loggerName string, minLevel AtomicLevelName) log.Logger {
	dt := utils.DateTime{}
	config := zap.NewProductionConfig()
	config.EncoderConfig.LevelKey = "l"
	config.EncoderConfig.TimeKey = "t"
	config.EncoderConfig.CallerKey = "c"
	config.OutputPaths = append(config.OutputPaths, getOrCreatelogFilePath(dt.Now()))

	zapLogger, _ := config.Build()

	// TODO: I am not sure what is the usage for the zapcore.InfoLevel
	kitLogger := kitzap.NewZapSugarLogger(zapLogger, zapcore.InfoLevel)
	kitLogger = log.With(kitLogger,
		// "level", loggerConfig.LevelName,
		"timestamp", log.DefaultTimestampUTC,
		"caller", log.DefaultCaller,
		"process", utils.ProcessName(),
		"loggerName", loggerName,
		"env", env,
	)

	kitLogger = level.NewFilter(kitLogger, toLevelOption(minLevel))

	return kitLogger
}

// TODO: fix the date so that the file will be created once every 3 hours
func getOrCreatelogFilePath(date time.Time) string {
	os.Mkdir("logs", os.ModePerm)
	fileName := "2006-01-02T02T15:04"
	return "logs/" + date.Format(fileName)
}

func toLevelOption(l AtomicLevelName) level.Option {
	switch l {
	case DebugLogLevel:
		return level.AllowDebug()
	case InfoLogLevel:
		return level.AllowInfo()
	case WarnLogLevel:
		return level.AllowWarn()
	case ErrorLogLevel:
		return level.AllowError()
	case PanicLogLevel:
		return level.AllowError()
	default:
		return level.AllowAll()
	}
}

// AtomicLevelName represent name of specific log level
type AtomicLevelName string

const (
	// DebugLogLevel contains name of debug level
	DebugLogLevel AtomicLevelName = "DEBUG"
	// InfoLogLevel contains name of info level
	InfoLogLevel AtomicLevelName = "INFO"
	// WarnLogLevel contains name of warn level
	WarnLogLevel AtomicLevelName = "WARN"
	// ErrorLogLevel contains name of error level
	ErrorLogLevel AtomicLevelName = "ERROR"
	// PanicLogLevel contains name of panic level
	PanicLogLevel AtomicLevelName = "PANIC"
)
