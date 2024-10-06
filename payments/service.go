package main

import (
	"context"
	pb "restaurant-backend/common/api"
)

type service struct {
	
}

func NewService() *service {
	return &service{}
}

func (s *service) CreatePayment(ctx context.Context, req *pb.Order) (string, error) {
	return "http://localhost", nil
}