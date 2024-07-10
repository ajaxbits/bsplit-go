create table if not exists Users (
    uuid text primary key not null unique
    , name text not null
    , venmo_id text
) strict;

create table if not exists Groups (
    uuid text primary key not null unique
    , name text not null
    , description text
) strict;

create table if not exists GroupMembers (
    uuid text primary key not null unique
    , group_uuid text not null
    , user_uuid text not null
    , foreign key (group_uuid) references Groups(uuid)
    , foreign key (user_uuid) references Users(uuid)
) strict;

create table if not exists Transactions (
    uuid text primary key not null unique
    , type text check(type in ('expense', 'settle')) not null
    , description text not null
    , amount integer not null
    , date integer not null
    , paid_by text not null
    , group_uuid text
    , foreign key (paid_by) references Users(uuid)
    , foreign key (group_uuid) references Groups(uuid)
) strict;

create table if not exists TransactionParticipants (
    uuid text primary key not null unique
    , txn_uuid text not null
    , user_uuid text not null
    , share integer not null
    , foreign key (txn_uuid) references Transactions(uuid)
    , foreign key (user_uuid) references Users(uuid)
) strict;
