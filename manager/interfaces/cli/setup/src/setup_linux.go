//go:build linux
// +build linux

package setup

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"time"

	"github.com/syntropy-cc/syntropy-cooperative-grid/manager/interfaces/cli/setup/src/internal/types"
)

// setupLinuxImpl implements Linux-specific setup
func setupLinuxImpl(options types.SetupOptions) (*types.SetupResult, error) {
	fmt.Println("Iniciando setup para Linux...")

	// Initialize security validator
	securityValidator := NewSecurityValidator()

	// Validate paths for security
	if err := securityValidator.ValidatePath(options.ConfigPath, ""); err != nil {
		return nil, fmt.Errorf("invalid config path: %v", err)
	}
	if err := securityValidator.ValidatePath(options.HomeDir, ""); err != nil {
		return nil, fmt.Errorf("invalid home directory: %v", err)
	}

	// Check for privilege escalation prevention
	if options.InstallService && !securityValidator.CheckAdminRights() {
		return nil, fmt.Errorf("service installation requires administrative privileges")
	}

	// Create result structure
	result := &types.SetupResult{
		Success:     false,
		StartTime:   time.Now(),
		Environment: "linux",
		Options:     options,
	}

	// Step 1: Validate Linux environment
	fmt.Println("Etapa 1/3: Validando ambiente Linux...")
	validationResult, err := ValidateLinuxEnvironment(true)
	if err != nil {
		result.Error = fmt.Errorf("falha na validação do ambiente: %w", err)
		result.EndTime = time.Now()
		return result, result.Error
	}

	if !validationResult.Valid && !options.Force {
		result.Error = fmt.Errorf("ambiente inválido, use --force para ignorar validações")
		result.EndTime = time.Now()
		return result, result.Error
	}

	// Step 2: Configure Linux environment
	fmt.Println("Etapa 2/3: Configurando ambiente Linux...")
	if err := ConfigureLinuxEnvironment(validationResult, options); err != nil {
		result.Error = fmt.Errorf("falha na configuração do ambiente: %w", err)
		result.EndTime = time.Now()
		return result, result.Error
	}

	// Step 3: Install Linux services and dependencies
	fmt.Println("Etapa 3/3: Instalando serviços e dependências...")
	if err := installLinuxServices(validationResult, options); err != nil {
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

// installLinuxServices installs required Linux services and dependencies
func installLinuxServices(validationResult *types.ValidationResult, options types.SetupOptions) error {
	fmt.Println("Instalando serviços e dependências Linux...")

	// Create service directory
	serviceDir := filepath.Join(validationResult.Environment.HomeDir, ".syntropy", "services")
	if err := os.MkdirAll(serviceDir, 0755); err != nil {
		return fmt.Errorf("falha ao criar diretório de serviços: %w", err)
	}

	// Create bin directory
	binDir := filepath.Join(validationResult.Environment.HomeDir, ".syntropy", "bin")
	if err := os.MkdirAll(binDir, 0755); err != nil {
		return fmt.Errorf("falha ao criar diretório bin: %w", err)
	}

	// Install systemd service for Syntropy Manager if requested
	if options.InstallService {
		if err := installSystemdService(validationResult, serviceDir, binDir); err != nil {
			return fmt.Errorf("falha ao instalar serviço systemd: %w", err)
		}
	}

	return nil
}

// installSystemdService installs and configures the systemd service
func installSystemdService(validationResult *types.ValidationResult, serviceDir, binDir string) error {
	fmt.Println("Instalando serviço systemd para Syntropy Manager...")

	// Check if systemd is available
	if _, err := exec.LookPath("systemctl"); err != nil {
		return fmt.Errorf("systemd não encontrado, não é possível instalar o serviço: %w", err)
	}

	// Create systemd service file
	serviceFilePath := filepath.Join(serviceDir, "syntropy-manager.service")
	serviceContent := fmt.Sprintf(`[Unit]
Description=Syntropy Manager Service
After=network.target

[Service]
Type=simple
User=%s
ExecStart=%s/syntropy-manager
Restart=on-failure
RestartSec=5
StandardOutput=journal
StandardError=journal

[Install]
WantedBy=multi-user.target
`, os.Getenv("USER"), binDir)

	if err := os.WriteFile(serviceFilePath, []byte(serviceContent), 0644); err != nil {
		return fmt.Errorf("falha ao criar arquivo de serviço: %w", err)
	}

	// Create installation script
	installScriptPath := filepath.Join(serviceDir, "install_service.sh")
	installScriptContent := fmt.Sprintf(`#!/bin/bash
# Script de instalação do serviço Syntropy Manager

# Verificar permissões de root
if [ "$EUID" -ne 0 ]; then
  echo "Este script precisa ser executado como root"
  exit 1
fi

# Copiar arquivo de serviço para o diretório systemd
cp %s /etc/systemd/system/syntropy-manager.service

# Recarregar configuração do systemd
systemctl daemon-reload

# Habilitar serviço para iniciar com o sistema
systemctl enable syntropy-manager.service

# Iniciar serviço
systemctl start syntropy-manager.service

echo "Serviço Syntropy Manager instalado e iniciado com sucesso"
`, serviceFilePath)

	if err := os.WriteFile(installScriptPath, []byte(installScriptContent), 0755); err != nil {
		return fmt.Errorf("falha ao criar script de instalação: %w", err)
	}

	fmt.Println("Arquivo de serviço systemd criado em:", serviceFilePath)
	fmt.Println("Para instalar o serviço, execute como root:")
	fmt.Printf("  sudo %s\n", installScriptPath)

	return nil
}

// statusLinux verifica o status da instalação no Linux
func statusLinux(options types.SetupOptions) (*types.SetupResult, error) {
	fmt.Println("Verificando status do Syntropy CLI no Linux...")

	result := &types.SetupResult{
		Success:     false,
		StartTime:   time.Now(),
		Environment: "linux",
		Options:     options,
	}

	// Verificar diretório .syntropy
	homeDir, err := os.UserHomeDir()
	if err != nil {
		result.Error = fmt.Errorf("falha ao obter diretório home: %w", err)
		result.EndTime = time.Now()
		return result, result.Error
	}

	syntropyDir := filepath.Join(homeDir, ".syntropy")
	if _, err := os.Stat(syntropyDir); os.IsNotExist(err) {
		result.Error = fmt.Errorf("diretório Syntropy não encontrado, execute o setup primeiro")
		result.EndTime = time.Now()
		return result, result.Error
	}

	// Verificar arquivo de configuração
	configPath := filepath.Join(syntropyDir, "config", "manager.yaml")
	if options.ConfigPath != "" {
		configPath = options.ConfigPath
	}

	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		result.Error = fmt.Errorf("arquivo de configuração não encontrado: %s", configPath)
		result.EndTime = time.Now()
		return result, result.Error
	}

	// Verificar serviço systemd se aplicável
	if options.InstallService {
		cmd := exec.Command("systemctl", "is-active", "syntropy-manager.service")
		output, _ := cmd.Output()

		if string(output) != "active\n" {
			result.Error = fmt.Errorf("serviço Syntropy Manager não está ativo")
			result.EndTime = time.Now()
			return result, result.Error
		}
	}

	// Status verificado com sucesso
	result.Success = true
	result.EndTime = time.Now()
	result.ConfigPath = configPath

	fmt.Println("Status do Syntropy CLI:")
	fmt.Printf("Diretório: %s\n", syntropyDir)
	fmt.Printf("Configuração: %s\n", configPath)

	return result, nil
}

// resetLinux redefine a configuração do Syntropy CLI no Linux
func resetLinux(options types.SetupOptions) (*types.SetupResult, error) {
	fmt.Println("Redefinindo configuração do Syntropy CLI no Linux...")

	result := &types.SetupResult{
		Success:     false,
		StartTime:   time.Now(),
		Environment: "linux",
		Options:     options,
	}

	// Obter diretório home
	homeDir, err := os.UserHomeDir()
	if err != nil {
		result.Error = fmt.Errorf("falha ao obter diretório home: %w", err)
		result.EndTime = time.Now()
		return result, result.Error
	}

	// Parar serviço se estiver instalado
	if options.InstallService {
		stopCmd := exec.Command("systemctl", "stop", "syntropy-manager.service")
		_ = stopCmd.Run() // Ignorar erro se o serviço não existir

		disableCmd := exec.Command("systemctl", "disable", "syntropy-manager.service")
		_ = disableCmd.Run()

		// Remover arquivo de serviço
		serviceFile := "/etc/systemd/system/syntropy-manager.service"
		if _, err := os.Stat(serviceFile); err == nil {
			rmCmd := exec.Command("sudo", "rm", serviceFile)
			_ = rmCmd.Run()

			reloadCmd := exec.Command("systemctl", "daemon-reload")
			_ = reloadCmd.Run()
		}
	}

	// Remover diretório .syntropy
	syntropyDir := filepath.Join(homeDir, ".syntropy")
	if _, err := os.Stat(syntropyDir); err == nil {
		if err := os.RemoveAll(syntropyDir); err != nil {
			result.Error = fmt.Errorf("falha ao remover diretório Syntropy: %w", err)
			result.EndTime = time.Now()
			return result, result.Error
		}
	}

	// Reset concluído com sucesso
	result.Success = true
	result.EndTime = time.Now()

	fmt.Println("Configuração do Syntropy CLI redefinida com sucesso")

	return result, nil
}
