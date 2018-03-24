package server

import (
	"context"
	"fmt"
	"handler"
	"time"

	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"github.com/labstack/gommon/log"
)

// NewServer : Instantiate a server
func NewServer(h *handler.Handler) (e *echo.Echo) {
	e = echo.New()
	e.HideBanner = true
	e.Logger.SetLevel(log.ERROR)
	e.Use(middleware.Logger())
	e.Use(middleware.JWTWithConfig(middleware.JWTConfig{
		SigningKey: []byte(handler.Key),
		Skipper: func(c echo.Context) bool {
			// Skip authentication for and signup login requests
			if c.Request().Method == "OPTIONS" || c.Path() == "/api/v1/login" || c.Path() == "/api/v1/signup" || c.Path() == "/api/v1/ws/:username" {
				return true
			}
			return false
		},
	}))

	// CORS config
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins:     []string{"http://localhost:3000"},
		AllowMethods:     []string{echo.GET, echo.PUT, echo.POST, echo.DELETE},
		AllowCredentials: true,
		AllowHeaders:     []string{echo.HeaderAuthorization, echo.HeaderContentType},
	}))

	// Routes
	e.GET("/api/v1/ws/:username", h.NotifHandler.GetConnection)
	e.POST("/api/v1/signup", h.UserHandler.Signup)
	e.POST("/api/v1/login", h.UserHandler.Login)
	e.GET("/api/v1/userInfo/:username", h.UserHandler.FetchUserInfo)
	e.POST("/api/v1/updateUserInfo", h.UserHandler.UpdateUserInfo)

	return e
}

// TerminalControl : Thread for terminal control
func TerminalControl(e *echo.Echo, h *handler.Handler, srvAddr string) {
	var op string

	for {
		fmt.Println("Listening on " + srvAddr)
		fmt.Print("Option('h' for help): ")
		fmt.Scanln(&op)
		if op == "q" {
			fmt.Println("Shutting down server.")
			ShutdownServer(e, h)
			break
		} else if op == "h" {
			fmt.Println("'q' to shutdown server")
			//fmt.Println("'i' to reconstruct database to default (w/ some initial collections)")
		} else if op == "r" {
			//dbReinsert()
			fmt.Println("Database reconstructed.")
		}
	}
}

// ShutdownServer : Shutdown the server
func ShutdownServer(e *echo.Echo, h *handler.Handler) {
	h.Shutdown()

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := e.Shutdown(ctx); err != nil {
		e.Logger.Fatal(err)
	}
}