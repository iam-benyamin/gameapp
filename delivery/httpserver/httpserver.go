package httpserver

import (
	"fmt"
	"gameapp/config"
	"gameapp/service/authservice"
	"gameapp/service/userservice"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

type Server struct {
	config  config.Config
	authSvc authservice.Service
	userSvc userservice.Service
}

func New(config config.Config, authSVC authservice.Service, userSVC userservice.Service) Server {
	return Server{
		config:  config,
		authSvc: authSVC,
		userSvc: userSVC,
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
	e.POST("/users/register", s.UserRegister)

	// Start server
	e.Logger.Fatal(e.Start(fmt.Sprintf(":%d", s.config.HTTPServer.Port)))
}
