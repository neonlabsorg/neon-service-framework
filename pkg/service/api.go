package service

import (
	"context"
	"reflect"
	"runtime"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/neonlabsorg/neon-service-framework/pkg/api"
	"github.com/neonlabsorg/neon-service-framework/pkg/echo/binder"
	"github.com/neonlabsorg/neon-service-framework/pkg/logger"
	"github.com/neonlabsorg/neon-service-framework/pkg/service/configuration"
	"github.com/neonlabsorg/neon-service-framework/pkg/tools/collections"
)

type ApiServer struct {
	name        string
	server      *echo.Echo
	ctx         context.Context
	cfg         *configuration.ApiServerConfiguration
	middlewares collections.SafeMapCollection[echo.MiddlewareFunc]
	extender    api.ApiContextExtender
	logger      logger.Logger
}

func NewApiServer(
	name string,
	ctx context.Context,
	cfg *configuration.ApiServerConfiguration,
	extender api.ApiContextExtender,
	log logger.Logger,
) *ApiServer {
	s := &ApiServer{
		name:     name,
		ctx:      ctx,
		cfg:      cfg,
		logger:   log,
		extender: extender,
	}

	s.server = s.newEcho()

	return s
}

func (s *ApiServer) newEcho() *echo.Echo {
	e := echo.New()

	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.BodyLimit(s.cfg.BodyLimit))
	if s.cfg.UseCORS {
		e.Use(middleware.CORS())
	}

	e.HTTPErrorHandler = api.HttpErrorHandler

	e.Binder = new(binder.ModelBinder)

	return e
}

func (s *ApiServer) RegisterRoutes(handler func(server *echo.Echo) error) (err error) {
	err = handler(s.server)
	if err != nil {
		return
	}

	return nil
}

func (s *ApiServer) UseMiddleware(middlware echo.MiddlewareFunc) {
	name := runtime.FuncForPC(reflect.ValueOf(middlware).Pointer()).Name()
	s.middlewares.Set(name, middlware)
}

func (s *ApiServer) SetCustomExtender(extender api.ApiContextExtender) {
	s.extender = extender
}

func (s *ApiServer) registerExtender() {
	s.server.Use(s.extender.ExtendDefaultApiContext)
}

func (s *ApiServer) Run() (err error) {
	s.logger.Info().Msgf("the api server %s is starting", s.name)

	s.registerExtender()

	err = s.middlewares.Iter(func(middleware echo.MiddlewareFunc) error {
		s.server.Use(middleware)
		return nil
	})
	if err != nil {
		s.logger.Error().Err(err).Msg("error on init middlewares")
		return err
	}

	go func() {
		if err := s.server.Start(s.cfg.ListenAddr); err != nil {
			s.server.Logger.Info(err)
			s.logger.Error().Err(err).Msg("error on start api server")
		}
		s.logger.Info().Msgf("the api server %s has been started", s.name)
	}()

	<-s.ctx.Done()

	if err := s.server.Shutdown(s.ctx); err != nil {
		s.logger.Error().Err(err).Msg("error on shutdown echo")
		s.server.Logger.Fatal(err)
	}

	s.logger.Info().Msgf("the api server %s has been stopped", s.name)

	return nil
}

type UnitedApiServerClient struct {
	Name        string `json:"name"`
	ServiceName string `json:"service_name"`
	IsReady     bool   `json:"is_ready"`
}

type UnitedApiServer struct {
	clients collections.SafeMapCollection[*UnitedApiServerClient]
	server  *ApiServer
}

func NewUnitedApiServer(server *ApiServer) *UnitedApiServer {
	return &UnitedApiServer{
		clients: collections.NewSafeMapCollection[*UnitedApiServerClient](),
		server:  server,
	}
}

func (s *UnitedApiServer) AddClient(client *UnitedApiServerClient) (err error) {
	_, ok := s.clients.Get(client.Name)
	if ok {
		return ErrUnitedApiServerClientAlreadyExists
	}

	s.clients.Set(client.Name, client)

	return nil
}

func (s *UnitedApiServer) ReadyToStart(client *UnitedApiServerClient) (err error) {
	client.IsReady = true
	s.clients.Set(client.Name, client)

	err = s.tryToStart()
	if err != nil {
		return err
	}

	return nil
}

func (s *UnitedApiServer) tryToStart() (err error) {
	isReady := true
	err = s.clients.Iter(func(item *UnitedApiServerClient) error {
		if !item.IsReady {
			isReady = false
		}

		return nil
	})
	if err != nil {
		return err
	}

	if isReady {
		err = s.server.Run()
		if err != nil {
			return err
		}
	}

	return nil
}

func (s *UnitedApiServer) RegisterRoutes(handler func(server *echo.Echo) error) (err error) {
	return s.server.RegisterRoutes(handler)
}

func (s *UnitedApiServer) UseMiddleware(middlware echo.MiddlewareFunc) {
	s.server.UseMiddleware(middlware)
}

func (s *UnitedApiServer) SetCustomExtender(extender api.ApiContextExtender) {
	s.server.SetCustomExtender(extender)
}
