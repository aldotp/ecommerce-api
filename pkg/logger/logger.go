package logger

import (
	"fmt"
	"strings"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func InitLogger(env string) (*zap.Logger, error) {
	var cfg zap.Config

	switch strings.ToLower(env) {
	case "local":
		cfg = zap.Config{
			Encoding:         "console",
			OutputPaths:      []string{"stdout"},
			ErrorOutputPaths: []string{"stderr"},
			Level:            zap.NewAtomicLevelAt(zap.DebugLevel),
			EncoderConfig: zapcore.EncoderConfig{
				MessageKey:     "message",
				LevelKey:       "level",
				TimeKey:        "time",
				CallerKey:      "caller",
				EncodeDuration: zapcore.StringDurationEncoder,
				EncodeLevel:    zapcore.CapitalColorLevelEncoder,
				EncodeTime:     zapcore.ISO8601TimeEncoder,
				EncodeCaller:   zapcore.ShortCallerEncoder,
			},
		}

	case "development", "staging":
		cfg = zap.Config{
			Encoding:         "json",
			OutputPaths:      []string{"stdout"},
			ErrorOutputPaths: []string{"stderr"},
			Level:            zap.NewAtomicLevelAt(zap.DebugLevel),
			EncoderConfig: zapcore.EncoderConfig{
				MessageKey:     "message",
				LevelKey:       "level",
				TimeKey:        "time",
				CallerKey:      "caller",
				EncodeDuration: zapcore.StringDurationEncoder,
				EncodeLevel:    zapcore.LowercaseLevelEncoder,
				EncodeTime:     zapcore.ISO8601TimeEncoder,
				EncodeCaller:   zapcore.ShortCallerEncoder,
			},
		}

	case "production":
		cfg = zap.Config{
			Encoding:         "json",
			OutputPaths:      []string{"stdout"},
			ErrorOutputPaths: []string{"stderr"},
			Level:            zap.NewAtomicLevelAt(zap.InfoLevel),
			EncoderConfig: zapcore.EncoderConfig{
				MessageKey:     "message",
				LevelKey:       "level",
				TimeKey:        "time",
				CallerKey:      "caller",
				EncodeDuration: zapcore.StringDurationEncoder,
				EncodeLevel:    zapcore.LowercaseLevelEncoder,
				EncodeTime:     zapcore.ISO8601TimeEncoder,
				EncodeCaller:   zapcore.ShortCallerEncoder,
			},
		}

	default:
		panic(fmt.Sprintf("Unknown app environment: %s", env))
	}

	return cfg.Build()
}
