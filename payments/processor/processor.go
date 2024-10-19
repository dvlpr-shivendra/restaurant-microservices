package processor

import (
	pb "restaurant-backend/common/api"
)

type PaymentProcessor interface {
	CreatePaymentLink(*pb.Order) (string, error)
}
