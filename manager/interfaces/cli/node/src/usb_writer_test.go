package node

import (
	"os"
	"path/filepath"
	"testing"
)

// TestValidateISOFileWithoutSizeValidation tests that ISO size validation is disabled
func TestValidateISOFileWithoutSizeValidation(t *testing.T) {
	// Create temporary directory for testing
	tempDir := t.TempDir()

	// Create a small mock ISO file (35MB - smaller than the previous 100MB minimum)
	mockISOPath := filepath.Join(tempDir, "ubuntu-24.04-live-server-amd64.iso")
	mockISOContent := make([]byte, 35*1024*1024) // 35MB
	if err := os.WriteFile(mockISOPath, mockISOContent, 0644); err != nil {
		t.Fatalf("Failed to create mock ISO: %v", err)
	}

	// Create USB writer manager with mock logger
	writerManager := &USBWriterManager{
		logger: &MockLogger{},
	}

	// Test that validation passes even with small ISO
	err := writerManager.validateISOFile(mockISOPath)
	if err != nil {
		t.Errorf("Expected validation to pass for small ISO, but got error: %v", err)
	}
}

// TestValidateISOFileWithNonExistentFile tests validation with non-existent file
func TestValidateISOFileWithNonExistentFile(t *testing.T) {
	// Create USB writer manager with mock logger
	writerManager := &USBWriterManager{
		logger: &MockLogger{},
	}

	// Test with non-existent file
	err := writerManager.validateISOFile("/non/existent/file.iso")
	if err == nil {
		t.Errorf("Expected validation to fail for non-existent file, but got no error")
	}
}

// TestValidateISOFileWithUnreadableFile tests validation with unreadable file
func TestValidateISOFileWithUnreadableFile(t *testing.T) {
	// Create temporary directory for testing
	tempDir := t.TempDir()

	// Create a file with no read permissions
	mockISOPath := filepath.Join(tempDir, "unreadable.iso")
	if err := os.WriteFile(mockISOPath, []byte("test"), 0000); err != nil {
		t.Fatalf("Failed to create unreadable file: %v", err)
	}

	// Create USB writer manager with mock logger
	writerManager := &USBWriterManager{
		logger: &MockLogger{},
	}

	// Test that validation fails for unreadable file
	err := writerManager.validateISOFile(mockISOPath)
	if err == nil {
		t.Errorf("Expected validation to fail for unreadable file, but got no error")
	}
}
