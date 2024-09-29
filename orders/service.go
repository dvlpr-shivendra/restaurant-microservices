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

func (s *service) CreateOrder(ctx context.Context) error {
	return s.store.Create(ctx)
}

func (s *service) ValidateOrder(ctx context.Context, req *pb.CreateOrderRequest) error {
	if len(req.Items) == 0 {
		return common.ErrNoItems
	}

	mergedItems := mergeItemsQuantities(req.Items)

	log.Print(mergedItems)

	return nil
}

func mergeItemsQuantities(items []*pb.ItemsWithQuantity) []*pb.ItemsWithQuantity {
	merged := make([]*pb.ItemsWithQuantity, 0)

	for _, item := range items {
		found := false
		for _, mergedItem := range merged {
			if mergedItem.ID == item.ID {
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
