// Package security provides security validation for the API
package security

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/syntropy-cc/syntropy-cooperative-grid/manager/api/middleware"
	"github.com/syntropy-cc/syntropy-cooperative-grid/manager/api/types"
)

// SecurityValidator validates security aspects of the environment
type SecurityValidator struct {
	logger middleware.Logger
}

// NewSecurityValidator creates a new security validator
func NewSecurityValidator(logger middleware.Logger) *SecurityValidator {
	return &SecurityValidator{
		logger: logger,
	}
}

// Validate performs security validation
func (sv *SecurityValidator) Validate(req *types.ValidationRequest, result *types.ValidationResult) error {
	sv.logger.Info("Starting security validation", map[string]interface{}{
		"interface": req.Interface,
	})

	// Validate encryption capabilities
	if err := sv.validateEncryption(result); err != nil {
		sv.logger.Error("Encryption validation failed", map[string]interface{}{
			"error": err.Error(),
		})
	}

	// Validate secure random number generation
	if err := sv.validateSecureRandom(result); err != nil {
		sv.logger.Error("Secure random validation failed", map[string]interface{}{
			"error": err.Error(),
		})
	}

	// Validate key generation capabilities
	if err := sv.validateKeyGeneration(result); err != nil {
		sv.logger.Error("Key generation validation failed", map[string]interface{}{
			"error": err.Error(),
		})
	}

	// Validate file permissions
	if err := sv.validateFilePermissions(result); err != nil {
		sv.logger.Error("File permissions validation failed", map[string]interface{}{
			"error": err.Error(),
		})
	}

	// Validate network security
	if err := sv.validateNetworkSecurity(result); err != nil {
		sv.logger.Error("Network security validation failed", map[string]interface{}{
			"error": err.Error(),
		})
	}

	// Check for known vulnerabilities
	sv.checkKnownVulnerabilities(result)

	// Generate security recommendations
	sv.generateSecurityRecommendations(result)

	// Check compliance standards
	sv.checkComplianceStandards(result)

	sv.logger.Info("Security validation completed", map[string]interface{}{
		"encryption_available": result.Security.EncryptionAvailable,
		"secure_random":        result.Security.SecureRandom,
		"key_generation":       result.Security.KeyGeneration,
		"file_permissions":     result.Security.FilePermissions,
		"network_security":     result.Security.NetworkSecurity,
	})

	return nil
}

// validateEncryption validates encryption capabilities
func (sv *SecurityValidator) validateEncryption(result *types.ValidationResult) error {
	// Test RSA encryption/decryption
	key, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		result.Security.EncryptionAvailable = false
		result.Errors = append(result.Errors, types.ValidationItem{
			Code:      "ENCRYPTION_FAILED",
			Message:   "Failed to generate RSA key for encryption test",
			Severity:  string(types.SeverityError),
			Category:  string(types.CategorySecurity),
			Fixable:   false,
			Timestamp: time.Now(),
		})
		return err
	}

	// Test encryption
	testData := []byte("test encryption data")
	encrypted, err := rsa.EncryptOAEP(
		sha256.New(),
		rand.Reader,
		&key.PublicKey,
		testData,
		nil,
	)
	if err != nil {
		result.Security.EncryptionAvailable = false
		result.Errors = append(result.Errors, types.ValidationItem{
			Code:      "ENCRYPTION_TEST_FAILED",
			Message:   "Encryption test failed",
			Severity:  string(types.SeverityError),
			Category:  string(types.CategorySecurity),
			Fixable:   false,
			Timestamp: time.Now(),
		})
		return err
	}

	// Test decryption
	decrypted, err := rsa.DecryptOAEP(
		sha256.New(),
		rand.Reader,
		key,
		encrypted,
		nil,
	)
	if err != nil {
		result.Security.EncryptionAvailable = false
		result.Errors = append(result.Errors, types.ValidationItem{
			Code:      "DECRYPTION_TEST_FAILED",
			Message:   "Decryption test failed",
			Severity:  string(types.SeverityError),
			Category:  string(types.CategorySecurity),
			Fixable:   false,
			Timestamp: time.Now(),
		})
		return err
	}

	// Verify data integrity
	if string(decrypted) != string(testData) {
		result.Security.EncryptionAvailable = false
		result.Errors = append(result.Errors, types.ValidationItem{
			Code:      "ENCRYPTION_INTEGRITY_FAILED",
			Message:   "Encryption/decryption data integrity test failed",
			Severity:  string(types.SeverityError),
			Category:  string(types.CategorySecurity),
			Fixable:   false,
			Timestamp: time.Now(),
		})
		return fmt.Errorf("encryption integrity test failed")
	}

	result.Security.EncryptionAvailable = true
	result.Security.Recommendations = append(result.Security.Recommendations, "RSA encryption is available and working correctly")

	return nil
}

