//go:build darwin
// +build darwin

package node

import (
	"context"
	"fmt"
	"os/exec"
	"regexp"
	"strconv"
	"strings"

	"github.com/syntropy-cc/syntropy-cooperative-grid/manager/interfaces/cli/node/src/internal/helpers"
	"github.com/syntropy-cc/syntropy-cooperative-grid/manager/interfaces/cli/node/src/internal/types"
)

// USBDetectorMacOS implements USB detection for macOS systems
type USBDetectorMacOS struct {
	*USBDetectorBase
}

// NewUSBDetectorMacOS creates a new macOS USB detector
func NewUSBDetectorMacOS(logger types.Logger) *USBDetectorMacOS {
	base := NewUSBDetectorBase("darwin", logger)
	return &USBDetectorMacOS{
		USBDetectorBase: base,
	}
}

// macOS-specific implementation will override the base method

// DetectDevices detects all available USB devices on macOS
func (udm *USBDetectorMacOS) DetectDevices(ctx context.Context) ([]types.USBDevice, error) {
	udm.logger.Debug("Detecting USB devices on macOS...")

	// Use diskutil to detect external devices
	devices, err := udm.detectExternalDevicesWithDiskutil(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to detect external devices: %w", err)
	}

	udm.logger.Debug("Detected USB devices", "count", len(devices))
	return devices, nil
}

// DetectInterchangeableDevices detects only removable USB devices
func (udm *USBDetectorMacOS) DetectRemovableDevices(ctx context.Context) ([]types.USBDevice, error) {
	devices, err := udm.DetectDevices(ctx)
	if err != nil {
		return nil, err
	}

	var removable []types.USBDevice
	for _, device := range devices {
		if device.IsRemovable && !device.IsSystem {
			removable = append(removable, device)
		}
	}

	udm.logger.Debug("Detected removable USB devices", "count", len(removable))
	return removable, nil
}

// GetDeviceInfo gets detailed information about a specific device
func (udm *USBDetectorMacOS) GetDeviceInfo(ctx context.Context, devicePath string) (*types.USBDevice, error) {
	// Get basic device info using diskutil
	device, err := udm.getDeviceInfoWithDiskutil(ctx, devicePath)
	if err != nil {
		return nil, fmt.Errorf("failed to get device info: %w", err)
	}

	// Get USB-specific info using system_profiler
	usbInfo, err := udm.getUSBDeviceInfoWithSystemProfiler(ctx, devicePath)
	if err != nil {
		udm.logger.Warn("Failed to get USB-specific info", "device", devicePath, "error", err)
	} else {
		// Merge USB info
		device.Vendor = usbInfo.Vendor
		device.Model = usbInfo.Model
		device.Serial = usbInfo.Serial
		device.Speed = usbInfo.Speed
	}

	// Check if device is system device
	systemDevices, err := udm.GetSystemDevices(ctx)
	if err != nil {
		udm.logger.Warn("Failed to get system devices", "error", err)
		systemDevices = []string{}
	}

	isSystem := false
	for _, sysDevice := range systemDevices {
		if sysDevice == devicePath {
			isSystem = true
			break
		}
	}

	device.IsSystem = isSystem
	return device, nil
}

// IsDeviceRemovable checks if a device is removable
func (udm *USBDetectorMacOS) IsDeviceRemovable(ctx context.Context, devicePath string) (bool, error) {
	// Use diskutil to check if device is removable
	cmd := exec.CommandContext(ctx, "diskutil", "info", devicePath)
	output, err := cmd.Output()
	if err != nil {
		return false, fmt.Errorf("failed to check device removability: %w", err)
	}

	// Parse output for removable information
	lines := strings.Split(string(output), "\n")
	for _, line := range lines {
		if strings.Contains(line, "Removable Media") {
			parts := strings.Split(line, ":")
			if len(parts) >= 2 {
				removable := strings.TrimSpace(parts[1])
				return removable == "Yes", nil
			}
		}
	}

	return false, nil
}

