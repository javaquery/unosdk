package installer

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	"github.com/javaquery/unosdk/internal/providers"
	"github.com/javaquery/unosdk/pkg/models"
	"github.com/javaquery/unosdk/pkg/utils"
	"go.uber.org/zap"
)

// Installer manages SDK installations
type Installer struct {
	registry   *providers.Registry
	downloader *Downloader
	extractor  *Extractor
	logger     *zap.Logger
}

// NewInstaller creates a new Installer
func NewInstaller(registry *providers.Registry) *Installer {
	return &Installer{
		registry:   registry,
		downloader: NewDownloader(),
		extractor:  NewExtractor(),
		logger:     utils.NewLogger(),
	}
}

// Install installs an SDK
func (i *Installer) Install(ctx context.Context, sdkType models.SDKType, providerName, version, arch string) (*models.SDK, error) {
	i.logger.Info("Starting installation",
		zap.String("type", string(sdkType)),
		zap.String("provider", providerName),
		zap.String("version", version),
	)

	// Get provider
	provider, ok := i.registry.Get(sdkType, providerName)
	if !ok {
		return nil, fmt.Errorf("provider not found: %s:%s", sdkType, providerName)
	}

	// Validate version
	if err := provider.Validate(version); err != nil {
		return nil, fmt.Errorf("invalid version: %w", err)
	}

	// Handle "latest" version
	if version == "latest" {
		latest, err := provider.GetLatestVersion(ctx)
		if err != nil {
			return nil, fmt.Errorf("failed to get latest version: %w", err)
		}
		version = latest
	}

	// Get download URL
	downloadURL, err := provider.GetDownloadURL(version, arch)
	if err != nil {
		return nil, fmt.Errorf("failed to get download URL: %w", err)
	}

	// Get install path
	installPath := provider.GetDefaultInstallPath(version)

	// Check if already installed by looking for actual content
	if _, err := os.Stat(installPath); err == nil {
		// Check if directory has content (not just an empty dir left from previous uninstall)
		entries, readErr := os.ReadDir(installPath)
		if readErr == nil && len(entries) > 0 {
			i.logger.Warn("SDK already installed at path", zap.String("path", installPath))
			// Find the actual install path (might be a subdirectory)
			actualPath, _ := i.findActualInstallPath(installPath)
			return &models.SDK{
				Type:        sdkType,
				Provider:    providerName,
				Version:     version,
				InstallPath: actualPath,
				Installed:   true,
			}, nil
		}
		// If directory is empty, remove it and continue with installation
		os.RemoveAll(installPath)
	}

	// Create temporary directory for download
	tempDir, err := os.MkdirTemp("", "unosdk-*")
	if err != nil {
		return nil, fmt.Errorf("failed to create temp directory: %w", err)
	}
	defer os.RemoveAll(tempDir)

	// Download
	i.logger.Info("Downloading SDK", zap.String("url", downloadURL))
	downloadPath := filepath.Join(tempDir, filepath.Base(downloadURL))
	
	if err := i.downloader.Download(ctx, downloadURL, downloadPath); err != nil {
		return nil, fmt.Errorf("download failed: %w", err)
	}

	// Create install directory
	if err := os.MkdirAll(installPath, 0755); err != nil {
		return nil, fmt.Errorf("failed to create install directory: %w", err)
	}

	// Extract
	i.logger.Info("Extracting SDK", zap.String("path", installPath))
	if err := i.extractor.Extract(downloadPath, installPath); err != nil {
		return nil, fmt.Errorf("extraction failed: %w", err)
	}

	// Check if extraction created a single root directory
	// If so, update installPath to point to that directory
	actualInstallPath, err := i.findActualInstallPath(installPath)
	if err != nil {
		return nil, fmt.Errorf("failed to determine actual install path: %w", err)
	}

	sdk := &models.SDK{
		Type:        sdkType,
		Provider:    providerName,
		Version:     version,
		InstallPath: actualInstallPath,
		DownloadURL: downloadURL,
		Installed:   true,
	}

	i.logger.Info("Installation completed successfully", zap.String("path", actualInstallPath))
	return sdk, nil
}

// findActualInstallPath checks if the extraction created a single root directory
// and returns the path to that directory, otherwise returns the original path
func (i *Installer) findActualInstallPath(installPath string) (string, error) {
	entries, err := os.ReadDir(installPath)
	if err != nil {
		return "", err
	}

	// If there's exactly one entry and it's a directory, use that as the install path
	if len(entries) == 1 && entries[0].IsDir() {
		return filepath.Join(installPath, entries[0].Name()), nil
	}

	// Otherwise, return the original path
	return installPath, nil
}

// Uninstall removes an installed SDK
func (i *Installer) Uninstall(installPath string) error {
	i.logger.Info("Uninstalling SDK", zap.String("path", installPath))

	if _, err := os.Stat(installPath); os.IsNotExist(err) {
		return fmt.Errorf("SDK not found at path: %s", installPath)
	}

	// Remove the SDK directory
	if err := os.RemoveAll(installPath); err != nil {
		return fmt.Errorf("failed to remove SDK: %w", err)
	}

	// Also clean up parent directory if it's empty (version directory)
	// This handles cases where installPath is like: .unosdk/java/openjdk/11/jdk-11.0.21+9
	// We want to also remove the "11" directory if it's now empty
	parentDir := filepath.Dir(installPath)
	if entries, err := os.ReadDir(parentDir); err == nil && len(entries) == 0 {
		os.RemoveAll(parentDir)
		i.logger.Info("Cleaned up empty parent directory", zap.String("path", parentDir))
	}

	i.logger.Info("Uninstallation completed successfully")
	return nil
}