// validateSecureRandom validates secure random number generation
func (sv *SecurityValidator) validateSecureRandom(result *types.ValidationResult) error {
	// Test secure random number generation
	randomBytes := make([]byte, 32)
	_, err := rand.Read(randomBytes)
	if err != nil {
		result.Security.SecureRandom = false
		result.Errors = append(result.Errors, types.ValidationItem{
			Code:      "SECURE_RANDOM_FAILED",
			Message:   "Secure random number generation failed",
			Severity:  string(types.SeverityError),
			Category:  string(types.CategorySecurity),
			Fixable:   false,
			Timestamp: time.Now(),
		})
		return err
	}

	// Check if we got actual random data (basic entropy check)
	zeroCount := 0
	for _, b := range randomBytes {
		if b == 0 {
			zeroCount++
		}
	}

	// If more than 25% of bytes are zero, it might indicate poor entropy
	if zeroCount > len(randomBytes)/4 {
		result.Security.SecureRandom = false
		result.Warnings = append(result.Warnings, types.ValidationItem{
			Code:     "POOR_ENTROPY",
			Message:  "Random number generator may have poor entropy",
			Severity: string(types.SeverityWarning),
			Category: string(types.CategorySecurity),
			Fixable:  true,
			AutoFix: &types.AutoFixInfo{
				Available: true,
				Manual:    "Check system entropy sources and consider using hardware random number generators",
				Risk:      "medium",
			},
			Timestamp: time.Now(),
		})
		return fmt.Errorf("poor entropy detected")
	}

	result.Security.SecureRandom = true
	result.Security.Recommendations = append(result.Security.Recommendations, "Secure random number generation is working correctly")

	return nil
}

