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
}

func New(config config.Config, authSVC authservice.Service,
	userSVC userservice.Service, userValidator uservalidator.Validator,
	backofficeUserHandlerSvc backofficeuserservice.Service, authorizationSvc authorizationservice.Service,
	matchingSvc matchingservice.Service, matchingValidator matchingvalidator.Validator,
) Server {

	return Server{
		config:                config,
		userHandler:           userhandler.New(config.Auth, authSVC, userSVC, userValidator),
		backofficeUserHandler: backofficeuserhandler.New(config.Auth, authSVC, backofficeUserHandlerSvc, authorizationSvc),
		matchingHandler:       matchinghandler.New(config.Auth, authSVC, matchingSvc, matchingValidator),
	}
}

func (s Server) Serve() {
	// Echo instance
	e := echo.New()

	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	// Routes
	e.GET("/heath-check/", s.healthCheck)

	s.userHandler.SetUserRouts(e)
	s.backofficeUserHandler.SetUserRouts(e)
	s.matchingHandler.SetSetMatchingRouts(e)

	// Start server
	address := fmt.Sprintf(":%d", s.config.HTTPServer.Port)
	fmt.Printf("start echo server on %s\n", address)
	e.Logger.Fatal(e.Start(address))
}
