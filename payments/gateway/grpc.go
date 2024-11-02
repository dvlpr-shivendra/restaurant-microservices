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

func (g *gateway) UpdateOrderAfterPaymentLink(ctx context.Context, orderId, paymentLink string) error {
	conn, err := discovery.ServiceConnection(ctx, "orders", g.registry)

	if err != nil {
		log.Fatalf("Failed to connect to orders service: %v", err);
	}

	defer conn.Close()

	orderClient := pb.NewOrderServiceClient(conn)

	_, err = orderClient.UpdateOrder(ctx, &pb.Order{
		Id: orderId,
		Status: "waiting_payment",
		PaymentLink: paymentLink,
	})

	return err
}