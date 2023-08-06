package alerts

import "github.com/neonlabsorg/neon-service-framework/pkg/errors"

var (
	ErrAdapterWasntInstalled         = errors.Critical.NewWithCode(100201001, "alert adapter wasn't installed")
	ErrReservedAdapterWasntInstalled = errors.Critical.NewWithCode(100201002, "reserved alert adapter wasn't installed")

	ErrPrometheusAlertManagerUrlIsEmpty = errors.Critical.New("prometheus alert manager url is empty")
)
