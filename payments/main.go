package main

import (
	"context"
	"log"
	"net"
	"restaurant-backend/common"
	"restaurant-backend/common/broker"
	"restaurant-backend/common/discovery"
	"restaurant-backend/common/discovery/consul"
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

	service := NewService(payuProcessor)

	amqpConsumer := NewConsumer(service)

	go amqpConsumer.Listen(channel)

	grpcServer := grpc.NewServer()

	l, err := net.Listen("tcp", grpcAddress)

	if err != nil {
		log.Fatal(err.Error())
	}

	defer l.Close()

	if err := grpcServer.Serve(l); err != nil {
		log.Fatal(err.Error())
	}
}
