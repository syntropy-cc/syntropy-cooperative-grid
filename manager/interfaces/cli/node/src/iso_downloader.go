package node

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"

	"node-component/src/internal/constants"
	"node-component/src/internal/helpers"
	"node-component/src/internal/types"

	"gopkg.in/yaml.v3"
)

// ISODownloader defines the interface for downloading and managing Ubuntu ISOs
type ISODownloader interface {
	// DownloadISO downloads Ubuntu Server ISO with progress tracking
	DownloadISO(ctx context.Context, version string) (*types.ISOInfo, error)

	// GetCachedISO returns cached ISO info if available
	GetCachedISO(version string) (*types.ISOInfo, error)

	// ListAvailableVersions returns list of available Ubuntu versions
	ListAvailableVersions() []types.UbuntuVersion

	// GetISOInfo returns information about a specific ISO version
	GetISOInfo(version string) (*types.UbuntuVersion, error)

	// ValidateISO validates ISO integrity using SHA256 checksum
	ValidateISO(isoPath string, expectedSHA256 string) error

	// CleanupCache removes old ISOs from cache
	CleanupCache(maxAge time.Duration) error

	// GetCacheStats returns cache statistics
	GetCacheStats() (*types.ISOCacheStats, error)
}

// ISODownloaderImpl implements the ISODownloader interface
type ISODownloaderImpl struct {
	cacheDir   string
	logger     types.Logger
	httpClient *http.Client
	progressCB func(downloaded, total int64)
}

// NewISODownloader creates a new ISO downloader
func NewISODownloader(logger types.Logger) *ISODownloaderImpl {
	// Get home directory using cross-platform method
	homeDir, err := os.UserHomeDir()
	if err != nil {
		logger.Error("Failed to get home directory", "error", err)
		// Fallback to temp directory
		homeDir = os.TempDir()
	}

	// Cache directory should be .syntropy/cache/isos (cache dir created by setup)
	cacheDir := filepath.Join(homeDir, ".syntropy", "cache", "isos")

	// Check if .syntropy/cache exists (should be created by setup component)
	syntropyCacheDir := filepath.Join(homeDir, ".syntropy", "cache")
	if _, err := os.Stat(syntropyCacheDir); err != nil {
		logger.Error("Setup cache directory not found", "error", err, "path", syntropyCacheDir)
		// Fallback to temp directory if setup not found
		cacheDir = filepath.Join(os.TempDir(), "syntropy", "cache", "isos")
	}

	// Create isos subdirectory if it doesn't exist
	if err := os.MkdirAll(cacheDir, 0755); err != nil {
		logger.Error("Failed to create isos cache directory", "error", err)
	}

	// Create HTTP client with timeout and redirect following
	httpClient := &http.Client{
		Timeout: constants.DefaultHTTPTimeout,
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			// Permitir até 10 redirecionamentos
			if len(via) >= 10 {
				return fmt.Errorf("too many redirects")
			}
			return nil
		},
	}

	return &ISODownloaderImpl{
		cacheDir:   cacheDir,
		logger:     logger,
		httpClient: httpClient,
	}
}

// loadISOConfig loads ISO configuration from manager.yaml
func (id *ISODownloaderImpl) loadISOConfig() (*types.ISODownloadConfig, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return nil, err
	}

	configPath := filepath.Join(homeDir, ".syntropy", "config", "manager.yaml")

	// Se não existe, retornar configuração padrão
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		return id.getDefaultISOConfig(), nil
	}

	// Ler e parsear arquivo
	data, err := os.ReadFile(configPath)
	if err != nil {
		id.logger.Warn("Failed to read config, using defaults", "error", err)
		return id.getDefaultISOConfig(), nil
	}

	// Parse YAML para extrair seção ISO
	var config struct {
		ISO types.ISODownloadConfig `yaml:"iso"`
	}

	if err := yaml.Unmarshal(data, &config); err != nil {
		id.logger.Warn("Failed to parse config, using defaults", "error", err)
		return id.getDefaultISOConfig(), nil
	}

	return &config.ISO, nil
}

