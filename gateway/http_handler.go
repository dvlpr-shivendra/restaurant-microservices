package main

import (
	"errors"
	"fmt"
	"net/http"
	"restaurant-backend/common"
	pb "restaurant-backend/common/api"
	"restaurant-backend/gateway/gateway"

	"go.opentelemetry.io/otel"
	otelCodes "go.opentelemetry.io/otel/codes"
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

	tracer := otel.Tracer("http")
	ctx, span := tracer.Start(r.Context(), fmt.Sprintf("%s %s", r.Method, r.URL.Path))
	defer span.End()

	order, err := h.gateway.GetOrder(ctx, orderId, customerId)

	rStatus := status.Convert(err)

	if rStatus != nil {
		span.SetStatus(otelCodes.Error, err.Error())
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

	tracer := otel.Tracer("http")
	ctx, span := tracer.Start(r.Context(), fmt.Sprintf("%s %s", r.Method, r.URL.Path))
	defer span.End()

	if err := validateItems(items); err != nil {
		common.WriteError(w, http.StatusBadRequest, err.Error())
		return
	}

	order, err := h.gateway.CreateOrder(ctx, &pb.CreateOrderRequest{
		CustomerId: r.PathValue("customerId"),
		Items:      items,
	})

	rStatus := status.Convert(err)

	if rStatus != nil {
		span.SetStatus(otelCodes.Error, err.Error())
		if rStatus.Code() != codes.InvalidArgument {
			common.WriteError(w, http.StatusBadRequest, rStatus.Message())
			return
		}
	}

	res := &CreateOrderResponse{
		Order:       order,
		RedirectUrl: "http://localhost:8080/payment/success.html?customerID=" + order.CustomerId + "&orderID=" + order.Id,
	}

	common.WriteJSON(w, http.StatusCreated, res)
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
