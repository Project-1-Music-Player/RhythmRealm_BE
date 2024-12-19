package handlers

import (
	"net/http"
	"rr-backend/internal/database"
	"rr-backend/internal/models"
	"time"

	"github.com/gocql/gocql"
	"github.com/labstack/echo/v4"
)

func FetchPlaylistsHandler(scyllaService database.ScyllaService) echo.HandlerFunc {
	return func(c echo.Context) error {
		userID := c.Get("userID").(string)

		playlists, err := scyllaService.FetchPlaylists(userID)
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, "Failed to fetch playlists")
		}

		return c.JSON(http.StatusOK, playlists)
	}
}

func AddPlaylistHandler(dbService database.ScyllaService) echo.HandlerFunc {
	return func(c echo.Context) error {
		userID := c.Get("userID").(string)
		playlist := new(models.Playlist)
		if err := c.Bind(playlist); err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, "Failed to bind playlist")
		}
		playlist.PlaylistID = gocql.TimeUUID()
		playlist.UserID = userID
		err := dbService.AddPlaylist(playlist.PlaylistID, playlist.UserID, playlist.Name, playlist.Description)
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, "Failed to add playlist")
		}
		return c.JSON(http.StatusCreated, playlist)
	}
}

func UpdatePlaylistHandler(scyllaService database.ScyllaService) echo.HandlerFunc {
	return func(c echo.Context) error {
		playlistID := c.Param("playlist_id")

		playlistUUID, err := gocql.ParseUUID(playlistID)
		if err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, "Invalid playlist ID")
		}

		var updateData models.Playlist
		if err := c.Bind(&updateData); err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, "Invalid request body")
		}

		err = scyllaService.UpdatePlaylist(playlistUUID, updateData.Name, updateData.Description)
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, "Failed to update playlist")
		}

		return c.JSON(http.StatusCreated, echo.Map{
			"message": "Update playlist successfully",
		})
	}
}

func AddSongToPlaylistHandler(scyllaService database.ScyllaService) echo.HandlerFunc {
	return func(c echo.Context) error {
		playlistID := c.Param("playlist_id")
		songID := c.Param("song_id")
		userID := c.Get("userID").(string)

		playlistUUID, err := gocql.ParseUUID(playlistID)
		if err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, "Invalid playlist ID")
		}

		songUUID, err := gocql.ParseUUID(songID)
		if err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, "Invalid song ID")
		}

		err = scyllaService.AddSongToPlaylist(playlistUUID, userID, songUUID, time.Now()) // use current time stamp for now.
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, "Failed to add song to playlist")
		}

		return c.JSON(http.StatusCreated, echo.Map{
			"message": "Song added to playlist successfully",
		})
	}
}

func RemoveSongFromPlaylistHandler(scyllaService database.ScyllaService) echo.HandlerFunc {
	return func(c echo.Context) error {
		playlistID := c.Param("playlist_id")
		songID := c.Param("song_id")

		playlistUUID, err := gocql.ParseUUID(playlistID)
		if err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, "Invalid playlist ID")
		}
		songUUID, err := gocql.ParseUUID(songID)
		if err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, "Invalid song ID")
		}
		err = scyllaService.RemoveSongFromPlaylist(playlistUUID, songUUID)
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, "Failed to remove song from playlist")
		}

		return c.JSON(http.StatusOK, echo.Map{
			"message": "Song removed from playlist successfully",
		})
	}
}

func GetSongsInPlaylistHandler(scyllaService database.ScyllaService) echo.HandlerFunc {
	return func(c echo.Context) error {
		playlistID := c.Param("playlist_id")

		playlistUUID, err := gocql.ParseUUID(playlistID)
		if err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, "Invalid playlist ID")
		}

		playllistSongs, err := scyllaService.GetSongsInPlaylist(playlistUUID)
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, "Failed to get liked songs")
		}

		return c.JSON(http.StatusOK, playllistSongs)
	}
}

func RemovePlaylistHandler(scyllaService database.ScyllaService) echo.HandlerFunc {
	return func(c echo.Context) error {
		playlistID := c.Param("playlist_id")
		playlistUUID, err := gocql.ParseUUID(playlistID)
		if err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, "Invalid playlist ID")
		}

		err = scyllaService.RemovePlaylist(playlistUUID)
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, "Failed to remove playlist")
		}

		return c.JSON(http.StatusOK, echo.Map{
			"message": "Playlist removed successfully",
		})
	}
}

func GetArtistPlaylistsHandler(scyllaService database.ScyllaService) echo.HandlerFunc {
	return func(c echo.Context) error {
		artistID := c.Param("artist_id")

		playlists, err := scyllaService.FetchPlaylistsByArtist(artistID)
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, "Failed to fetch artist playlists")
		}

		return c.JSON(http.StatusOK, playlists)
	}
}
