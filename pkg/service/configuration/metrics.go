package configuration

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/neonlabsorg/neon-service-framework/pkg/env"
)

type MetricsServerConfiguration struct {
	ServiceName   string
	ListenAddress string
	ListenPort    int
	Interval      int
}

// LOAD METRICS CONFIGURATION
func (c *ServiceConfiguration) loadMetricsServerConfiguration(serviceName string) error {
	serviceName = strings.ToUpper(serviceName)
	listenAddr := env.Get(fmt.Sprintf("NS_METRICS_%s_LISTEN_ADDRESS", serviceName))
	listenPortString := env.Get(fmt.Sprintf("NS_METRICS_%s_LISTEN_PORT", serviceName))
	intervalString := env.Get(fmt.Sprintf("NS_METRICS_%s_INTERVAL", serviceName))

	port, err := strconv.Atoi(listenPortString)
	if err != nil {
		port = 0
	}

	interval, err := strconv.Atoi(intervalString)
	if err != nil {
		interval = 0
	}

	cfg := &MetricsServerConfiguration{
		ServiceName:   serviceName,
		ListenAddress: listenAddr,
		ListenPort:    port,
		Interval:      interval,
	}

	c.MetricsServer = cfg

	return nil
}
