//go:build windows
// +build windows

package setup

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/syntropy-cc/syntropy-cooperative-grid/manager/interfaces/cli/setup/src/internal/types"
)

// setupWindows implementa a configuração específica para Windows
func setupWindows(options types.SetupOptions) (*types.SetupResult, error) {
	fmt.Println("Iniciando setup para Windows...")

	// Create result structure
	result := &types.SetupResult{
		Success:     false,
		StartTime:   time.Now(),
		Environment: "windows",
		Options:     options,
	}

	// Step 1: Validate Windows environment
	fmt.Println("Etapa 1/3: Validando ambiente Windows...")
	validationResult, err := ValidateWindowsEnvironment()
	if err != nil {
		result.Error = fmt.Errorf("falha na validação do ambiente: %w", err)
		result.EndTime = time.Now()
		return result, result.Error
	}

	if !validationResult.Valid && !options.Force {
		result.Error = fmt.Errorf("ambiente inválido: %s", validationResult.Message)
		result.EndTime = time.Now()
		return result, result.Error
	}

	// Step 2: Configure Windows environment
	fmt.Println("Etapa 2/3: Configurando ambiente Windows...")
	if err := ConfigureWindowsEnvironment(validationResult, options); err != nil {
		result.Error = fmt.Errorf("falha na configuração do ambiente: %w", err)
		result.EndTime = time.Now()
		return result, result.Error
	}

	// Step 3: Install Windows services and dependencies
	fmt.Println("Etapa 3/3: Instalando serviços e dependências...")
	if err := installWindowsServices(validationResult, options); err != nil {
		result.Error = fmt.Errorf("falha na instalação dos serviços: %w", err)
		result.EndTime = time.Now()
		return result, result.Error
	}

	// Setup completed successfully
	result.Success = true
	result.EndTime = time.Now()
	result.ConfigPath = filepath.Join(validationResult.Environment.HomeDir, ".syntropy", "config", "manager.yaml")
	if options.ConfigPath != "" {
		result.ConfigPath = options.ConfigPath
	}

	fmt.Println("Setup concluído com sucesso!")
	fmt.Printf("Arquivo de configuração: %s\n", result.ConfigPath)
	fmt.Printf("Duração do setup: %s\n", result.EndTime.Sub(result.StartTime))

	return result, nil
}

// installWindowsServices installs required Windows services and dependencies
func installWindowsServices(validationResult *types.ValidationResult, options types.SetupOptions) error {
	fmt.Println("Installing Windows services and dependencies...")

	// Create service directory
	serviceDir := filepath.Join(validationResult.Environment.HomeDir, ".syntropy", "services")
	if err := os.MkdirAll(serviceDir, 0755); err != nil {
		return fmt.Errorf("failed to create service directory: %w", err)
	}

	// Install Windows service for Syntropy Manager
	if options.InstallService {
		fmt.Println("Installing Syntropy Manager Windows service...")
		// This would typically use Windows Service Control Manager APIs
		// For this implementation, we'll just create a placeholder script
		serviceScript := filepath.Join(serviceDir, "install_service.ps1")
		serviceContent := `
# Syntropy Manager Service Installation Script
$ServiceName = "SyntropyManager"
$DisplayName = "Syntropy Manager Service"
$BinPath = "$env:USERPROFILE\.syntropy\bin\syntropy-manager.exe"
$Description = "Manages Syntropy network connections and configurations"

# Check if service exists
$service = Get-Service -Name $ServiceName -ErrorAction SilentlyContinue
if ($service) {
    Write-Host "Service $ServiceName already exists. Stopping and removing..."
    Stop-Service -Name $ServiceName -Force
    sc.exe delete $ServiceName
}

# Create new service
New-Service -Name $ServiceName -BinaryPathName $BinPath -DisplayName $DisplayName -Description $Description -StartupType Automatic
Write-Host "Service $ServiceName installed successfully"
`
		if err := os.WriteFile(serviceScript, []byte(serviceContent), 0644); err != nil {
			return fmt.Errorf("failed to create service installation script: %w", err)
		}

		fmt.Println("Service installation script created at:", serviceScript)
		fmt.Println("To install the service, run the script as Administrator")
	}

	// Create startup script
	startupScript := filepath.Join(serviceDir, "startup.ps1")
	startupContent := `
# Syntropy Manager Startup Script
$SyntropyHome = "$env:USERPROFILE\.syntropy"
$ConfigPath = "$SyntropyHome\config\manager.yaml"
$LogPath = "$SyntropyHome\logs\manager.log"

# Start Syntropy Manager
Start-Process -FilePath "$SyntropyHome\bin\syntropy-manager.exe" -ArgumentList "--config", $ConfigPath, "--log", $LogPath
Write-Host "Syntropy Manager started"
`
	if err := os.WriteFile(startupScript, []byte(startupContent), 0644); err != nil {
		return fmt.Errorf("failed to create startup script: %w", err)
	}

	fmt.Println("Startup script created at:", startupScript)

	return nil
}

