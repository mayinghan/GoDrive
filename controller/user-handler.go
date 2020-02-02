package controller

import (
	"GoDrive/db"
	"GoDrive/utils"
	"encoding/json"
	"fmt"
	"net/http"
)

const salt = "&6ty"

type userResponse struct {
	StatusCode int    `json:"code"`
	Msg        string `json:"msg"`
}

func userErrorResp(s int, msg string) userResponse {
	return userResponse{StatusCode: s, Msg: msg}
}

func returnUserRespJSON(w http.ResponseWriter, v userResponse) {
	js, err := json.Marshal(v)
	if err != nil {
		e := fmt.Sprintf("Failed to create json obj %s\n", err.Error())
		panic(e)
	}
	w.Header().Set("Content-Type", "application/json")
	if v.StatusCode != 200 {
		w.WriteHeader(v.StatusCode)
	}
	w.Write(js)
}

// RegisterHandler handles user registration. Method: POST
func RegisterHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		var regInfo db.RegInfo
		if r.Body == nil {
			resp := userErrorResp(http.StatusInternalServerError, "request body is empty")
			returnUserRespJSON(w, resp)
			return
		}

		// request body is a json object
		err := json.NewDecoder(r.Body).Decode(&regInfo)
		// using MD5 and the salt to hash the password
		hashedPassword := utils.MD5([]byte(regInfo.Password + salt))
		regInfo.Password = hashedPassword

		if err != nil {
			returnUserRespJSON(w, userErrorResp(http.StatusInternalServerError, "Failed to parse json body object"))
			return
		}

		status, msg, err := db.UserRegister(&regInfo)
		if err != nil {
			returnUserRespJSON(w, userErrorResp(500, msg))
			return
		}

		if status {
			returnUserRespJSON(w, userResponse{200, msg})
		} else {
			returnUserRespJSON(w, userResponse{http.StatusBadRequest, msg})
		}
		return
	}
}
