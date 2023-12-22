//

package fsm

import (
	"github.com/go-playground/validator/v10"
	"github.com/go-playground/validator/v10/non-standard/validators"
)

func (obj StateMachine[C]) Validate() error {
	s := struct {
		States  StateMap        `validate:"notblank"`
		Storage StateStorage[C] `validate:"required"`
	}{
		States:  obj.states,
		Storage: obj.storage,
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
