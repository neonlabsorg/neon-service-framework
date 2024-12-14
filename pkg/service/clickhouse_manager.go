package service

import (
	"context"
	"fmt"
	"net"

	"github.com/ClickHouse/clickhouse-go/v2"
	"github.com/ClickHouse/clickhouse-go/v2/lib/driver"
	"github.com/neonlabsorg/neon-service-framework/pkg/errors"
	"github.com/neonlabsorg/neon-service-framework/pkg/logger"
	"github.com/neonlabsorg/neon-service-framework/pkg/service/configuration"
	"github.com/neonlabsorg/neon-service-framework/pkg/tools/collections"
)

type ClickhouseManager struct {
	ctx     context.Context
	log     logger.Logger
	configs collections.BasicMapCollection[*configuration.ClickhouseConfiguration]
	conns   collections.BasicMapCollection[driver.Conn]
}

func NewClickhouseManager(
	ctx context.Context,
	log logger.Logger,
	configs collections.BasicMapCollection[*configuration.ClickhouseConfiguration],
) *ClickhouseManager {
	return &ClickhouseManager{
		ctx:     ctx,
		log:     log,
		configs: configs,
		conns:   make(collections.BasicMapCollection[driver.Conn]),
	}
}

func (m *ClickhouseManager) InitConnections() (err error) {
	for name, cfg := range m.configs {
		conn, err := m.initConnection(cfg)
		if err != nil {
			return err
		}
		m.conns.Set(name, conn)
	}

	return nil
}

func (m *ClickhouseManager) initConnection(cfg *configuration.ClickhouseConfiguration) (conn driver.Conn, err error) {
	return clickhouse.Open(&clickhouse.Options{
		Addr: cfg.Addr,
		Auth: clickhouse.Auth{
			Database: cfg.Database,
			Username: cfg.Username,
			Password: cfg.Password,
		},
		DialContext: func(ctx context.Context, addr string) (net.Conn, error) {
			var d net.Dialer
			return d.DialContext(ctx, "tcp", addr)
		},
		Debug: cfg.Debug,
		Debugf: func(format string, v ...any) {
			fmt.Printf(format, v)
		},
		Settings: clickhouse.Settings{
			"max_execution_time": cfg.MaxExecutionTime,
		},
		Compression: &clickhouse.Compression{
			Method: clickhouse.CompressionLZ4,
		},
		DialTimeout:          cfg.DialTimeout,
		MaxOpenConns:         cfg.MaxOpenConns,
		MaxIdleConns:         cfg.MaxIdleConns,
		ConnMaxLifetime:      cfg.ConnMaxLifetime,
		ConnOpenStrategy:     clickhouse.ConnOpenInOrder,
		BlockBufferSize:      uint8(cfg.BlockBufferSize),
		MaxCompressionBuffer: cfg.MaxCompressionBuffer,
	})
}

func (m *ClickhouseManager) GetConnection(db string) (conn driver.Conn, err error) {
	conn, ok := m.conns.Get(db)
	if !ok {
		return nil, errors.NotFound.Newf("clickhouse connection not found: %s", db)
	}

	return conn, nil
}

func (m *ClickhouseManager) MustGetConnection(db string) (conn driver.Conn) {
	conn, ok := m.conns.Get(db)
	if !ok {
		panic(fmt.Sprintf("clickhouse connection not found: %s", db))

	}

	return conn
}
