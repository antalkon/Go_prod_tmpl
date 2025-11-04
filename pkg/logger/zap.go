package logger

import (
	"context"
	"strings"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type ctxKeyType string

const CtxKey ctxKeyType = "logger"

func New(env string) (*zap.Logger, error) {
	envU := strings.ToUpper(strings.TrimSpace(env))

	var cfg zap.Config
	switch envU {
	case "DEV", "DEBUG":
		cfg = zap.NewDevelopmentConfig()
		cfg.Encoding = "console"
		cfg.Level = zap.NewAtomicLevelAt(zap.DebugLevel)
		cfg.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	case "PROD":
		cfg = zap.NewProductionConfig()
		cfg.Encoding = "json"
		cfg.Level = zap.NewAtomicLevelAt(zap.InfoLevel)
	default:
		cfg = zap.NewProductionConfig()
		cfg.Encoding = "json"
		cfg.Level = zap.NewAtomicLevelAt(zap.InfoLevel)
	}

	cfg.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	cfg.EncoderConfig.EncodeDuration = zapcore.StringDurationEncoder

	return cfg.Build(zap.AddCaller(), zap.AddStacktrace(zapcore.ErrorLevel))
}

func Named(l *zap.Logger, name string) *zap.Logger {
	if l == nil {
		return zap.NewNop()
	}
	return l.Named(name)
}

func WithContext(ctx context.Context, l *zap.Logger) context.Context {
	return context.WithValue(ctx, CtxKey, l)
}

func FromContext(ctx context.Context, fallback *zap.Logger) *zap.Logger {
	if ctx == nil {
		if fallback != nil {
			return fallback
		}
		return zap.NewNop()
	}
	if v := ctx.Value(CtxKey); v != nil {
		if lg, ok := v.(*zap.Logger); ok {
			return lg
		}
	}
	if fallback != nil {
		return fallback
	}
	return zap.NewNop()
}
