package setup

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	"gopkg.in/yaml.v3"

	"github.com/syntropy-cc/syntropy-cooperative-grid/manager/interfaces/cli/setup/src/internal/types"
)

// Configurator implementa a interface Configurator
type Configurator struct {
	configDir    string
	templatesDir string
	logger       *SetupLogger
}

// NewConfigurator cria um novo configurador
func NewConfigurator(logger *SetupLogger) *Configurator {
	homeDir, _ := os.UserHomeDir()
	configDir := filepath.Join(homeDir, ".syntropy", "config")
	templatesDir := filepath.Join(homeDir, ".syntropy", "templates")

	// Criar diretórios se não existirem
	os.MkdirAll(configDir, 0755)
	os.MkdirAll(templatesDir, 0755)

	return &Configurator{
		configDir:    configDir,
		templatesDir: templatesDir,
		logger:       logger,
	}
}

// GenerateConfig gera a configuração principal
func (c *Configurator) GenerateConfig(options *types.ConfigOptions) error {
	c.logger.LogStep("config_generation_start", map[string]interface{}{
		"owner_name":  options.OwnerName,
		"owner_email": options.OwnerEmail,
	})

	// Criar estrutura de diretórios se necessário
	if err := c.CreateStructure(); err != nil {
		return types.ErrConfigGenerationError(err)
	}

	// Gerar configuração principal
	config := &types.SetupConfig{
		Manager: types.ManagerConfig{
			HomeDir:     filepath.Join(os.Getenv("HOME"), ".syntropy"),
			LogLevel:    "info",
			APIEndpoint: "https://api.syntropy.network",
			Directories: map[string]string{
				"config":  filepath.Join(os.Getenv("HOME"), ".syntropy", "config"),
				"keys":    filepath.Join(os.Getenv("HOME"), ".syntropy", "keys"),
				"nodes":   filepath.Join(os.Getenv("HOME"), ".syntropy", "nodes"),
				"logs":    filepath.Join(os.Getenv("HOME"), ".syntropy", "logs"),
				"cache":   filepath.Join(os.Getenv("HOME"), ".syntropy", "cache"),
				"backups": filepath.Join(os.Getenv("HOME"), ".syntropy", "backups"),
			},
			DefaultPaths: map[string]string{
				"config": filepath.Join(os.Getenv("HOME"), ".syntropy", "config", "manager.yaml"),
				"log":    filepath.Join(os.Getenv("HOME"), ".syntropy", "logs", "manager.log"),
			},
		},
		OwnerKey: types.OwnerKey{
			Type: "ed25519",
			Path: filepath.Join(os.Getenv("HOME"), ".syntropy", "keys", "owner.key"),
		},
		Environment: types.Environment{
			OS:           runtime.GOOS,
			Architecture: runtime.GOARCH,
			HomeDir:      os.Getenv("HOME"),
		},
	}

	// Aplicar configurações personalizadas
	if options.NetworkConfig != nil {
		// Aplicar configuração de rede se fornecida
	}

	if options.SecurityConfig != nil {
		// Aplicar configuração de segurança se fornecida
	}

	// Salvar configuração
	configPath := filepath.Join(c.configDir, "manager.yaml")
	if err := c.saveConfig(config, configPath); err != nil {
		return types.ErrConfigGenerationError(err)
	}

	c.logger.LogStep("config_generation_completed", map[string]interface{}{
		"config_path": configPath,
	})

	return nil
}

// CreateStructure cria a estrutura de diretórios necessária
func (c *Configurator) CreateStructure() error {
	c.logger.LogStep("structure_creation_start", nil)

	homeDir, _ := os.UserHomeDir()
	baseDir := filepath.Join(homeDir, ".syntropy")

	directories := []string{
		filepath.Join(baseDir, "config"),
		filepath.Join(baseDir, "keys"),
		filepath.Join(baseDir, "nodes"),
		filepath.Join(baseDir, "logs"),
		filepath.Join(baseDir, "cache"),
		filepath.Join(baseDir, "backups"),
		filepath.Join(baseDir, "templates"),
		filepath.Join(baseDir, "state"),
	}

	for _, dir := range directories {
		if err := os.MkdirAll(dir, 0755); err != nil {
			return types.ErrStructureCreationError(dir, err)
		}
	}

	c.logger.LogStep("structure_creation_completed", map[string]interface{}{
		"base_dir":    baseDir,
		"directories": directories,
	})

	return nil
}