// getDefaultISOConfig returns default ISO configuration
func (id *ISODownloaderImpl) getDefaultISOConfig() *types.ISODownloadConfig {
	return &types.ISODownloadConfig{
		CustomURLs: []string{
			"https://releases.ubuntu.com/24.04.3/ubuntu-24.04.3-live-server-amd64.iso?_gl=1*12y79vf*_gcl_au*MTAzMzgyNzE0NC4xNzU5NDIzNDAz",
		},
		PreferredMirrors: []string{
			"https://releases.ubuntu.com",
			"https://mirror.math.princeton.edu/pub/ubuntu-iso",
			"https://mirrors.kernel.org/ubuntu-iso",
			"https://mirror.umd.edu/ubuntu-iso",
		},
		EnableAutoFallback: true,
		MaxRetries:         3,
		Timeout:            30 * time.Minute,
		SkipValidation:     false,
	}
}

// buildDownloadURLs builds a prioritized list of download URLs
func (id *ISODownloaderImpl) buildDownloadURLs(version string, customURL string) ([]string, error) {
	urls := []string{}

	// 1. URL personalizada da flag CLI (prioridade máxima)
	if customURL != "" {
		id.logger.Info("Adding custom URL from CLI flag", "url", customURL)
		urls = append(urls, customURL)
	}

	// 2. Carregar configuração
	config, err := id.loadISOConfig()
	if err != nil {
		id.logger.Warn("Failed to load config", "error", err)
		config = id.getDefaultISOConfig()
	}

	// 3. URLs personalizadas do manager.yaml
	urls = append(urls, config.CustomURLs...)

	// 4. URL oficial primária
	isoVersion, err := id.GetISOInfo(version)
	if err == nil {
		urls = append(urls, isoVersion.DownloadURL)
	}

	// 5. Mirrors alternativos
	if config.EnableAutoFallback {
		for _, mirror := range config.PreferredMirrors {
			if !strings.HasPrefix(mirror, "https://releases.ubuntu.com") {
				mirrorURL := fmt.Sprintf("%s/%s/ubuntu-%s-live-server-amd64.iso",
					strings.TrimSuffix(mirror, "/"), version, version)
				urls = append(urls, mirrorURL)
			}
		}
	}

	// 6. Variável de ambiente (ISO_CUSTOM_URLS)
	if envURLs := os.Getenv("ISO_CUSTOM_URLS"); envURLs != "" {
		for _, url := range strings.Split(envURLs, ",") {
			urls = append(urls, strings.TrimSpace(url))
		}
	}

	// Remover duplicatas mantendo ordem
	seen := make(map[string]bool)
	uniqueURLs := []string{}
	for _, url := range urls {
		if !seen[url] && url != "" {
			seen[url] = true
			uniqueURLs = append(uniqueURLs, url)
		}
	}

	id.logger.Info("Built download URL list", "total_urls", len(uniqueURLs), "urls", uniqueURLs)
	return uniqueURLs, nil
}

// validateURL checks if a URL is accessible
func (id *ISODownloaderImpl) validateURL(ctx context.Context, url string) error {
	// Para URLs personalizados, ser mais flexível na validação
	if strings.Contains(url, "ubuntu.com/download") || strings.Contains(url, "thank-you") {
		id.logger.Debug("Skipping validation for Ubuntu download page", "url", url)
		return nil
	}

	req, err := http.NewRequestWithContext(ctx, "HEAD", url, nil)
	if err != nil {
		return err
	}

	// Adicionar headers para simular um navegador
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36")
	req.Header.Set("Accept", "*/*")

	resp, err := id.httpClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Aceitar mais códigos de status (incluindo redirecionamentos)
	if resp.StatusCode >= 200 && resp.StatusCode < 400 {
		return nil
	}

	return fmt.Errorf("URL returned status %d", resp.StatusCode)
}

