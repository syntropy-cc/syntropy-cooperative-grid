// Package validation provides validation services for the API
package validation

import (
	"fmt"
	"runtime"
	"sync"
	"time"

	"manager/api/middleware"
	"manager/api/services/validation/compatibility"
	"manager/api/services/validation/dependencies"
	"manager/api/services/validation/environment"
	"manager/api/services/validation/performance"
	"manager/api/services/validation/security"
	"manager/api/types"
)

// ValidationService provides validation services for all interfaces
type ValidationService struct {
	environmentValidator   *environment.EnvironmentValidator
	securityValidator      *security.SecurityValidator
	performanceValidator   *performance.PerformanceValidator
	compatibilityValidator *compatibility.CompatibilityValidator
	dependenciesValidator  *dependencies.DependenciesValidator
	logger                 middleware.Logger
}

// NewValidationService creates a new validation service
func NewValidationService(logger middleware.Logger) *ValidationService {
	return &ValidationService{
		environmentValidator:   environment.NewEnvironmentValidator(logger),
		securityValidator:      security.NewSecurityValidator(logger),
		performanceValidator:   performance.NewPerformanceValidator(logger),
		compatibilityValidator: compatibility.NewCompatibilityValidator(logger),
		dependenciesValidator:  dependencies.NewDependenciesValidator(logger),
		logger:                 logger,
	}
}

// ValidateEnvironment validates the environment for setup compatibility
func (vs *ValidationService) ValidateEnvironment(req *types.ValidationRequest) (*types.ValidationResult, error) {
	startTime := time.Now()

	vs.logger.Info("Starting environment validation", map[string]interface{}{
		"interface": req.Interface,
		"user_id":   req.UserID,
		"os":        runtime.GOOS,
	})

	result := &types.ValidationResult{
		Interface:     req.Interface,
		Version:       "1.0.0",
		Timestamp:     startTime,
		Environment:   &types.EnvironmentInfo{},
		Resources:     &types.SystemResources{},
		Compatibility: &types.Compatibility{},
		Security:      &types.SecurityCheck{},
		Performance:   &types.PerformanceCheck{},
		Warnings:      []types.ValidationItem{},
		Errors:        []types.ValidationItem{},
	}

	// Validate environment based on OS
	switch runtime.GOOS {
	case "windows":
		if err := vs.environmentValidator.ValidateWindows(req, result); err != nil {
			vs.logger.Error("Windows environment validation failed", map[string]interface{}{
				"error":     err.Error(),
				"interface": req.Interface,
			})
			return nil, fmt.Errorf("windows environment validation failed: %w", err)
		}
	case "linux":
		if err := vs.environmentValidator.ValidateLinux(req, result); err != nil {
			vs.logger.Error("Linux environment validation failed", map[string]interface{}{
				"error":     err.Error(),
				"interface": req.Interface,
			})
			return nil, fmt.Errorf("linux environment validation failed: %w", err)
		}
	case "darwin":
		if err := vs.environmentValidator.ValidateDarwin(req, result); err != nil {
			vs.logger.Error("macOS environment validation failed", map[string]interface{}{
				"error":     err.Error(),
				"interface": req.Interface,
			})
			return nil, fmt.Errorf("macos environment validation failed: %w", err)
		}
	default:
		result.Errors = append(result.Errors, types.ValidationItem{
			Code:      "UNSUPPORTED_OS",
			Message:   fmt.Sprintf("Operating system '%s' is not supported", runtime.GOOS),
			Severity:  string(types.SeverityError),
			Category:  string(types.CategoryCompatibility),
			Fixable:   false,
			Timestamp: time.Now(),
		})
		result.Valid = false
	}

	// Set overall validation result
	result.Valid = len(result.Errors) == 0
	result.Duration = time.Since(startTime)

	vs.logger.Info("Environment validation completed", map[string]interface{}{
		"interface": req.Interface,
		"valid":     result.Valid,
		"errors":    len(result.Errors),
		"warnings":  len(result.Warnings),
		"duration":  result.Duration.String(),
	})

	return result, nil
}

