package config

import (
	"github.com/go-playground/validator/v10"
	"github.com/go-playground/validator/v10/non-standard/validators"
)

func validateStruct(cfg any) error {
	v := validator.New()
	if err := v.RegisterValidation("notblank", validators.NotBlank); err != nil {
		return err
	}

	if err := v.Struct(cfg); err != nil {
		return err
	}

	return nil
}