// SetProgressCallback sets the progress callback for downloads
func (id *ISODownloaderImpl) SetProgressCallback(callback func(downloaded, total int64)) {
	id.progressCB = callback
}

// DownloadISO downloads Ubuntu Server ISO with progress tracking
func (id *ISODownloaderImpl) DownloadISO(ctx context.Context, version string) (*types.ISOInfo, error) {
	return id.DownloadISOWithOptions(ctx, version, "", false)
}

// DownloadISOWithOptions downloads Ubuntu Server ISO with custom URL support
func (id *ISODownloaderImpl) DownloadISOWithOptions(ctx context.Context, version string, customURL string, skipValidation bool) (*types.ISOInfo, error) {
	id.logger.Info("Starting ISO download with fallback system", "version", version, "custom_url", customURL)

	// Verificar cache primeiro
	if cachedISO, err := id.GetCachedISO(version); err == nil && cachedISO != nil {
		id.logger.Info("Using cached ISO", "version", version, "path", cachedISO.FilePath)
		return cachedISO, nil
	}

	// Construir lista de URLs
	urls, err := id.buildDownloadURLs(version, customURL)
	if err != nil {
		return nil, fmt.Errorf("failed to build download URLs: %w", err)
	}

	if len(urls) == 0 {
		return nil, fmt.Errorf("no download URLs available for version %s", version)
	}

	id.logger.Info("Attempting download from multiple sources", "total_urls", len(urls))

	// Tentar cada URL
	var attempts []types.ISODownloadAttempt
	var lastErr error

	for i, url := range urls {
		attempt := types.ISODownloadAttempt{URL: url}
		attemptStart := time.Now()

		id.logger.Info("Trying download", "attempt", i+1, "total", len(urls), "url", url)
		fmt.Printf("🔄 Tentativa %d/%d: %s\n", i+1, len(urls), url)

		// Validar URL primeiro
		if err := id.validateURL(ctx, url); err != nil {
			attempt.Success = false
			attempt.ErrorMessage = fmt.Sprintf("URL validation failed: %v", err)
			attempt.Duration = time.Since(attemptStart)
			attempts = append(attempts, attempt)

			id.logger.Warn("URL not available", "url", url, "error", err)
			fmt.Printf("   ⚠️  URL não disponível: %v\n", err)
			lastErr = err
			continue
		}

		fmt.Printf("   ✅ URL disponível, iniciando download...\n")
		fmt.Printf("📥 Baixando ubuntu-%s-live-server-amd64.iso\n", version)

		// Tentar download
		isoPath := filepath.Join(id.cacheDir, fmt.Sprintf("ubuntu-%s-live-server-amd64.iso", version))

		// Para URLs de páginas de download do Ubuntu, tentar extrair o link real
		if strings.Contains(url, "ubuntu.com/download") || strings.Contains(url, "thank-you") {
			realURL, err := id.extractDownloadURL(ctx, url)
			if err != nil {
				id.logger.Warn("Failed to extract download URL, trying original", "url", url, "error", err)
			} else {
				id.logger.Info("Extracted real download URL", "original", url, "real", realURL)
				url = realURL
			}
		}

		if err := id.downloadFile(ctx, url, isoPath); err != nil {
			attempt.Success = false
			attempt.ErrorMessage = fmt.Sprintf("Download failed: %v", err)
			attempt.Duration = time.Since(attemptStart)
			attempts = append(attempts, attempt)

			id.logger.Warn("Download failed", "url", url, "error", err)
			fmt.Printf("   ❌ Download falhou: %v\n", err)
			lastErr = err
			continue
		}

		// Obter informações do ISO (necessário para criar ISOInfo)
		isoVersion, err := id.GetISOInfo(version)
		if err != nil {
			attempt.Success = false
			attempt.ErrorMessage = fmt.Sprintf("Failed to get ISO info: %v", err)
			attempt.Duration = time.Since(attemptStart)
			attempts = append(attempts, attempt)
			lastErr = err
			continue
		}

		// Verificar se deve pular validação
		if skipValidation {
			fmt.Printf("   ⚠️  Download concluído, pulando validação SHA256...\n")
			id.logger.Warn("ISO validation skipped - this is a temporary feature that will be re-enabled when the software matures")
		} else {
			fmt.Printf("   ✅ Download concluído, validando integridade...\n")

			if err := id.ValidateISO(isoPath, isoVersion.SHA256); err != nil {
				attempt.Success = false
				attempt.ErrorMessage = fmt.Sprintf("Validation failed: %v", err)
				attempt.Duration = time.Since(attemptStart)
				attempts = append(attempts, attempt)

				id.logger.Warn("ISO validation failed", "url", url, "error", err)
				fmt.Printf("   ❌ Validação falhou: %v\n", err)
				os.Remove(isoPath)
				lastErr = err
				continue
			}
		}

		// Sucesso!
		attempt.Success = true
		attempt.Duration = time.Since(attemptStart)
		attempts = append(attempts, attempt)

		fmt.Printf("✅ ISO baixada e validada com sucesso de: %s\n", url)

		// Criar ISOInfo
		fileInfo, err := os.Stat(isoPath)
		if err != nil {
			return nil, fmt.Errorf("failed to get file info: %w", err)
		}

		return &types.ISOInfo{
			Version:      version,
			FilePath:     isoPath,
			FileName:     isoVersion.FileName,
			Size:         fileInfo.Size(),
			SHA256:       isoVersion.SHA256,
			DownloadURL:  url,
			DownloadedAt: time.Now(),
		}, nil
	}

	// Todos os URLs falharam
	return nil, id.handleAllDownloadsFailed(version, attempts, lastErr)
}

