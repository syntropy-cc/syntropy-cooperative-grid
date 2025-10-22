//go:build linux
// +build linux

package node

import (
	"context"
	"fmt"
	"os/exec"
	"regexp"
	"strconv"
	"strings"

	"node-component/src/internal/helpers"
	"node-component/src/internal/types"
)

// USBDetectorLinux implements USB detection for Linux systems
type USBDetectorLinux struct {
	*USBDetectorBase
}

// NewUSBDetectorLinux creates a new Linux USB detector
func NewUSBDetectorLinux(logger types.Logger) *USBDetectorLinux {
	base := NewUSBDetectorBase("linux", logger)
	return &USBDetectorLinux{
		USBDetectorBase: base,
	}
}

// Linux-specific implementation will override the base method

// DetectDevices detects all available USB devices on Linux
func (udl *USBDetectorLinux) DetectDevices(ctx context.Context) ([]types.USBDevice, error) {
	udl.logger.Debug("Detecting USB devices on Linux...")

	// Use lsblk to detect block devices
	devices, err := udl.detectBlockDevices(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to detect block devices: %w", err)
	}

	// Filter USB devices
	var usbDevices []types.USBDevice
	for _, device := range devices {
		if udl.isUSBDevice(ctx, device) {
			usbDevices = append(usbDevices, device)
		}
	}

	udl.logger.Debug("Detected USB devices", "count", len(usbDevices))
	return usbDevices, nil
}

// DetectRemovableDevices detects only removable USB devices
func (udl *USBDetectorLinux) DetectRemovableDevices(ctx context.Context) ([]types.USBDevice, error) {
	devices, err := udl.DetectDevices(ctx)
	if err != nil {
		return nil, err
	}

	var removable []types.USBDevice
	for _, device := range devices {
		if device.IsRemovable && !device.IsSystem {
			removable = append(removable, device)
		}
	}

	udl.logger.Debug("Detected removable USB devices", "count", len(removable))
	return removable, nil
}

// GetDeviceInfo gets detailed information about a specific device
func (udl *USBDetectorLinux) GetDeviceInfo(ctx context.Context, devicePath string) (*types.USBDevice, error) {
	// Get basic device info
	device, err := udl.getDeviceInfoFromPath(ctx, devicePath)
	if err != nil {
		return nil, fmt.Errorf("failed to get device info: %w", err)
	}

	// Get additional USB-specific info
	usbInfo, err := udl.getUSBDeviceInfo(ctx, devicePath)
	if err != nil {
		udl.logger.Warn("Failed to get USB-specific info", "device", devicePath, "error", err)
	} else {
		// Merge USB info
		device.Vendor = usbInfo.Vendor
		device.Model = usbInfo.Model
		device.Serial = usbInfo.Serial
		device.Speed = usbInfo.Speed
	}

	return device, nil
}

// IsDeviceRemovable checks if a device is removable
func (udl *USBDetectorLinux) IsDeviceRemovable(ctx context.Context, devicePath string) (bool, error) {
	// Use lsblk to check if device is removable
	cmd := exec.CommandContext(ctx, "lsblk", "-o", "NAME,RM", "-n", "-r")
	output, err := cmd.Output()
	if err != nil {
		return false, fmt.Errorf("failed to check device removability: %w", err)
	}

	lines := strings.Split(string(output), "\n")
	deviceName := helpers.GetDeviceNameFromPath(devicePath)

	for _, line := range lines {
		fields := strings.Fields(line)
		if len(fields) >= 2 && fields[0] == deviceName {
			return fields[1] == "1", nil
		}
	}

	return false, fmt.Errorf("device not found: %s", devicePath)
}

// GetSystemDevices returns all system devices to avoid writing to them
func (udl *USBDetectorLinux) GetSystemDevices(ctx context.Context) ([]string, error) {
	var systemDevices []string

	// Get root filesystem device
	if rootDevice, err := udl.getRootFilesystemDevice(ctx); err == nil {
		systemDevices = append(systemDevices, rootDevice)
	}

	// Get boot device
	if bootDevice, err := udl.getBootDevice(ctx); err == nil {
		systemDevices = append(systemDevices, bootDevice)
	}

	// Get devices mounted at critical paths
	criticalPaths := []string{"/", "/boot", "/var", "/usr", "/home"}
	for _, path := range criticalPaths {
		if device, err := udl.getMountedDevice(ctx, path); err == nil && device != "" {
			systemDevices = append(systemDevices, device)
		}
	}

	udl.logger.Debug("Identified system devices", "devices", systemDevices)
	return systemDevices, nil
}

// Private helper methods

