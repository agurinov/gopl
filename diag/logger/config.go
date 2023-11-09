package logger

import (
	"github.com/go-playground/validator/v10"

	"github.com/agurinov/gopl/env/envvars"
)

type Config struct {
	Format          string `validate:"oneof=json console"`
	Enabled         bool
	EnableCaller    bool
	EnableTraceback bool
}

// TODO(a.gurinov): must be codegen part
func LoadConfig() (Config, error) {
	var (
		enabled bool
		format  string
	)

	enabled, err := envvars.LogEnabled.Value()
	if err != nil {
		return Config{}, err
	}

	format, err = envvars.LogFormat.Value()
	if err != nil {
		return Config{}, err
	}

	cfg := Config{
		Enabled: enabled,
		Format:  format,
	}

	if validateErr := validator.New().Struct(cfg); validateErr != nil {
		return Config{}, validateErr
	}

	return cfg, nil
}
