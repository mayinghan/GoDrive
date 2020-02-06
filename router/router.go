package router

import (
	"GoDrive/handler"

	"github.com/gin-gonic/gin"
)

// Router : routing rules
func Router() *gin.Engine {
	router := gin.Default()
	// APIs that don't need auth
	router.POST("/api/user/signup", handler.RegisterHandler)
	router.GET("/api/user/verify", handler.SendVerifyEmailHandler)
	return router
}