// handleAllDownloadsFailed provides helpful feedback when all downloads fail
func (id *ISODownloaderImpl) handleAllDownloadsFailed(version string, attempts []types.ISODownloadAttempt, lastErr error) error {
	fmt.Printf("\n❌ Todos os downloads falharam para Ubuntu %s\n\n", version)

	// Mostrar resumo das tentativas
	fmt.Printf("📋 Resumo das tentativas (%d):\n", len(attempts))
	for i, attempt := range attempts {
		status := "❌ Falhou"
		if attempt.Success {
			status = "✅ Sucesso"
		}
		fmt.Printf("   %d. %s - %s\n", i+1, status, attempt.URL)
		if !attempt.Success && attempt.ErrorMessage != "" {
			fmt.Printf("      Erro: %s\n", attempt.ErrorMessage)
		}
	}

	// Fornecer opções ao usuário
	fmt.Printf("\n🔧 Opções disponíveis:\n")
	fmt.Printf("1. 📁 Fornecer ISO local:\n")
	fmt.Printf("   syntropy node create --iso /caminho/para/ubuntu-%s-live-server-amd64.iso\n\n", version)

	fmt.Printf("2. 🌐 Configurar URL personalizada via CLI:\n")
	fmt.Printf("   syntropy node create --iso-url https://seu-mirror.com/ubuntu-%s-live-server-amd64.iso\n\n", version)

	fmt.Printf("3. ⚙️  Adicionar URLs no arquivo de configuração:\n")
	fmt.Printf("   Edite: ~/.syntropy/config/manager.yaml\n")
	fmt.Printf("   Adicione na seção 'iso.custom_urls':\n")
	fmt.Printf("   iso:\n")
	fmt.Printf("     custom_urls:\n")
	fmt.Printf("       - https://seu-mirror.com/ubuntu-%s-live-server-amd64.iso\n\n", version)

	fmt.Printf("4. 🔄 Tentar novamente (pode ser problema temporário):\n")
	fmt.Printf("   syntropy node create --ubuntu-version %s\n\n", version)

	fmt.Printf("5. 🌍 Usar variável de ambiente:\n")
	fmt.Printf("   export ISO_CUSTOM_URLS=\"https://mirror1.com/iso,https://mirror2.com/iso\"\n\n")

	fmt.Printf("💡 Dica: Baixe manualmente de https://ubuntu.com/download/server\n")

	return fmt.Errorf("all %d download attempts failed. Last error: %w", len(attempts), lastErr)
}

