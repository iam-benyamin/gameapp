package httpserver

import (
	"gameapp/service/userservice"
	"github.com/labstack/echo/v4"
	"net/http"
)

func (s Server) UserRegister(c echo.Context) error {
	var uReq userservice.RegisterRequest
	if err := c.Bind(&uReq); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}

	resp, err := s.userSvc.Register(uReq)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	return c.JSON(http.StatusCreated, resp)
}
