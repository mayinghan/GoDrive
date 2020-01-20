package meta

import (
	"sort"
)

// FileMeta contains file meta info struct
type FileMeta struct {
	FileSha1 string `json:"FileSha1"`
	FileName string `json:"FileName"`
	FileSize int64  `json:"FileSize"`
	Location string `json:"Location"`
	UploadAt string `json:"UploadAt"`
}

var fileMetas map[string]FileMeta

// when import this package, init() will be called
func init() {
	fileMetas = make(map[string]FileMeta)
}

// UpdateFileMeta : add/modify file meta info
func UpdateFileMeta(fm FileMeta) {
	fileMetas[fm.FileSha1] = fm
}

// GetFileMeta : get FileMeta struct based on give SHA1 hash code
func GetFileMeta(sha1 string) FileMeta {
	return fileMetas[sha1]
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

// RemoveMeta : remove the file meta, in the future, need to consider about multithreading security
func RemoveMeta(fileSha1 string) {
	delete(fileMetas, fileSha1)
}

func minInt(a, b int) int {
	if a < b {
		return a
	}
	return b
}
