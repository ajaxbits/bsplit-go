package db

import (
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
	"time"
)

func Initialize() (*sql.DB, error) {
	db, err := sql.Open("sqlite3", "./expenses.db")
	if err != nil {
		return nil, err
	}

	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(25)
	db.SetConnMaxLifetime(5 * time.Minute)

	_, err = db.Exec("PRAGMA journal_mode=WAL;")
	if err != nil {
		return nil, err
	}

	_, err = db.Exec("PRAGMA synchronous=NORMAL;")
	if err != nil {
		return nil, err
	}

	if err := createTables(db); err != nil {
		return nil, err
	}
	return db, nil
}

func createTables(db *sql.DB) error {
	tables := []string{
		`CREATE TABLE IF NOT EXISTS users (
            id TEXT PRIMARY KEY,
            name TEXT NOT NULL,
            created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
        );`,
		`CREATE TABLE IF NOT EXISTS transactions (
            id TEXT PRIMARY KEY,
            description TEXT NOT NULL,
            amount INTEGER NOT NULL,
            date TIMESTAMP NOT NULL,
            paid_by INTEGER,
            FOREIGN KEY (paid_by) REFERENCES users(id)
        );`,
		`CREATE TABLE IF NOT EXISTS expense_participants (
            id TEXT PRIMARY KEY,
            expense_id INTEGER,
            user_id INTEGER,
            share INTEGER,
            FOREIGN KEY (expense_id) REFERENCES transactions(id),
            FOREIGN KEY (user_id) REFERENCES users(id)
        );`,
		`CREATE TABLE IF NOT EXISTS ledgers (
            id TEXT PRIMARY KEY,
            from_user_id INTEGER,
            to_user_id INTEGER,
            amount INTEGER,
            FOREIGN KEY (from_user_id) REFERENCES users(id),
            FOREIGN KEY (to_user_id) REFERENCES users(id)
        );`,
	}

	for _, table := range tables {
		if _, err := db.Exec(table); err != nil {
			return err
		}
	}

	return nil
}
