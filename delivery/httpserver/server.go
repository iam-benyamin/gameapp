package httpserver

import (
	"fmt"
	"gameapp/config"
	"gameapp/delivery/httpserver/backofficeuserhandler"
	"gameapp/delivery/httpserver/matchinghandler"
	"gameapp/delivery/httpserver/userhandler"
	"gameapp/logger"
	"gameapp/service/authorizationservice"
	"gameapp/service/authservice"
	"gameapp/service/backofficeuserservice"
	"gameapp/service/matchingservice"
	"gameapp/service/presenceservice"
	"gameapp/service/userservice"
	"gameapp/validator/matchingvalidator"
	"gameapp/validator/uservalidator"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"go.uber.org/zap"
)

type Server struct {
	config                config.Config
	userHandler           userhandler.Handler
	backofficeUserHandler backofficeuserhandler.Handler
	matchingHandler       matchinghandler.Handler
	Router                *echo.Echo
}

func New(config config.Config, authSVC authservice.Service,
	userSVC userservice.Service, userValidator uservalidator.Validator,
	backofficeUserHandlerSvc backofficeuserservice.Service, authorizationSvc authorizationservice.Service,
	matchingSvc matchingservice.Service, matchingValidator matchingvalidator.Validator, presenceSvc presenceservice.Service,
) Server {

	return Server{
		Router:                echo.New(),
		config:                config,
		userHandler:           userhandler.New(config.Auth, authSVC, userSVC, userValidator, presenceSvc),
		backofficeUserHandler: backofficeuserhandler.New(config.Auth, authSVC, backofficeUserHandlerSvc, authorizationSvc),
		matchingHandler:       matchinghandler.New(config.Auth, authSVC, matchingSvc, matchingValidator, presenceSvc),
	}
}

func (s Server) Serve() {
	// Middleware
	//s.Router.Use(middleware.Logger())
	// use my zap logger not echo logger
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
			errMsg := ""
			if v.Error != nil {
				errMsg = v.Error.Error()
			}

			logger.Logger.Named("http-server").Info("request",
				zap.String("request_id", v.RequestID),
				zap.String("host", v.Host),
				zap.String("content-length", v.ContentLength),
				zap.String("protocol", v.Protocol),
				zap.String("method", v.Method),
				zap.Duration("latency", v.Latency),
				zap.String("error", errMsg),
				zap.String("remote_ip", v.RemoteIP),
				zap.Int64("response_size", v.ResponseSize),
				zap.String("uri", v.URI),
				zap.Int("status", v.Status),
			)

			return nil
		},
	}))

	s.Router.Use(middleware.Recover())

	// Routes
	s.Router.GET("/heath-check/", s.healthCheck)

	s.userHandler.SetUserRouts(s.Router)
	s.backofficeUserHandler.SetUserRouts(s.Router)
	s.matchingHandler.SetSetMatchingRouts(s.Router)

	// Start server
	address := fmt.Sprintf(":%d", s.config.HTTPServer.Port)
	fmt.Printf("start echo server on %s\n", address)
	if err := s.Router.Start(address); err != nil {
		fmt.Println("router start error ", err)
	}
}
