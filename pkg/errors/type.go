package errors

import (
	"fmt"

	"github.com/pkg/errors"
)

type ErrorType uint

func (t ErrorType) New(msg string) error {
	return Error{errorType: t, originalError: errors.New(msg)}
}

func (t ErrorType) Newf(msg string, args ...interface{}) error {
	err := fmt.Errorf(msg, args...)

	return Error{errorType: t, originalError: err}
}

func (t ErrorType) NewWithCode(code ErrorCode, msg string) error {
	return Error{code: code, errorType: t, originalError: errors.New(msg)}
}

func (t ErrorType) NewfWithCode(code ErrorCode, msg string, args ...interface{}) error {
	err := fmt.Errorf(msg, args...)

	return Error{code: code, errorType: t, originalError: err}
}

func (t ErrorType) Wrap(err error, msg string) error {
	return t.Wrapf(err, msg)
}

func (t ErrorType) Wrapf(err error, msg string, args ...interface{}) error {
	newErr := errors.Wrapf(err, msg, args...)

	return Error{errorType: t, originalError: newErr}
}

func (t ErrorType) WrapWithCode(code ErrorCode, err error, msg string) error {
	return t.WrapfWithCode(code, err, msg)
}

func (t ErrorType) WrapfWithCode(code ErrorCode, err error, msg string, args ...interface{}) error {
	newErr := errors.Wrapf(err, msg, args...)

	return Error{code: code, errorType: t, originalError: newErr}
}
