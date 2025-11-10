package repository

import "file-sentinel/cmd/file-sentinel/model"

type Repository interface {
	ClearFileList(dirPath string) error
	PutFileList(dirPath string, fileList []model.FileInfo) error
	GetFileList(dirPath string) ([]model.FileInfo, error)
}
