package db

import (
	"GoDrive/db/mydb"
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
func OnFileUploadFinished(filehash string, filename string, filesize int64, filelocation string) (bool, error) {
	// using prepared statement to prevent
	statement, err := mydb.DBConn().Prepare(
		/* insert ignore: if an error occured during a bacth of insertions,
		only the one with error will fail, the rest of insertions will succeed.
		*/
		"insert ignore into tbl_file (`sha1`, `name`, `size`, `location`) " +
			"values (?,?,?,?)",
	)

	if err != nil {
		fmt.Println("Failed to prepare statement, err: " + err.Error())
		return false, err
	}

	defer statement.Close()

	result, err := statement.Exec(filehash, filename, filesize, filelocation)
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
		fmt.Printf("File with hash %s was already in DB\n", filehash)
		// update copies value by adding one
		updateStmt, err := mydb.DBConn().Prepare("update tbl_file set copies=copies + 1 where sha1=?")
		if err != nil {
			fmt.Println("Update prepare stmt failed")
			fmt.Println(err.Error())
			return false, err
		}
		result, err := updateStmt.Exec(filehash)
		if err != nil {
			fmt.Println(err.Error())
			return false, err
		}
		if count, _ := result.RowsAffected(); count == 1 {
			fmt.Println("One record updated")
			return true, nil
		}
		return false, nil
	}

	return true, nil
}

// GetFileMeta : query from the DB and return the TableFile
func GetFileMeta(filehash string) (*TableFile, error) {
	// using prepared statement to query from the DB
	statement, err := mydb.DBConn().Prepare(
		"select sha1, name, size, location from tbl_file where sha1 = ?  limit 1")

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
	statement, err := mydb.DBConn().Prepare("select 1 from tbl_file where sha1=? limit 1")
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

// GetLastNMetaList : get last n file metas by batch. param: batch size
func GetLastNMetaList(batchSize int) ([]TableFile, error) {
	statement, err := mydb.DBConn().Prepare("select sha1, location, name, size from tbl_file limit ?")
	if err != nil {
		fmt.Println("Failed to prepare statement, err: " + err.Error())
		return nil, err
	}
	defer statement.Close()

	rows, err := statement.Query(batchSize)
	if err != nil {
		fmt.Println(err.Error())
		return nil, err
	}

	var files []TableFile
	for rows.Next() {
		file := TableFile{}
		err = rows.Scan(&file.FileHash, &file.FileLocation,
			&file.FileName, &file.FileSize)
		if err != nil {
			fmt.Println(err.Error())
			break
		}
		files = append(files, file)
	}

	fmt.Printf("Get %d files \n", len(files))
	return files, nil
}

// OnFileDecrementCopies : decrements the number of copies available by users
func OnFileDecrementCopies(numCopies int, filehash string) bool {
	statement, err := mydb.DBConn().Prepare("update tbl_file set copies = ? where sha1 = ?")
	if err != nil {
		fmt.Println("Failed to prepare statement, err: " + err.Error())
		return false
	}
	defer statement.Close()
	marked, err := statement.Exec(numCopies, filehash)
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

// OnFileDelete : deletes the file from the tbl_file db (when copies <= 1)
func OnFileDelete(filehash string) bool {
	statement, err := mydb.DBConn().Prepare("delete from tbl_file where sha1 = ?")
	if err != nil {
		fmt.Println("Failed to prepare statement, err: " + err.Error())
		return false
	}
	defer statement.Close()
	result, err := statement.Exec(filehash)
	if err != nil {
		fmt.Println("Failed to delete file, err: " + err.Error())
		return false
	}
	rows, err := result.RowsAffected()
	if err != nil {
		return false
	}
	fmt.Println("Updated table:", rows)
	return true
}

// OnFileRemoved : either decrements the number of copies or deletes the file in the tbl_file db (remove success, need to delete meta?)
func OnFileRemoved(filehash string) (bool, bool) {
	statement, err := mydb.DBConn().Prepare("select copies from tbl_file where sha1 = ?")
	if err != nil {
		fmt.Println("Failed to prepare statement, err: " + err.Error())
		return false, false
	}
	defer statement.Close()

	var numCopies int
	err = statement.QueryRow(filehash).Scan(&numCopies)
	if err != nil {
		fmt.Println(err.Error())
		return false, false
	}

	if numCopies >= 2 {
		return OnFileDecrementCopies(numCopies-1, filehash), false
	}

	return OnFileDelete(filehash), true

}
