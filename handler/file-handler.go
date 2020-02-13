package handler

import (
	"GoDrive/meta"
	"GoDrive/utils"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

// UploadHandler handels file upload
func UploadHandler(c *gin.Context) {
	head, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":  406,
			"msg":   "Could not receive file.",
			"error": err.Error(),
		})
		return
	}

	fileMeta := meta.FileMeta{
		FileName: head.Filename,
		Location: "C://Users/liuwi/Desktop/tmp/" + head.Filename,
		UploadAt: time.Now().Format("2006-01-02 15:04:05"),
	}

	err = c.SaveUploadedFile(head, fileMeta.Location)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":  500,
			"msg":   "Internal Server Error: Failed to save file to the DB.",
			"error": err.Error(),
		})
		return
	}

	newFile, err := os.Open(fileMeta.Location)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":  500,
			"msg":   "Internal Server Error: Failed to save file to the DB.",
			"error": err.Error(),
		})
		return
	}
	defer newFile.Close()
	newFileInfo, err := newFile.Stat()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":  500,
			"msg":   "Internal Server Error: Failed to save file to the DB.",
			"error": err.Error(),
		})
		return
	}

	// update file meta hashmap
	fileMeta.FileSize = newFileInfo.Size()
	fileMeta.FileSha1 = utils.FileSHA1(newFile)
	// upload meta data to DB
	_ = meta.UpdateFileMetaDB(fileMeta)

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"msg":  "File successfully uploaded!",
		"data": struct {
			FileMeta *meta.FileMeta `json:"meta"`
		}{
			FileMeta: &fileMeta,
		},
	})
	return
}

// UploadHandler handles file upload
/* func UploadHandler(w http.ResponseWriter, r *http.Request) {
	// if r.Method == "GET" {
	// 	// return upload file page
	// 	page, err := ioutil.ReadFile("./static/view/upload.html")
	// 	if err != nil {
	// 		io.WriteString(w, "internal server error")
	// 		return
	// 	}
	// 	io.WriteString(w, string(page))
	// } else if r.Method == "POST" {
	var fmr fileMetaResponse
	// get a file stream and save into local fs
	// fmt.Printf("%v\n", r)
	file, head, err := r.FormFile("file")
	if err != nil {
		fmr = fileErrorResponse(406, "could not receive file")
		returnFileRespJSON(w, fmr)
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
		fmr = fileErrorResponse(500, "failed to create tmp file for io")
		returnFileRespJSON(w, fmr)
		return
	}
	defer newFile.Close()

	fileMeta.FileSize, err = io.Copy(newFile, file)
	if err != nil {
		fmr = fileErrorResponse(500, "failed to copy content to tmp file")
		returnFileRespJSON(w, fmr)
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
	fmr = fileMetaResponse{
		FileMeta:   &fileMeta,
		StatusCode: 200,
		Msg:        "file successfully uploaded!",
	}
	returnFileRespJSON(w, fmr)
	// http.Redirect(w, r, "/file/upload/success", http.StatusFound)
	// }
} */

// GetFileMetaHandler gets the meta data of the given file from request.form
func GetFileMetaHandler(c *gin.Context) {
	var filehash string
	if err := c.ShouldBindJSON(&filehash); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 500,
			"msg":  err.Error(),
		})
		panic(err)
	}

	// c.Request.ParseForm()
	// fmt.Printf("%v\n", c.Request)
	// filehash := c.Request.Form["filehash"][0]
	filemeta, err := meta.GetFileMetaDB(filehash)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":  500,
			"msg":   "Internal Server Error: Failed to retrieve file meta.",
			"error": err.Error(),
		})
		return
	}

	data, err := json.Marshal(filemeta)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":  500,
			"msg":   "Internal Server Error: Failed to retrieve file meta.",
			"error": err.Error(),
		})
		return
	}
	c.Writer.Write(data)
}

// GetFileMetaHandler gets the meta data of the given file from request.form
/* func GetFileMetaHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	fmt.Printf("%v\n", r)
	filehash := r.Form["filehash"][0]
	filemeta, err := meta.GetFileMetaDB(filehash)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}

	data, err := json.Marshal(filemeta)
	if err != nil {
		fmr := fileErrorResponse(500, "failed to get file meta from hash")
		returnFileRespJSON(w, fmr)
		//w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Write(data)
} */

// QueryByBatchHandler : query the last `n` files' info. Query file meta by batch.
func QueryByBatchHandler(c *gin.Context) {
	var lim string
	if err := c.ShouldBindJSON(&lim); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 500,
			"msg":  err.Error(),
		})
		panic(err)
	}

	// "limit": how many files the user want to query
	count, _ := strconv.Atoi(lim)
	fMetas := meta.GetLastFileMetas(count)

	// return to client as a JSON
	data, err := json.Marshal(fMetas)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":  500,
			"msg":   "Internal Server Error: Failed to query file information.",
			"error": err.Error(),
		})
		return
	}
	c.Writer.Write(data)
}

// QueryByBatchHandler : query the last `n` files' info. Query file meta by batch.
/* func QueryByBatchHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()

	// "limit": how many files the user want to query
	count, _ := strconv.Atoi(r.Form.Get("limit"))
	fMetas := meta.GetLastFileMetas(count)

	// return to client as a JSON
	data, err := json.Marshal(fMetas)
	if err != nil {
		fmr := fileErrorResponse(500, "failed to query file information")
		returnFileRespJSON(w, fmr)
		//w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Write(data)
} */

