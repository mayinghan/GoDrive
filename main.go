package main

import (
	"GoDrive/controller"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func main() {
	fmt.Printf("The server running on 127.0.0.1:5050")
	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/file/upload", controller.UploadHandler)
	router.HandleFunc("/file/upload/success", controller.UploadSuccessHandler)
	router.HandleFunc("/file/meta", controller.GetFileMetaHandler)
	router.HandleFunc("/file/query", controller.QueryByBatchHandler)
	router.HandleFunc("/file/download", controller.DownloadHandler)
	router.HandleFunc("/file/update", controller.FileUpdateHandler)
	router.HandleFunc("/file/delete", controller.FileDeleteHandler)
	log.Fatal(http.ListenAndServe(":5050", router))
}
