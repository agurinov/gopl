package telegram

import (
	"github.com/go-playground/validator/v10"
	"github.com/go-playground/validator/v10/non-standard/validators"
	"go.uber.org/zap"
)

func (obj Auth) Validate() error {
	s := struct {
		Logger    *zap.Logger       `validate:"required"`
		BotTokens map[string]string `validate:"required,dive,keys,required,endkeys,required"`
	}{
		Logger:    obj.logger,
		BotTokens: obj.botTokens,
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
