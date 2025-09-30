package setup

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"time"

	"setup-component/src/internal/types"
)

// StateManager implementa a interface StateManager
type StateManager struct {
	statePath string
	lockFile  string
	mutex     sync.RWMutex
	logger    *SetupLogger
}

// NewStateManager cria um novo gerenciador de estado
func NewStateManager(logger *SetupLogger) *StateManager {
	homeDir, _ := os.UserHomeDir()
	stateDir := filepath.Join(homeDir, ".syntropy", "state")
	os.MkdirAll(stateDir, 0755)

	return &StateManager{
		statePath: filepath.Join(stateDir, "setup_state.json"),
		lockFile:  filepath.Join(stateDir, "setup_state.lock"),
		logger:    logger,
	}
}

// LoadState carrega o estado atual
func (sm *StateManager) LoadState() (*types.SetupState, error) {
	sm.mutex.RLock()
	defer sm.mutex.RUnlock()

	sm.logger.LogDebug("Carregando estado do setup", map[string]interface{}{
		"state_path": sm.statePath,
	})

	// Verificar se o arquivo de estado existe
	if _, err := os.Stat(sm.statePath); os.IsNotExist(err) {
		sm.logger.LogDebug("Arquivo de estado não encontrado", map[string]interface{}{
			"state_path": sm.statePath,
		})
		return nil, types.ErrStateLoadError(fmt.Errorf("arquivo de estado não encontrado: %s", sm.statePath))
	}

	// Ler arquivo de estado
	data, err := os.ReadFile(sm.statePath)
	if err != nil {
		return nil, types.ErrStateLoadError(err)
	}

	// Deserializar estado
	var state types.SetupState
	if err := json.Unmarshal(data, &state); err != nil {
		return nil, types.ErrStateCorruptedError(err)
	}

	sm.logger.LogDebug("Estado carregado com sucesso", map[string]interface{}{
		"version": state.Version,
		"status":  state.Status,
	})

	return &state, nil
}

// SaveState salva o estado de forma atômica
func (sm *StateManager) SaveState(state *types.SetupState) error {
	sm.mutex.Lock()
	defer sm.mutex.Unlock()

	sm.logger.LogDebug("Salvando estado do setup", map[string]interface{}{
		"version": state.Version,
		"status":  state.Status,
	})

	// Atualizar timestamp
	state.UpdatedAt = time.Now()

	// Serializar estado
	data, err := json.MarshalIndent(state, "", "  ")
	if err != nil {
		return types.ErrStateSaveError(err)
	}

	// Criar arquivo temporário para operação atômica
	tempPath := sm.statePath + ".tmp"

	// Escrever para arquivo temporário
	if err := os.WriteFile(tempPath, data, 0644); err != nil {
		return types.ErrStateSaveError(err)
	}

	// Renomear arquivo temporário para arquivo final (operação atômica)
	if err := os.Rename(tempPath, sm.statePath); err != nil {
		// Limpar arquivo temporário em caso de erro
		os.Remove(tempPath)
		return types.ErrStateSaveError(err)
	}

	sm.logger.LogInfo("Estado salvo com sucesso", map[string]interface{}{
		"version": state.Version,
		"status":  state.Status,
	})

	return nil
}

// UpdateState atualiza o estado de forma atômica
func (sm *StateManager) UpdateState(update func(*types.SetupState) error) error {
	sm.mutex.Lock()
	defer sm.mutex.Unlock()

	sm.logger.LogDebug("Atualizando estado do setup", nil)

	// Carregar estado atual
	state, err := sm.loadStateUnsafe()
	if err != nil {
		return err
	}

	// Aplicar atualização
	if err := update(state); err != nil {
		return err
	}

	// Salvar estado atualizado
	return sm.saveStateUnsafe(state)
}

