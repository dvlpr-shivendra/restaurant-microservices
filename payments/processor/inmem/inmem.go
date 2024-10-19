package inmem

import (
	pb "restaurant-backend/common/api"
)

type inmem struct {
}

func NewProcessor() *inmem {
	return &inmem{}
}

func (p *inmem) CreatePaymentLink(order *pb.Order) (string, error) {
	// TODO: Implement in-memory payment link generation logic
	return "dummy-link", nil
}
