package controller

import (
	"GoDisk/meta"
	"GoDisk/utils"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"time"
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

		fileMeta := meta.FileMeta{
			FileName: head.Filename,
			Location: "/tmp/" + head.Filename,
			UploadAt: time.Now().Format("2006-01-02 15:04:05"),
		}
		newFile, err := os.Create(fileMeta.Location)
		if err != nil {
			fmt.Printf("Failed to create tmp file %s\n", err.Error())
			return
		}
		defer newFile.Close()

		fileMeta.FileSize, err = io.Copy(newFile, file)
		if err != nil {
			fmt.Printf("Failed to copy content to temp file %s\n", err)
			return
		}
		// move the seek of new file to the start point
		newFile.Seek(0, 0)
		// update file meta hashmap
		fileMeta.FileSha1 = utils.FileSHA1(newFile)
		meta.UpdateFileMeta(fileMeta)
		// io.WriteString(w, "Upload Successfully")
		// redirect to /success
		http.Redirect(w, r, "/file/upload/success", http.StatusFound)
	}
}

// UploadSuccessHandler will return content when upload successfully
func UploadSuccessHandler(w http.ResponseWriter, r *http.Request) {
	io.WriteString(w, "Upload Successfully")
}
