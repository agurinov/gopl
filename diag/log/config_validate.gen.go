package log

import validator "github.com/go-playground/validator/v10"

func (obj Config) Validate() error {
	if err := validator.New().Struct(obj); err != nil {
		return err
	}

	return nil
}
