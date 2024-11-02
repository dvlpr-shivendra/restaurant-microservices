package main

import (
	pb "restaurant-backend/common/api"
)

type CreateOrderResponse struct {
	Order       *pb.Order `"json:"order"`
	RedirectUrl string    `"json:"redirectUrl"`
}
