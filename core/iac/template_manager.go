package iac

import (
	"fmt"
	"html/template"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"time"
)

// TemplateManager gerencia templates de Infrastructure as Code
type TemplateManager struct {
	templateDir string
}

// TemplateData contém dados para processamento de templates
type TemplateData struct {
	NodeName              string
	NodeDescription       string
	Coordinates           string
	CreatedAt             string
	AdminPasswordHash     string
	OwnerPublicKey        string
	CommunityPublicKey    string
	KeyInstallationCommands string
	MetadataCreationCommands string
	TemplateCreationCommands string
	StartupServiceCommands string
	NodeID                string
	LocationNodeID        string
	DetectionMethod       string
	DetectedCity          string
	DetectedCountry       string
	OwnerFingerprint      string
	CommunityFingerprint  string
}

// NewTemplateManager cria um novo gerenciador de templates
func NewTemplateManager(templateDir string) *TemplateManager {
	return &TemplateManager{
		templateDir: templateDir,
	}
}

// LoadTemplate carrega um template de arquivo
func (tm *TemplateManager) LoadTemplate(templatePath string) (string, error) {
	fullPath := filepath.Join(tm.templateDir, templatePath)
	
	content, err := ioutil.ReadFile(fullPath)
	if err != nil {
		return "", fmt.Errorf("falha ao ler template %s: %w", fullPath, err)
	}
	
	return string(content), nil
}

// ProcessTemplate processa um template com os dados fornecidos
func (tm *TemplateManager) ProcessTemplate(templateContent string, data *TemplateData) (string, error) {
	// Criar template com funções customizadas
	tmpl := template.New("iac-template").Funcs(template.FuncMap{
		"now": func() string {
			return time.Now().Format(time.RFC3339)
		},
		"split": func(s, sep string) []string {
			return strings.Split(s, sep)
		},
		"trim": strings.TrimSpace,
		"upper": strings.ToUpper,
		"lower": strings.ToLower,
	})
	
	// Parse do template
	parsedTemplate, err := tmpl.Parse(templateContent)
	if err != nil {
		return "", fmt.Errorf("falha ao fazer parse do template: %w", err)
	}
	
	// Processar template
	var result strings.Builder
	if err := parsedTemplate.Execute(&result, data); err != nil {
		return "", fmt.Errorf("falha ao executar template: %w", err)
	}
	
	return result.String(), nil
}

// GenerateUserData gera configuração user-data usando template
func (tm *TemplateManager) GenerateUserData(data *TemplateData) (string, error) {
	templateContent, err := tm.LoadTemplate("cloud-init/user-data-template.yaml")
	if err != nil {
		return "", err
	}
	
	return tm.ProcessTemplate(templateContent, data)
}

// GenerateMetaData gera configuração meta-data usando template
func (tm *TemplateManager) GenerateMetaData(data *TemplateData) (string, error) {
	templateContent, err := tm.LoadTemplate("cloud-init/meta-data-template.yaml")
	if err != nil {
		return "", err
	}
	
	return tm.ProcessTemplate(templateContent, data)
}

// GenerateNetworkConfig gera configuração de rede usando template
func (tm *TemplateManager) GenerateNetworkConfig(data *TemplateData) (string, error) {
	templateContent, err := tm.LoadTemplate("cloud-init/network-config-template.yaml")
	if err != nil {
		return "", err
	}
	
	return tm.ProcessTemplate(templateContent, data)
}

// SaveCloudInitFiles salva arquivos cloud-init em um diretório
func (tm *TemplateManager) SaveCloudInitFiles(outputDir string, data *TemplateData) error {
	// Criar diretório se não existir
	if err := os.MkdirAll(outputDir, 0755); err != nil {
		return fmt.Errorf("falha ao criar diretório %s: %w", outputDir, err)
	}
	
	// Gerar user-data
	userData, err := tm.GenerateUserData(data)
	if err != nil {
		return fmt.Errorf("falha ao gerar user-data: %w", err)
	}
	
	userDataPath := filepath.Join(outputDir, "user-data")
	if err := ioutil.WriteFile(userDataPath, []byte(userData), 0644); err != nil {
		return fmt.Errorf("falha ao salvar user-data: %w", err)
	}
	
	// Gerar meta-data
	metaData, err := tm.GenerateMetaData(data)
	if err != nil {
		return fmt.Errorf("falha ao gerar meta-data: %w", err)
	}
	
	metaDataPath := filepath.Join(outputDir, "meta-data")
	if err := ioutil.WriteFile(metaDataPath, []byte(metaData), 0644); err != nil {
		return fmt.Errorf("falha ao salvar meta-data: %w", err)
	}
	
	// Gerar network-config
	networkConfig, err := tm.GenerateNetworkConfig(data)
	if err != nil {
		return fmt.Errorf("falha ao gerar network-config: %w", err)
	}
	
	networkConfigPath := filepath.Join(outputDir, "network-config")
	if err := ioutil.WriteFile(networkConfigPath, []byte(networkConfig), 0644); err != nil {
		return fmt.Errorf("falha ao salvar network-config: %w", err)
	}
	
	return nil
}

// ValidateTemplate valida se um template está correto
func (tm *TemplateManager) ValidateTemplate(templateContent string) error {
	// Tentar fazer parse do template
	tmpl := template.New("validation").Funcs(template.FuncMap{
		"now": func() string { return time.Now().Format(time.RFC3339) },
		"split": func(s, sep string) []string { return strings.Split(s, sep) },
		"trim": strings.TrimSpace,
		"upper": strings.ToUpper,
		"lower": strings.ToLower,
	})
	
	_, err := tmpl.Parse(templateContent)
	if err != nil {
		return fmt.Errorf("template inválido: %w", err)
	}
	
	return nil
}

// ListTemplates lista todos os templates disponíveis
func (tm *TemplateManager) ListTemplates() ([]string, error) {
	var templates []string
	
	err := filepath.Walk(tm.templateDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		
		// Verificar se é um arquivo de template
		if !info.IsDir() && (strings.HasSuffix(path, ".yaml") || strings.HasSuffix(path, ".yml") || strings.HasSuffix(path, ".tpl")) {
			relPath, err := filepath.Rel(tm.templateDir, path)
			if err != nil {
				return err
			}
			templates = append(templates, relPath)
		}
		
		return nil
	})
	
	if err != nil {
		return nil, fmt.Errorf("falha ao listar templates: %w", err)
	}
	
	return templates, nil
}

// GetTemplateInfo obtém informações sobre um template
func (tm *TemplateManager) GetTemplateInfo(templatePath string) (map[string]interface{}, error) {
	templateContent, err := tm.LoadTemplate(templatePath)
	if err != nil {
		return nil, err
	}
	
	info := map[string]interface{}{
		"path":    templatePath,
		"size":    len(templateContent),
		"valid":   tm.ValidateTemplate(templateContent) == nil,
		"modified": time.Now().Format(time.RFC3339), // Placeholder
	}
	
	// Detectar tipo de template
	if strings.Contains(templateContent, "#cloud-config") {
		info["type"] = "cloud-init"
	} else if strings.Contains(templateContent, "autoinstall:") {
		info["type"] = "ubuntu-autoinstall"
	} else if strings.Contains(templateContent, "{{") {
		info["type"] = "template"
	} else {
		info["type"] = "yaml"
	}
	
	return info, nil
}
