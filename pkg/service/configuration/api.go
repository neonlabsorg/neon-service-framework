package configuration

import "github.com/neonlabsorg/neon-service-framework/pkg/env"

type ApiServerConfiguration struct {
	ListenAddr string
	UseCORS    bool
	BodyLimit  string
}

func (c *ServiceConfiguration) loadApiServerConfiguration() (err error) {
	c.ApiServer = &ApiServerConfiguration{
		ListenAddr: env.Get("NS_API_LISTEN_ADDR", "0.0.0.0:8080"),
		UseCORS:    env.GetBool("NS_API_USE_CORS", true),
		BodyLimit:  env.Get("NS_API_BODY_LIMIT", "2M"),
	}

	return nil
}
