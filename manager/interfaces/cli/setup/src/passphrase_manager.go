package setup

import (
	"fmt"
	"os"
	"regexp"
	"syscall"

	"setup-component/src/internal/types"

	"github.com/zalando/go-keyring"
	"golang.org/x/term"
)

const (
	KeyringServicePassphrase = "syntropy-owner-passphrase"
	KeyringUserPassphrase    = "owner-key-passphrase"
	MinPassphraseLength      = 12
)

// PassphraseManager manages user passphrases for key encryption
type PassphraseManager struct {
	logger types.SetupLogger
}

// NewPassphraseManager creates a new passphrase manager
func NewPassphraseManager(logger types.SetupLogger) *PassphraseManager {
	return &PassphraseManager{
		logger: logger,
	}
}

// PromptForPassphrase prompts the user for a passphrase with validation
func (pm *PassphraseManager) PromptForPassphrase() (string, error) {
	pm.logger.LogStep("passphrase_prompt_start", nil)

	fmt.Println("🔐 Configuração de Segurança")
	fmt.Println()
	fmt.Println("Para proteger suas chaves privadas, você precisa definir uma SENHA FORTE.")
	fmt.Println()
	fmt.Println("📋 Requisitos da senha:")
	fmt.Println("   • Mínimo 12 caracteres")
	fmt.Println("   • Pelo menos uma letra maiúscula")
	fmt.Println("   • Pelo menos uma letra minúscula")
	fmt.Println("   • Pelo menos um número")
	fmt.Println("   • Pelo menos um caractere especial (!@#$%^&*)")
	fmt.Println()
	fmt.Println("⚠️  BOAS PRÁTICAS DE SEGURANÇA:")
	fmt.Println("   • Use uma senha única (não reutilize de outros serviços)")
	fmt.Println("   • Considere usar um gerenciador de senhas")
	fmt.Println("   • NUNCA compartilhe sua senha")
	fmt.Println("   • Faça backup da senha em local seguro")
	fmt.Println("   • Troque a senha periodicamente (use: syntropy setup key migrate)")
	fmt.Println()

	for {
		fmt.Print("   Digite sua senha: ")
		passphrase, err := term.ReadPassword(int(syscall.Stdin))
		fmt.Println() // Nova linha após senha
		if err != nil {
			// Fallback para entrada normal se term.ReadPassword falhar
			fmt.Print("   (entrada sem echo não disponível, digite normalmente): ")
			var input string
			fmt.Scanln(&input)
			passphrase = []byte(input)
			err = nil // Reset error para continuar o loop
		}

		passphraseStr := string(passphrase)
		if err := pm.ValidatePassphraseStrength(passphraseStr); err != nil {
			fmt.Printf("❌ %s\n", err.Error())
			fmt.Println("   Tente novamente...")
			fmt.Println()
			continue
		}

		fmt.Print("   Confirme sua senha: ")
		confirmPassphrase, err := term.ReadPassword(int(syscall.Stdin))
		fmt.Println() // Nova linha após senha
		if err != nil {
			// Fallback para entrada normal se term.ReadPassword falhar
			fmt.Print("   (entrada sem echo não disponível, digite normalmente): ")
			var input string
			fmt.Scanln(&input)
			confirmPassphrase = []byte(input)
			err = nil // Reset error para continuar o loop
		}

		if passphraseStr != string(confirmPassphrase) {
			fmt.Println("❌ Senhas não correspondem. Tente novamente...")
			fmt.Println()
			continue
		}

		pm.logger.LogStep("passphrase_prompt_completed", map[string]interface{}{
			"length": len(passphraseStr),
		})

		return passphraseStr, nil
	}
}

// ValidatePassphraseStrength validates passphrase strength according to security requirements
func (pm *PassphraseManager) ValidatePassphraseStrength(passphrase string) error {
	if len(passphrase) < MinPassphraseLength {
		return fmt.Errorf("senha muito curta (mínimo %d caracteres)", MinPassphraseLength)
	}

	// Check for uppercase letter
	hasUpper, _ := regexp.MatchString(`[A-Z]`, passphrase)
	if !hasUpper {
		return fmt.Errorf("senha deve conter pelo menos uma letra maiúscula")
	}

	// Check for lowercase letter
	hasLower, _ := regexp.MatchString(`[a-z]`, passphrase)
	if !hasLower {
		return fmt.Errorf("senha deve conter pelo menos uma letra minúscula")
	}

	// Check for digit
	hasDigit, _ := regexp.MatchString(`[0-9]`, passphrase)
	if !hasDigit {
		return fmt.Errorf("senha deve conter pelo menos um número")
	}

	// Check for special character
	hasSpecial, _ := regexp.MatchString(`[!@#$%^&*()_+\-=\[\]{};':"\\|,.<>\/?]`, passphrase)
	if !hasSpecial {
		return fmt.Errorf("senha deve conter pelo menos um caractere especial (!@#$%^&*)")
	}

	return nil
}

