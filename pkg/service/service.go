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
	"github.com/labstack/echo/v4"
	"github.com/neonlabsorg/neon-service-framework/pkg/api"
	"github.com/neonlabsorg/neon-service-framework/pkg/env"
	"github.com/neonlabsorg/neon-service-framework/pkg/logger"
	"github.com/neonlabsorg/neon-service-framework/pkg/service/configuration"
	"github.com/urfave/cli/v2"
	"google.golang.org/grpc"
)

var Version string

type Service struct {
	env             string
	name            string
	version         string
	cfg             *configuration.ServiceConfiguration
	ctx             context.Context
	cliApp          *cli.App
	cliContext      *cli.Context
	loggerManager   *LoggerManager
	databaseManager *DatabaseManager
	solanaRpcClient *rpc.Client
	grpcServer      *GRPCServer
	apiServer       *ApiServer
	handlers        []func(service *Service)
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

	if configuration.UseGRPCServer {
		s.initGRPCServer(configuration.GRPCServer)
	}

	if configuration.UseAPIServer {
		s.initApiServer(configuration.ApiServer)
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

	if s.cfg.UseGRPCServer {
		wg.Add(1)
		go func(s *Service, wGroup *sync.WaitGroup) {
			defer wGroup.Done()

			s.GetLogger().Info().Msg("GRPC Server is starting")
			if s.grpcServer.services.Len() == 0 {
				s.GetLogger().Error().Msg("grpc server is running with no services")
			}
			err = s.grpcServer.Run()
			if err != nil {
				s.GetLogger().Error().Err(err).Msgf("error on running grpc server")
			} else {
				s.loggerManager.GetLogger().Info().Msg("GRPC Server has been started")
			}
		}(s, &wg)
	}

	if s.cfg.UseAPIServer {
		wg.Add(1)
		go func(s *Service, wGroup *sync.WaitGroup) {
			defer wGroup.Done()

			s.loggerManager.GetLogger().Info().Msg("API Server is starting")
			err = s.apiServer.Run()
			if err != nil {
				s.GetLogger().Error().Err(err).Msgf("error on running api server")
			} else {
				s.loggerManager.GetLogger().Info().Msg("API Server has been started")
			}
		}(s, &wg)
	}

	<-s.ctx.Done()
	wg.Wait()

	s.loggerManager.GetLogger().Info().Msgf("Service %s has been stopped", s.name)

	return
}

func (s *Service) initGRPCServer(cfg *configuration.GRPCServerConfiguration) {
	s.grpcServer = NewGRPCServer(cfg.ListenAddr)
}

func (s *Service) initApiServer(cfg *configuration.ApiServerConfiguration) {
	extender := api.NewDefaultApiContextExtender(api.NewValidator(), s.GetLogger())

	s.apiServer = NewApiServer(
		s.ctx,
		cfg,
		extender,
		s.GetLogger(),
	)
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

func (s *Service) RegisterApiRoutes(handler func(server *echo.Echo) error) (err error) {
	return s.apiServer.RegisterRoutes(handler)
}

func (s *Service) SetCustomExtenderForApiServer(extender api.ApiContextExtender) {
	s.apiServer.SetCustomExtender(extender)
}

func (s *Service) UseMiddlewareForApiServer(middlware echo.MiddlewareFunc) {
	s.apiServer.UseMiddleware(middlware)
}

func (s *Service) GetDatabaseManager() *DatabaseManager {
	return s.databaseManager
}
