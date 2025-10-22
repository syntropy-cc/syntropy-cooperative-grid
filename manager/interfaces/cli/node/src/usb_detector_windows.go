//go:build windows
// +build windows

package node

import (
	"context"
	"fmt"
	"os/exec"
	"strconv"
	"strings"

	"node-component/src/internal/types"
)

// USBDetectorWindows implements USB detection for Windows systems
type USBDetectorWindows struct {
	*USBDetectorBase
}

// NewUSBDetectorWindows creates a new Windows USB detector
func NewUSBDetectorWindows(logger types.Logger) *USBDetectorWindows {
	base := NewUSBDetectorBase("windows", logger)
	return &USBDetectorWindows{
		USBDetectorBase: base,
	}
}

// Windows-specific implementation will override the base method

// DetectDevices detects all available USB devices on Windows
func (udw *USBDetectorWindows) DetectDevices(ctx context.Context) ([]types.USBDevice, error) {
	udw.logger.Debug("Detecting USB devices on Windows...")

	// Use PowerShell to detect USB devices
	devices, err := udw.detectUSBDevicesWithPowerShell(ctx)
	if err != nil {
		// Fallback to WMIC if PowerShell fails
		udw.logger.Warn("PowerShell detection failed, trying WMIC", "error", err)
		devices, err = udw.detectUSBDevicesWithWMIC(ctx)
		if err != nil {
			return nil, fmt.Errorf("failed to detect USB devices: %w", err)
		}
	}

	udw.logger.Debug("Detected USB devices", "count", len(devices))
	return devices, nil
}

// DetectRemovableDevices detects only removable USB devices
func (udw *USBDetectorWindows) DetectRemovableDevices(ctx context.Context) ([]types.USBDevice, error) {
	devices, err := udw.DetectDevices(ctx)
	if err != nil {
		return nil, err
	}

	var removable []types.USBDevice
	for _, device := range devices {
		if device.IsRemovable && !device.IsSystem {
			removable = append(removable, device)
		}
	}

	udw.logger.Debug("Detected removable USB devices", "count", len(removable))
	return removable, nil
}

// GetDeviceInfo gets detailed information about a specific device
func (udw *USBDetectorWindows) GetDeviceInfo(ctx context.Context, devicePath string) (*types.USBDevice, error) {
	// Get basic device info using PowerShell
	device, err := udw.getDeviceInfoWithPowerShell(ctx, devicePath)
	if err != nil {
		// Fallback to WMIC
		device, err = udw.getDeviceInfoWithWMIC(ctx, devicePath)
		if err != nil {
			return nil, fmt.Errorf("failed to get device info: %w", err)
		}
	}

	// Check if device is system device
	systemDevices, err := udw.GetSystemDevices(ctx)
	if err != nil {
		udw.logger.Warn("Failed to get system devices", "error", err)
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
func (udw *USBDetectorWindows) IsDeviceRemovable(ctx context.Context, devicePath string) (bool, error) {
	// Use PowerShell to check if device is removable
	cmd := exec.CommandContext(ctx, "powershell", "-Command",
		fmt.Sprintf("Get-WmiObject -Class Win32_LogicalDisk | Where-Object {$_.DeviceID -eq '%s'} | Select-Object -ExpandProperty DriveType", devicePath))

	output, err := cmd.Output()
	if err != nil {
		return false, fmt.Errorf("failed to check device removability: %w", err)
	}

	driveType := strings.TrimSpace(string(output))
	// DriveType 2 = Removable disk
	return driveType == "2", nil
}

// GetSystemDevices returns all system devices to avoid writing to them
func (udw *USBDetectorWindows) GetSystemDevices(ctx context.Context) ([]string, error) {
	var systemDevices []string

	// Get system drive (usually C:)
	if systemDrive, err := udw.getSystemDrive(ctx); err == nil {
		systemDevices = append(systemDevices, systemDrive)
	}

	// Get boot drive
	if bootDrive, err := udw.getBootDrive(ctx); err == nil {
		systemDevices = append(systemDevices, bootDrive)
	}

	// Get all fixed drives
	fixedDrives, err := udw.getFixedDrives(ctx)
	if err == nil {
		systemDevices = append(systemDevices, fixedDrives...)
	}

	udw.logger.Debug("Identified system devices", "devices", systemDevices)
	return systemDevices, nil
}

// Private helper methods

// detectUSBDevicesWithPowerShell detects USB devices using PowerShell
func (udw *USBDetectorWindows) detectUSBDevicesWithPowerShell(ctx context.Context) ([]types.USBDevice, error) {
	// PowerShell command to get USB removable drives
	cmd := exec.CommandContext(ctx, "powershell", "-Command", `
		Get-WmiObject -Class Win32_LogicalDisk | 
		Where-Object {$_.DriveType -eq 2} | 
		ForEach-Object {
			$drive = $_.DeviceID
			$size = $_.Size
			$free = $_.FreeSpace
			$label = $_.VolumeLabel
			Write-Output "$drive|$size|$free|$label"
		}
	`)

	output, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("failed to run PowerShell command: %w", err)
	}

	var devices []types.USBDevice
	lines := strings.Split(string(output), "\n")

	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		fields := strings.Split(line, "|")
		if len(fields) < 4 {
			continue
		}

		drive := fields[0]
		sizeStr := fields[1]
		label := fields[3]

		size, err := strconv.ParseInt(sizeStr, 10, 64)
		if err != nil {
			udw.logger.Warn("Failed to parse device size", "device", drive, "size", sizeStr, "error", err)
			continue
		}

		// Get additional device info
		device := types.USBDevice{
			Path:        drive,
			RawPath:     drive,
			Capacity:    size,
			Speed:       0, // Will be filled later
			Vendor:      "",
			Model:       label,
			Serial:      "",
			IsSystem:    false,
			IsRemovable: true,
		}

		// Get USB-specific info
		usbInfo, err := udw.getUSBDeviceInfoWithPowerShell(ctx, drive)
		if err == nil {
			device.Vendor = usbInfo.Vendor
			device.Model = usbInfo.Model
			device.Serial = usbInfo.Serial
			device.Speed = usbInfo.Speed
		}

		devices = append(devices, device)
	}

	return devices, nil
}

