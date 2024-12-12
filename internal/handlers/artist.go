package handlers

import (
    "net/http"
    "rr-backend/internal/database"
    "github.com/labstack/echo/v4"
)

func GetAllArtistsHandler(dbService database.ScyllaService) echo.HandlerFunc {
    return func(c echo.Context) error {
        artists, err := dbService.GetAllArtists()
        if err != nil {
            return echo.NewHTTPError(http.StatusInternalServerError, "Failed to get artists")
        }
        return c.JSON(http.StatusOK, artists)
    }
}

func GetArtistWithSongsHandler(dbService database.ScyllaService) echo.HandlerFunc {
    return func(c echo.Context) error {
        artistID := c.Param("artist_id")
        
        artistWithSongs, err := dbService.GetArtistWithSongs(artistID)
        if err != nil {
            return echo.NewHTTPError(http.StatusInternalServerError, "Failed to get artist data")
        }
        
        return c.JSON(http.StatusOK, artistWithSongs)
    }
}

func FollowArtistHandler(dbService database.ScyllaService) echo.HandlerFunc {
    return func(c echo.Context) error {
        artistID := c.Param("artist_id")
        followerID := c.Get("userID").(string)

        err := dbService.FollowArtist(artistID, followerID)
        if err != nil {
            return echo.NewHTTPError(http.StatusInternalServerError, "Failed to follow artist")
        }

        return c.JSON(http.StatusOK, echo.Map{
            "message": "Successfully followed artist",
        })
    }
}

func UnfollowArtistHandler(dbService database.ScyllaService) echo.HandlerFunc {
    return func(c echo.Context) error {
        artistID := c.Param("artist_id")
        followerID := c.Get("userID").(string)

        err := dbService.UnfollowArtist(artistID, followerID)
        if err != nil {
            return echo.NewHTTPError(http.StatusInternalServerError, "Failed to unfollow artist")
        }

        return c.JSON(http.StatusOK, echo.Map{
            "message": "Successfully unfollowed artist",
        })
    }
}

func GetFollowedArtistsHandler(dbService database.ScyllaService) echo.HandlerFunc {
    return func(c echo.Context) error {
        userID := c.Get("userID").(string)
        
        artists, err := dbService.GetFollowedArtists(userID)
        if err != nil {
            return echo.NewHTTPError(http.StatusInternalServerError, "Failed to get followed artists")
        }
        
        return c.JSON(http.StatusOK, artists)
    }
} 