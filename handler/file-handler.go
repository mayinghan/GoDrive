package handler

import (
	"GoDrive/aws"
	"GoDrive/config"
	"GoDrive/db"
	"GoDrive/meta"
	"GoDrive/utils"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"runtime"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

const goos string = runtime.GOOS

// UploadHandler handles file upload
func UploadHandler(c *gin.Context) {
	head, err := c.FormFile("file")
	clientHash := c.PostForm("filehash")

	if err != nil {
		panic(err.Error())
	}

	var basepath string = config.WholeFileStoreLocation
	fileMeta := meta.FileMeta{
		FileName: head.Filename,
		Location: basepath + head.Filename,
		UploadAt: time.Now().Format("2006-01-02 15:04:05"),
	}

	err = c.SaveUploadedFile(head, fileMeta.Location)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":  1,
			"msg":   "Internal Server Error: Failed to save file to the DB.",
			"error": err.Error(),
		})
		return
	}

	newFile, err := os.Open(fileMeta.Location)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":  1,
			"msg":   "Internal Server Error: Failed to save file to the DB.",
			"error": err.Error(),
		})
		return
	}

	defer newFile.Close()
	newFileInfo, err := newFile.Stat()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":  1,
			"msg":   "Internal Server Error: Failed to save file to the DB.",
			"error": err.Error(),
		})
		return
	}

	// update file meta hashmap
	fileMeta.FileSize = newFileInfo.Size()
	fileMeta.FileMD5 = utils.FileMD5(newFile)

	// integrity checking
	if fileMeta.FileMD5 != clientHash {
		// integrity checking failed
		c.JSON(200, gin.H{
			"code": 1,
			"msg":  "Server integrity check failed, please try upload again later",
		})
		return
	}

	// getting username
	username, exist := c.Get("username")
	if !exist {
		fmt.Printf("Failed to find username.")
	}

	// upload meta data to databases
	uploadDB := meta.UpdateFileMetaDB(fileMeta, username.(string))

	if !uploadDB {
		c.JSON(200, gin.H{
			"code": 1,
			"msg":  "Internal Server Error: Failed to save metadata to the DB.",
		})
		return
	}

	uploadAWS, err := aws.UploadToAWS(fileMeta.Location, fileMeta.FileMD5)
	if !uploadAWS {
		c.JSON(200, gin.H{
			"code":  1,
			"msg":   "Internal Server Error: Failed to save file to the AWS.",
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"msg":  "File successfully uploaded!",
		"data": struct {
			FileMeta *meta.FileMeta `json:"meta"`
		}{
			FileMeta: &fileMeta,
		},
	})
	return
}

// GetFileMetaHandler gets the meta data of the given file from request.form
func GetFileMetaHandler(c *gin.Context) {
	var filehash string
	if err := c.ShouldBindJSON(&filehash); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 1,
			"msg":  err.Error(),
		})
		panic(err)
	}

	filemeta, err := meta.GetFileMetaDB(filehash)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":  1,
			"msg":   "Internal Server Error: Failed to retrieve file meta.",
			"error": err.Error(),
		})
		return
	}

	data, err := json.Marshal(filemeta)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":  1,
			"msg":   "Internal Server Error: Failed to retrieve file meta.",
			"error": err.Error(),
		})
		return
	}
	c.Writer.Write(data)
}

// QueryByBatchHandler : query the last `n` files' info. Query file meta by batch.
func QueryByBatchHandler(c *gin.Context) {
	var lim string
	if err := c.ShouldBindJSON(&lim); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 1,
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
			"code":  1,
			"msg":   "Internal Server Error: Failed to query file information.",
			"error": err.Error(),
		})
		return
	}
	c.Writer.Write(data)
}

// DownloadHandler : download file
func DownloadHandler(c *gin.Context) {
	filehash := c.Query("filehash")
	metaInfo := meta.GetFileMeta(filehash)

	f, err := os.Open(metaInfo.Location)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":     1,
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
			"code":  1,
			"msg":   "Internal Server Error: Failed to read file for download.",
			"error": err.Error(),
		})
		return
	}

	c.Writer.Header().Set("Content-Type", "appllication/octect-stream")
	c.Writer.Header().Set("Content-Disposition", "attatchment; filename=\""+metaInfo.FileName+"\"")
	c.Writer.Write(data)
}

// FileUpdateHandler : renames file
func FileUpdateHandler(c *gin.Context) {
	// only accept post request
	if c.Request.Method != "POST" {
		c.JSON(http.StatusMethodNotAllowed, gin.H{
			"code": 1,
			"msg":  "Status Method Not Allowed: Failed to update file - POST request only.",
		})
		return
	}
	c.Request.ParseForm()
	operationType := c.Request.Form.Get("op") // for future use: expand operation type to not only renaming file
	filehash := c.Request.Form.Get("filehash")
	newFileName := c.Request.Form.Get("filename")

	if operationType != "update-name" {
		c.JSON(http.StatusForbidden, gin.H{
			"code": 1,
			"msg":  "Status Forbidden: Failed to update file name.",
		})
		return
	}

	currFileMeta := meta.GetFileMeta(filehash)
	currFileMeta.FileName = newFileName
	meta.UpdateFileMeta(currFileMeta)

	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"msg":  "File successfully updated!",
		"data": struct {
			FileMeta *meta.FileMeta `json:"meta"`
		}{
			FileMeta: &currFileMeta,
		},
	})
	return
}

// FileDeleteHandler : delete the file (soft-delete by using a flag)
func FileDeleteHandler(c *gin.Context) {
	var filehash string = c.Query("filehash")
	fileMeta, err := meta.GetFileMetaDB(filehash)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":  1,
			"msg":   "Internal server error: Failed to delete file from the database.",
			"error": err.Error(),
		})
		return
	}

	// getting username
	username, exist := c.Get("username")
	if !exist {
		fmt.Printf("Failed to find username.")
	}

	removeFromDB, delFile := meta.RemoveMetaDB(username.(string), filehash)
	if !removeFromDB {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code": 1,
			"msg":  "Internal server error: Failed to delete file from the databases.",
		})
		return
	}
	if delFile {
		meta.RemoveMeta(filehash)
	}
	os.Remove(fileMeta.Location)
	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"msg":  "File successfully deleted!",
	})
	return
}

// InstantUpload : check if the file is already in the database by comparing the hash.
// If so, then instant upload is triggered
func InstantUpload(c *gin.Context) {
	fileHash := c.Query("filehash")
	fileHash = strings.TrimRight(fileHash, "\n")
	if fileHash == "" {
		c.JSON(200, gin.H{
			"code": 1,
			"msg":  "Empty filehash received, please wait until the file finish preprocess",
		})
		return
	}
	dup, err := db.IsFileUploaded(fileHash)
	if err != nil {
		panic(err.Error())
	}
	// if the file is already uploaded before
	if dup {
		// update the value `copies` in the database
		err := db.UpdateCopies(fileHash)
		if err != nil {
			panic(err.Error())
		}
		// update successfully
		c.JSON(200, gin.H{
			"code": 0,
			"msg":  "Duplicate file detected",
			"data": gin.H{
				"shouldUpload": false,
			},
		})
		return
	}
	// no duplicated file detected
	c.JSON(200, gin.H{
		"code": 0,
		"msg":  "No dup file detected",
		"data": gin.H{
			"shouldUpload": true,
		},
	})

}
