package main

import (
	"go_todo/database"
	"go_todo/route"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	e := echo.New()

	// middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"https://www.thunderclient.io"},
		AllowMethods: []string{http.MethodGet, http.MethodPost, http.MethodPut, http.MethodDelete},
	}))

	// db
	database.Init()
	db := database.GetDB()
	defer database.CloseDB()

	// routing
	route.NewRouter(e, db)
	// start
	e.Logger.Fatal(e.Start(":8080"))
}
