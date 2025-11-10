package repository

import (
	"context"
	"file-sentinel/cmd/file-sentinel/model"
	"file-sentinel/ent"
	"time"

	_ "github.com/lib/pq"
)

func NewRepositoryEnt(host, port, user, password, database string) (Repository, error) {
	client, err := ent.Open("postgres", "host="+host+
		" port="+port+
		" user="+user+
		" password="+password+
		" dbname="+database+
		" sslmode=disable")
	if err != nil {
		panic("failed to connect to database: " + err.Error())
	}

	// create the schema if it doesn't exist
	if err := client.Schema.Create(context.Background()); err != nil {
		panic(err)
	}
	return &RepositoryEnt{
		client: client,
	}, nil
}

type RepositoryEnt struct {
	client *ent.Client
}

// GetFileList implements Repository.
func (r *RepositoryEnt) GetFileList(dirPath string) ([]model.FileInfo, error) {
	panic("unimplemented")
}

// GetFileList implements Repository.
func (r *RepositoryEnt) InsertFileList(podName string, dirPath string, fileInfos []model.FileInfo) error {
	ctx := context.Background()
	var fileList []string
	for _, fileInfo := range fileInfos {
		fileList = append(fileList, fileInfo.FileName)
	}
	_, err := r.client.FileInfo.
		Create().
		SetNodeName(podName).
		SetMountPath(dirPath).
		SetFileList(fileList).
		SetCreatedAt(time.Now()).
		Save(ctx)
	if err != nil {
		return err
	}
	return nil
}
