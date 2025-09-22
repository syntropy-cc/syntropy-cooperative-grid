package errors

import (
	"fmt"
	"net/http"
)

// ErrorCode represents a specific error code
type ErrorCode string

const (
	// General errors
	ErrCodeInternal     ErrorCode = "INTERNAL_ERROR"
	ErrCodeNotFound     ErrorCode = "NOT_FOUND"
	ErrCodeInvalidInput ErrorCode = "INVALID_INPUT"
	ErrCodeUnauthorized ErrorCode = "UNAUTHORIZED"
	ErrCodeForbidden    ErrorCode = "FORBIDDEN"
	ErrCodeConflict     ErrorCode = "CONFLICT"
	ErrCodeTimeout      ErrorCode = "TIMEOUT"

	// Node errors
	ErrCodeNodeNotFound      ErrorCode = "NODE_NOT_FOUND"
	ErrCodeNodeAlreadyExists ErrorCode = "NODE_ALREADY_EXISTS"
	ErrCodeNodeInvalidStatus ErrorCode = "NODE_INVALID_STATUS"
	ErrCodeNodeCreationFailed ErrorCode = "NODE_CREATION_FAILED"

	// Container errors
	ErrCodeContainerNotFound      ErrorCode = "CONTAINER_NOT_FOUND"
	ErrCodeContainerAlreadyExists ErrorCode = "CONTAINER_ALREADY_EXISTS"
	ErrCodeContainerInvalidStatus ErrorCode = "CONTAINER_INVALID_STATUS"
	ErrCodeContainerDeployFailed  ErrorCode = "CONTAINER_DEPLOY_FAILED"

	// Network errors
	ErrCodeNetworkRouteNotFound ErrorCode = "NETWORK_ROUTE_NOT_FOUND"
	ErrCodeNetworkRouteExists   ErrorCode = "NETWORK_ROUTE_EXISTS"
	ErrCodeNetworkMeshDisabled  ErrorCode = "NETWORK_MESH_DISABLED"

	// Cooperative errors
	ErrCodeInsufficientCredits ErrorCode = "INSUFFICIENT_CREDITS"
	ErrCodeTransactionFailed   ErrorCode = "TRANSACTION_FAILED"
	ErrCodeProposalNotFound    ErrorCode = "PROPOSAL_NOT_FOUND"
	ErrCodeVoteAlreadyCast     ErrorCode = "VOTE_ALREADY_CAST"

	// USB errors
	ErrCodeUSBDeviceNotFound ErrorCode = "USB_DEVICE_NOT_FOUND"
	ErrCodeUSBDeviceInUse    ErrorCode = "USB_DEVICE_IN_USE"
	ErrCodeUSBFormatFailed   ErrorCode = "USB_FORMAT_FAILED"
)

// APIError represents an API error with additional context
type APIError struct {
	Code       ErrorCode `json:"code"`
	Message    string    `json:"message"`
	Details    string    `json:"details,omitempty"`
	HTTPStatus int       `json:"-"`
}

// Error implements the error interface
func (e *APIError) Error() string {
	if e.Details != "" {
		return fmt.Sprintf("%s: %s (%s)", e.Code, e.Message, e.Details)
	}
	return fmt.Sprintf("%s: %s", e.Code, e.Message)
}

// NewAPIError creates a new API error
func NewAPIError(code ErrorCode, message string, details ...string) *APIError {
	apiErr := &APIError{
		Code:       code,
		Message:    message,
		HTTPStatus: getHTTPStatus(code),
	}
	
	if len(details) > 0 {
		apiErr.Details = details[0]
	}
	
	return apiErr
}

// WrapAPIError wraps an existing error with API error context
func WrapAPIError(err error, code ErrorCode, message string) *APIError {
	apiErr := NewAPIError(code, message)
	if err != nil {
		apiErr.Details = err.Error()
	}
	return apiErr
}