// BackupState cria um backup do estado
func (sm *StateManager) BackupState(name string) error {
	sm.mutex.RLock()
	defer sm.mutex.RUnlock()

	sm.logger.LogStep("backup_state_start", map[string]interface{}{
		"backup_name": name,
	})

	// Carregar estado atual
	state, err := sm.loadStateUnsafe()
	if err != nil {
		return err
	}

	// Criar diretório de backups
	backupDir := filepath.Join(filepath.Dir(sm.statePath), "backups")
	os.MkdirAll(backupDir, 0755)

	// Gerar nome do arquivo de backup
	timestamp := time.Now().Format("20060102_150405")
	backupPath := filepath.Join(backupDir, fmt.Sprintf("state_%s_%s.json", name, timestamp))

	// Serializar estado
	data, err := json.MarshalIndent(state, "", "  ")
	if err != nil {
		return types.ErrBackupFailedError(err)
	}

	// Escrever backup
	if err := os.WriteFile(backupPath, data, 0644); err != nil {
		return types.ErrBackupFailedError(err)
	}

	// Atualizar informações de backup no estado
	state.LastBackup = &types.BackupInfo{
		ID:        name,
		Path:      backupPath,
		Size:      int64(len(data)),
		CreatedAt: time.Now(),
	}

	sm.logger.LogStep("backup_state_completed", map[string]interface{}{
		"backup_path": backupPath,
		"backup_size": len(data),
	})

	return nil
}

// RestoreState restaura o estado de um backup
func (sm *StateManager) RestoreState(backupPath string) error {
	sm.mutex.Lock()
	defer sm.mutex.Unlock()

	sm.logger.LogStep("restore_state_start", map[string]interface{}{
		"backup_path": backupPath,
	})

	// Verificar se o arquivo de backup existe
	if _, err := os.Stat(backupPath); os.IsNotExist(err) {
		return types.ErrRestoreFailedError(fmt.Errorf("arquivo de backup não encontrado: %s", backupPath))
	}

	// Ler arquivo de backup
	data, err := os.ReadFile(backupPath)
	if err != nil {
		return types.ErrRestoreFailedError(err)
	}

	// Deserializar estado do backup
	var state types.SetupState
	if err := json.Unmarshal(data, &state); err != nil {
		return types.ErrRestoreFailedError(err)
	}

	// Salvar estado restaurado
	if err := sm.saveStateUnsafe(&state); err != nil {
		return err
	}

	sm.logger.LogStep("restore_state_completed", map[string]interface{}{
		"backup_path": backupPath,
		"version":     state.Version,
		"status":      state.Status,
	})

	return nil
}

// VerifyIntegrity verifica a integridade do estado
func (sm *StateManager) VerifyIntegrity() error {
	sm.mutex.RLock()
	defer sm.mutex.RUnlock()

	sm.logger.LogDebug("Verificando integridade do estado", nil)

	// Carregar estado
	state, err := sm.loadStateUnsafe()
	if err != nil {
		return types.ErrIntegrityCheckError("state_load", err)
	}

	// Verificar campos obrigatórios
	if state.Version == "" {
		return types.ErrIntegrityCheckError("version_missing", fmt.Errorf("versão não encontrada"))
	}

	if state.CreatedAt.IsZero() {
		return types.ErrIntegrityCheckError("created_at_missing", fmt.Errorf("data de criação não encontrada"))
	}

	if state.Status == "" {
		return types.ErrIntegrityCheckError("status_missing", fmt.Errorf("status não encontrado"))
	}

	// Verificar se o arquivo de estado não está corrompido
	if _, err := os.Stat(sm.statePath); err != nil {
		return types.ErrIntegrityCheckError("file_access", err)
	}

	sm.logger.LogDebug("Integridade do estado verificada com sucesso", map[string]interface{}{
		"version": state.Version,
		"status":  state.Status,
	})

	return nil
}

// GetStatePath retorna o caminho do arquivo de estado
func (sm *StateManager) GetStatePath() string {
	return sm.statePath
}

