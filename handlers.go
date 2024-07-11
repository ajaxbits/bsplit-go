package main

import (
	"encoding/json"
	"html/template"
	"log"
	"net/http"
	"strconv"

	_ "embed"

	"ajaxbits.com/bsplit/db"
	"ajaxbits.com/bsplit/views"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

var templates = template.Must(template.ParseGlob("templates/*.html"))

var writeQueries = db.New(writeDb)
var readQueries = db.New(readDb)


func RootHandler(c echo.Context) error {
	return Render(c, http.StatusOK, views.Base())
}

func SplitHandler(c echo.Context) error {
	totalStr, participantsStr := c.FormValue("total"), c.FormValue("participants")
	total, err := strconv.ParseFloat(totalStr, 64)
	if err != nil {
		c.Logger().Errorf("invalid total amount: %+v", err)
		return c.String(http.StatusInternalServerError, "unable to create split")
	}
	
	participants, err := strconv.Atoi(participantsStr)
	if err != nil {
		c.Logger().Errorf("error parsing participants: %+v", err)
		return c.String(http.StatusInternalServerError, "unable to create split")
	} else if participants < 2 {
		c.Logger().Errorf("not enough participants: %+v", err)
		return c.String(http.StatusInternalServerError, "unable to create split, not enough participants")
	}

	split := Split(total, participants, Even)

	c.Logger().Infof("split: %+v", split)
	return Render(c, http.StatusOK, views.Result(participants, split))
}

func GetUsersHandler(c echo.Context) error {
	users, err := readQueries.GetAllUsers(c.Request().Context())
	if err != nil {
		c.Logger().Error("Could not list users from db")
		return c.String(http.StatusInternalServerError, "unable to list users")
	}

	return c.JSON(200, users)
}
func SearchUsersHandler(c echo.Context) error {
	users, err := readQueries.GetAllUsers(c.Request().Context())
	if err != nil {
		c.Logger().Error("Could not list users from db")
		return c.String(http.StatusInternalServerError, "unable to list users")
	}

	return Render(c, http.StatusOK, views.UsersResult(users))
}

func CreateUserHandler(c echo.Context) error {
	userName, venmoId := c.QueryParam("name"), c.QueryParam("venmo")
	if userName == "" {
		c.Logger().Error("User name field empty")
		return c.String(http.StatusInternalServerError, "unable to create user")
	}

	userUuid, err := uuid.NewV7()
	if err != nil {
		c.Logger().Errorf("could not create uuid: %+v", err)
		return c.String(http.StatusInternalServerError, "unable to create user")
	}

	user, err := writeQueries.CreateUser(ctx, db.CreateUserParams{
		Uuid:    userUuid.String(),
		Name:    userName,
		VenmoID: &venmoId,
	})
	if err != nil {
		c.Logger().Errorf("could not create user in db: %+v", err)
		return c.String(http.StatusInternalServerError, "unable to create user")
	}

	c.Logger().Infof("user: %+v", user)
	return c.NoContent(200)
}

func GetGroupsHandler(c echo.Context) error {
	groups, err := readQueries.GetAllGroups(c.Request().Context())
	if err != nil {
		c.Logger().Error("Could not list users from db")
		return c.String(http.StatusInternalServerError, "unable to list users")
	}

	return c.JSON(200, groups)
}

func CreateGroupHandler(c echo.Context) error {
	decoder := json.NewDecoder(c.Request().Body)
	var g struct {
		Name        string `json:"name"`
		Description string `json:"description"`
	}
	err := decoder.Decode(&g)
	if err != nil {
		c.Logger().Errorf("could not create json decoder: %+v", err)
		return c.String(http.StatusInternalServerError, "unable to create group")
	}

	if g.Name == "" {
		log.Println("User endpoint has no name in path")
		return c.String(http.StatusInternalServerError, "unable to create group")
	}

	groupUuid, err := uuid.NewV7()
	if err != nil {
		c.Logger().Errorf("could not create uuid: %+v", err)
		return c.String(http.StatusInternalServerError, "unable to create group")
	}

	group, err := writeQueries.CreateGroup(ctx, db.CreateGroupParams{
		Uuid:        groupUuid.String(),
		Name:        g.Name,
		Description: &g.Description,
	})
	if err != nil {
		c.Logger().Errorf("could not create group in db: %+v", err)
		return c.String(http.StatusInternalServerError, "unable to create group")
	}

	c.Logger().Infof("group: %+v", group)
	return c.NoContent(200)
}

func TransactionHandler(c echo.Context) error {
	decoder := json.NewDecoder(c.Request().Body)
	var t struct {
		Description  string  `json:"description"`
		Amount       int64   `json:"amount"`
		Date         int64   `json:"date"`
		PaidBy       string  `json:"paid_by"`
		GroupUuid    *string `json:"group_uuid"`
		Participants []struct {
			UserUuid string `json:"user_uuid"`
			Share    int64  `json:"share"`
		} `json:"participants"`
	}
	err := decoder.Decode(&t)
	if err != nil {
		c.Logger().Errorf("could not create json decoder: %+v", err)
		return c.String(http.StatusInternalServerError, "unable to create transaction")
	}

	if t.Description == "" || t.Amount == 0 || len(t.Participants) <= 0 {
		c.Logger().Errorf("transaction has invalid format: %+v", err)
		return c.String(http.StatusInternalServerError, "unable to create transaction")
	}

	txnUuid, err := uuid.NewV7()
	if err != nil {
		c.Logger().Errorf("could not create uuid: %+v", err)
		return c.String(http.StatusInternalServerError, "unable to create transaction")
	}

	txn, err := writeQueries.CreateTransaction(ctx, db.CreateTransactionParams{
		Uuid:        txnUuid.String(),
		Description: t.Description,
		Type:        "expense",
		Amount:      t.Amount,
		Date:        t.Date,
		PaidBy:      t.PaidBy,
		GroupUuid:   t.GroupUuid,
	})
	if err != nil {
		c.Logger().Errorf("could not create transaction in db: %+v", err)
		return c.String(http.StatusInternalServerError, "unable to create transaction")
	}

	for _, p := range t.Participants {
		txnParticipantUuid, err := uuid.NewV7()
		if err != nil {
			c.Logger().Errorf("could not create uuid: %+v", err)
			return c.String(http.StatusInternalServerError, "unable to create transaction")
		}
		_, err = writeQueries.CreateTransactionParticipant(ctx, db.CreateTransactionParticipantParams{
			Uuid:     txnParticipantUuid.String(),
			TxnUuid:  txnUuid.String(),
			UserUuid: p.UserUuid,
			Share:    p.Share,
		})
		if err != nil {
			c.Logger().Errorf("could not create transaction participant entry in db: %+v", err)
			return c.String(http.StatusInternalServerError, "unable to create transaction")
		}
	}

	c.Logger().Infof("transaction: %+v", txn)
	return c.NoContent(200)
}
