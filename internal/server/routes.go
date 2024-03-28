package server

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/markbates/goth/gothic"
)

func (s *Server) RegisterRoutes() http.Handler {
	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	e.GET("/", s.HelloWorldHandler)

	e.GET("/:bucketName/:objectName", s.getMusic)

	e.GET("/health", s.healthHandler)
	e.GET("/auth/:provider/callback", s.getAuthCallback)
	e.GET("/logout", s.getLogout)
	e.GET("/auth/:provider", s.getHandleAuth)
	return e
}

func (s *Server) HelloWorldHandler(c echo.Context) error {
	resp := map[string]string{
		"message": "Hello World",
	}

	return c.JSON(http.StatusOK, resp)
}

func (s *Server) getMusic(c echo.Context) error {
	bucketName := c.Param("bucketName")
	objectName := c.Param("objectName")
	o, err := s.musicService.StreamMusic(c, bucketName, objectName)

	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Error retrieving the object")
	}
	defer o.Close()

	// Set the appropriate headers for audio content
	c.Response().Header().Set(echo.HeaderContentType, "audio/mpeg")
	c.Response().Header().Set(echo.HeaderContentDisposition, "inline; filename="+objectName)

	// Use SendStream to stream the object to the response
	return c.Stream(http.StatusOK, "audio/mpeg", o)
}

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
