package main

import (
	"GoDrive/config"
	"GoDrive/mq"
	"encoding/json"
	"log"
	"os"
)

func TransferCallback(msg []byte) {
	// 1. parse msg
	pubData := mq.TransferData{}
	err := json.Unmarshal(msg, &pubData)
	if err != nil {
		panic(err)
	}

	// create file pointer based on tmp store location
	if pubData.IsLarge {
		f, err := os.Open(pubData.CurLocation)
		if err != nil {
			panic(err)
		}

	}

}

func main() {
	log.Println("start transferring tasks")
	mq.StartConsume(
		config.TransS3QueueName,
		"transfer_s3",
		TransferCallback,
	)
}
