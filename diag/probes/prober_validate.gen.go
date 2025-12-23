// Code generated: TODO

package probes

import (
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/go-playground/validator/v10/non-standard/validators"
	"go.uber.org/zap"

	"github.com/agurinov/gopl/run"
)

func (obj *Prober) Validate() error {
	s := struct {
		Logger          *zap.Logger   `validate:"required"`
		ReadinessProbes []run.Fn      `validate:"dive,required"`
		LivenessProbes  []run.Fn      `validate:"dive,required"`
		CheckInterval   time.Duration `validate:"min=1s"`
		CheckTimeout    time.Duration `validate:"min=1s"`
	}{
		Logger:          obj.logger,
		ReadinessProbes: obj.readinessProbes,
		LivenessProbes:  obj.livenessProbes,
		CheckInterval:   obj.checkInterval,
		CheckTimeout:    obj.checkTimeout,
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
