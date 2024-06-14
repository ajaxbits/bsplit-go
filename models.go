package main

import (
	"database/sql"
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

func (user User) createUser(tx *sql.Tx) error {
	defer tx.Rollback()

	stmt, err := tx.Prepare("INSERT INTO Users(id, name) VALUES(?, ?)")
	if err != nil {
		return err
	}

	_, err = stmt.Exec(user.id, user.name)
	if err != nil {
		return err
	}

	if err := tx.Commit(); err != nil {
		return err
	}

	return nil
}
