//go:build !integration && !e2e && !performance && !security
// +build !integration,!e2e,!performance,!security

package unit

import (
	"os"
	"path/filepath"
	"testing"

	setup "setup-component/src"
)

// TestNewConfigurator testa a criação do configurador
func TestNewConfigurator(t *testing.T) {
	// Criar diretório temporário para testes
	tempDir := t.TempDir()
	originalHome := os.Getenv("HOME")
	os.Setenv("HOME", tempDir)
	defer os.Setenv("HOME", originalHome)

	logger := setup.NewSetupLogger()
	defer logger.Close()

	tests := []struct {
		name    string
		logger  *setup.SetupLogger
		wantErr bool
	}{
		{
			name:    "should create configurator successfully",
			logger:  logger,
			wantErr: false,
		},
		{
			name:    "should create configurator with nil logger",
			logger:  nil,
			wantErr: false, // Logger pode ser nil
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			configurator := setup.NewConfigurator(tt.logger)
			if configurator == nil {
				t.Error("NewConfigurator() returned nil configurator")
			}
		})
	}
}

// TestConfigurator_GenerateConfig testa a geração de configuração
func TestConfigurator_GenerateConfig(t *testing.T) {
	// Criar diretório temporário para testes
	tempDir := t.TempDir()
	originalHome := os.Getenv("HOME")
	os.Setenv("HOME", tempDir)
	defer os.Setenv("HOME", originalHome)

	logger := setup.NewSetupLogger()
	defer logger.Close()

	manager, err := setup.NewSetupManager()
	if err != nil {
		t.Fatalf("Failed to create setup manager: %v", err)
	}

	tests := []struct {
		name    string
		options *setup.SetupOptions
		wantErr bool
	}{
		{
			name: "should generate config successfully with valid options",
			options: &setup.SetupOptions{
				Force:          true,
				SkipValidation: true,
				CustomSettings: map[string]string{
					"owner_name":  "Test User",
					"owner_email": "test@example.com",
				},
			},
			wantErr: false,
		},
		{
			name: "should generate config successfully with empty options",
			options: &setup.SetupOptions{
				Force:          true,
				SkipValidation: true,
			},
			wantErr: false,
		},
		{
			name:    "should fail with nil options",
			options: nil,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := manager.SetupWithPublicOptions(tt.options)
			if (err != nil) != tt.wantErr {
				t.Errorf("SetupManager.Setup() error = %v, wantErr %v", err, tt.wantErr)
			}

			// Verificar se o arquivo de configuração foi criado
			if !tt.wantErr {
				configPath := filepath.Join(tempDir, ".syntropy", "config", "manager.yaml")
				if _, err := os.Stat(configPath); os.IsNotExist(err) {
					t.Errorf("Config file not created: %s", configPath)
				}
			}
		})
	}
}

// TestConfigurator_CreateStructure testa a criação da estrutura de diretórios
func TestConfigurator_CreateStructure(t *testing.T) {
	// Criar diretório temporário para testes
	tempDir := t.TempDir()
	originalHome := os.Getenv("HOME")
	os.Setenv("HOME", tempDir)
	defer os.Setenv("HOME", originalHome)

	logger := setup.NewSetupLogger()
	defer logger.Close()

	configurator := setup.NewConfigurator(logger)

	tests := []struct {
		name    string
		wantErr bool
	}{
		{
			name:    "should create structure successfully",
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := configurator.CreateStructure()
			if (err != nil) != tt.wantErr {
				t.Errorf("Configurator.CreateStructure() error = %v, wantErr %v", err, tt.wantErr)
			}

			// Verificar se os diretórios foram criados
			if !tt.wantErr {
				directories := []string{
					"config",
					"keys",
					"nodes",
					"logs",
					"cache",
					"backups",
					"templates",
					"state",
				}

				for _, dir := range directories {
					dirPath := filepath.Join(tempDir, ".syntropy", dir)
					if _, err := os.Stat(dirPath); os.IsNotExist(err) {
						t.Errorf("Directory not created: %s", dirPath)
					}
				}
			}
		})
	}
}

