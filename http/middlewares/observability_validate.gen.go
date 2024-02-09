package middlewares

import (
	"github.com/go-playground/validator/v10"
	"github.com/go-playground/validator/v10/non-standard/validators"
)

func (obj observability) Validate() error {
	s := struct {
	}{}

	v := validator.New()
	if err := v.RegisterValidation("notblank", validators.NotBlank); err != nil {
		return err
	}

	if err := v.Struct(s); err != nil {
		return err
	}

	return nil
}
