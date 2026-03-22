package golang

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	"github.com/javaquery/unosdk/pkg/models"
)

// GoProvider implements the Provider interface for Go
type GoProvider struct{}

// NewGoProvider creates a new Go provider
func NewGoProvider() *GoProvider {
	return &GoProvider{}
}

func (p *GoProvider) Name() string {
	return "golang"
}

func (p *GoProvider) DisplayName() string {
	return "Go"
}

func (p *GoProvider) Type() models.SDKType {
	return models.GoSDK
}

func (p *GoProvider) GetVersions(ctx context.Context) ([]string, error) {
	return []string{
		"1.26.1",
		"1.26.0",
		"1.25.8",
		"1.25.7",
		"1.25.6",
		"1.25.5",
		"1.25.4",
		"1.25.3",
		"1.25.2",
		"1.25.1",
		"1.25.0",
		"1.24.13",
		"1.24.12",
		"1.24.11",
		"1.24.10",
		"1.24.9",
		"1.24.8",
		"1.24.7",
		"1.24.6",
		"1.24.5",
		"1.24.4",
		"1.24.3",
		"1.24.2",
		"1.24.1",
		"1.24.0",
	}, nil
}

func (p *GoProvider) GetLatestVersion(ctx context.Context) (string, error) {
	return "1.26.1", nil
}

func (p *GoProvider) GetDownloadURL(version string, arch string) (string, error) {
	// Map architecture to Go's naming convention
	goArch := arch
	switch arch {
	case "x64", "amd64":
		goArch = "amd64"
	case "x86", "386":
		goArch = "386"
	case "arm64":
		goArch = "arm64"
	default:
		goArch = "amd64"
	}

	// Go download URL format: https://go.dev/dl/go{version}.{os}-{arch}.zip
	baseURL := "https://go.dev/dl"
	fileName := fmt.Sprintf("go%s.windows-%s.zip", version, goArch)
	downloadURL := fmt.Sprintf("%s/%s", baseURL, fileName)

	return downloadURL, nil
}

func (p *GoProvider) GetChecksum(version string, arch string) (string, error) {
	// Checksums would typically be downloaded from Go's SHA256 files
	// For now, returning empty string (optional verification)
	return "", nil
}

func (p *GoProvider) GetDefaultInstallPath(version string) string {
	homeDir, _ := os.UserHomeDir()
	return filepath.Join(homeDir, ".unosdk", "go", "golang", version)
}

func (p *GoProvider) Validate(version string) error {
	if version == "" {
		return fmt.Errorf("version cannot be empty")
	}
	return nil
}
