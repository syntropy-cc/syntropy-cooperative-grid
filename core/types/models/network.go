package models

import (
	"time"
)

// NetworkRoute represents a network route between nodes
type NetworkRoute struct {
	ID               string                 `json:"id" gorm:"primaryKey;type:uuid;default:gen_random_uuid()"`
	SourceNodeID     string                 `json:"source_node_id" gorm:"not null;index"`
	DestinationNodeID string                 `json:"destination_node_id" gorm:"not null;index"`
	Config           map[string]interface{} `json:"config" gorm:"type:jsonb"`
	Status           string                 `json:"status" gorm:"not null;default:'active'"`
	CreatedAt        time.Time              `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt        time.Time              `json:"updated_at" gorm:"autoUpdateTime"`
	
	// Relationships
	SourceNode      *Node `json:"source_node,omitempty" gorm:"foreignKey:SourceNodeID"`
	DestinationNode *Node `json:"destination_node,omitempty" gorm:"foreignKey:DestinationNodeID"`
}

// NetworkRouteStatus represents the possible statuses of a network route
type NetworkRouteStatus string

const (
	NetworkRouteStatusActive   NetworkRouteStatus = "active"
	NetworkRouteStatusInactive NetworkRouteStatus = "inactive"
	NetworkRouteStatusError    NetworkRouteStatus = "error"
)

// IsValid checks if the network route status is valid
func (s NetworkRouteStatus) IsValid() bool {
	switch s {
	case NetworkRouteStatusActive, NetworkRouteStatusInactive, NetworkRouteStatusError:
		return true
	default:
		return false
	}
}

// NetworkRouteConfig represents configuration for a network route
type NetworkRouteConfig struct {
	Priority    int      `json:"priority,omitempty"`
	Protocol    string   `json:"protocol,omitempty"`
	Ports       []int    `json:"ports,omitempty"`
	Encryption  bool     `json:"encryption"`
	Compression bool     `json:"compression"`
	Tags        []string `json:"tags,omitempty"`
}

// ServiceMesh represents service mesh configuration
type ServiceMesh struct {
	ID        string                 `json:"id" gorm:"primaryKey;type:uuid;default:gen_random_uuid()"`
	Name      string                 `json:"name" gorm:"not null"`
	Config    map[string]interface{} `json:"config" gorm:"type:jsonb"`
	Status    string                 `json:"status" gorm:"not null;default:'disabled'"`
	CreatedAt time.Time              `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt time.Time              `json:"updated_at" gorm:"autoUpdateTime"`
}

// ServiceMeshStatus represents the possible statuses of a service mesh
type ServiceMeshStatus string

const (
	ServiceMeshStatusEnabled  ServiceMeshStatus = "enabled"
	ServiceMeshStatusDisabled ServiceMeshStatus = "disabled"
	ServiceMeshStatusError    ServiceMeshStatus = "error"
)

// IsValid checks if the service mesh status is valid
func (s ServiceMeshStatus) IsValid() bool {
	switch s {
	case ServiceMeshStatusEnabled, ServiceMeshStatusDisabled, ServiceMeshStatusError:
		return true
	default:
		return false
	}
}

// ServiceMeshConfig represents configuration for a service mesh
type ServiceMeshConfig struct {
	LoadBalancing string            `json:"load_balancing,omitempty"`
	Security      SecurityConfig    `json:"security,omitempty"`
	Monitoring    MonitoringConfig  `json:"monitoring,omitempty"`
	Policies      []MeshPolicy      `json:"policies,omitempty"`
}

// SecurityConfig represents security configuration for the service mesh
type SecurityConfig struct {
	MTLS           bool     `json:"mtls"`
	Encryption     bool     `json:"encryption"`
	Authentication bool     `json:"authentication"`
	Authorization  bool     `json:"authorization"`
	AllowedIPs     []string `json:"allowed_ips,omitempty"`
	BlockedIPs     []string `json:"blocked_ips,omitempty"`
}

// MonitoringConfig represents monitoring configuration for the service mesh
type MonitoringConfig struct {
	Metrics    bool `json:"metrics"`
	Logging    bool `json:"logging"`
	Tracing    bool `json:"tracing"`
	Alerting   bool `json:"alerting"`
}

// MeshPolicy represents a policy in the service mesh
type MeshPolicy struct {
	Name        string                 `json:"name"`
	Type        string                 `json:"type"`
	Config      map[string]interface{} `json:"config"`
	Priority    int                    `json:"priority"`
	Enabled     bool                   `json:"enabled"`
}

// TableName returns the table name for the NetworkRoute model
func (NetworkRoute) TableName() string {
	return "network_routes"
}

// TableName returns the table name for the ServiceMesh model
func (ServiceMesh) TableName() string {
	return "service_mesh"
}
