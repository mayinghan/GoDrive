package config

import (
	"runtime"
)

const goos string = runtime.GOOS

// WholeFileStoreLocation is the path for whole file storage
var WholeFileStoreLocation string

// ChunkFileStoreDirectory is the path for file chunk storage
var ChunkFileStoreDirectory string

// DataSourceName is the path for database
var DataSourceName string

// StoreMethod : either "AWS" or "Local"
var StoreMethod string = "Local"

func init() {
	if goos == "darwin" {
		WholeFileStoreLocation = "/tmp/"
		ChunkFileStoreDirectory = "/tmp/files/"
		DataSourceName = "root:123456@tcp(127.0.0.1:13306)/fileserver?charset=utf8"
	} else {
		WholeFileStoreLocation = "C://Users/liuwi/Desktop/tmp/"
		ChunkFileStoreDirectory = "C://Users/liuwi/Desktop/tmp/files/"
		DataSourceName = "root:123456@tcp(192.168.99.100:13306)/fileserver?charset=utf8"
	}
}
