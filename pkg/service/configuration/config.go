package configuration

type Config struct {
	Name          string
	Storage       *ConfigStorageList
	IsConsoleApp  bool
	IsUnitedApp   bool
	UseGRPCServer bool
	UseAPIServer  bool
}

type ConfigStorageList struct {
	Postgres   []string
	Clickhouse []string
}
