SELECT * from Transactions;
SELECT * from TransactionParticipants;

with NetOwed as (
    select 
        tp.user_uuid as debtor
        , t.paid_by as creditor
        , SUM(tp.share) as amount_owed
    from
        TransactionParticipants tp
    join
        Transactions t on tp.txn_uuid = t.uuid
    where
        t.type = 'expense' 
        and tp.user_uuid <> t.paid_by
        -- and tp.group_uuid = @group_uuid
    group by
        tp.user_uuid, t.paid_by
),
AggregateNetOwed as (
    select
        debtor_net_owed.debtor
        , debtor_net_owed.creditor
        , SUM(debtor_net_owed.amount_owed - COALESCE(creditor_net_owed.amount_owed, 0)) as net_amount
    from
        NetOwed debtor_net_owed
    left join
        NetOwed creditor_net_owed 
        on debtor_net_owed.debtor = creditor_net_owed.creditor 
        and debtor_net_owed.creditor = creditor_net_owed.debtor
    group by
        debtor_net_owed.debtor, debtor_net_owed.creditor
)
select 
    debtor
    , creditor
    , net_amount
from 
    AggregateNetOwed
where
    net_amount > 0
order by 
    debtor, creditor;
