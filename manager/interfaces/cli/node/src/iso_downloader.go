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
	"time"

	"node-component/src/internal/constants"
	"node-component/src/internal/helpers"
	"node-component/src/internal/types"
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
	// Create cache directory
	cacheDir := filepath.Join(os.Getenv("HOME"), ".syntropy", "cache", "isos")
	if err := os.MkdirAll(cacheDir, 0755); err != nil {
		logger.Error("Failed to create cache directory", "error", err)
	}

	// Create HTTP client with timeout
	httpClient := &http.Client{
		Timeout: constants.DefaultHTTPTimeout,
	}

	return &ISODownloaderImpl{
		cacheDir:   cacheDir,
		logger:     logger,
		httpClient: httpClient,
	}
}

// SetProgressCallback sets the progress callback for downloads
func (id *ISODownloaderImpl) SetProgressCallback(callback func(downloaded, total int64)) {
	id.progressCB = callback
}

// DownloadISO downloads Ubuntu Server ISO with progress tracking
func (id *ISODownloaderImpl) DownloadISO(ctx context.Context, version string) (*types.ISOInfo, error) {
	id.logger.Info("Starting ISO download", "version", version)

	// Get ISO information
	isoVersion, err := id.GetISOInfo(version)
	if err != nil {
		return nil, fmt.Errorf("failed to get ISO info: %w", err)
	}

	// Check if already cached
	if cachedISO, err := id.GetCachedISO(version); err == nil && cachedISO != nil {
		id.logger.Info("Using cached ISO", "version", version, "path", cachedISO.FilePath)
		return cachedISO, nil
	}

	// Download ISO
	isoPath := filepath.Join(id.cacheDir, isoVersion.FileName)
	if err := id.downloadFile(ctx, isoVersion.DownloadURL, isoPath); err != nil {
		return nil, fmt.Errorf("failed to download ISO: %w", err)
	}

	// Validate ISO
	if err := id.ValidateISO(isoPath, isoVersion.SHA256); err != nil {
		os.Remove(isoPath) // Clean up invalid file
		return nil, fmt.Errorf("ISO validation failed: %w", err)
	}

	// Get file info
	fileInfo, err := os.Stat(isoPath)
	if err != nil {
		return nil, fmt.Errorf("failed to get file info: %w", err)
	}

	isoInfo := &types.ISOInfo{
		Version:      version,
		FilePath:     isoPath,
		FileName:     isoVersion.FileName,
		Size:         fileInfo.Size(),
		SHA256:       isoVersion.SHA256,
		DownloadURL:  isoVersion.DownloadURL,
		DownloadedAt: time.Now(),
	}

	id.logger.Info("ISO download completed", "version", version, "size", fileInfo.Size())
	return isoInfo, nil
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
			FileName:    "ubuntu-24.04-server-amd64.iso",
			DownloadURL: "https://releases.ubuntu.com/24.04/ubuntu-24.04-server-amd64.iso",
			SHA256:      constants.Ubuntu2404ServerSHA256,
			Size:        int64(constants.Ubuntu2404ServerSize),
			ReleaseDate: time.Date(2024, 4, 25, 0, 0, 0, 0, time.UTC),
		},
		{
			Version:     "22.04",
			LTS:         true,
			FileName:    "ubuntu-22.04-server-amd64.iso",
			DownloadURL: "https://releases.ubuntu.com/22.04/ubuntu-22.04-server-amd64.iso",
			SHA256:      constants.Ubuntu2204ServerSHA256,
			Size:        int64(constants.Ubuntu2204ServerSize),
			ReleaseDate: time.Date(2022, 4, 21, 0, 0, 0, 0, time.UTC),
		},
		{
			Version:     "20.04",
			LTS:         true,
			FileName:    "ubuntu-20.04-server-amd64.iso",
			DownloadURL: "https://releases.ubuntu.com/20.04/ubuntu-20.04-server-amd64.iso",
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

// downloadFile downloads a file with progress tracking
func (id *ISODownloaderImpl) downloadFile(ctx context.Context, url string, filePath string) error {
	id.logger.Debug("Starting download", "url", url, "path", filePath)

	// Create request with context
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

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
	logger types.Logger
	start  time.Time
}

// NewISOProgressTracker creates a new progress tracker
func NewISOProgressTracker(logger types.Logger) *ISOProgressTracker {
	return &ISOProgressTracker{
		logger: logger,
		start:  time.Now(),
	}
}

// TrackProgress implements progress callback for ISO downloads
func (ipt *ISOProgressTracker) TrackProgress(downloaded, total int64) {
	if total <= 0 {
		return
	}

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

	// Log progress every 5%
	if int(percentage)%5 == 0 {
		ipt.logger.Info("Download progress",
			"progress", fmt.Sprintf("%.1f%%", percentage),
			"downloaded", fmt.Sprintf("%.1fMB", downloadedMB),
			"total", fmt.Sprintf("%.1fMB", totalMB),
			"speed", fmt.Sprintf("%.1fMB/s", speed),
			"eta", eta.Truncate(time.Second))
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
