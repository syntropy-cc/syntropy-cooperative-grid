// Package config provides configuration handlers for the API
package config

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/syntropy-cc/syntropy-cooperative-grid/manager/api/middleware"
	"github.com/syntropy-cc/syntropy-cooperative-grid/manager/api/services/config"
	"github.com/syntropy-cc/syntropy-cooperative-grid/manager/api/services/validation"
	"github.com/syntropy-cc/syntropy-cooperative-grid/manager/api/types"

	"github.com/gin-gonic/gin"
)

// ConfigHandler handles configuration-related HTTP requests
type ConfigHandler struct {
	configService     *config.ConfigService
	validationService *validation.ValidationService
	logger            middleware.Logger
}

// NewConfigHandler creates a new configuration handler
func NewConfigHandler(
	configService *config.ConfigService,
	validationService *validation.ValidationService,
	logger middleware.Logger,
) *ConfigHandler {
	return &ConfigHandler{
		configService:     configService,
		validationService: validationService,
		logger:            logger,
	}
}

// GenerateSetupConfig generates a setup configuration
// @Summary Generate setup configuration
// @Description Generate a complete setup configuration for the specified interface and environment
// @Tags configuration
// @Accept json
// @Produce json
// @Param request body types.ConfigRequest true "Configuration request"
// @Success 200 {object} types.ConfigResponse
// @Failure 400 {object} types.ConfigResponse
// @Failure 500 {object} types.ConfigResponse
// @Router /api/v1/config/generate [post]
func (h *ConfigHandler) GenerateSetupConfig(c *gin.Context) {
	startTime := time.Now()

	var req types.ConfigRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Error("Failed to bind configuration request", map[string]interface{}{
			"error":     err.Error(),
			"interface": c.GetHeader("X-Interface"),
		})
		c.JSON(http.StatusBadRequest, types.ConfigResponse{
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
		h.logger.Warn("Invalid interface type", map[string]interface{}{
			"interface": req.Interface,
			"user_id":   req.UserID,
		})
		c.JSON(http.StatusBadRequest, types.ConfigResponse{
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

	// Generate configuration
	generatedConfig, err := h.configService.GenerateConfig(&req)
	if err != nil {
		h.logger.Error("Failed to generate configuration", map[string]interface{}{
			"error":     err.Error(),
			"interface": req.Interface,
			"user_id":   req.UserID,
		})
		c.JSON(http.StatusInternalServerError, types.ConfigResponse{
			Success: false,
			Error: &types.ErrorDetail{
				Code:    "GENERATION_FAILED",
				Message: "Failed to generate configuration",
				Details: err.Error(),
			},
			Message: "Configuration generation failed",
			Code:    http.StatusInternalServerError,
		})
		return
	}

	duration := time.Since(startTime)
	h.logger.Info("Configuration generated successfully", map[string]interface{}{
		"interface": req.Interface,
		"user_id":   req.UserID,
		"duration":  duration.String(),
	})

	c.JSON(http.StatusOK, types.ConfigResponse{
		Success: true,
		Config:  generatedConfig,
		Message: "Configuration generated successfully",
		Code:    http.StatusOK,
		Metadata: &types.ConfigMetadata{
			Version:     "1.0.0",
			CreatedAt:   time.Now(),
			Interface:   req.Interface,
			Environment: req.Environment.OS,
			Checksum:    h.calculateChecksum(generatedConfig),
		},
	})
}

// ValidateConfig validates a configuration
// @Summary Validate configuration
// @Description Validate a configuration for correctness and completeness
// @Tags configuration
// @Accept json
// @Produce json
// @Param request body types.ConfigRequest true "Configuration validation request"
// @Success 200 {object} types.ValidationResponse
// @Failure 400 {object} types.ValidationResponse
// @Failure 500 {object} types.ValidationResponse
// @Router /api/v1/config/validate [post]
func (h *ConfigHandler) ValidateConfig(c *gin.Context) {
	startTime := time.Now()

	var req types.ConfigRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Error("Failed to bind validation request", map[string]interface{}{
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

	// Create validation request
	validationReq := &types.ValidationRequest{
		Type:        "configuration",
		Environment: req.Environment,
		Interface:   req.Interface,
		UserID:      req.UserID,
		SessionID:   req.SessionID,
		CustomData:  req.CustomData,
	}

	// Validate configuration
	validationResult, err := h.validationService.ValidateConfig(validationReq, req.Config)
	if err != nil {
		h.logger.Error("Failed to validate configuration", map[string]interface{}{
			"error":     err.Error(),
			"interface": req.Interface,
			"user_id":   req.UserID,
		})
		c.JSON(http.StatusInternalServerError, types.ValidationResponse{
			Success: false,
			Error: &types.ErrorDetail{
				Code:    "VALIDATION_FAILED",
				Message: "Failed to validate configuration",
				Details: err.Error(),
			},
			Message: "Configuration validation failed",
			Code:    http.StatusInternalServerError,
		})
		return
	}

	duration := time.Since(startTime)
	h.logger.Info("Configuration validated", map[string]interface{}{
		"interface": req.Interface,
		"user_id":   req.UserID,
		"valid":     validationResult.Valid,
		"duration":  duration.String(),
	})

	c.JSON(http.StatusOK, types.ValidationResponse{
		Success: true,
		Result:  validationResult,
		Message: "Configuration validation completed",
		Code:    http.StatusOK,
	})
}

// BackupConfig creates a backup of the current configuration
// @Summary Backup configuration
// @Description Create a backup of the current configuration
// @Tags configuration
// @Accept json
// @Produce json
// @Param request body types.ConfigRequest true "Backup request"
// @Success 200 {object} types.Response
// @Failure 400 {object} types.Response
// @Failure 500 {object} types.Response
// @Router /api/v1/config/backup [post]
func (h *ConfigHandler) BackupConfig(c *gin.Context) {
	startTime := time.Now()

	var req types.ConfigRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Error("Failed to bind backup request", map[string]interface{}{
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

	// Create backup
	backup, err := h.configService.CreateBackup(&req)
	if err != nil {
		h.logger.Error("Failed to create configuration backup", map[string]interface{}{
			"error":     err.Error(),
			"interface": req.Interface,
			"user_id":   req.UserID,
		})
		c.JSON(http.StatusInternalServerError, types.Response{
			Success: false,
			Error: &types.ErrorDetail{
				Code:    "BACKUP_FAILED",
				Message: "Failed to create configuration backup",
				Details: err.Error(),
			},
			Message: "Configuration backup failed",
			Code:    http.StatusInternalServerError,
		})
		return
	}

	duration := time.Since(startTime)
	h.logger.Info("Configuration backup created", map[string]interface{}{
		"interface": req.Interface,
		"user_id":   req.UserID,
		"backup_id": backup.ID,
		"duration":  duration.String(),
	})

	c.JSON(http.StatusOK, types.Response{
		Success: true,
		Data:    backup,
		Message: "Configuration backup created successfully",
		Code:    http.StatusOK,
	})
}

// RestoreConfig restores a configuration from backup
// @Summary Restore configuration
// @Description Restore a configuration from backup
// @Tags configuration
// @Accept json
// @Produce json
// @Param request body types.ConfigRestoreRequest true "Restore request"
// @Success 200 {object} types.ConfigRestoreResponse
// @Failure 400 {object} types.ConfigRestoreResponse
// @Failure 500 {object} types.ConfigRestoreResponse
// @Router /api/v1/config/restore [post]
func (h *ConfigHandler) RestoreConfig(c *gin.Context) {
	startTime := time.Now()

	var req types.ConfigRestoreRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Error("Failed to bind restore request", map[string]interface{}{
			"error":     err.Error(),
			"interface": c.GetHeader("X-Interface"),
		})
		c.JSON(http.StatusBadRequest, types.ConfigRestoreResponse{
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

	// Restore configuration
	restoreResult, err := h.configService.RestoreConfig(&req)
	if err != nil {
		h.logger.Error("Failed to restore configuration", map[string]interface{}{
			"error":     err.Error(),
			"interface": c.GetHeader("X-Interface"),
			"backup_id": req.BackupID,
		})
		c.JSON(http.StatusInternalServerError, types.ConfigRestoreResponse{
			Success: false,
			Error: &types.ErrorDetail{
				Code:    "RESTORE_FAILED",
				Message: "Failed to restore configuration",
				Details: err.Error(),
			},
			Message: "Configuration restore failed",
			Code:    http.StatusInternalServerError,
		})
		return
	}

	duration := time.Since(startTime)
	h.logger.Info("Configuration restored", map[string]interface{}{
		"interface": c.GetHeader("X-Interface"),
		"backup_id": req.BackupID,
		"duration":  duration.String(),
	})

	c.JSON(http.StatusOK, restoreResult)
}

// ListConfigs lists available configurations
// @Summary List configurations
// @Description List available configurations with optional filtering and pagination
// @Tags configuration
// @Accept json
// @Produce json
// @Param request query types.ConfigListRequest false "List request parameters"
// @Success 200 {object} types.ConfigListResponse
// @Failure 400 {object} types.ConfigListResponse
// @Failure 500 {object} types.ConfigListResponse
// @Router /api/v1/config/list [get]
func (h *ConfigHandler) ListConfigs(c *gin.Context) {
	startTime := time.Now()

	// Parse query parameters
	req := &types.ConfigListRequest{
		Interface: c.Query("interface"),
		UserID:    c.Query("user_id"),
		SessionID: c.Query("session_id"),
		Pagination: types.PaginationOptions{
			Page:     1,
			PageSize: 10,
		},
		Sort: types.SortOptions{
			Field: "created_at",
			Order: "desc",
		},
	}

	// Parse pagination parameters
	if pageStr := c.Query("page"); pageStr != "" {
		if page, err := strconv.Atoi(pageStr); err == nil && page > 0 {
			req.Pagination.Page = page
		}
	}
	if pageSizeStr := c.Query("page_size"); pageSizeStr != "" {
		if pageSize, err := strconv.Atoi(pageSizeStr); err == nil && pageSize > 0 {
			req.Pagination.PageSize = pageSize
		}
	}

	// Parse sort parameters
	if sortField := c.Query("sort_field"); sortField != "" {
		req.Sort.Field = sortField
	}
	if sortOrder := c.Query("sort_order"); sortOrder != "" {
		req.Sort.Order = sortOrder
	}

	// List configurations
	configs, err := h.configService.ListConfigs(req)
	if err != nil {
		h.logger.Error("Failed to list configurations", map[string]interface{}{
			"error":     err.Error(),
			"interface": req.Interface,
		})
		c.JSON(http.StatusInternalServerError, types.ConfigListResponse{
			Success: false,
			Error: &types.ErrorDetail{
				Code:    "LIST_FAILED",
				Message: "Failed to list configurations",
				Details: err.Error(),
			},
			Message: "Failed to list configurations",
			Code:    http.StatusInternalServerError,
		})
		return
	}

	duration := time.Since(startTime)
	h.logger.Info("Configurations listed", map[string]interface{}{
		"interface": req.Interface,
		"count":     len(configs),
		"duration":  duration.String(),
	})

	c.JSON(http.StatusOK, types.ConfigListResponse{
		Success:    true,
		Configs:    configs,
		Message:    "Configurations listed successfully",
		Code:       http.StatusOK,
		Pagination: &req.Pagination,
	})
}

// GetConfigTemplate gets a configuration template
// @Summary Get configuration template
// @Description Get a configuration template for the specified interface and environment
// @Tags configuration
// @Accept json
// @Produce json
// @Param interface query string true "Interface type"
// @Param environment query string true "Environment"
// @Param template query string false "Template name"
// @Success 200 {object} types.Response
// @Failure 400 {object} types.Response
// @Failure 404 {object} types.Response
// @Failure 500 {object} types.Response
// @Router /api/v1/config/template [get]
func (h *ConfigHandler) GetConfigTemplate(c *gin.Context) {
	interfaceType := c.Query("interface")
	environment := c.Query("environment")
	templateName := c.Query("template")

	if interfaceType == "" || environment == "" {
		c.JSON(http.StatusBadRequest, types.Response{
			Success: false,
			Error: &types.ErrorDetail{
				Code:    "MISSING_PARAMETERS",
				Message: "Interface and environment parameters are required",
			},
			Message: "Missing required parameters",
			Code:    http.StatusBadRequest,
		})
		return
	}

	// Get template
	template, err := h.configService.GetTemplate(interfaceType, environment, templateName)
	if err != nil {
		h.logger.Error("Failed to get configuration template", map[string]interface{}{
			"error":       err.Error(),
			"interface":   interfaceType,
			"environment": environment,
			"template":    templateName,
		})

		statusCode := http.StatusInternalServerError
		if err.Error() == "template not found" {
			statusCode = http.StatusNotFound
		}

		c.JSON(statusCode, types.Response{
			Success: false,
			Error: &types.ErrorDetail{
				Code:    "TEMPLATE_ERROR",
				Message: "Failed to get configuration template",
				Details: err.Error(),
			},
			Message: "Failed to get configuration template",
			Code:    statusCode,
		})
		return
	}

	c.JSON(http.StatusOK, types.Response{
		Success: true,
		Data:    template,
		Message: "Configuration template retrieved successfully",
		Code:    http.StatusOK,
	})
}

// isValidInterface validates if the interface type is supported
func (h *ConfigHandler) isValidInterface(interfaceType string) bool {
	switch interfaceType {
	case types.InterfaceCLI, types.InterfaceWeb, types.InterfaceDesktop, types.InterfaceMobile:
		return true
	default:
		return false
	}
}

// calculateChecksum calculates a simple checksum for the configuration
func (h *ConfigHandler) calculateChecksum(config *types.SetupConfig) string {
	data, err := json.Marshal(config)
	if err != nil {
		return "error"
	}

	// Simple hash calculation (in production, use proper hashing)
	hash := 0
	for _, b := range data {
		hash += int(b)
	}

	return fmt.Sprintf("%x", hash)
}
