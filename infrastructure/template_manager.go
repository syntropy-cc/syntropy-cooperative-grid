package infrastructure

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"text/template"
)

// TemplateManager gerencia templates de cloud-init
type TemplateManager struct {
	templateDir string
}

// TemplateData contém dados para renderização de templates
type TemplateData struct {
	NodeName             string
	NodeDescription      string
	Coordinates          string
	CreatedAt            string
	AdminPasswordHash    string
	OwnerPublicKey       string
	CommunityPublicKey   string
	KeyInstallationCmds  string
	MetadataCreationCmds string
}

// NewTemplateManager cria um novo gerenciador de templates
func NewTemplateManager(templateDir string) *TemplateManager {
	return &TemplateManager{
		templateDir: templateDir,
	}
}

// RenderTemplate renderiza um template com os dados fornecidos
func (tm *TemplateManager) RenderTemplate(templateName string, data *TemplateData) (string, error) {
	templatePath := filepath.Join(tm.templateDir, templateName)

	// Verificar se o arquivo existe
	if _, err := os.Stat(templatePath); os.IsNotExist(err) {
		return "", fmt.Errorf("template não encontrado: %s", templatePath)
	}

	// Ler o template
	templateContent, err := os.ReadFile(templatePath)
	if err != nil {
		return "", fmt.Errorf("erro ao ler template: %w", err)
	}

	// Criar e executar o template
	tmpl, err := template.New(templateName).Parse(string(templateContent))
	if err != nil {
		return "", fmt.Errorf("erro ao fazer parse do template: %w", err)
	}

	var buf bytes.Buffer
	err = tmpl.Execute(&buf, data)
	if err != nil {
		return "", fmt.Errorf("erro ao executar template: %w", err)
	}

	return buf.String(), nil
}

// SaveTemplate salva um template renderizado em um arquivo
func (tm *TemplateManager) SaveTemplate(templateName string, data *TemplateData, outputPath string) error {
	content, err := tm.RenderTemplate(templateName, data)
	if err != nil {
		return err
	}

	// Criar diretório se não existir
	dir := filepath.Dir(outputPath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("erro ao criar diretório: %w", err)
	}

	// Salvar arquivo
	if err := os.WriteFile(outputPath, []byte(content), 0644); err != nil {
		return fmt.Errorf("erro ao salvar template: %w", err)
	}

	return nil
}
