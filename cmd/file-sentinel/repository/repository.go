package repository

import "file-sentinel/cmd/file-sentinel/model"

type Repository interface {
	InsertFileList(podName string, dirPath string, fileList []model.FileInfo) error
	GetFileList(dirPath string) ([]model.FileInfo, error)
}
