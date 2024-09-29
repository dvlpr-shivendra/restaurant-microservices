package main

import (
	"log"
	"net/http"
	"restaurant-backend/common"

	_ "github.com/joho/godotenv/autoload"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	pb "restaurant-backend/common/api"
)

var (
	httpAddress          = common.Env("HTTP_ADDRESS", ":8080")
	ordersServiceAddress = "localhost:2000"
)

func main() {
	conn, err := grpc.NewClient(ordersServiceAddress, grpc.WithTransportCredentials(insecure.NewCredentials()))

	if err != nil {
		log.Fatalf("Could not connect to orders service: %v", err)
	}

	defer conn.Close()

	log.Printf("Orders service started at %s", ordersServiceAddress)

	c := pb.NewOrderServiceClient(conn)

	mux := http.NewServeMux()

	handler := NewHandler(c)
	handler.registerRoutes(mux)

	log.Printf("Starting server on %s", httpAddress)

	if err := http.ListenAndServe(httpAddress, mux); err != nil {
		panic(err)
	}
}
