package model

import "time"

type FileInfo struct {
	ID           uint32
	FileBaseName string
	FileFile     string
	OwnerUser    string
	Mode         string
	SizeBytes    int64
	RegularFile  bool
	ModifiedAt   time.Time
}
