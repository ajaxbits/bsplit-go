package models

import (
	"github.com/google/uuid"
	"time"
)

type User struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	Name      string    `json:"name"`
	VenmoID   *string   `json:"venmo_id,omitempty"`
}

type Group struct {
	ID          uuid.UUID `json:"id"`
	CreatedAt   time.Time `json:"created_at"`
	Name        string    `json:"name"`
	Description *string   `json:"description,omitempty"`
	Members     []User    `json:"members"`
}

type Transaction struct {
	ID           uuid.UUID    `json:"id"`
	CreatedAt    time.Time    `json:"created_at"`
	Type         string       `json:"type"`
	Description  string       `json:"description"`
	Amount       int          `json:"amount"`
	Date         time.Time    `json:"date"`
	PaidBy       User         `json:"paid_by"`
	Group        *Group       `json:"group,omitempty"`
	Participants map[User]int `json:"participants"`
}
