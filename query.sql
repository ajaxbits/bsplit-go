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



-- name: CreateTransaction :one
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



