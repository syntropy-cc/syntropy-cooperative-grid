package models

import (
	"time"
)

// Node represents a Syntropy node in the cooperative grid
type Node struct {
	ID           string                 `json:"id" gorm:"primaryKey;type:uuid;default:gen_random_uuid()"`
	Name         string                 `json:"name" gorm:"uniqueIndex;not null"`
	Description  string                 `json:"description"`
	Status       string                 `json:"status" gorm:"not null;default:'creating'"`
	HardwareInfo map[string]interface{} `json:"hardware_info" gorm:"type:jsonb"`
	NetworkConfig map[string]interface{} `json:"network_config" gorm:"type:jsonb"`
	CreatedAt    time.Time              `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt    time.Time              `json:"updated_at" gorm:"autoUpdateTime"`
}

// NodeStatus represents the possible statuses of a node
type NodeStatus string

const (
	NodeStatusCreating   NodeStatus = "creating"
	NodeStatusRunning    NodeStatus = "running"
	NodeStatusStopped    NodeStatus = "stopped"
	NodeStatusError      NodeStatus = "error"
	NodeStatusRestarting NodeStatus = "restarting"
	NodeStatusUpdating   NodeStatus = "updating"
)

// IsValid checks if the node status is valid
func (s NodeStatus) IsValid() bool {
	switch s {
	case NodeStatusCreating, NodeStatusRunning, NodeStatusStopped, 
		 NodeStatusError, NodeStatusRestarting, NodeStatusUpdating:
		return true
	default:
		return false
	}
}

// HardwareInfo represents hardware information for a node
type HardwareInfo struct {
	USBDevice  string `json:"usb_device,omitempty"`
	AutoDetect bool   `json:"auto_detect"`
	CPU        string `json:"cpu,omitempty"`
	Memory     string `json:"memory,omitempty"`
	Storage    string `json:"storage,omitempty"`
	Network    string `json:"network,omitempty"`
}

// NetworkConfig represents network configuration for a node
type NetworkConfig struct {
	IPAddress    string   `json:"ip_address,omitempty"`
	Port         int      `json:"port,omitempty"`
	Protocols    []string `json:"protocols,omitempty"`
	Firewall     bool     `json:"firewall"`
	LoadBalancer bool     `json:"load_balancer"`
	MeshEnabled  bool     `json:"mesh_enabled"`
}

// TableName returns the table name for the Node model
func (Node) TableName() string {
	return "nodes"
}
