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

// platformSpecificValidation executes additional Windows-specific validation
func (udw *USBDetectorWindows) platformSpecificValidation(ctx context.Context, device types.USBDevice) error {
	udw.logger.Debug("Windows-specific validation", "device", device.Path)

	// Extract disk number from path
	var diskNum int
	var err error

	if strings.Contains(device.Path, "PHYSICALDRIVE") {
		// Direct PHYSICALDRIVE path
		fmt.Sscanf(device.Path, "\\\\.\\PHYSICALDRIVE%d", &diskNum)
	} else if strings.HasSuffix(device.Path, ":") {
		// Drive letter (e.g., "E:") - convert to PHYSICALDRIVE
		physicalDrive, err := udw.GetPhysicalDriveFromLetter(ctx, device.Path)
		if err != nil {
			return fmt.Errorf("failed to convert drive letter to physical drive: %w", err)
		}
		fmt.Sscanf(physicalDrive, "\\\\.\\PHYSICALDRIVE%d", &diskNum)
	} else {
		return fmt.Errorf("invalid Windows device path: %s (expected PHYSICALDRIVE or drive letter)", device.Path)
	}

	// PowerShell script for independent validation
	psScript := fmt.Sprintf(`
	$disk = Get-Disk -Number %d -ErrorAction SilentlyContinue
	if (-not $disk) { exit 1 }
	
	# Critical checks
	if ($disk.IsSystem) { exit 2 }
	if ($disk.IsBoot) { exit 3 }
	if ($disk.IsOffline) { exit 4 }
	
	# Check partitions
	$partitions = Get-Partition -DiskNumber %d -ErrorAction SilentlyContinue
	foreach ($part in $partitions) {
		if ($part.DriveLetter -eq "C") { exit 5 }
		if ($part.IsSystem) { exit 6 }
		if ($part.IsBoot) { exit 7 }
	}
	
	# Check if has Windows installed
	$volumes = Get-Volume | Where-Object { $_.DriveLetter -ne $null }
	foreach ($vol in $volumes) {
		$partition = Get-Partition -DriveLetter $vol.DriveLetter -ErrorAction SilentlyContinue
		if ($partition -and $partition.DiskNumber -eq %d) {
			$winDir = $vol.DriveLetter + ":\Windows"
			if (Test-Path $winDir) { exit 8 }
		}
	}
	
	exit 0
	`, diskNum, diskNum, diskNum)

	cmd := exec.CommandContext(ctx, "powershell.exe", "-NoProfile", "-Command", psScript)
	output, err := cmd.CombinedOutput()

	if err != nil {
		exitCode := cmd.ProcessState.ExitCode()
		switch exitCode {
		case 2:
			return fmt.Errorf("device is system disk")
		case 3:
			return fmt.Errorf("device is boot disk")
		case 4:
			return fmt.Errorf("device is offline")
		case 5:
			return fmt.Errorf("device contains C: drive")
		case 6:
			return fmt.Errorf("device contains system partition")
		case 7:
			return fmt.Errorf("device contains boot partition")
		case 8:
			return fmt.Errorf("device contains Windows installation")
		default:
			return fmt.Errorf("validation failed: %s", string(output))
		}
	}

	return nil
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
	// PowerShell melhorado - nunca retorna valores null/vazios
	cmd := exec.CommandContext(ctx, "powershell", "-Command", `
		Get-WmiObject -Class Win32_LogicalDisk | 
		Where-Object {$_.DriveType -eq 2} | 
		ForEach-Object {
			$drive = $_.DeviceID
			$size = if ($null -ne $_.Size) { $_.Size } else { 0 }
			$free = if ($null -ne $_.FreeSpace) { $_.FreeSpace } else { 0 }
			$label = if ($_.VolumeLabel) { $_.VolumeLabel } else { "Unknown" }
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
			udw.logger.Warn("Invalid device info format", "line", line)
			continue
		}

		drive := strings.TrimSpace(fields[0])
		sizeStr := strings.TrimSpace(fields[1])
		label := strings.TrimSpace(fields[3])

		// Parse size com tratamento robusto
		size := int64(0)
		if sizeStr != "" && sizeStr != "0" {
			if parsedSize, err := strconv.ParseInt(sizeStr, 10, 64); err == nil {
				size = parsedSize
			} else {
				udw.logger.Warn("Failed to parse device size",
					"device", drive,
					"size_str", sizeStr,
					"error", err)
			}
		}

		// Se size ainda é 0, tentar método alternativo
		if size == 0 {
			udw.logger.Info("Device reported zero size, trying alternative detection", "device", drive)
			if altSize, err := udw.getDeviceSizeAlternative(ctx, drive); err == nil && altSize > 0 {
				size = altSize
				udw.logger.Info("Successfully obtained size via alternative method",
					"device", drive,
					"size", size)
			} else {
				udw.logger.Warn("Could not determine device size, device may be unformatted or faulty",
					"device", drive,
					"error", err)
				// Ainda assim adicionar o dispositivo, mas com aviso
			}
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
		udw.logger.Debug("Detected USB device",
			"path", device.Path,
			"capacity", device.Capacity,
			"model", label)
	}

	if len(devices) == 0 {
		udw.logger.Warn("No removable USB devices found via PowerShell")
		return nil, fmt.Errorf("no removable USB devices detected")
	}

	udw.logger.Info("Successfully detected USB devices", "count", len(devices))
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

// getDeviceSizeFromDisk obtém tamanho do disco físico associado à letra
func (udw *USBDetectorWindows) getDeviceSizeFromDisk(ctx context.Context, driveLetter string) (int64, error) {
	psScript := fmt.Sprintf(`
		$drive = "%s"
		
		# Tentar via Get-Disk
		try {
			$partition = Get-Partition | Where-Object { $_.DriveLetter -eq $drive.TrimEnd(':') } -ErrorAction Stop
			if ($partition) {
				$disk = Get-Disk -Number $partition.DiskNumber -ErrorAction Stop
				if ($disk) {
					Write-Output $disk.Size
					exit 0
				}
			}
		} catch {}
		
		# Fallback: via WMI DiskDrive
		try {
			$logicalDisk = Get-WmiObject -Class Win32_LogicalDisk -Filter "DeviceID='$drive'" -ErrorAction Stop
			$partition = Get-WmiObject -Query "ASSOCIATORS OF {Win32_LogicalDisk.DeviceID='$drive'} WHERE AssocClass = Win32_LogicalDiskToPartition" -ErrorAction Stop
			if ($partition) {
				$disk = Get-WmiObject -Query "ASSOCIATORS OF {Win32_DiskPartition.DeviceID='$($partition.DeviceID)'} WHERE AssocClass = Win32_DiskDriveToDiskPartition" -ErrorAction Stop
				if ($disk) {
					Write-Output $disk.Size
					exit 0
				}
			}
		} catch {}
		
		Write-Output "0"
	`, driveLetter)

	cmd := exec.CommandContext(ctx, "powershell", "-NoProfile", "-Command", psScript)
	output, err := cmd.Output()
	if err != nil {
		return 0, fmt.Errorf("failed to get disk size: %w", err)
	}

	sizeStr := strings.TrimSpace(string(output))
	if sizeStr == "" || sizeStr == "0" {
		return 0, fmt.Errorf("no size information available")
	}

	size, err := strconv.ParseInt(sizeStr, 10, 64)
	if err != nil {
		return 0, fmt.Errorf("failed to parse size: %w", err)
	}

	return size, nil
}

// getDeviceSizeAlternative é um wrapper que tenta múltiplos métodos
func (udw *USBDetectorWindows) getDeviceSizeAlternative(ctx context.Context, driveLetter string) (int64, error) {
	// Método 1: Via Get-Disk
	if size, err := udw.getDeviceSizeFromDisk(ctx, driveLetter); err == nil && size > 0 {
		return size, nil
	}

	// Método 2: Via WMIC
	cmd := exec.CommandContext(ctx, "wmic", "logicaldisk", "where",
		fmt.Sprintf("deviceid='%s'", driveLetter),
		"get", "size", "/format:value")

	output, err := cmd.Output()
	if err == nil {
		lines := strings.Split(string(output), "\n")
		for _, line := range lines {
			if strings.HasPrefix(line, "Size=") {
				sizeStr := strings.TrimPrefix(line, "Size=")
				sizeStr = strings.TrimSpace(sizeStr)
				if size, err := strconv.ParseInt(sizeStr, 10, 64); err == nil {
					return size, nil
				}
			}
		}
	}

	return 0, fmt.Errorf("all alternative methods failed to get device size")
}

// USBDeviceInfo represents USB-specific device information
type USBDeviceInfo struct {
	Vendor string
	Model  string
	Serial string
	Speed  int
}

// GetPhysicalDriveFromLetter converts a drive letter to PhysicalDrive path
func (udw *USBDetectorWindows) GetPhysicalDriveFromLetter(ctx context.Context, driveLetter string) (string, error) {
	// PowerShell command to map drive letter to physical drive number
	cmd := exec.CommandContext(ctx, "powershell", "-Command", fmt.Sprintf(`
		$drive = Get-WmiObject -Class Win32_LogicalDisk -Filter "DeviceID='%s'"
		if ($drive) {
			$partition = Get-WmiObject -Query "ASSOCIATORS OF {Win32_LogicalDisk.DeviceID='%s'} WHERE AssocClass = Win32_LogicalDiskToPartition"
			if ($partition) {
				$disk = Get-WmiObject -Query "ASSOCIATORS OF {Win32_DiskPartition.DeviceID='$($partition.DeviceID)'} WHERE AssocClass = Win32_DiskDriveToDiskPartition"
				if ($disk) {
					Write-Output $disk.DeviceID
				}
			}
		}
	`, driveLetter, driveLetter))

	output, err := cmd.Output()
	if err != nil {
		return "", fmt.Errorf("failed to map drive letter to physical drive: %w", err)
	}

	physicalDrive := strings.TrimSpace(string(output))
	if physicalDrive == "" {
		return "", fmt.Errorf("could not map drive %s to physical drive", driveLetter)
	}

	return physicalDrive, nil
}

// createWindowsDetector creates Windows detector - overrides stub in usb_detector.go
func createWindowsDetector(logger types.Logger) USBDetector {
	return NewUSBDetectorWindows(logger)
}

// createLinuxDetector stub for Linux detector (not available on Windows)
func createLinuxDetector(logger types.Logger) USBDetector {
	return nil
}

// createMacOSDetector stub for macOS detector (not available on Windows)
func createMacOSDetector(logger types.Logger) USBDetector {
	return nil
}

// ValidateDeviceDoubleCheck overrides the base implementation to use Windows-specific validation
func (udw *USBDetectorWindows) ValidateDeviceDoubleCheck(ctx context.Context, device types.USBDevice) error {
	udw.logger.Info("Starting Windows double validation check", "device", device.Path)

	// Validation 1: Check system flags
	if device.IsSystem {
		return fmt.Errorf("SECURITY: device %s is marked as system device", device.Path)
	}

	// Validation 2: Check if not removable
	if !device.IsRemovable {
		return fmt.Errorf("SECURITY: device %s is not removable", device.Path)
	}

	// Validation 3: Check suspicious capacity (system disks are usually > 100GB)
	minCapacity := int64(512 * 1024 * 1024)             // 512 MB
	maxCapacity := int64(2 * 1024 * 1024 * 1024 * 1024) // 2 TB

	// Permitir capacity = 0 para dispositivos não formatados ou com problemas de detecção
	if device.Capacity > 0 {
		if device.Capacity < minCapacity {
			return fmt.Errorf("SECURITY: device %s is too small (%d bytes)", device.Path, device.Capacity)
		}

		if device.Capacity > maxCapacity {
			return fmt.Errorf("SECURITY: device %s is suspiciously large (%d bytes)", device.Path, device.Capacity)
		}
	} else {
		udw.logger.Warn("Device has zero capacity - may be unformatted or faulty", "device", device.Path)
	}

	udw.logger.Info("First validation passed", "device", device.Path)

	// Independent validation via Windows-specific method
	if err := udw.platformSpecificValidation(ctx, device); err != nil {
		return fmt.Errorf("SECURITY: platform-specific validation failed: %w", err)
	}

	udw.logger.Info("Windows double validation check passed", "device", device.Path)
	return nil
}
