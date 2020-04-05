package mq

import (
	"GoDrive/config"
	"log"

	"github.com/streadway/amqp"
)

// Publish : publish a message
func Publish(exchange string, routingKey string, message []byte) bool {
	// 1. check channel is ok
	if !initChannel(config.RabbitURL) {
		return false
	}
	// 2. publish msg to mq
	err := channel.Publish(
		exchange,
		routingKey,
		false,
		false,
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        message,
		},
	)
	if err != nil {
		log.Println(err)
		return false
	}

	return true
}
