//go:build windows
// +build windows

package node

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"syscall"
	"time"

	"node-component/src/internal/types"
)

// USBWriterWindows implements USB writing for Windows systems
type USBWriterWindows struct {
	*USBWriterBase
}

// NewUSBWriterWindows creates a new Windows USB writer
func NewUSBWriterWindows(logger types.Logger) *USBWriterWindows {
	base := NewUSBWriterBase("windows", logger)
	return &USBWriterWindows{
		USBWriterBase: base,
	}
}

// WriteISO writes an ISO to a USB device with cloud-init injection
func (uww *USBWriterWindows) WriteISO(ctx context.Context, isoPath string, selectedDevice *SelectedUSBDevice, cloudInitConfig *types.CloudInitConfig) (*types.WriteResult, error) {
	uww.logger.Info("Writing ISO to USB device",
		"iso", isoPath,
		"device", selectedDevice.Device.Path,
		"token", selectedDevice.ValidationToken)

	startTime := time.Now()

	// Log detailed device information
	uww.logger.Info("Device details",
		"path", selectedDevice.Device.Path,
		"capacity", selectedDevice.Device.Capacity,
		"model", selectedDevice.Device.Model,
		"serial", selectedDevice.Device.Serial,
		"validation_token", selectedDevice.ValidationToken)

	// Validate inputs using the selected device
	if err := uww.validateInputs(isoPath, selectedDevice.Device.Path); err != nil {
		return nil, fmt.Errorf("input validation failed: %w", err)
	}

	// Check required tools before starting
	if err := uww.checkRequiredTools(ctx); err != nil {
		return nil, fmt.Errorf("required tools check failed: %w", err)
	}

	// Get ISO file size
	isoFile, err := os.Open(isoPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open ISO file: %w", err)
	}
	defer isoFile.Close()

	isoInfo, err := isoFile.Stat()
	if err != nil {
		return nil, fmt.Errorf("failed to get ISO file info: %w", err)
	}

	// Create progress tracker
	progressTracker := NewWriteProgressTracker(uww.logger, isoInfo.Size())

	// Write ISO to device using dd or diskpart
	actualISOPath := isoPath

	// Write ISO to device using multiple fallback methods
	if err := uww.writeISOWithFallback(ctx, actualISOPath, selectedDevice, progressTracker); err != nil {
		return &types.WriteResult{
			DevicePath:   selectedDevice.Device.Path,
			ISOPath:      isoPath,
			BytesWritten: 0,
			Duration:     time.Since(startTime),
			Success:      false,
			ErrorMessage: err.Error(),
		}, nil
	}

	// Create NoCloud partition for cloud-init if provided
	if cloudInitConfig != nil {
		uww.logger.Info("Creating NoCloud partition for cloud-init")
		cloudInitInjector := NewCloudInitInjector(uww.logger)
		if err := cloudInitInjector.CreateNoCloudPartition(selectedDevice.Device.Path, cloudInitConfig); err != nil {
			uww.logger.Error("Failed to create NoCloud partition", "error", err)
			return nil, fmt.Errorf("failed to create cloud-init partition: %w", err)
		}
	}

	// Clean up injected ISO if created
	if actualISOPath != isoPath {
		os.Remove(actualISOPath)
	}

	duration := time.Since(startTime)
	uww.logger.Info("ISO write completed successfully",
		"device", selectedDevice.Device.Path,
		"duration", duration,
		"bytes_written", isoInfo.Size())

	return &types.WriteResult{
		DevicePath:        selectedDevice.Device.Path,
		ISOPath:           isoPath,
		BytesWritten:      isoInfo.Size(),
		Duration:          duration,
		Success:           true,
		CloudInitInjected: cloudInitConfig != nil,
	}, nil
}

// UnmountDevice unmounts a device before writing
func (uww *USBWriterWindows) UnmountDevice(ctx context.Context, devicePath string) error {
	uww.logger.Debug("Unmounting device", "device", devicePath)

	// Get drive letter from device path (e.g., "C:" from "C:\")
	driveLetter := strings.TrimSuffix(devicePath, ":")
	if driveLetter == "" {
		return fmt.Errorf("invalid device path: %s", devicePath)
	}

	// Try PowerShell method first (more reliable)
	if err := uww.unmountWithPowerShell(ctx, devicePath); err == nil {
		uww.logger.Debug("Device unmounted successfully with PowerShell", "device", devicePath)
		return nil
	}

	// Fallback to diskpart if available
	if err := uww.unmountWithDiskpart(ctx, devicePath); err != nil {
		uww.logger.Warn("Failed to unmount device with both PowerShell and diskpart", "device", devicePath, "error", err)
		fmt.Printf("⚠️  Aviso: Falha ao desmontar dispositivo: %v\n", err)
		fmt.Println("   ℹ️  Continuando com a operação...")
		// Continue even if unmount fails
	}

	uww.logger.Debug("Device unmounted successfully", "device", devicePath)
	return nil
}

