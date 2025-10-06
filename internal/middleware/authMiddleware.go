package middleware

import (
	"crypto/rsa"
	"net/http"
	"strings"

	"github.com/FIAP-11SOAT/totem-de-pedidos/internal/helper"
	"github.com/labstack/echo/v4"
)

func JWTAuthMiddleware(keyMap map[string]*rsa.PublicKey) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			authHeader := c.Request().Header.Get("Authorization")
			if authHeader == "" {
				return echo.NewHTTPError(http.StatusUnauthorized, "Authorization header missing")
			}

			parts := strings.SplitN(authHeader, " ", 2)
			if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") {
				return echo.NewHTTPError(http.StatusUnauthorized, "Invalid Authorization header format")
			}

			tokenStr := parts[1]

			validator := helper.NewJWTValidator()
			token, err := validator.Validate(tokenStr, keyMap)
			if err != nil {
				return echo.NewHTTPError(http.StatusUnauthorized, "Invalid or expired token: "+err.Error())
			}

			c.Set("user_token", token)
			return next(c)
		}
	}
}
