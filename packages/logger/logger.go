package logger

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// New creates a zap logger. dev=true -> readable console format, else JSON production.
func New(dev bool) *zap.Logger {
	if dev {
		cfg := zap.NewDevelopmentConfig()
		cfg.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
		l, _ := cfg.Build()
		return l
	}
	l, _ := zap.NewProduction()
	return l
}

// MustNew panics if build fails (rarely happens).
func MustNew(dev bool) *zap.Logger {
	return New(dev)
}