// MountDevice mounts a device after writing
func (uww *USBWriterWindows) MountDevice(ctx context.Context, devicePath string) error {
	uww.logger.Debug("Mounting device", "device", devicePath)

	// Get drive letter from device path
	driveLetter := strings.TrimSuffix(devicePath, ":")
	if driveLetter == "" {
		return fmt.Errorf("invalid device path: %s", devicePath)
	}

	// Try PowerShell method first (more reliable)
	if err := uww.mountWithPowerShell(ctx, devicePath); err == nil {
		uww.logger.Debug("Device mounted successfully with PowerShell", "device", devicePath)
		return nil
	}

	// Fallback to diskpart if available
	if err := uww.mountWithDiskpart(ctx, devicePath); err != nil {
		uww.logger.Warn("Failed to mount device with both PowerShell and diskpart", "device", devicePath, "error", err)
		return nil // Continue even if mount fails
	}

	uww.logger.Debug("Device mounted successfully", "device", devicePath)
	return nil
}

// Private helper methods

// unmountWithPowerShell unmounts device using PowerShell
func (uww *USBWriterWindows) unmountWithPowerShell(ctx context.Context, devicePath string) error {
	driveLetter := strings.TrimSuffix(devicePath, ":")

	psScript := fmt.Sprintf(`
try {
    $volume = Get-Volume -DriveLetter %s -ErrorAction SilentlyContinue
    if ($volume) {
        $disk = Get-Disk -Number $volume.DiskNumber
        if ($disk) {
            Set-Disk -Number $disk.Number -IsOffline $true
            Set-Disk -Number $disk.Number -IsOffline $false
            Write-Output "Device unmounted successfully"
        } else {
            Write-Error "Disk not found for volume"
            exit 1
        }
    } else {
        Write-Error "Volume not found"
        exit 1
    }
} catch {
    Write-Error "Failed to unmount device: $_"
    exit 1
}
`, driveLetter)

	cmd := exec.CommandContext(ctx, "powershell", "-Command", psScript)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("PowerShell unmount failed: %w - %s", err, string(output))
	}
	return nil
}

// unmountWithDiskpart unmounts device using diskpart (fallback)
func (uww *USBWriterWindows) unmountWithDiskpart(ctx context.Context, devicePath string) error {
	driveLetter := strings.TrimSuffix(devicePath, ":")

	diskpartScript := fmt.Sprintf(`
select disk %s
offline disk
online disk
`, driveLetter)

	scriptPath := filepath.Join(os.TempDir(), "syntropy-unmount.txt")
	if err := os.WriteFile(scriptPath, []byte(diskpartScript), 0644); err != nil {
		return fmt.Errorf("failed to create diskpart script: %w", err)
	}
	defer os.Remove(scriptPath)

	cmd := exec.CommandContext(ctx, "diskpart", "/s", scriptPath)
	output, err := cmd.CombinedOutput()
	if err != nil {
		errorMsg := uww.parseDiskpartError(err, string(output), "unmount")
		return fmt.Errorf("diskpart unmount failed: %s", errorMsg)
	}
	return nil
}

// mountWithPowerShell mounts device using PowerShell
func (uww *USBWriterWindows) mountWithPowerShell(ctx context.Context, devicePath string) error {
	driveLetter := strings.TrimSuffix(devicePath, ":")

	psScript := fmt.Sprintf(`
try {
    $volume = Get-Volume -DriveLetter %s -ErrorAction SilentlyContinue
    if ($volume) {
        $disk = Get-Disk -Number $volume.DiskNumber
        if ($disk) {
            Set-Disk -Number $disk.Number -IsOffline $false
            Set-Disk -Number $disk.Number -IsReadOnly $false
            Write-Output "Device mounted successfully"
        } else {
            Write-Error "Disk not found for volume"
            exit 1
        }
    } else {
        Write-Error "Volume not found"
        exit 1
    }
} catch {
    Write-Error "Failed to mount device: $_"
    exit 1
}
`, driveLetter)

	cmd := exec.CommandContext(ctx, "powershell", "-Command", psScript)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("PowerShell mount failed: %w - %s", err, string(output))
	}
	return nil
}

