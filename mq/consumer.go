package mq

import "log"

var done chan bool

// StartConsume : consumer needs to keep listening to the message queue to handle requests
func StartConsume(queueName, consumerName string, callback func(msg []byte) bool) {
	// 1. get message channel
	msgs, err := channel.Consume(
		queueName,
		consumerName,
		true,
		false, // multi consumer on one queue
		false,
		false,
		nil,
	)
	if err != nil {
		panic(err)
	}
	// 2. looping the channel to get info
	done = make(chan bool)
	go func() {
		for d := range msgs {
			// 3. callback to handle msg
			processErr := callback(d.Body)
			if processErr {
				// @TODO: put the task to the err queue to recover the error
				log.Println("task failed.....")
			}
		}
	}()

	<-done

	channel.Close()
}

// StopConsume : stop listening to the queue
func StopConsume() {
	done <- true
}
