//go:build windows
// +build windows

package setup

import (
	"crypto/ed25519"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"gopkg.in/yaml.v3"

	"github.com/syntropy-cc/syntropy-cooperative-grid/manager/interfaces/cli/setup/internal/types"
)

// ConfigureWindowsEnvironment configura o ambiente Windows para o Syntropy CLI
func ConfigureWindowsEnvironment(syntropyDir string, validationResult *types.ValidationResult) (*types.SetupConfig, error) {
	fmt.Println("Configurando ambiente Windows...")

	// Criar diretórios necessários
	dirs := []string{
		filepath.Join(syntropyDir, "config"),
		filepath.Join(syntropyDir, "logs"),
		filepath.Join(syntropyDir, "data"),
		filepath.Join(syntropyDir, "bin"),
	}

	for _, dir := range dirs {
		if err := os.MkdirAll(dir, 0755); err != nil {
			return nil, fmt.Errorf("falha ao criar diretório %s: %w", dir, err)
		}
	}

	// Gerar par de chaves Ed25519 para o proprietário
	pubKey, privKey, err := ed25519.GenerateKey(rand.Reader)
	if err != nil {
		return nil, fmt.Errorf("falha ao gerar par de chaves: %w", err)
	}

	// Codificar chaves em base64
	pubKeyBase64 := base64.StdEncoding.EncodeToString(pubKey)
	privKeyBase64 := base64.StdEncoding.EncodeToString(privKey)

	// Criar configuração
	config := &types.SetupConfig{
		Manager: types.ManagerConfig{
			Paths: map[string]string{
				"config":      filepath.Join(syntropyDir, "config"),
				"logs":        filepath.Join(syntropyDir, "logs"),
				"data":        filepath.Join(syntropyDir, "data"),
				"bin":         filepath.Join(syntropyDir, "bin"),
				"syntropy":    syntropyDir,
				"home":        validationResult.Environment.HomeDir,
				"manager_exe": filepath.Join(syntropyDir, "bin", "syntropy-manager.exe"),
			},
			DefaultFiles: map[string]string{
				"config": filepath.Join(syntropyDir, "config", "manager.yaml"),
				"log":    filepath.Join(syntropyDir, "logs", "manager.log"),
			},
			OwnerKey: types.OwnerKey{
				Public:  pubKeyBase64,
				Private: privKeyBase64,
				Type:    "ed25519",
			},
			Environment: types.Environment{
				OS:                validationResult.Environment.OS,
				OSVersion:         validationResult.Environment.OSVersion,
				Architecture:      validationResult.Environment.Architecture,
				PowerShellVersion: validationResult.Environment.PowerShellVersion,
				IsAdmin:           validationResult.Environment.IsAdmin,
				HasInternet:       validationResult.Environment.HasInternet,
				SetupTimestamp:    time.Now().Format(time.RFC3339),
			},
		},
	}

	// Salvar configuração em YAML
	configPath := filepath.Join(syntropyDir, "config", "manager.yaml")
	configData, err := yaml.Marshal(config)
	if err != nil {
		return nil, fmt.Errorf("falha ao serializar configuração: %w", err)
	}

	if err := os.WriteFile(configPath, configData, 0644); err != nil {
		return nil, fmt.Errorf("falha ao salvar configuração: %w", err)
	}

	fmt.Printf("Configuração salva em: %s\n", configPath)
	return config, nil
}

// ConfigureWindowsEnvironment configures the Windows environment for Syntropy CLI
func ConfigureWindowsEnvironment(validationResult *types.ValidationResult, options SetupOptions) error {
	fmt.Println("Configuring Windows environment...")

	// Create Syntropy directory structure
	if err := createSyntropyDirStructure(validationResult.Environment.HomeDir); err != nil {
		return fmt.Errorf("failed to create Syntropy directory structure: %w", err)
	}

	// Generate owner key
	ownerKeyPath, err := generateOwnerKey(validationResult.Environment.HomeDir)
	if err != nil {
		return fmt.Errorf("failed to generate owner key: %w", err)
	}

	// Generate manager configuration
	if err := generateManagerConfig(validationResult.Environment.HomeDir, ownerKeyPath, options); err != nil {
		return fmt.Errorf("failed to generate manager configuration: %w", err)
	}

	fmt.Println("Windows environment configured successfully!")
	return nil
}

