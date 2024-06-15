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

func CreateTransaction(ctx context.Context, db *sql.DB, transaction *models.Transaction) error {
	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	query := `INSERT INTO transactions (id, description, amount, date, paid_by) VALUES (?, ?, ?, ?, ?)`
	_, err = tx.ExecContext(ctx, query, transaction.ID.String(), transaction.Description, transaction.Amount, transaction.Date, transaction.PaidBy.ID.String())
	if err != nil {
		return err
	}

	return tx.Commit()
}

func GetTransaction(ctx context.Context, db *sql.DB, id uuid.UUID) (*models.Transaction, error) {
	transaction := &models.Transaction{}
	var transactionIDStr, paidByIDStr string

	query := `SELECT id, description, amount, date, paid_by FROM transactions WHERE id = ?`

	err := db.QueryRowContext(ctx, query, id.String()).Scan(&transactionIDStr, &transaction.Description, &transaction.Amount, &transaction.Date, &paidByIDStr)
	if err != nil {
		return nil, err
	}

	transactionID, err := uuid.Parse(transactionIDStr)
	if err != nil {
		return nil, err
	}
	transaction.ID = transactionID

	paidByID, err := uuid.Parse(paidByIDStr)
	if err != nil {
		return nil, err
	}
	transaction.PaidBy.ID = paidByID

	return transaction, nil
}
