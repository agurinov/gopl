package graceful

import (
	"context"
	"sync/atomic"

	"github.com/go-playground/validator/v10"
	"github.com/go-playground/validator/v10/non-standard/validators"
	"go.uber.org/zap"
)

func (obj Wrapper) Validate() error {
	s := struct {
		Logger      *zap.Logger        `validate:"required"`
		Ctx         context.Context    `validate:"required"`
		ForceCancel context.CancelFunc `validate:"required"`
		Closed      *atomic.Bool       `validate:"required"`
	}{
		Logger:      obj.logger,
		Ctx:         obj.ctx,
		ForceCancel: obj.forceCancel,
		Closed:      obj.closed,
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
