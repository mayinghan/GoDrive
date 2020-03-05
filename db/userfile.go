package db

import (
	"GoDrive/db/mydb"
	"fmt"
)

// OnFileUploadUser returns a bool after the file is uploaded to tbl_userfile db
func OnFileUploadUser(username string, filehash string, filesize int64, filename string) (bool, error) {
	statement, err := mydb.DBConn().Prepare(
		/* insert ignore: if an error occured during a bacth of insertions,
		only the one with error will fail, the rest of insertions will succeed.
		*/
		"insert ignore into tbl_userfile (`username`, `hash`, `size`, `filename`) values (?,?,?,?)",
	)
	if err != nil {
		fmt.Println("Failed to prepare statement, err: " + err.Error())
		return false, err
	}
	defer statement.Close()

	result, err := statement.Exec(username, filehash, filesize, filename)
	if err != nil {
		fmt.Println(err.Error())
		return false, err
	}

	// check if the file is insert -> see how many row is affected
	row, err := result.RowsAffected()
	if err != nil {
		fmt.Println("Failed to perform database operation")
		return false, err
	}
	if row <= 0 {
		fmt.Printf("User already has file with hash %s in DB", filehash)
		return true, err
	}

	return true, nil
}

// OnFileRemoveUser : Use a delete flag to mark resources as deleted but not acctually deleted (change status from 0 to 1)
func OnFileRemoveUser(username string, filehash string) (bool, error) {
	statement, err := mydb.DBConn().Prepare("delete from tbl_userfile where username = ? and hash = ?")
	if err != nil {
		fmt.Println("Failed to prepare statement, err: " + err.Error())
		return false, err
	}
	defer statement.Close()
	results, err := statement.Exec(username, filehash)
	if err != nil {
		fmt.Println("Failed to delete data from user table, err: " + err.Error())
		return false, err
	}
	rows, err := results.RowsAffected()
	if err != nil {
		return false, err
	}
	if rows <= 0 {
		fmt.Println("Nothing deleted")
		return true, err
	}
	fmt.Println("Updated table:", rows)
	return true, nil
}
