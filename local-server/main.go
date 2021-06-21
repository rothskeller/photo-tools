package main

import (
	"log"
	"net/http"
)

func main() {
	server := http.FileServer(http.Dir("."))
	if err := http.ListenAndServe("localhost:4000", server); err != nil {
		log.Fatal(err)
	}
}
