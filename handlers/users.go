package handlers

import (
	"net/http"

	"ajaxbits.com/bsplit/db"
	"ajaxbits.com/bsplit/views"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

func GetUsersHandler(c echo.Context) error {
	users, err := db.ReadQueries.GetAllUsers(db.Ctx)
	if err != nil {
		c.Logger().Error("Could not list users from db")
		return c.String(http.StatusInternalServerError, "unable to list users")
	}

	return c.JSON(200, users)
}
func SearchUsersHandler(c echo.Context) error {
	users, err := db.ReadQueries.GetAllUsers(db.Ctx)
	if err != nil {
		c.Logger().Error("Could not list users from db")
		return c.String(http.StatusInternalServerError, "unable to list users")
	}

	return views.Render(c, http.StatusOK, views.UsersResult(users))
}

func CreateUserHandler(c echo.Context) error {
	userName, venmoId := c.FormValue("name"), c.FormValue("venmo_id")
	if userName == "" {
		c.Logger().Error("User name field empty")
		return c.String(http.StatusInternalServerError, "unable to create user")
	}

	userUuid, err := uuid.NewV7()
	if err != nil {
		c.Logger().Errorf("could not create uuid: %+v", err)
		return c.String(http.StatusInternalServerError, "unable to create user")
	}

	_, err = db.WriteQueries.CreateUser(db.Ctx, db.CreateUserParams{
		Uuid:    userUuid.String(),
		Name:    userName,
		VenmoID: &venmoId,
	})
	if err != nil {
		c.Logger().Errorf("could not create user in db: %+v", err)
		return c.String(http.StatusInternalServerError, "unable to create user")
	}

	return c.HTML(200, "Created!")
}