// TestConfigurator_GenerateKeys testa a geração de chaves
func TestConfigurator_GenerateKeys(t *testing.T) {
	// Criar diretório temporário para testes
	tempDir := t.TempDir()
	originalHome := os.Getenv("HOME")
	os.Setenv("HOME", tempDir)
	defer os.Setenv("HOME", originalHome)

	logger := setup.NewSetupLogger()
	defer logger.Close()

	configurator := setup.NewConfigurator(logger)

	tests := []struct {
		name    string
		wantErr bool
	}{
		{
			name:    "should generate keys successfully",
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			keyPair, err := configurator.GenerateKeys()
			if (err != nil) != tt.wantErr {
				t.Errorf("Configurator.GenerateKeys() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr {
				if keyPair == nil {
					t.Error("Configurator.GenerateKeys() returned nil key pair")
					return
				}

				// Verificar campos obrigatórios
				if keyPair.ID == "" {
					t.Error("Key pair missing ID")
				}
				if keyPair.Algorithm == "" {
					t.Error("Key pair missing algorithm")
				}
				if keyPair.PrivateKey == nil {
					t.Error("Key pair missing private key")
				}
				if keyPair.PublicKey == nil {
					t.Error("Key pair missing public key")
				}
				if keyPair.Fingerprint == "" {
					t.Error("Key pair missing fingerprint")
				}
			}
		})
	}
}

