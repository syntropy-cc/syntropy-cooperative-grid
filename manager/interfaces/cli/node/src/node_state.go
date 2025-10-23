package node

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"time"

	"node-component/src/internal/types"
)

// NodeStateManager gerencia estado de nós com sincronização arquivo↔memória
// Fonte de verdade: arquivos em ~/.syntropy/nodes/
// Cache: nodeStates em memória
// Sincronização: Automática a cada 10 segundos (Estratégia A)
type NodeStateManager struct {
	// Cache em memória para performance
	nodeStates  map[string]*types.NodeStatus
	nodeConfigs map[string]*types.NodeConfig
	mutex       sync.RWMutex

	// Persistência (fonte de verdade)
	nodesDir string // ~/.syntropy/nodes/
	logger   types.Logger

	// Sincronização automática
	syncInterval time.Duration
	syncTicker   *time.Ticker
	stopChan     chan struct{}
	isRunning    bool
}

// NewNodeStateManager cria gerenciador com caminho correto
func NewNodeStateManager(logger types.Logger, nodesDir string) *NodeStateManager {
	// Usar nodesDir fornecido ou derivar de getSyntropyDir
	if nodesDir == "" {
		homeDir, err := os.UserHomeDir()
		if err != nil {
			logger.Warn("Failed to get home directory, using /tmp", "error", err)
			homeDir = "/tmp"
		}
		nodesDir = filepath.Join(homeDir, ".syntropy", "nodes")
	}

	return &NodeStateManager{
		nodeStates:   make(map[string]*types.NodeStatus),
		nodeConfigs:  make(map[string]*types.NodeConfig),
		nodesDir:     nodesDir,
		logger:       logger,
		syncInterval: 10 * time.Second,
		stopChan:     make(chan struct{}),
		isRunning:    false,
	}
}

// LoadState carrega estado inicial do arquivo (fonte de verdade)
func (nsm *NodeStateManager) LoadState() error {
	nsm.mutex.Lock()
	defer nsm.mutex.Unlock()

	// Criar diretório se não existir (primeiro acesso)
	if err := os.MkdirAll(nsm.nodesDir, 0700); err != nil {
		if !os.IsExist(err) {
			return fmt.Errorf("failed to create nodes directory: %w", err)
		}
	}

	// Ler todos os diretórios de nó
	entries, err := os.ReadDir(nsm.nodesDir)
	if err != nil {
		if os.IsNotExist(err) {
			return nil // OK - não há nós ainda
		}
		return fmt.Errorf("failed to read nodes directory: %w", err)
	}

	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}

		nodeID := entry.Name()

		// Carregar status.json
		status, err := nsm.loadNodeStatus(nodeID)
		if err != nil {
			nsm.logger.Warn("Failed to load node status", "node_id", nodeID, "error", err)
			continue
		}

		// Carregar config.json
		config, err := nsm.loadNodeConfig(nodeID)
		if err != nil {
			nsm.logger.Warn("Failed to load node config", "node_id", nodeID, "error", err)
			continue
		}

		nsm.nodeStates[nodeID] = status
		if config != nil {
			nsm.nodeConfigs[nodeID] = config
		}
	}

	return nil
}

// CreateNode cria novo nó e persiste IMEDIATAMENTE em arquivo
func (nsm *NodeStateManager) CreateNode(nodeID string, config *types.NodeConfig) error {
	nsm.mutex.Lock()
	defer nsm.mutex.Unlock()

	// Criar status inicial
	status := &types.NodeStatus{
		NodeID:    nodeID,
		Status:    "pending",
		CreatedAt: time.Now(),
	}

	// Persistir em arquivo PRIMEIRO (não em memória primeiro!)
	if err := nsm.saveNodeStatus(nodeID, status); err != nil {
		return fmt.Errorf("failed to persist node status: %w", err)
	}

	if err := nsm.saveNodeConfig(nodeID, config); err != nil {
		// Limpar arquivo se salvar config falhar
		nsm.logger.Warn("Failed to save config, trying to clean up status file", nil)
		os.RemoveAll(filepath.Join(nsm.nodesDir, nodeID))
		return fmt.Errorf("failed to persist node config: %w", err)
	}

	// Depois atualizar cache
	nsm.nodeStates[nodeID] = status
	nsm.nodeConfigs[nodeID] = config

	return nil
}

