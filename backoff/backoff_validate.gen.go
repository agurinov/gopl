//

package backoff

import (
	strategies "github.com/agurinov/gopl/backoff/strategies"
	"github.com/go-playground/validator/v10"
	"github.com/go-playground/validator/v10/non-standard/validators"
)

func (obj Backoff) Validate() error {
	s := struct {
		Strategy   strategies.Interface `validate:"required"`
		MaxRetries uint32               `validate:"min=1"`
	}{
		Strategy:   obj.strategy,
		MaxRetries: obj.maxRetries,
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
