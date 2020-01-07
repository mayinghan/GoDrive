package util

import (
	"crypto/md5"
	"crypto/sha1"
	"encoding/hex"
	"hash"
	"io"
	"os"
	"path/filepath"
)

// Sha1Stream is a stream hashed by SHA1 algorithm
type Sha1Stream struct {
	_sha1 hash.Hash
}

// Update the Sha1Stream with the given value
func (s *Sha1Stream) Update(data []byte) {
	if s._sha1 == nil {
		s._sha1 = sha1.New()
	}

	s._sha1.Write(data)
}

// Sum returns a hex string of the checksum of current hash
func (s *Sha1Stream) Sum() string {
	return hex.EncodeToString(s._sha1.Sum([]byte("")))
}

// Sha1 returns the hex string of the SHA1 hased checksum of input []byte
func Sha1(data []byte) string {
	_sha1 := sha1.New()
	_sha1.Write(data)
	return hex.EncodeToString(_sha1.Sum([]byte("")))
}

// MD5 returns the hex string of the MD5 hashed checksum of input []byte
func MD5(data []byte) string {
	_md5 := md5.New()
	_md5.Write(data)
	return hex.EncodeToString(_md5.Sum([]byte("")))
}

// FileSHA1 returns the hex string of the SHA1 hashed checksum of file content
func FileSHA1(file *os.File) string {
	_sha1 := sha1.New()
	io.Copy(_sha1, file)
	return hex.EncodeToString(_sha1.Sum(nil))
}

// FileMD5 returns the hex string of the MD5 hashed checksum of file content
func FileMD5(file *os.File) string {
	_md5 := md5.New()
	io.Copy(_md5, file)
	return hex.EncodeToString(_md5.Sum(nil))
}

// PathExists check if path exists
func PathExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	// either case above, then return with the error message
	return false, err
}

// GetFileSize returns the size of the file
func GetFileSize(file string) int64 {
	var result int64
	filepath.Walk(file, func(path string, f os.FileInfo, err error) error {
		result = f.Size()
		return nil
	})

	return result
}
