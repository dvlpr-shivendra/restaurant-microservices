package main

import (
	"context"
	pb "restaurant-backend/common/api"
	"restaurant-backend/payments/processor/inmem"
	"testing"
)

func TestService(t *testing.T) {
	processor := inmem.NewProcessor()
	service := NewService(processor)

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
