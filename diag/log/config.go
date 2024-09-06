package log

import (
	validator "github.com/go-playground/validator/v10"

	c "github.com/agurinov/gopl/patterns/creational"
)

type (
	Config struct {
		Format          string `validate:"oneof=json console"`
		Level           string `validate:"oneof=debug info warn error"`
		EnableCaller    bool   `yaml:"enable_caller"`
		EnableTraceback bool   `yaml:"enable_traceback"`
	}
	Option = c.Option[Config]
)

var newConfig = c.NewWithValidate[Config, Option]

func (obj Config) Validate() error {
	if err := validator.New().Struct(obj); err != nil {
		return err
	}

	return nil
}
