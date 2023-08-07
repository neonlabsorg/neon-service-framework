package configuration

import (
	"time"

	"github.com/neonlabsorg/neon-service-framework/pkg/env"
)

type AlertsConfiguration struct {
	MainAdapter    string
	ReserveAdapter string
	Prometheus     *PrometheusAdapterConfig
}

type PrometheusAdapterConfig struct {
	URL      string
	Attempts int
	Interval time.Duration
	IsDemo   bool
}

func NewPrometheusAdapterConfigFromEnv() (cfg *PrometheusAdapterConfig, err error) {
	cfg = &PrometheusAdapterConfig{
		URL:      env.Get("NS_ALERTS_PROMETHEUS_ALERT_MANAGER_URL"),
		Attempts: env.GetInt("NS_ALERTS_PROMETHEUS_ATTEMPTS", 5),
		Interval: env.GetDuration("NS_ALERTS_PROMETHEUS_INTERVAL", time.Duration(time.Second)),
		IsDemo:   env.GetBool("NS_ALERTS_PROMETHEUS_DEMO_MODE", false),
	}

	return cfg, nil
}

func (c *ServiceConfiguration) loadAlertsConfiguration() (err error) {
	c.Alerts = &AlertsConfiguration{
		MainAdapter:    env.Get("NS_ALERTS_MAIN_ADAPTER", "prometheus"),
		ReserveAdapter: env.Get("NS_ALERTS_RESERVE_ADAPTER", "console"),
	}

	prometheusConfiguration, err := NewPrometheusAdapterConfigFromEnv()
	if err != nil {
		return err
	}

	c.Alerts.Prometheus = prometheusConfiguration

	return nil
}
