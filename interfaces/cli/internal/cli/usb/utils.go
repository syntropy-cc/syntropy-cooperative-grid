package usb

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

// formatSize formata bytes em uma string legível
func formatSize(bytes int64) string {
	const unit = 1024
	if bytes < unit {
		return fmt.Sprintf("%d B", bytes)
	}

	div, exp := int64(unit), 0
	for n := bytes / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}

	units := []string{"KB", "MB", "GB", "TB"}
	if exp >= len(units) {
		exp = len(units) - 1
	}

	return fmt.Sprintf("%.1f %s", float64(bytes)/float64(div), units[exp])
}

// outputTable exibe dispositivos em formato tabela
func outputTable(devices []USBDevice) error {
	platform := detectPlatform()

	if platform == "wsl" || platform == "windows" {
		fmt.Printf("%-15s %-8s %-30s %-15s %-10s\n",
			"DISPOSITIVO", "TAMANHO", "MODELO", "SERIAL", "PLATAFORMA")
		fmt.Println(strings.Repeat("─", 80))

		for _, device := range devices {
			fmt.Printf("%-15s %-8s %-30s %-15s %-10s\n",
				device.Path, device.Size, device.Model, device.Serial, device.Platform)
		}
	} else {
		fmt.Printf("%-12s %-8s %-20s %-15s %-10s %-10s\n",
			"DISPOSITIVO", "TAMANHO", "MODELO", "FABRICANTE", "REMOVÍVEL", "PLATAFORMA")
		fmt.Println(strings.Repeat("─", 80))

		for _, device := range devices {
			removable := "Não"
			if device.Removable {
				removable = "Sim"
			}
			fmt.Printf("%-12s %-8s %-20s %-15s %-10s %-10s\n",
				device.Path, device.Size, device.Model, device.Vendor, removable, device.Platform)
		}
	}

	return nil
}

// outputJSON exibe dispositivos em formato JSON
func outputJSON(devices []USBDevice) error {
	encoder := json.NewEncoder(os.Stdout)
	encoder.SetIndent("", "  ")
	return encoder.Encode(devices)
}

// outputYAML exibe dispositivos em formato YAML
func outputYAML(devices []USBDevice) error {
	fmt.Println("devices:")
	for _, device := range devices {
		fmt.Printf("  - path: %s\n", device.Path)
		fmt.Printf("    size: %s\n", device.Size)
		fmt.Printf("    size_gb: %d\n", device.SizeGB)
		fmt.Printf("    model: %s\n", device.Model)
		fmt.Printf("    vendor: %s\n", device.Vendor)
		fmt.Printf("    serial: %s\n", device.Serial)
		fmt.Printf("    removable: %t\n", device.Removable)
		fmt.Printf("    platform: %s\n", device.Platform)
		if device.DiskNumber > 0 {
			fmt.Printf("    disk_number: %d\n", device.DiskNumber)
		}
		if device.WindowsPath != "" {
			fmt.Printf("    windows_path: %s\n", device.WindowsPath)
		}
	}
	return nil
}

// downloadUbuntuISO baixa a ISO Ubuntu base
func downloadUbuntuISO(destPath string) error {
	// URL da ISO Ubuntu 24.04 LTS Server
	isoURL := "https://releases.ubuntu.com/24.04/ubuntu-24.04-live-server-amd64.iso"

	// Usar wget para baixar
	cmd := exec.Command("wget",
		"--progress=bar:force:noscroll",
		"--tries=3",
		"--timeout=30",
		"--continue",
		"-O", destPath,
		isoURL)

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("erro ao baixar ISO Ubuntu: %w", err)
	}

	return nil
}

