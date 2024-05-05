package server

import (
	"net/http"

	"rr-backend/internal/handlers"
	mdw "rr-backend/internal/middleware"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func (s *Server) RegisterRoutes() http.Handler {

	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"http://localhost:5173"},
		AllowMethods: []string{echo.GET, echo.PUT, echo.POST, echo.DELETE, echo.OPTIONS},
		AllowHeaders: []string{"Authorization", "Content-Type", "X-Requested-With"},
	}))

	e.GET("/", s.HelloWorldHandler)
	// TODO: Reformat/structure and group endpoints
	e.GET("/health", s.healthHandler)
	// e.GET("/auth/:provider", s.getHandleAuth)
	e.POST("/auth/google", handlers.UpsertUserHandler(s.db), mdw.JWTMiddleware)

	// TODO: Add endpoint for user profile

	e.POST("/music/upload", handlers.UploadMusicHandler(s.db, s.musicService), mdw.JWTMiddleware)
	e.GET("/music", handlers.GetSongsByUser(s.db), mdw.JWTMiddleware)
	e.GET("/music/stream/:song_id", handlers.StreamMusic(s.db, s.musicService))
	e.GET("/music/search", handlers.SearchSongs(s.db))

	e.GET("/playlists", handlers.FetchPlaylistsHandler(s.db), mdw.JWTMiddleware)
	e.POST("/playlists", handlers.AddPlaylistHandler(s.db), mdw.JWTMiddleware)
	e.DELETE("/playlists/:playlist_id", handlers.RemovePlaylistHandler(s.db), mdw.JWTMiddleware)
	e.POST("/playlists/:playlist_id/songs", handlers.AddSongToPlaylistHandler(s.db), mdw.JWTMiddleware)
	e.DELETE("/playlists/:playlist_id/songs/:song_id", handlers.RemoveSongFromPlaylistHandler(s.db), mdw.JWTMiddleware)

	return e
}

func (s *Server) HelloWorldHandler(c echo.Context) error {
	resp := map[string]string{
		"message": "Hello World",
	}

	return c.JSON(http.StatusOK, resp)
}

func (s *Server) healthHandler(c echo.Context) error {
	return c.JSON(http.StatusOK, s.db.Health())
}
