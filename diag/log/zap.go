package log

import (
	"testing"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"go.uber.org/zap/zaptest"
)

func NewZap(opts ...Option) (*zap.Logger, *zap.AtomicLevel, error) {
	cfg, err := newConfig(opts...)
	if err != nil {
		return nil, nil, err
	}

	lvl, err := zap.ParseAtomicLevel(cfg.Level)
	if err != nil {
		return nil, nil, err
	}

	var (
		zapConfig = zap.Config{
			Level:             lvl,
			Encoding:          cfg.Format,
			Development:       false,
			DisableCaller:     !cfg.EnableCaller,
			DisableStacktrace: !cfg.EnableTraceback,
			EncoderConfig: zapcore.EncoderConfig{
				TimeKey:        "@timestamp",
				LevelKey:       "level",
				NameKey:        "logger",
				CallerKey:      "caller",
				MessageKey:     "message",
				StacktraceKey:  "stacktrace",
				LineEnding:     zapcore.DefaultLineEnding,
				EncodeLevel:    zapcore.LowercaseLevelEncoder,
				EncodeTime:     zapcore.ISO8601TimeEncoder,
				EncodeDuration: zapcore.SecondsDurationEncoder,
				EncodeCaller:   zapcore.ShortCallerEncoder,
			},
			OutputPaths:      []string{"stdout"},
			ErrorOutputPaths: []string{"stdout"},
		}
		zapOptions []zap.Option
	)

	if cfg.EnableCaller {
		zapOptions = append(zapOptions,
			zap.AddCaller(),
			zap.AddCallerSkip(1),
		)
	}

	logger, err := zapConfig.Build(zapOptions...)
	if err != nil {
		return nil, nil, err
	}

	return logger, &lvl, nil
}

func MustNewZapSystem() *zap.Logger {
	logger, _, err := NewZap(
		WithLevel("debug"),
		WithFormat("json"),
	)
	if err != nil {
		panic(err)
	}

	return logger
}

func NewZapTest(t *testing.T) *zap.Logger {
	t.Helper()

	// func WrapOptions(zapOpts ...zap.Option) LoggerOption
	return zaptest.NewLogger(t,
		zaptest.Level(zap.DebugLevel),
	)
}
