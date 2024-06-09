package main

import (
	"log"
	"net/http"
)

func main() {
	http.HandleFunc("/", rootHandler)
	http.HandleFunc("/split", splitHandler)
	log.Println("Starting bsplit server on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
