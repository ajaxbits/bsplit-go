package main

import (
	"github.com/google/uuid"
	"time"
)

type User struct {
	id   uuid.UUID
	name string
}

type Expense struct {
	id          uuid.UUID
	description string
	amount      int64
	date        time.Time
	paid_by     User
}
