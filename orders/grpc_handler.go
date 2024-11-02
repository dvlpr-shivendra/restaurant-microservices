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

func (h *grpcHandler) GetOrder(ctx context.Context, p *pb.GetOrderRequest) (*pb.Order, error) {
	return h.service.GetOrder(ctx, p)
}

func (h *grpcHandler) CreateOrder(ctx context.Context, p *pb.CreateOrderRequest) (*pb.Order, error) {
	log.Println("New order received")

	items, err := h.service.ValidateOrder(ctx, p)

	if err != nil {
		return nil, err
	}

	order, err := h.service.CreateOrder(ctx, p, items)

	if err != nil {
		return nil, err
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
		ContentType:  "application/json",
		Body:         marshallOrder,
		DeliveryMode: amqp.Persistent,
	})

	return order, nil
}

func (h *grpcHandler) UpdateOrder(ctx context.Context, order *pb.Order) (*pb.Order, error) {
	return h.service.UpdateOrder(ctx, order)
}
