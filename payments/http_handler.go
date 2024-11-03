package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"restaurant-backend/common/broker"
	"time"

	pb "restaurant-backend/common/api"

	amqp "github.com/rabbitmq/amqp091-go"
	"go.opentelemetry.io/otel"
)

type PaymentHTTPHandler struct {
	channel *amqp.Channel
}

func NewPaymentHTTPHandler(channel *amqp.Channel) *PaymentHTTPHandler {
	return &PaymentHTTPHandler{channel}
}

func (h *PaymentHTTPHandler) registerRoutes(router *http.ServeMux) {
	router.HandleFunc("/webhook", h.handleCheckoutWebhook)
}

func (h *PaymentHTTPHandler) handleCheckoutWebhook(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)

	defer cancel()

	order := &pb.Order{
		Id:          "1",
		CustomerId:  "1",
		Status:      "paid",
		PaymentLink: "",
	}

	marshallOrder, _ := json.Marshal(order)

	tr := otel.Tracer("amqp")
	amqpContext, messageSpan := tr.Start(ctx, fmt.Sprintf("AMQP - publish - %s", broker.OrderPaidEvent))
	defer messageSpan.End()

	headers := broker.InjectAMQPHeaders(amqpContext)

	h.channel.PublishWithContext(ctx, broker.OrderPaidEvent, "", false, false, amqp.Publishing{
		ContentType:  "application/json",
		Body:         marshallOrder,
		DeliveryMode: amqp.Persistent,
		Headers:      headers,
	})

	log.Println("Message published order.paid")

	w.WriteHeader(http.StatusOK)
}
