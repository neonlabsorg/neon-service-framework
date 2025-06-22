package service

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/neonlabsorg/neon-service-framework/pkg/errors"
	"github.com/neonlabsorg/neon-service-framework/pkg/logger"
	"github.com/neonlabsorg/neon-service-framework/pkg/service/configuration"
	"github.com/neonlabsorg/neon-service-framework/pkg/tools/collections"
)

type PostgresManager struct {
	ctx     context.Context
	log     logger.Logger
	configs collections.BasicMapCollection[*configuration.PostgresConfiguration]
	pools   collections.BasicMapCollection[*pgxpool.Pool]
}

func NewPostgresManager(
	ctx context.Context,
	log logger.Logger,
	configs collections.BasicMapCollection[*configuration.PostgresConfiguration],
) *PostgresManager {
	return &PostgresManager{
		ctx:     ctx,
		log:     log,
		configs: configs,
		pools:   make(collections.BasicMapCollection[*pgxpool.Pool]),
	}
}

func (m *PostgresManager) InitConnectionPools() (err error) {
	for name, cfg := range m.configs {
		pool, err := m.initConnectionPool(cfg)
		if err != nil {
			return err
		}
		m.pools.Set(name, pool)
	}

	return nil
}

func (m *PostgresManager) initConnectionPool(cfg *configuration.PostgresConfiguration) (pool *pgxpool.Pool, err error) {
	return pgxpool.New(m.ctx, cfg.BuildConnectionString())
}

func (m *PostgresManager) GetConnectionPool(db string) (pool *pgxpool.Pool, err error) {
	pool, ok := m.pools.Get(db)
	if !ok {
		return nil, errors.NotFound.Newf("connection pool not found: %s", db)
	}

	return pool, nil
}

func (m *PostgresManager) MustGetConnectionPool(db string) (pool *pgxpool.Pool) {
	pool, ok := m.pools.Get(db)
	if !ok {
		panic(fmt.Sprintf("connection pool not found: %s", db))

	}

	return pool
}
