package main

import (
	"errors"
	"net/http"
	"restaurant-backend/common"
	pb "restaurant-backend/common/api"
	"restaurant-backend/gateway/gateway"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type handler struct {
	gateway gateway.OrdersGateway
}

func NewHandler(gateway gateway.OrdersGateway) *handler {
	return &handler{gateway}
}

func (h *handler) registerRoutes(mux *http.ServeMux) {
	// Serve public directory
	mux.Handle("/", http.FileServer(http.Dir("public")))
	mux.HandleFunc("POST /api/customers/{customerId}/orders", h.handleCreateOrder)
	mux.HandleFunc("GET /api/customers/{customerId}/orders/{orderId}", h.handleGetOrder)
}

func (h *handler) handleGetOrder(w http.ResponseWriter, r *http.Request) {
	customerId := r.PathValue("customerId")
	orderId := r.PathValue("orderId")

	order, err := h.gateway.GetOrder(r.Context(), orderId, customerId)

	rStatus := status.Convert(err)

	if rStatus != nil {
		if rStatus.Code() != codes.InvalidArgument {
			common.WriteError(w, http.StatusBadRequest, rStatus.Message())
			return
		}
	}

	common.WriteJSON(w, http.StatusOK, order)
}

func (h *handler) handleCreateOrder(w http.ResponseWriter, r *http.Request) {

	var items []*pb.ItemsWithQuantity

	if err := common.ReadJSON(r, &items); err != nil {
		common.WriteError(w, http.StatusBadRequest, err.Error())
		return
	}

	if err := validateItems(items); err != nil {
		common.WriteError(w, http.StatusBadRequest, err.Error())
		return
	}

	order, err := h.gateway.CreateOrder(r.Context(), &pb.CreateOrderRequest{
		CustomerId: r.PathValue("customerId"),
		Items:      items,
	})

	rStatus := status.Convert(err)

	if rStatus != nil {
		if rStatus.Code() != codes.InvalidArgument {
			common.WriteError(w, http.StatusBadRequest, rStatus.Message())
			return
		}
	}

	common.WriteJSON(w, http.StatusCreated, order)
}

func validateItems(items []*pb.ItemsWithQuantity) error {
	if len(items) == 0 {
		return common.ErrNoItems
	}

	for _, item := range items {
		if item.Id == "" {
			return errors.New("items Id is required")
		}

		if item.Quantity <= 0 {
			return errors.New("items must have a valid quantity")
		}
	}

	return nil
}
