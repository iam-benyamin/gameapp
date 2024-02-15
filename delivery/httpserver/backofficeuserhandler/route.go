package backofficeuserhandler

import (
	"gameapp/delivery/httpserver/middleware"
	"gameapp/entity"
	"github.com/labstack/echo/v4"
)

func (h Handler) SetUserRouts(e *echo.Echo) {
	userGroup := e.Group("/backoffice/users")

	userGroup.GET("/", h.listUsers, middleware.Auth(h.authSvc, h.authConfig),
		middleware.AccessCheck(h.authorizationSvc, entity.UserListPermission))
}
