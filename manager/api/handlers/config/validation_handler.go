// Package config provides validation handlers for the API
package config

import (
	"fmt"
	"net/http"
	"time"

	"manager/api/middleware"
	"manager/api/services/validation"
	"manager/api/types"

	"github.com/gin-gonic/gin"
)

// ValidationHandler handles validation-related HTTP requests
type ValidationHandler struct {
	validationService *validation.ValidationService
	logger            middleware.Logger
}

// NewValidationHandler creates a new validation handler
func NewValidationHandler(
	validationService *validation.ValidationService,
	logger middleware.Logger,
) *ValidationHandler {
	return &ValidationHandler{
		validationService: validationService,
		logger:            logger,
	}
}

// ValidateEnvironment validates the environment
// @Summary Validate environment
// @Description Validate the environment for setup compatibility
// @Tags validation
// @Accept json
// @Produce json
// @Param request body types.ValidationRequest true "Validation request"
// @Success 200 {object} types.ValidationResponse
// @Failure 400 {object} types.ValidationResponse
// @Failure 500 {object} types.ValidationResponse
// @Router /api/v1/validation/environment [post]
func (h *ValidationHandler) ValidateEnvironment(c *gin.Context) {
	startTime := time.Now()

	var req types.ValidationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Error("Failed to bind environment validation request", map[string]interface{}{
			"error":     err.Error(),
			"interface": c.GetHeader("X-Interface"),
		})
		c.JSON(http.StatusBadRequest, types.ValidationResponse{
			Success: false,
			Error: &types.ErrorDetail{
				Code:    "INVALID_REQUEST",
				Message: "Invalid request format",
				Details: err.Error(),
			},
			Message: "Invalid request format",
			Code:    http.StatusBadRequest,
		})
		return
	}

	// Validate environment
	validationResult, err := h.validationService.ValidateEnvironment(&req)
	if err != nil {
		h.logger.Error("Environment validation failed", map[string]interface{}{
			"error":     err.Error(),
			"interface": req.Interface,
			"user_id":   req.UserID,
		})
		c.JSON(http.StatusInternalServerError, types.ValidationResponse{
			Success: false,
			Error: &types.ErrorDetail{
				Code:    "VALIDATION_FAILED",
				Message: "Environment validation failed",
				Details: err.Error(),
			},
			Message: "Environment validation failed",
			Code:    http.StatusInternalServerError,
		})
		return
	}

	duration := time.Since(startTime)
	validationResult.Duration = duration

	h.logger.Info("Environment validation completed", map[string]interface{}{
		"interface": req.Interface,
		"user_id":   req.UserID,
		"valid":     validationResult.Valid,
		"errors":    len(validationResult.Errors),
		"warnings":  len(validationResult.Warnings),
		"duration":  duration.String(),
	})

	c.JSON(http.StatusOK, types.ValidationResponse{
		Success: true,
		Result:  validationResult,
		Message: "Environment validation completed",
		Code:    http.StatusOK,
	})
}

// ValidateSecurity validates security aspects
// @Summary Validate security
// @Description Validate security aspects of the environment
// @Tags validation
// @Accept json
// @Produce json
// @Param request body types.ValidationRequest true "Security validation request"
// @Success 200 {object} types.ValidationResponse
// @Failure 400 {object} types.ValidationResponse
// @Failure 500 {object} types.ValidationResponse
// @Router /api/v1/validation/security [post]
func (h *ValidationHandler) ValidateSecurity(c *gin.Context) {
	startTime := time.Now()

	var req types.ValidationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Error("Failed to bind security validation request", map[string]interface{}{
			"error":     err.Error(),
			"interface": c.GetHeader("X-Interface"),
		})
		c.JSON(http.StatusBadRequest, types.ValidationResponse{
			Success: false,
			Error: &types.ErrorDetail{
				Code:    "INVALID_REQUEST",
				Message: "Invalid request format",
				Details: err.Error(),
			},
			Message: "Invalid request format",
			Code:    http.StatusBadRequest,
		})
		return
	}

	// Validate security
	validationResult, err := h.validationService.ValidateSecurity(&req)
	if err != nil {
		h.logger.Error("Security validation failed", map[string]interface{}{
			"error":     err.Error(),
			"interface": req.Interface,
			"user_id":   req.UserID,
		})
		c.JSON(http.StatusInternalServerError, types.ValidationResponse{
			Success: false,
			Error: &types.ErrorDetail{
				Code:    "VALIDATION_FAILED",
				Message: "Security validation failed",
				Details: err.Error(),
			},
			Message: "Security validation failed",
			Code:    http.StatusInternalServerError,
		})
		return
	}

	duration := time.Since(startTime)
	validationResult.Duration = duration

	h.logger.Info("Security validation completed", map[string]interface{}{
		"interface": req.Interface,
		"user_id":   req.UserID,
		"valid":     validationResult.Valid,
		"errors":    len(validationResult.Errors),
		"warnings":  len(validationResult.Warnings),
		"duration":  duration.String(),
	})

	c.JSON(http.StatusOK, types.ValidationResponse{
		Success: true,
		Result:  validationResult,
		Message: "Security validation completed",
		Code:    http.StatusOK,
	})
}

