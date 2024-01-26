package httpserver

import (
	"gameapp/pkg/httpmsg"
	"github.com/labstack/echo/v4"
	"net/http"

	"gameapp/dto"
)

func (s Server) userProfile(c echo.Context) error {

	authToken := c.Request().Header.Get("Authorization")
	claims, err := s.authSvc.ParseToken(authToken)
	if err != nil {
		return echo.NewHTTPError(http.StatusUnauthorized, err.Error())
	}

	resp, err := s.userSvc.Profile(dto.ProfileRequest{UserID: claims.UserID})
	if err != nil {
		msg, code := httpmsg.Error(err)
		return echo.NewHTTPError(code, msg)
	}

	return c.JSON(http.StatusOK, resp)
}

func (s Server) userLogin(c echo.Context) error {
	var lReq dto.LoginRequest
	if err := c.Bind(&lReq); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}

	resp, err := s.userSvc.Login(lReq)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	return c.JSON(http.StatusOK, resp)
}

func (s Server) userRegister(c echo.Context) error {
	var uReq dto.RegisterRequest
	if err := c.Bind(&uReq); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}

	if fieldsErrors, err := s.userValidator.ValidateRegisterRequest(uReq); err != nil {
		msg, code := httpmsg.Error(err)
		return c.JSON(code, echo.Map{
			"message": msg,
			"errors":  fieldsErrors,
		})
		//return echo.NewHTTPError(code, msg, fieldsErrors)
	}

	resp, err := s.userSvc.Register(uReq)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	return c.JSON(http.StatusCreated, resp)
}
