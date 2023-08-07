package service

import "github.com/neonlabsorg/neon-service-framework/pkg/errors"

var (
	ErrUnitedApiServerClientAlreadyExists = errors.Critical.NewWithCode(100101001, "united api server client already exists")
	ErrGettingUnexpectedApiServer         = errors.Critical.NewWithCode(100101002, "getting unexpected api server")
	ErrUnitedApiServerIsNotInitialized    = errors.Critical.NewWithCode(100101003, "united api server is not initialized")
	ErrUnregisteredAlertAdapter           = errors.Critical.NewWithCode(100101004, "unregistered alert adapter")
)
