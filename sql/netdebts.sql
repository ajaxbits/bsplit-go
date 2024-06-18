WITH NetOwed AS (
    SELECT 
        tp1.user_id AS debtor,
        t.paid_by AS creditor,
        SUM(tp1.share) AS amount_owed
    FROM
        TransactionParticipants tp1
    JOIN
        Transactions t ON tp1.txn_id = t.id
    WHERE
        t.type = 'expense'
        -- and t.group_id = ? 
        -- AND (? IS NULL OR tp1.user_id IN (?))
        -- AND (? IS NULL OR t.paid_by IN (?))
    GROUP BY
        tp1.user_id, t.paid_by
)
SELECT 
    CASE
        WHEN n1.amount_owed > COALESCE(n2.amount_owed, 0) THEN n1.debtor
        ELSE n1.creditor
    END AS debtor,
    CASE
        WHEN n1.amount_owed > COALESCE(n2.amount_owed, 0) THEN n1.creditor
        ELSE n1.debtor
    END AS creditor,
    ABS(n1.amount_owed - COALESCE(n2.amount_owed, 0)) AS net_amount
FROM
    NetOwed n1
LEFT JOIN
    NetOwed n2 ON n1.debtor = n2.creditor AND n1.creditor = n2.debtor
WHERE
    n1.amount_owed != COALESCE(n2.amount_owed, 0)

