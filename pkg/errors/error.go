package errors

import (
	"fmt"

	"github.com/pkg/errors"
)

const (
	NoType = ErrorType(iota)
	Validation
	NotFound
	AccessDenied
	Unauthorized
	Logical
	Temporarily
	Internal
	Critical
)

type ErrorCode uint64

type Error struct {
	code          ErrorCode
	errorType     ErrorType
	originalError error
	context       ErrorContext
}

func (e Error) Error() string {
	return e.originalError.Error()
}

func (e *Error) GetType() ErrorType {
	return e.errorType
}

func (e *Error) GetCode() ErrorCode {
	return e.code
}

func (e *Error) AddToContext(key string, value string) {
	e.context.Set(key, value)
}

func (e *Error) GetContextValueByKey(key string) {
	e.context.Get(key)
}

func (e *Error) GetContext() ErrorContext {
	return e.context
}

func (e *Error) ClearContext() {
	e.context.Clear()
}

func (e Error) Cause() error {
	return errors.Cause(e)
}

func New(msg string) error {
	return Error{errorType: NoType, originalError: errors.New(msg)}
}

func Newf(msg string, args ...interface{}) error {
	return Error{errorType: NoType, originalError: errors.New(fmt.Sprintf(msg, args...))}
}

func NewWithCode(code ErrorCode, msg string) error {
	return Error{code: code, errorType: NoType, originalError: errors.New(msg)}
}

func NewfWithCode(code ErrorCode, msg string, args ...interface{}) error {
	return Error{code: code, errorType: NoType, originalError: errors.New(fmt.Sprintf(msg, args...))}
}

func Wrap(err error, msg string) error {
	return Wrapf(err, msg)
}

func WrapWithCode(code ErrorCode, err error, msg string) error {
	return WrapfWithCode(code, err, msg)
}

func Cause(err error) error {
	return errors.Cause(err)
}

func Wrapf(err error, msg string, args ...interface{}) error {
	wrappedError := errors.Wrapf(err, msg, args...)
	if customErr, ok := err.(Error); ok {
		return Error{
			code:          customErr.code,
			errorType:     customErr.errorType,
			originalError: wrappedError,
			context:       customErr.context,
		}
	}

	return Error{errorType: NoType, originalError: wrappedError}
}

func WrapfWithCode(code ErrorCode, err error, msg string, args ...interface{}) error {
	wrappedError := errors.Wrapf(err, msg, args...)
	if customErr, ok := err.(Error); ok {
		return Error{
			code:          customErr.code,
			errorType:     customErr.errorType,
			originalError: wrappedError,
			context:       customErr.context,
		}
	}

	return Error{code: code, errorType: NoType, originalError: wrappedError}
}

func GetType(err error) ErrorType {
	if customErr, ok := err.(Error); ok {
		return customErr.errorType
	}

	return NoType
}
