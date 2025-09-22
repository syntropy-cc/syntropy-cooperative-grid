package models

import (
	"time"
)

// Container represents a containerized application in the Syntropy grid
type Container struct {
	ID        string                 `json:"id" gorm:"primaryKey;type:uuid;default:gen_random_uuid()"`
	NodeID    string                 `json:"node_id" gorm:"not null;index"`
	Name      string                 `json:"name" gorm:"not null"`
	Image     string                 `json:"image" gorm:"not null"`
	Status    string                 `json:"status" gorm:"not null;default:'creating'"`
	Config    map[string]interface{} `json:"config" gorm:"type:jsonb"`
	CreatedAt time.Time              `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt time.Time              `json:"updated_at" gorm:"autoUpdateTime"`
	
	// Relationships
	Node *Node `json:"node,omitempty" gorm:"foreignKey:NodeID"`
}

// ContainerStatus represents the possible statuses of a container
type ContainerStatus string

const (
	ContainerStatusCreating   ContainerStatus = "creating"
	ContainerStatusRunning    ContainerStatus = "running"
	ContainerStatusStopped    ContainerStatus = "stopped"
	ContainerStatusError      ContainerStatus = "error"
	ContainerStatusRestarting ContainerStatus = "restarting"
	ContainerStatusUpdating   ContainerStatus = "updating"
)

// IsValid checks if the container status is valid
func (s ContainerStatus) IsValid() bool {
	switch s {
	case ContainerStatusCreating, ContainerStatusRunning, ContainerStatusStopped,
		 ContainerStatusError, ContainerStatusRestarting, ContainerStatusUpdating:
		return true
	default:
		return false
	}
}

// ContainerConfig represents configuration for a container
type ContainerConfig struct {
	Ports     []PortMapping     `json:"ports,omitempty"`
	EnvVars   map[string]string `json:"env_vars,omitempty"`
	Volumes   []VolumeMapping   `json:"volumes,omitempty"`
	Resources ResourceLimits    `json:"resources,omitempty"`
	Replicas  int               `json:"replicas,omitempty"`
}

// PortMapping represents a port mapping for a container
type PortMapping struct {
	HostPort      int    `json:"host_port"`
	ContainerPort int    `json:"container_port"`
	Protocol      string `json:"protocol"`
}

// VolumeMapping represents a volume mapping for a container
type VolumeMapping struct {
	HostPath      string `json:"host_path"`
	ContainerPath string `json:"container_path"`
	ReadOnly      bool   `json:"read_only"`
}

// ResourceLimits represents resource limits for a container
type ResourceLimits struct {
	CPU    string `json:"cpu,omitempty"`
	Memory string `json:"memory,omitempty"`
	Disk   string `json:"disk,omitempty"`
}

// TableName returns the table name for the Container model
func (Container) TableName() string {
	return "containers"
}
