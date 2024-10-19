package main

import (
	"context"
	pb "restaurant-backend/common/api"
	"restaurant-backend/payments/processor"
)

type service struct {
	processor processor.PaymentProcessor
}

func NewService(processor processor.PaymentProcessor) *service {
	return &service{processor}
}

func (s *service) CreatePayment(ctx context.Context, order *pb.Order) (string, error) {
	return s.processor.CreatePaymentLink(order)
}
