package logging

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func MustSetLevelLog(initLevel string) *zap.Logger {
	var logger *zap.Logger
	var logLevel zapcore.Level

	err := logLevel.Set(initLevel)
	if err != nil {
		panic(err)
	}
	zapCfg := zap.NewProductionConfig()
	zapCfg.Level.SetLevel(logLevel)
	zapCfg.EncoderConfig.TimeKey = "timestamp"
	zapCfg.EncoderConfig.MessageKey = "message"
	zapCfg.EncoderConfig.LevelKey = "severity"
	zapCfg.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	logger, err = zapCfg.Build()
	if err != nil {
		panic(err)
	}
	return logger
}
