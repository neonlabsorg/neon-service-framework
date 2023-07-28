package configuration

import (
	"fmt"
	"strings"
	"time"

	"github.com/davecgh/go-spew/spew"
	"github.com/neonlabsorg/neon-service-framework/pkg/env"
	"github.com/neonlabsorg/neon-service-framework/pkg/errors"
)

type ClickhouseConfiguration struct {
	Addr                 []string
	Database             string
	Username             string
	Password             string
	MaxExecutionTime     int
	DialTimeout          time.Duration
	MaxOpenConns         int
	MaxIdleConns         int
	ConnMaxLifetime      time.Duration
	BlockBufferSize      uint
	MaxCompressionBuffer int
	Debug                bool
}

func (c *ServiceConfiguration) loadClickhouseStorageConfig(dbName string) (cfg *ClickhouseConfiguration, err error) {
	config, err := c.loadCommonClickhouseStorageConfig()
	if err != nil {
		return nil, err
	}

	name := strings.ToUpper(dbName)

	config.Addr = env.GetStringList(fmt.Sprintf("NS_DB_CH_%s_NODES", name), ";", config.Addr)
	config.Database = env.Get(fmt.Sprintf("NS_DB_CH_%s_DATABASE", name))
	config.Username = env.Get(fmt.Sprintf("NS_DB_CH_%s_USERNAME", name))
	config.Password = env.Get(fmt.Sprintf("NS_DB_CH_%s_PASSWORD", name))

	config.DialTimeout = env.GetDuration(fmt.Sprintf("NS_DB_CH_%s_DIAL_TIMEOUT", name), config.DialTimeout)
	config.MaxOpenConns = env.GetInt(fmt.Sprintf("NS_DB_CH_%s_MAX_OPEN_CONNS", name), config.MaxOpenConns)
	config.MaxIdleConns = env.GetInt(fmt.Sprintf("NS_DB_CH_%s_MAX_IDLE_CONNS", name), config.MaxIdleConns)
	config.ConnMaxLifetime = env.GetDuration(fmt.Sprintf("NS_DB_CH_%s_CONN_MAX_LIFITIME", name), config.ConnMaxLifetime)
	config.BlockBufferSize = env.GetUint(fmt.Sprintf("NS_DB_CH_%s_BLOCK_BUFFER_SIZE", name), config.BlockBufferSize)
	config.MaxCompressionBuffer = env.GetInt(fmt.Sprintf("NS_DB_CH_%s_MAX_COMPRESSION_BUFFER", name), config.MaxCompressionBuffer)
	config.MaxExecutionTime = env.GetInt(fmt.Sprintf("NS_DB_CH_%s_MAX_EXECUTION_TIME", name), config.MaxExecutionTime)

	if config.Database == "" || len(config.Addr) == 0 {
		return nil, errors.Critical.Newf("invalid env parameters for database '%s': %s", name, spew.Sdump(config))
	}

	return config, nil
}

func (c *ServiceConfiguration) loadCommonClickhouseStorageConfig() (cfg *ClickhouseConfiguration, err error) {
	config := &ClickhouseConfiguration{
		Addr:                 env.GetStringList("NS_DB_CH_NODES", ";"),
		Username:             env.Get("NS_DB_CH_USERNAME", ""),
		Password:             env.Get("NS_DB_CH_PASSWORD", ""),
		DialTimeout:          env.GetDuration("NS_DB_CH_DIAL_TIMEOUT", time.Second*30),
		MaxOpenConns:         env.GetInt("NS_DB_CH_MAX_OPEN_CONNS", 100),
		MaxIdleConns:         env.GetInt("NS_DB_CH_MAX_IDLE_CONNS", 100),
		ConnMaxLifetime:      env.GetDuration("NS_DB_CH_CONN_MAX_LIFITIME", time.Minute*10),
		BlockBufferSize:      env.GetUint("NS_DB_CH_BLOCK_BUFFER_SIZE", 100),
		MaxCompressionBuffer: env.GetInt("NS_DB_CH_MAX_COMPRESSION_BUFFER", 10240),
		MaxExecutionTime:     env.GetInt("NS_DB_CH_MAX_EXECUTION_TIME", 360),
		Debug:                env.GetBool("NS_DB_CH_DEBUG", false),
	}

	return config, nil
}
