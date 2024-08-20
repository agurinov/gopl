package crontab

import (
	"time"

	gocron "github.com/go-co-op/gocron/v2"
	"github.com/go-playground/validator/v10"
	"github.com/go-playground/validator/v10/non-standard/validators"
	"go.uber.org/zap"
)

func (obj Scheduler) Validate() error {
	s := struct {
		Logger          *zap.Logger      `validate:"required"`
		Scheduler       gocron.Scheduler `validate:"required"`
		Jobs            map[string]Job   `validate:"gt=0,dive,keys,required,endkeys,required"`
		ShutdownTimeout time.Duration    `validate:"required"`
	}{
		Logger:          obj.logger,
		Scheduler:       obj.scheduler,
		Jobs:            obj.jobs,
		ShutdownTimeout: obj.shutdownTimeout,
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
