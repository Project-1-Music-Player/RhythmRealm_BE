package server

import (
	"fmt"
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

func (s *Server) healthHandler(c echo.Context) error {
	return c.JSON(http.StatusOK, s.db.Health())
}

func (s *Server) getAuthCallback(c echo.Context) error {
	provider := c.Param("provider")

	c.Request().URL.Path = fmt.Sprintf("/auth/%s/callback", provider)

	// handle the callback from the provider
	user, err := gothic.CompleteUserAuth(c.Response(), c.Request())
	if err != nil {
		fmt.Fprint(c.Response(), c.Request())
		return err
	}
	fmt.Print(user)

	http.Redirect(c.Response(), c.Request(), "/", http.StatusTemporaryRedirect)
	return nil
}

func (s *Server) getLogout(c echo.Context) error {
	gothic.Logout(c.Response(), c.Request())
	http.Redirect(c.Response(), c.Request(), "/", http.StatusTemporaryRedirect)
	return nil
}

func (s *Server) getHandleAuth(c echo.Context) error {
	if user, err := gothic.CompleteUserAuth(c.Response(), c.Request()); err == nil {
		fmt.Print(user)
	} else {
		gothic.BeginAuthHandler(c.Response(), c.Request())
	}
	return nil
}