// resetWindows implements the Windows-specific reset process
func resetWindows(options types.SetupOptions) (*types.SetupResult, error) {
	fmt.Println("Resetting Syntropy CLI setup for Windows...")

	// Create result structure
	result := &types.SetupResult{
		Success:     false,
		StartTime:   time.Now(),
		Environment: "windows",
		Options:     options,
	}

	// Get Syntropy directory
	homeDir, err := os.UserHomeDir()
	if err != nil {
		result.Error = fmt.Errorf("failed to get user home directory: %w", err)
		result.EndTime = time.Now()
		return result, result.Error
	}

	syntropyDir := filepath.Join(homeDir, ".syntropy")

	// Check if Syntropy directory exists
	if _, err := os.Stat(syntropyDir); os.IsNotExist(err) {
		result.Error = fmt.Errorf("syntropy directory does not exist: %s", syntropyDir)
		result.EndTime = time.Now()
		return result, result.Error
	}

	// Stop services if running
	fmt.Println("Stopping Syntropy services...")
	// This would typically use Windows Service Control Manager APIs
	// For this implementation, we'll just create a placeholder script
	resetScript := filepath.Join(syntropyDir, "services", "reset.ps1")
	resetContent := `
# Syntropy Manager Service Reset Script
$ServiceName = "SyntropyManager"

# Check if service exists
$service = Get-Service -Name $ServiceName -ErrorAction SilentlyContinue
if ($service) {
    Write-Host "Stopping service $ServiceName..."
    Stop-Service -Name $ServiceName -Force
    Write-Host "Removing service $ServiceName..."
    sc.exe delete $ServiceName
}

Write-Host "Syntropy services stopped and removed"
`
	if err := os.WriteFile(resetScript, []byte(resetContent), 0644); err != nil {
		fmt.Println("Warning: Failed to create reset script:", err)
	} else {
		fmt.Println("Reset script created at:", resetScript)
		fmt.Println("To remove services, run the script as Administrator")
	}

	// Remove Syntropy directory if requested
	if options.Force {
		fmt.Println("Removing Syntropy directory...")
		if err := os.RemoveAll(syntropyDir); err != nil {
			result.Error = fmt.Errorf("failed to remove syntropy directory: %w", err)
			result.EndTime = time.Now()
			return result, result.Error
		}
		fmt.Println("Syntropy directory removed:", syntropyDir)
	} else {
		fmt.Println("Keeping Syntropy directory. Use --force to remove it.")
	}

	// Reset completed successfully
	result.Success = true
	result.EndTime = time.Now()

	fmt.Println("Syntropy CLI reset completed successfully!")
	fmt.Printf("Reset duration: %s\n", result.EndTime.Sub(result.StartTime))

	return result, nil
}

// statusWindows implementa a verificação de status específica para Windows
func statusWindows(options types.SetupOptions) (*types.SetupResult, error) {
	fmt.Println("Verificando status para Windows...")

	// Create result structure
	result := &types.SetupResult{
		Success:     false,
		StartTime:   time.Now(),
		Environment: "windows",
		Options:     options,
	}

	// Get Syntropy directory
	homeDir, err := os.UserHomeDir()
	if err != nil {
		result.Error = fmt.Errorf("falha ao obter diretório do usuário: %w", err)
		result.EndTime = time.Now()
		return result, result.Error
	}

	syntropyDir := filepath.Join(homeDir, ".syntropy")
	if options.SyntropyDir != "" {
		syntropyDir = options.SyntropyDir
	}

	// Check if Syntropy directory exists
	if _, err := os.Stat(syntropyDir); os.IsNotExist(err) {
		result.Error = fmt.Errorf("Syntropy CLI não está configurado (diretório não encontrado): %s", syntropyDir)
		result.EndTime = time.Now()
		return result, result.Error
	}

	// Check configuration file
	configPath := filepath.Join(syntropyDir, "config", "manager.yaml")
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		result.Error = fmt.Errorf("configuração do Syntropy CLI não encontrada: %s", configPath)
		result.EndTime = time.Now()
		return result, result.Error
	}

	// Check owner key
	keyPath := filepath.Join(syntropyDir, "keys", "owner.key")
	if _, err := os.Stat(keyPath); os.IsNotExist(err) {
		result.Error = fmt.Errorf("chave do proprietário do Syntropy CLI não encontrada: %s", keyPath)
		result.EndTime = time.Now()
		return result, result.Error
	}

	// Check service status
	// This would typically use Windows Service Control Manager APIs
	// For this implementation, we'll just report that it needs to be checked manually
	fmt.Println("Nota: Para verificar o status do serviço, execute 'Get-Service SyntropyManager' no PowerShell como Administrador")

	// Status check completed successfully
	result.Success = true
	result.EndTime = time.Now()
	result.ConfigPath = configPath

	fmt.Println("Syntropy CLI está configurado corretamente!")
	fmt.Printf("Arquivo de configuração: %s\n", result.ConfigPath)
	fmt.Printf("Chave do proprietário: %s\n", keyPath)

	return result, nil
}
