package main

import (
	"encoding/json"
	"html/template"
	"log"
	"net/http"
	"strconv"

	_ "embed"

	"ajaxbits.com/bsplit/internal/db"
	"ajaxbits.com/bsplit/internal/splits"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

var templates = template.Must(template.ParseGlob("templates/*.html"))

var writeQueries = db.New(writeDb)
var readQueries = db.New(readDb)

func RootHandler(w http.ResponseWriter, r *http.Request) {
	templates.ExecuteTemplate(w, "home.html", nil)
}

func SplitHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		r.ParseForm()
		totalStr := r.FormValue("total")
		total, err := strconv.ParseFloat(totalStr, 64)
		if err != nil {
			http.Error(w, "Invalid total amount", http.StatusBadRequest)
			return
		}

		split := splits.Split(total, 3, splits.Even)

		log.Println("Split:", split)

		templates.ExecuteTemplate(w, "result.html", total)
	} else {
		log.Println("Split endpoint called without post command")
		http.Redirect(w, r, "/", http.StatusSeeOther)
	}
}

func UserHandler(c echo.Context) error {
	userName, venmoId := c.QueryParam("name"), c.QueryParam("venmo")
	if userName != "" {
		userUuid, err := uuid.NewV7()
		if err != nil {
			log.Fatal(err)
		}

		user, err := writeQueries.CreateUser(ctx, db.CreateUserParams{
			Uuid:    userUuid.String(),
			Name:    userName,
			VenmoID: &venmoId,
		})
		if err != nil {
			c.Logger().Errorf("could not create user: %+v", err)
			return c.String(http.StatusInternalServerError, "unable to create user")
		} else {
			c.Logger().Infof("user: %+v", user)
		}
	} else {
		c.Logger().Error("User name field empty")
		return c.String(http.StatusInternalServerError, "unable to create user")
	}

	return c.NoContent(200)
}

func GroupHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		log.Println("Group endpoint called with wrong method")
		http.Redirect(w, r, "/", http.StatusSeeOther)
	}

	decoder := json.NewDecoder(r.Body)
	var g struct {
		Name        string `json:"name"`
		Description string `json:"description"`
	}
	err := decoder.Decode(&g)
	if err != nil {
		log.Fatal(err)
	}

	if g.Name == "" {
		log.Println("User endpoint has no name in path")
		http.Redirect(w, r, "/", http.StatusSeeOther)
	}

	groupUuid, err := uuid.NewV7()
	if err != nil {
		log.Fatal(err)
	}

	group, err := writeQueries.CreateGroup(ctx, db.CreateGroupParams{
		Uuid:        groupUuid.String(),
		Name:        g.Name,
		Description: &g.Description,
	})
	if err != nil {
		log.Fatal(err)
	} else {
		log.Println(group)
	}
}

type NewTransaction struct {
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

func TransactionHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		log.Println("Group endpoint called with wrong method")
		http.Redirect(w, r, "/", http.StatusSeeOther)
	}

	decoder := json.NewDecoder(r.Body)
	var t NewTransaction
	err := decoder.Decode(&t)
	if err != nil {
		log.Fatal(err)
	}

	if t.Description == "" || t.Amount == 0 || len(t.Participants) > 0 {
		log.Println("Something went wrong")
		http.Redirect(w, r, "/", http.StatusSeeOther)
	}
	txnUuid, err := uuid.NewV7()
	if err != nil {
		log.Fatal(err)
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

	for _, p := range t.Participants {
		txnParticipantUuid, err := uuid.NewV7()
		if err != nil {
			log.Fatal(err)
		}
		txnParticipant, err := writeQueries.CreateTransactionParticipant(ctx, db.CreateTransactionParticipantParams{
			Uuid:     txnParticipantUuid.String(),
			TxnUuid:  txnUuid.String(),
			UserUuid: p.UserUuid,
			Share:    p.Share,
		})
		if err != nil {
			log.Fatal(err)
		} else {
			log.Println(txnParticipant)
		}
	}

	if err != nil {
		log.Fatal(err)
	} else {
		log.Println(txn)
	}
}