// GenerateKeys gera chaves criptográficas
func (c *Configurator) GenerateKeys() (*types.KeyPair, error) {
	c.logger.LogStep("key_generation_start", nil)

	// Usar KeyManager para gerar chaves
	keyManager := NewKeyManager(c.logger)
	keyPair, err := keyManager.GenerateKeyPair("ed25519")
	if err != nil {
		return nil, types.ErrKeyGenerationError("ed25519", err)
	}

	// Armazenar chaves
	if err := keyManager.StoreKeyPair(keyPair, "default_passphrase"); err != nil {
		return nil, types.ErrKeyStorageError(keyPair.ID, err)
	}

	c.logger.LogStep("key_generation_completed", map[string]interface{}{
		"key_id":      keyPair.ID,
		"algorithm":   keyPair.Algorithm,
		"fingerprint": keyPair.Fingerprint,
	})

	return keyPair, nil
}

// ValidateConfig valida a configuração
func (c *Configurator) ValidateConfig() error {
	c.logger.LogStep("config_validation_start", nil)

	configPath := filepath.Join(c.configDir, "manager.yaml")

	// Verificar se o arquivo de configuração existe
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		return types.ErrConfigCorruptedError(configPath, err)
	}

	// Ler e validar configuração
	config, err := c.loadConfig(configPath)
	if err != nil {
		return types.ErrConfigCorruptedError(configPath, err)
	}

	// Validar campos obrigatórios
	if config.Manager.HomeDir == "" {
		return fmt.Errorf("campo obrigatório ausente: manager.home_dir")
	}

	if config.OwnerKey.Type == "" {
		return fmt.Errorf("campo obrigatório ausente: owner_key.type")
	}

	if config.Environment.OS == "" {
		return fmt.Errorf("campo obrigatório ausente: environment.os")
	}

	c.logger.LogStep("config_validation_completed", map[string]interface{}{
		"config_path": configPath,
		"valid":       true,
	})

	return nil
}

// BackupConfig cria um backup da configuração
func (c *Configurator) BackupConfig(name string) error {
	c.logger.LogStep("config_backup_start", map[string]interface{}{
		"backup_name": name,
	})

	configPath := filepath.Join(c.configDir, "manager.yaml")

	// Verificar se o arquivo de configuração existe
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		return types.ErrBackupFailedError(fmt.Errorf("arquivo de configuração não encontrado: %s", configPath))
	}

	// Criar diretório de backups
	backupDir := filepath.Join(filepath.Dir(c.configDir), "backups")
	os.MkdirAll(backupDir, 0755)

	// Gerar nome do arquivo de backup
	timestamp := time.Now().Format("20060102_150405")
	backupPath := filepath.Join(backupDir, fmt.Sprintf("config_%s_%s.yaml", name, timestamp))

	// Ler arquivo de configuração
	configData, err := os.ReadFile(configPath)
	if err != nil {
		return types.ErrBackupFailedError(err)
	}

	// Escrever backup
	if err := os.WriteFile(backupPath, configData, 0644); err != nil {
		return types.ErrBackupFailedError(err)
	}

	c.logger.LogStep("config_backup_completed", map[string]interface{}{
		"backup_path": backupPath,
		"backup_size": len(configData),
	})

	return nil
}

// RestoreConfig restaura a configuração de um backup
func (c *Configurator) RestoreConfig(backupPath string) error {
	c.logger.LogStep("config_restore_start", map[string]interface{}{
		"backup_path": backupPath,
	})

	// Verificar se o arquivo de backup existe
	if _, err := os.Stat(backupPath); os.IsNotExist(err) {
		return types.ErrRestoreFailedError(fmt.Errorf("arquivo de backup não encontrado: %s", backupPath))
	}

	// Ler arquivo de backup
	backupData, err := os.ReadFile(backupPath)
	if err != nil {
		return types.ErrRestoreFailedError(err)
	}

	// Validar configuração do backup
	var config types.SetupConfig
	if err := yaml.Unmarshal(backupData, &config); err != nil {
		return types.ErrRestoreFailedError(fmt.Errorf("configuração de backup inválida: %w", err))
	}

	// Criar backup da configuração atual antes de restaurar
	currentConfigPath := filepath.Join(c.configDir, "manager.yaml")
	if _, err := os.Stat(currentConfigPath); err == nil {
		if err := c.BackupConfig("pre_restore"); err != nil {
			c.logger.LogWarning("Falha ao criar backup da configuração atual", map[string]interface{}{
				"error": err.Error(),
			})
		}
	}

	// Escrever configuração restaurada
	if err := os.WriteFile(currentConfigPath, backupData, 0644); err != nil {
		return types.ErrRestoreFailedError(err)
	}

	c.logger.LogStep("config_restore_completed", map[string]interface{}{
		"backup_path": backupPath,
		"config_path": currentConfigPath,
	})

	return nil
}