// ValidatePerformance validates performance aspects
// @Summary Validate performance
// @Description Validate performance aspects of the environment
// @Tags validation
// @Accept json
// @Produce json
// @Param request body types.ValidationRequest true "Performance validation request"
// @Success 200 {object} types.ValidationResponse
// @Failure 400 {object} types.ValidationResponse
// @Failure 500 {object} types.ValidationResponse
// @Router /api/v1/validation/performance [post]
func (h *ValidationHandler) ValidatePerformance(c *gin.Context) {
	startTime := time.Now()

	var req types.ValidationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Error("Failed to bind performance validation request", map[string]interface{}{
			"error":     err.Error(),
			"interface": c.GetHeader("X-Interface"),
		})
		c.JSON(http.StatusBadRequest, types.ValidationResponse{
			Success: false,
			Error: &types.ErrorDetail{
				Code:    "INVALID_REQUEST",
				Message: "Invalid request format",
				Details: err.Error(),
			},
			Message: "Invalid request format",
			Code:    http.StatusBadRequest,
		})
		return
	}

	// Validate performance
	validationResult, err := h.validationService.ValidatePerformance(&req)
	if err != nil {
		h.logger.Error("Performance validation failed", map[string]interface{}{
			"error":     err.Error(),
			"interface": req.Interface,
			"user_id":   req.UserID,
		})
		c.JSON(http.StatusInternalServerError, types.ValidationResponse{
			Success: false,
			Error: &types.ErrorDetail{
				Code:    "VALIDATION_FAILED",
				Message: "Performance validation failed",
				Details: err.Error(),
			},
			Message: "Performance validation failed",
			Code:    http.StatusInternalServerError,
		})
		return
	}

	duration := time.Since(startTime)
	validationResult.Duration = duration

	h.logger.Info("Performance validation completed", map[string]interface{}{
		"interface": req.Interface,
		"user_id":   req.UserID,
		"valid":     validationResult.Valid,
		"errors":    len(validationResult.Errors),
		"warnings":  len(validationResult.Warnings),
		"duration":  duration.String(),
	})

	c.JSON(http.StatusOK, types.ValidationResponse{
		Success: true,
		Result:  validationResult,
		Message: "Performance validation completed",
		Code:    http.StatusOK,
	})
}

// ValidateCompatibility validates compatibility aspects
// @Summary Validate compatibility
// @Description Validate compatibility aspects of the environment
// @Tags validation
// @Accept json
// @Produce json
// @Param request body types.ValidationRequest true "Compatibility validation request"
// @Success 200 {object} types.ValidationResponse
// @Failure 400 {object} types.ValidationResponse
// @Failure 500 {object} types.ValidationResponse
// @Router /api/v1/validation/compatibility [post]
func (h *ValidationHandler) ValidateCompatibility(c *gin.Context) {
	startTime := time.Now()

	var req types.ValidationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Error("Failed to bind compatibility validation request", map[string]interface{}{
			"error":     err.Error(),
			"interface": c.GetHeader("X-Interface"),
		})
		c.JSON(http.StatusBadRequest, types.ValidationResponse{
			Success: false,
			Error: &types.ErrorDetail{
				Code:    "INVALID_REQUEST",
				Message: "Invalid request format",
				Details: err.Error(),
			},
			Message: "Invalid request format",
			Code:    http.StatusBadRequest,
		})
		return
	}

	// Validate compatibility
	validationResult, err := h.validationService.ValidateCompatibility(&req)
	if err != nil {
		h.logger.Error("Compatibility validation failed", map[string]interface{}{
			"error":     err.Error(),
			"interface": req.Interface,
			"user_id":   req.UserID,
		})
		c.JSON(http.StatusInternalServerError, types.ValidationResponse{
			Success: false,
			Error: &types.ErrorDetail{
				Code:    "VALIDATION_FAILED",
				Message: "Compatibility validation failed",
				Details: err.Error(),
			},
			Message: "Compatibility validation failed",
			Code:    http.StatusInternalServerError,
		})
		return
	}

	duration := time.Since(startTime)
	validationResult.Duration = duration

	h.logger.Info("Compatibility validation completed", map[string]interface{}{
		"interface": req.Interface,
		"user_id":   req.UserID,
		"valid":     validationResult.Valid,
		"errors":    len(validationResult.Errors),
		"warnings":  len(validationResult.Warnings),
		"duration":  duration.String(),
	})

	c.JSON(http.StatusOK, types.ValidationResponse{
		Success: true,
		Result:  validationResult,
		Message: "Compatibility validation completed",
		Code:    http.StatusOK,
	})
}

