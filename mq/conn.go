package mq

import (
	"GoDrive/config"
	"log"

	"github.com/streadway/amqp"
)

var conn *amqp.Connection
var channel *amqp.Channel

// receive error info when error happened
var notifyClose chan *amqp.Error

func UpdateMQHost(host string) {
	config.RabbitURL = host
}

func Init() {
	if !config.AsyncTransferEnable {
		return
	}

	if initChannel(config.RabbitURL) {
		channel.NotifyClose(notifyClose)
	}

	// if disconnected, try reconnect automatically
	go func() {
		for {
			select {
			case msg := <-notifyClose:
				conn = nil
				channel = nil
				log.Printf(" MQ closed due to %+v\n", msg)
				initChannel(config.RabbitURL)
			}
		}
	}()
}

func initChannel(url string) bool {
	// 1. check if channel is established
	if channel != nil {
		return true
	}

	// 2. get mq connection
	conn, err := amqp.Dial(url)
	if err != nil {
		log.Println(err)
		return false
	}

	// 3. open channel
	channel, err = conn.Channel()
	if err != nil {
		log.Println(err)
		return false
	}

	return true
}
