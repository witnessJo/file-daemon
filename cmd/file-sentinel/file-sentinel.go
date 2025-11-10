package main

import (
	"file-sentinel/cmd/file-sentinel/model"
	"file-sentinel/cmd/file-sentinel/repository"
	"log/slog"
	"os"
	"time"
)

type fileSentinel struct {
	dirPath     string
	minuteCycle int
	repo        repository.Repository
}

func NewFileSentinel(repo repository.Repository, dirPath string, minuteCycle int) *fileSentinel {
	return &fileSentinel{
		repo:        repo,
		dirPath:     dirPath,
		minuteCycle: minuteCycle,
	}
}

func (f *fileSentinel) Start() error {

	// Initial file list update
	updateFileHandler := func() {
		fileList, err := f.getFileList(f.dirPath)
		if err != nil {
			slog.Error("Error getting file list", "error", err)
			return
		}
		err = f.updateFileList(f.dirPath, fileList)
		if err != nil {
			slog.Error("Error updating file list", "error", err)
			return
		}
	}

	updateFileHandler()

	ticker := time.NewTicker(time.Duration(f.minuteCycle) * time.Minute)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			updateFileHandler()
		}
	}
}

func (f *fileSentinel) setMinuteCycle(minuteCycle int) {
	f.minuteCycle = minuteCycle
}

func (f *fileSentinel) setDirPath(dirPath string) {
	f.dirPath = dirPath
}

func (f *fileSentinel) getFileList(dirPath string) ([]model.FileInfo, error) {
	entries, err := os.ReadDir(dirPath)
	if err != nil {
		return nil, err
	}
	var fileList []model.FileInfo
	for _, entry := range entries {
		info, err := entry.Info()
		if err != nil {
			return nil, err
		}
		fileList = append(fileList, model.FileInfo{
			ID:           0,
			FileBaseName: entry.Name(),
			FileFile:     dirPath + "/" + entry.Name(),
			OwnerUser:    "", // OwnerUser retrieval is platform-dependent and omitted here
			Mode:         info.Mode().String(),
			SizeBytes:    info.Size(),
			RegularFile:  info.Mode().IsRegular(),
			ModifiedAt:   info.ModTime(),
		})
	}

	for _, file := range fileList {
		slog.Debug("Found file", "name", file.FileBaseName, "size", file.SizeBytes)
	}

	return fileList, nil
}

func (f *fileSentinel) updateFileList(dirPath string, fileList []model.FileInfo) error {
	// Placeholder for updating file list in the repository

	err := f.repo.PutFileList(dirPath, fileList) // Simplified for illustration
	if err != nil {
		return err
	}

	return nil
}
