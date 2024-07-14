package main

import (
	"context"
	_ "embed"
	"log/slog"
	"os"

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

	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	e.Use(middleware.RequestLoggerWithConfig(middleware.RequestLoggerConfig{
		LogStatus:   true,
		LogURI:      true,
		LogError:    true,
		HandleError: false,
		LogValuesFunc: func(c echo.Context, v middleware.RequestLoggerValues) error {
			if v.Error == nil {
				logger.LogAttrs(context.Background(), slog.LevelInfo, "REQUEST",
					slog.String("uri", v.URI),
					slog.Int("status", v.Status),
				)
			} else {
				logger.LogAttrs(context.Background(), slog.LevelError, "REQUEST_ERROR",
					slog.String("uri", v.URI),
					slog.Int("status", v.Status),
					slog.String("err", v.Error.Error()),
				)
			}
			return nil
		},
	}))

	e.Use(middleware.GzipWithConfig(middleware.GzipConfig{
		Level: 5,
	}))

	e.GET("/", handlers.RootHandler)
	e.POST("/user", handlers.CreateUserHandler)
	e.GET("/users", handlers.GetUsersHandler)
	e.POST("/users", handlers.SearchUsersHandler)
	e.GET("/groups", handlers.GetGroupsHandler)
	e.POST("/group", handlers.CreateGroupHandler)
	e.PUT("/txn", handlers.TransactionHandler)
	e.POST("/split", handlers.SplitHandler)
	
	e.Logger.Fatal(e.Start("localhost:8080"))
}
