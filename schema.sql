create table if not exists Users (
    uuid blob primary key not null unique
    , created_at timestamp default current_timestamp
    , name text not null
    , venmo_id text
);

create table if not exists Groups (
    uuid blob primary key not null unique
    , created_at timestamp default current_timestamp
    , name text not null
    , description text
);

create table if not exists GroupMembers (
    uuid blob primary key not null unique
    , group_uuid blob not null
    , user_uuid blob not null
    , foreign key (group_uuid) references groups(id)
    , foreign key (user_uuid) references users(id)
);

create table if not exists Transactions (
    uuid blob primary key not null unique
    , created_at timestamp default current_timestamp
    , type text check(type in ('expense', 'settle')) not null
    , description text not null
    , amount integer not null
    , date timestamp not null
    , paid_by blob not null
    , group_uuid blob
    , foreign key (paid_by) references users(id)
    , foreign key (group_uuid) references groups(id)
);

create table if not exists TransactionParticipants (
    uuid blob primary key not null unique
    , txn_uuid blob not null
    , user_uuid blob not null
    , share integer not null
    , foreign key (txn_uuid) references transactions(id)
    , foreign key (user_uuid) references users(id)
);
