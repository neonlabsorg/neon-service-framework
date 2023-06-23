package configuration

import (
	"fmt"
	"strings"

	"github.com/davecgh/go-spew/spew"
	"github.com/neonlabsorg/neon-service-framework/pkg/env"
	"github.com/neonlabsorg/neon-service-framework/pkg/errors"
)

// POSTGRESQL DATABASE
type PostgresConfiguration struct {
	Hostname string
	Port     int
	SSLMode  string
	Username string
	Password string
	Database string
}

func (c *PostgresConfiguration) BuildConnectionString() string {
	return fmt.Sprintf("host=%s port=%d user=%s dbname=%s sslmode=%s password=%s",
		c.Hostname, c.Port, c.Username, c.Database, c.SSLMode, c.Password)
}

// LOAD POSTGRES CONFIGURATION
func (c *ServiceConfiguration) loadPostgresStorageConfig(name string) (cfg *PostgresConfiguration, err error) {
	cfg = c.loadDefaultPostgresStorageConfig()

	name = strings.ToUpper(name)

	cfg.Hostname = env.Get(fmt.Sprintf("NS_DB_PG_%s_HOSTNAME", name), cfg.Hostname)
	cfg.Port = env.GetInt(fmt.Sprintf("NS_DB_PG_%s_PORT", name), cfg.Port)
	cfg.SSLMode = env.Get(fmt.Sprintf("NS_DB_PG_%s_SSLMODE", name), cfg.SSLMode)
	cfg.Username = env.Get(fmt.Sprintf("NS_DB_PG_%s_USERNAME", name), cfg.Username)
	cfg.Password = env.Get(fmt.Sprintf("NS_DB_PG_%s_PASSWORD", name), cfg.Password)
	cfg.Database = env.Get(fmt.Sprintf("NS_DB_PG_%s_DATABASE", name), cfg.Database)

	if cfg.Database == "" || cfg.Hostname == "" || cfg.Username == "" {
		return nil, errors.Critical.Newf("invalid env parameters for database '%s': %s", name, spew.Sdump(cfg))
	}

	return cfg, nil
}

func (c *ServiceConfiguration) loadDefaultPostgresStorageConfig() (cfg *PostgresConfiguration) {
	return &PostgresConfiguration{
		Hostname: env.Get("NS_DB_PG_HOSTNAME"),
		Port:     env.GetInt("NS_DB_PG_PORT", 5432),
		SSLMode:  env.Get("NS_DB_PG_SSLMODE", "disable"),
		Username: env.Get("NS_DB_PG_USERNAME"),
		Password: env.Get("NS_DB_PG_PASSWORD"),
	}
}

type PostgresConfigCollection map[string]*PostgresConfiguration

func (c *PostgresConfigCollection) init() {
	if c == nil {
		*c = make(map[string]*PostgresConfiguration)
	}
}

func (c PostgresConfigCollection) Add(name string, config *PostgresConfiguration) {
	c.init()
	c[name] = config
}

func (c PostgresConfigCollection) Get(name string) (config *PostgresConfiguration, ok bool) {
	c.init()
	config, ok = c[name]
	return config, ok
}
