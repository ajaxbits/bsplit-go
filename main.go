package main

import (
	"context"
	"log"
	"net/http"
	"time"

	"ajaxbits.com/bsplit/internal/db"
	"ajaxbits.com/bsplit/internal/handlers"
	"ajaxbits.com/bsplit/internal/models"
)

func main() {
	database, err := db.Initialize()
	if err != nil {
		log.Fatal(err)
	}
	defer database.Close()

	ctx := context.Background()

	alice, err := db.CreateUser(ctx, database, "Alice", nil)
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

	tx, err := db.CreateTransaction(ctx, database, "expense", "dinner", 10000, time.Now(), *alice, group, map[models.User]int{*alice: 500, *bob: 500})
	if err != nil {
		log.Fatal(err)
	} else {
		log.Println(tx)
	}

	debts, err := db.GetDebts(ctx, database, []models.User{*alice, *bob}, group)
	if err != nil {
		log.Fatal(err)
	} else {
		for debtor, debtMap := range *debts {
			log.Printf("user %s owes the following:", debtor.Name)
			for creditor, amount := range debtMap {
				log.Printf("  %s: $%.2f", creditor.Name, float64(amount)/100)
			}
		}
	}

	http.HandleFunc("/", handlers.RootHandler)
	http.HandleFunc("/split", handlers.SplitHandler)
	log.Println("Starting bsplit server on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