// mountWithDiskpart mounts device using diskpart (fallback)
func (uww *USBWriterWindows) mountWithDiskpart(ctx context.Context, devicePath string) error {
	driveLetter := strings.TrimSuffix(devicePath, ":")

	diskpartScript := fmt.Sprintf(`
select disk %s
online disk
attributes disk clear readonly
`, driveLetter)

	scriptPath := filepath.Join(os.TempDir(), "syntropy-mount.txt")
	if err := os.WriteFile(scriptPath, []byte(diskpartScript), 0644); err != nil {
		return fmt.Errorf("failed to create diskpart script: %w", err)
	}
	defer os.Remove(scriptPath)

	cmd := exec.CommandContext(ctx, "diskpart", "/s", scriptPath)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("diskpart mount failed: %w - %s", err, string(output))
	}
	return nil
}

// validateInputs validates the input parameters
func (uww *USBWriterWindows) validateInputs(isoPath string, devicePath string) error {
	// Validate ISO file
	if _, err := os.Stat(isoPath); err != nil {
		return fmt.Errorf("ISO file does not exist: %s", isoPath)
	}

	// Validate device path format (should be like "C:")
	if !strings.HasSuffix(devicePath, ":") {
		return fmt.Errorf("invalid device path format: %s (should be like 'C:')", devicePath)
	}

	// Check if drive exists
	cmd := exec.Command("wmic", "logicaldisk", "where", fmt.Sprintf("deviceid='%s'", devicePath), "get", "deviceid")
	output, err := cmd.Output()
	if err != nil {
		return fmt.Errorf("failed to check if drive exists: %w", err)
	}

	if !strings.Contains(string(output), devicePath) {
		return fmt.Errorf("drive does not exist: %s", devicePath)
	}

	return nil
}

// writeISOWithFallback writes ISO using multiple fallback methods
func (uww *USBWriterWindows) writeISOWithFallback(ctx context.Context, isoPath string, selectedDevice *SelectedUSBDevice, progressTracker *WriteProgressTracker) error {
	uww.logger.Debug("Writing ISO with fallback methods",
		"iso", isoPath,
		"device", selectedDevice.Device.Path,
		"token", selectedDevice.ValidationToken)

	// Method 1: Try dd first (fastest if available)
	fmt.Printf("   🔄 Método 1: dd (Device: %s, Token: %s)\n", selectedDevice.Device.Path, selectedDevice.ValidationToken[:8])
	if err := uww.writeISOWithDD(ctx, isoPath, selectedDevice, progressTracker); err == nil {
		fmt.Println("   ✅ Sucesso com dd!")
		return nil
	}

	// Method 2: Try PowerShell with robocopy
	fmt.Printf("   🔄 Método 2: PowerShell + robocopy (Device: %s, Token: %s)\n", selectedDevice.Device.Path, selectedDevice.ValidationToken[:8])
	if err := uww.writeISOWithPowerShell(ctx, isoPath, selectedDevice, progressTracker); err == nil {
		fmt.Println("   ✅ Sucesso com PowerShell!")
		return nil
	}

	// Method 3: Try PowerShell with direct copy
	fmt.Printf("   🔄 Método 3: PowerShell + copy (Device: %s, Token: %s)\n", selectedDevice.Device.Path, selectedDevice.ValidationToken[:8])
	if err := uww.writeISOWithPowerShellCopy(ctx, isoPath, selectedDevice, progressTracker); err == nil {
		fmt.Println("   ✅ Sucesso com PowerShell copy!")
		return nil
	}

	// Method 4: Try diskpart + copy
	fmt.Printf("   🔄 Método 4: diskpart + copy (Device: %s, Token: %s)\n", selectedDevice.Device.Path, selectedDevice.ValidationToken[:8])
	if err := uww.writeISOWithDiskpart(ctx, isoPath, selectedDevice, progressTracker); err == nil {
		fmt.Println("   ✅ Sucesso com diskpart!")
		return nil
	}

	return fmt.Errorf("todos os métodos de escrita falharam - verifique permissões e conectividade do dispositivo")
}

