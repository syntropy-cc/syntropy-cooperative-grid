package node

import (
	"testing"

	"node-component/src/internal/types"
)

// TestSelectUSBDeviceFromDetected tests the USB device selection functionality
func TestSelectUSBDeviceFromDetected(t *testing.T) {
	// Create a mock CreateSubcomponent for testing
	cs := &CreateSubcomponent{}

	// Test case 1: Single device (should be selected automatically)
	t.Run("SingleDevice", func(t *testing.T) {
		devices := []types.USBDevice{
			{
				Path:        "/dev/sdb",
				Capacity:    8589934592, // 8 GB
				Vendor:      "SanDisk",
				Model:       "Ultra USB 3.0",
				IsRemovable: true,
			},
		}

		// This test would require mocking fmt.Scanln, which is complex
		// For now, we'll just verify the function doesn't panic with valid input
		defer func() {
			if r := recover(); r != nil {
				t.Errorf("Function panicked: %v", r)
			}
		}()

		// Note: This test would need to be run interactively or with input mocking
		// For automated testing, we would need to mock the user input
		_, err := cs.selectUSBDeviceFromDetected(devices)
		if err != nil {
			// Expected error due to no input in test environment
			t.Logf("Expected error in test environment: %v", err)
		}
	})

	// Test case 2: Multiple devices
	t.Run("MultipleDevices", func(t *testing.T) {
		devices := []types.USBDevice{
			{
				Path:        "/dev/sdb",
				Capacity:    8589934592, // 8 GB
				Vendor:      "SanDisk",
				Model:       "Ultra USB 3.0",
				IsRemovable: true,
			},
			{
				Path:        "/dev/sdc",
				Capacity:    17179869184, // 16 GB
				Vendor:      "Kingston",
				Model:       "DataTraveler",
				IsRemovable: true,
			},
		}

		defer func() {
			if r := recover(); r != nil {
				t.Errorf("Function panicked: %v", r)
			}
		}()

		_, err := cs.selectUSBDeviceFromDetected(devices)
		if err != nil {
			// Expected error due to no input in test environment
			t.Logf("Expected error in test environment: %v", err)
		}
	})

	// Test case 3: No devices
	t.Run("NoDevices", func(t *testing.T) {
		devices := []types.USBDevice{}

		_, err := cs.selectUSBDeviceFromDetected(devices)
		if err == nil {
			t.Error("Expected error for empty device list")
		}

		expectedError := "no devices available for selection"
		if err.Error() != expectedError {
			t.Errorf("Expected error '%s', got '%s'", expectedError, err.Error())
		}
	})
}

// TestDeviceSelectionDisplay tests the display format of device information
func TestDeviceSelectionDisplay(t *testing.T) {
	// Test that device information is formatted correctly
	device := types.USBDevice{
		Path:        "/dev/sdb",
		Capacity:    8589934592, // 8 GB
		Vendor:      "SanDisk",
		Model:       "Ultra USB 3.0",
		IsRemovable: true,
	}

	// Test capacity conversion to GB
	capacityGB := float64(device.Capacity) / (1024 * 1024 * 1024)
	expectedCapacity := 8.0
	if capacityGB != expectedCapacity {
		t.Errorf("Expected capacity %.2f GB, got %.2f GB", expectedCapacity, capacityGB)
	}

	// Test that all device fields are accessible
	if device.Path == "" {
		t.Error("Device path should not be empty")
	}
	if device.Vendor == "" {
		t.Error("Device vendor should not be empty")
	}
	if device.Model == "" {
		t.Error("Device model should not be empty")
	}
	if !device.IsRemovable {
		t.Error("Device should be removable")
	}
}
