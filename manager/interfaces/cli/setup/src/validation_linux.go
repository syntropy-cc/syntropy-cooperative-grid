//go:build linux
// +build linux

package setup

import (
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"runtime"
	"strings"
	"time"

	"github.com/shirou/gopsutil/v3/disk"
	"github.com/shirou/gopsutil/v3/host"

	"github.com/syntropy-cc/syntropy-cooperative-grid/manager/interfaces/cli/setup/internal/types"
)

// ValidateLinuxEnvironment valida o ambiente Linux para o setup do Syntropy CLI
func ValidateLinuxEnvironment(verbose bool) (*types.ValidationResult, error) {
	fmt.Println("Validando ambiente Linux...")

	// Criar resultado da validação
	result := &types.ValidationResult{
		Valid:    true,
		Warnings: []string{},
		Errors:   []string{},
		Environment: types.EnvironmentInfo{
			OS:           "linux",
			Architecture: runtime.GOARCH,
		},
	}

	// Verificar versão do Linux
	if err := checkLinuxVersion(result); err != nil {
		result.Errors = append(result.Errors, fmt.Sprintf("Falha ao verificar versão do Linux: %v", err))
		result.Valid = false
	}

	// Verificar permissões de administrador
	if err := checkRootPermissions(result); err != nil {
		result.Warnings = append(result.Warnings, fmt.Sprintf("Sem permissões de root: %v", err))
	}

	// Verificar espaço em disco
	if err := checkDiskSpace(result); err != nil {
		result.Errors = append(result.Errors, fmt.Sprintf("Falha ao verificar espaço em disco: %v", err))
		result.Valid = false
	}

	// Verificar dependências do sistema
	if err := checkSystemDependencies(result); err != nil {
		result.Warnings = append(result.Warnings, fmt.Sprintf("Dependências do sistema não satisfeitas: %v", err))
	}

	// Verificar conectividade com a internet
	if err := checkInternetConnectivity(result); err != nil {
		result.Warnings = append(result.Warnings, fmt.Sprintf("Falha ao verificar conectividade com a internet: %v", err))
	}

	// Obter diretório home
	homeDir, err := os.UserHomeDir()
	if err != nil {
		result.Errors = append(result.Errors, fmt.Sprintf("Falha ao obter diretório home: %v", err))
		result.Valid = false
	} else {
		result.Environment.HomeDir = homeDir
	}

	// Imprimir resultados da validação se verbose
	if verbose {
		printValidationResults(result)
	}

	return result, nil
}

// checkLinuxVersion verifica a versão do Linux
func checkLinuxVersion(result *types.ValidationResult) error {
	info, err := host.Info()
	if err != nil {
		return err
	}

	result.Environment.OS = info.Platform
	result.Environment.OSVersion = info.PlatformVersion

	// Verificar se a distribuição é suportada
	supportedDistros := []string{"ubuntu", "debian", "centos", "fedora", "rhel"}
	distroSupported := false

	for _, distro := range supportedDistros {
		if strings.ToLower(info.Platform) == distro {
			distroSupported = true
			break
		}
	}

	if !distroSupported {
		result.Warnings = append(result.Warnings, fmt.Sprintf("Distribuição Linux %s pode não ser totalmente suportada. Ubuntu, Debian, CentOS, Fedora ou RHEL são recomendados.", info.Platform))
	}

	return nil
}

// checkRootPermissions verifica se o usuário tem permissões de root
func checkRootPermissions(result *types.ValidationResult) error {
	// Verificar se o usuário é root (UID 0)
	if os.Geteuid() == 0 {
		result.Environment.HasAdminRights = true
		return nil
	}

	// Tentar executar um comando que requer privilégios
	cmd := exec.Command("sudo", "-n", "true")
	err := cmd.Run()

	if err == nil {
		result.Environment.HasAdminRights = true
		return nil
	}

	result.Environment.HasAdminRights = false
	result.Warnings = append(result.Warnings, "Executando sem privilégios de administrador. Algumas funcionalidades podem não funcionar corretamente.")

	return fmt.Errorf("sem permissões de root")
}

// checkDiskSpace verifica se há espaço em disco suficiente
func checkDiskSpace(result *types.ValidationResult) error {
	// Obter uso do disco no diretório home
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return err
	}

	usage, err := disk.Usage(homeDir)
	if err != nil {
		return err
	}

	// Converter para GB
	freeGB := float64(usage.Free) / (1024 * 1024 * 1024)
	result.Environment.AvailableDiskGB = freeGB

	// Verificar se há pelo menos 1GB de espaço livre
	if freeGB < 1.0 {
		result.Errors = append(result.Errors, fmt.Sprintf("Espaço em disco insuficiente. %.2f GB disponíveis, 1.0 GB necessários.", freeGB))
		result.Valid = false
	}

	return nil
}

// checkSystemDependencies verifica as dependências do sistema
func checkSystemDependencies(result *types.ValidationResult) error {
	// Lista de dependências necessárias
	dependencies := []string{"curl", "wget", "git"}
	missingDeps := []string{}

	for _, dep := range dependencies {
		cmd := exec.Command("which", dep)
		if err := cmd.Run(); err != nil {
			missingDeps = append(missingDeps, dep)
		}
	}

	if len(missingDeps) > 0 {
		return fmt.Errorf("dependências ausentes: %s", strings.Join(missingDeps, ", "))
	}

	return nil
}

// checkInternetConnectivity verifica se há conectividade com a internet
func checkInternetConnectivity(result *types.ValidationResult) error {
	// Tentar conectar a um host confiável
	client := http.Client{
		Timeout: 5 * time.Second,
	}

	_, err := client.Get("https://www.google.com")
	if err != nil {
		result.Environment.HasInternet = false
		result.Warnings = append(result.Warnings, "Sem conectividade com a internet detectada. Algumas funcionalidades podem não funcionar corretamente.")
		return err
	}

	result.Environment.HasInternet = true
	return nil
}

// printValidationResults imprime os resultados da validação
func printValidationResults(result *types.ValidationResult) {
	fmt.Println("\nResultados da Validação:")
	fmt.Println("===================")
	fmt.Printf("Válido: %v\n", result.Valid)
	fmt.Printf("SO: %s %s\n", result.Environment.OS, result.Environment.OSVersion)
	fmt.Printf("Arquitetura: %s\n", result.Environment.Architecture)
	fmt.Printf("Permissões de Administrador: %v\n", result.Environment.HasAdminRights)
	fmt.Printf("Espaço em Disco Disponível: %.2f GB\n", result.Environment.AvailableDiskGB)
	fmt.Printf("Conectividade com a Internet: %v\n", result.Environment.HasInternet)
	fmt.Printf("Diretório Home: %s\n", result.Environment.HomeDir)

	if len(result.Warnings) > 0 {
		fmt.Println("\nAvisos:")
		for _, warning := range result.Warnings {
			fmt.Printf("- %s\n", warning)
		}
	}

	if len(result.Errors) > 0 {
		fmt.Println("\nErros:")
		for _, err := range result.Errors {
			fmt.Printf("- %s\n", err)
		}
	}
}
