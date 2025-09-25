// Package config provides setup-specific handlers for the API
package config

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"manager/api/middleware"
	"manager/api/services/config"
	"manager/api/services/validation"
	"manager/api/types"

	"github.com/gin-gonic/gin"
)

// SetupHandler handles setup-specific HTTP requests
type SetupHandler struct {
	configService     *config.ConfigService
	validationService *validation.ValidationService
	setupService      *config.SetupService
	logger            middleware.Logger
}

// NewSetupHandler creates a new setup handler
func NewSetupHandler(
	configService *config.ConfigService,
	validationService *validation.ValidationService,
	setupService *config.SetupService,
	logger middleware.Logger,
) *SetupHandler {
	return &SetupHandler{
		configService:     configService,
		validationService: validationService,
		setupService:      setupService,
		logger:            logger,
	}
}

// Setup performs a complete setup process
// @Summary Perform setup
// @Description Perform a complete setup process for the specified interface
// @Tags setup
// @Accept json
// @Produce json
// @Param request body types.SetupRequest true "Setup request"
// @Success 200 {object} types.SetupResponse
// @Failure 400 {object} types.SetupResponse
// @Failure 500 {object} types.SetupResponse
// @Router /api/v1/setup/execute [post]
func (h *SetupHandler) Setup(c *gin.Context) {
	startTime := time.Now()

	var req types.SetupRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Error("Failed to bind setup request", map[string]interface{}{
			"error":     err.Error(),
			"interface": c.GetHeader("X-Interface"),
		})
		c.JSON(http.StatusBadRequest, types.SetupResponse{
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

	// Validate interface type
	if !h.isValidInterface(req.Interface) {
		h.logger.Warn("Invalid interface type for setup", map[string]interface{}{
			"interface": req.Interface,
			"user_id":   req.UserID,
		})
		c.JSON(http.StatusBadRequest, types.SetupResponse{
			Success: false,
			Error: &types.ErrorDetail{
				Code:    "INVALID_INTERFACE",
				Message: "Invalid interface type",
				Details: fmt.Sprintf("Interface '%s' is not supported", req.Interface),
			},
			Message: "Invalid interface type",
			Code:    http.StatusBadRequest,
		})
		return
	}

	// Check if setup already exists (unless force is specified)
	if !req.Options.Force {
		existingSetup, err := h.setupService.GetExistingSetup(req.Interface, req.UserID)
		if err == nil && existingSetup != nil {
			h.logger.Info("Setup already exists", map[string]interface{}{
				"interface": req.Interface,
				"user_id":   req.UserID,
			})
			c.JSON(http.StatusConflict, types.SetupResponse{
				Success: false,
				Error: &types.ErrorDetail{
					Code:    "SETUP_EXISTS",
					Message: "Setup already exists",
					Details: "A setup already exists for this interface and user. Use force=true to overwrite.",
				},
				Message: "Setup already exists",
				Code:    http.StatusConflict,
			})
			return
		}
	}

	// Perform setup
	setupResult, err := h.setupService.ExecuteSetup(&req)
	if err != nil {
		h.logger.Error("Setup execution failed", map[string]interface{}{
			"error":     err.Error(),
			"interface": req.Interface,
			"user_id":   req.UserID,
		})
		c.JSON(http.StatusInternalServerError, types.SetupResponse{
			Success: false,
			Error: &types.ErrorDetail{
				Code:    "SETUP_FAILED",
				Message: "Setup execution failed",
				Details: err.Error(),
			},
			Message: "Setup execution failed",
			Code:    http.StatusInternalServerError,
		})
		return
	}

	duration := time.Since(startTime)
	setupResult.Duration = duration

	h.logger.Info("Setup completed successfully", map[string]interface{}{
		"interface": req.Interface,
		"user_id":   req.UserID,
		"duration":  duration.String(),
		"success":   setupResult.Success,
	})

	c.JSON(http.StatusOK, types.SetupResponse{
		Success: true,
		Result:  setupResult,
		Message: "Setup completed successfully",
		Code:    http.StatusOK,
	})
}

// ValidateSetup validates the current setup
// @Summary Validate setup
// @Description Validate the current setup for the specified interface
// @Tags setup
// @Accept json
// @Produce json
// @Param request body types.SetupRequest true "Setup validation request"
// @Success 200 {object} types.ValidationResponse
// @Failure 400 {object} types.ValidationResponse
// @Failure 500 {object} types.ValidationResponse
// @Router /api/v1/setup/validate [post]
func (h *SetupHandler) ValidateSetup(c *gin.Context) {
	startTime := time.Now()

	var req types.SetupRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Error("Failed to bind setup validation request", map[string]interface{}{
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

	// Validate setup
	validationResult, err := h.setupService.ValidateSetup(&req)
	if err != nil {
		h.logger.Error("Setup validation failed", map[string]interface{}{
			"error":     err.Error(),
			"interface": req.Interface,
			"user_id":   req.UserID,
		})
		c.JSON(http.StatusInternalServerError, types.ValidationResponse{
			Success: false,
			Error: &types.ErrorDetail{
				Code:    "VALIDATION_FAILED",
				Message: "Setup validation failed",
				Details: err.Error(),
			},
			Message: "Setup validation failed",
			Code:    http.StatusInternalServerError,
		})
		return
	}

	duration := time.Since(startTime)
	validationResult.Duration = duration

	h.logger.Info("Setup validation completed", map[string]interface{}{
		"interface": req.Interface,
		"user_id":   req.UserID,
		"valid":     validationResult.Valid,
		"duration":  duration.String(),
	})

	c.JSON(http.StatusOK, types.ValidationResponse{
		Success: true,
		Result:  validationResult,
		Message: "Setup validation completed",
		Code:    http.StatusOK,
	})
}

// GetSetupStatus gets the current setup status
// @Summary Get setup status
// @Description Get the current setup status for the specified interface
// @Tags setup
// @Accept json
// @Produce json
// @Param interface query string true "Interface type"
// @Param user_id query string false "User ID"
// @Success 200 {object} types.Response
// @Failure 400 {object} types.Response
// @Failure 404 {object} types.Response
// @Failure 500 {object} types.Response
// @Router /api/v1/setup/status [get]
func (h *SetupHandler) GetSetupStatus(c *gin.Context) {
	interfaceType := c.Query("interface")
	userID := c.Query("user_id")

	if interfaceType == "" {
		c.JSON(http.StatusBadRequest, types.Response{
			Success: false,
			Error: &types.ErrorDetail{
				Code:    "MISSING_PARAMETERS",
				Message: "Interface parameter is required",
			},
			Message: "Missing required parameters",
			Code:    http.StatusBadRequest,
		})
		return
	}

	// Get setup status
	status, err := h.setupService.GetSetupStatus(interfaceType, userID)
	if err != nil {
		h.logger.Error("Failed to get setup status", map[string]interface{}{
			"error":     err.Error(),
			"interface": interfaceType,
			"user_id":   userID,
		})

		statusCode := http.StatusInternalServerError
		if err.Error() == "setup not found" {
			statusCode = http.StatusNotFound
		}

		c.JSON(statusCode, types.Response{
			Success: false,
			Error: &types.ErrorDetail{
				Code:    "STATUS_ERROR",
				Message: "Failed to get setup status",
				Details: err.Error(),
			},
			Message: "Failed to get setup status",
			Code:    statusCode,
		})
		return
	}

	c.JSON(http.StatusOK, types.Response{
		Success: true,
		Data:    status,
		Message: "Setup status retrieved successfully",
		Code:    http.StatusOK,
	})
}

// ResetSetup resets the setup for the specified interface
// @Summary Reset setup
// @Description Reset the setup for the specified interface
// @Tags setup
// @Accept json
// @Produce json
// @Param request body types.SetupRequest true "Setup reset request"
// @Success 200 {object} types.Response
// @Failure 400 {object} types.Response
// @Failure 500 {object} types.Response
// @Router /api/v1/setup/reset [post]
func (h *SetupHandler) ResetSetup(c *gin.Context) {
	startTime := time.Now()

	var req types.SetupRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Error("Failed to bind setup reset request", map[string]interface{}{
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

	// Reset setup
	err := h.setupService.ResetSetup(&req)
	if err != nil {
		h.logger.Error("Setup reset failed", map[string]interface{}{
			"error":     err.Error(),
			"interface": req.Interface,
			"user_id":   req.UserID,
		})
		c.JSON(http.StatusInternalServerError, types.Response{
			Success: false,
			Error: &types.ErrorDetail{
				Code:    "RESET_FAILED",
				Message: "Setup reset failed",
				Details: err.Error(),
			},
			Message: "Setup reset failed",
			Code:    http.StatusInternalServerError,
		})
		return
	}

	duration := time.Since(startTime)
	h.logger.Info("Setup reset completed", map[string]interface{}{
		"interface": req.Interface,
		"user_id":   req.UserID,
		"duration":  duration.String(),
	})

	c.JSON(http.StatusOK, types.Response{
		Success: true,
		Message: "Setup reset completed successfully",
		Code:    http.StatusOK,
	})
}

// GetSetupHistory gets the setup history for the specified interface
// @Summary Get setup history
// @Description Get the setup history for the specified interface
// @Tags setup
// @Accept json
// @Produce json
// @Param interface query string true "Interface type"
// @Param user_id query string false "User ID"
// @Param limit query int false "Limit number of entries"
// @Success 200 {object} types.Response
// @Failure 400 {object} types.Response
// @Failure 500 {object} types.Response
// @Router /api/v1/setup/history [get]
func (h *SetupHandler) GetSetupHistory(c *gin.Context) {
	interfaceType := c.Query("interface")
	userID := c.Query("user_id")
	limitStr := c.Query("limit")

	if interfaceType == "" {
		c.JSON(http.StatusBadRequest, types.Response{
			Success: false,
			Error: &types.ErrorDetail{
				Code:    "MISSING_PARAMETERS",
				Message: "Interface parameter is required",
			},
			Message: "Missing required parameters",
			Code:    http.StatusBadRequest,
		})
		return
	}

	// Parse limit
	limit := 10 // default limit
	if limitStr != "" {
		if parsedLimit, err := strconv.Atoi(limitStr); err == nil && parsedLimit > 0 {
			limit = parsedLimit
		}
	}

	// Get setup history
	history, err := h.setupService.GetSetupHistory(interfaceType, userID, limit)
	if err != nil {
		h.logger.Error("Failed to get setup history", map[string]interface{}{
			"error":     err.Error(),
			"interface": interfaceType,
			"user_id":   userID,
		})
		c.JSON(http.StatusInternalServerError, types.Response{
			Success: false,
			Error: &types.ErrorDetail{
				Code:    "HISTORY_ERROR",
				Message: "Failed to get setup history",
				Details: err.Error(),
			},
			Message: "Failed to get setup history",
			Code:    http.StatusInternalServerError,
		})
		return
	}

	c.JSON(http.StatusOK, types.Response{
		Success: true,
		Data:    history,
		Message: "Setup history retrieved successfully",
		Code:    http.StatusOK,
	})
}

// isValidInterface validates if the interface type is supported
func (h *SetupHandler) isValidInterface(interfaceType string) bool {
	switch interfaceType {
	case types.InterfaceCLI, types.InterfaceWeb, types.InterfaceDesktop, types.InterfaceMobile:
		return true
	default:
		return false
	}
}
