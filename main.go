package main

import (
	"ajaxbits.com/bsplit/db"
	"github.com/google/uuid"
	"log"
	"net/http"
)

func main() {
	db, err := db.Initialize()
	defer db.Close()
	if err != nil {
		log.Fatal(err)
	}

	tx, err := db.Begin()
	if err != nil {
		log.Fatal(err)
	}

	user := User{name: "Alex", id: uuid.New()}
	if err := user.createUser(tx); err != nil {
		log.Fatal(err)
	}

	http.HandleFunc("/", rootHandler)
	http.HandleFunc("/split", splitHandler)
	log.Println("Starting bsplit server on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
