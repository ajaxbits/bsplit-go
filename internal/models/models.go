package models

import (
	"github.com/google/uuid"
	"time"
)

type User struct {
	ID        uuid.UUID
	Name      string
	CreatedAt time.Time
}

type Expense struct {
	ID          uuid.UUID
	Description string
	Amount      int64
	Date        time.Time
	PaidBy      User
}
