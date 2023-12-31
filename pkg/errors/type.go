package errors

import (
	"fmt"

	"github.com/pkg/errors"
)

type ErrorType uint

func (t ErrorType) String() string {
	switch t {
	case Validation:
		return "Validation"
	case NotFound:
		return "NotFound"
	case AccessDenied:
		return "AccessDenied"
	case Unauthorized:
		return "Unauthorized"
	case Logical:
		return "Logical"
	case Temporarily:
		return "Temporarily"
	case Internal:
		return "Internal"
	case Critical:
		return "Critical"
	default:
		return "no_type"
	}
}

func (t ErrorType) New(msg string) Error {
	return Error{errorType: t, originalError: errors.New(msg)}
}

func (t ErrorType) Newf(msg string, args ...interface{}) Error {
	err := fmt.Errorf(msg, args...)

	return Error{errorType: t, originalError: err}
}

func (t ErrorType) NewWithCode(code ErrorCode, msg string) Error {
	return Error{code: code, errorType: t, originalError: errors.New(msg)}
}

func (t ErrorType) NewfWithCode(code ErrorCode, msg string, args ...interface{}) Error {
	err := fmt.Errorf(msg, args...)

	return Error{code: code, errorType: t, originalError: err}
}

func (t ErrorType) Wrap(err error, msg string) Error {
	return t.Wrapf(err, msg)
}

func (t ErrorType) Wrapf(err error, msg string, args ...interface{}) Error {
	newErr := errors.Wrapf(err, msg, args...)

	return Error{errorType: t, originalError: newErr}
}

func (t ErrorType) WrapWithCode(code ErrorCode, err error, msg string) Error {
	return t.WrapfWithCode(code, err, msg)
}

func (t ErrorType) WrapfWithCode(code ErrorCode, err error, msg string, args ...interface{}) Error {
	newErr := errors.Wrapf(err, msg, args...)

	return Error{code: code, errorType: t, originalError: newErr}
}