// GetNode retorna status de um nó
func (nsm *NodeStateManager) GetNode(nodeID string) (*types.NodeStatus, error) {
	nsm.mutex.RLock()
	defer nsm.mutex.RUnlock()

	status, exists := nsm.nodeStates[nodeID]
	if !exists {
		return nil, fmt.Errorf("node not found: %s", nodeID)
	}

	return status, nil
}

// UpdateNodeStatus atualiza status progressivamente (hardware após conexão)
func (nsm *NodeStateManager) UpdateNodeStatus(nodeID string, status *types.NodeStatus) error {
	nsm.mutex.Lock()
	defer nsm.mutex.Unlock()

	// Garantir que não perde informações anteriores
	if existing, exists := nsm.nodeStates[nodeID]; exists {
		if status.CreatedAt.IsZero() {
			status.CreatedAt = existing.CreatedAt
		}
		if status.NodeID == "" {
			status.NodeID = existing.NodeID
		}
	}

	// Persistir em arquivo (fonte de verdade)
	if err := nsm.saveNodeStatus(nodeID, status); err != nil {
		return fmt.Errorf("failed to update node status in file: %w", err)
	}

	// Atualizar cache
	nsm.nodeStates[nodeID] = status

	return nil
}

// TransitionToActive marca nó como ativo e coleta hardware inicial
func (nsm *NodeStateManager) TransitionToActive(nodeID string, ipAddress string, hardware *types.HardwareInfo) error {
	nsm.mutex.Lock()
	defer nsm.mutex.Unlock()

	status, exists := nsm.nodeStates[nodeID]
	if !exists {
		return fmt.Errorf("node not found: %s", nodeID)
	}

	status.Status = "active"
	status.RegisteredAt = time.Now()
	status.IPAddress = ipAddress
	status.LastHeartbeat = time.Now()

	if hardware != nil {
		status.Hardware = hardware
	}

	if err := nsm.saveNodeStatus(nodeID, status); err != nil {
		return fmt.Errorf("failed to persist active state: %w", err)
	}

	return nil
}

// StartAutoSync inicia sincronização periódica com arquivo
func (nsm *NodeStateManager) StartAutoSync() {
	nsm.mutex.Lock()
	defer nsm.mutex.Unlock()

	if nsm.isRunning {
		return // Já está rodando
	}

	nsm.isRunning = true
	nsm.syncTicker = time.NewTicker(nsm.syncInterval)

	go func() {
		for {
			select {
			case <-nsm.syncTicker.C:
				// Sincronizar: detectar mudanças externas
				if err := nsm.detectExternalChanges(); err != nil {
					nsm.logger.Warn("Auto-sync detect changes failed", "error", err)
				}

			case <-nsm.stopChan:
				nsm.syncTicker.Stop()
				return
			}
		}
	}()
}

// StopAutoSync para a sincronização automática
func (nsm *NodeStateManager) StopAutoSync() {
	nsm.mutex.Lock()
	defer nsm.mutex.Unlock()

	if !nsm.isRunning {
		return
	}

	nsm.isRunning = false
	close(nsm.stopChan)
}

// detectExternalChanges detecta nós adicionados externamente (arquivo)
func (nsm *NodeStateManager) detectExternalChanges() error {
	nsm.mutex.Lock()
	defer nsm.mutex.Unlock()

	entries, err := os.ReadDir(nsm.nodesDir)
	if err != nil {
		return fmt.Errorf("failed to read nodes directory: %w", err)
	}

	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}

		nodeID := entry.Name()

		// Se nó em arquivo mas não em cache, recarregar
		if _, exists := nsm.nodeStates[nodeID]; !exists {
			status, err := nsm.loadNodeStatus(nodeID)
			if err != nil {
				nsm.logger.Warn("Failed to load external node", "node_id", nodeID, "error", err)
				continue
			}

			nsm.nodeStates[nodeID] = status
			nsm.logger.Info("Loaded externally added node from file", "node_id", nodeID)
		}
	}

	return nil
}

