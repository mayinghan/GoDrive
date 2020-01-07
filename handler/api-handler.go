package handler

import (
	"io"
	"io/ioutil"
	"net/http"
)

// Handling file upload
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
	}
}
