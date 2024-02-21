package matchinghandler

import (
	"gameapp/service/authservice"
	"gameapp/service/matchingservice"
	"gameapp/service/presenceservice"
	"gameapp/validator/matchingvalidator"
)

type Handler struct {
	authConfig        authservice.Config
	authSvc           authservice.Service
	matchingSvc       matchingservice.Service
	matchingValidator matchingvalidator.Validator
	presenceSvc       presenceservice.Service
}

func New(authConfig authservice.Config, authSVC authservice.Service, matchingSvc matchingservice.Service, matchingValidator matchingvalidator.Validator, presenceSvc presenceservice.Service) Handler {
	return Handler{
		authConfig:        authConfig,
		authSvc:           authSVC,
		matchingSvc:       matchingSvc,
		matchingValidator: matchingValidator,
		presenceSvc:       presenceSvc,
	}
}
