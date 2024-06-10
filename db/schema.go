package db

import (
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
)

func Initialize() error {
	db, err := sql.Open("sqlite3", "./expenses.db")
	if err != nil {
		return err
	}
	defer db.Close()

	if err := createTables(db); err != nil {
		return err
	}

	return nil
}

func createTables(db *sql.DB) error {
	tables := []string{
		`CREATE TABLE IF NOT EXISTS Users (
            id TEXT PRIMARY KEY,
            name TEXT NOT NULL
        );`,
		`CREATE TABLE IF NOT EXISTS Expenses (
            id INTEGER PRIMARY KEY AUTOINCREMENT,
            description TEXT NOT NULL,
            amount REAL NOT NULL,
            date DATETIME DEFAULT CURRENT_TIMESTAMP,
            paid_by INTEGER,
            FOREIGN KEY (paid_by) REFERENCES Users(id)
        );`,
		`CREATE TABLE IF NOT EXISTS ExpenseParticipants (
            id INTEGER PRIMARY KEY AUTOINCREMENT,
            expense_id INTEGER,
            user_id INTEGER,
            share REAL,
            FOREIGN KEY (expense_id) REFERENCES Expenses(id),
            FOREIGN KEY (user_id) REFERENCES Users(id)
        );`,
		`CREATE TABLE IF NOT EXISTS Ledgers (
            id INTEGER PRIMARY KEY AUTOINCREMENT,
            from_user_id INTEGER,
            to_user_id INTEGER,
            amount REAL,
            FOREIGN KEY (from_user_id) REFERENCES Users(id),
            FOREIGN KEY (to_user_id) REFERENCES Users(id)
        );`,
	}

	for _, table := range tables {
		if _, err := db.Exec(table); err != nil {
			return err
		}
	}

	return nil
}
