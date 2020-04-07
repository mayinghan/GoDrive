package mq

// TransferData : struct for message queue
type TransferData struct {
	FileHash     string
	CurLocation  string
	DestLocation string
	StoreType    string
	IsLarge      bool
}
