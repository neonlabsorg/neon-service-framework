package api

import (
	"fmt"
	"runtime"

	"github.com/labstack/echo/v4"
	echoMiddleware "github.com/labstack/echo/v4/middleware"
)

type (
	// RecoverConfig defines the config for Recover middleware.
	RecoverConfig struct {
		// Skipper defines a function to skip middleware.
		Skipper echoMiddleware.Skipper

		// Size of the stack to be printed.
		// Optional. Default value 4KB.
		StackSize int `yaml:"stack_size"`

		// DisableStackAll disables formatting stack traces of all other goroutines
		// into buffer after the trace for the current goroutine.
		// Optional. Default value false.
		DisableStackAll bool `yaml:"disable_stack_all"`

		// DisablePrintStack disables printing stack trace.
		// Optional. Default value as false.
		DisablePrintStack bool `yaml:"disable_print_stack"`
	}
)

var (
	// DefaultRecoverConfig is the default Recover middleware config.
	DefaultRecoverConfig = RecoverConfig{
		Skipper:           echoMiddleware.DefaultSkipper,
		StackSize:         4 << 10, // 4 KB
		DisableStackAll:   true,
		DisablePrintStack: false,
	}
)

// Recover returns a middleware which recovers from panics anywhere in the chain
// and handles the control to the centralized HTTPErrorHandler.
func (e *DefaultApiContextExtender) Recover() echo.MiddlewareFunc {
	return e.RecoverWithConfig(DefaultRecoverConfig)
}

type ErrorRecoverWithStackTrace struct {
	Message string
	Stack   string
}

func (e *ErrorRecoverWithStackTrace) Error() string {
	return e.Message
}

func (e *ErrorRecoverWithStackTrace) StackTrace() string {
	return e.Stack
}

// RecoverWithConfig returns a Recover middleware with config.
// See: `Recover()`.
func (e *DefaultApiContextExtender) RecoverWithConfig(config RecoverConfig) echo.MiddlewareFunc {
	// Defaults
	if config.Skipper == nil {
		config.Skipper = DefaultRecoverConfig.Skipper
	}
	if config.StackSize == 0 {
		config.StackSize = DefaultRecoverConfig.StackSize
	}

	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			if config.Skipper(c) {
				return next(c)
			}

			defer func() {
				if r := recover(); r != nil {
					err, ok := r.(error)
					if !ok {
						err = fmt.Errorf("%v", r)
					}
					stack := make([]byte, config.StackSize)
					length := runtime.Stack(stack, !config.DisableStackAll)
					if !config.DisablePrintStack {
						msg := fmt.Sprintf("[PANIC RECOVER] %v %s\n", err, stack[:length])
						e.logger.Error().Err(err).Msg(msg)
					}
					newError := &ErrorRecoverWithStackTrace{
						Message: err.Error(),
						Stack:   string(stack[:length]),
					}
					c.Error(newError)
				}
			}()
			return next(c)
		}
	}
}
