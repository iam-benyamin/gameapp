package userhandler

import (
	"gameapp/service/authservice"
	"gameapp/service/presenceservice"
	"gameapp/service/userservice"
	"gameapp/validator/uservalidator"
)

type Handler struct {
	authConfig    authservice.Config
	authSvc       authservice.Service
	userSvc       userservice.Service
	userValidator uservalidator.Validator
	presenceSvc   presenceservice.Service
}

func New(authConfig authservice.Config, authSVC authservice.Service, userSVC userservice.Service, userValidator uservalidator.Validator, presenceSvc presenceservice.Service) Handler {
	return Handler{
		authConfig:    authConfig,
		authSvc:       authSVC,
		userSvc:       userSVC,
		userValidator: userValidator,
		presenceSvc:   presenceSvc,
	}
}