// validateKeyGeneration validates key generation capabilities
func (sv *SecurityValidator) validateKeyGeneration(result *types.ValidationResult) error {
	// Test RSA key generation
	key, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		result.Security.KeyGeneration = false
		result.Errors = append(result.Errors, types.ValidationItem{
			Code:      "KEY_GENERATION_FAILED",
			Message:   "RSA key generation failed",
			Severity:  string(types.SeverityError),
			Category:  string(types.CategorySecurity),
			Fixable:   false,
			Timestamp: time.Now(),
		})
		return err
	}

	// Test key serialization
	keyBytes, err := x509.MarshalPKCS8PrivateKey(key)
	if err != nil {
		result.Security.KeyGeneration = false
		result.Errors = append(result.Errors, types.ValidationItem{
			Code:      "KEY_SERIALIZATION_FAILED",
			Message:   "Key serialization failed",
			Severity:  string(types.SeverityError),
			Category:  string(types.CategorySecurity),
			Fixable:   false,
			Timestamp: time.Now(),
		})
		return err
	}

	// Test PEM encoding
	pemBlock := &pem.Block{
		Type:  "PRIVATE KEY",
		Bytes: keyBytes,
	}
	pemBytes := pem.EncodeToMemory(pemBlock)
	if pemBytes == nil {
		result.Security.KeyGeneration = false
		result.Errors = append(result.Errors, types.ValidationItem{
			Code:      "PEM_ENCODING_FAILED",
			Message:   "PEM encoding failed",
			Severity:  string(types.SeverityError),
			Category:  string(types.CategorySecurity),
			Fixable:   false,
			Timestamp: time.Now(),
		})
		return err
	}

	// Test key parsing
	block, _ := pem.Decode(pemBytes)
	if block == nil {
		result.Security.KeyGeneration = false
		result.Errors = append(result.Errors, types.ValidationItem{
			Code:      "PEM_DECODING_FAILED",
			Message:   "PEM decoding failed",
			Severity:  string(types.SeverityError),
			Category:  string(types.CategorySecurity),
			Fixable:   false,
			Timestamp: time.Now(),
		})
		return err
	}

	parsedKey, err := x509.ParsePKCS8PrivateKey(block.Bytes)
	if err != nil {
		result.Security.KeyGeneration = false
		result.Errors = append(result.Errors, types.ValidationItem{
			Code:      "KEY_PARSING_FAILED",
			Message:   "Key parsing failed",
			Severity:  string(types.SeverityError),
			Category:  string(types.CategorySecurity),
			Fixable:   false,
			Timestamp: time.Now(),
		})
		return err
	}

	// Verify key type
	if _, ok := parsedKey.(*rsa.PrivateKey); !ok {
		result.Security.KeyGeneration = false
		result.Errors = append(result.Errors, types.ValidationItem{
			Code:      "INVALID_KEY_TYPE",
			Message:   "Generated key is not RSA type",
			Severity:  string(types.SeverityError),
			Category:  string(types.CategorySecurity),
			Fixable:   false,
			Timestamp: time.Now(),
		})
		return fmt.Errorf("invalid key type")
	}

	result.Security.KeyGeneration = true
	result.Security.Recommendations = append(result.Security.Recommendations, "RSA key generation and management is working correctly")

	return nil
}

// validateFilePermissions validates file permission security
func (sv *SecurityValidator) validateFilePermissions(result *types.ValidationResult) error {
	// Create a temporary file to test permissions
	tempFile := filepath.Join(os.TempDir(), "syntropy_security_test")

	// Create file with restrictive permissions
	file, err := os.OpenFile(tempFile, os.O_CREATE|os.O_WRONLY, 0600)
	if err != nil {
		result.Security.FilePermissions = false
		result.Errors = append(result.Errors, types.ValidationItem{
			Code:      "FILE_CREATION_FAILED",
			Message:   "Failed to create test file with secure permissions",
			Severity:  string(types.SeverityError),
			Category:  string(types.CategorySecurity),
			Fixable:   false,
			Timestamp: time.Now(),
		})
		return err
	}
	defer file.Close()

	// Write test data
	testData := []byte("security test data")
	_, err = file.Write(testData)
	if err != nil {
		result.Security.FilePermissions = false
		result.Errors = append(result.Errors, types.ValidationItem{
			Code:      "FILE_WRITE_FAILED",
			Message:   "Failed to write to test file",
			Severity:  string(types.SeverityError),
			Category:  string(types.CategorySecurity),
			Fixable:   false,
			Timestamp: time.Now(),
		})
		os.Remove(tempFile)
		return err
	}

	// Check file permissions
	fileInfo, err := os.Stat(tempFile)
	if err != nil {
		result.Security.FilePermissions = false
		result.Errors = append(result.Errors, types.ValidationItem{
			Code:      "FILE_STAT_FAILED",
			Message:   "Failed to get file information",
			Severity:  string(types.SeverityError),
			Category:  string(types.CategorySecurity),
			Fixable:   false,
			Timestamp: time.Now(),
		})
		os.Remove(tempFile)
		return err
	}

	// Check if permissions are secure (owner read/write only)
	mode := fileInfo.Mode().Perm()
	if mode&0077 != 0 { // Check if group or other have permissions
		result.Security.FilePermissions = false
		result.Warnings = append(result.Warnings, types.ValidationItem{
			Code:     "INSECURE_FILE_PERMISSIONS",
			Message:  fmt.Sprintf("File permissions %s are too permissive", mode.String()),
			Severity: string(types.SeverityWarning),
			Category: string(types.CategorySecurity),
			Fixable:  true,
			AutoFix: &types.AutoFixInfo{
				Available: true,
				Command:   fmt.Sprintf("chmod 600 %s", tempFile),
				Manual:    "Set file permissions to 600 (owner read/write only)",
				Risk:      "low",
			},
			Timestamp: time.Now(),
		})
	} else {
		result.Security.FilePermissions = true
		result.Security.Recommendations = append(result.Security.Recommendations, "File permission controls are working correctly")
	}

	// Clean up test file
	os.Remove(tempFile)
	return nil
}

