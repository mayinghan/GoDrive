package config

var (
	// AsyncTransferEnable : whether enable mq when uploading to aws s3
	AsyncTransferEnable = true
	// RabbitURL : rabbitmq url
	RabbitURL = "amqp://guest:guest@127.0.0.1:5672/"
	// TransExchangeName : name of exchange
	TransExchangeName = "upload-server.trans"
	// TransS3QueueName : queue name for s3
	TransS3QueueName = "upload-server.trans.aws"
	// TransS3ErrQueueName : queue name for s3 error
	TransS3ErrQueueName = "upload-server.trans.aws.err"
	// TransS3RoutingKey : routing key for using aws s3
	TransS3RoutingKey = "aws"
)
