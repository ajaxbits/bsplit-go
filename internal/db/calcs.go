package db

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	"ajaxbits.com/bsplit/internal/models"
	"github.com/google/uuid"
)

func GetDebts(ctx context.Context, database *sql.DB, users []models.User, group *models.Group) (*map[models.User]map[models.User]int, error) {
	placeholders := make([]string, len(users))
	args := make([]interface{}, len(users)+1)
	args[0] = group.ID.String()
	for i, user := range users {
		placeholders[i] = "?"
		args[i+1] = user.ID.String()
	}
	sqlUserIDsPlaceholders := strings.Join(placeholders, ", ")

	query := fmt.Sprintf(`
        WITH NetOwed AS (
            SELECT 
                tp.user_id AS debtor,
                t.paid_by AS creditor,
                SUM(tp.share) AS amount_owed
            FROM
                TransactionParticipants tp
            JOIN
                Transactions t ON tp.txn_id = t.id
            WHERE
                t.group_id = ? AND t.type = 'expense' AND tp.user_id IN (%s)
            GROUP BY
                tp.user_id, t.paid_by
        )
        SELECT 
            n1.debtor,
            n1.creditor,
            n1.amount_owed - COALESCE(n2.amount_owed, 0) AS net_amount
        FROM
            NetOwed n1
        LEFT JOIN
            NetOwed n2 ON n1.debtor = n2.creditor AND n1.creditor = n2.debtor
        WHERE
            n1.amount_owed > COALESCE(n2.amount_owed, 0)
    `, sqlUserIDsPlaceholders)

	rows, err := database.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	results := make(map[models.User]map[models.User]int)
	for rows.Next() {
		var debtorUUIDStr string
		var creditorUUIDStr string
		var totalOwed int

		if err := rows.Scan(&debtorUUIDStr, &creditorUUIDStr, &totalOwed); err != nil {
			return nil, err
		}

		debtorUUID, err := uuid.Parse(debtorUUIDStr)
		if err != nil {
			return nil, err
		}
		debtor, err := GetUser(ctx, database, &debtorUUID)
		if err != nil {
			return nil, err
		}

		creditorUUID, err := uuid.Parse(creditorUUIDStr)
		if err != nil {
			return nil, err
		}
		creditor, err := GetUser(ctx, database, &creditorUUID)
		if err != nil {
			return nil, err
		}

		if _, exists := results[*debtor]; !exists {
			results[*debtor] = make(map[models.User]int)
		}

		results[*debtor][*creditor] = totalOwed
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return &results, nil
}
