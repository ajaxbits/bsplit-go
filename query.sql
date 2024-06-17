-- name: GetUsers :many
SELECT
    *
FROM
    Users;

-- name: CreateTransaction :one
INSERT INTO Transactions (
    id
    , created_at
    , type
    , description
    , amount
    , date
    , paid_by
    , group_id
) VALUES (
    @id
    , CURRENT_TIMESTAMP
    , @type
    , @description
    , @amount
    , @date
    , @paid_by
    , @group_id
) RETURNING *;


-- name: GetDebts :many
WITH net_owed AS (
    SELECT 
        tp1.user_id AS debtor,
        t.paid_by AS creditor,
        SUM(tp1.share) AS amount_owed
    FROM
        TransactionParticipants tp1
    JOIN
        Transactions t ON tp1.txn_id = t.id
    WHERE
        t.type = 'expense' AND tp1.user_id IN (sqlc.slice('users'))
    GROUP BY
        tp1.user_id, t.paid_by
) 
SELECT 
    n1.debtor,
    n1.creditor,
    n1.amount_owed - COALESCE(n2.amount_owed, 0) AS net_amount
FROM
    net_owed n1
LEFT JOIN
    net_owed n2 ON n1.debtor = n2.creditor AND n1.creditor = n2.debtor
WHERE
    n1.amount_owed > COALESCE(n2.amount_owed, 0);
