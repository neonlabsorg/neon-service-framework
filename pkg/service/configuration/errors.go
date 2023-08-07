package configuration

import "github.com/neonlabsorg/neon-service-framework/pkg/errors"

var (
	ErrTryingToUseUnitedApiServerName   = errors.Critical.NewWithCode(100102001, "trying to use united API server name")
	ErrListenAddressForAPIServerisEmpty = errors.Critical.NewWithCode(100102002, "listen address for API server is empty")
	ErrPrometheusAlertManagerUrlIsEmpty = errors.Critical.NewWithCode(100102003, "prometheus alert manager url is empty")
)
