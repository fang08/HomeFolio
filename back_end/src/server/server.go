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
	const urlPrefix = "/api/v1"

	e = echo.New()
	e.HideBanner = true
	e.Logger.SetLevel(log.ERROR)
	e.Use(middleware.Logger())
	e.Use(middleware.JWTWithConfig(middleware.JWTConfig{
		SigningKey: []byte(handler.Key),
		Skipper: func(c echo.Context) bool {
			// Skip authentication for and signup login requests
			if c.Request().Method == "OPTIONS" || c.Path() == urlPrefix+"/login" || c.Path() == urlPrefix+"/signup" || c.Path() == urlPrefix+"/tuplecount" || c.Path() == urlPrefix+"/ws/:username" {
				return true
			}
			return false
		},
	}))

	// CORS config
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins:     []string{"http://localhost:3000", "http://fang08.github.io/HomeFolio", "https://fang08.github.io/HomeFolio", "http://homefolio.fang.hojason117.com", "https://homefolio.fang.hojason117.com"},
		AllowMethods:     []string{echo.GET, echo.PUT, echo.POST, echo.DELETE},
		AllowCredentials: true,
		AllowHeaders:     []string{echo.HeaderAuthorization, echo.HeaderContentType},
	}))

	// Routes
	e.GET("/api/v1/ws/:username", h.NotifHandler.GetConnection)
	e.POST("/api/v1/signup", h.UserHandler.Signup)
	e.POST("/api/v1/login", h.UserHandler.Login)
	e.GET("/api/v1/userInfo/:uid", h.UserHandler.FetchUserInfo)
	e.POST("/api/v1/updateUserInfo", h.UserHandler.UpdateUserInfo)
	e.GET("/api/v1/houseInfo/:hid", h.HouseHandler.FetchSingleHouseInfo)
	e.GET("/api/v1/houseInfo", h.HouseHandler.FetchRegionHouseInfo)
	e.GET("/api/v1/topliked", h.HouseHandler.FetchTopLikedHouses)
	e.GET("/api/v1/topviewed", h.HouseHandler.FetchTopViewedHouses)
	e.GET("/api/v1/tuplecount", h.HouseHandler.GetTupleCount)
	e.DELETE("/api/v1/deletehouse/:hid", h.HouseHandler.DeleteHouse)
	e.GET("/api/v1/searchhouse", h.HouseHandler.SearchHouse)
	e.GET("/api/v1/userInfo/sellHouses/:uid", h.UserHandler.FetchSellHouse)
	e.GET("/api/v1/userInfo/likedHouses/:uid", h.UserHandler.FetchLikedHouse)
	e.GET("/api/v1/userInfo/viewedHouses/:uid", h.UserHandler.FetchViewedHouse)
	e.POST("/api/v1/sell", h.HouseHandler.Sell)
	e.POST("/api/v1/updateHouseInfo/:hid", h.HouseHandler.UpdateHouseInfo)
	e.POST("/api/v1/like", h.HouseHandler.AddLike)
	e.DELETE("/api/v1/unlike/:uidhid", h.HouseHandler.RemoveLike)
	e.GET("/api/v1/houseInfo/likedUsers/:hid", h.HouseHandler.FetchLikedUser)
	e.GET("/api/v1/houseInfo/viewedUsers/:hid", h.HouseHandler.FetchViewedUser)
	e.POST("/api/v1/viewed", h.HouseHandler.AddViewed)
	e.POST("/api/v1/buyhouse/:hid", h.HouseHandler.BuyHouse)
	e.GET("/api/v1/userInfo/boughtHouses/:uid", h.UserHandler.FetchBoughtHouse)
	e.GET("/api/v1/userInfo/isPopular/:uid", h.UserHandler.IsPopularUser)
	e.GET("/api/v1/userInfo/isActive/:uid", h.UserHandler.IsActiveUser)

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
