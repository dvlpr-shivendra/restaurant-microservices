package gateway

import (
	"context"
	pb "restaurant-backend/common/api"
)

type OrdersGateway interface {
	CreateOrder(ctx context.Context, order *pb.CreateOrderRequest) (*pb.Order, error)
	GetOrder(ctx context.Context, orderId, customerId string) (*pb.Order, error)
}
