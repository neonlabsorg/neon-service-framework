package api

import (
	"reflect"

	"gopkg.in/go-playground/validator.v9"
)

type Validator struct {
	validator *validator.Validate
}

func (v *Validator) Validate(i interface{}) error {
	err := v.validator.Struct(i)
	if err != nil {
		return NewValidatorError(err, i)
	}
	return nil
}

func (v *Validator) Var(i interface{}, tag string) error {
	err := v.validator.Var(i, tag)
	if err != nil {
		return NewValidatorError(err, nil)
	}
	return nil
}

func (v *Validator) RegisterValidation(tag string, fn Func, callValidationEvenIfNull ...bool) error {
	return v.validator.RegisterValidation(tag, func(fl validator.FieldLevel) bool {
		return fn(fl.Field())
	}, callValidationEvenIfNull...)
}

type Func func(value reflect.Value) bool

func NewValidator() *Validator {
	return &Validator{validator: validator.New()}
}

func NewValidatorError(err error, model interface{}) error {
	return &ValidateError{err, model}
}

type ValidateError struct {
	Veer  error
	Model interface{}
}

func (e *ValidateError) Error() string {
	return e.Veer.Error()
}
