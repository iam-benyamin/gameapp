package middleware

import (
	"gameapp/param"
	"gameapp/pkg/claim"
	"gameapp/pkg/errmsg"
	"gameapp/pkg/timestamp"
	"gameapp/service/presenceservice"
	"github.com/labstack/echo/v4"
	"net/http"
)

func UpsertPresence(service presenceservice.Service) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) (err error) {
			claims := claim.GetClaimsFromEchoContext(c)

			_, err = service.UpsertPresence(c.Request().Context(), param.UpsertPresenceRequest{
				UserID:    claims.UserID,
				Timestamp: timestamp.Now(),
			})

			if err != nil {
				// TODO: log unexpected error
				// we can just log the error and go to the next step
				return c.JSON(http.StatusInternalServerError, echo.Map{
					"message": errmsg.ErrorSomeThingWentWrong,
				})
			}

			return next(c)
		}
	}
}
