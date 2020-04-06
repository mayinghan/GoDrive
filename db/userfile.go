package db

import (
	"GoDrive/db/mydb"
	"database/sql"
	"fmt"
)

// UserFileDB : a struct that correspond to a file's info in the database
type UserFileDB struct {
	FileHash sql.NullString
	FileName sql.NullString
	FileSize sql.NullInt64
}

// UserFile : a struct that convert the DB model to regular Golang type model (NullStrig -> String)
type UserFile struct {
	FileHash string `json:"key"`
	FileName string `json:"filename"`
	FileSize int    `json:"filesize"`
}

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
	fmt.Printf("username %s, filehash: %s", username, filehash)
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

// CheckDuplicateFile searches tbl_userfile for existing hash.
func CheckDuplicateFile(user string, filehash string) (bool, error) {
	var hash string
	statement, err := mydb.DBConn().Prepare("select hash from tbl_userfile where username = ? and hash = ?")
	if err != nil {
		fmt.Println("Failed to prepare statement, err: " + err.Error())
		return false, err
	}
	defer statement.Close()

	err = statement.QueryRow(user, filehash).Scan(&hash)
	if err != nil {
		if err == sql.ErrNoRows {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

// GetAllUserFiles : Returns all files uploaded by user 'username'
func GetAllUserFiles(user string) (bool, []*UserFile, error) {
	statement, err := mydb.DBConn().Prepare("select hash, filename, size from tbl_userfile where username = ?")
	if err != nil {
		fmt.Println("Failed to prepare statement, err: " + err.Error())
		return false, nil, err
	}
	defer statement.Close()

	rows, err := statement.Query(user)
	if err != nil {
		fmt.Println(err.Error())
		return false, nil, err
	}

	var files []*UserFile
	for rows.Next() {
		fileDB := UserFileDB{}
		err = rows.Scan(&fileDB.FileHash, &fileDB.FileName, &fileDB.FileSize)
		if err != nil {
			fmt.Println(err.Error())
			panic(err.Error())
		}

		file := &UserFile{}
		if fileDB.FileHash.Valid {
			file.FileHash = fileDB.FileHash.String
		}
		if fileDB.FileName.Valid {
			file.FileName = fileDB.FileName.String
		}
		if fileDB.FileSize.Valid {
			file.FileSize = int(fileDB.FileSize.Int64)
		}
		files = append(files, file)
	}
	return true, files, nil
}
