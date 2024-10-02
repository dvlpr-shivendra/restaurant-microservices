package main

import (
	"context"
	"log"
	"net"
	"restaurant-backend/common"
	"restaurant-backend/common/discovery"
	"restaurant-backend/common/discovery/consul"
	"time"

	"google.golang.org/grpc"
)

var (
	serviceName   = "orders"
	grpcAddress   = common.Env("GRPC_ADDRESS", "localhost:2000")
	consulAddress = common.Env("CONSUL_ADDRESS", "localhost:8500")
)

func main() {

	registry, err := consul.NewRegistry(consulAddress, serviceName)

	if err != nil {
		panic(err)
	}

	ctx := context.Background()
	instanceId := discovery.GenerateInstanceID(serviceName)

	if err := registry.Register(ctx, instanceId, serviceName, grpcAddress); err != nil {
		panic(err)
	}

	go func() {
		for {
			if err := registry.HealthCheck(instanceId, serviceName); err != nil {
				log.Printf("Failed to check health: %v", err)
			}

			time.Sleep(5 * time.Second)
		}
	}()

	defer registry.DeRegister(ctx, instanceId, serviceName)

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
