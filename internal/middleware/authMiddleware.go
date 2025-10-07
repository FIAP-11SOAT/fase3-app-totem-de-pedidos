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

			c.Set("cognito_claims", token)
			return next(c)
		}
	}
}

func RoleMiddleware(allowedRoles ...string) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			claims := c.Get("cognito_claims").(*helper.CognitoClaims)
			if claims == nil {
				return echo.NewHTTPError(http.StatusUnauthorized, "Token não encontrado")
			}

			for _, userRole := range claims.CognitoGroups {
				for _, allowedRole := range allowedRoles {
					if userRole == allowedRole {
						return next(c)
					}
				}
			}

			return echo.NewHTTPError(http.StatusForbidden, "Acesso negado - role necessária: "+strings.Join(allowedRoles, ", "))
		}
	}
}

func RequireAdminRole() echo.MiddlewareFunc {
	return RoleMiddleware("admin", "employees")
}
