package configuration

// SERVICE CONFIGURATION
type ServiceConfiguration struct {
	Name          string
	IsConsoleApp  bool
	IsUnitedApp   bool
	UseGRPCServer bool
	UseAPIServer  bool
	Logger        *LoggerConfiguration
	Storage       *StorageConfiguration
	MetricsServer *MetricsServerConfiguration
	GRPCServer    *GRPCServerConfiguration
	ApiServer     *ApiServerConfiguration
}

// INIT CONFIGURATION
func NewServiceConfiguration(cfg *Config) (serviceConfiguration *ServiceConfiguration, err error) {
	serviceConfiguration = &ServiceConfiguration{
		Name:          cfg.Name,
		IsConsoleApp:  cfg.IsConsoleApp,
		IsUnitedApp:   cfg.IsUnitedApp,
		UseGRPCServer: cfg.UseGRPCServer,
		UseAPIServer:  cfg.UseAPIServer,
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

	if err = serviceConfiguration.loadApiServerConfiguration(); err != nil {
		return nil, err
	}

	return serviceConfiguration, nil
}
