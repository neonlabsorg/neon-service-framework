package service

import (
	"context"

	"github.com/neonlabsorg/neon-service-framework/pkg/logger"
	"github.com/neonlabsorg/neon-service-framework/pkg/service/configuration"
)

type DatabaseManager struct {
	ctx                  context.Context
	storageConfiguration *configuration.StorageConfiguration
	logger               logger.Logger
	postgresManager      *PostgresManager
	clickhouseManager    *ClickhouseManager
}

func NewDatabaseManager(
	ctx context.Context,
	storageConfiguration *configuration.StorageConfiguration,
	log logger.Logger,
) (manager *DatabaseManager, err error) {
	manager = &DatabaseManager{
		ctx:                  ctx,
		storageConfiguration: storageConfiguration,
		logger:               log,
	}

	manager.postgresManager = NewPostgresManager(ctx, log, storageConfiguration.Postgres)

	err = manager.postgresManager.InitConnectionPools()
	if err != nil {
		return nil, err
	}

	manager.clickhouseManager = NewClickhouseManager(ctx, log, storageConfiguration.Clichouse)
	err = manager.clickhouseManager.InitConnections()
	if err != nil {
		return nil, err
	}

	return manager, nil
}

func (m *DatabaseManager) GetPostgresManager() *PostgresManager {
	return m.postgresManager
}

func (m *DatabaseManager) GetClickhouseManager() *ClickhouseManager {
	return m.clickhouseManager
}
