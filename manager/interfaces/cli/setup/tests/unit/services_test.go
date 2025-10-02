//go:build !integration && !e2e && !performance && !security
// +build !integration,!e2e,!performance,!security

package unit

import (
	"os"
	"path/filepath"
	"testing"

	"setup-component/tests/mocks"
)

// TestNewFileService testa a criação do serviço de arquivos
func TestNewFileService(t *testing.T) {
	tests := []struct {
		name    string
		wantErr bool
	}{
		{
			name:    "should create file service successfully",
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			service := mocks.NewFileService()
			if service == nil {
				t.Error("NewFileService() returned nil service")
			}
		})
	}
}

// TestFileService_Exists testa o método Exists
func TestFileService_Exists(t *testing.T) {
	// Criar diretório temporário para testes
	tempDir := t.TempDir()

	service := mocks.NewFileService()

	tests := []struct {
		name     string
		filePath string
		setup    bool
		want     bool
	}{
		{
			name:     "should return false when file does not exist",
			filePath: filepath.Join(tempDir, "nonexistent.txt"),
			setup:    false,
			want:     false,
		},
		{
			name:     "should return true when file exists",
			filePath: filepath.Join(tempDir, "test.txt"),
			setup:    true,
			want:     true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Criar arquivo se necessário
			if tt.setup {
				err := os.WriteFile(tt.filePath, []byte("test content"), 0644)
				if err != nil {
					t.Fatalf("Failed to create test file: %v", err)
				}
			}

			result := service.Exists(tt.filePath)
			if result != tt.want {
				t.Errorf("FileService.Exists() = %v, want %v", result, tt.want)
			}
		})
	}
}

// TestFileService_Read testa o método Read
func TestFileService_Read(t *testing.T) {
	// Criar diretório temporário para testes
	tempDir := t.TempDir()

	service := mocks.NewFileService()

	tests := []struct {
		name     string
		filePath string
		content  []byte
		setup    bool
		wantErr  bool
	}{
		{
			name:     "should fail read when file does not exist",
			filePath: filepath.Join(tempDir, "nonexistent.txt"),
			content:  nil,
			setup:    false,
			wantErr:  true,
		},
		{
			name:     "should read file successfully when file exists",
			filePath: filepath.Join(tempDir, "test.txt"),
			content:  []byte("test content"),
			setup:    true,
			wantErr:  false,
		},
		{
			name:     "should read empty file successfully",
			filePath: filepath.Join(tempDir, "empty.txt"),
			content:  []byte(""),
			setup:    true,
			wantErr:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Criar arquivo se necessário
			if tt.setup {
				err := os.WriteFile(tt.filePath, tt.content, 0644)
				if err != nil {
					t.Fatalf("Failed to create test file: %v", err)
				}
			}

			content, err := service.Read(tt.filePath)
			if (err != nil) != tt.wantErr {
				t.Errorf("FileService.Read() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr {
				if string(content) != string(tt.content) {
					t.Errorf("FileService.Read() content = %s, want %s", string(content), string(tt.content))
				}
			}
		})
	}
}

// TestFileService_Write testa o método Write
func TestFileService_Write(t *testing.T) {
	// Criar diretório temporário para testes
	tempDir := t.TempDir()

	service := mocks.NewFileService()

	tests := []struct {
		name     string
		filePath string
		content  []byte
		wantErr  bool
	}{
		{
			name:     "should write file successfully",
			filePath: filepath.Join(tempDir, "test.txt"),
			content:  []byte("test content"),
			wantErr:  false,
		},
		{
			name:     "should write file with empty content",
			filePath: filepath.Join(tempDir, "empty.txt"),
			content:  []byte(""),
			wantErr:  false,
		},
		{
			name:     "should write file with binary content",
			filePath: filepath.Join(tempDir, "binary.bin"),
			content:  []byte{0x00, 0x01, 0x02, 0x03},
			wantErr:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := service.Write(tt.filePath, tt.content)
			if (err != nil) != tt.wantErr {
				t.Errorf("FileService.Write() error = %v, wantErr %v", err, tt.wantErr)
			}

			// Verificar se o arquivo foi criado
			if !tt.wantErr {
				if _, err := os.Stat(tt.filePath); os.IsNotExist(err) {
					t.Errorf("File not created: %s", tt.filePath)
				}

				// Verificar conteúdo
				content, err := os.ReadFile(tt.filePath)
				if err != nil {
					t.Errorf("Failed to read file: %v", err)
				} else if string(content) != string(tt.content) {
					t.Errorf("File content mismatch: got %s, want %s", string(content), string(tt.content))
				}
			}
		})
	}
}

