package main

import (
	"file-sentinel/cmd/file-sentinel/model"
	"file-sentinel/cmd/file-sentinel/repository"
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestNewFileSentinel(t *testing.T) {
	repo, err := repository.NewRepositoryNative()
	if err != nil {
		t.Fatalf("Failed to create repository: %v", err)
	}

	dirPath := "/test/dir"
	minuteCycle := 5

	fs := NewFileSentinel(repo, dirPath, minuteCycle)

	if fs.dirPath != dirPath {
		t.Errorf("dirPath not set correctly: got %q, want %q", fs.dirPath, dirPath)
	}
	if fs.minuteCycle != minuteCycle {
		t.Errorf("minuteCycle not set correctly: got %d, want %d", fs.minuteCycle, minuteCycle)
	}
	if fs.repo != repo {
		t.Error("repo not set correctly")
	}
}

func TestGetFiles(t *testing.T) {
	repo, err := repository.NewRepositoryNative()
	if err != nil {
		t.Fatalf("Failed to create repository: %v", err)
	}

	t.Run("valid directory with files", func(t *testing.T) {
		// Create a temporary directory with test files
		tmpDir := t.TempDir()

		// Create test files
		testFiles := []string{"file1.txt", "file2.log", "file3.dat"}
		for _, filename := range testFiles {
			filePath := filepath.Join(tmpDir, filename)
			if err := os.WriteFile(filePath, []byte("test content"), 0644); err != nil {
				t.Fatalf("Failed to create test file: %v", err)
			}
		}

		fs := NewFileSentinel(repo, tmpDir, 1)

		fileList, err := fs.getFileList(tmpDir)
		if err != nil {
			t.Fatalf("getFileList failed: %v", err)
		}

		if len(fileList) != len(testFiles) {
			t.Errorf("Expected %d files, got %d", len(testFiles), len(fileList))
		}

		// Verify file properties
		for _, fileInfo := range fileList {
			if fileInfo.FileBaseName == "" {
				t.Error("FileBaseName should not be empty")
			}
			if fileInfo.FileFile == "" {
				t.Error("FileFile (full path) should not be empty")
			}
			if fileInfo.SizeBytes < 0 {
				t.Error("SizeBytes should not be negative")
			}
			if !fileInfo.RegularFile {
				t.Error("Test files should be regular files")
			}
			if fileInfo.ModifiedAt.IsZero() {
				t.Error("ModifiedAt should be set")
			}
		}
	})

	t.Run("nonexistent directory", func(t *testing.T) {
		fs := NewFileSentinel(repo, "/nonexistent/path", 1)

		_, err := fs.getFileList("/nonexistent/path")
		if err == nil {
			t.Error("Expected error for nonexistent directory, got nil")
		}
	})

	t.Run("empty directory", func(t *testing.T) {
		// Create an empty temporary directory
		tmpDir := t.TempDir()

		fs := NewFileSentinel(repo, tmpDir, 1)

		fileList, err := fs.getFileList(tmpDir)
		if err != nil {
			t.Fatalf("getFileList failed: %v", err)
		}

		if len(fileList) != 0 {
			t.Errorf("Expected empty file list, got %d files", len(fileList))
		}
	})

	t.Run("directory with subdirectories", func(t *testing.T) {
		// Create a temporary directory with files and subdirectories
		tmpDir := t.TempDir()

		// Create test file
		if err := os.WriteFile(filepath.Join(tmpDir, "file1.txt"), []byte("test"), 0644); err != nil {
			t.Fatalf("Failed to create test file: %v", err)
		}

		// Create subdirectory
		subDir := filepath.Join(tmpDir, "subdir")
		if err := os.Mkdir(subDir, 0755); err != nil {
			t.Fatalf("Failed to create subdirectory: %v", err)
		}

		fs := NewFileSentinel(repo, tmpDir, 1)

		fileList, err := fs.getFileList(tmpDir)
		if err != nil {
			t.Fatalf("getFileList failed: %v", err)
		}

		// Should have 2 entries: 1 file and 1 directory
		if len(fileList) != 2 {
			t.Errorf("Expected 2 entries (1 file + 1 dir), got %d", len(fileList))
		}

		// Find the directory entry
		foundDir := false
		foundFile := false
		for _, fileInfo := range fileList {
			if fileInfo.FileBaseName == "subdir" {
				foundDir = true
				if fileInfo.RegularFile {
					t.Error("Subdirectory should not be marked as RegularFile")
				}
			}
			if fileInfo.FileBaseName == "file1.txt" {
				foundFile = true
				if !fileInfo.RegularFile {
					t.Error("Regular file should be marked as RegularFile")
				}
			}
		}

		if !foundDir {
			t.Error("Subdirectory not found in file list")
		}
		if !foundFile {
			t.Error("Regular file not found in file list")
		}
	})
}

func TestStart(t *testing.T) {
	repo, err := repository.NewRepositoryNative()
	if err != nil {
		t.Fatalf("Failed to create repository: %v", err)
	}

	t.Run("start with valid directory", func(t *testing.T) {
		// Create a temporary directory
		tmpDir := t.TempDir()

		// Create a test file
		if err := os.WriteFile(filepath.Join(tmpDir, "test.txt"), []byte("test"), 0644); err != nil {
			t.Fatalf("Failed to create test file: %v", err)
		}

		fs := NewFileSentinel(repo, tmpDir, 1)

		// Start will panic because repository.PutFileList is not implemented
		// We expect this panic
		defer func() {
			if r := recover(); r != nil {
				// Expected panic from unimplemented repository
				if r != "unimplemented" {
					t.Errorf("Unexpected panic: %v", r)
				}
			} else {
				t.Error("Expected panic from unimplemented repository, but got none")
			}
		}()

		_ = fs.Start()
	})
}

func TestSetMinuteCycle(t *testing.T) {
	repo, err := repository.NewRepositoryNative()
	if err != nil {
		t.Fatalf("Failed to create repository: %v", err)
	}

	fs := NewFileSentinel(repo, "/test/dir", 1)

	testCases := []struct {
		name        string
		minuteCycle int
	}{
		{"positive value", 5},
		{"zero value", 0},
		{"large value", 60},
		{"negative value", -1},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			fs.setMinuteCycle(tc.minuteCycle)
			if fs.minuteCycle != tc.minuteCycle {
				t.Errorf("setMinuteCycle failed: got %d, want %d", fs.minuteCycle, tc.minuteCycle)
			}
		})
	}
}

