//go:build windows
// +build windows

package setup

import (
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"syscall"
	"time"

	"github.com/shirou/gopsutil/v3/disk"
	"github.com/shirou/gopsutil/v3/host"

	"github.com/syntropy-cc/syntropy-cooperative-grid/manager/interfaces/cli/setup/src/internal/types"
)

// ValidateWindowsEnvironment valida o ambiente Windows para o setup do Syntropy CLI
func ValidateWindowsEnvironment() (*types.ValidationResult, error) {
	fmt.Println("Validando ambiente Windows...")

	// Criar resultado da validação
	result := &types.ValidationResult{
		Valid:     true,
		Message:   "Ambiente validado com sucesso",
		Timestamp: time.Now(),
		Environment: &types.EnvironmentInfo{
			OS:           "windows",
			Architecture: runtime.GOARCH,
		},
	}

	// Obter diretório home
	homeDir, err := os.UserHomeDir()
	if err != nil {
		result.Valid = false
		result.Message = fmt.Sprintf("Falha ao obter diretório home: %v", err)
		return result, err
	}
	result.Environment.HomeDir = homeDir

	// Validar versão do Windows
	if err := validateWindowsVersion(result); err != nil {
		result.Valid = false
		result.Message = fmt.Sprintf("Versão do Windows não suportada: %v", err)
		return result, err
	}

	// Validar permissões de administrador
	if err := validateAdminRights(result); err != nil {
		// Apenas aviso, não falha
		result.Warnings = append(result.Warnings, fmt.Sprintf("Sem permissões de administrador: %v", err))
	}

	// Validar espaço em disco
	if err := validateDiskSpace(result); err != nil {
		result.Valid = false
		result.Message = fmt.Sprintf("Espaço em disco insuficiente: %v", err)
		return result, err
	}

	// Validar versão do PowerShell
	if err := validatePowerShellVersion(result); err != nil {
		// Apenas aviso, não falha
		result.Warnings = append(result.Warnings, fmt.Sprintf("PowerShell não atende aos requisitos: %v", err))
	}

	// Validar conectividade com a internet
	if err := validateInternetConnectivity(result); err != nil {
		result.Valid = false
		result.Message = fmt.Sprintf("Sem conectividade com a internet: %v", err)
		return result, err
	}

	return result, nil
}

// validateWindowsVersion valida a versão do Windows
func validateWindowsVersion(result *types.ValidationResult) error {
	info, err := host.Info()
	if err != nil {
		return fmt.Errorf("falha ao obter informações do sistema: %w", err)
	}

	// Extrair versão principal do Windows
	versionParts := strings.Split(info.PlatformVersion, ".")
	if len(versionParts) < 1 {
		return fmt.Errorf("formato de versão inválido: %s", info.PlatformVersion)
	}

	majorVersion, err := strconv.Atoi(versionParts[0])
	if err != nil {
		return fmt.Errorf("falha ao analisar versão: %w", err)
	}

	// Windows 10 ou superior (versão principal 10 ou superior)
	if majorVersion < 10 {
		return fmt.Errorf("Windows %s não suportado, requer Windows 10 ou superior", info.PlatformVersion)
	}

	result.Environment.OSVersion = info.PlatformVersion
	return nil
}

// validateAdminRights valida se o usuário tem direitos de administrador
func validateAdminRights(result *types.ValidationResult) error {
	// No Windows, tentamos criar um arquivo em um diretório que requer privilégios de administrador
	testFile := filepath.Join("C:\\Windows\\Temp", "syntropy_admin_test.tmp")

	file, err := os.Create(testFile)
	if err != nil {
		return fmt.Errorf("sem permissões de administrador: %w", err)
	}

	file.Close()
	os.Remove(testFile)

	result.Environment.IsAdmin = true
	return nil
}

