package mocks

import (
	"os"
	"path/filepath"
)

// FileExists checks if a file exists
func FileExists(path string) bool {
	_, err := os.Stat(path)
	return !os.IsNotExist(err)
}

// DirExists checks if a directory exists
func DirExists(path string) bool {
	info, err := os.Stat(path)
	return err == nil && info.IsDir()
}

// CreateDir creates a directory
func CreateDir(path string) error {
	return os.MkdirAll(path, 0755)
}

// WriteFile writes content to a file
func WriteFile(path string, content []byte) error {
	dir := filepath.Dir(path)
	if err := CreateDir(dir); err != nil {
		return err
	}
	return os.WriteFile(path, content, 0644)
}

// ReadFile reads content from a file
func ReadFile(path string) ([]byte, error) {
	return os.ReadFile(path)
}

// RemoveFile removes a file or directory
func RemoveFile(path string) error {
	return os.RemoveAll(path)
}

// GetFileSize returns the size of a file
func GetFileSize(path string) (int64, error) {
	info, err := os.Stat(path)
	if err != nil {
		return 0, err
	}
	return info.Size(), nil
}

// GetFileMode returns the mode of a file
func GetFileMode(path string) (os.FileMode, error) {
	info, err := os.Stat(path)
	if err != nil {
		return 0, err
	}
	return info.Mode(), nil
}

// SetFileMode sets the mode of a file
func SetFileMode(path string, mode os.FileMode) error {
	return os.Chmod(path, mode)
}
