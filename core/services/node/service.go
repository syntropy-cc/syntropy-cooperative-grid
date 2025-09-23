package node

import (
	"context"
	"fmt"

	"github.com/syntropy-cc/cooperative-grid/core/types/models"
)

// Service handles node management operations
type Service struct {
	repo Repository
	log  Logger
}

// Repository defines the interface for node data access
type Repository interface {
	Create(ctx context.Context, node *models.Node) error
	GetByID(ctx context.Context, id string) (*models.Node, error)
	GetByName(ctx context.Context, name string) (*models.Node, error)
	List(ctx context.Context, filter *Filter) ([]*models.Node, error)
	Update(ctx context.Context, node *models.Node) error
	Delete(ctx context.Context, id string) error
}

// Logger defines the interface for logging
type Logger interface {
	Info(args ...interface{})
	Error(args ...interface{})
	Debug(args ...interface{})
	WithFields(fields map[string]interface{}) Logger
}

// Filter defines filtering options for node queries
type Filter struct {
	Status string
	Name   string
	Limit  int
	Offset int
}

// CreateNodeRequest represents a request to create a new node
type CreateNodeRequest struct {
	Name        string `json:"name"`
	Description string `json:"description,omitempty"`
	USBDevice   string `json:"usb_device,omitempty"`
	AutoDetect  bool   `json:"auto_detect"`
}

// NewService creates a new node service
func NewService(repo Repository, log Logger) *Service {
	return &Service{
		repo: repo,
		log:  log,
	}
}

// CreateNode creates a new node
func (s *Service) CreateNode(ctx context.Context, req *CreateNodeRequest) (*models.Node, error) {
	s.log.WithFields(map[string]interface{}{
		"name":       req.Name,
		"usb_device": req.USBDevice,
		"auto_detect": req.AutoDetect,
	}).Info("Creating new node")

	// Validate request
	if err := s.validateCreateRequest(req); err != nil {
		return nil, fmt.Errorf("invalid request: %w", err)
	}

	// Check if node with same name already exists
	existing, err := s.repo.GetByName(ctx, req.Name)
	if err == nil && existing != nil {
		return nil, fmt.Errorf("node with name '%s' already exists", req.Name)
	}

	// Create node model
	node := &models.Node{
		Name:        req.Name,
		Description: req.Description,
		Status:      "creating",
		HardwareInfo: map[string]interface{}{
			"usb_device": req.USBDevice,
			"auto_detect": req.AutoDetect,
		},
	}

	// Save to repository
	if err := s.repo.Create(ctx, node); err != nil {
		return nil, fmt.Errorf("failed to create node: %w", err)
	}

	s.log.WithFields(map[string]interface{}{
		"node_id": node.ID,
		"name":    node.Name,
	}).Info("Node created successfully")

	return node, nil
}

// GetNode retrieves a node by ID
func (s *Service) GetNode(ctx context.Context, id string) (*models.Node, error) {
	s.log.WithFields(map[string]interface{}{
		"node_id": id,
	}).Debug("Getting node")

	node, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get node: %w", err)
	}

	return node, nil
}

// ListNodes retrieves a list of nodes with optional filtering
func (s *Service) ListNodes(ctx context.Context, filter *Filter) ([]*models.Node, error) {
	s.log.WithFields(map[string]interface{}{
		"filter": filter,
	}).Debug("Listing nodes")

	nodes, err := s.repo.List(ctx, filter)
	if err != nil {
		return nil, fmt.Errorf("failed to list nodes: %w", err)
	}

	return nodes, nil
}

// UpdateNode updates an existing node
func (s *Service) UpdateNode(ctx context.Context, id string, updates map[string]interface{}) (*models.Node, error) {
	s.log.WithFields(map[string]interface{}{
		"node_id": id,
		"updates": updates,
	}).Info("Updating node")

	// Get existing node
	node, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get node: %w", err)
	}

	// Apply updates
	if name, ok := updates["name"].(string); ok {
		node.Name = name
	}
	if description, ok := updates["description"].(string); ok {
		node.Description = description
	}
	if status, ok := updates["status"].(string); ok {
		node.Status = status
	}

	// Save updated node
	if err := s.repo.Update(ctx, node); err != nil {
		return nil, fmt.Errorf("failed to update node: %w", err)
	}

	s.log.WithFields(map[string]interface{}{
		"node_id": node.ID,
		"name":    node.Name,
	}).Info("Node updated successfully")

	return node, nil
}

// DeleteNode deletes a node
func (s *Service) DeleteNode(ctx context.Context, id string) error {
	s.log.WithFields(map[string]interface{}{
		"node_id": id,
	}).Info("Deleting node")

	// Check if node exists
	_, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return fmt.Errorf("failed to get node: %w", err)
	}

	// Delete node
	if err := s.repo.Delete(ctx, id); err != nil {
		return fmt.Errorf("failed to delete node: %w", err)
	}

	s.log.WithFields(map[string]interface{}{
		"node_id": id,
	}).Info("Node deleted successfully")

	return nil
}

// RestartNode restarts a node and its services
func (s *Service) RestartNode(ctx context.Context, id string) error {
	s.log.WithFields(map[string]interface{}{
		"node_id": id,
	}).Info("Restarting node")

	// Get node
	node, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return fmt.Errorf("failed to get node: %w", err)
	}

	// Update status to restarting
	node.Status = "restarting"
	if err := s.repo.Update(ctx, node); err != nil {
		return fmt.Errorf("failed to update node status: %w", err)
	}

	// TODO: Implement actual restart logic
	// This would involve:
	// 1. Stopping all services on the node
	// 2. Restarting the node
	// 3. Starting services again
	// 4. Updating status to "running"

	// For now, just update status to running
	node.Status = "running"
	if err := s.repo.Update(ctx, node); err != nil {
		return fmt.Errorf("failed to update node status: %w", err)
	}

	s.log.WithFields(map[string]interface{}{
		"node_id": id,
	}).Info("Node restarted successfully")

	return nil
}

// validateCreateRequest validates a create node request
func (s *Service) validateCreateRequest(req *CreateNodeRequest) error {
	if req.Name == "" {
		return fmt.Errorf("node name is required")
	}

	if len(req.Name) < 3 {
		return fmt.Errorf("node name must be at least 3 characters long")
	}

	if len(req.Name) > 50 {
		return fmt.Errorf("node name must be less than 50 characters")
	}

	// Validate name format (alphanumeric, hyphens, underscores)
	for _, char := range req.Name {
		if !((char >= 'a' && char <= 'z') || 
			 (char >= 'A' && char <= 'Z') || 
			 (char >= '0' && char <= '9') || 
			 char == '-' || char == '_') {
			return fmt.Errorf("node name can only contain alphanumeric characters, hyphens, and underscores")
		}
	}

	return nil
}


