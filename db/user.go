package db

import (
	"GoDrive/db/mydb"
	"fmt"
)

// UserRegister handles user registration. Return a bool and a server message
func UserRegister(username string, password string) (bool, string, error) {
	stmt, err := mydb.DBConn().Prepare(
		"insert ignore into tbl_user (`username`, `password`) values(?, ?)")

	if err != nil {
		e := fmt.Sprint("Internal server error: Failed to insert to DB.")
		fmt.Println(e + err.Error())
		return false, e, err
	}
	defer stmt.Close()

	result, err := stmt.Exec(username, password)
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