// TestFileService_Delete testa o método Delete
func TestFileService_Delete(t *testing.T) {
	// Criar diretório temporário para testes
	tempDir := t.TempDir()

	service := mocks.NewFileService()

	tests := []struct {
		name     string
		filePath string
		setup    bool
		wantErr  bool
	}{
		{
			name:     "should fail delete when file does not exist",
			filePath: filepath.Join(tempDir, "nonexistent.txt"),
			setup:    false,
			wantErr:  true,
		},
		{
			name:     "should delete file successfully when file exists",
			filePath: filepath.Join(tempDir, "test.txt"),
			setup:    true,
			wantErr:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Criar arquivo se necessário
			if tt.setup {
				err := os.WriteFile(tt.filePath, []byte("test content"), 0644)
				if err != nil {
					t.Fatalf("Failed to create test file: %v", err)
				}
			}

			err := service.Delete(tt.filePath)
			if (err != nil) != tt.wantErr {
				t.Errorf("FileService.Delete() error = %v, wantErr %v", err, tt.wantErr)
			}

			// Verificar se o arquivo foi excluído
			if !tt.wantErr {
				if _, err := os.Stat(tt.filePath); !os.IsNotExist(err) {
					t.Errorf("File not deleted: %s", tt.filePath)
				}
			}
		})
	}
}

// TestFileService_Copy testa o método Copy
func TestFileService_Copy(t *testing.T) {
	// Criar diretório temporário para testes
	tempDir := t.TempDir()

	service := mocks.NewFileService()

	tests := []struct {
		name    string
		srcPath string
		dstPath string
		content []byte
		setup   bool
		wantErr bool
	}{
		{
			name:    "should fail copy when source file does not exist",
			srcPath: filepath.Join(tempDir, "nonexistent.txt"),
			dstPath: filepath.Join(tempDir, "copy.txt"),
			content: nil,
			setup:   false,
			wantErr: true,
		},
		{
			name:    "should copy file successfully when source file exists",
			srcPath: filepath.Join(tempDir, "test.txt"),
			dstPath: filepath.Join(tempDir, "copy.txt"),
			content: []byte("test content"),
			setup:   true,
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Criar arquivo fonte se necessário
			if tt.setup {
				err := os.WriteFile(tt.srcPath, tt.content, 0644)
				if err != nil {
					t.Fatalf("Failed to create source file: %v", err)
				}
			}

			err := service.Copy(tt.srcPath, tt.dstPath)
			if (err != nil) != tt.wantErr {
				t.Errorf("FileService.Copy() error = %v, wantErr %v", err, tt.wantErr)
			}

			// Verificar se o arquivo foi copiado
			if !tt.wantErr {
				if _, err := os.Stat(tt.dstPath); os.IsNotExist(err) {
					t.Errorf("File not copied: %s", tt.dstPath)
				}

				// Verificar conteúdo
				content, err := os.ReadFile(tt.dstPath)
				if err != nil {
					t.Errorf("Failed to read copied file: %v", err)
				} else if string(content) != string(tt.content) {
					t.Errorf("Copied file content mismatch: got %s, want %s", string(content), string(tt.content))
				}
			}
		})
	}
}

// TestFileService_Move testa o método Move
func TestFileService_Move(t *testing.T) {
	// Criar diretório temporário para testes
	tempDir := t.TempDir()

	service := mocks.NewFileService()

	tests := []struct {
		name    string
		srcPath string
		dstPath string
		content []byte
		setup   bool
		wantErr bool
	}{
		{
			name:    "should fail move when source file does not exist",
			srcPath: filepath.Join(tempDir, "nonexistent.txt"),
			dstPath: filepath.Join(tempDir, "moved.txt"),
			content: nil,
			setup:   false,
			wantErr: true,
		},
		{
			name:    "should move file successfully when source file exists",
			srcPath: filepath.Join(tempDir, "test.txt"),
			dstPath: filepath.Join(tempDir, "moved.txt"),
			content: []byte("test content"),
			setup:   true,
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Criar arquivo fonte se necessário
			if tt.setup {
				err := os.WriteFile(tt.srcPath, tt.content, 0644)
				if err != nil {
					t.Fatalf("Failed to create source file: %v", err)
				}
			}

			err := service.Move(tt.srcPath, tt.dstPath)
			if (err != nil) != tt.wantErr {
				t.Errorf("FileService.Move() error = %v, wantErr %v", err, tt.wantErr)
			}

			// Verificar se o arquivo foi movido
			if !tt.wantErr {
				// Arquivo fonte não deve existir
				if _, err := os.Stat(tt.srcPath); !os.IsNotExist(err) {
					t.Errorf("Source file not moved: %s", tt.srcPath)
				}

				// Arquivo destino deve existir
				if _, err := os.Stat(tt.dstPath); os.IsNotExist(err) {
					t.Errorf("File not moved to destination: %s", tt.dstPath)
				}

				// Verificar conteúdo
				content, err := os.ReadFile(tt.dstPath)
				if err != nil {
					t.Errorf("Failed to read moved file: %v", err)
				} else if string(content) != string(tt.content) {
					t.Errorf("Moved file content mismatch: got %s, want %s", string(content), string(tt.content))
				}
			}
		})
	}
}

