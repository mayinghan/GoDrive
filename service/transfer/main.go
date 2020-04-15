package main

import (
	"GoDrive/config"
	"GoDrive/mq"
	"encoding/json"
	"io/ioutil"
	"log"
	"strconv"
	"strings"
)

// TransferCallback : handle file transfer task
func TransferCallback(msg []byte) bool {
	// 1. parse msg
	pubData := mq.TransferData{}
	err := json.Unmarshal(msg, &pubData)
	if err != nil {
		log.Println(err.Error())
		return false
	}

	// create file pointer based on tmp store location
	if pubData.IsLarge {
		// for large file, open the directory
		fileList, err := ioutil.ReadDir(pubData.CurLocation)
		if err != nil {
			log.Println(err.Error())
			return false
		}

		for _, file := range fileList {
			fileName := file.Name()
			chunkIdx, _ := strconv.Atoi(strings.Split(fileName, "_")[1])

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
