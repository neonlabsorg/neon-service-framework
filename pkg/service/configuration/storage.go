package configuration

type PostgresConfigCollection map[string]*PostgresConfiguration

func (c PostgresConfigCollection) init() {
	if c == nil {
		c = make(map[string]*PostgresConfiguration)
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

// DATABASES
type StorageConfiguration struct {
	Postgres PostgresConfigCollection
}

func (c *ServiceConfiguration) loadStorageConfigurations(storageList *ConfigStorageList) (err error) {
	if storageList == nil {
		return nil
	}

	if err = c.loadPostgresStorageConfigs(storageList.Postgres); err != nil {
		return err
	}

	return nil
}

func (c *ServiceConfiguration) loadPostgresStorageConfigs(list []string) (err error) {
	for _, db := range list {
		postgresConfig, err := c.loadPostgresStorageConfig(db)
		if err != nil {
			return err
		}
		c.Storage.Postgres.Add(db, postgresConfig)
	}

	return nil
}
