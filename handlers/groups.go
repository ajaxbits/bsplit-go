package handlers

import (
	"net/http"

	"ajaxbits.com/bsplit/db"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

func GetGroupsHandler(c echo.Context) error {
	groups, err := db.ReadQueries.GetAllGroups(db.Ctx)
	if err != nil {
		c.Logger().Error("Could not list users from db")
		return c.String(http.StatusInternalServerError, "unable to list users")
	}

	return c.JSON(200, groups)
}

func CreateGroupHandler(c echo.Context) error {
	groupName, description := c.FormValue("name"), c.FormValue("description")

	if groupName == "" {
		c.Logger().Error("group name is empty")
		return c.String(http.StatusInternalServerError, "unable to create group")
	}

	groupUuid, err := uuid.NewV7()
	if err != nil {
		c.Logger().Errorf("could not create uuid: %+v", err)
		return c.String(http.StatusInternalServerError, "unable to create group")
	}

	group, err := db.WriteQueries.CreateGroup(db.Ctx, db.CreateGroupParams{
		Uuid:        groupUuid.String(),
		Name:        groupName,
		Description: &description,
	})
	if err != nil {
		c.Logger().Errorf("could not create group in db: %+v", err)
		return c.String(http.StatusInternalServerError, "unable to create group")
	}

	c.Logger().Infof("group: %+v", group)

	return c.HTML(200, "Created!")
}
