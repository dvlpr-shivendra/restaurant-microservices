package main

import (
	"context"
	"log"

	"restaurant-backend/common"
	pb "restaurant-backend/common/api"
)

type service struct {
	store OrdersStore
}

func NewService(store OrdersStore) *service {
	return &service{store}
}

func (s *service) GetOrder(ctx context.Context, req *pb.GetOrderRequest) (*pb.Order, error) {
	return s.store.Get(ctx, req.OrderId, req.CustomerId)
}

func (s *service) CreateOrder(ctx context.Context, req *pb.CreateOrderRequest, items []*pb.Item) (*pb.Order, error) {

	id, err := s.store.Create(ctx, req, items)

	if err != nil {
		return nil, err
	}

	order := &pb.Order{
		Id:         id,
		CustomerId: req.CustomerId,
		Items:      items,
		Status:     "pending",
		PaymentLink: "",
	}

	return order, nil
}

func (s *service) UpdateOrder(ctx context.Context, order *pb.Order) (*pb.Order, error) {
	err := s.store.Update(ctx, order.Id, order)

	if err != nil {
		return nil, err
	}

	return order, nil
}

func (s *service) ValidateOrder(ctx context.Context, req *pb.CreateOrderRequest) ([]*pb.Item, error) {
	if len(req.Items) == 0 {
		return nil, common.ErrNoItems
	}

	mergedItems := mergeItemsQuantities(req.Items)

	log.Printf("Merged items: %v", mergedItems)

	// temp
	var items []*pb.Item

	for _, item := range mergedItems {
		items = append(items, &pb.Item{
			Id:       item.Id,
			Quantity: item.Quantity,
			PriceId:  "1",
		})
	}

	return items, nil
}

func mergeItemsQuantities(items []*pb.ItemsWithQuantity) []*pb.ItemsWithQuantity {
	merged := make([]*pb.ItemsWithQuantity, 0)

	for _, item := range items {
		found := false
		for _, mergedItem := range merged {
			if mergedItem.Id == item.Id {
				mergedItem.Quantity += item.Quantity
				found = true
				break
			}
		}

		if !found {
			merged = append(merged, item)
		}

	}
	return merged
}
