package main

import (
	"context"
	_ "embed"

	"ajaxbits.com/bsplit/db"

	"github.com/a-h/templ"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/labstack/gommon/log"
)

//go:embed schema.sql
var ddl string

var ctx = context.Background()
var readDb, writeDb = db.Init()

func Render(ctx echo.Context, statusCode int, t templ.Component) error {
	buf := templ.GetBuffer()
	defer templ.ReleaseBuffer(buf)

	if err := t.Render(ctx.Request().Context(), buf); err != nil {
		return err
	}

	return ctx.HTML(statusCode, buf.String())
}

func main() {
	defer readDb.Close()
	defer writeDb.Close()

	if _, err := writeDb.ExecContext(ctx, ddl); err != nil {
		log.Fatal(err)
	}

	e := echo.New()
	e.HideBanner = true
	e.Logger.SetLevel(log.INFO)
	e.Pre(middleware.RemoveTrailingSlash())
	e.Use(middleware.Logger()) // TODO: use the newer logger eventually
	e.Use(middleware.GzipWithConfig(middleware.GzipConfig{
		Level: 5,
	}))

	e.GET("/", RootHandler)
	e.GET("/users", GetUsersHandler)
	e.PUT("/user", CreateUserHandler)
	e.GET("/groups", GetGroupsHandler)
	e.PUT("/group", CreateGroupHandler)
	e.PUT("/txn", TransactionHandler)
	e.POST("/split", SplitHandler)

	e.Logger.Fatal(e.Start("localhost:8080"))
}
