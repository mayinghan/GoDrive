package main

import (
	"GoDrive/router"
	"fmt"
)

func main() {
	fmt.Printf("The server running on 127.0.0.1:5050\n")
	router := router.Router()

	err := router.Run(":5050")
	if err != nil {
		fmt.Printf("Failed to start server, err:%s\n", err.Error())
	}

}
