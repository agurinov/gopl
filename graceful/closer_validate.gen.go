package graceful

import (
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/go-playground/validator/v10/non-standard/validators"
	"go.uber.org/zap"
)

func (obj Closer) Validate() error {
	s := struct {
		Logger  *zap.Logger   `validate:"required"`
		Timeout time.Duration `validate:"required"`
	}{
		Logger:  obj.logger,
		Timeout: obj.timeout,
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