// LoadTemplate carrega um template de configuração
func (c *Configurator) LoadTemplate(templateName string) (string, error) {
	c.logger.LogDebug("Carregando template", map[string]interface{}{
		"template_name": templateName,
	})

	templatePath := filepath.Join(c.templatesDir, templateName)

	// Verificar se o template existe
	if _, err := os.Stat(templatePath); os.IsNotExist(err) {
		return "", fmt.Errorf("template não encontrado: %s", templateName)
	}

	// Ler template
	templateData, err := os.ReadFile(templatePath)
	if err != nil {
		return "", fmt.Errorf("falha ao ler template: %w", err)
	}

	return string(templateData), nil
}

// SaveTemplate salva um template de configuração
func (c *Configurator) SaveTemplate(templateName string, content string) error {
	c.logger.LogDebug("Salvando template", map[string]interface{}{
		"template_name": templateName,
	})

	templatePath := filepath.Join(c.templatesDir, templateName)

	// Escrever template
	if err := os.WriteFile(templatePath, []byte(content), 0644); err != nil {
		return fmt.Errorf("falha ao salvar template: %w", err)
	}

	return nil
}

// ProcessTemplate processa um template com variáveis
func (c *Configurator) ProcessTemplate(templateName string, variables map[string]interface{}) (string, error) {
	c.logger.LogDebug("Processando template", map[string]interface{}{
		"template_name": templateName,
		"variables":     variables,
	})

	// Carregar template
	templateContent, err := c.LoadTemplate(templateName)
	if err != nil {
		return "", err
	}

	// Processar template (implementação simplificada)
	// Em produção, usar um processador de templates mais robusto
	processedContent := templateContent

	// Substituir variáveis simples
	for key, value := range variables {
		placeholder := fmt.Sprintf("{{.%s}}", key)
		processedContent = strings.ReplaceAll(processedContent, placeholder, fmt.Sprintf("%v", value))
	}

	return processedContent, nil
}

// GetConfigPath retorna o caminho do arquivo de configuração
func (c *Configurator) GetConfigPath() string {
	return filepath.Join(c.configDir, "manager.yaml")
}

// GetTemplatesDir retorna o diretório de templates
func (c *Configurator) GetTemplatesDir() string {
	return c.templatesDir
}

// ListTemplates lista os templates disponíveis
func (c *Configurator) ListTemplates() ([]string, error) {
	files, err := os.ReadDir(c.templatesDir)
	if err != nil {
		return nil, fmt.Errorf("falha ao ler diretório de templates: %w", err)
	}

	var templates []string
	for _, file := range files {
		if !file.IsDir() && (filepath.Ext(file.Name()) == ".yaml" || filepath.Ext(file.Name()) == ".yml") {
			templates = append(templates, file.Name())
		}
	}

	return templates, nil
}

// Métodos auxiliares

// saveConfig salva a configuração em um arquivo
func (c *Configurator) saveConfig(config *types.SetupConfig, path string) error {
	// Serializar para YAML
	data, err := yaml.Marshal(config)
	if err != nil {
		return fmt.Errorf("falha ao serializar configuração: %w", err)
	}

	// Escrever arquivo
	if err := os.WriteFile(path, data, 0644); err != nil {
		return fmt.Errorf("falha ao escrever arquivo de configuração: %w", err)
	}

	return nil
}

// loadConfig carrega a configuração de um arquivo
func (c *Configurator) loadConfig(path string) (*types.SetupConfig, error) {
	// Ler arquivo
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("falha ao ler arquivo de configuração: %w", err)
	}

	// Deserializar de YAML
	var config types.SetupConfig
	if err := yaml.Unmarshal(data, &config); err != nil {
		return nil, fmt.Errorf("falha ao deserializar configuração: %w", err)
	}

	return &config, nil
}

// validateConfigStructure valida a estrutura da configuração
func (c *Configurator) validateConfigStructure(config *types.SetupConfig) error {
	// Validar campos obrigatórios
	if config.Manager.HomeDir == "" {
		return fmt.Errorf("campo obrigatório ausente: manager.home_dir")
	}

	if config.OwnerKey.Type == "" {
		return fmt.Errorf("campo obrigatório ausente: owner_key.type")
	}

	if config.Environment.OS == "" {
		return fmt.Errorf("campo obrigatório ausente: environment.os")
	}

	// Validar diretórios
	for name, path := range config.Manager.Directories {
		if path == "" {
			return fmt.Errorf("caminho vazio para diretório: %s", name)
		}
	}

	// Validar caminhos padrão
	for name, path := range config.Manager.DefaultPaths {
		if path == "" {
			return fmt.Errorf("caminho vazio para arquivo padrão: %s", name)
		}
	}

	return nil
}
