package mq

var done chan bool

// StartConsume : consumer needs to keep listening to the message queue to handle requests
func StartConsume(queueName, consumerName string, callback func(msg []byte)) {
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
			callback(d.Body)
		}
	}()

	<-done

	channel.Close()
}

// StopConsume : stop listening to the queue
func StopConsume() {
	done <- true
}
