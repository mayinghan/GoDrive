package database

import (
	"GoDrive/database/mydb"
	"fmt"
)

// TableFile : a struct that correspond to a file's info in the database
type TableFile struct {
	FileHash     string
	FileName     string
	FileSize     int64
	FileLocation string
}

// OnFileUploadFinished returns a bool after the file is uploaded to db
func OnFileUploadFinished(filehash string, filename string, filesize int64, filelocation string) bool {
	// using prepared statement to prevent
	statement, err := mydb.DBConn().Prepare(
		/* insert ignore: if an error occured during a bacth of insertions,
		only the one with error will fail, the rest of insertions will succeed.
		*/
		"insert ignore into tbl_file (`sha1`, `name`, `size`, `location`, `status`) " +
			"values (?,?,?,?,0)",
	)

	if err != nil {
		fmt.Println("Failed to prepare statement, err: " + err.Error())
		return false
	}

	defer statement.Close()

	result, err := statement.Exec(filehash, filename, filesize, filelocation)
	if err != nil {
		fmt.Println(err.Error())
		return false
	}

	// check if the file is insert -> see how many row is affected
	if row, err := result.RowsAffected(); err == nil {
		if row <= 0 {
			fmt.Printf("File with hash %s was already in DB", filehash)
		}
		return true
	}

	return false
}

// GetFileMeta : query from the DB and return the TableFile
func GetFileMeta(filehash string) (*TableFile, error) {
	// `TODO: implement this function
	return &TableFile{}, nil
}
