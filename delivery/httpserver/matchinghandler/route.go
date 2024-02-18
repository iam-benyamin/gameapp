package matchinghandler

import (
	"gameapp/delivery/httpserver/middleware"
	"github.com/labstack/echo/v4"
)

func (h Handler) SetSetMatchingRouts(e *echo.Echo) {
	userGroup := e.Group("/matching")

	userGroup.POST("/add-to-waiting-list", h.addToWaitingList, middleware.Auth(h.authSvc, h.authConfig))
}
