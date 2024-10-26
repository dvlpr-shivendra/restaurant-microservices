package main

import (
	"context"
	pb "restaurant-backend/common/api"
)

type OrdersService interface {
	GetOrder(context.Context, *pb.GetOrderRequest) (*pb.Order, error)
	CreateOrder(context.Context, *pb.CreateOrderRequest, []*pb.Item) (*pb.Order, error)
	ValidateOrder(context.Context, *pb.CreateOrderRequest) ([]*pb.Item, error)
}

type OrdersStore interface {
	Create(context.Context, *pb.CreateOrderRequest, []*pb.Item) (string, error)
	Get(ctx context.Context, id, customerId string) (*pb.Order, error)
}
