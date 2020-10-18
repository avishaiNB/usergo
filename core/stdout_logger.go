package core

import (
	"context"
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type stdoutLogger struct {
	zapLogger *zap.Logger
	Ctx       Ctx
}

func NewStdOutLogger(params map[string]interface{}) Logger {
	atom := getAtomicLevel(params["stdOutAtomicLevel"])
	ctx := NewCtx()
	encoderCfg := zap.NewProductionEncoderConfig()
	encoderCfg.TimeKey = ""

	log := zap.New(zapcore.NewCore(
		zapcore.NewJSONEncoder(encoderCfg),
		zapcore.Lock(os.Stdout),
		atom,
	))
	return &stdoutLogger{
		zapLogger: log,
		Ctx:       ctx,
	}
}

func (stdoutLogger stdoutLogger) Log(ctx context.Context, message string, params ...interface{}) {

}
