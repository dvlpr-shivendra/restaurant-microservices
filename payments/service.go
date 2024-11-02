package main

import (
	"context"
	pb "restaurant-backend/common/api"
	"restaurant-backend/payments/gateway"
	"restaurant-backend/payments/processor"
)

type service struct {
	processor processor.PaymentProcessor
	gateway   gateway.OrdersGateway
}

func NewService(processor processor.PaymentProcessor, gateway gateway.OrdersGateway) *service {
	return &service{processor, gateway}
}

func (s *service) CreatePayment(ctx context.Context, order *pb.Order) (string, error) {
	link, err := s.processor.CreatePaymentLink(order)

	if err != nil {
		return "", err
	}

	return link, s.gateway.UpdateOrderAfterPaymentLink(ctx, order.Id, link)
}