// GetBackupPath retorna o caminho do diretório de backups
func (sm *StateManager) GetBackupPath() string {
	return filepath.Join(filepath.Dir(sm.statePath), "backups")
}

// ListBackups lista os backups disponíveis
func (sm *StateManager) ListBackups() ([]types.BackupInfo, error) {
	backupDir := sm.GetBackupPath()

	// Verificar se o diretório de backups existe
	if _, err := os.Stat(backupDir); os.IsNotExist(err) {
		return []types.BackupInfo{}, nil
	}

	// Ler arquivos de backup
	files, err := os.ReadDir(backupDir)
	if err != nil {
		return nil, fmt.Errorf("falha ao ler diretório de backups: %w", err)
	}

	var backups []types.BackupInfo
	for _, file := range files {
		if filepath.Ext(file.Name()) == ".json" {
			filePath := filepath.Join(backupDir, file.Name())
			info, err := file.Info()
			if err != nil {
				continue
			}

			backups = append(backups, types.BackupInfo{
				ID:        file.Name(),
				Path:      filePath,
				Size:      info.Size(),
				CreatedAt: info.ModTime(),
			})
		}
	}

	return backups, nil
}

// CleanupOldBackups remove backups antigos
func (sm *StateManager) CleanupOldBackups(retentionDays int) error {
	sm.logger.LogStep("cleanup_backups_start", map[string]interface{}{
		"retention_days": retentionDays,
	})

	backups, err := sm.ListBackups()
	if err != nil {
		return err
	}

	cutoffTime := time.Now().AddDate(0, 0, -retentionDays)
	removedCount := 0

	for _, backup := range backups {
		if backup.CreatedAt.Before(cutoffTime) {
			if err := os.Remove(backup.Path); err != nil {
				sm.logger.LogWarning("Falha ao remover backup antigo", map[string]interface{}{
					"backup_path": backup.Path,
					"error":       err.Error(),
				})
			} else {
				removedCount++
			}
		}
	}

	sm.logger.LogStep("cleanup_backups_completed", map[string]interface{}{
		"removed_count": removedCount,
	})

	return nil
}

// Métodos auxiliares (sem locks - devem ser chamados com locks já adquiridos)

// loadStateUnsafe carrega o estado sem adquirir locks
func (sm *StateManager) loadStateUnsafe() (*types.SetupState, error) {
	// Verificar se o arquivo de estado existe
	if _, err := os.Stat(sm.statePath); os.IsNotExist(err) {
		return nil, types.ErrStateLoadError(fmt.Errorf("arquivo de estado não encontrado: %s", sm.statePath))
	}

	// Ler arquivo de estado
	data, err := os.ReadFile(sm.statePath)
	if err != nil {
		return nil, types.ErrStateLoadError(err)
	}

	// Deserializar estado
	var state types.SetupState
	if err := json.Unmarshal(data, &state); err != nil {
		return nil, types.ErrStateCorruptedError(err)
	}

	return &state, nil
}

// saveStateUnsafe salva o estado sem adquirir locks
func (sm *StateManager) saveStateUnsafe(state *types.SetupState) error {
	// Atualizar timestamp
	state.UpdatedAt = time.Now()

	// Serializar estado
	data, err := json.MarshalIndent(state, "", "  ")
	if err != nil {
		return types.ErrStateSaveError(err)
	}

	// Criar arquivo temporário para operação atômica
	tempPath := sm.statePath + ".tmp"

	// Escrever para arquivo temporário
	if err := os.WriteFile(tempPath, data, 0644); err != nil {
		return types.ErrStateSaveError(err)
	}

	// Renomear arquivo temporário para arquivo final (operação atômica)
	if err := os.Rename(tempPath, sm.statePath); err != nil {
		// Limpar arquivo temporário em caso de erro
		os.Remove(tempPath)
		return types.ErrStateSaveError(err)
	}

	return nil
}
