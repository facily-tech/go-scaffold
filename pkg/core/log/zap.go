package log

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func NewLogger(disableStackTrace bool) (*zap.Logger, error) {
	loggerConfig := zap.NewProductionConfig()
	loggerConfig.EncoderConfig.TimeKey = "timestamp"
	loggerConfig.EncoderConfig.EncodeTime = zapcore.RFC3339NanoTimeEncoder
	loggerConfig.DisableStacktrace = disableStackTrace

	// your log setup goes here

	logger, err := loggerConfig.Build()
	if err != nil {
		return nil, err
	}

	return logger, nil
}