// ValidateSecurity validates security aspects of the environment
func (vs *ValidationService) ValidateSecurity(req *types.ValidationRequest) (*types.ValidationResult, error) {
	startTime := time.Now()

	vs.logger.Info("Starting security validation", map[string]interface{}{
		"interface": req.Interface,
		"user_id":   req.UserID,
	})

	result := &types.ValidationResult{
		Interface: req.Interface,
		Version:   "1.0.0",
		Timestamp: startTime,
		Security:  &types.SecurityCheck{},
		Warnings:  []types.ValidationItem{},
		Errors:    []types.ValidationItem{},
	}

	// Validate security aspects
	if err := vs.securityValidator.Validate(req, result); err != nil {
		vs.logger.Error("Security validation failed", map[string]interface{}{
			"error":     err.Error(),
			"interface": req.Interface,
		})
		return nil, fmt.Errorf("security validation failed: %w", err)
	}

	// Set overall validation result
	result.Valid = len(result.Errors) == 0
	result.Duration = time.Since(startTime)

	vs.logger.Info("Security validation completed", map[string]interface{}{
		"interface": req.Interface,
		"valid":     result.Valid,
		"errors":    len(result.Errors),
		"warnings":  len(result.Warnings),
		"duration":  result.Duration.String(),
	})

	return result, nil
}

// ValidatePerformance validates performance aspects of the environment
func (vs *ValidationService) ValidatePerformance(req *types.ValidationRequest) (*types.ValidationResult, error) {
	startTime := time.Now()

	vs.logger.Info("Starting performance validation", map[string]interface{}{
		"interface": req.Interface,
		"user_id":   req.UserID,
	})

	result := &types.ValidationResult{
		Interface:   req.Interface,
		Version:     "1.0.0",
		Timestamp:   startTime,
		Performance: &types.PerformanceCheck{},
		Warnings:    []types.ValidationItem{},
		Errors:      []types.ValidationItem{},
	}

	// Validate performance aspects
	if err := vs.performanceValidator.Validate(req, result); err != nil {
		vs.logger.Error("Performance validation failed", map[string]interface{}{
			"error":     err.Error(),
			"interface": req.Interface,
		})
		return nil, fmt.Errorf("performance validation failed: %w", err)
	}

	// Set overall validation result
	result.Valid = len(result.Errors) == 0
	result.Duration = time.Since(startTime)

	vs.logger.Info("Performance validation completed", map[string]interface{}{
		"interface": req.Interface,
		"valid":     result.Valid,
		"errors":    len(result.Errors),
		"warnings":  len(result.Warnings),
		"duration":  result.Duration.String(),
	})

	return result, nil
}

// ValidateCompatibility validates compatibility aspects
func (vs *ValidationService) ValidateCompatibility(req *types.ValidationRequest) (*types.ValidationResult, error) {
	startTime := time.Now()

	vs.logger.Info("Starting compatibility validation", map[string]interface{}{
		"interface": req.Interface,
		"user_id":   req.UserID,
	})

	result := &types.ValidationResult{
		Interface:     req.Interface,
		Version:       "1.0.0",
		Timestamp:     startTime,
		Compatibility: &types.Compatibility{},
		Warnings:      []types.ValidationItem{},
		Errors:        []types.ValidationItem{},
	}

	// Validate compatibility aspects
	if err := vs.compatibilityValidator.Validate(req, result); err != nil {
		vs.logger.Error("Compatibility validation failed", map[string]interface{}{
			"error":     err.Error(),
			"interface": req.Interface,
		})
		return nil, fmt.Errorf("compatibility validation failed: %w", err)
	}

	// Set overall validation result
	result.Valid = len(result.Errors) == 0
	result.Duration = time.Since(startTime)

	vs.logger.Info("Compatibility validation completed", map[string]interface{}{
		"interface": req.Interface,
		"valid":     result.Valid,
		"errors":    len(result.Errors),
		"warnings":  len(result.Warnings),
		"duration":  result.Duration.String(),
	})

	return result, nil
}

