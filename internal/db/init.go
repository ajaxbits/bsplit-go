package db

import (
	"database/sql"
	"log"
	"net/url"
	"runtime"
)

func Init() (*sql.DB, *sql.DB) {
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

	return writeDb, readDb
}
