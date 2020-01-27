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

// IsFileUploaded : check if the file is already uploaded
func IsFileUploaded(filehash string) (bool, error) {
	statement, err := mydb.DBConn().Prepare("select 1 from tbl_file where sha1=? and status=0 limit 1")
	if err != nil {
		fmt.Println("Failed to prepare statement, err: " + err.Error())
		return false, err
	}

	defer statement.Close()
	rows, err := statement.Query(filehash)
	if err != nil {
		return false, err
	} else if rows == nil || !rows.Next() {
		return false, nil
	}

	return true, nil
}

// GetBatchFileMeta : get file metas by batch. param: batch size
func GetBatchFileMeta(batchSize int) ([]TableFile, error) {
	statement, err := mydb.DBConn().Prepare("select sha1, location, name, size from tbl_file where status=1 limit ?")
	if err != nil {
		fmt.Println("Failed to prepare statement, err: " + err.Error())
		return nil, err
	}
	defer statement.Close()
	// todo: query `limit`
	return nil, nil
}

// OnFileRemoved : Use a delete flag to mark resources as deleted but not acctually deleted (change `status` from 0 to 1)
func OnFileRemoved(filehash string) bool {
	statement, err := mydb.DBConn().Prepare("update tbl_file set status = '1' where sha1 = ? and status = 0")
	if err != nil {
		fmt.Println("Failed to prepare statement, err: " + err.Error())
		return false
	}
	defer statement.Close()
	marked, err := statement.Exec(filehash)
	if err != nil {
		fmt.Println("Failed to update table, err: " + err.Error())
		return false
	}
	rows, err := marked.RowsAffected()
	if err != nil {
		return false
	}
	fmt.Println("Updated table:", rows)
	return true
}
