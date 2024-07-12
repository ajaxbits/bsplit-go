package handlers

import (
	"net/http"

	"ajaxbits.com/bsplit/views"
	"github.com/labstack/echo/v4"
)

func RootHandler(c echo.Context) error {
	return views.Render(c, http.StatusOK, views.Base())
}