// TestFileService_GetSize testa o método GetSize
func TestFileService_GetSize(t *testing.T) {
	// Criar diretório temporário para testes
	tempDir := t.TempDir()

	service := mocks.NewFileService()

	tests := []struct {
		name     string
		filePath string
		content  []byte
		setup    bool
		wantErr  bool
	}{
		{
			name:     "should fail get size when file does not exist",
			filePath: filepath.Join(tempDir, "nonexistent.txt"),
			content:  nil,
			setup:    false,
			wantErr:  true,
		},
		{
			name:     "should get file size successfully when file exists",
			filePath: filepath.Join(tempDir, "test.txt"),
			content:  []byte("test content"),
			setup:    true,
			wantErr:  false,
		},
		{
			name:     "should get zero size for empty file",
			filePath: filepath.Join(tempDir, "empty.txt"),
			content:  []byte(""),
			setup:    true,
			wantErr:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Criar arquivo se necessário
			if tt.setup {
				err := os.WriteFile(tt.filePath, tt.content, 0644)
				if err != nil {
					t.Fatalf("Failed to create test file: %v", err)
				}
			}

			size, err := service.GetSize(tt.filePath)
			if (err != nil) != tt.wantErr {
				t.Errorf("FileService.GetSize() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr {
				expectedSize := int64(len(tt.content))
				if size != expectedSize {
					t.Errorf("FileService.GetSize() = %v, want %v", size, expectedSize)
				}
			}
		})
	}
}

// TestFileService_GetHash testa o método GetHash
func TestFileService_GetHash(t *testing.T) {
	// Criar diretório temporário para testes
	tempDir := t.TempDir()

	service := mocks.NewFileService()

	tests := []struct {
		name     string
		filePath string
		content  []byte
		setup    bool
		wantErr  bool
	}{
		{
			name:     "should fail get hash when file does not exist",
			filePath: filepath.Join(tempDir, "nonexistent.txt"),
			content:  nil,
			setup:    false,
			wantErr:  true,
		},
		{
			name:     "should get file hash successfully when file exists",
			filePath: filepath.Join(tempDir, "test.txt"),
			content:  []byte("test content"),
			setup:    true,
			wantErr:  false,
		},
		{
			name:     "should get hash for empty file",
			filePath: filepath.Join(tempDir, "empty.txt"),
			content:  []byte(""),
			setup:    true,
			wantErr:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Criar arquivo se necessário
			if tt.setup {
				err := os.WriteFile(tt.filePath, tt.content, 0644)
				if err != nil {
					t.Fatalf("Failed to create test file: %v", err)
				}
			}

			hash, err := service.GetHash(tt.filePath)
			if (err != nil) != tt.wantErr {
				t.Errorf("FileService.GetHash() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr {
				if hash == "" {
					t.Error("FileService.GetHash() returned empty hash")
				}
			}
		})
	}
}

// TestFileService_List testa o método List
func TestFileService_List(t *testing.T) {
	// Criar diretório temporário para testes
	tempDir := t.TempDir()

	service := mocks.NewFileService()

	tests := []struct {
		name    string
		dirPath string
		setup   bool
		wantErr bool
	}{
		{
			name:    "should fail list when directory does not exist",
			dirPath: filepath.Join(tempDir, "nonexistent"),
			setup:   false,
			wantErr: true,
		},
		{
			name:    "should list files successfully when directory exists",
			dirPath: filepath.Join(tempDir, "test"),
			setup:   true,
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Criar diretório e arquivos se necessário
			if tt.setup {
				err := os.MkdirAll(tt.dirPath, 0755)
				if err != nil {
					t.Fatalf("Failed to create test directory: %v", err)
				}

				files := []string{"file1.txt", "file2.txt", "file3.txt"}
				for _, file := range files {
					filePath := filepath.Join(tt.dirPath, file)
					err := os.WriteFile(filePath, []byte("test content"), 0644)
					if err != nil {
						t.Fatalf("Failed to create test file %s: %v", file, err)
					}
				}
			}

			files, err := service.List(tt.dirPath)
			if (err != nil) != tt.wantErr {
				t.Errorf("FileService.List() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr {
				if files == nil {
					t.Error("FileService.List() returned nil files")
				}
			}
		})
	}
}
