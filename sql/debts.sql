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
    GROUP BY
        tp1.user_id, t.paid_by
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
