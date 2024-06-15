SELECT
    tp.user_id,
    u.name,
    SUM(tp.share) AS total_owed
FROM
    TransactionParticipants tp
JOIN
    Transactions t ON tp.txn_id = t.id
JOIN
    Users u ON tp.user_id = u.id
WHERE
    t.type = 'expense'
GROUP BY
    tp.user_id;
