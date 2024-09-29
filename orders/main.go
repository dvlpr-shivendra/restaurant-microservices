package main

import (
	"context"
	"log"
	"net"
	"restaurant-backend/common"

	"google.golang.org/grpc"
)

var (
	grpcAddress = common.Env("GRPC_ADDRESS", "localhost:2000")
)

func main() {
	grpcServer := grpc.NewServer()

	l, err := net.Listen("tcp", grpcAddress)

	if err != nil {
		log.Fatal(err.Error())
	}

	defer l.Close()

	store := NewStore()
	service := NewService(store)
	NewGrpcHandler(grpcServer, service)
	service.CreateOrder(context.Background())

	log.Printf("Orders service started at %s", grpcAddress)

	if err := grpcServer.Serve(l); err != nil {
		log.Fatal(err.Error())
	}
}
