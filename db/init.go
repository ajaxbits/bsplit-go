package db

import (
	"context"
	"database/sql"
	_ "embed"
	"log"
	"net/url"
	"runtime"

	_ "github.com/mattn/go-sqlite3"
)

//go:embed sql/schema.sql
var ddl string

var Ctx = context.Background()

var (
	readDb       *sql.DB
	writeDb      *sql.DB
	ReadQueries  *Queries
	WriteQueries *Queries
)

func Init() {
	// https://kerkour.com/sqlite-for-servers
	dbConnectionUrlParams := make(url.Values)
	dbConnectionUrlParams.Add("_txlock", "immediate")
	dbConnectionUrlParams.Add("_journal_mode", "WAL")
	dbConnectionUrlParams.Add("_busy_timeout", "5000")
	dbConnectionUrlParams.Add("_synchronous", "NORMAL")
	dbConnectionUrlParams.Add("_cache_size", "1000000000")
	dbConnectionUrlParams.Add("_foreign_keys", "true")
	dbConnectionUrlParams.Add("_temp_store", "memory")
	dbConnectionUrl := "file:expenses.db?" + dbConnectionUrlParams.Encode()

	writeDb, err := sql.Open("sqlite3", dbConnectionUrl)
	if err != nil {
		log.Fatal(err)
	}
	writeDb.SetMaxOpenConns(1)

	readDb, err := sql.Open("sqlite3", dbConnectionUrl)
	if err != nil {
		log.Fatal(err)
	}
	readDb.SetMaxOpenConns(max(4, runtime.NumCPU()))

	if _, err := writeDb.ExecContext(Ctx, ddl); err != nil {
		log.Fatal(err)
	}

	ReadQueries = New(readDb)
	WriteQueries = New(writeDb)
}

func Close() {
	readDb.Close()
	writeDb.Close()
}
