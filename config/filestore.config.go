package config

import (
	"runtime"
)

const goos string = runtime.GOOS

var WholeFileStoreLocation string
var ChunkFileStoreDirectory string

func init() {
	if goos == "darwin" {
		WholeFileStoreLocation = "/tmp/"
		ChunkFileStoreDirectory = "/tmp/files/"
	} else {
		WholeFileStoreLocation = "C://Users/liuwi/Desktop/tmp/"
		ChunkFileStoreDirectory = "C://Users/liuwi/Desktop/tmp/files/"
	}
}
