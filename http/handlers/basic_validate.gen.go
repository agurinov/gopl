// Code generated: TODO

package handlers

import (
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/go-playground/validator/v10/non-standard/validators"
	"go.uber.org/zap"
)

func (obj Basic) Validate() error {
	s := struct {
		Logger   *zap.Logger             `validate:"required"`
		Handlers map[string]http.Handler `validate:"dive,keys,required,endkeys,required"`
	}{
		Logger:   obj.logger,
		Handlers: obj.handlers,
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
