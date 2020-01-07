package controller

import (
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
)

// UploadHandler handles file upload
func UploadHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		// return upload file page
		page, err := ioutil.ReadFile("./static/view/upload.html")
		if err != nil {
			io.WriteString(w, "internal server error")
			return
		}
		io.WriteString(w, string(page))
	} else if r.Method == "POST" {
		// get a file stream and save into local fs
		file, head, err := r.FormFile("file")
		if err != nil {
			fmt.Printf("Failed to get file %s\n", err.Error())
		}
		// make sure the file handler is closed
		defer file.Close()

		newFile, err := os.Create("/tmp/" + head.Filename)
		if err != nil {
			fmt.Printf("Failed to create tmp file %s\n", err.Error())
			return
		}
		defer newFile.Close()

		_, err = io.Copy(newFile, file)
		if err != nil {
			fmt.Printf("Failed to copy content to temp file %s\n", err)
			return
		}

		//io.WriteString(w, "Upload Successfully")
		// redirect to /success
		http.Redirect(w, r, "/file/upload/success", http.StatusFound)
	}
}

// UploadSuccessHandler will return content when upload successfully
func UploadSuccessHandler(w http.ResponseWriter, r *http.Request) {
	io.WriteString(w, "Upload Successfully")
}