// getHTTPStatus returns the appropriate HTTP status code for an error code
func getHTTPStatus(code ErrorCode) int {
	switch code {
	case ErrCodeNotFound, ErrCodeNodeNotFound, ErrCodeContainerNotFound, 
		 ErrCodeNetworkRouteNotFound, ErrCodeProposalNotFound, ErrCodeUSBDeviceNotFound:
		return http.StatusNotFound
		
	case ErrCodeInvalidInput, ErrCodeNodeInvalidStatus, ErrCodeContainerInvalidStatus:
		return http.StatusBadRequest
		
	case ErrCodeUnauthorized:
		return http.StatusUnauthorized
		
	case ErrCodeForbidden:
		return http.StatusForbidden
		
	case ErrCodeConflict, ErrCodeNodeAlreadyExists, ErrCodeContainerAlreadyExists, 
		 ErrCodeNetworkRouteExists, ErrCodeVoteAlreadyCast, ErrCodeUSBDeviceInUse:
		return http.StatusConflict
		
	case ErrCodeTimeout:
		return http.StatusRequestTimeout
		
	case ErrCodeInsufficientCredits, ErrCodeTransactionFailed, ErrCodeUSBFormatFailed,
		 ErrCodeNodeCreationFailed, ErrCodeContainerDeployFailed:
		return http.StatusUnprocessableEntity
		
	case ErrCodeNetworkMeshDisabled:
		return http.StatusServiceUnavailable
		
	default:
		return http.StatusInternalServerError
	}
}

// Predefined errors for common cases
var (
	// General errors
	ErrInternal     = NewAPIError(ErrCodeInternal, "An internal error occurred")
	ErrNotFound     = NewAPIError(ErrCodeNotFound, "Resource not found")
	ErrInvalidInput = NewAPIError(ErrCodeInvalidInput, "Invalid input provided")
	ErrUnauthorized = NewAPIError(ErrCodeUnauthorized, "Unauthorized access")
	ErrForbidden    = NewAPIError(ErrCodeForbidden, "Access forbidden")
	ErrConflict     = NewAPIError(ErrCodeConflict, "Resource conflict")
	ErrTimeout      = NewAPIError(ErrCodeTimeout, "Operation timed out")

	// Node errors
	ErrNodeNotFound      = NewAPIError(ErrCodeNodeNotFound, "Node not found")
	ErrNodeAlreadyExists = NewAPIError(ErrCodeNodeAlreadyExists, "Node already exists")
	ErrNodeInvalidStatus = NewAPIError(ErrCodeNodeInvalidStatus, "Invalid node status")
	ErrNodeCreationFailed = NewAPIError(ErrCodeNodeCreationFailed, "Failed to create node")

	// Container errors
	ErrContainerNotFound      = NewAPIError(ErrCodeContainerNotFound, "Container not found")
	ErrContainerAlreadyExists = NewAPIError(ErrCodeContainerAlreadyExists, "Container already exists")
	ErrContainerInvalidStatus = NewAPIError(ErrCodeContainerInvalidStatus, "Invalid container status")
	ErrContainerDeployFailed  = NewAPIError(ErrCodeContainerDeployFailed, "Failed to deploy container")

	// Network errors
	ErrNetworkRouteNotFound = NewAPIError(ErrCodeNetworkRouteNotFound, "Network route not found")
	ErrNetworkRouteExists   = NewAPIError(ErrCodeNetworkRouteExists, "Network route already exists")
	ErrNetworkMeshDisabled  = NewAPIError(ErrCodeNetworkMeshDisabled, "Service mesh is disabled")

	// Cooperative errors
	ErrInsufficientCredits = NewAPIError(ErrCodeInsufficientCredits, "Insufficient credits")
	ErrTransactionFailed   = NewAPIError(ErrCodeTransactionFailed, "Transaction failed")
	ErrProposalNotFound    = NewAPIError(ErrCodeProposalNotFound, "Proposal not found")
	ErrVoteAlreadyCast     = NewAPIError(ErrCodeVoteAlreadyCast, "Vote already cast")

	// USB errors
	ErrUSBDeviceNotFound = NewAPIError(ErrCodeUSBDeviceNotFound, "USB device not found")
	ErrUSBDeviceInUse    = NewAPIError(ErrCodeUSBDeviceInUse, "USB device is in use")
	ErrUSBFormatFailed   = NewAPIError(ErrCodeUSBFormatFailed, "Failed to format USB device")
)

// IsAPIError checks if an error is an API error
func IsAPIError(err error) bool {
	_, ok := err.(*APIError)
	return ok
}

// GetAPIError extracts API error from an error, returns nil if not an API error
func GetAPIError(err error) *APIError {
	if apiErr, ok := err.(*APIError); ok {
		return apiErr
	}
	return nil
}

// GetErrorCode extracts the error code from an error
func GetErrorCode(err error) ErrorCode {
	if apiErr := GetAPIError(err); apiErr != nil {
		return apiErr.Code
	}
	return ErrCodeInternal
}

// GetHTTPStatus extracts the HTTP status code from an error
func GetHTTPStatus(err error) int {
	if apiErr := GetAPIError(err); apiErr != nil {
		return apiErr.HTTPStatus
	}
	return http.StatusInternalServerError
}
