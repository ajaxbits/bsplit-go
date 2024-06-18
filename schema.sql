create table if not exists Users (
    uuid text primary key not null unique
    , created_at timestamp default current_timestamp
    , name text not null
    , venmo_id text
);

create table if not exists Groups (
    uuid text primary key not null unique
    , created_at timestamp default current_timestamp
    , name text not null
    , description text
);

create table if not exists GroupMembers (
    uuid text primary key not null unique
    , group_uuid text not null
    , user_uuid text not null
    , foreign key (group_uuid) references groups(id)
    , foreign key (user_uuid) references users(id)
);

create table if not exists Transactions (
    uuid text primary key not null unique
    , created_at timestamp default current_timestamp
    , type text check(type in ('expense', 'settle')) not null
    , description text not null
    , amount integer not null
    , date timestamp not null
    , paid_by text not null
    , group_uuid text
    , foreign key (paid_by) references users(id)
    , foreign key (group_uuid) references groups(id)
);

create table if not exists TransactionParticipants (
    uuid text primary key not null unique
    , txn_uuid text not null
    , user_uuid text not null
    , share integer not null
    , foreign key (txn_uuid) references transactions(id)
    , foreign key (user_uuid) references users(id)
);
