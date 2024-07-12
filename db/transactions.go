package db

import (
	"context"
	"database/sql"
)

type Xn struct {
	Tx  *sql.Tx
	Qtx *Queries
}

func createTx(ctx context.Context, db *sql.DB, queries *Queries) (*Xn, error) {
	tx, err := db.Begin()
	if err != nil {
		return nil, err
	}
	qtx := queries.WithTx(tx)

	return &Xn{
		Tx:  tx,
		Qtx: qtx,
	}, nil
}

func ReadBegin(ctx context.Context) (*Xn, error) {
	return createTx(ctx, readDb, ReadQueries)
}

func WriteBegin(ctx context.Context) (*Xn, error) {
	return createTx(ctx, writeDb, WriteQueries)
}
