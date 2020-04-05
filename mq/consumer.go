package mq

var done chan bool

// consumer needs to keep listening to the message queue to handle requests
func startConsume(queueName, consumerName string, callback func(msg []byte) bool) {
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
			result := callback(d.Body)
			if !result {
				panic(err)
			}
		}
	}()
	// 3. callback to handle msg
	<-done

	channel.Close()
}
