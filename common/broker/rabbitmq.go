package broker

import (
	"fmt"
	"log"

	amqp "github.com/rabbitmq/amqp091-go"
)

func Connect(user, password, host, port string) (*amqp.Channel, func() error) {
	address := fmt.Sprintf("amqp://%s:%s@%s:%s", user, password, host, port)

	conn, err := amqp.Dial(address)

	if err != nil {
		log.Fatal(err)
	}

	channel, err := conn.Channel()

	if err != nil {
		log.Fatal(err)
	}

	err = channel.ExchangeDeclare(OrderCreatedEvent, "direct", true, false, false, false, nil)

	if err != nil {
		log.Fatal(err)
	}

	err = channel.ExchangeDeclare(OrderPaidEvent, "fanout", true, false, false, false, nil)

	if err != nil {
		log.Fatal(err)
	}

	return channel, conn.Close
}