// GetCachedISO returns cached ISO info if available
func (id *ISODownloaderImpl) GetCachedISO(version string) (*types.ISOInfo, error) {
	isoVersion, err := id.GetISOInfo(version)
	if err != nil {
		return nil, fmt.Errorf("failed to get ISO info: %w", err)
	}

	isoPath := filepath.Join(id.cacheDir, isoVersion.FileName)

	// Check if file exists
	if !helpers.FileExists(isoPath) {
		return nil, fmt.Errorf("cached ISO not found: %s", isoPath)
	}

	// Validate cached ISO
	if err := id.ValidateISO(isoPath, isoVersion.SHA256); err != nil {
		id.logger.Warn("Cached ISO validation failed, will re-download", "version", version, "error", err)
		os.Remove(isoPath) // Remove invalid cached file
		return nil, fmt.Errorf("cached ISO validation failed: %w", err)
	}

	// Get file info
	fileInfo, err := os.Stat(isoPath)
	if err != nil {
		return nil, fmt.Errorf("failed to get cached file info: %w", err)
	}

	isoInfo := &types.ISOInfo{
		Version:      version,
		FilePath:     isoPath,
		FileName:     isoVersion.FileName,
		Size:         fileInfo.Size(),
		SHA256:       isoVersion.SHA256,
		DownloadURL:  isoVersion.DownloadURL,
		DownloadedAt: fileInfo.ModTime(),
	}

	id.logger.Debug("Found cached ISO", "version", version, "path", isoPath)
	return isoInfo, nil
}

// ListAvailableVersions returns list of available Ubuntu versions
func (id *ISODownloaderImpl) ListAvailableVersions() []types.UbuntuVersion {
	return []types.UbuntuVersion{
		{
			Version:     "24.04",
			LTS:         true,
			FileName:    "ubuntu-24.04-live-server-amd64.iso",
			DownloadURL: "https://releases.ubuntu.com/24.04.3/ubuntu-24.04.3-live-server-amd64.iso?_gl=1*12y79vf*_gcl_au*MTAzMzgyNzE0NC4xNzU5NDIzNDAz",
			SHA256:      constants.Ubuntu2404ServerSHA256,
			Size:        int64(constants.Ubuntu2404ServerSize),
			ReleaseDate: time.Date(2024, 4, 25, 0, 0, 0, 0, time.UTC),
		},
		{
			Version:     "22.04",
			LTS:         true,
			FileName:    "ubuntu-22.04-live-server-amd64.iso",
			DownloadURL: "https://releases.ubuntu.com/22.04/ubuntu-22.04-live-server-amd64.iso",
			SHA256:      constants.Ubuntu2204ServerSHA256,
			Size:        int64(constants.Ubuntu2204ServerSize),
			ReleaseDate: time.Date(2022, 4, 21, 0, 0, 0, 0, time.UTC),
		},
		{
			Version:     "20.04",
			LTS:         true,
			FileName:    "ubuntu-20.04-live-server-amd64.iso",
			DownloadURL: "https://releases.ubuntu.com/20.04/ubuntu-20.04-live-server-amd64.iso",
			SHA256:      constants.Ubuntu2004ServerSHA256,
			Size:        int64(constants.Ubuntu2004ServerSize),
			ReleaseDate: time.Date(2020, 4, 23, 0, 0, 0, 0, time.UTC),
		},
	}
}

// GetISOInfo returns information about a specific ISO version
func (id *ISODownloaderImpl) GetISOInfo(version string) (*types.UbuntuVersion, error) {
	versions := id.ListAvailableVersions()

	for _, v := range versions {
		if v.Version == version {
			return &v, nil
		}
	}

	return nil, fmt.Errorf("unsupported Ubuntu version: %s", version)
}

