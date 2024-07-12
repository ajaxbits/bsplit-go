package handlers

import (
	"encoding/json"
	"net/http"

	"ajaxbits.com/bsplit/db"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

func TransactionHandler(c echo.Context) error {
	x, err := db.WriteBegin(db.Ctx)
	if err != nil {
		c.Logger().Errorf("could not start sql txn: %+v", err)
		return c.String(http.StatusInternalServerError, "unable to create transaction")
	}

	tx, qtx := x.Tx, x.Qtx
	defer tx.Rollback()

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
	err = decoder.Decode(&t)
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

	txn, err := qtx.CreateTransaction(db.Ctx, db.CreateTransactionParams{
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
		_, err = qtx.CreateTransactionParticipant(db.Ctx, db.CreateTransactionParticipantParams{
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
	
	err = tx.Commit()
	if err != nil {
		c.Logger().Errorf("could not commit db transaction: %+v", err)
		return c.String(http.StatusInternalServerError, "unable to create transaction")
	}

	c.Logger().Infof("transaction: %+v", txn)
	return c.NoContent(200)
}
