package middleware

import (
	"gameapp/pkg/constant"
	"gameapp/service/authservice"
	"github.com/labstack/echo/v4"

	mw "github.com/labstack/echo-jwt/v4"
)

// Closure: a function which return function or higher order function
func Auth(service authservice.Service, config authservice.Config) echo.MiddlewareFunc {
	return mw.WithConfig(mw.Config{
		ContextKey:    constant.AuthMiddlewareContextKey,
		SigningKey:    []byte(config.SignKey),
		SigningMethod: "HS256",
		ParseTokenFunc: func(c echo.Context, auth string) (interface{}, error) {
			claims, err := service.ParseToken(auth)
			if err != nil {
				return nil, err
			}

			return claims, nil
		},
	})

}
