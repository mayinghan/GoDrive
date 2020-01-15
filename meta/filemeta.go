package meta

import "sort"

// FileMeta contains file meta info struct
type FileMeta struct {
	FileSha1 string
	FileName string
	FileSize int64
	Location string
	UploadAt string
}

var fileMetas map[string]FileMeta

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
	fMetaSlice := make([]FileMeta, len(fileMetas))
	for _, v := range fileMetas {
		fMetaSlice = append(fMetaSlice, v)
	}
	// sorted by 'uploadAt'
	sort.Sort(SortedByUploadTime(fMetaSlice))
	return fMetaSlice[0:count]
}
