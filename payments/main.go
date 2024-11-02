package main

import (
	"context"
	"log"
	"net"
	"net/http"
	"restaurant-backend/common"
	"restaurant-backend/common/broker"
	"restaurant-backend/common/discovery"
	"restaurant-backend/common/discovery/consul"
	"restaurant-backend/payments/gateway"
	"restaurant-backend/payments/processor/payu"
	"time"

	_ "github.com/joho/godotenv/autoload"

	"google.golang.org/grpc"
)

var (
	serviceName   = "payment"
	grpcAddress   = common.Env("GRPC_ADDRESS", "localhost:2001")
	consulAddress = common.Env("CONSUL_ADDRESS", "localhost:8500")
	amqpUser      = common.Env("RABBITMQ_USER", "guest")
	amqpPassword  = common.Env("RABBITMQ_PASSWORD", "guest")
	amqpHost      = common.Env("RABBITMQ_HOST", "localhost")
	amqpPort      = common.Env("RABBITMQ_PORT", "5672")
	httpAddress   = common.Env("HTTP_ADDRESS", "localhost:8081")
)

func main() {

	registry, err := consul.NewRegistry(consulAddress, serviceName)

	if err != nil {
		panic(err)
	}

	ctx := context.Background()
	instanceId := discovery.GenerateInstanceId(serviceName)

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

	payuProcessor := payu.NewProcessor()

	channel, close := broker.Connect(amqpUser, amqpPassword, amqpHost, amqpPort)

	defer func() {
		close()
		channel.Close()
	}()

	gateway := gateway.NewGRPCGateway(registry)

	service := NewService(payuProcessor, gateway)

	amqpConsumer := NewConsumer(service)

	go amqpConsumer.Listen(channel)

	mux := http.NewServeMux()
	httpServer := NewPaymentHTTPHandler(channel)
	httpServer.registerRoutes(mux)

	go func() {
		log.Println("Starting HTTP server on", httpAddress)
		if err := http.ListenAndServe(httpAddress, mux); err != nil {
			log.Fatal(err)
		}
	}()

	grpcServer := grpc.NewServer()

	l, err := net.Listen("tcp", grpcAddress)

	log.Println("Starting gRPC server on", grpcAddress)
	
	if err != nil {
		log.Fatal(err.Error())
	}

	defer l.Close()

	if err := grpcServer.Serve(l); err != nil {
		log.Fatal(err.Error())
	}
}
