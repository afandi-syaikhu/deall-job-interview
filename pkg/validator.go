package pkg

import (
	"github.com/go-playground/validator"
)

type Validator struct {
	validator *validator.Validate
}

func (_v *Validator) Validate(i interface{}) error {
	err := _v.validator.Struct(i)
	if err != nil {
		return err
	}

	return nil
}

func GetValidator() *Validator {
	return &Validator{
		validator: validator.New(),
	}
}
