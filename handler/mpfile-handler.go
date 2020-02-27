package handler

import (
	"GoDrive/config"
	"GoDrive/utils"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/gin-gonic/gin"
)

// GetFileChunk receives the file chunks
func GetFileChunk(c *gin.Context) {
	username, exist := c.Get("username")
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
	log.Println("current user's username", fileuser)

	if username != fileuser {
		log.Println("Authentication error, uploadId belonger is not current user")
		return
	}

	chunkRootPath := config.ChunkFileStoreDirectory
	log.Printf("filename : %s\nuploadId: %s\n", filename, uploadID)
	tempPath := chunkRootPath + filehash + "/" + chunkID
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

// CheckIntegrity checks the file hash again to make sure the file is not modified
func CheckIntegrity(c *gin.Context) {
	type body struct {
		Filehash    string `json:"filehash"`
		Filename    string `json:"filename"`
		UploadID    string `json:"uploadId"`
		ChunkLength int    `json:"chunkLength,string"`
	}

	var b body
	if err := c.ShouldBindJSON(&b); err != nil {
		c.JSON(200, gin.H{
			"code": 1,
			"msg":  err.Error(),
		})
		panic(err)
	}

	log.Printf("Checking integrity.. Filename: %s, Filehash: %s", b.Filename, b.Filehash)

	mdhash := new(utils.MD5Stream)
	folder := config.ChunkFileStoreDirectory + b.Filehash + "/"
	counter := 0

	filepath.Walk(folder, func(path string, info os.FileInfo, err error) error {
		counter++
		return nil
	})

	// missing chunks
	if counter != b.ChunkLength {
		c.JSON(200, gin.H{
			"code": 1,
			"msg":  "Missing chunks",
		})
		return
	}

	// iterate files
	filepath.Walk(folder, func(path string, info os.FileInfo, err error) error {
		f, err := ioutil.ReadFile(path)
		mdhash.Update(f)
		return nil
	})
	hash := mdhash.Sum()
	log.Printf("hash after server calculation is: %s\n", hash)

	c.JSON(200, gin.H{
		"code": 0,
		"msg":  "",
		"data": gin.H{
			"Hash server calculation": hash,
		},
	})
}
