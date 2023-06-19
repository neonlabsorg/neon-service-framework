package configuration

// DATABASES
type StorageConfiguration struct {
	Postgres  PostgresConfigCollection
	Clichouse ClickhouseConfigCollection
}

func (c *ServiceConfiguration) loadStorageConfigurations(storageList *ConfigStorageList) (err error) {
	if storageList == nil {
		return nil
	}

	if err = c.loadPostgresStorageConfigs(storageList.Postgres); err != nil {
		return err
	}

	if err = c.loadClickhouseStorageConfigs(storageList.Clickhouse); err != nil {
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

func (c *ServiceConfiguration) loadClickhouseStorageConfigs(list []string) (err error) {
	for _, db := range list {
		clickhouseConfig, err := c.loadClickhouseStorageConfig(db)
		if err != nil {
			return err
		}
		c.Storage.Clichouse.Add(db, clickhouseConfig)
	}

	return nil
}
