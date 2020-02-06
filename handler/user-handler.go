package handler

import (
	"GoDrive/cache"
	"GoDrive/db"
	"GoDrive/utils"
	"fmt"
	"math/rand"
	"net/http"
	"strconv"
	"time"

	"github.com/garyburd/redigo/redis"
	"github.com/gin-gonic/gin"
)

const salt = "&6ty"

// RegisterHandler handles user registration. Method: POST
func RegisterHandler(c *gin.Context) {
	var regInput db.RegInfo
	if err := c.ShouldBindJSON(&regInput); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"msg":  err.Error(),
			"code": 1,
		})

		return
	}

	fmt.Printf("%v\n", regInput)

	// encrypt the password
	encryptedPwd := utils.MD5([]byte(regInput.Password + salt))
	regInput.Password = encryptedPwd
	suc, msg, err := db.UserRegister(&regInput)

	if suc {
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

// RegisVerifyEmailHandler : send verify code to user email to finish registration
func RegisVerifyEmailHandler(c *gin.Context) {
	type verifyEmail struct {
		Email string `json:"email" binding:"required"`
	}

	var vrfEmail verifyEmail
	if err := c.ShouldBindJSON(&vrfEmail); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"msg":   "Internal error happened",
			"code":  1,
			"error": err.Error(),
		})
		return
	}

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

	if storedTime != 0 && currTimestamp-int64(storedTime) < 180 {
		fmt.Println("dont send email again")
	}

	redisConn.Do("HSET", vrfEmail.Email, "create_at", currTimestamp)

	code := rand.Intn(899999) + 100000
	s := strconv.Itoa(code)
	fmt.Println(s)
	err = utils.SendMail(vrfEmail.Email)
	if err != nil {
		panic(err)
	}
}
