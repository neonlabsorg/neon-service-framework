package configuration

import (
	"fmt"
	"strings"

	"github.com/neonlabsorg/neon-service-framework/pkg/env"
	"github.com/neonlabsorg/neon-service-framework/pkg/errors"
	"github.com/neonlabsorg/neon-service-framework/pkg/tools/collections"
)

const UNITED_API_SERVER = "UNITED"

type ApiServersConfiguration struct {
	Servers collections.BasicMapCollection[*ApiServerConfiguration]
}

type ApiServerConfiguration struct {
	Name       string
	ListenAddr string
	UseCORS    bool
	BodyLimit  string
}

func (c *ServiceConfiguration) loadApiServersConfiguration(useUnitedAPIServer bool, servers []string) (err error) {
	c.ApiServers = &ApiServersConfiguration{}

	if useUnitedAPIServer {
		unitedServer, err := c.loadUnitedApiServerConfiguration()
		if err != nil {
			return err
		}

		c.ApiServers.Servers.Set(UNITED_API_SERVER, unitedServer)
	}

	for _, name := range servers {
		if strings.ToUpper(name) == UNITED_API_SERVER {
			return ErrTryingToUseUnitedApiServerName
		}

		server, err := c.loadApiServerConfigurationByName(name)
		if err != nil {
			return err
		}

		c.ApiServers.Servers.Set(name, server)
	}

	return nil
}

func (c *ServiceConfiguration) loadApiServerConfigurationByName(name string) (cfg *ApiServerConfiguration, err error) {
	name = strings.ToUpper(name)
	cfg = &ApiServerConfiguration{
		Name:       name,
		ListenAddr: env.Get(fmt.Sprintf("NS_API_%s_LISTEN_ADDR", name), ""),
		UseCORS:    env.GetBool(fmt.Sprintf("NS_API_%s_USE_CORS", name), true),
		BodyLimit:  env.Get(fmt.Sprintf("NS_API_%s_BODY_LIMIT", name), "2M"),
	}

	if cfg.ListenAddr == "" {
		return nil, errors.Critical.Wrap(ErrListenAddressForAPIServerisEmpty, "ListenAddress is "+name)
	}

	return cfg, nil
}

func (c *ServiceConfiguration) loadUnitedApiServerConfiguration() (cfg *ApiServerConfiguration, err error) {
	cfg = &ApiServerConfiguration{
		Name:       UNITED_API_SERVER,
		ListenAddr: env.Get("NS_API_LISTEN_ADDR", "0.0.0.0:8080"),
		UseCORS:    env.GetBool("NS_API_USE_CORS", true),
		BodyLimit:  env.Get("NS_API_BODY_LIMIT", "2M"),
	}

	return cfg, nil
}
