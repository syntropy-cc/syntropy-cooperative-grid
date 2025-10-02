//go:build !integration && !e2e && !performance && !security
// +build !integration,!e2e,!performance,!security

package unit

import (
	"os"
	"path/filepath"
	"testing"

	"setup-component/tests/mocks"
)

// TestFileExists testa a função FileExists
func TestFileExists(t *testing.T) {
	// Criar diretório temporário para testes
	tempDir := t.TempDir()

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

			result := mocks.FileExists(tt.filePath)
			if result != tt.want {
				t.Errorf("FileExists() = %v, want %v", result, tt.want)
			}
		})
	}
}

// TestDirExists testa a função DirExists
func TestDirExists(t *testing.T) {
	// Criar diretório temporário para testes
	tempDir := t.TempDir()

	tests := []struct {
		name    string
		dirPath string
		setup   bool
		want    bool
	}{
		{
			name:    "should return false when directory does not exist",
			dirPath: filepath.Join(tempDir, "nonexistent"),
			setup:   false,
			want:    false,
		},
		{
			name:    "should return true when directory exists",
			dirPath: filepath.Join(tempDir, "test"),
			setup:   true,
			want:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Criar diretório se necessário
			if tt.setup {
				err := os.MkdirAll(tt.dirPath, 0755)
				if err != nil {
					t.Fatalf("Failed to create test directory: %v", err)
				}
			}

			result := mocks.DirExists(tt.dirPath)
			if result != tt.want {
				t.Errorf("DirExists() = %v, want %v", result, tt.want)
			}
		})
	}
}

// TestCreateDir testa a função CreateDir
func TestCreateDir(t *testing.T) {
	// Criar diretório temporário para testes
	tempDir := t.TempDir()

	tests := []struct {
		name    string
		dirPath string
		wantErr bool
	}{
		{
			name:    "should create directory successfully",
			dirPath: filepath.Join(tempDir, "test"),
			wantErr: false,
		},
		{
			name:    "should create nested directory successfully",
			dirPath: filepath.Join(tempDir, "nested", "test"),
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := mocks.CreateDir(tt.dirPath)
			if (err != nil) != tt.wantErr {
				t.Errorf("CreateDir() error = %v, wantErr %v", err, tt.wantErr)
			}

			// Verificar se o diretório foi criado
			if !tt.wantErr {
				if _, err := os.Stat(tt.dirPath); os.IsNotExist(err) {
					t.Errorf("Directory not created: %s", tt.dirPath)
				}
			}
		})
	}
}

// TestWriteFile testa a função WriteFile
func TestWriteFile(t *testing.T) {
	// Criar diretório temporário para testes
	tempDir := t.TempDir()

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
			err := mocks.WriteFile(tt.filePath, tt.content)
			if (err != nil) != tt.wantErr {
				t.Errorf("WriteFile() error = %v, wantErr %v", err, tt.wantErr)
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

// TestReadFile testa a função ReadFile
func TestReadFile(t *testing.T) {
	// Criar diretório temporário para testes
	tempDir := t.TempDir()

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

			content, err := mocks.ReadFile(tt.filePath)
			if (err != nil) != tt.wantErr {
				t.Errorf("ReadFile() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr {
				if string(content) != string(tt.content) {
					t.Errorf("ReadFile() content = %s, want %s", string(content), string(tt.content))
				}
			}
		})
	}
}

// TestDeleteFile testa a função DeleteFile
func TestDeleteFile(t *testing.T) {
	// Criar diretório temporário para testes
	tempDir := t.TempDir()

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

			err := mocks.DeleteFile(tt.filePath)
			if (err != nil) != tt.wantErr {
				t.Errorf("DeleteFile() error = %v, wantErr %v", err, tt.wantErr)
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

// TestCopyFile testa a função CopyFile
func TestCopyFile(t *testing.T) {
	// Criar diretório temporário para testes
	tempDir := t.TempDir()

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

			err := mocks.CopyFile(tt.srcPath, tt.dstPath)
			if (err != nil) != tt.wantErr {
				t.Errorf("CopyFile() error = %v, wantErr %v", err, tt.wantErr)
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

// TestMoveFile testa a função MoveFile
func TestMoveFile(t *testing.T) {
	// Criar diretório temporário para testes
	tempDir := t.TempDir()

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

			err := mocks.MoveFile(tt.srcPath, tt.dstPath)
			if (err != nil) != tt.wantErr {
				t.Errorf("MoveFile() error = %v, wantErr %v", err, tt.wantErr)
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

// TestGetFileSize testa a função GetFileSize
func TestGetFileSize(t *testing.T) {
	// Criar diretório temporário para testes
	tempDir := t.TempDir()

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

			size, err := mocks.GetFileSize(tt.filePath)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetFileSize() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr {
				expectedSize := int64(len(tt.content))
				if size != expectedSize {
					t.Errorf("GetFileSize() = %v, want %v", size, expectedSize)
				}
			}
		})
	}
}

// TestGetFileHash testa a função GetFileHash
func TestGetFileHash(t *testing.T) {
	// Criar diretório temporário para testes
	tempDir := t.TempDir()

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

			hash, err := mocks.GetFileHash(tt.filePath)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetFileHash() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr {
				if hash == "" {
					t.Error("GetFileHash() returned empty hash")
				}
			}
		})
	}
}

// TestListFiles testa a função ListFiles
func TestListFiles(t *testing.T) {
	// Criar diretório temporário para testes
	tempDir := t.TempDir()

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

			files, err := mocks.ListFiles(tt.dirPath)
			if (err != nil) != tt.wantErr {
				t.Errorf("ListFiles() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr {
				if files == nil {
					t.Error("ListFiles() returned nil files")
				}
			}
		})
	}
}

// TestListDirs testa a função ListDirs
func TestListDirs(t *testing.T) {
	// Criar diretório temporário para testes
	tempDir := t.TempDir()

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
			name:    "should list directories successfully when directory exists",
			dirPath: filepath.Join(tempDir, "test"),
			setup:   true,
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Criar diretório e subdiretórios se necessário
			if tt.setup {
				err := os.MkdirAll(tt.dirPath, 0755)
				if err != nil {
					t.Fatalf("Failed to create test directory: %v", err)
				}

				dirs := []string{"dir1", "dir2", "dir3"}
				for _, dir := range dirs {
					dirPath := filepath.Join(tt.dirPath, dir)
					err := os.MkdirAll(dirPath, 0755)
					if err != nil {
						t.Fatalf("Failed to create test directory %s: %v", dir, err)
					}
				}
			}

			dirs, err := mocks.ListDirs(tt.dirPath)
			if (err != nil) != tt.wantErr {
				t.Errorf("ListDirs() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr {
				if dirs == nil {
					t.Error("ListDirs() returned nil directories")
				}
			}
		})
	}
}
