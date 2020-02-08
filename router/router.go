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

	router.POST("/api/user/signup", handler.RegisterHandler)
	router.GET("/api/user/verify", handler.SendVerifyEmailHandler)
	router.POST("/api/user/login", handler.LoginHandler)

	router.Use(middleware.JWT())
	router.GET("/api/user/info", handler.UserInfo)
	return router
}
