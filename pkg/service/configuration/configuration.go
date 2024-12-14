package configuration

// SERVICE CONFIGURATION
type ServiceConfiguration struct {
	Name               string
	IsConsoleApp       bool
	IsUnitedApp        bool
	UseGRPCServer      bool
	UseUnitedAPIServer bool
	Logger             *LoggerConfiguration
	Storage            *StorageConfiguration
	MetricsServer      *MetricsServerConfiguration
	GRPCServer         *GRPCServerConfiguration
	ApiServers         *ApiServersConfiguration
}

// INIT CONFIGURATION
func NewServiceConfiguration(cfg *Config) (serviceConfiguration *ServiceConfiguration, err error) {
	if cfg.ApiServers == nil {
		cfg.ApiServers = &ApiServersConfig{}
	}

	serviceConfiguration = &ServiceConfiguration{
		Name:               cfg.Name,
		IsConsoleApp:       cfg.IsConsoleApp,
		IsUnitedApp:        cfg.IsUnitedApp,
		UseGRPCServer:      cfg.UseGRPCServer,
		UseUnitedAPIServer: cfg.UseUnitedAPIServer,
		Storage: &StorageConfiguration{
			Postgres: make(map[string]*PostgresConfiguration),
		},
	}

	if err = serviceConfiguration.loadLoggerConfiguration(); err != nil {
		return nil, err
	}

	if err = serviceConfiguration.loadStorageConfigurations(cfg.Storage); err != nil {
		return nil, err
	}

	if err = serviceConfiguration.loadMetricsServerConfiguration(cfg.Name); err != nil {
		return nil, err
	}

	if err = serviceConfiguration.loadGRPCServerConfiguration(); err != nil {
		return nil, err
	}

	if err = serviceConfiguration.loadApiServersConfiguration(cfg.UseUnitedAPIServer, cfg.ApiServers.Names); err != nil {
		return nil, err
	}

	return serviceConfiguration, nil
}
