package main

import (
	"context"
	pb "restaurant-backend/common/api"
)

type PaymentsService interface {
	CreatePayment(context.Context, *pb.Order) (string, error)
}
