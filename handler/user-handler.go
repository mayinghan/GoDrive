package handler

import (
	"GoDrive/cache"
	"GoDrive/config"
	"GoDrive/db"
	"GoDrive/utils"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"strconv"
	"time"

	"github.com/garyburd/redigo/redis"
	"github.com/gin-gonic/gin"
)

const salt = "&6ty"

// LoginHandler handles user login.
func LoginHandler(c *gin.Context) {

	var userInput db.LoginInfo
	if err := c.ShouldBindJSON(&userInput); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 1,
			"msg":  err.Error(),
		})
		panic(err)
	}

	log.Printf("%v\n", userInput)

	suc, username, msg, err := db.UserLogin(&userInput)

	if suc {
		//Create the expiration time (1 hour) and the JWT claim
		tokenStr, err := utils.Gentoken(username)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"code":  1,
				"msg":   "Internal server error: Failed to create JWT token.",
				"error": err.Error(),
			})
		} else {
			c.SetCookie(
				"token",           //name
				tokenStr,          //value
				config.CookieLife, //max age
				"/",               //path
				"localhost",       //domain
				false,             //secure
				false,             //httponly
			)
			c.JSON(http.StatusOK, gin.H{
				"code": 0,
				"msg":  msg,
				"data": struct {
					Username string `json:"username"`
				}{
					Username: username,
				},
			})

		}
	} else if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":  1,
			"msg":   msg,
			"error": err.Error(),
		})
	} else {
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"code": 1,
			"msg":  msg,
		})
	}

	return
}

// RegisterHandler handles user registration. Method: POST
func RegisterHandler(c *gin.Context) {
	var regInput db.RegInfo
	if err := c.ShouldBindJSON(&regInput); err != nil {
		fmt.Printf("%v\n", err)
		c.JSON(http.StatusBadRequest, gin.H{
			"msg":  err.Error(),
			"code": 1,
		})
		panic(err)
	}

	fmt.Printf("%v\n", regInput)

	veriCode := regInput.Code
	rc := cache.EmailVeriPool().Get()
	code, err := redis.Uint64(rc.Do("HGET", regInput.Email, "code"))
	if err != nil {
		fmt.Println(err.Error())
		c.JSON(500, gin.H{
			"code": 1,
			"msg":  "The code is expired!",
		})
		return
	}
	fmt.Println(veriCode)
	if int64(code)-veriCode != 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 1,
			"msg":  "Invalid verification code!",
		})
		return
	}
	// encrypt the password
	encryptedPwd := utils.MD5([]byte(regInput.Password + salt))
	regInput.Password = encryptedPwd
	suc, msg, err := db.UserRegister(&regInput)

	if suc {
		//Create the expiration time (1 hour) and the JWT claim
		tokenStr, err := utils.Gentoken(regInput.Username)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"code":  1,
				"msg":   "Internal server error: Failed to create JWT token.",
				"error": err.Error(),
			})
		} else {
			c.SetCookie(
				"token",           //name
				tokenStr,          //value
				config.CookieLife, //max age
				"/",               //path
				"localhost",       //domain
				false,             //secure
				false,             //httponly
			)
			c.JSON(http.StatusOK, gin.H{
				"code": 0,
				"msg":  msg,
				"data": struct {
					Username string `json:"username"`
					Email    string `json:"email"`
				}{
					Username: regInput.Username,
					Email:    regInput.Email,
				},
			})
		}
	} else if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":  1,
			"msg":   msg,
			"error": err.Error(),
		})
	} else {
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"code": 1,
			"msg":  msg,
		})
	}
	return
}

// SendVerifyEmailHandler : send verify code to user email to finish registration
func SendVerifyEmailHandler(c *gin.Context) {

	var vrfEmail db.VerifyEmail
	if err := c.ShouldBind(&vrfEmail); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"msg":   "Internal error happened",
			"code":  1,
			"error": err.Error(),
		})
		panic(err)
	}
	suc, msg, err := db.CheckEmail(&vrfEmail)
	if suc {
		// get redis pool connection
		redisConn := cache.EmailVeriPool().Get()
		defer redisConn.Close()

		// check current user email
		currTimestamp := time.Now().UTC().Unix()
		storedTime, err := redis.Uint64(redisConn.Do("HGET", vrfEmail.Email, "create_at"))
		if err != nil {
			fmt.Printf("redis get previous created time failed %v\n", err)
			storedTime = 0
		}

		if storedTime != 0 && currTimestamp-int64(storedTime) < config.SendCodeCoolDown {
			fmt.Println("dont send email again")
			c.JSON(http.StatusBadRequest, gin.H{
				"code": 1,
				"msg":  "Send request too fast! Please wait " + strconv.FormatInt(config.SendCodeCoolDown+10, 10) + "s to resend the code",
			})
			return
		}

		rand.Seed(currTimestamp)
		code := rand.Intn(899999) + 100000
		s := strconv.Itoa(code)
		redisConn.Do("HMSET", vrfEmail.Email, "create_at", currTimestamp, "code", code)
		// code expires after 10 min
		redisConn.Do("EXPIRE", vrfEmail.Email, 600)
		fmt.Println(s)
		err = utils.SendMail(vrfEmail.Email, s)
		if err != nil {
			panic(err)
		}
	} else {
		fmt.Println("email already exists")
		c.JSON(http.StatusBadRequest, gin.H{
			"code":  1,
			"msg":   msg,
			"error": err.Error(),
		})
		return
	}
}

// UserInfo : Query user info
func UserInfo(c *gin.Context) {
	// getting username
	username, exist := c.Get("username")
	if exist {
		fmt.Printf("username: %s\n", username)
	}
	token, _ := c.Cookie("token")
	fmt.Printf("Got user token: %s\n", token)
	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"msg":  "",
	})
}
