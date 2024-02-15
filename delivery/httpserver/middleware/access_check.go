package middleware

import (
	"gameapp/entity"
	"gameapp/pkg/claim"
	"gameapp/pkg/errmsg"
	"gameapp/service/authorizationservice"
	"github.com/labstack/echo/v4"
	"net/http"
)

func AccessCheck(service authorizationservice.Service, permissions ...entity.PermissionTitle) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) (err error) {
			claims := claim.GetClaimsFromEchoContext(c)
			isAllowed, err := service.CheckAccess(claims.UserID, claims.Role, permissions...)
			if err != nil {
				return c.JSON(http.StatusForbidden, echo.Map{
					"message": errmsg.ErrorSomeThingWentWrong,
				})
			}

			if !isAllowed {
				return c.JSON(http.StatusForbidden, echo.Map{
					"message": errmsg.ErrorMsgAccessDenied,
				})
			}

			return nil
		}
	}
}
