package api

import (
	"github.com/labstack/echo/v4"
	"github.com/neonlabsorg/neon-service-framework/pkg/logger"
)

type DefaultApiContext struct {
	echo.Context
	validator *Validator
	logger    logger.Logger
}

func (c *DefaultApiContext) GetValidator() *Validator {
	return c.validator
}

func (c *DefaultApiContext) bindAndValidateModel(model interface{}) error {
	if err := c.Bind(model); err != nil {
		return err
	}
	if err := c.Validate(model); err != nil && err != echo.ErrValidatorNotRegistered {
		return err
	}
	return nil
}

func (c *DefaultApiContext) BindAndValidateModel(model interface{}) error {
	if err := c.bindAndValidateModel(model); err != nil {
		return err
	}
	if err := c.GetValidator().Validate(model); err != nil {
		return err
	}
	return nil
}

func ExtendApiContext(echoCtx echo.Context) *DefaultApiContext {
	return echoCtx.(*DefaultApiContext)
}

type ApiContextExtender interface {
	ExtendDefaultApiContext(h echo.HandlerFunc) echo.HandlerFunc
}

type DefaultApiContextExtender struct {
	validatorInstance *Validator
	logger            logger.Logger
}

func NewDefaultApiContextExtender(
	validatorInstance *Validator,
	log logger.Logger,
) *DefaultApiContextExtender {
	return &DefaultApiContextExtender{
		validatorInstance: validatorInstance,
		logger:            log,
	}
}

func (e *DefaultApiContextExtender) ExtendDefaultApiContext(h echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) (err error) {
		ctx := &DefaultApiContext{
			Context:   c,
			validator: e.validatorInstance,
			logger:    e.logger,
		}

		return h(ctx)
	}
}