// GetSystemDevices returns all system devices to avoid writing to them
func (udm *USBDetectorMacOS) GetSystemDevices(ctx context.Context) ([]string, error) {
	var systemDevices []string

	// Get root filesystem device
	if rootDevice, err := udm.getRootFilesystemDevice(ctx); err == nil {
		systemDevices = append(systemDevices, rootDevice)
	}

	// Get boot device
	if bootDevice, err := udm.getBootDevice(ctx); err == nil {
		systemDevices = append(systemDevices, bootDevice)
	}

	// Get devices mounted at critical paths
	criticalPaths := []string{"/", "/System", "/Applications", "/Users"}
	for _, path := range criticalPaths {
		if device, err := udm.getMountedDevice(ctx, path); err == nil && device != "" {
			systemDevices = append(systemDevices, device)
		}
	}

	udm.logger.Debug("Identified system devices", "devices", systemDevices)
	return systemDevices, nil
}

// Private helper methods

// detectExternalDevicesWithDiskutil detects external devices using diskutil
func (udm *USBDetectorMacOS) detectExternalDevicesWithDiskutil(ctx context.Context) ([]types.USBDevice, error) {
	// Use diskutil to list all external devices
	cmd := exec.CommandContext(ctx, "diskutil", "list", "-plist")
	output, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("failed to run diskutil list: %w", err)
	}

	// Parse plist output to find external devices
	devices, err := udm.parseDiskutilPlist(output)
	if err != nil {
		return nil, fmt.Errorf("failed to parse diskutil output: %w", err)
	}

	return devices, nil
}

// getDeviceInfoWithDiskutil gets device info using diskutil
func (udm *USBDetectorMacOS) getDeviceInfoWithDiskutil(ctx context.Context, devicePath string) (*types.USBDevice, error) {
	// Check if device exists
	if !helpers.FileExists(devicePath) {
		return nil, fmt.Errorf("device does not exist: %s", devicePath)
	}

	// Use diskutil info to get device information
	cmd := exec.CommandContext(ctx, "diskutil", "info", devicePath)
	output, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("failed to get device info: %w", err)
	}

	// Parse diskutil info output
	device, err := udm.parseDiskutilInfo(output, devicePath)
	if err != nil {
		return nil, fmt.Errorf("failed to parse device info: %w", err)
	}

	return device, nil
}

// getUSBDeviceInfoWithSystemProfiler gets USB-specific device information
func (udm *USBDetectorMacOS) getUSBDeviceInfoWithSystemProfiler(ctx context.Context, devicePath string) (*USBDeviceInfo, error) {
	// Use system_profiler to get USB device information
	cmd := exec.CommandContext(ctx, "system_profiler", "SPUSBDataType", "-xml")
	output, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("failed to get USB device info: %w", err)
	}

	// Parse XML output to find matching device
	usbInfo, err := udm.parseSystemProfilerUSB(output, devicePath)
	if err != nil {
		return nil, fmt.Errorf("failed to parse USB device info: %w", err)
	}

	return usbInfo, nil
}

// parseDiskutilPlist parses diskutil list plist output
func (udm *USBDetectorMacOS) parseDiskutilPlist(output []byte) ([]types.USBDevice, error) {
	// This is a simplified parser - in production, use proper plist parsing
	var devices []types.USBDevice

	// For now, use a simpler approach with diskutil list
	cmd := exec.Command("diskutil", "list")
	output, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("failed to run diskutil list: %w", err)
	}

	lines := strings.Split(string(output), "\n")
	var currentDevice string

	for _, line := range lines {
		line = strings.TrimSpace(line)

		// Check for device identifier (e.g., /dev/disk2)
		if strings.HasPrefix(line, "/dev/disk") && !strings.Contains(line, " ") {
			currentDevice = line
			continue
		}

		// Check for external device indicators
		if currentDevice != "" && (strings.Contains(line, "External") ||
			strings.Contains(line, "USB") ||
			strings.Contains(line, "Removable")) {

			// Get device size
			size, err := udm.getDeviceSizeWithDiskutil(currentDevice)
			if err != nil {
				udm.logger.Warn("Failed to get device size", "device", currentDevice, "error", err)
				continue
			}

			device := types.USBDevice{
				Path:        currentDevice,
				RawPath:     currentDevice,
				Capacity:    size,
				Speed:       0,
				Vendor:      "",
				Model:       "",
				Serial:      "",
				IsSystem:    false,
				IsRemovable: true,
			}

			devices = append(devices, device)
		}
	}

	return devices, nil
}

