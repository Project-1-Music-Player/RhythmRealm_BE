package handlers

import (
	"net/http"
	"rr-backend/internal/database"
	"rr-backend/internal/models"

	"github.com/labstack/echo/v4"
)

func UpsertUserHandler(scyllaService database.ScyllaService) echo.HandlerFunc {
	return func(c echo.Context) error {
		// Get the verified UID from the JWT middleware
		userID := c.Get("userID").(string)

		// Bind the rest of the user data from the request body
		user := new(models.User)
		if err := c.Bind(user); err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, err.Error())
		}
		user.UserID = userID // Set the UID from the verified token

		// Upsert the user in the database
		err := scyllaService.UpsertUser(user.UserID, user.Username, user.Email, user.Role)
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, "Failed to upsert user")
		}

		return c.JSON(http.StatusOK, echo.Map{
			"message": "User upserted successfully",
		})
	}
}

func PromoteListenerToArtistHandler(scyllaService database.ScyllaService) echo.HandlerFunc {
	return func(c echo.Context) error {
		listenerID := c.Get("userID").(string)

		err := scyllaService.UpdateUserRole(listenerID, "artist")
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, "Failed to promote listener to artist")
		}

		return c.JSON(http.StatusOK, echo.Map{
			"message": "Listener promoted to artist successfully",
		})
	}
}