// detectBlockDevices detects all block devices using lsblk
func (udl *USBDetectorLinux) detectBlockDevices(ctx context.Context) ([]types.USBDevice, error) {
	cmd := exec.CommandContext(ctx, "lsblk", "-o", "NAME,SIZE,TYPE,MOUNTPOINT", "-n", "-r")
	output, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("failed to run lsblk: %w", err)
	}

	var devices []types.USBDevice
	lines := strings.Split(string(output), "\n")

	for _, line := range lines {
		if line == "" {
			continue
		}

		fields := strings.Fields(line)
		if len(fields) < 3 {
			continue
		}

		// Only process disk devices (not partitions)
		if fields[2] != "disk" {
			continue
		}

		devicePath := "/dev/" + fields[0]
		capacity, err := udl.parseSize(fields[1])
		if err != nil {
			udl.logger.Warn("Failed to parse device size", "device", devicePath, "size", fields[1], "error", err)
			continue
		}

		device := types.USBDevice{
			Path:        devicePath,
			RawPath:     devicePath,
			Capacity:    capacity,
			Speed:       0, // Will be filled later
			Vendor:      "",
			Model:       "",
			Serial:      "",
			IsSystem:    false,
			IsRemovable: false, // Will be filled later
		}

		devices = append(devices, device)
	}

	return devices, nil
}

// isUSBDevice checks if a device is a USB device
func (udl *USBDetectorLinux) isUSBDevice(ctx context.Context, device types.USBDevice) bool {
	// Check if device is in /sys/block and has USB attributes
	deviceName := helpers.GetDeviceNameFromPath(device.Path)
	sysPath := fmt.Sprintf("/sys/block/%s/device", deviceName)

	// Check if device has USB attributes
	if helpers.FileExists(sysPath + "/subsystem") {
		subsystem, err := helpers.ReadFileSafely(sysPath + "/subsystem")
		if err == nil && strings.Contains(string(subsystem), "usb") {
			return true
		}
	}

	// Check if device is in USB directory structure
	if helpers.FileExists(sysPath + "/../subsystem") {
		subsystem, err := helpers.ReadFileSafely(sysPath + "/../subsystem")
		if err == nil && strings.Contains(string(subsystem), "usb") {
			return true
		}
	}

	return false
}

// getDeviceInfoFromPath gets basic device information from device path
func (udl *USBDetectorLinux) getDeviceInfoFromPath(ctx context.Context, devicePath string) (*types.USBDevice, error) {
	// Check if device exists
	if !helpers.FileExists(devicePath) {
		return nil, fmt.Errorf("device does not exist: %s", devicePath)
	}

	// Get device size using lsblk
	cmd := exec.CommandContext(ctx, "lsblk", "-o", "NAME,SIZE,TYPE", "-n", "-r", devicePath)
	output, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("failed to get device info: %w", err)
	}

	lines := strings.Split(string(output), "\n")
	if len(lines) == 0 || lines[0] == "" {
		return nil, fmt.Errorf("no device info found for: %s", devicePath)
	}

	fields := strings.Fields(lines[0])
	if len(fields) < 3 {
		return nil, fmt.Errorf("invalid device info format: %s", lines[0])
	}

	capacity, err := udl.parseSize(fields[1])
	if err != nil {
		return nil, fmt.Errorf("failed to parse device size: %w", err)
	}

	// Check if device is removable
	isRemovable, err := udl.IsDeviceRemovable(ctx, devicePath)
	if err != nil {
		udl.logger.Warn("Failed to check device removability", "device", devicePath, "error", err)
		isRemovable = false
	}

	// Check if device is system device
	systemDevices, err := udl.GetSystemDevices(ctx)
	if err != nil {
		udl.logger.Warn("Failed to get system devices", "error", err)
		systemDevices = []string{}
	}

	isSystem := false
	for _, sysDevice := range systemDevices {
		if sysDevice == devicePath {
			isSystem = true
			break
		}
	}

	device := &types.USBDevice{
		Path:        devicePath,
		RawPath:     devicePath,
		Capacity:    capacity,
		Speed:       0,
		Vendor:      "",
		Model:       "",
		Serial:      "",
		IsSystem:    isSystem,
		IsRemovable: isRemovable,
	}

	return device, nil
}

