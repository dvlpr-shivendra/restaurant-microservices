package main

import (
	"context"
	"log"
	"net/http"
	"restaurant-backend/common"
	"restaurant-backend/gateway/gateway"
	"time"

	_ "github.com/joho/godotenv/autoload"

	"restaurant-backend/common/discovery"
	"restaurant-backend/common/discovery/consul"
)

var (
	serviceName   = "gateway"
	httpAddress   = common.Env("HTTP_ADDRESS", ":8080")
	consulAddress = common.Env("CONSUL_ADDRESS", "localhost:8500")
	jaegerAddress = common.Env("JAEGER_ADDRESS", "localhost:4318")
)

func main() {
	err := common.SetGlobalTracer(context.TODO(), serviceName, jaegerAddress)

	if err != nil {
		log.Fatalf("Failed to set tracer: %v", err)
	}

	registry, err := consul.NewRegistry(consulAddress, serviceName)

	if err != nil {
		panic(err)
	}

	ctx := context.Background()
	instanceId := discovery.GenerateInstanceId(serviceName)

	if err := registry.Register(ctx, instanceId, serviceName, httpAddress); err != nil {
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

	mux := http.NewServeMux()

	ordersGateway := gateway.NewGRPCGateway(registry)

	handler := NewHandler(ordersGateway)

	handler.registerRoutes(mux)

	log.Printf("Starting server on %s", httpAddress)

	if err := http.ListenAndServe(httpAddress, mux); err != nil {
		panic(err)
	}
}
