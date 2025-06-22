package configuration

import "github.com/neonlabsorg/neon-service-framework/pkg/tools/collections"

// DATABASES
type StorageConfiguration struct {
	Postgres  collections.BasicMapCollection[*PostgresConfiguration]
	Clichouse collections.BasicMapCollection[*ClickhouseConfiguration]
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
		c.Storage.Postgres.Set(db, postgresConfig)
	}

	return nil
}

func (c *ServiceConfiguration) loadClickhouseStorageConfigs(list []string) (err error) {
	for _, db := range list {
		clickhouseConfig, err := c.loadClickhouseStorageConfig(db)
		if err != nil {
			return err
		}
		c.Storage.Clichouse.Set(db, clickhouseConfig)
	}

	return nil
}