// TestConfigurator_ValidateConfig testa a validação de configuração
func TestConfigurator_ValidateConfig(t *testing.T) {
	// Criar diretório temporário para testes
	tempDir := t.TempDir()
	originalHome := os.Getenv("HOME")
	os.Setenv("HOME", tempDir)
	defer os.Setenv("HOME", originalHome)

	logger := setup.NewSetupLogger()
	defer logger.Close()

	configurator := setup.NewConfigurator(logger)

	tests := []struct {
		name    string
		setup   bool
		wantErr bool
	}{
		{
			name:    "should fail validation when no config exists",
			setup:   false,
			wantErr: true,
		},
		{
			name:    "should validate config successfully when config exists",
			setup:   true,
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Criar configuração se necessário
			if tt.setup {
				options := &setup.SetupOptions{
					Force:          true,
					SkipValidation: true,
					CustomSettings: map[string]string{
						"owner_name":  "Test User",
						"owner_email": "test@example.com",
					},
				}
				manager, err := setup.NewSetupManager()
				if err != nil {
					t.Fatalf("Failed to create setup manager: %v", err)
				}
				err = manager.SetupWithPublicOptions(options)
				if err != nil {
					t.Fatalf("Failed to generate config: %v", err)
				}
			}

			err := configurator.ValidateConfig()
			if (err != nil) != tt.wantErr {
				t.Errorf("Configurator.ValidateConfig() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

// TestConfigurator_BackupConfig testa o backup de configuração
func TestConfigurator_BackupConfig(t *testing.T) {
	// Criar diretório temporário para testes
	tempDir := t.TempDir()
	originalHome := os.Getenv("HOME")
	os.Setenv("HOME", tempDir)
	defer os.Setenv("HOME", originalHome)

	logger := setup.NewSetupLogger()
	defer logger.Close()

	configurator := setup.NewConfigurator(logger)

	tests := []struct {
		name    string
		setup   bool
		wantErr bool
	}{
		{
			name:    "should fail backup when no config exists",
			setup:   false,
			wantErr: true,
		},
		{
			name:    "should backup config successfully when config exists",
			setup:   true,
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Criar configuração se necessário
			if tt.setup {
				options := &setup.SetupOptions{
					Force:          true,
					SkipValidation: true,
					CustomSettings: map[string]string{
						"owner_name":  "Test User",
						"owner_email": "test@example.com",
					},
				}
				manager, err := setup.NewSetupManager()
				if err != nil {
					t.Fatalf("Failed to create setup manager: %v", err)
				}
				err = manager.SetupWithPublicOptions(options)
				if err != nil {
					t.Fatalf("Failed to generate config: %v", err)
				}
			}

			err := configurator.BackupConfig("test_backup")
			if (err != nil) != tt.wantErr {
				t.Errorf("Configurator.BackupConfig() error = %v, wantErr %v", err, tt.wantErr)
			}

			// Verificar se o backup foi criado
			if !tt.wantErr {
				backupDir := filepath.Join(tempDir, ".syntropy", "backups")
				if _, err := os.Stat(backupDir); os.IsNotExist(err) {
					t.Errorf("Backup directory not created: %s", backupDir)
				}
			}
		})
	}
}

// TestConfigurator_RestoreConfig testa a restauração de configuração
func TestConfigurator_RestoreConfig(t *testing.T) {
	// Criar diretório temporário para testes
	tempDir := t.TempDir()
	originalHome := os.Getenv("HOME")
	os.Setenv("HOME", tempDir)
	defer os.Setenv("HOME", originalHome)

	logger := setup.NewSetupLogger()
	defer logger.Close()

	configurator := setup.NewConfigurator(logger)

	tests := []struct {
		name    string
		setup   bool
		wantErr bool
	}{
		{
			name:    "should fail restore when no backup exists",
			setup:   false,
			wantErr: true,
		},
		{
			name:    "should restore config successfully when backup exists",
			setup:   true,
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var backupPath string

			// Criar backup se necessário
			if tt.setup {
				options := &setup.SetupOptions{
					Force:          true,
					SkipValidation: true,
					CustomSettings: map[string]string{
						"owner_name":  "Test User",
						"owner_email": "test@example.com",
					},
				}
				manager, err := setup.NewSetupManager()
				if err != nil {
					t.Fatalf("Failed to create setup manager: %v", err)
				}
				err = manager.SetupWithPublicOptions(options)
				if err != nil {
					t.Fatalf("Failed to generate config: %v", err)
				}

				err = configurator.BackupConfig("test_backup")
				if err != nil {
					t.Fatalf("Failed to backup config: %v", err)
				}

				// Encontrar o arquivo de backup
				backupDir := filepath.Join(tempDir, ".syntropy", "backups")
				files, err := os.ReadDir(backupDir)
				if err != nil {
					t.Fatalf("Failed to read backup directory: %v", err)
				}

				for _, file := range files {
					if filepath.Ext(file.Name()) == ".yaml" {
						backupPath = filepath.Join(backupDir, file.Name())
						break
					}
				}

				if backupPath == "" {
					t.Fatalf("Backup file not found")
				}
			} else {
				backupPath = filepath.Join(tempDir, "nonexistent_backup.yaml")
			}

			err := configurator.RestoreConfig(backupPath)
			if (err != nil) != tt.wantErr {
				t.Errorf("Configurator.RestoreConfig() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

// TestConfigurator_LoadTemplate testa o carregamento de template
func TestConfigurator_LoadTemplate(t *testing.T) {
	// Criar diretório temporário para testes
	tempDir := t.TempDir()
	originalHome := os.Getenv("HOME")
	os.Setenv("HOME", tempDir)
	defer os.Setenv("HOME", originalHome)

	logger := setup.NewSetupLogger()
	defer logger.Close()

	configurator := setup.NewConfigurator(logger)

	tests := []struct {
		name         string
		templateName string
		setup        bool
		wantErr      bool
	}{
		{
			name:         "should fail load when template does not exist",
			templateName: "nonexistent.yaml",
			setup:        false,
			wantErr:      true,
		},
		{
			name:         "should load template successfully when template exists",
			templateName: "test.yaml",
			setup:        true,
			wantErr:      false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Criar template se necessário
			if tt.setup {
				templatePath := filepath.Join(tempDir, ".syntropy", "templates", tt.templateName)
				templateContent := "test: template content"
				err := os.WriteFile(templatePath, []byte(templateContent), 0644)
				if err != nil {
					t.Fatalf("Failed to create template: %v", err)
				}
			}

			content, err := configurator.LoadTemplate(tt.templateName)
			if (err != nil) != tt.wantErr {
				t.Errorf("Configurator.LoadTemplate() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr {
				if content == "" {
					t.Error("Configurator.LoadTemplate() returned empty content")
				}
			}
		})
	}
}

// TestConfigurator_SaveTemplate testa o salvamento de template
func TestConfigurator_SaveTemplate(t *testing.T) {
	// Criar diretório temporário para testes
	tempDir := t.TempDir()
	originalHome := os.Getenv("HOME")
	os.Setenv("HOME", tempDir)
	defer os.Setenv("HOME", originalHome)

	logger := setup.NewSetupLogger()
	defer logger.Close()

	configurator := setup.NewConfigurator(logger)

	tests := []struct {
		name         string
		templateName string
		content      string
		wantErr      bool
	}{
		{
			name:         "should save template successfully",
			templateName: "test.yaml",
			content:      "test: template content",
			wantErr:      false,
		},
		{
			name:         "should save template with empty content",
			templateName: "empty.yaml",
			content:      "",
			wantErr:      false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := configurator.SaveTemplate(tt.templateName, tt.content)
			if (err != nil) != tt.wantErr {
				t.Errorf("Configurator.SaveTemplate() error = %v, wantErr %v", err, tt.wantErr)
			}

			// Verificar se o template foi salvo
			if !tt.wantErr {
				templatePath := filepath.Join(tempDir, ".syntropy", "templates", tt.templateName)
				if _, err := os.Stat(templatePath); os.IsNotExist(err) {
					t.Errorf("Template not saved: %s", templatePath)
				}

				// Verificar conteúdo
				content, err := os.ReadFile(templatePath)
				if err != nil {
					t.Errorf("Failed to read template: %v", err)
				} else if string(content) != tt.content {
					t.Errorf("Template content mismatch: got %s, want %s", string(content), tt.content)
				}
			}
		})
	}
}

// TestConfigurator_ProcessTemplate testa o processamento de template
func TestConfigurator_ProcessTemplate(t *testing.T) {
	// Criar diretório temporário para testes
	tempDir := t.TempDir()
	originalHome := os.Getenv("HOME")
	os.Setenv("HOME", tempDir)
	defer os.Setenv("HOME", originalHome)

	logger := setup.NewSetupLogger()
	defer logger.Close()

	configurator := setup.NewConfigurator(logger)

	tests := []struct {
		name         string
		templateName string
		variables    map[string]interface{}
		setup        bool
		wantErr      bool
	}{
		{
			name:         "should fail process when template does not exist",
			templateName: "nonexistent.yaml",
			variables:    map[string]interface{}{},
			setup:        false,
			wantErr:      true,
		},
		{
			name:         "should process template successfully when template exists",
			templateName: "test.yaml",
			variables: map[string]interface{}{
				"name":  "Test User",
				"email": "test@example.com",
			},
			setup:   true,
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Criar template se necessário
			if tt.setup {
				templatePath := filepath.Join(tempDir, ".syntropy", "templates", tt.templateName)
				templateContent := "name: {{.name}}\nemail: {{.email}}"
				err := os.WriteFile(templatePath, []byte(templateContent), 0644)
				if err != nil {
					t.Fatalf("Failed to create template: %v", err)
				}
			}

			content, err := configurator.ProcessTemplate(tt.templateName, tt.variables)
			if (err != nil) != tt.wantErr {
				t.Errorf("Configurator.ProcessTemplate() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr {
				if content == "" {
					t.Error("Configurator.ProcessTemplate() returned empty content")
				}
			}
		})
	}
}

// TestConfigurator_GetConfigPath testa a obtenção do caminho de configuração
func TestConfigurator_GetConfigPath(t *testing.T) {
	// Criar diretório temporário para testes
	tempDir := t.TempDir()
	originalHome := os.Getenv("HOME")
	os.Setenv("HOME", tempDir)
	defer os.Setenv("HOME", originalHome)

	logger := setup.NewSetupLogger()
	defer logger.Close()

	configurator := setup.NewConfigurator(logger)

	tests := []struct {
		name string
		want string
	}{
		{
			name: "should return config path",
			want: filepath.Join(tempDir, ".syntropy", "config", "manager.yaml"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := configurator.GetConfigPath()
			if result != tt.want {
				t.Errorf("Configurator.GetConfigPath() = %v, want %v", result, tt.want)
			}
		})
	}
}

// TestConfigurator_GetTemplatesDir testa a obtenção do diretório de templates
func TestConfigurator_GetTemplatesDir(t *testing.T) {
	// Criar diretório temporário para testes
	tempDir := t.TempDir()
	originalHome := os.Getenv("HOME")
	os.Setenv("HOME", tempDir)
	defer os.Setenv("HOME", originalHome)

	logger := setup.NewSetupLogger()
	defer logger.Close()

	configurator := setup.NewConfigurator(logger)

	tests := []struct {
		name string
		want string
	}{
		{
			name: "should return templates directory path",
			want: filepath.Join(tempDir, ".syntropy", "templates"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := configurator.GetTemplatesDir()
			if result != tt.want {
				t.Errorf("Configurator.GetTemplatesDir() = %v, want %v", result, tt.want)
			}
		})
	}
}

// TestConfigurator_ListTemplates testa a listagem de templates
func TestConfigurator_ListTemplates(t *testing.T) {
	// Criar diretório temporário para testes
	tempDir := t.TempDir()
	originalHome := os.Getenv("HOME")
	os.Setenv("HOME", tempDir)
	defer os.Setenv("HOME", originalHome)

	logger := setup.NewSetupLogger()
	defer logger.Close()

	configurator := setup.NewConfigurator(logger)

	tests := []struct {
		name    string
		setup   bool
		wantErr bool
	}{
		{
			name:    "should list templates successfully when templates exist",
			setup:   true,
			wantErr: false,
		},
		{
			name:    "should list templates successfully when no templates exist",
			setup:   false,
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Criar templates se necessário
			if tt.setup {
				templates := []string{"test1.yaml", "test2.yml", "test3.txt"}
				for _, template := range templates {
					templatePath := filepath.Join(tempDir, ".syntropy", "templates", template)
					err := os.WriteFile(templatePath, []byte("test content"), 0644)
					if err != nil {
						t.Fatalf("Failed to create template %s: %v", template, err)
					}
				}
			}

			templates, err := configurator.ListTemplates()
			if (err != nil) != tt.wantErr {
				t.Errorf("Configurator.ListTemplates() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr {
				if templates == nil {
					t.Error("Configurator.ListTemplates() returned nil templates")
				}
			}
		})
	}
}
