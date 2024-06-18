package main

import (
	"context"
	"database/sql"
	_ "embed"
	"log"

	// "net/http"
	"time"

	"ajaxbits.com/bsplit/internal/db"
	"ajaxbits.com/bsplit/internal/dbc"

	// "ajaxbits.com/bsplit/internal/handlers"
	"ajaxbits.com/bsplit/internal/models"
)

//go:embed schema.sql
var ddl string

func main() {
	database, err := sql.Open("sqlite3", "./expenses.db")
	if err != nil {
		log.Fatal(err)
	}
	defer database.Close()

	ctx := context.Background()
	queries := dbc.New(database)

	if _, err := database.ExecContext(ctx, ddl); err != nil {
		log.Fatal(err)
	}

	// aliceUuid, err := uuid.NewV7()
	// if err != nil {
	// 	log.Fatal(err)
	// }

	alice, err := db.CreateUser(ctx, database, "Alice", nil)
	// alice, err := queries.CreateUser(ctx, dbc.CreateUserParams{
	// 	ID:      aliceUuid[:],
	// 	Name:    "Alice",
	// 	VenmoID: nil,
	// })
	if err != nil {
		log.Fatal(err)
	} else {
		log.Println(alice)
	}

	bob, err := db.CreateUser(ctx, database, "Bob", nil)
	if err != nil {
		log.Fatal(err)
	} else {
		log.Println(bob)
	}

	group, err := db.CreateGroup(ctx, database, "underthecovers", nil, []models.User{*alice, *bob})
	if err != nil {
		log.Fatal(err)
	} else {
		log.Println(group)
	}

	tx, err := db.CreateTransaction(ctx, database, "expense", "dinner", 1000, time.Now(), *alice, group, map[models.User]int{*alice: 500, *bob: 500})
	if err != nil {
		log.Fatal(err)
	} else {
		log.Println(tx)
	}

	tx2, err := db.CreateTransaction(ctx, database, "expense", "new car", 2000, time.Now(), *bob, group, map[models.User]int{*alice: 1000, *bob: 1000})
	if err != nil {
		log.Fatal(err)
	} else {
		log.Println(tx2)
	}

	debts, err := queries.GetDebts(ctx)
	if err != nil {
		log.Fatal(err)
	}

	for _, debt := range debts {
		debtor, err := queries.GetUser(ctx, debt.Debtor)
		if err != nil {
			log.Fatal(err)
		}
		creditor, err := queries.GetUser(ctx, debt.Creditor)
		if err != nil {
			log.Fatal(err)
		}

		log.Default().Printf("%s owes %s: %v", debtor.Name, creditor.Name, debt.NetAmount)
	}

	// http.HandleFunc("/", handlers.RootHandler)
	// http.HandleFunc("/split", handlers.SplitHandler)
	// log.Println("Starting bsplit server on :8080")
	// log.Fatal(http.ListenAndServe(":8080", nil))
}
