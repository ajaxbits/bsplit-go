CREATE TABLE IF NOT EXISTS Users (
    id BLOB PRIMARY KEY NOT NULL UNIQUE
    , created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
    , name TEXT NOT NULL
    , venmo_id TEXT
);

CREATE TABLE IF NOT EXISTS Groups (
    id BLOB PRIMARY KEY NOT NULL UNIQUE
    , created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
    , name TEXT NOT NULL
    , description TEXT
);

CREATE TABLE IF NOT EXISTS GroupMembers (
    id BLOB PRIMARY KEY NOT NULL UNIQUE
    , group_id BLOB NOT NULL
    , user_id BLOB NOT NULL
    , FOREIGN KEY (group_id) REFERENCES Groups(id)
    , FOREIGN KEY (user_id) REFERENCES Users(id)
);

CREATE TABLE IF NOT EXISTS Transactions (
    id BLOB PRIMARY KEY NOT NULL UNIQUE
    , created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
    , type TEXT CHECK(type IN ('expense', 'settle')) NOT NULL
    , description TEXT NOT NULL
    , amount INTEGER NOT NULL
    , date TIMESTAMP NOT NULL
    , paid_by INTEGER NOT NULL
    , group_id BLOB
    , FOREIGN KEY (paid_by) REFERENCES Users(id)
    , FOREIGN KEY (group_id) REFERENCES Groups(id)
);

CREATE TABLE IF NOT EXISTS TransactionParticipants (
    id BLOB PRIMARY KEY NOT NULL UNIQUE
    , txn_id BLOB NOT NULL
    , user_id BLOB NOT NULL
    , share INTEGER NOT NULL
    , FOREIGN KEY (txn_id) REFERENCES Transactions(id)
    , FOREIGN KEY (user_id) REFERENCES Users(id)
);
