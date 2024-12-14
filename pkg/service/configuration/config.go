package configuration

type Config struct {
	Name               string
	Storage            *ConfigStorageList
	IsConsoleApp       bool
	IsUnitedApp        bool
	UseGRPCServer      bool
	UseUnitedAPIServer bool
	ApiServers         *ApiServersConfig
}

type ConfigStorageList struct {
	Postgres   []string
	Clickhouse []string
}

type ApiServersConfig struct {
	Names []string
}