// getUSBDeviceInfo gets USB-specific information about a device
func (udl *USBDetectorLinux) getUSBDeviceInfo(ctx context.Context, devicePath string) (*types.USBDevice, error) {
	deviceName := helpers.GetDeviceNameFromPath(devicePath)

	// Try to get USB info from sysfs
	usbInfo := &types.USBDevice{}

	// Get vendor and model from /sys/block/device/device
	sysPath := fmt.Sprintf("/sys/block/%s/device", deviceName)

	if helpers.FileExists(sysPath + "/vendor") {
		vendor, err := helpers.ReadFileSafely(sysPath + "/vendor")
		if err == nil {
			usbInfo.Vendor = strings.TrimSpace(string(vendor))
		}
	}

	if helpers.FileExists(sysPath + "/model") {
		model, err := helpers.ReadFileSafely(sysPath + "/model")
		if err == nil {
			usbInfo.Model = strings.TrimSpace(string(model))
		}
	}

	if helpers.FileExists(sysPath + "/serial") {
		serial, err := helpers.ReadFileSafely(sysPath + "/serial")
		if err == nil {
			usbInfo.Serial = strings.TrimSpace(string(serial))
		}
	}

	// Try to get speed from USB subsystem
	if helpers.FileExists(sysPath + "/../speed") {
		speed, err := helpers.ReadFileSafely(sysPath + "/../speed")
		if err == nil {
			if speedInt, err := strconv.Atoi(strings.TrimSpace(string(speed))); err == nil {
				usbInfo.Speed = speedInt
			}
		}
	}

	return usbInfo, nil
}

// getRootFilesystemDevice gets the device where root filesystem is mounted
func (udl *USBDetectorLinux) getRootFilesystemDevice(ctx context.Context) (string, error) {
	cmd := exec.CommandContext(ctx, "df", "/", "--output=source")
	output, err := cmd.Output()
	if err != nil {
		return "", fmt.Errorf("failed to get root filesystem device: %w", err)
	}

	lines := strings.Split(string(output), "\n")
	if len(lines) < 2 {
		return "", fmt.Errorf("invalid df output")
	}

	device := strings.TrimSpace(lines[1])
	return device, nil
}

// getBootDevice gets the device where boot partition is mounted
func (udl *USBDetectorLinux) getBootDevice(ctx context.Context) (string, error) {
	cmd := exec.CommandContext(ctx, "df", "/boot", "--output=source")
	output, err := cmd.Output()
	if err != nil {
		// Boot might not be on separate partition
		return "", nil
	}

	lines := strings.Split(string(output), "\n")
	if len(lines) < 2 {
		return "", nil
	}

	device := strings.TrimSpace(lines[1])
	return device, nil
}

// getMountedDevice gets the device mounted at a specific path
func (udl *USBDetectorLinux) getMountedDevice(ctx context.Context, path string) (string, error) {
	cmd := exec.CommandContext(ctx, "df", path, "--output=source")
	output, err := cmd.Output()
	if err != nil {
		return "", fmt.Errorf("failed to get mounted device for %s: %w", path, err)
	}

	lines := strings.Split(string(output), "\n")
	if len(lines) < 2 {
		return "", fmt.Errorf("invalid df output for %s", path)
	}

	device := strings.TrimSpace(lines[1])
	return device, nil
}

// parseSize parses a size string from lsblk (e.g., "8G", "500M")
func (udl *USBDetectorLinux) parseSize(sizeStr string) (int64, error) {
	sizeStr = strings.TrimSpace(sizeStr)

	// Regular expression to match size with unit
	re := regexp.MustCompile(`^(\d+(?:\.\d+)?)\s*([KMGTPE]?)$`)
	matches := re.FindStringSubmatch(sizeStr)

	if len(matches) != 3 {
		return 0, fmt.Errorf("invalid size format: %s", sizeStr)
	}

	value, err := strconv.ParseFloat(matches[1], 64)
	if err != nil {
		return 0, fmt.Errorf("invalid size value: %s", matches[1])
	}

	unit := matches[2]
	var multiplier int64 = 1

	switch unit {
	case "K":
		multiplier = 1024
	case "M":
		multiplier = 1024 * 1024
	case "G":
		multiplier = 1024 * 1024 * 1024
	case "T":
		multiplier = 1024 * 1024 * 1024 * 1024
	case "P":
		multiplier = 1024 * 1024 * 1024 * 1024 * 1024
	case "E":
		multiplier = 1024 * 1024 * 1024 * 1024 * 1024 * 1024
	}

	return int64(value * float64(multiplier)), nil
}

// USBDeviceInfo represents USB-specific device information
type USBDeviceInfo struct {
	Vendor string
	Model  string
	Serial string
	Speed  int
}

// NewPlatformUSBDetector creates the platform-specific USB detector
func NewPlatformUSBDetector(logger types.Logger) USBDetector {
	return NewUSBDetectorLinux(logger)
}
