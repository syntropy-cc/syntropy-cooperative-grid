//go:build linux
// +build linux

package setup

import (
	"crypto/ed25519"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"

	"github.com/syntropy-cc/syntropy-cooperative-grid/manager/interfaces/cli/setup/src/internal/types"
)

// ConfigureLinuxEnvironment configures the Linux environment for Syntropy CLI
func ConfigureLinuxEnvironment(validationResult *types.ValidationResult, options types.SetupOptions) error {
	fmt.Println("Configurando ambiente Linux...")

	// Initialize security validator
	securityValidator := NewSecurityValidator()

	// Pré-declarações para evitar shadowing de err e chaves
	var (
		err     error
		pubKey  ed25519.PublicKey
		privKey ed25519.PrivateKey
	)

	// Determine Syntropy directory
	syntropyDir := filepath.Join(validationResult.Environment.HomeDir, ".syntropy")
	if options.HomeDir != "" {
		syntropyDir = filepath.Join(options.HomeDir, ".syntropy")
	}

	// Create necessary directories
	dirs := []string{
		filepath.Join(syntropyDir, "config"),
		filepath.Join(syntropyDir, "logs"),
		filepath.Join(syntropyDir, "data"),
		filepath.Join(syntropyDir, "bin"),
		filepath.Join(syntropyDir, "services"),
	}

	for _, dir := range dirs {
		err = os.MkdirAll(dir, 0o755)
		if err != nil {
			return fmt.Errorf("falha ao criar diretório %s: %w", dir, err)
		}
	}

	// Generate Ed25519 key pair for owner
	pubKey, privKey, err = ed25519.GenerateKey(rand.Reader)
	if err != nil {
		return fmt.Errorf("falha ao gerar par de chaves: %w", err)
	}

	// Encode keys in base64
	pubKeyBase64 := base64.StdEncoding.EncodeToString(pubKey)
	// Store private key if needed in the future
	_ = base64.StdEncoding.EncodeToString(privKey)

	// Create configuration
	config := &types.SetupConfig{
		Manager: types.ManagerConfig{
			HomeDir:     syntropyDir,
			LogLevel:    "info",
			APIEndpoint: "https://api.syntropy.network",
			Directories: map[string]string{
				"config":   filepath.Join(syntropyDir, "config"),
				"logs":     filepath.Join(syntropyDir, "logs"),
				"data":     filepath.Join(syntropyDir, "data"),
				"bin":      filepath.Join(syntropyDir, "bin"),
				"syntropy": syntropyDir,
				"home":     validationResult.Environment.HomeDir,
			},
			DefaultPaths: map[string]string{
				"config": filepath.Join(syntropyDir, "config", "manager.yaml"),
				"log":    filepath.Join(syntropyDir, "logs", "manager.log"),
			},
		},
		OwnerKey: types.OwnerKey{
			Type:      "ed25519",
			Path:      filepath.Join(syntropyDir, "config", "owner.key"),
			PublicKey: pubKeyBase64,
		},
		Environment: types.Environment{
			OS:           validationResult.Environment.OS,
			Architecture: validationResult.Environment.Architecture,
			HomeDir:      validationResult.Environment.HomeDir,
		},
	}

	// Save configuration as YAML
	configPath := filepath.Join(syntropyDir, "config", "manager.yaml")
	if options.ConfigPath != "" {
		configPath = options.ConfigPath
		// Validate config path for security
		if err := securityValidator.ValidatePath(configPath, syntropyDir); err != nil {
			return fmt.Errorf("invalid config path: %v", err)
		}
		// Ensure directory exists
		configDir := filepath.Dir(configPath)
		err = os.MkdirAll(configDir, 0o755)
		if err != nil {
			return fmt.Errorf("falha ao criar diretório de configuração: %w", err)
		}
	}

	configData, err := yaml.Marshal(config)
	if err != nil {
		return fmt.Errorf("falha ao serializar configuração: %w", err)
	}

	// Validate configuration content for security
	if err := securityValidator.ValidateConfiguration(string(configData)); err != nil {
		return fmt.Errorf("invalid configuration content: %v", err)
	}

	err = os.WriteFile(configPath, configData, 0o644)
	if err != nil {
		return fmt.Errorf("falha ao salvar configuração: %w", err)
	}

	// Set secure file permissions
	if err := securityValidator.SetSecureFilePermissions(configPath); err != nil {
		fmt.Printf("Warning: Could not set secure permissions for config file: %v\n", err)
	}

	fmt.Printf("Configuração salva em: %s\n", configPath)

	// Create symbolic links for binaries if needed
	if validationResult.Environment.HasAdminRights {
		// Create symlink in /usr/local/bin if we have admin rights
		binPath := filepath.Join(syntropyDir, "bin", "syntropy-manager")
		symlinkPath := "/usr/local/bin/syntropy-manager"

		// Remove existing symlink if it exists
		_ = os.Remove(symlinkPath)

		// Create new symlink
		err = os.Symlink(binPath, symlinkPath)
		if err != nil {
			fmt.Printf("Aviso: Não foi possível criar link simbólico em %s: %v\n", symlinkPath, err)
		} else {
			fmt.Printf("Link simbólico criado em: %s\n", symlinkPath)
		}
	}

	return nil
}