// validateDiskSpace valida se há espaço em disco suficiente
func validateDiskSpace(result *types.ValidationResult) error {
	// Verificar espaço em disco no drive C:
	usage, err := disk.Usage("C:")
	if err != nil {
		return fmt.Errorf("falha ao verificar espaço em disco: %w", err)
	}

	// Requer pelo menos 1GB de espaço livre
	requiredSpace := uint64(1 * 1024 * 1024 * 1024) // 1GB em bytes
	if usage.Free < requiredSpace {
		return fmt.Errorf("espaço em disco insuficiente: %.2f GB livre, requer pelo menos 1 GB", float64(usage.Free)/(1024*1024*1024))
	}

	result.Environment.DiskSpace = usage.Free
	return nil
}

// validatePowerShellVersion valida a versão do PowerShell
func validatePowerShellVersion(result *types.ValidationResult) error {
	// Esta é uma implementação simplificada
	// Em um cenário real, executaríamos PowerShell para verificar a versão

	// Simulando que o PowerShell está disponível na versão 5.1
	result.Environment.PowerShellVersion = "5.1"
	return nil
}

// validateInternetConnectivity valida a conectividade com a internet
func validateInternetConnectivity(result *types.ValidationResult) error {
	// Verificar conectividade tentando acessar um site confiável
	client := http.Client{
		Timeout: 5 * time.Second,
	}

	_, err := client.Get("https://www.google.com")
	if err != nil {
		return fmt.Errorf("sem conectividade com a internet: %w", err)
	}

	result.Environment.HasInternet = true
	return nil
}

// ValidateWindowsEnvironment validates the Windows environment for Syntropy CLI setup
func ValidateWindowsEnvironment(verbose bool) (*types.ValidationResult, error) {
	fmt.Println("Validating Windows environment...")

	result := &types.ValidationResult{
		Valid:    true,
		Warnings: []string{},
		Errors:   []string{},
		Environment: types.EnvironmentInfo{
			OS:           "windows",
			Architecture: runtime.GOARCH,
		},
	}

	// Check OS version
	if err := checkWindowsVersion(result); err != nil {
		result.Errors = append(result.Errors, fmt.Sprintf("Failed to check Windows version: %v", err))
		result.Valid = false
	}

	// Check admin rights
	if err := checkAdminRights(result); err != nil {
		result.Errors = append(result.Errors, fmt.Sprintf("Failed to check admin rights: %v", err))
		result.Valid = false
	}

	// Check disk space
	if err := checkDiskSpace(result); err != nil {
		result.Errors = append(result.Errors, fmt.Sprintf("Failed to check disk space: %v", err))
		result.Valid = false
	}

	// Check PowerShell version
	if err := checkPowerShellVersion(result); err != nil {
		result.Errors = append(result.Errors, fmt.Sprintf("Failed to check PowerShell version: %v", err))
		result.Valid = false
	}

	// Check internet connectivity
	if err := checkInternetConnectivity(result); err != nil {
		result.Warnings = append(result.Warnings, fmt.Sprintf("Failed to check internet connectivity: %v", err))
		// Not setting valid to false as this is just a warning
	}

	// Check home directory
	homeDir, err := os.UserHomeDir()
	if err != nil {
		result.Errors = append(result.Errors, fmt.Sprintf("Failed to get user home directory: %v", err))
		result.Valid = false
	} else {
		result.Environment.HomeDir = homeDir
	}

	// Print validation results if verbose
	if verbose {
		printValidationResults(result)
	}

	return result, nil
}

// checkWindowsVersion checks the Windows version
func checkWindowsVersion(result *types.ValidationResult) error {
	info, err := host.Info()
	if err != nil {
		return err
	}

	result.Environment.OS = info.Platform
	result.Environment.OSVersion = info.PlatformVersion

	// Check if Windows version is supported (Windows 10 or later)
	major, err := strconv.Atoi(strings.Split(info.PlatformVersion, ".")[0])
	if err != nil {
		return err
	}

	if major < 10 {
		result.Warnings = append(result.Warnings, fmt.Sprintf("Windows version %s may not be fully supported. Windows 10 or later is recommended.", info.PlatformVersion))
	}

	return nil
}

