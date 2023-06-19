package configuration

import "github.com/neonlabsorg/neon-service-framework/pkg/env"

const DEFAULT_LISTEN_ADDRESS = ":50051"

type GRPCServerConfig struct {
	ListenAddr string
}

func (c *ServiceConfiguration) loadGRPCServerConfigFromInvironment() (err error) {
	listenAddr := env.Get("NS_GRPC_LISTEN_ADDR")
	if len(listenAddr) == 0 {
		listenAddr = DEFAULT_LISTEN_ADDRESS
	}

	c.GRPCServerConfig = &GRPCServerConfig{
		ListenAddr: listenAddr,
	}

	return nil
}
