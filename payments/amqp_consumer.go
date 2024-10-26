package main

import (
	"context"
	"log"
	"restaurant-backend/common/broker"

	amqp "github.com/rabbitmq/amqp091-go"

	"encoding/json"

	pb "restaurant-backend/common/api"
)

type consumer struct {
	service PaymentService
}

func NewConsumer(service PaymentService) *consumer {
	return &consumer{service}
}

func (c *consumer) Listen(channel *amqp.Channel) {
	queue, err := channel.QueueDeclare(broker.OrderCreatedEvent, true, false, false, false, nil)

	if err != nil {
		log.Fatal(err)
	}

	messages, err := channel.Consume(queue.Name, "", true, false, false, false, nil)

	if err != nil {
		log.Fatal(err)
	}

	var forever chan struct{}

	go func() {
		for message := range messages {
			log.Printf("Received a message: %v", message)

			order := &pb.Order{}
			err := json.Unmarshal(message.Body, order)

			if err != nil {
				log.Printf("Failed to unmarshal message: %v", err)
				continue
			}

			paymentLink, err := c.service.CreatePayment(context.Background(), order)

			if err != nil {
				log.Printf("Failed to create payment: %v", err)
				continue
			}

			log.Printf("Payment created for order with Id %s : %v", order.Id, paymentLink)
		}
	}()

	<-forever
}
