package gateway

import (
	"context"
	"log"
	"restaurant-backend/common/discovery"

	pb "restaurant-backend/common/api"
)

type gateway struct {
	registry discovery.Registry
}

func NewGRPCGateway(registry discovery.Registry) *gateway {
	return &gateway{registry}
}

func (g *gateway) CreateOrder(ctx context.Context, req *pb.CreateOrderRequest) (*pb.Order, error) {

	conn, err := discovery.ServiceConnection(ctx, "orders", g.registry)

	if err != nil {
		log.Fatalf("Failed to connect to orders service: %v", err);
	}
	
	c := pb.NewOrderServiceClient(conn)

	return c.CreateOrder(ctx, &pb.CreateOrderRequest{
		CustomerId: req.CustomerId,
		Items:      req.Items,
	})
}

func (g *gateway) GetOrder(ctx context.Context, orderId, customerId string) (*pb.Order, error) {
	conn, err := discovery.ServiceConnection(ctx, "orders", g.registry)

	if err != nil {
		log.Fatalf("Failed to connect to orders service: %v", err);
	}

	c := pb.NewOrderServiceClient(conn)

	return c.GetOrder(ctx, &pb.GetOrderRequest{
		OrderId:    orderId,
		CustomerId: customerId,
	})
}