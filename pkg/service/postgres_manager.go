package service

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/neonlabsorg/neon-service-framework/pkg/errors"
	"github.com/neonlabsorg/neon-service-framework/pkg/logger"
	"github.com/neonlabsorg/neon-service-framework/pkg/service/configuration"
)

type PostgresPoolCollection map[string]*pgxpool.Pool

func (p PostgresPoolCollection) init() {
	if p == nil {
		p = make(map[string]*pgxpool.Pool)
	}
}

func (p PostgresPoolCollection) Add(name string, pool *pgxpool.Pool) {
	p.init()
	p[name] = pool
}

func (p PostgresPoolCollection) Get(name string) (pool *pgxpool.Pool, ok bool) {
	p.init()
	pool, ok = p[name]
	return pool, ok
}

type PostgresManager struct {
	ctx     context.Context
	log     logger.Logger
	configs configuration.PostgresConfigCollection
	pools   PostgresPoolCollection
}

func NewPostgresManager(
	ctx context.Context,
	log logger.Logger,
	configs configuration.PostgresConfigCollection,
) *PostgresManager {
	return &PostgresManager{
		ctx:     ctx,
		log:     log,
		configs: configs,
		pools:   make(PostgresPoolCollection),
	}
}

func (m *PostgresManager) InitConnectionPools() (err error) {
	for name, cfg := range m.configs {
		pool, err := m.initConnectionPool(cfg)
		if err != nil {
			return err
		}
		m.pools.Add(name, pool)
	}

	return nil
}

func (m *PostgresManager) initConnectionPool(cfg *configuration.PostgresConfiguration) (pool *pgxpool.Pool, err error) {
	return pgxpool.New(m.ctx, cfg.BuildConnectionString())
}

func (m *PostgresManager) GetConnectionPool(db string) (pool *pgxpool.Pool, err error) {
	pool, ok := m.pools.Get(db)
	if !ok {
		return nil, errors.NotFound.New("connection pool not found")
	}

	return pool, nil
}
