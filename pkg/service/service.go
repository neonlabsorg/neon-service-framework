package service

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"strings"
	"sync"
	"syscall"

	"github.com/gagliardetto/solana-go/rpc"
	"github.com/neonlabsorg/neon-service-framework/pkg/alerts"
	"github.com/neonlabsorg/neon-service-framework/pkg/api"
	"github.com/neonlabsorg/neon-service-framework/pkg/env"
	"github.com/neonlabsorg/neon-service-framework/pkg/errors"
	"github.com/neonlabsorg/neon-service-framework/pkg/logger"
	"github.com/neonlabsorg/neon-service-framework/pkg/service/configuration"
	"github.com/urfave/cli/v2"
	"google.golang.org/grpc"
)

var Version string

type Service struct {
	env              string
	name             string
	version          string
	cfg              *configuration.ServiceConfiguration
	ctx              context.Context
	cliApp           *cli.App
	cliContext       *cli.Context
	loggerManager    *LoggerManager
	databaseManager  *DatabaseManager
	solanaRpcClient  *rpc.Client
	grpcServer       *GRPCServer
	apiServerManager *ApiServerManager
	alertDispatcher  *alerts.Dispatcher
	handlers         []func(service *Service)
}

func CreateService(
	config *configuration.Config,
) *Service {
	configuration, err := configuration.NewServiceConfiguration(config)
	if err != nil {
		panic(err)
	}

	env := env.Get("NS_ENV")
	if env == "" {
		env = "development"
	}

	if Version == "v." {
		Version = "v0.0.1"
	}

	s := &Service{
		env:     env,
		cfg:     configuration,
		name:    configuration.Name,
		version: Version,
	}

	s.initContext()
	s.initCliApp(configuration.IsConsoleApp)
	s.initLoggerManager(configuration.Logger)
	s.initSolana()
	s.initDatabases(configuration.Storage)
	s.initAlertDispatcher(configuration.Alerts)

	if configuration.UseGRPCServer {
		s.initGRPCServer(configuration.GRPCServer)
	}

	if len(configuration.ApiServers.Servers) > 0 {
		s.initApiServers(configuration.ApiServers)
	}

	if !configuration.IsConsoleApp {
		s.initMetrics(configuration.MetricsServer)
	}

	return s
}

func (s *Service) Run() {
	err := s.cliApp.Run(os.Args)
	if err != nil {
		panic(err.Error())
	}
}

func (s *Service) run(cliContext *cli.Context) (err error) {
	s.cliContext = cliContext
	s.loggerManager.GetLogger().Info().Msgf("Service %s version %s started", s.name, s.version)

	var wg sync.WaitGroup
	wg.Add(len(s.handlers))

	for _, handler := range s.handlers {
		go func(h func(s *Service), wGroup *sync.WaitGroup) {
			defer wGroup.Done()
			h(s)
		}(handler, &wg)
	}

	<-s.ctx.Done()
	wg.Wait()

	s.loggerManager.GetLogger().Info().Msgf("Service %s has been stopped", s.name)

	return
}

func (s *Service) initGRPCServer(cfg *configuration.GRPCServerConfiguration) {
	s.grpcServer = NewGRPCServer(cfg.ListenAddr)
}

func (s *Service) initApiServers(cfg *configuration.ApiServersConfiguration) {
	extender := api.NewDefaultApiContextExtender(api.NewValidator(), s.GetLogger())

	s.apiServerManager = NewApiServerManager(
		cfg,
		s.GetContext(),
		s.GetLogger(),
		extender,
	)

	err := s.apiServerManager.Init()
	if err != nil {
		s.GetLogger().Error().Err(err).Msg("can't initialize api server manager")
		panic(err)
	}
}

func (s *Service) initAlertDispatcher(cfg *configuration.AlertsConfiguration) {
	alertsContext := alerts.NewContext("project", s.name, "")

	var mainAdapter alerts.Adapter
	var reserveAdapter alerts.Adapter

	switch cfg.MainAdapter {
	case "prometheus":
		if cfg.Prometheus.URL == "" {
			break
		}
		mainAdapter = alerts.NewPrometheusAdapter(
			cfg.Prometheus,
			alertsContext,
			s.GetLogger(),
		)
	case "console":
		mainAdapter = alerts.NewConsoleAdapter(s.GetLogger())
	case "":
		mainAdapter = nil
	default:
		err := ErrUnregisteredAlertAdapter
		s.GetLogger().Error().Err(err).Msg("can't initialize main alerts adapter")
		panic(err)
	}

	switch cfg.ReserveAdapter {
	case "prometheus":
		reserveAdapter = alerts.NewPrometheusAdapter(
			cfg.Prometheus,
			alertsContext,
			s.GetLogger(),
		)
	case "console":
		reserveAdapter = alerts.NewConsoleAdapter(s.GetLogger())
	case "":
		reserveAdapter = nil
	default:
		err := ErrUnregisteredAlertAdapter
		s.GetLogger().Error().Err(err).Msg("can't initialize reserve alerts adapter")
		panic(err)
	}

	if mainAdapter == nil {
		mainAdapter = reserveAdapter
		reserveAdapter = nil
	}

	if mainAdapter == nil {
		mainAdapter = alerts.NewConsoleAdapter(s.GetLogger())
	}

	s.alertDispatcher = alerts.NewAlertDispatcher(mainAdapter, s.GetLogger())

	if reserveAdapter != nil {
		s.alertDispatcher.UseReservedAdapter(reserveAdapter)
	}
}