// DownloadHandler : download file
func DownloadHandler(c *gin.Context) {
	var fileSha1 string
	if err := c.ShouldBindJSON(&fileSha1); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 500,
			"msg":  err.Error(),
		})
		panic(err)
	}
	metaInfo := meta.GetFileMeta(fileSha1)

	f, err := os.Open(metaInfo.Location)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":     500,
			"msg":      "Internal Server Error: Failed to open file for download.",
			"error":    err.Error(),
			"location": metaInfo.Location,
		})
		return
	}
	defer f.Close()

	// read file into RAM. Assuming the file size is not large
	data, err := ioutil.ReadAll(f)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":  500,
			"msg":   "Internal Server Error: Failed to read file for download.",
			"error": err.Error(),
		})
		return
	}

	c.Writer.Header().Set("Content-Type", "appllication/octect-stream")
	c.Writer.Header().Set("Content-Disposition", "attatchment; filename=\""+metaInfo.FileName+"\"")
	c.Writer.Write(data)
}

// DownloadHandler : download file
/*func DownloadHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()

	fsha1 := r.Form.Get("filehash")
	metaInfo := meta.GetFileMeta(fsha1)

	f, err := os.Open(metaInfo.Location)
	if err != nil {
		fmr := fileErrorResponse(500, "failed to open file for download")
		returnFileRespJSON(w, fmr)
		//w.WriteHeader(http.StatusInternalServerError)
		return
	}

	defer f.Close()

	// read file into RAM. Assuming the file size is not large
	data, err := ioutil.ReadAll(f)
	if err != nil {
		fmr := fileErrorResponse(500, "failed to read file for download")
		returnFileRespJSON(w, fmr)
		//w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "appllication/octect-stream")
	w.Header().Set("Content-Disposition", "attatchment; filename=\""+metaInfo.FileName+"\"")
	w.Write(data)
} */

// FileUpdateHandler : renames file
func FileUpdateHandler(c *gin.Context) {
	// only accept post request
	if c.Request.Method != "POST" {
		c.JSON(http.StatusMethodNotAllowed, gin.H{
			"code": 405,
			"msg":  "Status Method Not Allowed: Failed to update file - POST request only.",
		})
		return
	}
	c.Request.ParseForm()
	operationType := c.Request.Form.Get("op") // for future use: expand operation type to not only renaming file
	fileSha1 := c.Request.Form.Get("filehash")
	newFileName := c.Request.Form.Get("filename")

	if operationType != "update-name" {
		c.JSON(http.StatusForbidden, gin.H{
			"code": 403,
			"msg":  "Status Forbidden: Failed to update file name.",
		})
		return
	}

	currFileMeta := meta.GetFileMeta(fileSha1)
	currFileMeta.FileName = newFileName
	meta.UpdateFileMeta(currFileMeta)

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"msg":  "File successfully updated!",
		"data": struct {
			FileMeta *meta.FileMeta `json:"meta"`
		}{
			FileMeta: &currFileMeta,
		},
	})
	return
}

// FileUpdateHandler : rename file
/*func FileUpdateHandler(w http.ResponseWriter, r *http.Request) {
	// only accept post request
	if r.Method != "POST" {
		fmr := fileErrorResponse(405, "failed to update file. post request only")
		returnFileRespJSON(w, fmr)
		//w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	r.ParseForm()

	operationType := r.Form.Get("op") // for future use: expand operation type to not only renaming file
	fileSha1 := r.Form.Get("filehash")
	newFileName := r.Form.Get("filename")

	if operationType != "update-name" {
		fmr := fileErrorResponse(403, "failed to update name")
		returnFileRespJSON(w, fmr)
		//w.WriteHeader(http.StatusForbidden)
		return
	}

	currFileMeta := meta.GetFileMeta(fileSha1)
	currFileMeta.FileName = newFileName
	meta.UpdateFileMeta(currFileMeta)

	fmr := fileMetaResponse{
		FileMeta:   &currFileMeta,
		StatusCode: 200,
		Msg:        "file successfully updated!",
	}
	returnFileRespJSON(w, fmr)

	// data, err := json.Marshal(currFileMeta)
	// if err != nil {
	// 	w.WriteHeader(http.StatusInternalServerError)
	// 	return
	// }

	//w.WriteHeader(http.StatusOK)
	//w.Write(data)
}
*/

// FileDeleteHandler : delete the file (soft-delete by using a flag)
/*func FileDeleteHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	fileSha1 := r.Form.Get("filehash")
	fileMeta, err := meta.GetFileMetaDB(fileSha1)
	if err != nil {
		fmr := fileErrorResponse(500, "failed to delete file from DB")
		returnFileRespJSON(w, fmr)
		return
	}
	// remove the file locally
	os.Remove(fileMeta.Location)
	meta.RemoveMeta(fileSha1)
	w.WriteHeader(http.StatusOK)
} */

// FileDeleteHandler : delete the file (soft-delete by using a flag)
func FileDeleteHandler(c *gin.Context) {
	var fileSha1 string
	if err := c.ShouldBindJSON(&fileSha1); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 500,
			"msg":  err.Error(),
		})
		panic(err)
	}
	fileMeta, err := meta.GetFileMetaDB(fileSha1)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":  500,
			"msg":   "Internal server error: Failed to delete file from the database.",
			"error": err.Error(),
		})
		return
	}
	removeFromDB := meta.RemoveMetaDB(fileSha1)
	if !removeFromDB {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":  500,
			"msg":   "Internal server error: Failed to delete file from the database.",
			"error": err.Error(),
		})
		return
	}
	os.Remove(fileMeta.Location)
	meta.RemoveMeta(fileSha1)
	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"msg":  "File successfully deleted!",
	})
	return
}
