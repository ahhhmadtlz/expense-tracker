package httpserver

import (
	"fmt"
	"log/slog"

	"github.com/ahhhmadtlz/expense-tracker/internal/config"
	"github.com/ahhhmadtlz/expense-tracker/internal/delivery/httpserver/userhandler"
	"github.com/ahhhmadtlz/expense-tracker/internal/domain/auth"
	userservice "github.com/ahhhmadtlz/expense-tracker/internal/domain/user/service"
	uservalidator "github.com/ahhhmadtlz/expense-tracker/internal/domain/user/validator"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

type Server struct {
	config          config.Config
	userHandler     userhandler.Handler
	logger *slog.Logger
	Router          *echo.Echo
}


func New(
	config config.Config,
	auth auth.Service,
	userSvc userservice.Service,
	userValidator uservalidator.Validator,
	logger        *slog.Logger,
)Server{
	return  Server{
		Router: echo.New(),
		config: config,
		logger:      logger,
		userHandler: userhandler.New(
			auth,userSvc,userValidator,config.Auth,logger,
		),
	}
}


func (s Server)Serve(){
	s.Router = echo.New()
	s.Router.Use(middleware.RequestID())
	s.Router.Use(middleware.RequestLoggerWithConfig(middleware.RequestLoggerConfig{
		LogURI:           true,
		LogStatus:        true,
		LogHost:          true,
		LogRemoteIP:      true,
		LogRequestID:     true,
		LogMethod:        true,
		LogContentLength: true,
		LogResponseSize:  true,
		LogLatency:       true,
		LogError:         true,
		LogProtocol:      true,
		LogValuesFunc: func(c echo.Context, v middleware.RequestLoggerValues) error {
			if v.Error == nil {
				s.logger.Info("request",
					"request_id", v.RequestID,
					"host", v.Host,
					"content_length", v.ContentLength,
					"protocol", v.Protocol,
					"method", v.Method,
					"latency", v.Latency,
					"remote_ip", v.RemoteIP,
					"response_size", v.ResponseSize,
					"uri", v.URI,
					"status", v.Status,
				)
			} else {
				s.logger.Error("request",
					"request_id", v.RequestID,
					"host", v.Host,
					"content_length", v.ContentLength,
					"protocol", v.Protocol,
					"method", v.Method,
					"latency", v.Latency,
					"error", v.Error.Error(),
					"remote_ip", v.RemoteIP,
					"response_size", v.ResponseSize,
					"uri", v.URI,
					"status", v.Status,
				)
			}
			return nil
		},
	}))

	s.Router.Use(middleware.Recover())

	s.Router.GET("/health-check", s.healthCheck)


	s.userHandler.SetRoutes(s.Router)


		// Start server
	address := fmt.Sprintf(":%d", s.config.HTTPServer.Port)
	fmt.Printf("start echo server on %s\n", address)
	if err := s.Router.Start(address); err != nil {
		fmt.Println("router start error", err)
	}

}