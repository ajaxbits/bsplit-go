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

	http.HandleFunc("/", handlers.RootHandler)
	http.HandleFunc("/split", handlers.SplitHandler)
	log.Println("Starting bsplit server on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
