package db

import (
	"ajaxbits.com/bsplit/internal/models"
	"context"
	"database/sql"
	"github.com/google/uuid"
)

func CreateUser(ctx context.Context, db *sql.DB, user *models.User) error {
	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	query := `INSERT INTO users (id, name, created_at) VALUES (?, ?, CURRENT_TIMESTAMP) RETURNING created_at`
	err = tx.QueryRowContext(ctx, query, user.ID.String(), user.Name).Scan(&user.CreatedAt)
	if err != nil {
		return err
	}

	return tx.Commit()
}

func GetUser(ctx context.Context, db *sql.DB, id *uuid.UUID) (*models.User, error) {
	user := &models.User{}
	query := `SELECT id, name, created_at FROM users WHERE id = ?`
	var uuidStr string
	err := db.QueryRowContext(ctx, query, id.String()).Scan(&uuidStr, &user.Name, &user.CreatedAt)
	if err != nil {
		return nil, err
	}

	userId, err := uuid.Parse(uuidStr)
	if err != nil {
		return nil, err
	}
	user.ID = userId

	return user, nil
}

// Create a new transaction
// func CreateTransaction(ctx context.Context, db *sql.DB, transaction *models.Transaction) error {
// 	tx, err := db.BeginTx(ctx, nil)
// 	if err != nil {
// 		return err
// 	}
// 	defer tx.Rollback()
//
// 	query := `INSERT INTO transactions (user_id, amount, created_at) VALUES (?, ?, ?)`
// 	result, err := tx.ExecContext(ctx, query, transaction.UserID, transaction.Amount, transaction.CreatedAt)
// 	if err != nil {
// 		return err
// 	}
// 	id, err := result.LastInsertId()
// 	if err != nil {
// 		return err
// 	}
// 	transaction.ID = id
//
// 	return tx.Commit()
// }

// Get a transaction by ID
// func GetTransaction(ctx context.Context, db *sql.DB, id int64) (*models.Transaction, error) {
// 	transaction := &models.Transaction{}
// 	query := `SELECT id, user_id, amount, created_at FROM transactions WHERE id = ?`
// 	err := db.QueryRowContext(ctx, query, id).Scan(&transaction.ID, &transaction.UserID, &transaction.Amount, &transaction.CreatedAt)
// 	if err != nil {
// 		return nil, err
// 	}
// 	return transaction, nil
// }
