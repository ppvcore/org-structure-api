package logger

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var Logger *zap.Logger

func InitZap() error {
	config := zap.NewProductionConfig()

	config.OutputPaths = []string{"logs/app.log"}

	config.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder

	var err error
	Logger, err = config.Build()
	if err != nil {
		return err
	}

	zap.ReplaceGlobals(Logger)
	return nil
}

func Close() {
	if Logger != nil {
		_ = Logger.Sync()
	}
}