// writeISOWithDD writes ISO to device using dd command (if available) or alternative method
func (uww *USBWriterWindows) writeISOWithDD(ctx context.Context, isoPath string, selectedDevice *SelectedUSBDevice, progressTracker *WriteProgressTracker) error {
	uww.logger.Debug("Writing ISO with dd", "iso", isoPath, "device", selectedDevice.Device.Path, "token", selectedDevice.ValidationToken)

	// Convert drive letter to PhysicalDrive if needed
	actualDevicePath := selectedDevice.Device.Path
	if strings.HasSuffix(selectedDevice.Device.Path, ":") && len(selectedDevice.Device.Path) == 2 {
		factory := NewUSBDetectorFactory()
		detector, err := factory.CreateUSBDetector(uww.logger)
		if err != nil {
			return fmt.Errorf("failed to create detector: %w", err)
		}

		windowsDetector, ok := detector.(*USBDetectorWindows)
		if !ok {
			return fmt.Errorf("expected Windows detector")
		}

		physicalDrive, err := windowsDetector.GetPhysicalDriveFromLetter(ctx, selectedDevice.Device.Path)
		if err != nil {
			return fmt.Errorf("failed to get physical drive: %w", err)
		}
		actualDevicePath = physicalDrive
	}

	// Try to use dd if available (common in Windows with WSL or Git Bash)
	cmd := exec.CommandContext(ctx, "dd",
		"if="+isoPath,
		"of="+actualDevicePath,
		"bs=4M",
		"status=progress",
		"conv=fsync")

	output, err := cmd.CombinedOutput()
	if err != nil {
		// Parse dd error for better user experience
		errorMsg := uww.parseDDError(err, string(output))
		uww.logger.Warn("dd command failed, trying alternative method", "error", errorMsg)
		fmt.Printf("⚠️  Aviso: Comando dd falhou: %s\n", errorMsg)
		fmt.Println("   🔄 Tentando método alternativo com PowerShell...")
		fmt.Println()

		// Fallback to PowerShell method
		return uww.writeISOWithPowerShell(ctx, isoPath, selectedDevice, progressTracker)
	}

	uww.logger.Debug("dd command completed successfully")
	return nil
}

// writeISOWithPowerShellCopy writes ISO using PowerShell with direct copy
func (uww *USBWriterWindows) writeISOWithPowerShellCopy(ctx context.Context, isoPath string, selectedDevice *SelectedUSBDevice, progressTracker *WriteProgressTracker) error {
	uww.logger.Debug("Writing ISO with PowerShell copy", "iso", isoPath, "device", selectedDevice.Device.Path)

	// PowerShell script to write ISO using direct copy
	psScript := fmt.Sprintf(`
$isoPath = "%s"
$selectedDevice.Device.Path = "%s"

try {
    # Mount ISO
    $isoMount = Mount-DiskImage -ImagePath $isoPath -PassThru
    $isoDrive = ($isoMount | Get-Volume).DriveLetter + ":"
    
    # Copy ISO contents to USB device using Copy-Item
    Copy-Item -Path "$isoDrive\*" -Destination $selectedDevice.Device.Path -Recurse -Force
    
    # Dismount ISO
    Dismount-DiskImage -ImagePath $isoPath
    
    Write-Output "ISO copy completed successfully"
} catch {
    Write-Error "Failed to copy ISO: $_"
    exit 1
}
`, isoPath, selectedDevice.Device.Path)

	// Execute PowerShell script
	cmd := exec.CommandContext(ctx, "powershell", "-Command", psScript)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("PowerShell copy failed: %w - %s", err, string(output))
	}

	uww.logger.Debug("PowerShell copy completed successfully")
	return nil
}

// writeISOWithDiskpart writes ISO using diskpart and copy
func (uww *USBWriterWindows) writeISOWithDiskpart(ctx context.Context, isoPath string, selectedDevice *SelectedUSBDevice, progressTracker *WriteProgressTracker) error {
	uww.logger.Debug("Writing ISO with diskpart", "iso", isoPath, "device", selectedDevice.Device.Path)

	// First, format the device with diskpart
	diskpartScript := fmt.Sprintf(`
select disk %s
clean
create partition primary
active
format fs=fat32 quick
assign
`, strings.TrimSuffix(selectedDevice.Device.Path, ":"))

	// Create temporary script file
	scriptPath := filepath.Join(os.TempDir(), "syntropy-format.txt")
	if err := os.WriteFile(scriptPath, []byte(diskpartScript), 0644); err != nil {
		return fmt.Errorf("failed to create diskpart script: %w", err)
	}
	defer os.Remove(scriptPath)

	// Execute diskpart
	cmd := exec.CommandContext(ctx, "diskpart", "/s", scriptPath)
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("diskpart format failed: %w", err)
	}

	// Then use PowerShell to copy ISO contents
	return uww.writeISOWithPowerShellCopy(ctx, isoPath, selectedDevice, progressTracker)
}

