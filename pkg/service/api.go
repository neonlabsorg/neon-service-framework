package service

import (
	"context"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/neonlabsorg/neon-service-framework/pkg/api"
	"github.com/neonlabsorg/neon-service-framework/pkg/echo/binder"
	"github.com/neonlabsorg/neon-service-framework/pkg/logger"
	"github.com/neonlabsorg/neon-service-framework/pkg/service/configuration"
)

type ApiServer struct {
	server   *echo.Echo
	ctx      context.Context
	cfg      *configuration.ApiServerConfiguration
	extender api.ApiContextExtender
	logger   logger.Logger
}

func NewApiServer(
	ctx context.Context,
	cfg *configuration.ApiServerConfiguration,
	extender api.ApiContextExtender,
	log logger.Logger,
) *ApiServer {
	s := &ApiServer{
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
	s.server.Use(middlware)
}

func (s *ApiServer) SetCustomExtender(extender api.ApiContextExtender) {
	s.extender = extender
	s.server.Use(s.extender.ExtendDefaultApiContext)
}

func (s *ApiServer) registerExtender() {
	s.server.Use(s.extender.ExtendDefaultApiContext)
}

func (s *ApiServer) Run() (err error) {

	s.registerExtender()

	go func() {
		if err := s.server.Start(s.cfg.ListenAddr); err != nil {
			s.server.Logger.Info(err)
			s.logger.Error().Err(err).Msg("error on start api server")
		}
	}()

	<-s.ctx.Done()

	if err := s.server.Shutdown(s.ctx); err != nil {
		s.logger.Error().Err(err).Msg("error on shutdown echo")
		s.server.Logger.Fatal(err)
	}

	return nil
}
