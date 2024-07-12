package main

import (
	_ "embed"

	"ajaxbits.com/bsplit/db"
	"ajaxbits.com/bsplit/handlers"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/labstack/gommon/log"
)


func main() {
	db.Init()
	defer db.Close()

	e := echo.New()
	e.HideBanner = true
	e.Logger.SetLevel(log.INFO)
	e.Pre(middleware.RemoveTrailingSlash())
	e.Use(middleware.Logger()) // TODO: use the newer logger eventually
	e.Use(middleware.GzipWithConfig(middleware.GzipConfig{
		Level: 5,
	}))

	e.GET("/", handlers.RootHandler)
	e.POST("/user", handlers.CreateUserHandler)
	e.GET("/users", handlers.GetUsersHandler)
	e.POST("/users", handlers.SearchUsersHandler)
	e.GET("/groups", handlers.GetGroupsHandler)
	e.PUT("/group", handlers.CreateGroupHandler)
	e.PUT("/txn", handlers.TransactionHandler)
	e.POST("/split", handlers.SplitHandler)

	e.Logger.Fatal(e.Start("localhost:8080"))
}
