package api

import "github.com/neonlabsorg/neon-service-framework/pkg/errors"

var (
	ErrInternal     = errors.Internal.NewWithCode(100000001, "internal error")
	ErrUnauthorized = errors.Unauthorized.NewWithCode(100000002, "unauthorized error")
	ErrValidation   = errors.Validation.NewWithCode(100000003, "validation error")
)
