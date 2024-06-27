package middleware

import (
	"net/http"
	"rr-backend/internal/database"

	"github.com/labstack/echo/v4"
)

func ArtistMiddleware(dbService database.ScyllaService) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			userID := c.Get("userID").(string)

			user, err := dbService.GetUserByID(userID)
			if err != nil {
				return echo.NewHTTPError(http.StatusInternalServerError, "Failed to get user")
			}

			if user == nil || (user.Role != "artist" && user.Role != "admin") {
				return echo.NewHTTPError(http.StatusForbidden, "Unauthorized. Only artists or admins are allowed.")
			}

			return next(c)
		}
	}
}
