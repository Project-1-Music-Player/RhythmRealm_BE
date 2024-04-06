package handlers

import (
	"fmt"
	"net/http"
	"time"

	"rr-backend/internal/database"

	"github.com/gocql/gocql"
	"github.com/labstack/echo/v4"
)

func UploadMusicHandler(dbService database.ScyllaService, minioService database.MinIOService) echo.HandlerFunc {
	return func(c echo.Context) error {
		// Verify the JWT and get the user ID
		userID := c.Get("userID").(string)

		// Read form fields
		title := c.FormValue("title")
		album := c.FormValue("album")
		releaseDate := c.FormValue("releaseDate")
		genre := c.FormValue("genre")

		// Read form files
		songFile, err := c.FormFile("song")
		if err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, "Song file is required")
		}
		thumbnailFile, err := c.FormFile("thumbnail")
		if err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, "Thumbnail file is required")
		}

		// Generate UUID for song
		songID := gocql.TimeUUID()

		// Upload song file to MinIO
		songSrc, err := songFile.Open()
		if err != nil {
			return err
		}
		defer songSrc.Close()

		songObjectName := fmt.Sprintf("songs/%s/%s", userID, songFile.Filename)
		_, err = minioService.UploadObject("music", songObjectName, songSrc, songFile.Size, songFile.Header.Get("Content-Type"))
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, "Failed to upload song")
		}

		// Upload thumbnail file to MinIO
		thumbnailSrc, err := thumbnailFile.Open()
		if err != nil {
			return err
		}
		defer thumbnailSrc.Close()

		thumbnailObjectName := fmt.Sprintf("thumbnails/%s/%s", userID, thumbnailFile.Filename)
		_, err = minioService.UploadObject("music", thumbnailObjectName, thumbnailSrc, thumbnailFile.Size, thumbnailFile.Header.Get("Content-Type"))
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, "Failed to upload thumbnail")
		}

		// Parse release date
		parsedReleaseDate, err := time.Parse("2006-01-02", releaseDate)
		if err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, "Invalid release date format")
		}

		// Write metadata to database
		err = dbService.InsertSong(songID, title, userID, album, parsedReleaseDate, genre, songObjectName, thumbnailObjectName)
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, "Failed to save song metadata")
		}

		return c.JSON(http.StatusOK, echo.Map{
			"message": "Music uploaded successfully",
		})
	}
}

func GetSongsByUser(dbService database.ScyllaService) echo.HandlerFunc {
	return func(c echo.Context) error {
		userID := c.Get("userID").(string) // Get the userID from JWT middleware

		songs, err := dbService.GetSongsByUserID(userID)
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, "Failed to get songs")
		}

		return c.JSON(http.StatusOK, songs)
	}
}
