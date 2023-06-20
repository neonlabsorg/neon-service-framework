package configuration

import "github.com/neonlabsorg/neon-service-framework/pkg/env"

const DEFAULT_LISTEN_ADDRESS = ":50051"

type GRPCServerConfiguration struct {
	ListenAddr string
}

func (c *ServiceConfiguration) loadGRPCServerConfiguration() (err error) {
	listenAddr := env.Get("NS_GRPC_LISTEN_ADDR")
	if len(listenAddr) == 0 {
		listenAddr = DEFAULT_LISTEN_ADDRESS
	}

	c.GRPCServer = &GRPCServerConfiguration{
		ListenAddr: listenAddr,
	}

	return nil
}
