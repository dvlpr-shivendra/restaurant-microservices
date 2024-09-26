package main

import (
	_ "github.com/joho/godotenv/autoload"
	"log"
	"net/http"
	"restaurant-backend/common"
)

var (
	httpAddress = common.Env("HTTP_ADDRESS", ":8080")
)

func main() {
	mux := http.NewServeMux()

	handler := NewHandler()
	handler.registerRoutes(mux)

	log.Printf("Starting server on %s", httpAddress)

	if err := http.ListenAndServe(httpAddress, mux); err != nil {
		panic(err)
	}
}
