package handler

import (
	"fmt"
	"net/http"
	"os"
	"path"
	"strings"

	"github.com/gin-gonic/gin"
)

// GetFileChunk receives the file chunks
func GetFileChunk(c *gin.Context) {
	_, exist := c.Get("username")
	if !exist {
		fmt.Println("Get username from context failed")
		c.JSON(200, gin.H{
			"code": 1,
			"msg":  "Auth failed",
		})
		return
	}

	chunk, err := c.FormFile("chunk")
	if err != nil {
		fmt.Println(err.Error())
		c.JSON(200, gin.H{
			"code": 1,
			"msg":  "server parse chunk file failed",
		})
		return
	}
	uploadID := c.PostForm("uploadId")
	chunkID := c.PostForm("chunkId")
	filename := c.PostForm("filename")
	filehash := c.PostForm("filehash")

	fileuser := strings.Split(uploadID, "-")[0]
	fmt.Println(fileuser)
	// if username != fileuser {
	// 	fmt.Println("Authentication error, uploadId belonger is not current user")
	// 	return
	// }

	fmt.Printf("filename : %s\nuploadId: %s\n", filename, uploadID)
	tempPath := "/tmp/" + filehash + "/" + chunkID
	os.MkdirAll(path.Dir(tempPath), 0744)

	if err := c.SaveUploadedFile(chunk, tempPath); err != nil {
		c.String(http.StatusBadRequest, "failed to save chunk")
		return
	}

	c.JSON(200, gin.H{
		"code": 0,
		"msg":  chunkID + " upload suc",
	})
	return
}
