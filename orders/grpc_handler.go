package main

import (
	"context"
	"encoding/json"
	"log"
	pb "restaurant-backend/common/api"
	"restaurant-backend/common/broker"

	amqp "github.com/rabbitmq/amqp091-go"
	"google.golang.org/grpc"
)

type grpcHandler struct {
	pb.UnimplementedOrderServiceServer
	service OrdersService
	channel *amqp.Channel
}

func NewGrpcHandler(grpcServer *grpc.Server, service OrdersService, channel *amqp.Channel) *grpcHandler {
	handler := &grpcHandler{
		service: service,
		channel: channel,
	}

	pb.RegisterOrderServiceServer(grpcServer, handler)

	return &grpcHandler{}
}

func (h *grpcHandler) CreateOrder(ctx context.Context, p *pb.CreateOrderRequest) (*pb.Order, error) {
	log.Println("New order received")

	order := &pb.Order{
		ID: "123",
	}

	marshallOrder, err := json.Marshal(order)

	if err != nil {
		return nil, err
	}

	queue, err := h.channel.QueueDeclare(broker.OrderCreatedEvent, true, false, false, false, nil)

	if err != nil {
		log.Fatal(err)
	}

	h.channel.PublishWithContext(ctx, "", queue.Name, false, false, amqp.Publishing{
		ContentType: "application/json",
		Body:        marshallOrder,
		DeliveryMode: amqp.Persistent,
	})

	return order, nil
}
