// Code generated: TODO

package kafka

import (
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/go-playground/validator/v10/non-standard/validators"
	"go.uber.org/zap"

	irun "github.com/agurinov/gopl/internal/run"
)

func (obj consumer) Validate() error {
	s := struct {
		RecordDiscarder     RecordDiscarder        `validate:"required"`
		RecordMapper        RecordMapper           `validate:"required"`
		Logger              *zap.Logger            `validate:"required"`
		PartitionDispatcher irun.Dispatcher[int32] `validate:"required"`
		Handler             Handler                `validate:"required_without=HandlerBatch"`
		HandlerBatch        HandlerBatch           `validate:"required_without=Handler"`
		Topic               string                 `validate:"required"`
		MaxPollRecords      int                    `validate:"required"`
		MaxPollDuration     time.Duration          `validate:"required"`
	}{
		Logger:              obj.logger,
		Handler:             obj.handler,
		HandlerBatch:        obj.handlerBatch,
		PartitionDispatcher: obj.partitionDispatcher,
		RecordDiscarder:     obj.recordDiscarder,
		RecordMapper:        obj.recordMapper,
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