// ValidateDependencies validates required and optional dependencies
func (vs *ValidationService) ValidateDependencies(req *types.ValidationRequest) (*types.ValidationResult, error) {
	startTime := time.Now()

	vs.logger.Info("Starting dependencies validation", map[string]interface{}{
		"interface": req.Interface,
		"user_id":   req.UserID,
	})

	result := &types.ValidationResult{
		Interface: req.Interface,
		Version:   "1.0.0",
		Timestamp: startTime,
		Warnings:  []types.ValidationItem{},
		Errors:    []types.ValidationItem{},
	}

	// Validate dependencies
	if err := vs.dependenciesValidator.Validate(req, result); err != nil {
		vs.logger.Error("Dependencies validation failed", map[string]interface{}{
			"error":     err.Error(),
			"interface": req.Interface,
		})
		return nil, fmt.Errorf("dependencies validation failed: %w", err)
	}

	// Set overall validation result
	result.Valid = len(result.Errors) == 0
	result.Duration = time.Since(startTime)

	vs.logger.Info("Dependencies validation completed", map[string]interface{}{
		"interface": req.Interface,
		"valid":     result.Valid,
		"errors":    len(result.Errors),
		"warnings":  len(result.Warnings),
		"duration":  result.Duration.String(),
	})

	return result, nil
}

// ValidateConfig validates a configuration
func (vs *ValidationService) ValidateConfig(req *types.ValidationRequest, config *types.SetupConfig) (*types.ValidationResult, error) {
	startTime := time.Now()

	vs.logger.Info("Starting configuration validation", map[string]interface{}{
		"interface": req.Interface,
		"user_id":   req.UserID,
	})

	result := &types.ValidationResult{
		Interface: req.Interface,
		Version:   "1.0.0",
		Timestamp: startTime,
		Warnings:  []types.ValidationItem{},
		Errors:    []types.ValidationItem{},
	}

	// Validate configuration structure
	if config == nil {
		result.Errors = append(result.Errors, types.ValidationItem{
			Code:      "NULL_CONFIG",
			Message:   "Configuration is null",
			Severity:  string(types.SeverityError),
			Category:  string(types.CategoryConfiguration),
			Fixable:   false,
			Timestamp: time.Now(),
		})
		result.Valid = false
		result.Duration = time.Since(startTime)
		return result, nil
	}

	// Validate manager configuration
	if err := vs.validateManagerConfig(config.Manager, result); err != nil {
		vs.logger.Error("Manager configuration validation failed", map[string]interface{}{
			"error":     err.Error(),
			"interface": req.Interface,
		})
		return nil, fmt.Errorf("manager configuration validation failed: %w", err)
	}

	// Validate owner key configuration
	if err := vs.validateOwnerKeyConfig(config.OwnerKey, result); err != nil {
		vs.logger.Error("Owner key configuration validation failed", map[string]interface{}{
			"error":     err.Error(),
			"interface": req.Interface,
		})
		return nil, fmt.Errorf("owner key configuration validation failed: %w", err)
	}

	// Validate security configuration
	if err := vs.validateSecurityConfig(config.Security, result); err != nil {
		vs.logger.Error("Security configuration validation failed", map[string]interface{}{
			"error":     err.Error(),
			"interface": req.Interface,
		})
		return nil, fmt.Errorf("security configuration validation failed: %w", err)
	}

	// Validate network configuration
	if err := vs.validateNetworkConfig(config.Network, result); err != nil {
		vs.logger.Error("Network configuration validation failed", map[string]interface{}{
			"error":     err.Error(),
			"interface": req.Interface,
		})
		return nil, fmt.Errorf("network configuration validation failed: %w", err)
	}

	// Set overall validation result
	result.Valid = len(result.Errors) == 0
	result.Duration = time.Since(startTime)

	vs.logger.Info("Configuration validation completed", map[string]interface{}{
		"interface": req.Interface,
		"valid":     result.Valid,
		"errors":    len(result.Errors),
		"warnings":  len(result.Warnings),
		"duration":  result.Duration.String(),
	})

	return result, nil
}

