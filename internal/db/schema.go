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
		`CREATE TABLE IF NOT EXISTS Users (
            id TEXT PRIMARY KEY NOT NULL UNIQUE,
            created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
            name TEXT NOT NULL,
            venmo_id TEXT
        );`,
		`CREATE TABLE IF NOT EXISTS Groups (
            id TEXT PRIMARY KEY NOT NULL UNIQUE,
            created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
            name TEXT NOT NULL,
            description TEXT
        );`,
		`CREATE TABLE IF NOT EXISTS GroupMembers (
            id TEXT PRIMARY KEY NOT NULL UNIQUE,
            group_id TEXT NOT NULL,
            user_id TEXT NOT NULL,
            FOREIGN KEY (group_id) REFERENCES Groups(id),
            FOREIGN KEY (user_id) REFERENCES Users(id)
        );`,
		`CREATE TABLE IF NOT EXISTS Transactions (
            id TEXT PRIMARY KEY NOT NULL UNIQUE,
            created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
            type TEXT CHECK(type IN ('expense', 'settle')) NOT NULL,
            description TEXT NOT NULL,
            amount INTEGER NOT NULL,
            date TIMESTAMP NOT NULL,
            paid_by INTEGER NOT NULL,
            group_id TEXT,
            FOREIGN KEY (paid_by) REFERENCES Users(id)
            FOREIGN KEY (group_id) REFERENCES Groups(id)
        );`,
		`CREATE TABLE IF NOT EXISTS TransactionParticipants (
            id TEXT PRIMARY KEY NOT NULL UNIQUE,
            txn_id INTEGER NOT NULL,
            user_id INTEGER NOT NULL,
            share INTEGER NOT NULL,
            FOREIGN KEY (txn_id) REFERENCES Transactions(id),
            FOREIGN KEY (user_id) REFERENCES Users(id)
        );`,
	}

	for _, table := range tables {
		if _, err := db.Exec(table); err != nil {
			return err
		}
	}

	return nil
}
