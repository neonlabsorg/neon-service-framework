package configuration

import (
	"fmt"
	"strings"
	"time"

	"github.com/neonlabsorg/neon-service-framework/pkg/env"
)

type MetricsServerConfiguration struct {
	Enable        bool
	ServiceName   string
	ListenAddress string
	ListenPort    int
	Interval      time.Duration
}

// LOAD METRICS CONFIGURATION
func (c *ServiceConfiguration) loadMetricsServerConfiguration(serviceName string) (err error) {
	cfg := c.loadDefaultMetricsServerConfiguration()

	serviceName = strings.ToUpper(serviceName)
	cfg.Enable = env.GetBool(fmt.Sprintf("NS_METRICS_%s_ENABLE", serviceName), cfg.Enable)
	cfg.ListenAddress = env.Get(fmt.Sprintf("NS_METRICS_%s_LISTEN_ADDRESS", serviceName), cfg.ListenAddress)
	cfg.ListenPort = env.GetInt(fmt.Sprintf("NS_METRICS_%s_LISTEN_PORT", serviceName), cfg.ListenPort)
	cfg.Interval = env.GetDuration(fmt.Sprintf("NS_METRICS_%s_INTERVAL", serviceName), cfg.Interval)

	c.MetricsServer = cfg

	return nil
}

func (c *ServiceConfiguration) loadDefaultMetricsServerConfiguration() (cfg *MetricsServerConfiguration) {
	return &MetricsServerConfiguration{
		Enable:        env.GetBool("NS_METRICS_ENABLE", false),
		ListenAddress: env.Get("NS_METRICS_LISTEN_ADDRESS", "0.0.0.0"),
		ListenPort:    env.GetInt("NS_METRICS_LISTEN_PORT", 20500),
		Interval:      env.GetDuration("NS_METRICS_INTERVAL", time.Second*5),
	}
}