// createSyntropyDirStructure creates the Syntropy directory structure
func createSyntropyDirStructure(homeDir string) error {
	syntropyDir := filepath.Join(homeDir, ".syntropy")
	
	// Create main directories
	dirs := []string{
		syntropyDir,
		filepath.Join(syntropyDir, "config"),
		filepath.Join(syntropyDir, "keys"),
		filepath.Join(syntropyDir, "nodes"),
		filepath.Join(syntropyDir, "logs"),
		filepath.Join(syntropyDir, "cache"),
		filepath.Join(syntropyDir, "backups"),
		filepath.Join(syntropyDir, "scripts"),
	}

	for _, dir := range dirs {
		if err := os.MkdirAll(dir, 0755); err != nil {
			return fmt.Errorf("failed to create directory %s: %w", dir, err)
		}
	}

	return nil
}

// generateOwnerKey generates an Ed25519 key pair for the owner
func generateOwnerKey(homeDir string) (string, error) {
	keysDir := filepath.Join(homeDir, ".syntropy", "keys")
	privateKeyPath := filepath.Join(keysDir, "owner.key")
	publicKeyPath := filepath.Join(keysDir, "owner.key.pub")

	// Check if keys already exist
	if _, err := os.Stat(privateKeyPath); err == nil {
		fmt.Println("Owner key already exists, skipping generation")
		return privateKeyPath, nil
	}

	// Generate Ed25519 key pair
	publicKey, privateKey, err := ed25519.GenerateKey(rand.Reader)
	if err != nil {
		return "", fmt.Errorf("failed to generate Ed25519 key pair: %w", err)
	}

	// Save private key
	privateKeyPEM := pem.EncodeToMemory(&pem.Block{
		Type:  "PRIVATE KEY",
		Bytes: privateKey,
	})
	if err := os.WriteFile(privateKeyPath, privateKeyPEM, 0600); err != nil {
		return "", fmt.Errorf("failed to write private key: %w", err)
	}

	// Save public key
	publicKeyPEM := pem.EncodeToMemory(&pem.Block{
		Type:  "PUBLIC KEY",
		Bytes: publicKey,
	})
	if err := os.WriteFile(publicKeyPath, publicKeyPEM, 0644); err != nil {
		return "", fmt.Errorf("failed to write public key: %w", err)
	}

	fmt.Println("Generated owner key pair")
	return privateKeyPath, nil
}

// generateManagerConfig generates the manager configuration
func generateManagerConfig(homeDir, ownerKeyPath string, options SetupOptions) error {
	syntropyDir := filepath.Join(homeDir, ".syntropy")
	configPath := filepath.Join(syntropyDir, "config", "manager.yaml")

	// Use custom config path if provided
	if options.ConfigPath != "" {
		configPath = options.ConfigPath
	}

	// Check if config already exists
	if _, err := os.Stat(configPath); err == nil && !options.Force {
		fmt.Println("Manager configuration already exists, skipping generation")
		return nil
	}

	// Create manager configuration
	config := types.SetupConfig{
		Manager: types.ManagerConfig{
			HomeDir:     syntropyDir,
			LogLevel:    "info",
			APIEndpoint: "https://api.syntropy.network",
			Directories: map[string]string{
				"config":  filepath.Join(syntropyDir, "config"),
				"keys":    filepath.Join(syntropyDir, "keys"),
				"nodes":   filepath.Join(syntropyDir, "nodes"),
				"logs":    filepath.Join(syntropyDir, "logs"),
				"cache":   filepath.Join(syntropyDir, "cache"),
				"backups": filepath.Join(syntropyDir, "backups"),
				"scripts": filepath.Join(syntropyDir, "scripts"),
			},
			DefaultPaths: map[string]string{
				"owner_key": ownerKeyPath,
				"config":    configPath,
			},
		},
		OwnerKey: types.OwnerKey{
			Type: "Ed25519",
			Path: ownerKeyPath,
		},
		Environment: types.Environment{
			OS:           "windows",
			Architecture: "amd64",
			HomeDir:      homeDir,
		},
	}

	// Marshal configuration to YAML
	yamlData, err := yaml.Marshal(config)
	if err != nil {
		return fmt.Errorf("failed to marshal configuration: %w", err)
	}

	// Write configuration to file
	if err := os.WriteFile(configPath, yamlData, 0644); err != nil {
		return fmt.Errorf("failed to write configuration: %w", err)
	}

	fmt.Printf("Generated manager configuration at %s\n", configPath)
	return nil
}