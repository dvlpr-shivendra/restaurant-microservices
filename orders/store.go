package main

import (
	"context"
	"errors"
	"fmt"
	pb "restaurant-backend/common/api"
)

var orders = make([]*pb.Order, 0)

type store struct {
}

func NewStore() *store {
	return &store{}
}

func (s *store) Create(ctx context.Context, p *pb.CreateOrderRequest, items []*pb.Item) (string, error) {
	id := fmt.Sprint(len(orders) + 1)

	orders = append(orders, &pb.Order{
		Id:         id,
		CustomerId: p.CustomerId,
		Status:     "pending",
	})

	return id, nil
}

func (s *store) Get(ctx context.Context, id, customerId string) (*pb.Order, error) {
	for _, o := range orders {
		if o.Id == id && o.CustomerId == customerId {
			return o, nil
		}
	}

	return nil, errors.New("order not found")
}

func (s *store) Update(ctx context.Context, id string, order *pb.Order) error {
	for i, o := range orders {
		if o.Id == id {
			orders[i].Status = o.Status
			orders[i].PaymentLink = o.PaymentLink
			return nil
		}
	}

	return nil
}
