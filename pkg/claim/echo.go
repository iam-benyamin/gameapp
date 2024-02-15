package claim

import (
	"gameapp/config"
	"gameapp/service/authservice"
	"github.com/labstack/echo/v4"
)

func GetClaimsFromEchoContext(c echo.Context) *authservice.Claims {
	// let it crash and defencive programming are two vision for programming
	return c.Get(config.AuthMiddlewareContextKey).(*authservice.Claims)
}
