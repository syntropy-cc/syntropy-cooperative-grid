package node

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"node-component/src/internal/types"
)

// ProjectResourceManager gerencia acesso a recursos do projeto
type ProjectResourceManager struct {
	projectRoot string
	logger      types.Logger
}

// NewProjectResourceManager cria um novo gerenciador de recursos
func NewProjectResourceManager(logger types.Logger) (*ProjectResourceManager, error) {
	projectRoot, err := findProjectRoot()
	if err != nil {
		return nil, fmt.Errorf("failed to find project root: %w", err)
	}

	logger.Debug("Project root found", "path", projectRoot)

	return &ProjectResourceManager{
		projectRoot: projectRoot,
		logger:      logger,
	}, nil
}

// findProjectRoot encontra o diretório raiz do projeto
func findProjectRoot() (string, error) {
	// Começar do diretório atual
	currentDir, err := os.Getwd()
	if err != nil {
		return "", fmt.Errorf("failed to get current directory: %w", err)
	}

	// Subir na hierarquia até encontrar o diretório raiz
	searchDir := currentDir
	for i := 0; i < 10; i++ { // Máximo 10 níveis para evitar loop infinito
		// Verificar se este diretório contém infrastructure/cloud-init/
		testPath := filepath.Join(searchDir, "infrastructure", "cloud-init")
		if _, err := os.Stat(testPath); err == nil {
			return searchDir, nil
		}

		// Subir um nível
		parentDir := filepath.Dir(searchDir)
		if parentDir == searchDir {
			// Chegou na raiz do filesystem
			break
		}
		searchDir = parentDir
	}

	return "", fmt.Errorf("project root not found - infrastructure/cloud-init/ directory not found")
}

// GetCloudInitTemplate retorna o caminho para um template cloud-init
func (prm *ProjectResourceManager) GetCloudInitTemplate(templateName string) (string, error) {
	templatePath := filepath.Join(prm.projectRoot, "infrastructure", "cloud-init", templateName)

	if _, err := os.Stat(templatePath); os.IsNotExist(err) {
		return "", fmt.Errorf("template not found: %s", templatePath)
	}

	prm.logger.Debug("Cloud-init template found", "template", templateName, "path", templatePath)
	return templatePath, nil
}

// GetCloudInitScript retorna o caminho para um script cloud-init
func (prm *ProjectResourceManager) GetCloudInitScript(scriptName string) (string, error) {
	scriptPath := filepath.Join(prm.projectRoot, "infrastructure", "cloud-init", "scripts", scriptName)

	if _, err := os.Stat(scriptPath); os.IsNotExist(err) {
		return "", fmt.Errorf("script not found: %s", scriptPath)
	}

	prm.logger.Debug("Cloud-init script found", "script", scriptName, "path", scriptPath)
	return scriptPath, nil
}

// ReadCloudInitTemplate lê o conteúdo de um template cloud-init
func (prm *ProjectResourceManager) ReadCloudInitTemplate(templateName string) ([]byte, error) {
	templatePath, err := prm.GetCloudInitTemplate(templateName)
	if err != nil {
		return nil, err
	}

	content, err := os.ReadFile(templatePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read template %s: %w", templateName, err)
	}

	prm.logger.Debug("Cloud-init template read", "template", templateName, "size", len(content))
	return content, nil
}

// ReadCloudInitScript lê o conteúdo de um script cloud-init
func (prm *ProjectResourceManager) ReadCloudInitScript(scriptName string) ([]byte, error) {
	scriptPath, err := prm.GetCloudInitScript(scriptName)
	if err != nil {
		return nil, err
	}

	content, err := os.ReadFile(scriptPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read script %s: %w", scriptName, err)
	}

	prm.logger.Debug("Cloud-init script read", "script", scriptName, "size", len(content))
	return content, nil
}

// ListCloudInitScripts lista todos os scripts disponíveis
func (prm *ProjectResourceManager) ListCloudInitScripts() ([]string, error) {
	scriptsDir := filepath.Join(prm.projectRoot, "infrastructure", "cloud-init", "scripts")

	// Verificar se o diretório existe
	if _, err := os.Stat(scriptsDir); os.IsNotExist(err) {
		return nil, fmt.Errorf("scripts directory not found: %s", scriptsDir)
	}

	// Ler diretório
	entries, err := os.ReadDir(scriptsDir)
	if err != nil {
		return nil, fmt.Errorf("failed to read scripts directory: %w", err)
	}

	var scripts []string
	for _, entry := range entries {
		if !entry.IsDir() && strings.HasSuffix(entry.Name(), ".sh") {
			scripts = append(scripts, entry.Name())
		}
	}

	prm.logger.Debug("Cloud-init scripts listed", "count", len(scripts), "scripts", scripts)
	return scripts, nil
}

// GetProjectRoot retorna o caminho raiz do projeto
func (prm *ProjectResourceManager) GetProjectRoot() string {
	return prm.projectRoot
}

// GetInfrastructurePath retorna o caminho para o diretório infrastructure
func (prm *ProjectResourceManager) GetInfrastructurePath() string {
	return filepath.Join(prm.projectRoot, "infrastructure")
}

// GetCloudInitPath retorna o caminho para o diretório cloud-init
func (prm *ProjectResourceManager) GetCloudInitPath() string {
	return filepath.Join(prm.projectRoot, "infrastructure", "cloud-init")
}

// GetCloudInitScriptsPath retorna o caminho para o diretório de scripts cloud-init
func (prm *ProjectResourceManager) GetCloudInitScriptsPath() string {
	return filepath.Join(prm.projectRoot, "infrastructure", "cloud-init", "scripts")
}

// ValidateInfrastructureStructure valida se a estrutura do projeto está correta
func (prm *ProjectResourceManager) ValidateInfrastructureStructure() error {
	requiredPaths := []string{
		prm.GetInfrastructurePath(),
		prm.GetCloudInitPath(),
		prm.GetCloudInitScriptsPath(),
	}

	requiredTemplates := []string{
		"user-data-template.yaml",
		"meta-data-template.yaml",
		"network-config-template.yaml",
	}

	// Verificar diretórios
	for _, path := range requiredPaths {
		if _, err := os.Stat(path); os.IsNotExist(err) {
			return fmt.Errorf("required directory not found: %s", path)
		}
	}

	// Verificar templates
	for _, template := range requiredTemplates {
		if _, err := prm.GetCloudInitTemplate(template); err != nil {
			return fmt.Errorf("required template not found: %s", template)
		}
	}

	prm.logger.Debug("Infrastructure structure validated successfully")
	return nil
}