// StorePassphraseInKeyring stores passphrase in system keyring
func (pm *PassphraseManager) StorePassphraseInKeyring(passphrase string) error {
	pm.logger.LogStep("passphrase_storage_start", map[string]interface{}{
		"keyring_available": isKeyringAvailable(),
	})

	if !isKeyringAvailable() {
		return pm.storePassphraseToFile(passphrase)
	}

	// Store in system keyring
	if err := keyring.Set(KeyringServicePassphrase, KeyringUserPassphrase, passphrase); err != nil {
		pm.logger.LogWarning("keyring_save_failed", map[string]interface{}{
			"error":            err.Error(),
			"fallback_to_file": true,
		})
		return pm.storePassphraseToFile(passphrase)
	}

	pm.logger.LogStep("passphrase_storage_completed", map[string]interface{}{
		"storage_method": "keyring",
	})

	return nil
}

// LoadPassphraseFromKeyring loads passphrase from system keyring
func (pm *PassphraseManager) LoadPassphraseFromKeyring() (string, error) {
	pm.logger.LogStep("passphrase_load_start", map[string]interface{}{
		"keyring_available": isKeyringAvailable(),
	})

	if !isKeyringAvailable() {
		return pm.loadPassphraseFromFile()
	}

	// Try to load from keyring first
	passphrase, err := keyring.Get(KeyringServicePassphrase, KeyringUserPassphrase)
	if err != nil {
		pm.logger.LogWarning("keyring_load_failed", map[string]interface{}{
			"error":            err.Error(),
			"fallback_to_file": true,
		})
		return pm.loadPassphraseFromFile()
	}

	pm.logger.LogStep("passphrase_load_completed", map[string]interface{}{
		"storage_method": "keyring",
	})

	return passphrase, nil
}

// PassphraseExists checks if passphrase exists in keyring or file
func (pm *PassphraseManager) PassphraseExists() bool {
	if isKeyringAvailable() {
		_, err := keyring.Get(KeyringServicePassphrase, KeyringUserPassphrase)
		if err == nil {
			return true
		}
	}

	// Check file fallback
	homeDir, _ := os.UserHomeDir()
	passphrasePath := fmt.Sprintf("%s/.syntropy/passphrase", homeDir)
	_, err := os.Stat(passphrasePath)
	return !os.IsNotExist(err)
}

// PromptForExistingPassphrase prompts for existing passphrase (for key operations)
func (pm *PassphraseManager) PromptForExistingPassphrase() (string, error) {
	fmt.Println("🔒 REQUISITO DE SEGURANÇA:")
	fmt.Println("   Para acessar a chave privada, você precisa fornecer a senha")
	fmt.Println("   definida durante o setup inicial.")
	fmt.Println()
	fmt.Print("   Digite sua senha: ")

	passphrase, err := term.ReadPassword(int(syscall.Stdin))
	fmt.Println() // Nova linha após senha
	if err != nil {
		// Fallback para entrada normal se term.ReadPassword falhar
		fmt.Print("   (entrada sem echo não disponível, digite normalmente): ")
		var input string
		fmt.Scanln(&input)
		passphrase = []byte(input)
	}

	return string(passphrase), nil
}

// ValidatePassphrase validates provided passphrase against stored one
func (pm *PassphraseManager) ValidatePassphrase(providedPassphrase string) error {
	storedPassphrase, err := pm.LoadPassphraseFromKeyring()
	if err != nil {
		// Se não há senha no keyring, permitir validação (para casos de migração)
		// A validação real será feita na descriptografia da chave
		return nil
	}

	if providedPassphrase != storedPassphrase {
		return fmt.Errorf("senha incorreta")
	}

	return nil
}