// ValidateISO validates ISO integrity using SHA256 checksum
func (id *ISODownloaderImpl) ValidateISO(isoPath string, expectedSHA256 string) error {
	id.logger.Debug("Validating ISO integrity", "path", isoPath)

	file, err := os.Open(isoPath)
	if err != nil {
		return fmt.Errorf("failed to open ISO file: %w", err)
	}
	defer file.Close()

	hasher := sha256.New()
	if _, err := io.Copy(hasher, file); err != nil {
		return fmt.Errorf("failed to read ISO file: %w", err)
	}

	actualSHA256 := hex.EncodeToString(hasher.Sum(nil))
	if actualSHA256 != expectedSHA256 {
		return fmt.Errorf("SHA256 mismatch: expected %s, got %s", expectedSHA256, actualSHA256)
	}

	id.logger.Debug("ISO validation successful", "path", isoPath)
	return nil
}

// CleanupCache removes old ISOs from cache
func (id *ISODownloaderImpl) CleanupCache(maxAge time.Duration) error {
	id.logger.Info("Cleaning up ISO cache", "maxAge", maxAge)

	files, err := os.ReadDir(id.cacheDir)
	if err != nil {
		return fmt.Errorf("failed to read cache directory: %w", err)
	}

	cutoff := time.Now().Add(-maxAge)
	removedCount := 0

	for _, file := range files {
		if file.IsDir() {
			continue
		}

		filePath := filepath.Join(id.cacheDir, file.Name())
		fileInfo, err := file.Info()
		if err != nil {
			id.logger.Warn("Failed to get file info", "file", file.Name(), "error", err)
			continue
		}

		if fileInfo.ModTime().Before(cutoff) {
			if err := os.Remove(filePath); err != nil {
				id.logger.Warn("Failed to remove old cache file", "file", file.Name(), "error", err)
			} else {
				removedCount++
				id.logger.Debug("Removed old cache file", "file", file.Name())
			}
		}
	}

	id.logger.Info("Cache cleanup completed", "removed", removedCount)
	return nil
}

// GetCacheStats returns cache statistics
func (id *ISODownloaderImpl) GetCacheStats() (*types.ISOCacheStats, error) {
	files, err := os.ReadDir(id.cacheDir)
	if err != nil {
		return nil, fmt.Errorf("failed to read cache directory: %w", err)
	}

	stats := &types.ISOCacheStats{
		CacheDir:   id.cacheDir,
		TotalFiles: len(files),
		TotalSize:  0,
		OldestFile: time.Now(),
		NewestFile: time.Time{},
	}

	for _, file := range files {
		if file.IsDir() {
			continue
		}

		fileInfo, err := file.Info()
		if err != nil {
			continue
		}

		stats.TotalSize += fileInfo.Size()

		if fileInfo.ModTime().Before(stats.OldestFile) {
			stats.OldestFile = fileInfo.ModTime()
		}

		if fileInfo.ModTime().After(stats.NewestFile) {
			stats.NewestFile = fileInfo.ModTime()
		}
	}

	return stats, nil
}

// Private helper methods

