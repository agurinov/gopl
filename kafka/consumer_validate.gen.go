// Code generated: TODO

package kafka

import (
	"github.com/go-playground/validator/v10"
	"github.com/go-playground/validator/v10/non-standard/validators"
	"github.com/twmb/franz-go/pkg/kgo"
	"go.uber.org/zap"
)

func (obj consumer[T]) Validate() error {
	s := struct {
		Logger          *zap.Logger                  `validate:"required"`
		Client          *kgo.Client                  `validate:"required"`
		PartitionHolder *partitionHolder             `validate:"required"`
		Handler         Handler[T]                   `validate:"required_without=HandlerBatch"`
		HandlerBatch    HandlerBatch[T]              `validate:"required_without=Handler"`
		RecordDiscarder RecordDiscarder[*kgo.Record] `validate:"required"`
		RecordMapper    RecordMapper[*kgo.Record, T] `validate:"required"`
	}{
		Logger:          obj.logger,
		Client:          obj.client,
		Handler:         obj.handler,
		HandlerBatch:    obj.handlerBatch,
		PartitionHolder: obj.partitionHolder,
		RecordDiscarder: obj.recordDiscarder,
		RecordMapper:    obj.recordMapper,
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
