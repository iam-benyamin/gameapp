package httpserver

import (
	"fmt"
	"gameapp/config"
	"gameapp/delivery/httpserver/backofficeuserhandler"
	"gameapp/delivery/httpserver/matchinghandler"
	"gameapp/delivery/httpserver/userhandler"
	"gameapp/service/authorizationservice"
	"gameapp/service/authservice"
	"gameapp/service/backofficeuserservice"
	"gameapp/service/matchingservice"
	"gameapp/service/userservice"
	"gameapp/validator/matchingvalidator"
	"gameapp/validator/uservalidator"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
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
	matchingSvc matchingservice.Service, matchingValidator matchingvalidator.Validator,
) Server {

	return Server{
		Router:                echo.New(),
		config:                config,
		userHandler:           userhandler.New(config.Auth, authSVC, userSVC, userValidator),
		backofficeUserHandler: backofficeuserhandler.New(config.Auth, authSVC, backofficeUserHandlerSvc, authorizationSvc),
		matchingHandler:       matchinghandler.New(config.Auth, authSVC, matchingSvc, matchingValidator),
	}
}

func (s Server) Serve() {
	// Middleware
	s.Router.Use(middleware.Logger())
	s.Router.Use(middleware.Recover())

	// Routes
	s.Router.GET("/heath-check/", s.healthCheck)

	s.userHandler.SetUserRouts(s.Router)
	s.backofficeUserHandler.SetUserRouts(s.Router)
	s.matchingHandler.SetSetMatchingRouts(s.Router)

	// Start server
	address := fmt.Sprintf(":%d", s.config.HTTPServer.Port)
	fmt.Printf("start echo server on %s\n", address)
	s.Router.Logger.Fatal(s.Router.Start(address))
}