func TestSetDirPath(t *testing.T) {
	repo, err := repository.NewRepositoryNative()
	if err != nil {
		t.Fatalf("Failed to create repository: %v", err)
	}

	fs := NewFileSentinel(repo, "/test/dir", 1)

	testCases := []struct {
		name    string
		dirPath string
	}{
		{"absolute path", "/new/test/path"},
		{"empty path", ""},
		{"relative path", "./relative/path"},
		{"path with spaces", "/path with spaces/dir"},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			fs.setDirPath(tc.dirPath)
			if fs.dirPath != tc.dirPath {
				t.Errorf("setDirPath failed: got %q, want %q", fs.dirPath, tc.dirPath)
			}
		})
	}
}

func TestUpdateFileList(t *testing.T) {
	repo, err := repository.NewRepositoryNative()
	if err != nil {
		t.Fatalf("Failed to create repository: %v", err)
	}

	fs := NewFileSentinel(repo, "/test/dir", 1)

	testFileList := []model.FileInfo{
		{
			ID:           1,
			FileBaseName: "test.txt",
			FileFile:     "/test/dir/test.txt",
			OwnerUser:    "testuser",
			Mode:         "-rw-r--r--",
			SizeBytes:    1024,
			RegularFile:  true,
			ModifiedAt:   time.Now(),
		},
		{
			ID:           2,
			FileBaseName: "test2.log",
			FileFile:     "/test/dir/test2.log",
			OwnerUser:    "testuser",
			Mode:         "-rw-r--r--",
			SizeBytes:    2048,
			RegularFile:  true,
			ModifiedAt:   time.Now(),
		},
	}

	// This will panic because the repository is not implemented
	// We expect this panic
	defer func() {
		if r := recover(); r != nil {
			// Expected panic from unimplemented repository
			if r != "unimplemented" {
				t.Errorf("Unexpected panic: %v", r)
			}
		} else {
			t.Error("Expected panic from unimplemented repository, but got none")
		}
	}()

	_ = fs.updateFileList("/test/dir", testFileList)
}

func Test_fileSentinel_Start(t *testing.T) {
	tests := []struct {
		name string // description of this test case
		// Named input parameters for receiver constructor.
		repo        repository.Repository
		dirPath     string
		minuteCycle int
		wantErr     bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := NewFileSentinel(tt.repo, tt.dirPath, tt.minuteCycle)
			gotErr := f.Start()
			if gotErr != nil {
				if !tt.wantErr {
					t.Errorf("Start() failed: %v", gotErr)
				}
				return
			}
			if tt.wantErr {
				t.Fatal("Start() succeeded unexpectedly")
			}
		})
	}
}
