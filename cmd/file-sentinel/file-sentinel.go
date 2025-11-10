package main

import (
	"file-sentinel/cmd/file-sentinel/model"
	"file-sentinel/cmd/file-sentinel/repository"
	"flag"
	"log/slog"
	"os"
	"path/filepath"
	"time"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
)

type fileSentinel struct {
	repo         repository.Repository
	nodeName     string
	targetDir    string
	minuteCycle  int
	clientset    *kubernetes.Clientset
	hostMountPod *corev1.Pod
	namespace    string
}

func NewFileSentinel(repo repository.Repository, guestDir string, minuteCycle int) *fileSentinel {
	var kubeconfig *string
	if home := homedir.HomeDir(); home != "" {
		kubeconfig = flag.String("kubeconfig", filepath.Join(home, ".kube", "config"), "absolute path to the kubeconfig file")
	} else {
		kubeconfig = flag.String("kubeconfig", "", "absolute path to the kubeconfig file")
	}
	flag.Parse()

	// get node from kubernetes pod environment
	nodeName := os.Getenv("MY_NODE_NAME")

	// Initialize Kubernetes clientset
	config, err := clientcmd.BuildConfigFromFlags("", *kubeconfig)
	if err != nil {
		slog.Error("Error building kubeconfig", "error", err)
		panic(err)
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		slog.Error("Error creating Kubernetes client", "error", err)
		panic(err)
	}

	return &fileSentinel{
		repo:         repo,
		nodeName:     nodeName,
		minuteCycle:  minuteCycle,
		namespace:    "default",
		hostMountPod: nil,
		clientset:    clientset,
	}
}

func (f *fileSentinel) Start() error {
	// Initial file list update
	updateFileHandler := func() {
		fileList, err := f.getFileList(f.targetDir)
		if err != nil {
			slog.Error("Error getting file list", "error", err)
			return
		}
		err = f.addNewFileList(fileList)
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
	f.targetDir = dirPath
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
			ID:          0,
			FileName:    entry.Name(),
			OwnerUser:   "", // OwnerUser retrieval is platform-dependent and omitted here
			Mode:        info.Mode().String(),
			SizeBytes:   info.Size(),
			RegularFile: info.Mode().IsRegular(),
			ModifiedAt:  info.ModTime(),
		})
	}

	for _, file := range fileList {
		slog.Debug("Found file", "name", file.FileName, "", file.SizeBytes)
	}

	return fileList, nil
}

func (f *fileSentinel) addNewFileList(fileList []model.FileInfo) error {
	// Placeholder for updating file list in the repository
	err := f.repo.InsertFileList(f.nodeName, f.targetDir, fileList)
	if err != nil {
		return err
	}

	return nil
}
