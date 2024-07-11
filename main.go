package main

import (
	"context"
	_ "embed"

	"ajaxbits.com/bsplit/db"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/labstack/gommon/log"
)

//go:embed schema.sql
var ddl string

var ctx = context.Background()
var readDb, writeDb = db.Init()

func main() {
	defer readDb.Close()
	defer writeDb.Close()

	if _, err := writeDb.ExecContext(ctx, ddl); err != nil {
		log.Fatal(err)
	}

	e := echo.New()
	e.Logger.SetLevel(log.INFO)
	e.Pre(middleware.RemoveTrailingSlash())
	e.Use(middleware.Logger()) // TODO: use the newer logger eventually
	e.Use(middleware.GzipWithConfig(middleware.GzipConfig{
		Level: 5,
	}))

	e.GET("/users", GetUsersHandler)
	e.PUT("/user", CreateUserHandler)
	e.GET("/groups", GetGroupsHandler)
	e.PUT("/group", CreateGroupHandler)
	e.PUT("/txn", TransactionHandler)

	e.Logger.Fatal(e.Start("localhost:8080"))
}
