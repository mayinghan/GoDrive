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
			log.Println("No token in the cookie, auth denied")
			c.Abort()
			return
		}

		log.Println("get token: ", token)

		clm, err := utils.ParseToken(token)

		var msg string
		if err != nil {
			suc = false
			switch err.(*jwt.ValidationError).Errors {
			case jwt.ValidationErrorExpired:
				{
					fmt.Println("token expired")
					msg = "token expired, please login again"
				}
			default:
				msg = "token auth check failed"
			}
		}
		// if auth failed
		if !suc {
			c.JSON(http.StatusOK, gin.H{
				"code": 2,
				"msg":  msg,
			})
			log.Println("auth failed due to ", msg)
			c.Abort()
			return
		}

		log.Printf("old token expires at %d\n", clm.ExpiresAt)
		// auth suc
		username := clm.Username
		c.Set("username", username)
		// refresh cookie
		cookie, _ := c.Request.Cookie("token")
		// refresh jwt token
		tokenStr, err := utils.Gentoken(username)
		newClm, _ := utils.ParseToken(tokenStr)
		log.Printf("new token expires at %d\n", newClm.ExpiresAt)
		c.SetCookie(cookie.Name, tokenStr, config.CookieLife, cookie.Path, cookie.Domain, false, false)

		c.Next()
	}
}
