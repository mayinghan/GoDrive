package middleware

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

// ErrHandler : error handler middleware
func ErrHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				var e gin.Error
				if gErr, ok := err.(gin.Error); ok {
					e = gErr
				} else {
					e.Err = errors.New("Internal Server Error")
				}
				fmt.Printf("Err handler middlware recovering, err: %v\n", e.Error())
				c.JSON(http.StatusInternalServerError, gin.H{
					"code": 1,
					"msg": e.Error(),
				})
				
				return
			}
		}()
		c.Next()
	}
}