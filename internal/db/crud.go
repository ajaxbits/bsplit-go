package db

import (
	"context"
	"database/sql"
	"time"

	"ajaxbits.com/bsplit/internal/models"
	"github.com/google/uuid"
)

func CreateUser(ctx context.Context, db *sql.DB, name string, venmoID *string) (*models.User, error) {
	user := &models.User{
		ID:      uuid.New(),
		Name:    name,
		VenmoID: venmoID,
	}

	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	query := `INSERT INTO Users (id, created_at, name, venmo_id) VALUES (?, CURRENT_TIMESTAMP, ?, ?) RETURNING created_at`
	err = tx.QueryRowContext(ctx, query, user.ID.String(), user.Name, user.VenmoID).Scan(&user.CreatedAt)
	if err != nil {
		return nil, err
	}

	err = tx.Commit()
	if err != nil {
		return nil, err
	}

	return user, nil
}

func GetUser(ctx context.Context, db *sql.DB, id *uuid.UUID) (*models.User, error) {
	user := &models.User{}
	query := `SELECT id, name, created_at FROM Users WHERE id = ?`
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

func UpdateUser(ctx context.Context, db *sql.DB, user *models.User) error {
	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	query := `UPDATE Users SET name = ? WHERE id = ?`
	_, err = tx.ExecContext(ctx, query, user.Name, user.ID.String())
	if err != nil {
		return err
	}

	return tx.Commit()
}

func CreateGroup(ctx context.Context, db *sql.DB, name string, description *string, members []models.User) (*models.Group, error) {
	group := &models.Group{
		ID:          uuid.New(),
		Name:        name,
		Description: description,
		Members:     members,
	}

	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	query := `INSERT INTO Groups (id, created_at, name, description) VALUES (?, CURRENT_TIMESTAMP, ?, ?) RETURNING created_at`
	err = tx.QueryRowContext(ctx, query, group.ID.String(), group.Name, group.Description).Scan(&group.CreatedAt)
	if err != nil {
		return nil, err
	}

	for _, member := range group.Members {
		query := `INSERT INTO GroupMembers (id, group_id, user_id) VALUES (?, ?, ?)`
		_, err := tx.ExecContext(ctx, query, uuid.New(), group.ID.String(), member.ID.String())
		if err != nil {
			return nil, err
		}
	}

	err = tx.Commit()
	if err != nil {
		return nil, err
	}

	return group, nil
}

func GetGroup(ctx context.Context, db *sql.DB, id *uuid.UUID) (*models.Group, error) {
	group := &models.Group{}
	query := `SELECT id, name, created_at FROM Users WHERE id = ?`
	var uuidStr string
	err := db.QueryRowContext(ctx, query, id.String()).Scan(&uuidStr, &group.Name, &group.CreatedAt)
	if err != nil {
		return nil, err
	}

	userId, err := uuid.Parse(uuidStr)
	if err != nil {
		return nil, err
	}
	group.ID = userId

	return group, nil
}

func CreateTransaction(ctx context.Context, db *sql.DB, txnType string, desc string, amt int, date time.Time, paidBy models.User, group *models.Group, participants map[models.User]int) (*models.Transaction, error) {
	tx := &models.Transaction{
		ID:           uuid.New(),
		Type:         txnType,
		Description:  desc,
		Amount:       amt,
		Date:         date,
		PaidBy:       paidBy,
		Group:        group,
		Participants: participants,
	}

	dbTx, err := db.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}
	defer dbTx.Rollback()

	query := `INSERT INTO Transactions (id, created_at, type, description, amount, date, paid_by, group_id) VALUES (?, CURRENT_TIMESTAMP, ?, ?, ?, ?, ?, ?) RETURNING created_at`
	err = dbTx.QueryRowContext(ctx, query, tx.ID.String(), tx.Type, tx.Description, tx.Amount, tx.Date, tx.PaidBy.ID.String(), tx.Group.ID.String()).Scan(&tx.CreatedAt)
	if err != nil {
		return nil, err
	}

	for participant, share := range tx.Participants {
		query := `INSERT INTO TransactionParticipants (id, txn_id, user_id, share) VALUES (?, ?, ?, ?)`
		_, err := dbTx.ExecContext(ctx, query, uuid.New(), tx.ID.String(), participant.ID.String(), share)
		if err != nil {
			return nil, err
		}
	}

	err = dbTx.Commit()
	if err != nil {
		return nil, err
	}

	return tx, nil
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
