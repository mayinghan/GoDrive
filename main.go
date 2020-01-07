package main

import (
	"GoDisk/controller"
	"fmt"
	"log"
	"net/http"
)

func listenHandler() {
	fmt.Printf("Server running on 127.0.0.1:8080")
}

func main() {
	http.HandleFunc("/file/upload", controller.UploadHandler)
	http.HandleFunc("/file/upload/success", controller.UploadSuccessHandler)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
