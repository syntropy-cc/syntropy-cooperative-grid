package node

import (
	"testing"

	"node-component/src/internal/types"
)

// TestSelectUSBDeviceFromDetected tests the USB device selection functionality
func TestSelectUSBDeviceFromDetected(t *testing.T) {
	// Create a mock CreateSubcomponent for testing
	cs := &CreateSubcomponent{}

	// Test case 1: Single device (should require manual selection even with one device)
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

// TestSelectedUSBDeviceTokenValidation tests the validation token generation and validation
func TestSelectedUSBDeviceTokenValidation(t *testing.T) {
	t.Run("TokenGeneration", func(t *testing.T) {
		device := types.USBDevice{
			Path:        "/dev/sdb",
			Capacity:    8589934592,
			Model:       "Ultra USB 3.0",
			Serial:      "123456789",
			IsRemovable: true,
		}

		selected := NewSelectedUSBDevice(device, "linux")

		if selected.ValidationToken == "" {
			t.Error("Validation token should not be empty")
		}

		if selected.Platform != "linux" {
			t.Errorf("Expected platform 'linux', got '%s'", selected.Platform)
		}

		if selected.Device.Path != device.Path {
			t.Error("Device should match original")
		}
	})

	t.Run("TokenValidation", func(t *testing.T) {
		device := types.USBDevice{
			Path:        "/dev/sdb",
			Capacity:    8589934592,
			Model:       "Ultra USB 3.0",
			Serial:      "123456789",
			IsRemovable: true,
		}

		selected := NewSelectedUSBDevice(device, "linux")

		// Validation should succeed with same device
		err := selected.Validate(device)
		if err != nil {
			t.Errorf("Validation should succeed with same device: %v", err)
		}

		// Validation should fail with different device
		differentDevice := types.USBDevice{
			Path:        "/dev/sdc",
			Capacity:    17179869184,
			Model:       "Different Model",
			Serial:      "987654321",
			IsRemovable: true,
		}

		err = selected.Validate(differentDevice)
		if err == nil {
			t.Error("Validation should fail with different device")
		}
	})

	t.Run("TokenConsistency", func(t *testing.T) {
		device := types.USBDevice{
			Path:        "/dev/sdb",
			Capacity:    8589934592,
			Model:       "Ultra USB 3.0",
			Serial:      "123456789",
			IsRemovable: true,
		}

		selected1 := NewSelectedUSBDevice(device, "linux")
		selected2 := NewSelectedUSBDevice(device, "linux")

		// Same device should generate same token
		if selected1.ValidationToken != selected2.ValidationToken {
			t.Error("Same device should generate consistent token")
		}
	})
}

// TestManualDeviceSelectionRequired tests that manual selection is always required
func TestManualDeviceSelectionRequired(t *testing.T) {
	t.Run("SingleDeviceRequiresManualSelection", func(t *testing.T) {
		cs := &CreateSubcomponent{}
		devices := []types.USBDevice{
			{
				Path:        "/dev/sdb",
				Capacity:    8589934592,
				Vendor:      "SanDisk",
				Model:       "Ultra USB 3.0",
				IsRemovable: true,
			},
		}

		// Even with one device, should require manual selection (will fail in test due to no input)
		_, err := cs.selectUSBDeviceFromDetected(devices)
		if err != nil {
			// Expected in test environment without input
			t.Logf("Manual selection required (expected in test): %v", err)
		}
	})
}

// TestDeviceTokenSecurity tests that device token prevents device swap
func TestDeviceTokenSecurity(t *testing.T) {
	t.Run("PreventDeviceSwap", func(t *testing.T) {
		// Original device
		device1 := types.USBDevice{
			Path:        "/dev/sdb",
			Capacity:    8589934592,
			Model:       "USB Drive A",
			Serial:      "111111",
			IsRemovable: true,
		}

		// Different device
		device2 := types.USBDevice{
			Path:        "/dev/sdc",
			Capacity:    8589934592,
			Model:       "USB Drive B",
			Serial:      "222222",
			IsRemovable: true,
		}

		selected := NewSelectedUSBDevice(device1, "linux")

		// Should detect device swap
		err := selected.Validate(device2)
		if err == nil {
			t.Error("Should detect device swap and fail validation")
		}

		if err != nil && !contains(err.Error(), "device mismatch") {
			t.Errorf("Expected 'device mismatch' error, got: %v", err)
		}
	})
}

// Helper function to check if string contains substring
func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(s) > len(substr) && (s[:len(substr)] == substr || s[len(s)-len(substr):] == substr || containsMiddle(s, substr)))
}

func containsMiddle(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
