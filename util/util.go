package util

import (
	"errors"

	"github.com/renjugeo/go-server/config"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var logConfig = zap.Config{
	Encoding:         "json",
	Level:            zap.NewAtomicLevelAt(zapcore.DebugLevel),
	OutputPaths:      []string{"stdout"},
	ErrorOutputPaths: []string{"stderr"},
	EncoderConfig: zapcore.EncoderConfig{
		MessageKey:  "message",
		LevelKey:    "level",
		TimeKey:     "timestamp",
		EncodeTime:  zapcore.ISO8601TimeEncoder,
		EncodeLevel: zapcore.LowercaseLevelEncoder,
	},
}

func GetLogger(cfg *config.Configuration) (*zap.Logger, error) {
	logConfig.Encoding = cfg.LogFormat
	switch cfg.LogLevel {
	case "debug":
		logConfig.Level = zap.NewAtomicLevelAt(zapcore.DebugLevel)
	case "warn":
		logConfig.Level = zap.NewAtomicLevelAt(zapcore.WarnLevel)
	case "info":
		logConfig.Level = zap.NewAtomicLevelAt(zapcore.InfoLevel)
	default:
		return nil, errors.New("unsupported log level")
	}
	return logConfig.Build()
}