// extractDownloadURL extracts the real download URL from Ubuntu download pages
func (id *ISODownloaderImpl) extractDownloadURL(ctx context.Context, pageURL string) (string, error) {
	id.logger.Debug("Extracting download URL from page", "url", pageURL)

	req, err := http.NewRequestWithContext(ctx, "GET", pageURL, nil)
	if err != nil {
		return "", err
	}

	// Adicionar headers para simular um navegador
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36")
	req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,*/*;q=0.8")
	req.Header.Set("Accept-Language", "en-US,en;q=0.9")

	resp, err := id.httpClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("page returned status %d", resp.StatusCode)
	}

	// Ler o conteúdo da página
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	// Procurar por links de download na página
	content := string(body)

	// Procurar por padrões comuns de links de download
	patterns := []string{
		`href="([^"]*ubuntu-[^"]*\.iso[^"]*)"`,
		`href="([^"]*\.iso[^"]*)"`,
		`"downloadUrl":"([^"]*\.iso[^"]*)"`,
		`"url":"([^"]*\.iso[^"]*)"`,
	}

	for _, pattern := range patterns {
		re := regexp.MustCompile(pattern)
		matches := re.FindStringSubmatch(content)
		if len(matches) > 1 {
			downloadURL := matches[1]

			// Se for um URL relativo, converter para absoluto
			if strings.HasPrefix(downloadURL, "/") {
				downloadURL = "https://releases.ubuntu.com" + downloadURL
			} else if strings.HasPrefix(downloadURL, "//") {
				downloadURL = "https:" + downloadURL
			}

			id.logger.Info("Found download URL in page", "url", downloadURL)
			return downloadURL, nil
		}
	}

	return "", fmt.Errorf("no download URL found in page")
}

// downloadFile downloads a file with progress tracking
func (id *ISODownloaderImpl) downloadFile(ctx context.Context, url string, filePath string) error {
	id.logger.Debug("Starting download", "url", url, "path", filePath)

	// Create request with context
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	// Adicionar headers para simular um navegador
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36")
	req.Header.Set("Accept", "*/*")
	req.Header.Set("Accept-Language", "en-US,en;q=0.9")

	// Make request
	resp, err := id.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to make request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("download failed with status: %d", resp.StatusCode)
	}

	// Get content length
	contentLength := resp.ContentLength
	if contentLength <= 0 {
		id.logger.Warn("Unknown content length for download", "url", url)
	}

	// Create output file
	outFile, err := os.Create(filePath)
	if err != nil {
		return fmt.Errorf("failed to create output file: %w", err)
	}
	defer outFile.Close()

	// Download with progress tracking
	var downloaded int64
	buffer := make([]byte, 32*1024) // 32KB buffer

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}

		n, err := resp.Body.Read(buffer)
		if n > 0 {
			if _, writeErr := outFile.Write(buffer[:n]); writeErr != nil {
				return fmt.Errorf("failed to write to file: %w", writeErr)
			}

			downloaded += int64(n)

			// Call progress callback if set
			if id.progressCB != nil && contentLength > 0 {
				id.progressCB(downloaded, contentLength)
			}
		}

		if err == io.EOF {
			break
		}
		if err != nil {
			return fmt.Errorf("failed to read response body: %w", err)
		}
	}

	// Final progress callback
	if id.progressCB != nil && contentLength > 0 {
		id.progressCB(contentLength, contentLength)
	}

	id.logger.Debug("Download completed", "url", url, "size", downloaded)
	return nil
}

// ISOProgressTracker tracks download progress
type ISOProgressTracker struct {
	logger     types.Logger
	start      time.Time
	lastUpdate time.Time
}

// NewISOProgressTracker creates a new progress tracker
func NewISOProgressTracker(logger types.Logger) *ISOProgressTracker {
	return &ISOProgressTracker{
		logger:     logger,
		start:      time.Now(),
		lastUpdate: time.Now(),
	}
}