// ValidateDependencies validates dependencies
// @Summary Validate dependencies
// @Description Validate required and optional dependencies
// @Tags validation
// @Accept json
// @Produce json
// @Param request body types.ValidationRequest true "Dependencies validation request"
// @Success 200 {object} types.ValidationResponse
// @Failure 400 {object} types.ValidationResponse
// @Failure 500 {object} types.ValidationResponse
// @Router /api/v1/validation/dependencies [post]
func (h *ValidationHandler) ValidateDependencies(c *gin.Context) {
	startTime := time.Now()

	var req types.ValidationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Error("Failed to bind dependencies validation request", map[string]interface{}{
			"error":     err.Error(),
			"interface": c.GetHeader("X-Interface"),
		})
		c.JSON(http.StatusBadRequest, types.ValidationResponse{
			Success: false,
			Error: &types.ErrorDetail{
				Code:    "INVALID_REQUEST",
				Message: "Invalid request format",
				Details: err.Error(),
			},
			Message: "Invalid request format",
			Code:    http.StatusBadRequest,
		})
		return
	}

	// Validate dependencies
	validationResult, err := h.validationService.ValidateDependencies(&req)
	if err != nil {
		h.logger.Error("Dependencies validation failed", map[string]interface{}{
			"error":     err.Error(),
			"interface": req.Interface,
			"user_id":   req.UserID,
		})
		c.JSON(http.StatusInternalServerError, types.ValidationResponse{
			Success: false,
			Error: &types.ErrorDetail{
				Code:    "VALIDATION_FAILED",
				Message: "Dependencies validation failed",
				Details: err.Error(),
			},
			Message: "Dependencies validation failed",
			Code:    http.StatusInternalServerError,
		})
		return
	}

	duration := time.Since(startTime)
	validationResult.Duration = duration

	h.logger.Info("Dependencies validation completed", map[string]interface{}{
		"interface": req.Interface,
		"user_id":   req.UserID,
		"valid":     validationResult.Valid,
		"errors":    len(validationResult.Errors),
		"warnings":  len(validationResult.Warnings),
		"duration":  duration.String(),
	})

	c.JSON(http.StatusOK, types.ValidationResponse{
		Success: true,
		Result:  validationResult,
		Message: "Dependencies validation completed",
		Code:    http.StatusOK,
	})
}

