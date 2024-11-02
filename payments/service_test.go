package main

import (
	"context"
	pb "restaurant-backend/common/api"
	"restaurant-backend/payments/gateway"
	"restaurant-backend/payments/processor/inmem"
	inmemRegistry "restaurant-backend/common/discovery/inmem"
	"testing"
)

func TestService(t *testing.T) {
	processor := inmem.NewProcessor()
	registry := inmemRegistry.NewRegistry()
	gateway := gateway.NewGRPCGateway(registry)
	service := NewService(processor, gateway)

	t.Run("should create payment link", func(t *testing.T) {
		link, err := service.CreatePayment(context.Background(), &pb.Order{})

		if err != nil {
			t.Errorf("Failed to create payment link: %v", err)
		}

		if link == "" {
			t.Error("CreatePayment() link is empty")
		}
	})
}
