package main

import (
	"context"
	"log"
	pb "restaurant-backend/common/api"

	"google.golang.org/grpc"
)

type grpcHandler struct {
	pb.UnimplementedOrderServiceServer
	service OrdersService
}

func NewGrpcHandler(grpcServer *grpc.Server, service OrdersService) *grpcHandler {
	handler := &grpcHandler{
		service: service,
	}

	pb.RegisterOrderServiceServer(grpcServer, handler)

	return &grpcHandler{}
}

func (h *grpcHandler) CreateOrder(context.Context, *pb.CreateOrderRequest) (*pb.Order, error) {
	log.Println("New order received")

	o := &pb.Order{
		ID: "123",
	}
	return o, nil
}
