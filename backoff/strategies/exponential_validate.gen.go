package strategies

import (
	"time"

	validator "github.com/go-playground/validator/v10"
	"github.com/go-playground/validator/v10/non-standard/validators"
)

func (obj exponential) Validate() error {
	s := struct {
		MinDelay   time.Duration `validate:"min=0s"`
		MaxDelay   time.Duration `validate:"min=1s"`
		Multiplier float64       `validate:"gte=1.0"`
		Jitter     float64       `validate:"gte=0.00,lte=1.0"`
	}{
		MinDelay:   obj.minDelay,
		MaxDelay:   obj.maxDelay,
		Multiplier: obj.multiplier,
		Jitter:     obj.jitter,
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