// DeletePassphrase removes passphrase from keyring
func (pm *PassphraseManager) DeletePassphrase() error {
	pm.logger.LogStep("passphrase_delete_start", map[string]interface{}{
		"keyring_available": isKeyringAvailable(),
	})

	if !isKeyringAvailable() {
		return pm.deletePassphraseFile()
	}

	// Try to remove from keyring
	if err := keyring.Delete(KeyringServicePassphrase, KeyringUserPassphrase); err != nil {
		pm.logger.LogWarning("keyring_delete_failed", map[string]interface{}{
			"error": err.Error(),
		})
		return pm.deletePassphraseFile()
	}

	pm.logger.LogStep("passphrase_delete_completed", map[string]interface{}{
		"storage_method": "keyring",
	})

	return nil
}

// Fallback methods for file storage

// storePassphraseToFile stores passphrase in file as fallback
func (pm *PassphraseManager) storePassphraseToFile(passphrase string) error {
	homeDir, _ := os.UserHomeDir()
	passphraseDir := fmt.Sprintf("%s/.syntropy", homeDir)
	passphrasePath := fmt.Sprintf("%s/passphrase", passphraseDir)

	// Create directory if it doesn't exist
	if err := os.MkdirAll(passphraseDir, 0700); err != nil {
		return fmt.Errorf("failed to create passphrase directory: %w", err)
	}

	// Write passphrase to file with restricted permissions
	if err := os.WriteFile(passphrasePath, []byte(passphrase), 0600); err != nil {
		return fmt.Errorf("failed to write passphrase file: %w", err)
	}

	pm.logger.LogStep("passphrase_save_fallback", map[string]interface{}{
		"storage_method":  "file",
		"passphrase_path": passphrasePath,
	})

	return nil
}

// loadPassphraseFromFile loads passphrase from file as fallback
func (pm *PassphraseManager) loadPassphraseFromFile() (string, error) {
	homeDir, _ := os.UserHomeDir()
	passphrasePath := fmt.Sprintf("%s/.syntropy/passphrase", homeDir)

	// Check if file exists
	if _, err := os.Stat(passphrasePath); os.IsNotExist(err) {
		return "", fmt.Errorf("passphrase not found in keyring or file")
	}

	// Read file
	data, err := os.ReadFile(passphrasePath)
	if err != nil {
		return "", fmt.Errorf("failed to read passphrase file: %w", err)
	}

	pm.logger.LogStep("passphrase_load_fallback", map[string]interface{}{
		"storage_method":  "file",
		"passphrase_path": passphrasePath,
	})

	return string(data), nil
}

// deletePassphraseFile removes passphrase file
func (pm *PassphraseManager) deletePassphraseFile() error {
	homeDir, _ := os.UserHomeDir()
	passphrasePath := fmt.Sprintf("%s/.syntropy/passphrase", homeDir)

	if err := os.Remove(passphrasePath); err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("failed to delete passphrase file: %w", err)
	}

	pm.logger.LogStep("passphrase_delete_fallback", map[string]interface{}{
		"storage_method":  "file",
		"passphrase_path": passphrasePath,
	})

	return nil
}

// Global functions for backward compatibility

// PromptForPassphrase is a global function that creates a manager and prompts for passphrase
func PromptForPassphrase() (string, error) {
	logger := NewSetupLogger()
	pm := NewPassphraseManager(logger)
	return pm.PromptForPassphrase()
}

// StorePassphraseInKeyring is a global function that creates a manager and stores passphrase
func StorePassphraseInKeyring(passphrase string) error {
	logger := NewSetupLogger()
	pm := NewPassphraseManager(logger)
	return pm.StorePassphraseInKeyring(passphrase)
}

// LoadPassphraseFromKeyring is a global function that creates a manager and loads passphrase
func LoadPassphraseFromKeyring() (string, error) {
	logger := NewSetupLogger()
	pm := NewPassphraseManager(logger)
	return pm.LoadPassphraseFromKeyring()
}

// PassphraseExists is a global function that creates a manager and checks if passphrase exists
func PassphraseExists() bool {
	logger := NewSetupLogger()
	pm := NewPassphraseManager(logger)
	return pm.PassphraseExists()
}

// PromptForExistingPassphrase is a global function that creates a manager and prompts for existing passphrase
func PromptForExistingPassphrase() (string, error) {
	logger := NewSetupLogger()
	pm := NewPassphraseManager(logger)
	return pm.PromptForExistingPassphrase()
}

// ValidatePassphrase is a global function that creates a manager and validates passphrase
func ValidatePassphrase(providedPassphrase string) error {
	logger := NewSetupLogger()
	pm := NewPassphraseManager(logger)
	return pm.ValidatePassphrase(providedPassphrase)
}
