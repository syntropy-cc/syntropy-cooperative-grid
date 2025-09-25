// Package types provides common type definitions for the API
package types

import (
	"time"
)

// ErrorDetail represents detailed error information
type ErrorDetail struct {
	Code    string                 `json:"code"`    // Error code
	Message string                 `json:"message"` // Error message
	Details string                 `json:"details"` // Error details
	Field   string                 `json:"field"`   // Field that caused the error
	Stack   []string               `json:"stack"`   // Stack trace
	Context map[string]interface{} `json:"context"` // Error context
}

// Response represents a generic API response
type Response struct {
	Success bool          `json:"success"`         // Success status
	Data    interface{}   `json:"data,omitempty"`  // Response data
	Error   *ErrorDetail  `json:"error,omitempty"` // Error details
	Message string        `json:"message"`         // Response message
	Code    int           `json:"code"`            // Response code
	Meta    *ResponseMeta `json:"meta,omitempty"`  // Response metadata
}

// ResponseMeta represents response metadata
type ResponseMeta struct {
	Timestamp  time.Time              `json:"timestamp"`   // Response timestamp
	RequestID  string                 `json:"request_id"`  // Request identifier
	Version    string                 `json:"version"`     // API version
	Interface  string                 `json:"interface"`   // Interface type
	UserID     string                 `json:"user_id"`     // User identifier
	SessionID  string                 `json:"session_id"`  // Session identifier
	Duration   time.Duration          `json:"duration"`    // Request duration
	CustomData map[string]interface{} `json:"custom_data"` // Custom metadata
}

// HealthCheck represents a health check response
type HealthCheck struct {
	Status      string                     `json:"status"`      // Health status
	Timestamp   time.Time                  `json:"timestamp"`   // Check timestamp
	Version     string                     `json:"version"`     // Service version
	Uptime      time.Duration              `json:"uptime"`      // Service uptime
	Components  map[string]ComponentHealth `json:"components"`  // Component health
	Performance *PerformanceMetrics        `json:"performance"` // Performance metrics
}

// ComponentHealth represents health of a component
type ComponentHealth struct {
	Status    string                 `json:"status"`     // Component status
	Message   string                 `json:"message"`    // Health message
	LastCheck time.Time              `json:"last_check"` // Last check timestamp
	Details   map[string]interface{} `json:"details"`    // Component details
}

// PerformanceMetrics represents performance metrics
type PerformanceMetrics struct {
	CPUUsage     float64 `json:"cpu_usage"`     // CPU usage percentage
	MemoryUsage  float64 `json:"memory_usage"`  // Memory usage percentage
	DiskUsage    float64 `json:"disk_usage"`    // Disk usage percentage
	NetworkIO    float64 `json:"network_io"`    // Network I/O rate
	RequestRate  float64 `json:"request_rate"`  // Request rate per second
	ResponseTime float64 `json:"response_time"` // Average response time
	ErrorRate    float64 `json:"error_rate"`    // Error rate percentage
}

// LogEntry represents a log entry
type LogEntry struct {
	Timestamp time.Time              `json:"timestamp"`  // Log timestamp
	Level     string                 `json:"level"`      // Log level
	Message   string                 `json:"message"`    // Log message
	Source    string                 `json:"source"`     // Log source
	RequestID string                 `json:"request_id"` // Request identifier
	UserID    string                 `json:"user_id"`    // User identifier
	SessionID string                 `json:"session_id"` // Session identifier
	Interface string                 `json:"interface"`  // Interface type
	Context   map[string]interface{} `json:"context"`    // Log context
	Stack     []string               `json:"stack"`      // Stack trace
}

// Metric represents a metric
type Metric struct {
	Name        string                 `json:"name"`        // Metric name
	Value       float64                `json:"value"`       // Metric value
	Unit        string                 `json:"unit"`        // Metric unit
	Type        string                 `json:"type"`        // Metric type (counter, gauge, histogram)
	Timestamp   time.Time              `json:"timestamp"`   // Metric timestamp
	Labels      map[string]string      `json:"labels"`      // Metric labels
	Description string                 `json:"description"` // Metric description
	Metadata    map[string]interface{} `json:"metadata"`    // Metric metadata
}

// Event represents an event
type Event struct {
	ID        string                 `json:"id"`         // Event ID
	Type      string                 `json:"type"`       // Event type
	Source    string                 `json:"source"`     // Event source
	Timestamp time.Time              `json:"timestamp"`  // Event timestamp
	Data      map[string]interface{} `json:"data"`       // Event data
	UserID    string                 `json:"user_id"`    // User identifier
	SessionID string                 `json:"session_id"` // Session identifier
	Interface string                 `json:"interface"`  // Interface type
	Severity  string                 `json:"severity"`   // Event severity
	Message   string                 `json:"message"`    // Event message
	Metadata  map[string]interface{} `json:"metadata"`   // Event metadata
}

