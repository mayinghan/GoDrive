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
	router.HandleFunc("/api/file/upload", controller.UploadHandler)
	router.HandleFunc("/api/file/upload/success", controller.UploadSuccessHandler)
	router.HandleFunc("/api/file/meta", controller.GetFileMetaHandler)
	router.HandleFunc("/api/file/query", controller.QueryByBatchHandler)
	router.HandleFunc("/api/file/download", controller.DownloadHandler)
	router.HandleFunc("/api/file/update", controller.FileUpdateHandler)
	router.HandleFunc("/api/file/delete", controller.FileDeleteHandler)
	log.Fatal(http.ListenAndServe(":5050", router))
}
