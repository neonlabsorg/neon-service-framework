package service

import (
	"context"

	"github.com/neonlabsorg/neon-service-framework/pkg/api"
	"github.com/neonlabsorg/neon-service-framework/pkg/logger"
	"github.com/neonlabsorg/neon-service-framework/pkg/service/configuration"
	"github.com/neonlabsorg/neon-service-framework/pkg/tools/collections"
)

type ApiServerManager struct {
	cfg      *configuration.ApiServersConfiguration
	ctx      context.Context
	logger   logger.Logger
	extender api.ApiContextExtender
	united   *UnitedApiServer
	servers  collections.BasicMapCollection[*ApiServer]
}

func NewApiServerManager(
	cfg *configuration.ApiServersConfiguration,
	ctx context.Context,
	logger logger.Logger,
	extender api.ApiContextExtender,
) *ApiServerManager {
	return &ApiServerManager{
		cfg:      cfg,
		ctx:      ctx,
		logger:   logger,
		extender: extender,
		servers:  make(collections.BasicMapCollection[*ApiServer]),
	}
}

func (m *ApiServerManager) Init() (err error) {
	for _, cfg := range m.cfg.Servers {

		srv := NewApiServer(
			cfg.Name,
			m.ctx,
			cfg,
			m.extender,
			m.logger,
		)

		if cfg.Name == configuration.UNITED_API_SERVER {
			m.united = NewUnitedApiServer(srv)
		} else {
			m.servers.Set(srv.name, srv)
		}
	}

	return nil
}

func (m *ApiServerManager) GetApiServerByName(name string) (srv *ApiServer, err error) {
	srv, ok := m.servers.Get(name)
	if !ok {
		return nil, ErrGettingUnexpectedApiServer
	}

	return srv, nil
}

func (m *ApiServerManager) GetUnitedApiServer() (srv *UnitedApiServer, err error) {
	if m.united == nil {
		return nil, ErrUnitedApiServerIsNotInitialized
	}

	return m.united, nil
}
