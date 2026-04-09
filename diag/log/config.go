package log

import (
	"go.uber.org/zap"

	c "github.com/agurinov/gopl/patterns/creational"
)

type (
	Config struct {
		Format          string `validate:"oneof=json console"`
		Level           string `validate:"oneof=debug info warn error"`
		zapFields       []zap.Field
		EnableCaller    bool `json:"enable_caller" yaml:"enable_caller"`
		EnableTraceback bool `json:"enable_traceback" yaml:"enable_traceback"`
	}
	Option = c.Option[Config]
)

func (cfg Config) New(opts ...Option) (*zap.Logger, *zap.AtomicLevel, error) {
	defaults := []Option{
		WithFormat(cfg.Format),
		WithLevel(cfg.Level),
		WithCaller(cfg.EnableCaller),
		WithTraceback(cfg.EnableTraceback),
		WithFields(cfg.zapFields...),
	}

	opts = append(defaults, opts...)

	return NewZap(opts...)
}
