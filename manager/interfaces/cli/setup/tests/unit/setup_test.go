package tests

import (
	"os"
	"testing"

	"github.com/syntropy-cc/syntropy-cooperative-grid/manager/interfaces/cli/setup"
	"github.com/syntropy-cc/syntropy-cooperative-grid/manager/interfaces/cli/setup/internal/types"
)

func TestGetSyntropyDir(t *testing.T) {
	// Test with default home directory
	dir := setup.GetSyntropyDir()
	
	// Check if the returned directory ends with .syntropy or Syntropy (for Windows)
	if len(dir) == 0 || (dir[len(dir)-9:] != ".syntropy" && dir[len(dir)-8:] != "Syntropy") {
		t.Errorf("Expected directory to end with .syntropy or Syntropy, got: %s", dir)
	}
}

func TestSetupOptions(t *testing.T) {
	// Test default options
	options := types.SetupOptions{}
	
	if options.Force {
		t.Errorf("Expected Force to be false by default")
	}
	
	if options.InstallService {
		t.Errorf("Expected InstallService to be false by default")
	}
	
	// Test setting options
	options = types.SetupOptions{
		Force:          true,
		InstallService: true,
		ConfigPath:     "/custom/config/path",
		HomeDir:        "/custom/home",
	}
	
	if !options.Force {
		t.Errorf("Expected Force to be true")
	}
	
	if !options.InstallService {
		t.Errorf("Expected InstallService to be true")
	}
	
	if options.ConfigPath != "/custom/config/path" {
		t.Errorf("Expected ConfigPath to be /custom/config/path, got: %s", options.ConfigPath)
	}
	
	if options.HomeDir != "/custom/home" {
		t.Errorf("Expected HomeDir to be /custom/home, got: %s", options.HomeDir)
	}
}

func TestSetupResult(t *testing.T) {
	// Test result with success
	result := types.SetupResult{
		Success:     true,
		ConfigPath:  "/path/to/config",
	}
	
	if !result.Success {
		t.Errorf("Expected Success to be true")
	}
	
	if result.ConfigPath != "/path/to/config" {
		t.Errorf("Expected ConfigPath to be /path/to/config, got: %s", result.ConfigPath)
	}
	
	// Test result with error
	errorResult := types.SetupResult{
		Success:     false,
		Error:       os.ErrNotExist,
	}
	
	if errorResult.Success {
		t.Errorf("Expected Success to be false")
	}
	
	if errorResult.Error != os.ErrNotExist {
		t.Errorf("Expected Error to be os.ErrNotExist")
	}
}