// writeISOWithPowerShell writes ISO using PowerShell as fallback
func (uww *USBWriterWindows) writeISOWithPowerShell(ctx context.Context, isoPath string, selectedDevice *SelectedUSBDevice, progressTracker *WriteProgressTracker) error {
	uww.logger.Debug("Writing ISO with PowerShell", "iso", isoPath, "device", selectedDevice.Device.Path)

	// PowerShell script to write ISO to USB
	psScript := fmt.Sprintf(`
$isoPath = "%s"
$selectedDevice.Device.Path = "%s"

try {
    # Mount ISO
    $isoMount = Mount-DiskImage -ImagePath $isoPath -PassThru
    $isoDrive = ($isoMount | Get-Volume).DriveLetter + ":"
    
    # Copy ISO contents to USB device
    robocopy $isoDrive $selectedDevice.Device.Path /E /COPY:DAT /R:3 /W:1
    
    # Dismount ISO
    Dismount-DiskImage -ImagePath $isoPath
    
    Write-Output "ISO write completed successfully"
} catch {
    Write-Error "Failed to write ISO: $_"
    exit 1
}
`, isoPath, selectedDevice.Device.Path)

	// Execute PowerShell script
	cmd := exec.CommandContext(ctx, "powershell", "-Command", psScript)
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("PowerShell ISO write failed: %w", err)
	}

	uww.logger.Debug("PowerShell ISO write completed successfully")
	return nil
}

// ValidateDevice validates if a device is suitable for writing (override base)
func (uww *USBWriterWindows) ValidateDevice(ctx context.Context, devicePath string) error {
	uww.logger.Debug("Validating Windows device", "device", devicePath)

	// Convert drive letter to PhysicalDrive if needed
	actualDevicePath := devicePath
	if strings.HasSuffix(devicePath, ":") && len(devicePath) == 2 {
		// This is a drive letter, convert to PhysicalDrive
		factory := NewUSBDetectorFactory()
		detector, err := factory.CreateUSBDetector(uww.logger)
		if err != nil {
			return fmt.Errorf("failed to create detector: %w", err)
		}

		windowsDetector, ok := detector.(*USBDetectorWindows)
		if !ok {
			return fmt.Errorf("expected Windows detector")
		}

		physicalDrive, err := windowsDetector.GetPhysicalDriveFromLetter(ctx, devicePath)
		if err != nil {
			return fmt.Errorf("failed to get physical drive: %w", err)
		}
		actualDevicePath = physicalDrive
	}

	// Validate physical drive format
	if !strings.Contains(actualDevicePath, "PHYSICALDRIVE") {
		return fmt.Errorf("invalid device path format: %s", actualDevicePath)
	}

	// Check if we have administrative privileges
	if !uww.hasAdminPrivileges(ctx) {
		return fmt.Errorf("administrative privileges required to write to USB device")
	}

	// Try to open the physical drive for writing
	// Note: This requires admin privileges
	handle, err := syscall.CreateFile(
		syscall.StringToUTF16Ptr(actualDevicePath),
		syscall.GENERIC_WRITE,
		syscall.FILE_SHARE_READ|syscall.FILE_SHARE_WRITE,
		nil,
		syscall.OPEN_EXISTING,
		0,
		0,
	)
	if err != nil {
		return fmt.Errorf("device is not writable (admin required): %s", devicePath)
	}
	syscall.CloseHandle(handle)

	uww.logger.Debug("Device validation passed", "device", devicePath)
	return nil
}

