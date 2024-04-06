package server

import (
	"net/http"

	"rr-backend/internal/handlers"
	mdw "rr-backend/internal/middleware"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/markbates/goth/gothic"
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

	//	e.GET("/:bucketName/:objectName", s.getMusic)

	e.GET("/health", s.healthHandler)
	e.GET("/auth/:provider/callback", s.getAuthCallback)
	e.GET("/logout", s.getLogout)
	e.GET("/auth/:provider", s.getHandleAuth)
	e.POST("/auth/google", handlers.UpsertUserHandler(s.db), mdw.JWTMiddleware)
	e.POST("/music/upload", handlers.UploadMusicHandler(s.db, s.musicService), mdw.JWTMiddleware)
	e.GET("/music", handlers.GetSongsByUser(s.db), mdw.JWTMiddleware)
	return e
}

func (s *Server) HelloWorldHandler(c echo.Context) error {
	resp := map[string]string{
		"message": "Hello World",
	}

	return c.JSON(http.StatusOK, resp)
}

// func (s *Server) getMusic(c echo.Context) error {
// 	bucketName := c.Param("bucketName")
// 	objectName := c.Param("objectName")
// 	o, err := s.musicService.ServeMusic(c, bucketName, objectName)

// 	if err != nil {
// 		return echo.NewHTTPError(http.StatusInternalServerError, "Error retrieving the object")
// 	}
// 	defer o.Close()

// 	objectInfo, err := o.Stat()
// 	if err != nil {
// 		return echo.NewHTTPError(http.StatusInternalServerError, "Error getting the object info")
// 	}

// 	// Use the helper function to serve the content
// 	helper.ServeContent(c.Response().Writer, c.Request(), objectName, objectInfo.LastModified, o)

// 	return nil
// }

func (s *Server) healthHandler(c echo.Context) error {
	return c.JSON(http.StatusOK, s.db.Health())
}

func (s *Server) getAuthCallback(c echo.Context) error {
	// No need to set c.Request().URL.Path manually
	user, err := gothic.CompleteUserAuth(c.Response().Writer, c.Request())
	if err != nil {
		return c.String(http.StatusBadRequest, err.Error()) // Return error as HTTP response
	}
	return c.JSON(http.StatusOK, user) // For simplicity, returning user as JSON
}

func (s *Server) getLogout(c echo.Context) error {
	gothic.Logout(c.Response().Writer, c.Request())
	return c.Redirect(http.StatusTemporaryRedirect, "/") // Redirect to home
}

func (s *Server) getHandleAuth(c echo.Context) error {
	gothic.GetProviderName = func(r *http.Request) (string, error) { return c.Param("provider"), nil }

	gothic.BeginAuthHandler(c.Response().Writer, c.Request())

	return nil
}
