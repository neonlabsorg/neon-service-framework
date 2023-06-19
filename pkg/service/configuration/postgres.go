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
	postgresConfig := &PostgresConfiguration{}

	name = strings.ToUpper(name)

	postgresConfig.Hostname = env.Get(fmt.Sprintf("NS_DB_PG_%s_HOSTNAME", name))
	postgresConfig.Port = env.GetInt(fmt.Sprintf("NS_DB_PG_%s_PORT", name), 5432)
	postgresConfig.SSLMode = env.Get(fmt.Sprintf("NS_DB_PG_%s_SSLMODE", name))
	postgresConfig.Username = env.Get(fmt.Sprintf("NS_DB_PG_%s_USERNAME", name))
	postgresConfig.Password = env.Get(fmt.Sprintf("NS_DB_PG_%s_PASSWORD", name))
	postgresConfig.Database = env.Get(fmt.Sprintf("NS_DB_PG_%s_DATABASE", name))

	if postgresConfig.Database == "" || postgresConfig.Hostname == "" || postgresConfig.Username == "" {
		return nil, errors.Critical.Newf("invalid env parameters for database '%s': %s", name, spew.Sdump(postgresConfig))
	}

	return postgresConfig, nil
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
