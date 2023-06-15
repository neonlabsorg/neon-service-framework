package configuration

import (
	"fmt"
	"os"
	"strings"

	"github.com/davecgh/go-spew/spew"
)

// POSTGRESQL DATABASE
type PostgresConfiguration struct {
	Hostname string
	Port     string
	SSLMode  string
	Username string
	Password string
	Database string
}

func (c *PostgresConfiguration) BuildConnectionString() string {
	return fmt.Sprintf("host=%s port=%s user=%s dbname=%s sslmode=%s password=%s",
		c.Hostname, c.Port, c.Username, c.Database, c.SSLMode, c.Password)
}

// LOAD POSTGRES CONFIGURATION
func (c *ServiceConfiguration) loadPostgresStorageConfig(name string) (cfg *PostgresConfiguration, err error) {
	postgresConfig := &PostgresConfiguration{}

	name = strings.ToUpper(name)

	postgresConfig.Hostname = os.Getenv(fmt.Sprintf("NS_DB_PG_%s_HOSTNAME", name))
	postgresConfig.Port = os.Getenv(fmt.Sprintf("NS_DB_PG_%s_PORT", name))
	postgresConfig.SSLMode = os.Getenv(fmt.Sprintf("NS_DB_PG_%s_SSLMODE", name))
	postgresConfig.Username = os.Getenv(fmt.Sprintf("NS_DB_PG_%s_USERNAME", name))
	postgresConfig.Password = os.Getenv(fmt.Sprintf("NS_DB_PG_%s_PASSWORD", name))
	postgresConfig.Database = os.Getenv(fmt.Sprintf("NS_DB_PG_%s_DATABASE", name))

	if postgresConfig.Port == "" {
		postgresConfig.Port = "5432"
	}

	if postgresConfig.Database == "" || postgresConfig.Hostname == "" || postgresConfig.Username == "" {
		return nil, fmt.Errorf("invalid env parameters for database '%s': %s", name, spew.Sdump(postgresConfig))
	}

	return postgresConfig, nil
}
