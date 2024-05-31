package sql

import (
	"github.com/go-playground/validator/v10"
	"github.com/go-playground/validator/v10/non-standard/validators"
)

func (obj cqrs) Validate() error {
	s := struct {
		RW DB `validate:"notblank"`
		RO DB `validate:"notblank"`
	}{
		RW: obj.rw,
		RO: obj.ro,
	}

	v := validator.New()
	if err := v.RegisterValidation("notblank", validators.NotBlank); err != nil {
		return err
	}

	if err := v.Struct(s); err != nil {
		return err
	}

	return nil
}
