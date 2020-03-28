package router

import (
	"GoDrive/handler"
	"GoDrive/middleware"

	"github.com/gin-gonic/gin"
)

// Router : routing rules
func Router() *gin.Engine {
	router := gin.Default()
	// APIs that don't need auth
	router.Use(middleware.ErrHandler())
	router.POST("/api/user/signup", handler.RegisterHandler)
	router.GET("/api/user/verify", handler.SendVerifyEmailHandler)
	router.POST("/api/user/login", handler.LoginHandler)

	router.Use(middleware.JWT())

	router.GET("/api/user/info", handler.UserInfo)
	router.GET("/api/user/filelist", handler.UserFileList)
	router.GET("/api/file/instantupload", handler.InstantUpload)
	router.POST("/api/file/uploadchunk", handler.GetFileChunk)
	router.POST("/api/file/checkIntegrity", handler.CheckIntegrity)
	router.POST("/api/file/upload", handler.UploadHandler)
	router.GET("/api/file/getfilemeta", handler.GetFileMetaHandler)
	router.GET("/api/file/querybatch", handler.QueryByBatchHandler)
	router.GET("/api/file/prevChunks", handler.GetPrevChunks)
	router.GET("/api/file/download", handler.DownloadHandler)
	router.POST("/api/file/update", handler.FileUpdateHandler)
	router.DELETE("/api/user/file", handler.FileDeleteHandler)

	return router
}
