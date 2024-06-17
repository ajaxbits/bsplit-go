// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.26.0

package dbc

import (
	"database/sql"
	"time"

	"github.com/google/uuid"
)

type Group struct {
	ID          uuid.UUID      `json:"id"`
	CreatedAt   sql.NullTime   `json:"created_at"`
	Name        string         `json:"name"`
	Description sql.NullString `json:"description"`
}

type GroupMember struct {
	ID      uuid.UUID `json:"id"`
	GroupID uuid.UUID `json:"group_id"`
	UserID  uuid.UUID `json:"user_id"`
}

type Transaction struct {
	ID          uuid.UUID    `json:"id"`
	CreatedAt   sql.NullTime `json:"created_at"`
	Type        string       `json:"type"`
	Description string       `json:"description"`
	Amount      int64        `json:"amount"`
	Date        time.Time    `json:"date"`
	PaidBy      int64        `json:"paid_by"`
	GroupID     uuid.UUID    `json:"group_id"`
}

type TransactionParticipant struct {
	ID     uuid.UUID `json:"id"`
	TxnID  uuid.UUID `json:"txn_id"`
	UserID uuid.UUID `json:"user_id"`
	Share  int64     `json:"share"`
}

type User struct {
	ID        uuid.UUID      `json:"id"`
	CreatedAt sql.NullTime   `json:"created_at"`
	Name      string         `json:"name"`
	VenmoID   sql.NullString `json:"venmo_id"`
}
