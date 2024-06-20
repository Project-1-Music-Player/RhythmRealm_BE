package handlers

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"rr-backend/internal/database"
	"rr-backend/internal/helper"

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

		parsedReleaseDate, err := time.Parse("2006-01-02", releaseDate)
		if err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, "Invalid release date format")
		}

		err = dbService.InsertSong(songID, title, userID, album, parsedReleaseDate, genre, songObjectName, thumbnailObjectName)
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, "Failed to save song metadata")
		}

		return c.JSON(http.StatusOK, echo.Map{
			"message": "Music uploaded successfully",
		})
	}
}

func RemoveSongHandler(dbService database.ScyllaService, minioService database.MinIOService) echo.HandlerFunc {
	return func(c echo.Context) error {
		songID := c.Param("song_id")
		objectName, err := dbService.GetObjectNameBySongID(songID)
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, "Failed to get song")
		}

		songUUID, err := gocql.ParseUUID(songID)
		if err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, "Invalid song ID")
		}

		err = dbService.RemoveSong(songUUID)
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, "Failed to remove song")
		}

		err = minioService.RemoveObject("music", objectName)
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, "Failed to remove song from storage")
		}

		return c.JSON(http.StatusOK, echo.Map{
			"message": "Song removed successfully",
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

func GetAllSongs(dbService database.ScyllaService) echo.HandlerFunc {
	return func(c echo.Context) error {
		songs, err := dbService.GetAllSongs()
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, "Failed to get songs")
		}

		return c.JSON(http.StatusOK, songs)
	}

}
func GetSongThumbnail(dbService database.ScyllaService, minioService database.MinIOService) echo.HandlerFunc {
	return func(c echo.Context) error {
		songID := c.Param("song_id")
		thumbnailName, err := dbService.GetSongThumbnailBySongID(songID)
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, "Failed to get thumbnail")
		}
		object, err := minioService.GetObject("music", thumbnailName)
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, "Failed to get thumbnail from storage")
		}
		defer object.Close()
		c.Stream(http.StatusOK, "image/jpeg", object)
		return nil
	}
}
func StreamMusic(dbService database.ScyllaService, minioService database.MinIOService) echo.HandlerFunc {
	return func(c echo.Context) error {
		songID := c.Param("song_id")
		objectName, err := dbService.GetObjectNameBySongID(songID)
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, "Failed to get song")
		}

		object, err := minioService.GetObject("music", objectName)
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, "Failed to get song from storage")
		}
		defer object.Close()

		helper.ServeContent(c.Response().Writer, c.Request(), objectName, time.Now(), object)
		go func() {
			if err := dbService.IncrementPlayCount(songID); err != nil {
				fmt.Printf("Error incrementing play count: %v\n", err)
				// TODO: Implement proper logging or error tracking mechanism
			}
		}()

		return nil
	}
}
func SearchSongs(dbService database.ScyllaService) echo.HandlerFunc {
	return func(c echo.Context) error {
		searchQuery := c.QueryParam("q")
		pageStr := c.QueryParam("page")
		limitStr := c.QueryParam("limit")

		page, _ := strconv.Atoi(pageStr)
		if page < 1 {
			page = 1
		}
		limit, _ := strconv.Atoi(limitStr)
		if limit <= 0 {
			limit = 10 // 10 items/page
		}
		offset := (page - 1) * limit

		songs, err := dbService.SearchSongs(searchQuery, limit, offset)
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, "Failed to search songs")
		}

		return c.JSON(http.StatusOK, songs)
	}
}

func LikeSongHandler(dbService database.ScyllaService) echo.HandlerFunc {
	return func(c echo.Context) error {
		songID := c.Param("song_id")
		userID := c.Get("userID").(string)

		songUUID, err := gocql.ParseUUID(songID)
		if err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, "Invalid song ID")
		}

		err = dbService.LikeSong(userID, songUUID)
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, "Failed to like song")
		}

		return c.JSON(http.StatusOK, echo.Map{
			"message": "Song liked successfully",
		})
	}
}

func UnlikeSongHandler(dbService database.ScyllaService) echo.HandlerFunc {
	return func(c echo.Context) error {
		songID := c.Param("song_id")
		userID := c.Get("userID").(string)

		songUUID, err := gocql.ParseUUID(songID)
		if err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, "Invalid song ID")
		}

		err = dbService.UnlikeSong(userID, songUUID)
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, "Failed to unlike song")
		}

		return c.JSON(http.StatusOK, echo.Map{
			"message": "Song unliked successfully",
		})
	}
}
