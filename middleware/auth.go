// Package middleware to check auth status
package middleware

import (
	"GoDrive/config"
	"GoDrive/utils"
	"fmt"
	"log"
	"net/http"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

// JWT : a middleware the parse the jwt in the cookie
func JWT() gin.HandlerFunc {
	return func(c *gin.Context) {
		var suc bool
		suc = true

		token, err := c.Cookie("token")
		if err != nil || token == "" {
			c.JSON(http.StatusOK, gin.H{
				"code": 1,
				"msg":  "No token in the cookie, auth denied",
			})
			c.Abort()
			return
		}

		log.Println("get token: ", token)

		_, err = utils.ParseToken(token)
		var msg string
		if err != nil {
			suc = false
			switch err.(*jwt.ValidationError).Errors {
			case jwt.ValidationErrorExpired:
				{
					fmt.Println("token expired")
					msg = "token expired"
				}
			default:
				msg = "token auth check failed"
			}
		}

		// if auth failed
		if !suc {
			c.JSON(http.StatusUnauthorized, gin.H{
				"code": 1,
				"msg":  msg,
			})
			fmt.Println("auth failed")
			c.Abort()
			return
		}

		// auth suc, refresh token
		cookie, _ := c.Request.Cookie("token")
		c.SetCookie(cookie.Name, cookie.Value, config.CookieLife, cookie.Path, cookie.Domain, false, false)
		c.Next()
	}
}
