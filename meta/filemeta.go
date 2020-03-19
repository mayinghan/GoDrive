package meta

import (
	"GoDrive/db"
	"fmt"
	"sort"
)

// FileMeta contains file meta info struct
type FileMeta struct {
	FileMD5  string `json:"hashkey"`
	FileName string `json:"name"`
	FileSize int64  `json:"size"`
	Location string `json:"location"`
	UploadAt string `json:"date"`
}

var fileMetas map[string]FileMeta

// when import this package, init() will be called
func init() {
	fileMetas = make(map[string]FileMeta)
}

// UpdateFileMeta : add/modify file meta info in RAM
func UpdateFileMeta(fm FileMeta) {
	fileMetas[fm.FileMD5] = fm
}

// UpdateFileMetaDB : add/modify file meta into tbl_file and tbl_userfile DBs
func UpdateFileMetaDB(fm FileMeta, username string) bool {

	fileSucc, err := db.OnFileUploadFinished(fm.FileMD5, fm.FileName, fm.FileSize, fm.Location)
	userSucc, errr := db.OnFileUploadUser(username, fm.FileMD5, fm.FileSize, fm.FileName)

	if err != nil {
		fmt.Println(err.Error())
		panic(err.Error())
	}
	if errr != nil {
		fmt.Println(errr.Error())
		panic(err.Error())
	}
	if fileSucc && userSucc {
		return true
	}
	return false
}

// GetFileMeta : get FileMeta struct based on give SHA1 hash code
func GetFileMeta(sha1 string) FileMeta {
	return fileMetas[sha1]
}

// GetFileMetaDB : get file meta info from DB
func GetFileMetaDB(hash string) (FileMeta, error) {
	tFile, err := db.GetFileMeta(hash)
	if err != nil {
		return FileMeta{}, err
	}

	fMeta := FileMeta{
		FileMD5:  tFile.FileHash,
		FileName: tFile.FileName.String,
		FileSize: tFile.FileSize.Int64,
		Location: tFile.FileLocation.String,
	}

	return fMeta, nil
}

// GetLastFileMetas : get the last `count` files' meta datas
func GetLastFileMetas(count int) []FileMeta {
	count = minInt(count, len(fileMetas))
	fMetaSlice := make([]FileMeta, len(fileMetas))
	for _, v := range fileMetas {
		fMetaSlice = append(fMetaSlice, v)
	}
	// sorted by 'uploadAt'
	sort.Sort(SortedByUploadTime(fMetaSlice))
	return fMetaSlice[0:count]
}

// GetLastFileMetasDB : get last `limit` files meta from DB
func GetLastFileMetasDB(limit int) ([]FileMeta, error) {
	files, err := db.GetLastNMetaList(limit)
	if err != nil {
		return make([]FileMeta, 0), err
	}

	fMetas := make([]FileMeta, len(files))
	for i := 0; i < len(fMetas); i++ {
		fMetas[i] = FileMeta{
			FileMD5:  files[i].FileHash,
			FileName: files[i].FileName.String,
			FileSize: files[i].FileSize.Int64,
			Location: files[i].FileLocation.String,
		}
	}

	return fMetas, nil
}

// RemoveMeta : remove the file meta, in the future, need to consider about multithreading security
func RemoveMeta(FileMD5 string) {
	delete(fileMetas, FileMD5)
}

// RemoveMetaDB removes a file meta from the db (remove success, delete meta)
func RemoveMetaDB(username string, filesha string) (bool, bool) {
	succ, err := db.OnFileRemoveUser(username, filesha)
	if err != nil {
		fmt.Println(err.Error())
	}
	if succ {
		return db.OnFileRemoved(filesha)
	}
	return false, false
}

func minInt(a, b int) int {
	if a < b {
		return a
	}
	return b
}
