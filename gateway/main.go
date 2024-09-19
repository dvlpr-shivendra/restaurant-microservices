package main

import (
	"log"
	"net/http"
)

const (
	port = ":8080"
)

func main() {
	mux := http.NewServeMux()

	handler := NewHandler()
	handler.registerRoutes(mux)

	log.Printf("starting server on %s", port)

	if err := http.ListenAndServe(port, mux); err != nil {
		panic(err)
	}
}