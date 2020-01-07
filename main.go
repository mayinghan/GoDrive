package main

import (
	"GoDisk/handler"
	"fmt"
	"log"
	"net/http"
)

func listenHandler() {
	fmt.Printf("Server running on 127.0.0.1:8080")
}

func main() {
	http.HandleFunc("/file/upload", handler.UploadHandler)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
