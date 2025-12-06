// Code generated: TODO

package kafka

import (
	"sync/atomic"

	"github.com/go-playground/validator/v10"
	"github.com/go-playground/validator/v10/non-standard/validators"
	"github.com/twmb/franz-go/pkg/kgo"
	"go.uber.org/zap"
)

func (obj consumer) Validate() error {
	s := struct {
		Logger          *zap.Logger       `validate:"required"`
		Closed          *atomic.Bool      `validate:"required"`
		Client          *kgo.Client       `validate:"required"`
		Discarder       consumerDiscarder `validate:"required"`
		PartitionHolder *partitionHolder  `validate:"required"`
		Handler         Handler           `validate:"required_without=HandlerBatch"`
		HandlerBatch    HandlerBatch      `validate:"required_without=Handler"`
	}{
		Logger:          obj.logger,
		Closed:          obj.closed,
		Client:          obj.client,
		Discarder:       obj.discarder,
		Handler:         obj.handler,
		HandlerBatch:    obj.handlerBatch,
		PartitionHolder: obj.partitionHolder,
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
