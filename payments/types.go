package main

import (
	"context"
	pb "restaurant-backend/common/api"
)

type PaymentService interface {
	CreatePayment(context.Context, *pb.Order) (string, error)
}
