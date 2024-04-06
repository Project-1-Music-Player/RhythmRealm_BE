package middleware

import (
	"net/http"
	firebase "rr-backend/internal/auth"
	"strings"

	"github.com/labstack/echo/v4"
)

func JWTMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		authHeader := c.Request().Header.Get("Authorization")
		if authHeader == "" {
			return echo.NewHTTPError(http.StatusUnauthorized, "No ID token provided")
		}

		idToken := strings.TrimPrefix(authHeader, "Bearer ")
		token, err := firebase.AuthClient.VerifyIDToken(c.Request().Context(), idToken)
		if err != nil {
			return echo.NewHTTPError(http.StatusUnauthorized, "Invalid or expired ID token")
		}

		// Store UID for handler
		c.Set("userID", token.UID)

		return next(c)
	}
}
