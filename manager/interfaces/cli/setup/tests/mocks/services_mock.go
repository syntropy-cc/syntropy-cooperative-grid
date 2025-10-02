package mocks

import (
	"os"
	"path/filepath"
)

// FileService provides file operations
type FileService struct{}

// NewFileService creates a new file service
func NewFileService() *FileService {
	return &FileService{}
}

// Exists checks if a file or directory exists
func (fs *FileService) Exists(path string) bool {
	_, err := os.Stat(path)
	return !os.IsNotExist(err)
}

// CreateDir creates a directory
func (fs *FileService) CreateDir(path string) error {
	return os.MkdirAll(path, 0755)
}

// WriteFile writes content to a file
func (fs *FileService) WriteFile(path string, content []byte) error {
	dir := filepath.Dir(path)
	if err := fs.CreateDir(dir); err != nil {
		return err
	}
	return os.WriteFile(path, content, 0644)
}

// ReadFile reads content from a file
func (fs *FileService) ReadFile(path string) ([]byte, error) {
	return os.ReadFile(path)
}

// RemoveFile removes a file or directory
func (fs *FileService) RemoveFile(path string) error {
	return os.RemoveAll(path)
}

// NetworkService provides network operations
type NetworkService struct{}

// NewNetworkService creates a new network service
func NewNetworkService() *NetworkService {
	return &NetworkService{}
}

// CheckConnectivity checks network connectivity
func (ns *NetworkService) CheckConnectivity() bool {
	// Simple connectivity check - in real implementation this would be more sophisticated
	return true
}

// SystemService provides system operations
type SystemService struct{}

// NewSystemService creates a new system service
func NewSystemService() *SystemService {
	return &SystemService{}
}

// GetOS returns the operating system name
func (ss *SystemService) GetOS() string {
	return "linux" // Simplified for testing
}