// ValidateAll performs comprehensive validation of all aspects
func (vs *ValidationService) ValidateAll(req *types.ValidationRequest) (*types.ValidationResult, error) {
	startTime := time.Now()

	vs.logger.Info("Starting comprehensive validation", map[string]interface{}{
		"interface": req.Interface,
		"user_id":   req.UserID,
		"parallel":  req.Options.Parallel,
	})

	// Create comprehensive result
	result := &types.ValidationResult{
		Interface:     req.Interface,
		Version:       "1.0.0",
		Timestamp:     startTime,
		Environment:   &types.EnvironmentInfo{},
		Resources:     &types.SystemResources{},
		Compatibility: &types.Compatibility{},
		Security:      &types.SecurityCheck{},
		Performance:   &types.PerformanceCheck{},
		Warnings:      []types.ValidationItem{},
		Errors:        []types.ValidationItem{},
	}

	// Perform validations
	if req.Options.Parallel {
		if err := vs.validateAllParallel(req, result); err != nil {
			return nil, fmt.Errorf("parallel validation failed: %w", err)
		}
	} else {
		if err := vs.validateAllSequential(req, result); err != nil {
			return nil, fmt.Errorf("sequential validation failed: %w", err)
		}
	}

	// Set overall validation result
	result.Valid = len(result.Errors) == 0
	result.Duration = time.Since(startTime)

	vs.logger.Info("Comprehensive validation completed", map[string]interface{}{
		"interface": req.Interface,
		"valid":     result.Valid,
		"errors":    len(result.Errors),
		"warnings":  len(result.Warnings),
		"duration":  result.Duration.String(),
	})

	return result, nil
}

// validateAllParallel performs all validations in parallel
func (vs *ValidationService) validateAllParallel(req *types.ValidationRequest, result *types.ValidationResult) error {
	var wg sync.WaitGroup
	var mu sync.Mutex
	var validationErrors []error

	// Environment validation
	wg.Add(1)
	go func() {
		defer wg.Done()
		envResult, err := vs.ValidateEnvironment(req)
		if err != nil {
			mu.Lock()
			validationErrors = append(validationErrors, fmt.Errorf("environment validation: %w", err))
			mu.Unlock()
			return
		}

		mu.Lock()
		result.Environment = envResult.Environment
		result.Resources = envResult.Resources
		result.Compatibility = envResult.Compatibility
		result.Warnings = append(result.Warnings, envResult.Warnings...)
		result.Errors = append(result.Errors, envResult.Errors...)
		mu.Unlock()
	}()

	// Security validation
	wg.Add(1)
	go func() {
		defer wg.Done()
		secResult, err := vs.ValidateSecurity(req)
		if err != nil {
			mu.Lock()
			validationErrors = append(validationErrors, fmt.Errorf("security validation: %w", err))
			mu.Unlock()
			return
		}

		mu.Lock()
		result.Security = secResult.Security
		result.Warnings = append(result.Warnings, secResult.Warnings...)
		result.Errors = append(result.Errors, secResult.Errors...)
		mu.Unlock()
	}()

	// Performance validation
	wg.Add(1)
	go func() {
		defer wg.Done()
		perfResult, err := vs.ValidatePerformance(req)
		if err != nil {
			mu.Lock()
			validationErrors = append(validationErrors, fmt.Errorf("performance validation: %w", err))
			mu.Unlock()
			return
		}

		mu.Lock()
		result.Performance = perfResult.Performance
		result.Warnings = append(result.Warnings, perfResult.Warnings...)
		result.Errors = append(result.Errors, perfResult.Errors...)
		mu.Unlock()
	}()

	// Dependencies validation
	wg.Add(1)
	go func() {
		defer wg.Done()
		depResult, err := vs.ValidateDependencies(req)
		if err != nil {
			mu.Lock()
			validationErrors = append(validationErrors, fmt.Errorf("dependencies validation: %w", err))
			mu.Unlock()
			return
		}

		mu.Lock()
		result.Warnings = append(result.Warnings, depResult.Warnings...)
		result.Errors = append(result.Errors, depResult.Errors...)
		mu.Unlock()
	}()

	wg.Wait()

	if len(validationErrors) > 0 {
		return fmt.Errorf("validation errors: %v", validationErrors)
	}

	return nil
}