// manageISOCache gerencia o download e cache de ISOs Ubuntu
func manageISOCache(cacheDir string) (string, error) {
	// Criar diretório de cache para ISOs
	isoDir := filepath.Join(cacheDir, "iso")
	if err := os.MkdirAll(isoDir, 0755); err != nil {
		return "", fmt.Errorf("erro ao criar diretório de cache: %w", err)
	}

	// Lista de URLs para tentar (do mais recente para o mais antigo)
	isoOptions := []struct {
		url      string
		filename string
		name     string
	}{
		{
			url:      "https://releases.ubuntu.com/24.04/ubuntu-24.04-live-server-amd64.iso",
			filename: "ubuntu-24.04-server.iso",
			name:     "Ubuntu 24.04 LTS Server",
		},
		{
			url:      "https://releases.ubuntu.com/22.04.5/ubuntu-22.04.5-live-server-amd64.iso",
			filename: "ubuntu-22.04.5-server-amd64.iso",
			name:     "Ubuntu 22.04.5 LTS Server",
		},
		{
			url:      "https://releases.ubuntu.com/22.04/ubuntu-22.04.4-live-server-amd64.iso",
			filename: "ubuntu-22.04-server.iso",
			name:     "Ubuntu 22.04 LTS Server",
		},
		{
			url:      "https://releases.ubuntu.com/20.04/ubuntu-20.04.6-live-server-amd64.iso",
			filename: "ubuntu-20.04-server.iso",
			name:     "Ubuntu 20.04.6 LTS Server",
		},
	}

	// Verificar se alguma ISO já existe no cache
	fmt.Println("🔍 Verificando cache de ISOs...")
	for _, iso := range isoOptions {
		isoPath := filepath.Join(isoDir, iso.filename)
		if fileInfo, err := os.Stat(isoPath); err == nil {
			// Verificar se o arquivo tem tamanho razoável (> 500MB)
			if fileInfo.Size() > 500*1024*1024 {
				fmt.Printf("✅ ISO encontrada no cache: %s\n", iso.name)
				fmt.Printf("   Arquivo: %s\n", isoPath)
				fmt.Printf("   Tamanho: %.2f GB\n", float64(fileInfo.Size())/(1024*1024*1024))
				return isoPath, nil
			} else {
				fmt.Printf("⚠️  ISO corrompida encontrada (tamanho: %d bytes), removendo...\n", fileInfo.Size())
				os.Remove(isoPath)
			}
		}
	}

	// Nenhuma ISO no cache, tentar baixar
	fmt.Println("\n📥 Nenhuma ISO encontrada no cache. Iniciando download...")
	fmt.Println("   Cache: " + isoDir)

	for _, iso := range isoOptions {
		isoPath := filepath.Join(isoDir, iso.filename)
		fmt.Printf("\n🌐 Tentando baixar: %s\n", iso.name)
		fmt.Printf("   URL: %s\n", iso.url)

		// Verificar se a URL existe antes de baixar (HEAD request)
		if !checkURLExists(iso.url) {
			fmt.Printf("   ❌ URL não disponível, tentando próxima opção...\n")
			continue
		}

		// Baixar com progresso
		if err := downloadWithProgress(iso.url, isoPath); err != nil {
			fmt.Printf("   ❌ Erro ao baixar: %v\n", err)
			// Remover arquivo parcial se existir
			os.Remove(isoPath)
			continue
		}

		// Verificar tamanho do arquivo baixado
		if fileInfo, err := os.Stat(isoPath); err == nil {
			if fileInfo.Size() > 500*1024*1024 {
				fmt.Printf("\n✅ ISO baixada com sucesso!\n")
				fmt.Printf("   Arquivo: %s\n", isoPath)
				fmt.Printf("   Tamanho: %.2f GB\n", float64(fileInfo.Size())/(1024*1024*1024))
				return isoPath, nil
			} else {
				fmt.Printf("   ❌ Arquivo baixado muito pequeno, tentando próxima opção...\n")
				os.Remove(isoPath)
			}
		}
	}

	return "", fmt.Errorf("não foi possível baixar nenhuma ISO Ubuntu. Verifique sua conexão ou forneça uma ISO com --iso")
}

// checkURLExists verifica se uma URL existe fazendo um HEAD request
func checkURLExists(url string) bool {
	cmd := exec.Command("curl", "-I", "--silent", "--head", "--fail", url)
	err := cmd.Run()
	return err == nil
}

// downloadWithProgress baixa um arquivo com indicador de progresso
func downloadWithProgress(url, destPath string) error {
	// Criar arquivo temporário
	tmpPath := destPath + ".tmp"

	// Usar wget com barra de progresso
	cmd := exec.Command("wget",
		"--progress=bar:force:noscroll",
		"--tries=3",
		"--timeout=30",
		"--continue",
		"-O", tmpPath,
		url)

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		return err
	}

	// Mover arquivo temporário para destino final
	if err := os.Rename(tmpPath, destPath); err != nil {
		return fmt.Errorf("erro ao mover arquivo: %w", err)
	}

	return nil
}
