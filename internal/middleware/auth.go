package middleware

import (
	"net/http"
	"strings"
	"university/internal/service"

	"github.com/labstack/echo/v4"
)

// AuthMiddleware validates JWT token from Authorization header
// Uses: Authorization: Bearer <token>
func AuthMiddleware(svc *service.Service) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			// Get Authorization header
			authHeader := c.Request().Header.Get("Authorization")
			if authHeader == "" {
				return c.JSON(http.StatusUnauthorized, map[string]string{
					"error": "missing authorization header",
				})
			}

			// Extract bearer token
			parts := strings.Split(authHeader, " ")
			if len(parts) != 2 || parts[0] != "Bearer" {
				return c.JSON(http.StatusUnauthorized, map[string]string{
					"error": "invalid authorization header format",
				})
			}

			tokenString := parts[1]

			// Validate token
			userID, err := svc.ValidateToken(tokenString)
			if err != nil {
				return c.JSON(http.StatusUnauthorized, map[string]string{
					"error": "invalid or expired token",
				})
			}

			// Store user ID in context for handler access
			c.Set("user_id", userID)

			return next(c)
		}
	}
}
