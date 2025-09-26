// Package mocks provides mock implementations for testing the setup component
package mocks

import (
	"os"
	"path/filepath"
	"strings"
	"time"
)

// MockFileSystem provides a mock implementation of filesystem operations
type MockFileSystem struct {
	Files       map[string][]byte // Virtual file system
	Directories map[string]bool   // Virtual directories
	Permissions map[string]os.FileMode // File permissions
	
	// Error simulation
	ReadErrors  map[string]error
	WriteErrors map[string]error
	StatErrors  map[string]error
	
	// Call tracking
	ReadCalls  []string
	WriteCalls []WriteCall
	StatCalls  []string
	MkdirCalls []string
}

// WriteCall tracks write operation calls
type WriteCall struct {
	Path string
	Data []byte
	Perm os.FileMode
}

// NewMockFileSystem creates a new mock filesystem
func NewMockFileSystem() *MockFileSystem {
	return &MockFileSystem{
		Files:       make(map[string][]byte),
		Directories: make(map[string]bool),
		Permissions: make(map[string]os.FileMode),
		ReadErrors:  make(map[string]error),
		WriteErrors: make(map[string]error),
		StatErrors:  make(map[string]error),
		ReadCalls:   make([]string, 0),
		WriteCalls:  make([]WriteCall, 0),
		StatCalls:   make([]string, 0),
		MkdirCalls:  make([]string, 0),
	}
}

// ReadFile mocks reading a file
func (m *MockFileSystem) ReadFile(filename string) ([]byte, error) {
	m.ReadCalls = append(m.ReadCalls, filename)
	
	if err, exists := m.ReadErrors[filename]; exists {
		return nil, err
	}
	
	if data, exists := m.Files[filename]; exists {
		return data, nil
	}
	
	return nil, &os.PathError{
		Op:   "open",
		Path: filename,
		Err:  os.ErrNotExist,
	}
}

// WriteFile mocks writing a file
func (m *MockFileSystem) WriteFile(filename string, data []byte, perm os.FileMode) error {
	m.WriteCalls = append(m.WriteCalls, WriteCall{
		Path: filename,
		Data: data,
		Perm: perm,
	})
	
	if err, exists := m.WriteErrors[filename]; exists {
		return err
	}
	
	// Create directory structure
	dir := filepath.Dir(filename)
	m.Directories[dir] = true
	
	// Write file
	m.Files[filename] = data
	m.Permissions[filename] = perm
	
	return nil
}

// Stat mocks file stat operation
func (m *MockFileSystem) Stat(name string) (os.FileInfo, error) {
	m.StatCalls = append(m.StatCalls, name)
	
	if err, exists := m.StatErrors[name]; exists {
		return nil, err
	}
	
	// Check if it's a file
	if data, exists := m.Files[name]; exists {
		return &MockFileInfo{
			name:    filepath.Base(name),
			size:    int64(len(data)),
			mode:    m.Permissions[name],
			modTime: time.Now(),
			isDir:   false,
		}, nil
	}
	
	// Check if it's a directory
	if _, exists := m.Directories[name]; exists {
		return &MockFileInfo{
			name:    filepath.Base(name),
			size:    0,
			mode:    os.ModeDir | 0755,
			modTime: time.Now(),
			isDir:   true,
		}, nil
	}
	
	return nil, &os.PathError{
		Op:   "stat",
		Path: name,
		Err:  os.ErrNotExist,
	}
}

// MkdirAll mocks creating directories
func (m *MockFileSystem) MkdirAll(path string, perm os.FileMode) error {
	m.MkdirCalls = append(m.MkdirCalls, path)
	
	if err, exists := m.WriteErrors[path]; exists {
		return err
	}
	
	// Create all parent directories
	parts := strings.Split(filepath.Clean(path), string(filepath.Separator))
	currentPath := ""
	
	for _, part := range parts {
		if part == "" {
			currentPath = string(filepath.Separator)
			continue
		}
		
		if currentPath == string(filepath.Separator) {
			currentPath = filepath.Join(currentPath, part)
		} else {
			currentPath = filepath.Join(currentPath, part)
		}
		
		m.Directories[currentPath] = true
		m.Permissions[currentPath] = perm
	}
	
	return nil
}

// Remove mocks file/directory removal
func (m *MockFileSystem) Remove(name string) error {
	delete(m.Files, name)
	delete(m.Directories, name)
	delete(m.Permissions, name)
	return nil
}

// RemoveAll mocks recursive removal
func (m *MockFileSystem) RemoveAll(path string) error {
	// Remove all files and directories under path
	for filePath := range m.Files {
		if strings.HasPrefix(filePath, path) {
			delete(m.Files, filePath)
			delete(m.Permissions, filePath)
		}
	}
	
	for dirPath := range m.Directories {
		if strings.HasPrefix(dirPath, path) {
			delete(m.Directories, dirPath)
			delete(m.Permissions, dirPath)
		}
	}
	
	return nil
}

// Exists checks if a file or directory exists
func (m *MockFileSystem) Exists(path string) bool {
	_, fileExists := m.Files[path]
	_, dirExists := m.Directories[path]
	return fileExists || dirExists
}

// IsDir checks if path is a directory
func (m *MockFileSystem) IsDir(path string) bool {
	_, exists := m.Directories[path]
	return exists
}

// SetReadError configures the mock to return an error on read
func (m *MockFileSystem) SetReadError(path string, err error) {
	m.ReadErrors[path] = err
}

// SetWriteError configures the mock to return an error on write
func (m *MockFileSystem) SetWriteError(path string, err error) {
	m.WriteErrors[path] = err
}

// SetStatError configures the mock to return an error on stat
func (m *MockFileSystem) SetStatError(path string, err error) {
	m.StatErrors[path] = err
}

// AddFile adds a file to the mock filesystem
func (m *MockFileSystem) AddFile(path string, content []byte, perm os.FileMode) {
	m.Files[path] = content
	m.Permissions[path] = perm
	
	// Ensure parent directory exists
	dir := filepath.Dir(path)
	m.Directories[dir] = true
}

// AddDirectory adds a directory to the mock filesystem
func (m *MockFileSystem) AddDirectory(path string, perm os.FileMode) {
	m.Directories[path] = true
	m.Permissions[path] = perm
}

// Reset clears all files, directories, and call tracking
func (m *MockFileSystem) Reset() {
	m.Files = make(map[string][]byte)
	m.Directories = make(map[string]bool)
	m.Permissions = make(map[string]os.FileMode)
	m.ReadErrors = make(map[string]error)
	m.WriteErrors = make(map[string]error)
	m.StatErrors = make(map[string]error)
	m.ReadCalls = make([]string, 0)
	m.WriteCalls = make([]WriteCall, 0)
	m.StatCalls = make([]string, 0)
	m.MkdirCalls = make([]string, 0)
}

// MockFileInfo implements os.FileInfo for testing
type MockFileInfo struct {
	name    string
	size    int64
	mode    os.FileMode
	modTime time.Time
	isDir   bool
}

func (m *MockFileInfo) Name() string       { return m.name }
func (m *MockFileInfo) Size() int64        { return m.size }
func (m *MockFileInfo) Mode() os.FileMode  { return m.mode }
func (m *MockFileInfo) ModTime() time.Time { return m.modTime }
func (m *MockFileInfo) IsDir() bool        { return m.isDir }
func (m *MockFileInfo) Sys() interface{}   { return nil }