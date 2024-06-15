package main

import (
	"ajaxbits.com/bsplit/internal/db"
	"ajaxbits.com/bsplit/internal/handlers"
	"ajaxbits.com/bsplit/internal/models"
	"context"
	"github.com/google/uuid"
	"log"
	"net/http"
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

	http.HandleFunc("/", handlers.RootHandler)
	http.HandleFunc("/split", handlers.SplitHandler)
	log.Println("Starting bsplit server on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