// UpdateNodeHardware atualiza informações de hardware (não destrutivo)
func (nsm *NodeStateManager) UpdateNodeHardware(nodeID string, hardware *types.HardwareInfo) error {
	nsm.mutex.Lock()
	defer nsm.mutex.Unlock()

	status, exists := nsm.nodeStates[nodeID]
	if !exists {
		return fmt.Errorf("node not found: %s", nodeID)
	}

	// Mesclar hardware (se já existe, atualizar campos)
	if status.Hardware == nil {
		status.Hardware = hardware
	} else {
		// Atualizar apenas campos preenchidos (não sobrescrever com nil)
		if hardware.CPUCores > 0 {
			status.Hardware.CPUCores = hardware.CPUCores
		}
		if hardware.MemoryGB > 0 {
			status.Hardware.MemoryGB = hardware.MemoryGB
		}
		if hardware.DiskGB > 0 {
			status.Hardware.DiskGB = hardware.DiskGB
		}
		if hardware.Hostname != "" {
			status.Hardware.Hostname = hardware.Hostname
		}
		if hardware.IPAddress != "" {
			status.Hardware.IPAddress = hardware.IPAddress
		}
		if hardware.OSVersion != "" {
			status.Hardware.OSVersion = hardware.OSVersion
		}
		if hardware.KernelVersion != "" {
			status.Hardware.KernelVersion = hardware.KernelVersion
		}
	}

	// Persistir atualização em arquivo
	if err := nsm.saveNodeStatus(nodeID, status); err != nil {
		return fmt.Errorf("failed to persist hardware info: %w", err)
	}

	return nil
}

// Métodos auxiliares de persistência

// saveNodeStatus persiste status em arquivo
func (nsm *NodeStateManager) saveNodeStatus(nodeID string, status *types.NodeStatus) error {
	nodeDir := filepath.Join(nsm.nodesDir, nodeID)
	if err := os.MkdirAll(nodeDir, 0700); err != nil {
		return fmt.Errorf("failed to create node directory: %w", err)
	}

	statusPath := filepath.Join(nodeDir, "status.json")
	data, err := json.MarshalIndent(status, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal status: %w", err)
	}

	// Permissões restritivas (apenas proprietário)
	if err := os.WriteFile(statusPath, data, 0600); err != nil {
		return fmt.Errorf("failed to write status file: %w", err)
	}

	return nil
}

// loadNodeStatus carrega status do arquivo
func (nsm *NodeStateManager) loadNodeStatus(nodeID string) (*types.NodeStatus, error) {
	nodeDir := filepath.Join(nsm.nodesDir, nodeID)
	statusPath := filepath.Join(nodeDir, "status.json")

	data, err := os.ReadFile(statusPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read status file: %w", err)
	}

	var status types.NodeStatus
	if err := json.Unmarshal(data, &status); err != nil {
		return nil, fmt.Errorf("failed to unmarshal status: %w", err)
	}

	return &status, nil
}

// saveNodeConfig persiste configuração em arquivo
func (nsm *NodeStateManager) saveNodeConfig(nodeID string, config *types.NodeConfig) error {
	nodeDir := filepath.Join(nsm.nodesDir, nodeID)
	if err := os.MkdirAll(nodeDir, 0700); err != nil {
		return fmt.Errorf("failed to create node directory: %w", err)
	}

	configPath := filepath.Join(nodeDir, "config.json")
	data, err := json.MarshalIndent(config, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal config: %w", err)
	}

	// Permissões restritivas (apenas proprietário)
	if err := os.WriteFile(configPath, data, 0600); err != nil {
		return fmt.Errorf("failed to write config file: %w", err)
	}

	return nil
}

// loadNodeConfig carrega configuração do arquivo
func (nsm *NodeStateManager) loadNodeConfig(nodeID string) (*types.NodeConfig, error) {
	nodeDir := filepath.Join(nsm.nodesDir, nodeID)
	configPath := filepath.Join(nodeDir, "config.json")

	data, err := os.ReadFile(configPath)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil // Config pode não existir ainda
		}
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	var config types.NodeConfig
	if err := json.Unmarshal(data, &config); err != nil {
		return nil, fmt.Errorf("failed to unmarshal config: %w", err)
	}

	return &config, nil
}
