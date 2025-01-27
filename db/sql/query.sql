-- name: CreateUser :one
insert into Users (
    uuid
    , name
    , venmo_id
) values (
    @uuid
    , @name
    , @venmo_id
) returning *;

-- name: GetUser :one
select
    *
from
    Users
where
    uuid = @uuid;

-- name: GetAllUsers :many
select
    *
from
    Users;

-- name: DeleteUser :exec
delete from Users
where
    uuid = @uuid;



-- name: CreateGroup :one
insert into Groups (
    uuid
    , name
    , description
) values (
    @uuid
    , @name
    , @description
) returning *;

-- name: GetGroup :one
select
    *
from
    Groups
where
    uuid = @uuid;

-- name: GetAllGroups :many
select
    *
from
    Groups;

-- name: DeleteGroup :exec
delete from Groups
where
    uuid = @uuid;



-- name: CreateTransaction :one
insert into Transactions (
    uuid
    , type
    , description
    , amount
    , date
    , paid_by
    , group_uuid
) values (
    @uuid
    , @type
    , @description
    , @amount
    , @date
    , @paid_by
    , @group_uuid
) returning *;



-- name: CreateTransactionParticipant :one
insert into TransactionParticipants (
    uuid
    , txn_uuid
    , user_uuid
    , share
) values (
    @uuid
    , @txn_uuid
    , @user_uuid
    , @share
) returning *;


-- name: GetDebts :many
with net_owed as (
    select 
        tp.user_uuid as debtor
        , t.paid_by as creditor
        , sum(tp.share) as amount_owed
    from
        TransactionParticipants tp
    join
        Transactions t on tp.txn_uuid = t.uuid
    where
        t.type = 'expense' 
        and tp.user_uuid <> t.paid_by
    group by
        tp.user_uuid, t.paid_by
), aggregate_net_owed as (
    select
        debtor_net_owed.debtor
        , debtor_net_owed.creditor
        , sum(debtor_net_owed.amount_owed - coalesce(creditor_net_owed.amount_owed, 0)) as net_amount
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
    net_amount > 0;
