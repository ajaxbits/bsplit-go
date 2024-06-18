-- name: CreateUser :one
insert into Users (
    id
    , created_at
    , name
    , venmo_id
) values (
    @id
    , current_timestamp
    , @name
    , @venmo_id
) returning *;

-- name: GetUser :one
select
    *
from
    Users
where
    id = @id;

-- name: GetAllUsers :many
select
    *
from
    Users;

-- name: DeleteUser :exec
delete from Users
where
    id = @id;



-- name: CreateGroup :one
insert into Groups (
    id
    , created_at
    , name
    , description
) values (
    @id
    , current_timestamp
    , @name
    , @description
) returning *;

-- name: GetGroup :one
select
    *
from
    Groups
where
    id = @id;

-- name: GetAllGroups :many
select
    *
from
    Groups;

-- name: DeleteGroup :exec
delete from Groups
where
    id = @id;



-- name: CreateTransactionRaw :one
insert into Transactions (
    id
    , created_at
    , type
    , description
    , amount
    , date
    , paid_by
    , group_id
) values (
    @id
    , current_timestamp
    , @type
    , @description
    , @amount
    , @date
    , @paid_by
    , @group_id
) returning *;



-- name: CreateTransactionParticipants :one
insert into TransactionParticipants (
    id
    , txn_id
    , user_id
    , share
) values (
    @id
    , @txn_id
    , @user_id
    , @share
) returning *;


-- name: GetDebts :many
with net_owed as (
    select 
        tp.user_id as debtor
        , t.paid_by as creditor
        , SUM(tp.share) as amount_owed
    from
        TransactionParticipants tp
    join
        Transactions t on tp.txn_id = t.id
    where
        t.type = 'expense' 
        and tp.user_id <> t.paid_by
    group by
        tp.user_id, t.paid_by
),
aggregate_net_owed as (
    select
        debtor_net_owed.debtor
        , debtor_net_owed.creditor
        , SUM(debtor_net_owed.amount_owed - COALESCE(creditor_net_owed.amount_owed, 0)) as net_amount
    from
        net_owed debtor_net_owed
    left join
        net_owed creditor_net_owed 
        on debtor_net_owed.debtor = creditor_net_owed.creditor 
        and debtor_net_owed.creditor = creditor_net_owed.debtor
    group by
        debtor_net_owed.debtor, debtor_net_owed.creditor
)
select 
    debtor
    , creditor
    , cast(net_amount as integer)
from 
    aggregate_net_owed
where
    net_amount > 0
order by 
    debtor, creditor;
