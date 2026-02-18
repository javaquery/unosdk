package c

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	"github.com/javaquery/unosdk/internal/providers/mingw"
	"github.com/javaquery/unosdk/pkg/models"
)

// MinGWProvider implements the Provider interface for MinGW-w64 C toolchain
type MinGWProvider struct{}

// NewMinGWProvider creates a new MinGW-w64 provider for C
func NewMinGWProvider() *MinGWProvider {
	return &MinGWProvider{}
}

func (p *MinGWProvider) Name() string {
	return "mingw"
}

func (p *MinGWProvider) DisplayName() string {
	return "MinGW-w64 (GCC)"
}

func (p *MinGWProvider) Type() models.SDKType {
	return models.CSDK
}

func (p *MinGWProvider) GetVersions(ctx context.Context) ([]string, error) {
	return mingw.Versions(), nil
}

func (p *MinGWProvider) GetLatestVersion(ctx context.Context) (string, error) {
	return mingw.LatestVersion(), nil
}

func (p *MinGWProvider) GetDownloadURL(version string, arch string) (string, error) {
	return mingw.GetDownloadURL(version, arch)
}

func (p *MinGWProvider) GetChecksum(version string, arch string) (string, error) {
	// Checksums would typically be downloaded from the release page
	// For now, returning empty string (optional verification)
	return "", nil
}

func (p *MinGWProvider) GetDefaultInstallPath(version string) string {
	homeDir, _ := os.UserHomeDir()
	return filepath.Join(homeDir, ".unosdk", "c", "mingw", version)
}

func (p *MinGWProvider) Validate(version string) error {
	if version == "" {
		return fmt.Errorf("version cannot be empty")
	}
	return nil
}
