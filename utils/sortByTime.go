package utils

import (
	"GoDisk/meta"
	"time"
)

const base = "2006-01-02 15:04:05"

// SortedByUploadTime : a slice of FileMetas sorted by uploaded time
type SortedByUploadTime []meta.FileMeta

/**
Implement the comparator interface
Len(), Swap(), Less()
*/

func (a SortedByUploadTime) Len() int {
	return len(a)
}

func (a SortedByUploadTime) Swap(i, j int) {
	a[i], a[j] = a[j], a[i]
}

func (a SortedByUploadTime) Less(i, j int) bool {
	iTime, _ := time.Parse(base, a[i].UploadAt)
	jTime, _ := time.Parse(base, a[j].UploadAt)

	return iTime.UnixNano() > jTime.UnixNano()
}