func (s *Service) initDatabases(cfg *configuration.StorageConfiguration) {
	var err error
	s.databaseManager, err = NewDatabaseManager(s.ctx, cfg, s.GetLogger())
	if err != nil {
		s.GetLogger().Error().Err(err).Msgf("error on init databases")
		panic(err)
	}
}

func (s *Service) initSolana() {
	solanaURL := env.Get("NS_SOLANA_URL")
	s.solanaRpcClient = rpc.New(solanaURL)
}

func (s *Service) initContext() {
	ctx, cancel := context.WithCancel(context.Background())
	sigquit := make(chan os.Signal, 1)
	signal.Ignore(syscall.SIGHUP, syscall.SIGPIPE)
	signal.Notify(sigquit, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-sigquit
		cancel()
	}()

	s.ctx = ctx
}

func (s *Service) initLoggerManager(cfg *configuration.LoggerConfiguration) {
	if cfg.Level == "" {
		if s.env == "development" {
			cfg.Level = "debug"
		} else {
			cfg.Level = "info"
		}
	}

	var log logger.Logger
	var err error

	if cfg.UseFile {
		log, err = logger.NewLogger(s.name, logger.LogSettings{
			Level: strings.ToLower(cfg.Level),
			Path:  strings.ToLower(cfg.FilePath),
		})

		if err != nil {
			panic(err)
		}
	} else {
		logger.InitDefaultLogger()
		log = logger.Get()
	}

	logger.SetDefaultLogger(log)

	s.loggerManager = NewLoggerManager(log)
}

func (s *Service) initCliApp(isConsoleApp bool) {
	s.cliApp = cli.NewApp()
	s.cliApp.Name = s.name
	s.cliApp.Version = s.version

	if !isConsoleApp {
		s.cliApp.Action = s.run
	}
}

func (s *Service) initMetrics(cfg *configuration.MetricsServerConfiguration) {
	if !cfg.Enable || cfg.ListenAddress == "" || cfg.ListenPort == 0 || cfg.Interval == 0 {
		s.GetLogger().Info().Msg("Metrics server inicialization has been skipped")
		return
	}

	metricsServer := NewMetricsServer(
		s.GetContext(),
		cfg.ServiceName,
		cfg.Interval,
		fmt.Sprintf("%s:%d", cfg.ListenAddress, cfg.ListenPort),
	)

	if err := metricsServer.Init(); err != nil {
		s.GetLogger().Error().Err(err).Msg("can't initialize metrics")
		panic(err)
	}

	go func() {
		if err := metricsServer.RunServer(); err != nil {
			s.GetLogger().Error().Err(err).Msg("can't start metrics server")
			panic(err)
		}
	}()
}

func (s *Service) ModifyCliApp(handler func(cliApp *cli.App)) {
	handler(s.cliApp)
}

func (s *Service) AddHandler(handler func(service *Service)) {
	s.handlers = append(s.handlers, handler)
}

func (s *Service) GetName() string {
	return s.name
}

func (s *Service) GetEnvironment() string {
	return s.env
}

func (s *Service) GetContext() context.Context {
	return s.ctx
}

func (s *Service) GetLogger() logger.Logger {
	return s.loggerManager.GetLogger()
}

func (s *Service) GetSolanaRpcClient() *rpc.Client {
	return s.solanaRpcClient
}

func (s *Service) RegisterGRPCService(svc *grpc.ServiceDesc, srv interface{}) {
	s.grpcServer.RegisterService(svc, srv)
}

func (s *Service) GetDatabaseManager() *DatabaseManager {
	return s.databaseManager
}

func (s *Service) GetAlertsDispatcher() *alerts.Dispatcher {
	return s.alertDispatcher
}

func (s *Service) GetApiServerManager() *ApiServerManager {
	return s.apiServerManager
}

func (s *Service) RunGRPCServer() (err error) {
	if s.grpcServer == nil {
		s.GetLogger().Error().Msg("the grpc server is not initialized")
		return errors.Critical.New("the grpc server is not initialized")
	}

	if !s.cfg.UseGRPCServer {
		s.GetLogger().Error().Msg("the grpc server is not activated")
		return errors.Critical.New("the grpc server is not activated")
	}

	s.GetLogger().Info().Msg("GRPC Server is starting")
	if s.grpcServer.services.Len() == 0 {
		s.GetLogger().Error().Msg("grpc server is running with no services")
	}

	err = s.grpcServer.Run()
	if err != nil {
		s.GetLogger().Error().Err(err).Msgf("error on running grpc server")
		return err
	}

	s.loggerManager.GetLogger().Info().Msg("GRPC Server has been started")

	return nil
}
