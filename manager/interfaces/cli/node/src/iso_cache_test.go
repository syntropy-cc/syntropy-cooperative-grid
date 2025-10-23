package node

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"node-component/src/internal/types"
)

// MockLogger implements types.Logger for testing
type MockLogger struct{}

func (m *MockLogger) Debug(msg string, fields ...interface{}) {}
func (m *MockLogger) Info(msg string, fields ...interface{})  {}
func (m *MockLogger) Warn(msg string, fields ...interface{})  {}
func (m *MockLogger) Error(msg string, fields ...interface{}) {}
func (m *MockLogger) Fatal(msg string, fields ...interface{}) {}
func (m *MockLogger) SetLevel(level string)                   {}
func (m *MockLogger) WithFields(fields map[string]interface{}) types.Logger {
	return m
}

// TestListCachedISOs tests the ListCachedISOs functionality
func TestListCachedISOs(t *testing.T) {
	// Create temporary directory for testing
	tempDir := t.TempDir()

	// Create ISO downloader with temp directory
	downloader := &ISODownloaderImpl{
		cacheDir: tempDir,
		logger:   &MockLogger{},
	}

	// Test with empty cache directory
	cachedISOs, err := downloader.ListCachedISOs()
	if err != nil {
		t.Fatalf("ListCachedISOs failed: %v", err)
	}
	if len(cachedISOs) != 0 {
		t.Errorf("Expected 0 cached ISOs, got %d", len(cachedISOs))
	}

	// Create a mock ISO file
	mockISOPath := filepath.Join(tempDir, "ubuntu-24.04-live-server-amd64.iso")
	mockISOContent := []byte("mock ISO content for testing")
	if err := os.WriteFile(mockISOPath, mockISOContent, 0644); err != nil {
		t.Fatalf("Failed to create mock ISO: %v", err)
	}

	// Test with mock ISO file (should be included since we skip validation)
	cachedISOs, err = downloader.ListCachedISOs()
	if err != nil {
		t.Fatalf("ListCachedISOs failed: %v", err)
	}
	if len(cachedISOs) != 1 {
		t.Errorf("Expected 1 cached ISO (validation skipped), got %d", len(cachedISOs))
	}
	if len(cachedISOs) > 0 && cachedISOs[0].Version != "24.04" {
		t.Errorf("Expected version 24.04, got %s", cachedISOs[0].Version)
	}
}

// TestExtractVersionFromFilename tests version extraction from filenames
func TestExtractVersionFromFilename(t *testing.T) {
	downloader := &ISODownloaderImpl{}

	testCases := []struct {
		filename string
		expected string
	}{
		{"ubuntu-24.04-live-server-amd64.iso", "24.04"},
		{"ubuntu-22.04-live-server-amd64.iso", "22.04"},
		{"ubuntu-20.04-live-server-amd64.iso", "20.04"},
		{"ubuntu-24.04.3-live-server-amd64.iso", "24.04.3"},
		{"not-an-iso.txt", ""},
		{"ubuntu-invalid.iso", ""},
	}

	for _, tc := range testCases {
		result := downloader.extractVersionFromFilename(tc.filename)
		if result != tc.expected {
			t.Errorf("extractVersionFromFilename(%s) = %s, expected %s", tc.filename, result, tc.expected)
		}
	}
}

// TestSelectISOFromCache tests ISO selection logic
func TestSelectISOFromCache(t *testing.T) {
	downloader := &ISODownloaderImpl{}

	// Test with empty cache
	result, err := downloader.selectISOFromCache([]*types.ISOInfo{}, "24.04")
	if err != nil {
		t.Fatalf("selectISOFromCache failed: %v", err)
	}
	if result != nil {
		t.Errorf("Expected nil result for empty cache, got %v", result)
	}

	// Test with single matching ISO
	mockISO := &types.ISOInfo{
		Version:      "24.04",
		FileName:     "ubuntu-24.04-live-server-amd64.iso",
		Size:         1024 * 1024 * 1024, // 1GB
		DownloadedAt: time.Now(),
	}

	result, err = downloader.selectISOFromCache([]*types.ISOInfo{mockISO}, "24.04")
	if err != nil {
		t.Fatalf("selectISOFromCache failed: %v", err)
	}
	if result == nil {
		t.Errorf("Expected non-nil result for matching ISO, got nil")
	}
	if result != nil && result.Version != "24.04" {
		t.Errorf("Expected version 24.04, got %s", result.Version)
	}
}