// validateAllSequential performs all validations sequentially
func (vs *ValidationService) validateAllSequential(req *types.ValidationRequest, result *types.ValidationResult) error {
	// Environment validation
	envResult, err := vs.ValidateEnvironment(req)
	if err != nil {
		return fmt.Errorf("environment validation: %w", err)
	}
	result.Environment = envResult.Environment
	result.Resources = envResult.Resources
	result.Compatibility = envResult.Compatibility
	result.Warnings = append(result.Warnings, envResult.Warnings...)
	result.Errors = append(result.Errors, envResult.Errors...)

	// Security validation
	secResult, err := vs.ValidateSecurity(req)
	if err != nil {
		return fmt.Errorf("security validation: %w", err)
	}
	result.Security = secResult.Security
	result.Warnings = append(result.Warnings, secResult.Warnings...)
	result.Errors = append(result.Errors, secResult.Errors...)

	// Performance validation
	perfResult, err := vs.ValidatePerformance(req)
	if err != nil {
		return fmt.Errorf("performance validation: %w", err)
	}
	result.Performance = perfResult.Performance
	result.Warnings = append(result.Warnings, perfResult.Warnings...)
	result.Errors = append(result.Errors, perfResult.Errors...)

	// Dependencies validation
	depResult, err := vs.ValidateDependencies(req)
	if err != nil {
		return fmt.Errorf("dependencies validation: %w", err)
	}
	result.Warnings = append(result.Warnings, depResult.Warnings...)
	result.Errors = append(result.Errors, depResult.Errors...)

	return nil
}

// validateManagerConfig validates manager configuration
func (vs *ValidationService) validateManagerConfig(config types.ManagerConfig, result *types.ValidationResult) error {
	if config.HomeDir == "" {
		result.Errors = append(result.Errors, types.ValidationItem{
			Code:      "MISSING_HOME_DIR",
			Message:   "Manager home directory is required",
			Severity:  string(types.SeverityError),
			Category:  string(types.CategoryConfiguration),
			Field:     "manager.home_dir",
			Fixable:   true,
			Timestamp: time.Now(),
		})
	}

	if config.LogLevel == "" {
		result.Warnings = append(result.Warnings, types.ValidationItem{
			Code:      "MISSING_LOG_LEVEL",
			Message:   "Log level not specified, using default",
			Severity:  string(types.SeverityWarning),
			Category:  string(types.CategoryConfiguration),
			Field:     "manager.log_level",
			Fixable:   true,
			Timestamp: time.Now(),
		})
	}

	return nil
}

// validateOwnerKeyConfig validates owner key configuration
func (vs *ValidationService) validateOwnerKeyConfig(config types.OwnerKey, result *types.ValidationResult) error {
	if config.Path == "" {
		result.Errors = append(result.Errors, types.ValidationItem{
			Code:      "MISSING_KEY_PATH",
			Message:   "Owner key path is required",
			Severity:  string(types.SeverityError),
			Category:  string(types.CategorySecurity),
			Field:     "owner_key.path",
			Fixable:   true,
			Timestamp: time.Now(),
		})
	}

	if config.PublicKey == "" {
		result.Errors = append(result.Errors, types.ValidationItem{
			Code:      "MISSING_PUBLIC_KEY",
			Message:   "Owner public key is required",
			Severity:  string(types.SeverityError),
			Category:  string(types.CategorySecurity),
			Field:     "owner_key.public_key",
			Fixable:   true,
			Timestamp: time.Now(),
		})
	}

	return nil
}

// validateSecurityConfig validates security configuration
func (vs *ValidationService) validateSecurityConfig(config types.SecurityConfig, result *types.ValidationResult) error {
	if config.EncryptionAlgorithm == "" {
		result.Warnings = append(result.Warnings, types.ValidationItem{
			Code:      "MISSING_ENCRYPTION_ALGORITHM",
			Message:   "Encryption algorithm not specified, using default",
			Severity:  string(types.SeverityWarning),
			Category:  string(types.CategorySecurity),
			Field:     "security.encryption_algorithm",
			Fixable:   true,
			Timestamp: time.Now(),
		})
	}

	return nil
}

// validateNetworkConfig validates network configuration
func (vs *ValidationService) validateNetworkConfig(config types.NetworkConfig, result *types.ValidationResult) error {
	if config.Port == 0 {
		result.Warnings = append(result.Warnings, types.ValidationItem{
			Code:      "MISSING_PORT",
			Message:   "Network port not specified, using default",
			Severity:  string(types.SeverityWarning),
			Category:  string(types.CategoryNetwork),
			Field:     "network.port",
			Fixable:   true,
			Timestamp: time.Now(),
		})
	}

	return nil
}
