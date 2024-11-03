package main

import (
	"context"
	"log"
	"net"
	"restaurant-backend/common"
	"restaurant-backend/common/broker"
	"restaurant-backend/common/discovery"
	"restaurant-backend/common/discovery/consul"
	"time"

	"go.uber.org/zap"
	"google.golang.org/grpc"
)

var (
	serviceName   = "orders"
	grpcAddress   = common.Env("GRPC_ADDRESS", "localhost:2000")
	consulAddress = common.Env("CONSUL_ADDRESS", "localhost:8500")
	amqpUser      = common.Env("RABBITMQ_USER", "guest")
	amqpPassword  = common.Env("RABBITMQ_PASSWORD", "guest")
	amqpHost      = common.Env("RABBITMQ_HOST", "localhost")
	amqpPort      = common.Env("RABBITMQ_PORT", "5672")
	jaegerAddress = common.Env("JAEGER_ADDRESS", "localhost:4318")
)

func main() {
	logger, _ := zap.NewProduction()
	defer logger.Sync()

	zap.ReplaceGlobals(logger)

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

	channel, close := broker.Connect(amqpUser, amqpPassword, amqpHost, amqpPort)

	defer func() {
		close()
		channel.Close()
	}()

	grpcServer := grpc.NewServer()

	l, err := net.Listen("tcp", grpcAddress)

	if err != nil {
		log.Fatal(err.Error())
	}

	defer l.Close()

	store := NewStore()

	service := NewService(store)

	serviceWithTelemetry := NewTelemetryMiddleware(service)
	serviceWithLogging := NewLoggingMiddleware(serviceWithTelemetry)

	NewGrpcHandler(grpcServer, serviceWithLogging, channel)

	consumer := NewConsumer(serviceWithLogging)

	go consumer.Listen(channel)

	log.Printf("Orders service started at %s", grpcAddress)

	if err := grpcServer.Serve(l); err != nil {
		log.Fatal(err.Error())
	}
}