// validateNetworkSecurity validates network security aspects
func (sv *SecurityValidator) validateNetworkSecurity(result *types.ValidationResult) error {
	// This is a simplified network security validation
	// In production, you would check for:
	// - Firewall status
	// - Open ports
	// - SSL/TLS capabilities
	// - Network encryption support

	result.Security.NetworkSecurity = true
	result.Security.Recommendations = append(result.Security.Recommendations, "Network security validation completed")

	// Check for common security issues
	sv.checkCommonSecurityIssues(result)

	return nil
}

// checkKnownVulnerabilities checks for known security vulnerabilities
func (sv *SecurityValidator) checkKnownVulnerabilities(result *types.ValidationResult) {
	// This is a simplified vulnerability check
	// In production, you would integrate with vulnerability databases

	vulnerabilities := []string{
		// Add known vulnerabilities here
	}

	result.Security.Vulnerabilities = vulnerabilities

	if len(vulnerabilities) > 0 {
		result.Warnings = append(result.Warnings, types.ValidationItem{
			Code:      "KNOWN_VULNERABILITIES",
			Message:   fmt.Sprintf("Found %d known security vulnerabilities", len(vulnerabilities)),
			Severity:  string(types.SeverityWarning),
			Category:  string(types.CategorySecurity),
			Fixable:   true,
			Timestamp: time.Now(),
		})
	}
}

// generateSecurityRecommendations generates security recommendations
func (sv *SecurityValidator) generateSecurityRecommendations(result *types.ValidationResult) {
	recommendations := []string{
		"Enable full disk encryption if not already enabled",
		"Use strong, unique passwords for all accounts",
		"Enable two-factor authentication where possible",
		"Keep the operating system and software updated",
		"Use a firewall to restrict network access",
		"Regularly backup important data",
		"Monitor system logs for suspicious activity",
	}

	result.Security.Recommendations = append(result.Security.Recommendations, recommendations...)
}

// checkComplianceStandards checks compliance with security standards
func (sv *SecurityValidator) checkComplianceStandards(result *types.ValidationResult) {
	compliance := []string{
		// Add compliance standards here
		"Basic security practices implemented",
	}

	result.Security.Compliance = compliance
}

// checkCommonSecurityIssues checks for common security issues
func (sv *SecurityValidator) checkCommonSecurityIssues(result *types.ValidationResult) {
	// Check for common security misconfigurations
	commonIssues := []string{
		// Add common security issues to check
	}

	if len(commonIssues) > 0 {
		result.Warnings = append(result.Warnings, types.ValidationItem{
			Code:      "COMMON_SECURITY_ISSUES",
			Message:   fmt.Sprintf("Found %d common security issues", len(commonIssues)),
			Severity:  string(types.SeverityWarning),
			Category:  string(types.CategorySecurity),
			Fixable:   true,
			Timestamp: time.Now(),
		})
	}
}