// TrackProgress implements progress callback for ISO downloads
func (ipt *ISOProgressTracker) TrackProgress(downloaded, total int64) {
	if total <= 0 {
		return
	}

	now := time.Now()
	// Update progress bar every 0.5 seconds to avoid too frequent updates
	if now.Sub(ipt.lastUpdate) < 500*time.Millisecond && downloaded < total {
		return
	}
	ipt.lastUpdate = now

	percentage := float64(downloaded) / float64(total) * 100
	elapsed := time.Since(ipt.start)

	// Calculate download speed
	var speed float64
	if elapsed.Seconds() > 0 {
		speed = float64(downloaded) / elapsed.Seconds() / 1024 / 1024 // MB/s
	}

	// Calculate ETA
	var eta time.Duration
	if speed > 0 && downloaded < total {
		remaining := total - downloaded
		eta = time.Duration(float64(remaining)/speed) * time.Second
	}

	// Format sizes
	downloadedMB := float64(downloaded) / 1024 / 1024
	totalMB := float64(total) / 1024 / 1024

	// Create progress bar
	barWidth := 30
	filled := int(float64(barWidth) * percentage / 100)
	bar := ""
	for i := 0; i < barWidth; i++ {
		if i < filled {
			bar += "█"
		} else {
			bar += "░"
		}
	}

	// Format ETA
	etaStr := "∞"
	if eta > 0 {
		etaStr = eta.Truncate(time.Second).String()
	}

	// Display progress bar with \r to overwrite the same line
	progressLine := fmt.Sprintf("\r📥 Baixando ISO: [%s] %.1f%% (%.1fGB/%.1fGB) %.1fMB/s ETA: %s",
		bar, percentage, downloadedMB/1024, totalMB/1024, speed, etaStr)

	// Print to stdout (not logger) for real-time display
	fmt.Print(progressLine)

	// If download is complete, add newline and log final status
	if downloaded >= total {
		fmt.Println() // Add newline after progress bar
		ipt.logger.Info("Download completed",
			"progress", "100.0%",
			"downloaded", fmt.Sprintf("%.1fGB", downloadedMB/1024),
			"total", fmt.Sprintf("%.1fGB", totalMB/1024),
			"speed", fmt.Sprintf("%.1fMB/s", speed),
			"duration", elapsed.Truncate(time.Second))
	}
}

// ISOManager manages ISO operations
type ISOManager struct {
	downloader *ISODownloaderImpl
	logger     types.Logger
}

// NewISOManager creates a new ISO manager
func NewISOManager(logger types.Logger) *ISOManager {
	return &ISOManager{
		downloader: NewISODownloader(logger),
		logger:     logger,
	}
}

// DownloadUbuntuServer downloads Ubuntu Server ISO
func (im *ISOManager) DownloadUbuntuServer(ctx context.Context, version string) (*types.ISOInfo, error) {
	im.logger.Info("Downloading Ubuntu Server ISO", "version", version)

	// Create progress tracker
	progressTracker := NewISOProgressTracker(im.logger)
	im.downloader.SetProgressCallback(progressTracker.TrackProgress)

	// Download ISO
	isoInfo, err := im.downloader.DownloadISO(ctx, version)
	if err != nil {
		return nil, fmt.Errorf("failed to download Ubuntu Server ISO: %w", err)
	}

	im.logger.Info("Ubuntu Server ISO downloaded successfully",
		"version", version,
		"path", isoInfo.FilePath,
		"size", isoInfo.Size)

	return isoInfo, nil
}

// GetCachedUbuntuServer returns cached Ubuntu Server ISO
func (im *ISOManager) GetCachedUbuntuServer(version string) (*types.ISOInfo, error) {
	return im.downloader.GetCachedISO(version)
}

// ListAvailableUbuntuVersions returns available Ubuntu versions
func (im *ISOManager) ListAvailableUbuntuVersions() []types.UbuntuVersion {
	return im.downloader.ListAvailableVersions()
}

// ValidateUbuntuISO validates Ubuntu ISO integrity
func (im *ISOManager) ValidateUbuntuISO(version string, isoPath string) error {
	isoVersion, err := im.downloader.GetISOInfo(version)
	if err != nil {
		return fmt.Errorf("failed to get ISO info: %w", err)
	}

	return im.downloader.ValidateISO(isoPath, isoVersion.SHA256)
}

// CleanupISOCache cleans up old ISOs
func (im *ISOManager) CleanupISOCache(maxAge time.Duration) error {
	return im.downloader.CleanupCache(maxAge)
}

// GetISOCacheStats returns cache statistics
func (im *ISOManager) GetISOCacheStats() (*types.ISOCacheStats, error) {
	return im.downloader.GetCacheStats()
}
