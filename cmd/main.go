package main

import (
	"ajaxbits.com/bsplit/internal/db"
	"ajaxbits.com/bsplit/internal/handlers"
	"ajaxbits.com/bsplit/internal/models"
	"context"
	"github.com/google/uuid"
	"log"
	"net/http"
	"time"
)

func main() {
	database, err := db.Initialize()
	if err != nil {
		log.Fatal(err)
	}
	defer database.Close()

	ctx := context.Background()

	userId := uuid.New()
	log.Println("User ID:", userId.String())

	db.CreateUser(ctx, database, &models.User{ID: userId, Name: "Alex"})

	user, err := db.GetUser(ctx, database, &userId)
	if err != nil {
		log.Fatal(err)
	} else {
		log.Println(user)
	}

	txn := models.Transaction{
		ID:          uuid.New(),
		Description: "A new transaction",
		Amount:      100 * 100,
		Date:        time.Now(),
		PaidBy:      *user,
	}
	db.CreateTransaction(ctx, database, &txn)

	transaction, err := db.GetTransaction(ctx, database, txn.ID)
	if err != nil {
		log.Fatal(err)
	} else {
		log.Println(transaction)
	}

	http.HandleFunc("/", handlers.RootHandler)
	http.HandleFunc("/split", handlers.SplitHandler)
	log.Println("Starting bsplit server on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
