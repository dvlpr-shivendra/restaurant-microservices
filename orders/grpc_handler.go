package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	pb "restaurant-backend/common/api"
	"restaurant-backend/common/broker"

	amqp "github.com/rabbitmq/amqp091-go"
	"go.opentelemetry.io/otel"
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
	q, err := h.channel.QueueDeclare(broker.OrderCreatedEvent, true, false, false, false, nil)

	if err != nil {
		log.Fatal(err)
	}

	tr := otel.Tracer("amqp")
	amqpContext, messageSpan := tr.Start(ctx, fmt.Sprintf("AMQP - publish - %s", q.Name))
	defer messageSpan.End()

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

	headers := broker.InjectAMQPHeaders(amqpContext)

	h.channel.PublishWithContext(ctx, "", queue.Name, false, false, amqp.Publishing{
		ContentType:  "application/json",
		Body:         marshallOrder,
		DeliveryMode: amqp.Persistent,
		Headers:      headers,
	})

	return order, nil
}

func (h *grpcHandler) UpdateOrder(ctx context.Context, order *pb.Order) (*pb.Order, error) {
	return h.service.UpdateOrder(ctx, order)
}