// ValidateAll performs comprehensive validation
// @Summary Validate all aspects
// @Description Perform comprehensive validation of all aspects
// @Tags validation
// @Accept json
// @Produce json
// @Param request body types.ValidationRequest true "Comprehensive validation request"
// @Success 200 {object} types.ValidationResponse
// @Failure 400 {object} types.ValidationResponse
// @Failure 500 {object} types.ValidationResponse
// @Router /api/v1/validation/all [post]
func (h *ValidationHandler) ValidateAll(c *gin.Context) {
	startTime := time.Now()

	var req types.ValidationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Error("Failed to bind comprehensive validation request", map[string]interface{}{
			"error":     err.Error(),
			"interface": c.GetHeader("X-Interface"),
		})
		c.JSON(http.StatusBadRequest, types.ValidationResponse{
			Success: false,
			Error: &types.ErrorDetail{
				Code:    "INVALID_REQUEST",
				Message: "Invalid request format",
				Details: err.Error(),
			},
			Message: "Invalid request format",
			Code:    http.StatusBadRequest,
		})
		return
	}

	// Perform comprehensive validation
	validationResult, err := h.validationService.ValidateAll(&req)
	if err != nil {
		h.logger.Error("Comprehensive validation failed", map[string]interface{}{
			"error":     err.Error(),
			"interface": req.Interface,
			"user_id":   req.UserID,
		})
		c.JSON(http.StatusInternalServerError, types.ValidationResponse{
			Success: false,
			Error: &types.ErrorDetail{
				Code:    "VALIDATION_FAILED",
				Message: "Comprehensive validation failed",
				Details: err.Error(),
			},
			Message: "Comprehensive validation failed",
			Code:    http.StatusInternalServerError,
		})
		return
	}

	duration := time.Since(startTime)
	validationResult.Duration = duration

	h.logger.Info("Comprehensive validation completed", map[string]interface{}{
		"interface": req.Interface,
		"user_id":   req.UserID,
		"valid":     validationResult.Valid,
		"errors":    len(validationResult.Errors),
		"warnings":  len(validationResult.Warnings),
		"duration":  duration.String(),
	})

	c.JSON(http.StatusOK, types.ValidationResponse{
		Success: true,
		Result:  validationResult,
		Message: "Comprehensive validation completed",
		Code:    http.StatusOK,
	})
}

// AutoFix attempts to automatically fix validation issues
// @Summary Auto-fix validation issues
// @Description Attempt to automatically fix validation issues
// @Tags validation
// @Accept json
// @Produce json
// @Param request body types.ValidationRequest true "Auto-fix request"
// @Success 200 {object} types.Response
// @Failure 400 {object} types.Response
// @Failure 500 {object} types.Response
// @Router /api/v1/validation/autofix [post]
func (h *ValidationHandler) AutoFix(c *gin.Context) {
	startTime := time.Now()

	var req types.ValidationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Error("Failed to bind auto-fix request", map[string]interface{}{
			"error":     err.Error(),
			"interface": c.GetHeader("X-Interface"),
		})
		c.JSON(http.StatusBadRequest, types.Response{
			Success: false,
			Error: &types.ErrorDetail{
				Code:    "INVALID_REQUEST",
				Message: "Invalid request format",
				Details: err.Error(),
			},
			Message: "Invalid request format",
			Code:    http.StatusBadRequest,
		})
		return
	}

	// Enable auto-fix in options
	if req.Options.AutoFix == false {
		req.Options.AutoFix = true
	}

	// Perform validation with auto-fix
	validationResult, err := h.validationService.ValidateAll(&req)
	if err != nil {
		h.logger.Error("Auto-fix validation failed", map[string]interface{}{
			"error":     err.Error(),
			"interface": req.Interface,
			"user_id":   req.UserID,
		})
		c.JSON(http.StatusInternalServerError, types.Response{
			Success: false,
			Error: &types.ErrorDetail{
				Code:    "AUTOFIX_FAILED",
				Message: "Auto-fix validation failed",
				Details: err.Error(),
			},
			Message: "Auto-fix validation failed",
			Code:    http.StatusInternalServerError,
		})
		return
	}

	duration := time.Since(startTime)
	validationResult.Duration = duration

	// Count fixed issues
	fixedCount := 0
	for _, item := range validationResult.Warnings {
		if item.AutoFix != nil && item.AutoFix.Available {
			fixedCount++
		}
	}
	for _, item := range validationResult.Errors {
		if item.AutoFix != nil && item.AutoFix.Available {
			fixedCount++
		}
	}

	h.logger.Info("Auto-fix validation completed", map[string]interface{}{
		"interface":   req.Interface,
		"user_id":     req.UserID,
		"valid":       validationResult.Valid,
		"fixed_count": fixedCount,
		"errors":      len(validationResult.Errors),
		"warnings":    len(validationResult.Warnings),
		"duration":    duration.String(),
	})

	c.JSON(http.StatusOK, types.Response{
		Success: true,
		Data: map[string]interface{}{
			"validation_result": validationResult,
			"fixed_count":       fixedCount,
			"duration":          duration.String(),
		},
		Message: fmt.Sprintf("Auto-fix completed. Fixed %d issues.", fixedCount),
		Code:    http.StatusOK,
	})
}
