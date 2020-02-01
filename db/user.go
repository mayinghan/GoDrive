package db

import (
	"GoDrive/db/mydb"
	"fmt"
)

// RegInfo is the registration input: username password and email
type RegInfo struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Email    string `json:"email"`
}

// UserRegister handles user registration. Return a bool and a server message
func UserRegister(regInfo *RegInfo) (bool, string, error) {
	username := regInfo.Username
	password := regInfo.Password
	email := regInfo.Email

	stmt, err := mydb.DBConn().Prepare(
		"insert ignore into tbl_user (`username`, `password`, `email`) values(?, ?, ?)")

	if err != nil {
		e := fmt.Sprint("Internal server error: Failed to insert to DB.")
		fmt.Println(e + err.Error())
		return false, e, err
	}
	defer stmt.Close()

	result, err := stmt.Exec(username, password, email)
	if err != nil {
		e := fmt.Sprint("Internal server error: Failed to insert to DB.")
		fmt.Println(e + err.Error())
		return false, e, err
	}
	// check how many row is affected
	if ra, err := result.RowsAffected(); err == nil && ra > 0 {
		return true, "Registered Successfully!", nil
	} else if err == nil && ra <= 0 {
		return false, "Failed to register. Duplicated user!", nil
	} else {
		return false, "Internal server error: Failed to insert to DB", err
	}
}