// parseDiskutilInfo parses diskutil info output
func (udm *USBDetectorMacOS) parseDiskutilInfo(output []byte, devicePath string) (*types.USBDevice, error) {
	lines := strings.Split(string(output), "\n")

	device := &types.USBDevice{
		Path:        devicePath,
		RawPath:     devicePath,
		Capacity:    0,
		Speed:       0,
		Vendor:      "",
		Model:       "",
		Serial:      "",
		IsSystem:    false,
		IsRemovable: false,
	}

	for _, line := range lines {
		line = strings.TrimSpace(line)

		if strings.Contains(line, "Disk Size") {
			// Parse disk size
			re := regexp.MustCompile(`\((\d+)\s*Bytes\)`)
			matches := re.FindStringSubmatch(line)
			if len(matches) >= 2 {
				if size, err := strconv.ParseInt(matches[1], 10, 64); err == nil {
					device.Capacity = size
				}
			}
		}

		if strings.Contains(line, "Removable Media") {
			parts := strings.Split(line, ":")
			if len(parts) >= 2 {
				removable := strings.TrimSpace(parts[1])
				device.IsRemovable = removable == "Yes"
			}
		}

		if strings.Contains(line, "Device / Media Name") {
			parts := strings.Split(line, ":")
			if len(parts) >= 2 {
				device.Model = strings.TrimSpace(parts[1])
			}
		}
	}

	return device, nil
}

// parseSystemProfilerUSB parses system_profiler USB output
func (udm *USBDetectorMacOS) parseSystemProfilerUSB(output []byte, devicePath string) (*USBDeviceInfo, error) {
	// This is a simplified parser - in production, use proper XML parsing
	// For now, return empty info
	return &USBDeviceInfo{
		Vendor: "",
		Model:  "",
		Serial: "",
		Speed:  0,
	}, nil
}

// getDeviceSizeWithDiskutil gets device size using diskutil
func (udm *USBDetectorMacOS) getDeviceSizeWithDiskutil(devicePath string) (int64, error) {
	cmd := exec.Command("diskutil", "info", devicePath)
	output, err := cmd.Output()
	if err != nil {
		return 0, fmt.Errorf("failed to get device size: %w", err)
	}

	lines := strings.Split(string(output), "\n")
	for _, line := range lines {
		if strings.Contains(line, "Disk Size") {
			re := regexp.MustCompile(`\((\d+)\s*Bytes\)`)
			matches := re.FindStringSubmatch(line)
			if len(matches) >= 2 {
				return strconv.ParseInt(matches[1], 10, 64)
			}
		}
	}

	return 0, fmt.Errorf("could not parse device size")
}

// getRootFilesystemDevice gets the device where root filesystem is mounted
func (udm *USBDetectorMacOS) getRootFilesystemDevice(ctx context.Context) (string, error) {
	cmd := exec.CommandContext(ctx, "df", "/", "-H")
	output, err := cmd.Output()
	if err != nil {
		return "", fmt.Errorf("failed to get root filesystem device: %w", err)
	}

	lines := strings.Split(string(output), "\n")
	if len(lines) < 2 {
		return "", fmt.Errorf("invalid df output")
	}

	fields := strings.Fields(lines[1])
	if len(fields) == 0 {
		return "", fmt.Errorf("invalid df output format")
	}

	return fields[0], nil
}

// getBootDevice gets the device where boot partition is mounted
func (udm *USBDetectorMacOS) getBootDevice(ctx context.Context) (string, error) {
	cmd := exec.CommandContext(ctx, "df", "/System", "-H")
	output, err := cmd.Output()
	if err != nil {
		// System might not be on separate partition
		return "", nil
	}

	lines := strings.Split(string(output), "\n")
	if len(lines) < 2 {
		return "", nil
	}

	fields := strings.Fields(lines[1])
	if len(fields) == 0 {
		return "", nil
	}

	return fields[0], nil
}

// getMountedDevice gets the device mounted at a specific path
func (udm *USBDetectorMacOS) getMountedDevice(ctx context.Context, path string) (string, error) {
	cmd := exec.CommandContext(ctx, "df", path, "-H")
	output, err := cmd.Output()
	if err != nil {
		return "", fmt.Errorf("failed to get mounted device for %s: %w", path, err)
	}

	lines := strings.Split(string(output), "\n")
	if len(lines) < 2 {
		return "", fmt.Errorf("invalid df output for %s", path)
	}

	fields := strings.Fields(lines[1])
	if len(fields) == 0 {
		return "", fmt.Errorf("invalid df output format for %s", path)
	}

	return fields[0], nil
}

// USBDeviceInfo represents USB-specific device information
type USBDeviceInfo struct {
	Vendor string
	Model  string
	Serial string
	Speed  int
}
