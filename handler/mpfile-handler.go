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
	"sort"
	"strconv"
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
		Filehash    interface{} `json:"filehash"`
		Filename    string      `json:"filename"`
		ChunkLength int         `json:"chunkLength"`
	}

	var b body
	if err := c.ShouldBindJSON(&b); err != nil {
		c.JSON(200, gin.H{
			"code": 1,
			"msg":  err.Error(),
		})
		panic(err)
	}

	fileHash := fmt.Sprintf("%v", b.Filehash)
	log.Printf("Checking integrity.. Filename: %s, Filehash: %v", b.Filename, b.Filehash)

	mdhash := new(utils.MD5Stream)
	folder := config.ChunkFileStoreDirectory + fileHash + "/"
	counter := 0

	filepath.Walk(folder, func(path string, info os.FileInfo, err error) error {
		if !info.IsDir() {
			counter++
		}

		return nil
	})
	fmt.Println(counter)
	fmt.Println(b.ChunkLength)
	// missing chunks
	if counter != b.ChunkLength {
		c.JSON(200, gin.H{
			"code": 1,
			"msg":  "Missing chunks",
		})
		return
	}

	// iterate files
	var chunks sortedChunk
	chunks, _ = ioutil.ReadDir(folder)
	log.Printf("chunk count: %d\n", len(chunks))

	// sort chunks based on name
	sort.Sort(chunks)
	for _, v := range chunks {
		chunkContent, err := ioutil.ReadFile(folder + "/" + v.Name())
		if err != nil {
			panic(err)
		}
		mdhash.Update(chunkContent)
	}
	hash := mdhash.Sum()

	// panic(gin.Error{Err: errors.New("123123")})
	if hash != fileHash {
		c.JSON(200, gin.H{
			"code": 1,
			"msg":  "server file integrity checking failed! Please try to reupload",
		})
		return
	}
	log.Printf("hash after server calculation is: %s\n", hash)

	c.JSON(200, gin.H{
		"code": 0,
		"msg":  "File uploaded successfully",
	})
}

type sortedChunk []os.FileInfo

/**
Comparator interface for SortedChunk
*/
func (a sortedChunk) Len() int {
	return len(a)
}

func (a sortedChunk) Swap(i, j int) {
	a[i], a[j] = a[j], a[i]
}

func (a sortedChunk) Less(i, j int) bool {
	idx1, _ := strconv.Atoi(strings.Split(a[i].Name(), "_")[1])
	idx2, _ := strconv.Atoi(strings.Split(a[j].Name(), "_")[1])

	return idx1 < idx2
}