// detectUSBDevicesWithWMIC detects USB devices using WMIC (fallback)
func (udw *USBDetectorWindows) detectUSBDevicesWithWMIC(ctx context.Context) ([]types.USBDevice, error) {
	// WMIC command to get logical disks
	cmd := exec.CommandContext(ctx, "wmic", "logicaldisk", "where", "drivetype=2", "get", "deviceid,size,freespace,volumename", "/format:csv")

	output, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("failed to run WMIC command: %w", err)
	}

	var devices []types.USBDevice
	lines := strings.Split(string(output), "\n")

	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" || strings.Contains(line, "Node") {
			continue
		}

		fields := strings.Split(line, ",")
		if len(fields) < 5 {
			continue
		}

		drive := fields[1]
		sizeStr := fields[2]
		label := fields[4]

		if drive == "" {
			continue
		}

		size, err := strconv.ParseInt(sizeStr, 10, 64)
		if err != nil {
			udw.logger.Warn("Failed to parse device size", "device", drive, "size", sizeStr, "error", err)
			continue
		}

		device := types.USBDevice{
			Path:        drive,
			RawPath:     drive,
			Capacity:    size,
			Speed:       0,
			Vendor:      "",
			Model:       label,
			Serial:      "",
			IsSystem:    false,
			IsRemovable: true,
		}

		devices = append(devices, device)
	}

	return devices, nil
}

// getDeviceInfoWithPowerShell gets device info using PowerShell
func (udw *USBDetectorWindows) getDeviceInfoWithPowerShell(ctx context.Context, devicePath string) (*types.USBDevice, error) {
	// PowerShell command to get detailed device info
	cmd := exec.CommandContext(ctx, "powershell", "-Command",
		fmt.Sprintf(`
			$drive = Get-WmiObject -Class Win32_LogicalDisk -Filter "DeviceID='%s'"
			if ($drive) {
				Write-Output "$($drive.DeviceID)|$($drive.Size)|$($drive.FreeSpace)|$($drive.VolumeLabel)|$($drive.DriveType)"
			}
		`, devicePath))

	output, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("failed to get device info: %w", err)
	}

	line := strings.TrimSpace(string(output))
	if line == "" {
		return nil, fmt.Errorf("device not found: %s", devicePath)
	}

	fields := strings.Split(line, "|")
	if len(fields) < 5 {
		return nil, fmt.Errorf("invalid device info format")
	}

	size, err := strconv.ParseInt(fields[1], 10, 64)
	if err != nil {
		return nil, fmt.Errorf("failed to parse device size: %w", err)
	}

	driveType, err := strconv.Atoi(fields[4])
	if err != nil {
		return nil, fmt.Errorf("failed to parse drive type: %w", err)
	}

	device := &types.USBDevice{
		Path:        devicePath,
		RawPath:     devicePath,
		Capacity:    size,
		Speed:       0,
		Vendor:      "",
		Model:       fields[3],
		Serial:      "",
		IsSystem:    false,
		IsRemovable: driveType == 2,
	}

	return device, nil
}