// checkRequiredTools checks if required tools are available for USB writing
func (uww *USBWriterWindows) checkRequiredTools(ctx context.Context) error {
	uww.logger.Debug("Checking required tools for Windows USB writing")

	// Check if PowerShell is available (always available on Windows)
	cmd := exec.CommandContext(ctx, "powershell", "-Command", "Get-Host")
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("PowerShell não está disponível - necessário para escrita USB")
	}

	// Check if diskpart is available (optional, will use PowerShell alternatives)
	cmd = exec.CommandContext(ctx, "diskpart")
	diskpartAvailable := cmd.Run() == nil

	if !diskpartAvailable {
		uww.logger.Info("diskpart command not available, will use PowerShell alternatives")
		fmt.Println("ℹ️  Comando diskpart não encontrado - usando alternativas PowerShell")
	}

	// Check if dd is available (optional, will fallback to PowerShell)
	cmd = exec.CommandContext(ctx, "dd", "--version")
	ddAvailable := cmd.Run() == nil

	if !ddAvailable {
		uww.logger.Info("dd command not available, will use PowerShell fallback")
		fmt.Println("ℹ️  Comando dd não encontrado - usando PowerShell como método alternativo")
	}

	uww.logger.Debug("Required tools check completed", "dd_available", ddAvailable, "diskpart_available", diskpartAvailable)
	return nil
}

// parseDDError parses dd error output for better user experience
func (uww *USBWriterWindows) parseDDError(err error, output string) string {
	// Check for common dd error patterns
	if strings.Contains(strings.ToLower(output), "executable file not found") {
		return "Comando dd não encontrado - instale WSL, Git Bash ou use PowerShell"
	}
	if strings.Contains(strings.ToLower(output), "permission denied") {
		return "Permissão negada - execute como administrador"
	}
	if strings.Contains(strings.ToLower(output), "no such file or directory") {
		return "Arquivo ou diretório não encontrado - verifique o caminho do dispositivo"
	}
	if strings.Contains(strings.ToLower(output), "device or resource busy") {
		return "Dispositivo está sendo usado - feche outros programas que possam estar acessando o USB"
	}
	if strings.Contains(strings.ToLower(output), "read-only file system") {
		return "Sistema de arquivos somente leitura - dispositivo pode estar protegido"
	}
	if strings.Contains(strings.ToLower(output), "no space left on device") {
		return "Sem espaço no dispositivo - verifique se o USB tem espaço suficiente"
	}
	if strings.Contains(strings.ToLower(output), "invalid argument") {
		return "Argumento inválido - dispositivo pode não suportar escrita direta"
	}

	// Return generic error with context
	return fmt.Sprintf("Erro do dd: %v", err)
}

// parseDiskpartError parses diskpart error output for better user experience
func (uww *USBWriterWindows) parseDiskpartError(err error, output, operation string) string {
	// Common diskpart error codes and their meanings
	errorMap := map[string]string{
		"0x80070057": "Parâmetro inválido - dispositivo pode estar em uso ou protegido",
		"0x80070005": "Acesso negado - execute como administrador",
		"0x8007001F": "Dispositivo não está pronto",
		"0x80070032": "Dispositivo não suporta o comando",
		"0x80070015": "Dispositivo não encontrado",
		"0x80070020": "Dispositivo está sendo usado por outro processo",
	}

	// Check for specific error codes in output
	for code, message := range errorMap {
		if strings.Contains(output, code) {
			return fmt.Sprintf("%s (Código: %s)", message, code)
		}
	}

	// Check for common error patterns
	if strings.Contains(strings.ToLower(output), "access denied") {
		return "Acesso negado - execute o comando como administrador"
	}
	if strings.Contains(strings.ToLower(output), "device not ready") {
		return "Dispositivo não está pronto - verifique se está conectado corretamente"
	}
	if strings.Contains(strings.ToLower(output), "device in use") {
		return "Dispositivo está sendo usado por outro processo - feche outros programas"
	}
	if strings.Contains(strings.ToLower(output), "invalid parameter") {
		return "Parâmetro inválido - dispositivo pode estar protegido ou em uso"
	}

	// Return generic error with operation context
	return fmt.Sprintf("Erro durante %s: %v", operation, err)
}

// hasAdminPrivileges checks if the current process has admin privileges
func (uww *USBWriterWindows) hasAdminPrivileges(ctx context.Context) bool {
	cmd := exec.CommandContext(ctx, "net", "session")
	if err := cmd.Run(); err != nil {
		return false
	}
	return true
}

// Windows-specific implementation will override the base method

// NewPlatformUSBWriter creates the platform-specific USB writer
func NewPlatformUSBWriter(logger types.Logger) USBWriter {
	return NewUSBWriterWindows(logger)
}
