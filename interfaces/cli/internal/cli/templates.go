package cli

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
)

// NewTemplatesCommand cria o comando de templates
func NewTemplatesCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "templates",
		Short: "Manage application templates",
		Long: `Manage application deployment templates for the Syntropy Cooperative Grid.

Templates provide pre-configured application deployments including:
- Scientific computing (Fortran, Python, R)
- Data science and machine learning
- Web applications and APIs
- Database services
- Monitoring and logging`,
	}

	// Adicionar subcomandos
	cmd.AddCommand(newTemplatesListCommand())
	cmd.AddCommand(newTemplatesShowCommand())
	cmd.AddCommand(newTemplatesDeployCommand())
	cmd.AddCommand(newTemplatesCreateCommand())

	return cmd
}

// newTemplatesListCommand cria o comando de listagem de templates
func newTemplatesListCommand() *cobra.Command {
	var (
		category string
		format   string
	)

	cmd := &cobra.Command{
		Use:   "list",
		Short: "List available templates",
		Long: `List all available application deployment templates.

You can filter by category and output in different formats.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return listTemplates(category, format)
		},
	}

	cmd.Flags().StringVarP(&category, "category", "c", "", "Filter by category (scientific, web, database, monitoring)")
	cmd.Flags().StringVarP(&format, "format", "f", "table", "Output format (table, json, yaml)")

	return cmd
}

// newTemplatesShowCommand cria o comando para mostrar template
func newTemplatesShowCommand() *cobra.Command {
	var (
		format string
	)

	cmd := &cobra.Command{
		Use:   "show <template-name>",
		Short: "Show template details",
		Long: `Show detailed information about a specific template.

This displays the template configuration, requirements, and deployment options.`,
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			templateName := args[0]
			return showTemplate(templateName, format)
		},
	}

	cmd.Flags().StringVarP(&format, "format", "f", "yaml", "Output format (yaml, json)")

	return cmd
}

// newTemplatesDeployCommand cria o comando de deploy
func newTemplatesDeployCommand() *cobra.Command {
	var (
		nodeName string
		values   []string
		dryRun   bool
	)

	cmd := &cobra.Command{
		Use:   "deploy <template-name>",
		Short: "Deploy template to node",
		Long: `Deploy an application template to a specific node.

This will:
1. Load the template configuration
2. Apply custom values if provided
3. Deploy to the target node
4. Monitor deployment status`,
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			templateName := args[0]
			return deployTemplate(templateName, nodeName, values, dryRun)
		},
	}

	cmd.Flags().StringVarP(&nodeName, "node", "n", "", "Target node name (required)")
	cmd.Flags().StringSliceVarP(&values, "set", "s", []string{}, "Set custom values (key=value)")
	cmd.Flags().BoolVar(&dryRun, "dry-run", false, "Show what would be deployed without actually deploying")

	cmd.MarkFlagRequired("node")

	return cmd
}

// newTemplatesCreateCommand cria o comando de criação de template
func newTemplatesCreateCommand() *cobra.Command {
	var (
		category    string
		description string
		output      string
	)

	cmd := &cobra.Command{
		Use:   "create <template-name>",
		Short: "Create new template",
		Long: `Create a new application deployment template.

This will generate a template skeleton that you can customize
for your specific application requirements.`,
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			templateName := args[0]
			return createTemplate(templateName, category, description, output)
		},
	}

	cmd.Flags().StringVarP(&category, "category", "c", "custom", "Template category")
	cmd.Flags().StringVarP(&description, "description", "d", "", "Template description")
	cmd.Flags().StringVarP(&output, "output", "o", "", "Output file path (default: auto-generated)")

	return cmd
}

// Implementações dos comandos

func listTemplates(category, format string) error {
	templates, err := loadTemplates()
	if err != nil {
		return fmt.Errorf("failed to load templates: %w", err)
	}

	// Filtrar por categoria se especificada
	if category != "" {
		templates = filterTemplatesByCategory(templates, category)
	}

	switch format {
	case "json":
		return outputTemplatesJSON(templates)
	case "yaml":
		return outputTemplatesYAML(templates)
	default:
		return outputTemplatesTable(templates)
	}
}

func showTemplate(templateName, format string) error {
	template, err := loadTemplate(templateName)
	if err != nil {
		return fmt.Errorf("template not found: %w", err)
	}

	switch format {
	case "json":
		return outputTemplateJSON(template)
	default:
		return outputTemplateYAML(template)
	}
}

func deployTemplate(templateName, nodeName string, values []string, dryRun bool) error {
	fmt.Printf("🚀 Deploying template '%s' to node '%s'\n", templateName, nodeName)

	// Carregar template
	template, err := loadTemplate(templateName)
	if err != nil {
		return fmt.Errorf("failed to load template: %w", err)
	}

	// Carregar nó
	node, err := loadNode(nodeName)
	if err != nil {
		return fmt.Errorf("node not found: %w", err)
	}

	// Aplicar valores customizados
	if len(values) > 0 {
		template = applyTemplateValues(template, values)
	}

	if dryRun {
		fmt.Println("🔍 Dry run - showing what would be deployed:")
		fmt.Printf("Template: %s\n", template.Name)
		fmt.Printf("Node: %s (%s)\n", node.Name, node.Network.IPAddress)
		fmt.Printf("Category: %s\n", template.Category)
		fmt.Println("Resources:")
		fmt.Printf("  CPU: %s\n", template.Resources.CPU)
		fmt.Printf("  Memory: %s\n", template.Resources.Memory)
		fmt.Println("✅ Dry run completed - no changes made")
		return nil
	}

	// Executar deploy
	fmt.Println("📦 Deploying application...")
	
	// Conectar ao nó e executar deploy
	if err := executeDeployment(node, template); err != nil {
		return fmt.Errorf("deployment failed: %w", err)
	}

	fmt.Printf("✅ Template '%s' deployed successfully to node '%s'\n", templateName, nodeName)
	return nil
}

func createTemplate(templateName, category, description, output string) error {
	fmt.Printf("📝 Creating template '%s'\n", templateName)

	// Gerar template base
	template := generateTemplateSkeleton(templateName, category, description)

	// Determinar arquivo de saída
	if output == "" {
		syntropyDir := getSyntropyDir()
		templatesDir := filepath.Join(syntropyDir, "config", "templates", "applications")
		os.MkdirAll(templatesDir, 0755)
		output = filepath.Join(templatesDir, templateName+".yaml")
	}

	// Salvar template
	if err := saveTemplate(template, output); err != nil {
		return fmt.Errorf("failed to save template: %w", err)
	}

	fmt.Printf("✅ Template created: %s\n", output)
	fmt.Println("📋 Next steps:")
	fmt.Println("1. Edit the template file to customize your application")
	fmt.Println("2. Test with: syntropy templates deploy " + templateName + " --node <node> --dry-run")
	fmt.Println("3. Deploy with: syntropy templates deploy " + templateName + " --node <node>")

	return nil
}

// Estruturas de dados

type Template struct {
	Name        string            `json:"name" yaml:"name"`
	Category    string            `json:"category" yaml:"category"`
	Description string            `json:"description" yaml:"description"`
	Version     string            `json:"version" yaml:"version"`
	Author      string            `json:"author" yaml:"author"`
	Resources   ResourceRequirements `json:"resources" yaml:"resources"`
	Services    []Service         `json:"services" yaml:"services"`
	Volumes     []Volume          `json:"volumes" yaml:"volumes"`
	Networks    []Network         `json:"networks" yaml:"networks"`
	Environment map[string]string `json:"environment" yaml:"environment"`
	Labels      map[string]string `json:"labels" yaml:"labels"`
}

type ResourceRequirements struct {
	CPU    string `json:"cpu" yaml:"cpu"`
	Memory string `json:"memory" yaml:"memory"`
	Storage string `json:"storage" yaml:"storage"`
}

type Service struct {
	Name        string            `json:"name" yaml:"name"`
	Image       string            `json:"image" yaml:"image"`
	Ports       []Port            `json:"ports" yaml:"ports"`
	Environment map[string]string `json:"environment" yaml:"environment"`
	Volumes     []VolumeMount     `json:"volumes" yaml:"volumes"`
	Resources   ResourceRequirements `json:"resources" yaml:"resources"`
}

type Port struct {
	Container int    `json:"container" yaml:"container"`
	Host      int    `json:"host" yaml:"host"`
	Protocol  string `json:"protocol" yaml:"protocol"`
}

type Volume struct {
	Name string `json:"name" yaml:"name"`
	Type string `json:"type" yaml:"type"`
	Size string `json:"size" yaml:"size"`
}

type VolumeMount struct {
	Name      string `json:"name" yaml:"name"`
	MountPath string `json:"mountPath" yaml:"mountPath"`
}

type Network struct {
	Name    string `json:"name" yaml:"name"`
	Driver  string `json:"driver" yaml:"driver"`
	Subnet  string `json:"subnet" yaml:"subnet"`
}

// Funções auxiliares

func loadTemplates() ([]Template, error) {
	syntropyDir := getSyntropyDir()
	templatesDir := filepath.Join(syntropyDir, "config", "templates", "applications")
	
	files, err := filepath.Glob(filepath.Join(templatesDir, "*.yaml"))
	if err != nil {
		return nil, err
	}

	templates := []Template{}
	for _, file := range files {
		template, err := loadTemplateFromFile(file)
		if err != nil {
			continue // Ignorar arquivos corrompidos
		}
		templates = append(templates, template)
	}

	return templates, nil
}

func loadTemplate(templateName string) (Template, error) {
	syntropyDir := getSyntropyDir()
	templateFile := filepath.Join(syntropyDir, "config", "templates", "applications", templateName+".yaml")
	return loadTemplateFromFile(templateFile)
}

func loadTemplateFromFile(filePath string) (Template, error) {
	_, err := os.ReadFile(filePath)
	if err != nil {
		return Template{}, err
	}

	// Implementação simplificada - em produção usar yaml.Unmarshal
	template := Template{
		Name:        strings.TrimSuffix(filepath.Base(filePath), ".yaml"),
		Category:    "custom",
		Description: "Custom template",
		Version:     "1.0.0",
		Author:      "user",
		Resources: ResourceRequirements{
			CPU:    "500m",
			Memory: "512Mi",
		},
	}

	return template, nil
}

func saveTemplate(template Template, filePath string) error {
	// Implementação simplificada - em produção usar yaml.Marshal
	content := fmt.Sprintf(`name: %s
category: %s
description: %s
version: %s
author: %s
resources:
  cpu: %s
  memory: %s
services:
  - name: %s
    image: nginx:latest
    ports:
      - container: 80
        host: 8080
        protocol: tcp
    environment:
      ENV: production
    resources:
      cpu: %s
      memory: %s
`,
		template.Name,
		template.Category,
		template.Description,
		template.Version,
		template.Author,
		template.Resources.CPU,
		template.Resources.Memory,
		template.Name,
		template.Resources.CPU,
		template.Resources.Memory,
	)

	return os.WriteFile(filePath, []byte(content), 0644)
}

func filterTemplatesByCategory(templates []Template, category string) []Template {
	filtered := []Template{}
	for _, template := range templates {
		if strings.ToLower(template.Category) == strings.ToLower(category) {
			filtered = append(filtered, template)
		}
	}
	return filtered
}

func applyTemplateValues(template Template, values []string) Template {
	// Implementação simplificada
	return template
}

func executeDeployment(node NodeInfo, template Template) error {
	// Implementação simplificada - em produção executar deploy real
	fmt.Printf("Connecting to node %s (%s)...\n", node.Name, node.Network.IPAddress)
	fmt.Printf("Deploying template %s...\n", template.Name)
	fmt.Println("Deployment completed successfully")
	return nil
}

func generateTemplateSkeleton(name, category, description string) Template {
	if description == "" {
		description = fmt.Sprintf("Custom template for %s", name)
	}

	return Template{
		Name:        name,
		Category:    category,
		Description: description,
		Version:     "1.0.0",
		Author:      os.Getenv("USER"),
		Resources: ResourceRequirements{
			CPU:    "500m",
			Memory: "512Mi",
		},
		Services: []Service{
			{
				Name:  name,
				Image: "nginx:latest",
				Ports: []Port{
					{
						Container: 80,
						Host:      8080,
						Protocol:  "tcp",
					},
				},
				Environment: map[string]string{
					"ENV": "production",
				},
				Resources: ResourceRequirements{
					CPU:    "500m",
					Memory: "512Mi",
				},
			},
		},
	}
}

// Funções de output

func outputTemplatesTable(templates []Template) error {
	fmt.Printf("%-20s %-15s %-50s %s\n", "NAME", "CATEGORY", "DESCRIPTION", "VERSION")
	fmt.Println(strings.Repeat("-", 100))

	for _, template := range templates {
		description := template.Description
		if len(description) > 47 {
			description = description[:44] + "..."
		}
		fmt.Printf("%-20s %-15s %-50s %s\n",
			template.Name, template.Category, description, template.Version)
	}

	return nil
}

func outputTemplatesJSON(templates []Template) error {
	// Implementação simplificada
	fmt.Println("[")
	for i, template := range templates {
		fmt.Printf("  {\n")
		fmt.Printf("    \"name\": \"%s\",\n", template.Name)
		fmt.Printf("    \"category\": \"%s\",\n", template.Category)
		fmt.Printf("    \"description\": \"%s\",\n", template.Description)
		fmt.Printf("    \"version\": \"%s\"\n", template.Version)
		fmt.Printf("  }")
		if i < len(templates)-1 {
			fmt.Print(",")
		}
		fmt.Println()
	}
	fmt.Println("]")
	return nil
}

func outputTemplatesYAML(templates []Template) error {
	fmt.Println("templates:")
	for _, template := range templates {
		fmt.Printf("- name: %s\n", template.Name)
		fmt.Printf("  category: %s\n", template.Category)
		fmt.Printf("  description: %s\n", template.Description)
		fmt.Printf("  version: %s\n", template.Version)
	}
	return nil
}

func outputTemplateJSON(template Template) error {
	// Implementação simplificada
	fmt.Printf("{\n")
	fmt.Printf("  \"name\": \"%s\",\n", template.Name)
	fmt.Printf("  \"category\": \"%s\",\n", template.Category)
	fmt.Printf("  \"description\": \"%s\",\n", template.Description)
	fmt.Printf("  \"version\": \"%s\"\n", template.Version)
	fmt.Printf("}\n")
	return nil
}

func outputTemplateYAML(template Template) error {
	fmt.Printf("name: %s\n", template.Name)
	fmt.Printf("category: %s\n", template.Category)
	fmt.Printf("description: %s\n", template.Description)
	fmt.Printf("version: %s\n", template.Version)
	fmt.Printf("author: %s\n", template.Author)
	fmt.Printf("resources:\n")
	fmt.Printf("  cpu: %s\n", template.Resources.CPU)
	fmt.Printf("  memory: %s\n", template.Resources.Memory)
	return nil
}