// Notification represents a notification
type Notification struct {
	ID        string                 `json:"id"`        // Notification ID
	Type      string                 `json:"type"`      // Notification type
	Title     string                 `json:"title"`     // Notification title
	Message   string                 `json:"message"`   // Notification message
	Timestamp time.Time              `json:"timestamp"` // Notification timestamp
	Read      bool                   `json:"read"`      // Whether notification is read
	UserID    string                 `json:"user_id"`   // User identifier
	Interface string                 `json:"interface"` // Interface type
	Priority  string                 `json:"priority"`  // Notification priority
	Actions   []NotificationAction   `json:"actions"`   // Available actions
	Metadata  map[string]interface{} `json:"metadata"`  // Notification metadata
}

// NotificationAction represents a notification action
type NotificationAction struct {
	ID      string `json:"id"`      // Action ID
	Label   string `json:"label"`   // Action label
	Type    string `json:"type"`    // Action type
	URL     string `json:"url"`     // Action URL
	Command string `json:"command"` // Action command
}

// User represents a user
type User struct {
	ID          string                 `json:"id"`          // User ID
	Username    string                 `json:"username"`    // Username
	Email       string                 `json:"email"`       // Email address
	FullName    string                 `json:"full_name"`   // Full name
	Role        string                 `json:"role"`        // User role
	Permissions []string               `json:"permissions"` // User permissions
	Preferences map[string]interface{} `json:"preferences"` // User preferences
	CreatedAt   time.Time              `json:"created_at"`  // Creation timestamp
	UpdatedAt   time.Time              `json:"updated_at"`  // Last update timestamp
	LastLogin   time.Time              `json:"last_login"`  // Last login timestamp
	Status      string                 `json:"status"`      // User status
	Metadata    map[string]interface{} `json:"metadata"`    // User metadata
}

// Session represents a user session
type Session struct {
	ID           string                 `json:"id"`            // Session ID
	UserID       string                 `json:"user_id"`       // User ID
	Interface    string                 `json:"interface"`     // Interface type
	CreatedAt    time.Time              `json:"created_at"`    // Session creation timestamp
	ExpiresAt    time.Time              `json:"expires_at"`    // Session expiration timestamp
	LastActivity time.Time              `json:"last_activity"` // Last activity timestamp
	IPAddress    string                 `json:"ip_address"`    // IP address
	UserAgent    string                 `json:"user_agent"`    // User agent
	Active       bool                   `json:"active"`        // Whether session is active
	Metadata     map[string]interface{} `json:"metadata"`      // Session metadata
}

// AuditLog represents an audit log entry
type AuditLog struct {
	ID        string                 `json:"id"`         // Audit log ID
	Action    string                 `json:"action"`     // Action performed
	Resource  string                 `json:"resource"`   // Resource affected
	UserID    string                 `json:"user_id"`    // User ID
	SessionID string                 `json:"session_id"` // Session ID
	Interface string                 `json:"interface"`  // Interface type
	Timestamp time.Time              `json:"timestamp"`  // Action timestamp
	IPAddress string                 `json:"ip_address"` // IP address
	UserAgent string                 `json:"user_agent"` // User agent
	Result    string                 `json:"result"`     // Action result
	Details   map[string]interface{} `json:"details"`    // Action details
	Metadata  map[string]interface{} `json:"metadata"`   // Audit metadata
}

// Status constants
const (
	StatusSuccess  = "success"
	StatusError    = "error"
	StatusWarning  = "warning"
	StatusInfo     = "info"
	StatusPending  = "pending"
	StatusRunning  = "running"
	StatusComplete = "complete"
	StatusFailed   = "failed"
	StatusActive   = "active"
	StatusInactive = "inactive"
	StatusDisabled = "disabled"
	StatusEnabled  = "enabled"
)

// Interface types
const (
	InterfaceCLI     = "cli"
	InterfaceWeb     = "web"
	InterfaceDesktop = "desktop"
	InterfaceMobile  = "mobile"
)

// Log levels
const (
	LogLevelDebug   = "debug"
	LogLevelInfo    = "info"
	LogLevelWarning = "warning"
	LogLevelError   = "error"
	LogLevelFatal   = "fatal"
)

// Priority levels
const (
	PriorityLow      = "low"
	PriorityMedium   = "medium"
	PriorityHigh     = "high"
	PriorityCritical = "critical"
)

// Severity levels
const (
	SeverityInfo     = "info"
	SeverityWarning  = "warning"
	SeverityError    = "error"
	SeverityCritical = "critical"
)

// User roles
const (
	RoleAdmin    = "admin"
	RoleUser     = "user"
	RoleGuest    = "guest"
	RoleOperator = "operator"
	RoleViewer   = "viewer"
)

// User status
const (
	UserStatusActive    = "active"
	UserStatusInactive  = "inactive"
	UserStatusPending   = "pending"
	UserStatusSuspended = "suspended"
	UserStatusDeleted   = "deleted"
)
