package database

import (
	"GoDrive/database/mydb"
	"database/sql"
	"fmt"
)

// TableFile : a struct that correspond to a file's info in the database
type TableFile struct {
	FileHash     string
	FileName     sql.NullString
	FileSize     sql.NullInt64
	FileLocation sql.NullString
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
	// using prepared statement to query from the DB
	statement, err := mydb.DBConn().Prepare(
		"select sha1, name, size, location from tbl_file where sha1 = ? and status=0 limit 1")

	if err != nil {
		fmt.Println("Failed to prepare statement, err: " + err.Error())
		return nil, err
	}

	defer statement.Close()

	var metaFile TableFile
	err = statement.QueryRow(filehash).Scan(
		&metaFile.FileHash, &metaFile.FileName, &metaFile.FileSize, &metaFile.FileLocation)

	if err != nil {
		fmt.Println(err.Error())
		return nil, err
	}

	return &metaFile, nil
}

// check if the file is already uploaded
// func IsFileUploaded(filehash string) bool {
// 	statement, err := mydb.DBConn().Prepare("select * from tbl_")
// }