// getDeviceInfoWithWMIC gets device info using WMIC (fallback)
func (udw *USBDetectorWindows) getDeviceInfoWithWMIC(ctx context.Context, devicePath string) (*types.USBDevice, error) {
	// WMIC command to get device info
	cmd := exec.CommandContext(ctx, "wmic", "logicaldisk", "where", fmt.Sprintf("deviceid='%s'", devicePath), "get", "deviceid,size,freespace,volumename,drivetype", "/format:csv")

	output, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("failed to get device info: %w", err)
	}

	lines := strings.Split(string(output), "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" || strings.Contains(line, "Node") {
			continue
		}

		fields := strings.Split(line, ",")
		if len(fields) >= 5 && fields[1] == devicePath {
			size, err := strconv.ParseInt(fields[2], 10, 64)
			if err != nil {
				return nil, fmt.Errorf("failed to parse device size: %w", err)
			}

			driveType, err := strconv.Atoi(fields[5])
			if err != nil {
				return nil, fmt.Errorf("failed to parse drive type: %w", err)
			}

			device := &types.USBDevice{
				Path:        devicePath,
				RawPath:     devicePath,
				Capacity:    size,
				Speed:       0,
				Vendor:      "",
				Model:       fields[4],
				Serial:      "",
				IsSystem:    false,
				IsRemovable: driveType == 2,
			}

			return device, nil
		}
	}

	return nil, fmt.Errorf("device not found: %s", devicePath)
}

// getUSBDeviceInfoWithPowerShell gets USB-specific device information
func (udw *USBDetectorWindows) getUSBDeviceInfoWithPowerShell(ctx context.Context, devicePath string) (*USBDeviceInfo, error) {
	// PowerShell command to get USB device info
	cmd := exec.CommandContext(ctx, "powershell", "-Command",
		fmt.Sprintf(`
			$drive = Get-WmiObject -Class Win32_LogicalDisk -Filter "DeviceID='%s'"
			if ($drive) {
				$physicalDrive = $drive.DeviceID.Replace(':', '')
				$disk = Get-WmiObject -Class Win32_DiskDrive | Where-Object {$_.DeviceID -like "*$physicalDrive*"}
				if ($disk) {
					Write-Output "$($disk.Manufacturer)|$($disk.Model)|$($disk.SerialNumber)|$($disk.SCSIBus)"
				}
			}
		`, devicePath))

	output, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("failed to get USB device info: %w", err)
	}

	line := strings.TrimSpace(string(output))
	if line == "" {
		return &USBDeviceInfo{}, nil
	}

	fields := strings.Split(line, "|")
	if len(fields) < 4 {
		return &USBDeviceInfo{}, nil
	}

	speed := 0
	if fields[3] != "" {
		if speedInt, err := strconv.Atoi(fields[3]); err == nil {
			speed = speedInt
		}
	}

	return &USBDeviceInfo{
		Vendor: strings.TrimSpace(fields[0]),
		Model:  strings.TrimSpace(fields[1]),
		Serial: strings.TrimSpace(fields[2]),
		Speed:  speed,
	}, nil
}

// getSystemDrive gets the system drive (usually C:)
func (udw *USBDetectorWindows) getSystemDrive(ctx context.Context) (string, error) {
	cmd := exec.CommandContext(ctx, "powershell", "-Command", "$env:SystemDrive")
	output, err := cmd.Output()
	if err != nil {
		return "", fmt.Errorf("failed to get system drive: %w", err)
	}

	systemDrive := strings.TrimSpace(string(output))
	if systemDrive == "" {
		return "C:", nil // Default fallback
	}

	return systemDrive, nil
}

// getBootDrive gets the boot drive
func (udw *USBDetectorWindows) getBootDrive(ctx context.Context) (string, error) {
	cmd := exec.CommandContext(ctx, "powershell", "-Command",
		"Get-WmiObject -Class Win32_OperatingSystem | Select-Object -ExpandProperty SystemDrive")
	output, err := cmd.Output()
	if err != nil {
		return "", fmt.Errorf("failed to get boot drive: %w", err)
	}

	bootDrive := strings.TrimSpace(string(output))
	if bootDrive == "" {
		return "C:", nil // Default fallback
	}

	return bootDrive, nil
}

// getFixedDrives gets all fixed drives
func (udw *USBDetectorWindows) getFixedDrives(ctx context.Context) ([]string, error) {
	cmd := exec.CommandContext(ctx, "powershell", "-Command",
		"Get-WmiObject -Class Win32_LogicalDisk | Where-Object {$_.DriveType -eq 3} | Select-Object -ExpandProperty DeviceID")
	output, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("failed to get fixed drives: %w", err)
	}

	var drives []string
	lines := strings.Split(string(output), "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line != "" {
			drives = append(drives, line)
		}
	}

	return drives, nil
}

// USBDeviceInfo represents USB-specific device information
type USBDeviceInfo struct {
	Vendor string
	Model  string
	Serial string
	Speed  int
}
