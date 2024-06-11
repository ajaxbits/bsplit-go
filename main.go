package main

import (
	"ajaxbits.com/bsplit/db"
	"log"
	"net/http"
)

func main() {
	err := db.Initialize()
	if err != nil {
		log.Fatal(err)
	}

	http.HandleFunc("/", rootHandler)
	http.HandleFunc("/split", splitHandler)
	log.Println("Starting bsplit server on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
