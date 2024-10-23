// Code generated: TODO

package handlers

import (
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/go-playground/validator/v10/non-standard/validators"
	"go.uber.org/zap"
)

func (obj basic) Validate() error {
	s := struct {
		Logger  *zap.Logger  `validate:"required"`
		Handler http.Handler `validate:"required"`
	}{
		Logger:  obj.logger,
		Handler: obj.handler,
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
