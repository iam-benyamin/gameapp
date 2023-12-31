package httpserver

import (
	"fmt"
	"gameapp/config"
	"gameapp/service/authservice"
	"gameapp/service/userservice"
	"gameapp/validator/uservalidator"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

type Server struct {
	config        config.Config
	authSvc       authservice.Service
	userSvc       userservice.Service
	userValidator uservalidator.Validator
}

func New(config config.Config, authSVC authservice.Service, userSVC userservice.Service, userValidator uservalidator.Validator) Server {
	return Server{
		config:        config,
		authSvc:       authSVC,
		userSvc:       userSVC,
		userValidator: userValidator,
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

	userGroup := e.Group("/users")
	userGroup.POST("/register", s.userRegister)
	userGroup.POST("/login", s.userLogin)
	userGroup.GET("/profile", s.userProfile)

	// Start server
	e.Logger.Fatal(e.Start(fmt.Sprintf(":%d", s.config.HTTPServer.Port)))
}
