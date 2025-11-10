package repository

import (
	"file-sentinel/cmd/file-sentinel/model"
)

func NewRepositoryNative() (Repository, error) {
	return &RepositoryNative{}, nil
}

type RepositoryNative struct {
}

// ClearFileList implements Repository.
func (r *RepositoryNative) ClearFileList(dirPath string) error {
	panic("unimplemented")
}

// GetFileList implements Repository.
func (r *RepositoryNative) GetFileList(dirPath string) ([]model.FileInfo, error) {
	panic("unimplemented")
}

// PutFileList implements Repository.
func (r *RepositoryNative) PutFileList(dirPath string, fileList []model.FileInfo) error {
	panic("unimplemented")
}
