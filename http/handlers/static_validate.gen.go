// Code generated: TODO

package handlers

import (
	"io/fs"

	"github.com/go-playground/validator/v10"
	"github.com/go-playground/validator/v10/non-standard/validators"
	"go.uber.org/zap"
)

func (obj static) Validate() error {
	s := struct {
		Logger *zap.Logger `validate:"required"`
		FS     fs.FS       `validate:"required"`
	}{
		Logger: obj.logger,
		FS:     obj.fs,
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