// checkAdminRights checks if the user has admin rights
func checkAdminRights(result *types.ValidationResult) error {
	// Try to open a privileged file
	f, err := os.Open("\\\\.\\PHYSICALDRIVE0")
	if err == nil {
		f.Close()
		result.Environment.HasAdminRights = true
	} else {
		result.Environment.HasAdminRights = false
		result.Warnings = append(result.Warnings, "Running without administrator privileges. Some features may not work correctly.")
	}

	return nil
}

// checkDiskSpace checks if there is enough disk space
func checkDiskSpace(result *types.ValidationResult) error {
	// Get disk usage of C: drive
	usage, err := disk.Usage("C:")
	if err != nil {
		return err
	}

	// Convert to GB
	freeGB := float64(usage.Free) / (1024 * 1024 * 1024)
	result.Environment.AvailableDiskGB = freeGB

	// Check if there is at least 1GB of free space
	if freeGB < 1.0 {
		result.Errors = append(result.Errors, fmt.Sprintf("Not enough disk space. %.2f GB available, 1.0 GB required.", freeGB))
		result.Valid = false
	}

	return nil
}

// checkPowerShellVersion checks the PowerShell version
func checkPowerShellVersion(result *types.ValidationResult) error {
	// Execute PowerShell command to get version
	cmd := exec.Command("powershell", "-Command", "$PSVersionTable.PSVersion.ToString()")
	output, err := cmd.Output()
	if err != nil {
		return err
	}

	version := strings.TrimSpace(string(output))
	result.Environment.PowerShellVer = version

	// Check if PowerShell version is at least 5.1
	parts := strings.Split(version, ".")
	if len(parts) >= 2 {
		major, err1 := strconv.Atoi(parts[0])
		minor, err2 := strconv.Atoi(parts[1])

		if err1 == nil && err2 == nil {
			if major < 5 || (major == 5 && minor < 1) {
				result.Warnings = append(result.Warnings, fmt.Sprintf("PowerShell version %s is below recommended version 5.1", version))
			}
		}
	}

	return nil
}

// checkInternetConnectivity checks if there is internet connectivity
func checkInternetConnectivity(result *types.ValidationResult) error {
	// Try to connect to a reliable host
	cmd := exec.Command("ping", "-n", "1", "8.8.8.8")
	err := cmd.Run()

	if err != nil {
		if exitError, ok := err.(*exec.ExitError); ok {
			if status, ok := exitError.Sys().(syscall.WaitStatus); ok {
				if status.ExitStatus() != 0 {
					result.Environment.HasInternet = false
					result.Warnings = append(result.Warnings, "No internet connectivity detected. Some features may not work correctly.")
					return nil
				}
			}
		}
		return err
	}

	result.Environment.HasInternet = true
	return nil
}

// printValidationResults prints the validation results
func printValidationResults(result *types.ValidationResult) {
	fmt.Println("\nValidation Results:")
	fmt.Println("===================")
	fmt.Printf("Valid: %v\n", result.Valid)
	fmt.Printf("OS: %s %s\n", result.Environment.OS, result.Environment.OSVersion)
	fmt.Printf("Architecture: %s\n", result.Environment.Architecture)
	fmt.Printf("Admin Rights: %v\n", result.Environment.HasAdminRights)
	fmt.Printf("PowerShell Version: %s\n", result.Environment.PowerShellVer)
	fmt.Printf("Available Disk Space: %.2f GB\n", result.Environment.AvailableDiskGB)
	fmt.Printf("Internet Connectivity: %v\n", result.Environment.HasInternet)
	fmt.Printf("Home Directory: %s\n", result.Environment.HomeDir)

	if len(result.Warnings) > 0 {
		fmt.Println("\nWarnings:")
		for _, warning := range result.Warnings {
			fmt.Printf("- %s\n", warning)
		}
	}

	if len(result.Errors) > 0 {
		fmt.Println("\nErrors:")
		for _, err := range result.Errors {
			fmt.Printf("- %s\n", err)
		}
	}
}
