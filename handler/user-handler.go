package handler

import (
	"GoDrive/db"
	"GoDrive/utils"
	"fmt"
	"net/http"

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
		utils.SendMail(regInput.Email)
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
