package configuration

type Config struct {
	Name         string
	Storage      *ConfigStorageList
	IsConsoleApp bool
	IsUnitedApp  bool
}

type ConfigStorageList struct {
	Postgres []string
}
