package model

import "time"

type FileInfo struct {
	ID          uint32
	FileName    string
	OwnerUser   string
	Mode        string
	SizeBytes   int64
	RegularFile bool
	ModifiedAt  time.Time
}
