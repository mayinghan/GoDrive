package controller

import (
	"GoDrive/meta"
	"GoDrive/utils"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"time"
)

// FileMetaResponse contains the file meta info struct and error messages
type FileMetaResponse struct {
	FileMeta   meta.FileMeta `json:"meta,omitempty"`
	StatusCode int           `json:"code,omitempty"`
	Msg        string        `json:"msg,omitempty"`
}

//returnErrorResponse creates a file meta reponse object that contains the error
func returnErrorResponse(c int, msg string) (fmr FileMetaResponse) {
	fmr = FileMetaResponse{
		StatusCode: c,
		Msg:        msg,
	}
	return
}

//returnJson writes Json message to front-end
func returnJSON(w http.ResponseWriter, v FileMetaResponse) {
	js, err := json.Marshal(v)
	if err != nil {
		e := fmt.Sprintf("Failed to create json object %s\n", err)
		panic(e)
		// return
	}
	if v.StatusCode != 200 {
		w.WriteHeader(v.StatusCode)
	} else {
		w.WriteHeader(200)
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(js)
}

// UploadHandler handles file upload
func UploadHandler(w http.ResponseWriter, r *http.Request) {
	// if r.Method == "GET" {
	// 	// return upload file page
	// 	page, err := ioutil.ReadFile("./static/view/upload.html")
	// 	if err != nil {
	// 		io.WriteString(w, "internal server error")
	// 		return
	// 	}
	// 	io.WriteString(w, string(page))
	// } else if r.Method == "POST" {
	var fmr FileMetaResponse
	// get a file stream and save into local fs
	// fmt.Printf("%v\n", r)
	file, head, err := r.FormFile("file")
	if err != nil {
		fmr = returnErrorResponse(406, "could not receive file")
		returnJSON(w, fmr)
		return
		//fmt.Printf("Failed to get file %s\n", err.Error())
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
		fmr = returnErrorResponse(500, "failed to create tmp file for io")
		returnJSON(w, fmr)
		return
	}
	defer newFile.Close()

	fileMeta.FileSize, err = io.Copy(newFile, file)
	if err != nil {
		fmr = returnErrorResponse(500, "failed to copy content to tmp file")
		returnJSON(w, fmr)
		//fmt.Printf("Failed to copy content to temp file %s\n", err)
		return
	}
	// move the seek of new file to the start point
	newFile.Seek(0, 0)
	// update file meta hashmap
	fileMeta.FileSha1 = utils.FileSHA1(newFile)
	//debug
	//fmt.Printf("%v\n", fileMeta)
	// upload meta data to DB
	_ = meta.UpdateFileMetaDB(fileMeta)

	// io.WriteString(w, "Upload Successfully")
	// redirect to /success
	fmr = FileMetaResponse{
		FileMeta:   fileMeta,
		StatusCode: 200,
		Msg:        "file successfully uploaded!",
	}
	returnJSON(w, fmr)
	// http.Redirect(w, r, "/file/upload/success", http.StatusFound)
	// }
}

// GetFileMetaHandler gets the meta data of the given file from request.form
func GetFileMetaHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	fmt.Printf("%v\n", r)
	filehash := r.Form["filehash"][0]
	filemeta, err := meta.GetFileMetaDB(filehash)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}

	data, err := json.Marshal(filemeta)
	if err != nil {
		fmr := returnErrorResponse(500, "failed to get file meta from hash")
		returnJSON(w, fmr)
		//w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Write(data)
}

// QueryByBatchHandler : query the last `n` files' info. Query file meta by batch.
func QueryByBatchHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()

	// "limit": how many files the user want to query
	count, _ := strconv.Atoi(r.Form.Get("limit"))
	fMetas := meta.GetLastFileMetas(count)

	// return to client as a JSON
	data, err := json.Marshal(fMetas)
	if err != nil {
		fmr := returnErrorResponse(500, "failed to query file information")
		returnJSON(w, fmr)
		//w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Write(data)
}

// DownloadHandler : download file
func DownloadHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()

	fsha1 := r.Form.Get("filehash")
	metaInfo := meta.GetFileMeta(fsha1)

	f, err := os.Open(metaInfo.Location)
	if err != nil {
		fmr := returnErrorResponse(500, "failed to open file for download")
		returnJSON(w, fmr)
		//w.WriteHeader(http.StatusInternalServerError)
		return
	}

	defer f.Close()

	// read file into RAM. Assuming the file size is not large
	data, err := ioutil.ReadAll(f)
	if err != nil {
		fmr := returnErrorResponse(500, "failed to read file for download")
		returnJSON(w, fmr)
		//w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "appllication/octect-stream")
	w.Header().Set("Content-Disposition", "attatchment; filename=\""+metaInfo.FileName+"\"")
	w.Write(data)
}

// FileUpdateHandler : rename file
func FileUpdateHandler(w http.ResponseWriter, r *http.Request) {
	// only accept post request
	if r.Method != "POST" {
		fmr := returnErrorResponse(405, "failed to update file. post request only")
		returnJSON(w, fmr)
		//w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	r.ParseForm()

	operationType := r.Form.Get("op") // for future use: expand operation type to not only renaming file
	fileSha1 := r.Form.Get("filehash")
	newFileName := r.Form.Get("filename")

	if operationType != "update-name" {
		fmr := returnErrorResponse(403, "failed to update name")
		returnJSON(w, fmr)
		//w.WriteHeader(http.StatusForbidden)
		return
	}

	currFileMeta := meta.GetFileMeta(fileSha1)
	currFileMeta.FileName = newFileName
	meta.UpdateFileMeta(currFileMeta)

	fmr := FileMetaResponse{
		FileMeta:   currFileMeta,
		StatusCode: 200,
		Msg:        "file successfully updated!",
	}
	returnJSON(w, fmr)

	// data, err := json.Marshal(currFileMeta)
	// if err != nil {
	// 	w.WriteHeader(http.StatusInternalServerError)
	// 	return
	// }

	//w.WriteHeader(http.StatusOK)
	//w.Write(data)
}

// FileDeleteHandler : delete the file (soft-delete by using a flag)
func FileDeleteHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	fileSha1 := r.Form.Get("filehash")
	fileMeta, err := meta.GetFileMetaDB(fileSha1)
	if err != nil {
		fmr := returnErrorResponse(500, "failed to delete file from DB")
		returnJSON(w, fmr)
		return
	}
	// remove the file locally
	os.Remove(fileMeta.Location)
	meta.RemoveMeta(fileSha1)
	w.WriteHeader(http.StatusOK)
}
