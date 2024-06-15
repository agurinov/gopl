// Code generated: TODO

package trace

import (
	"github.com/go-playground/validator/v10"
	"github.com/go-playground/validator/v10/non-standard/validators"
	"go.opentelemetry.io/otel/sdk/trace"
)

func (obj sampler) Validate() error {
	s := struct {
		RatioSampler trace.Sampler `validate:"required"`
	}{
		RatioSampler: obj.ratioSampler,
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
