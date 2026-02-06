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
		"1.23.5",
		"1.23.4",
		"1.23.3",
		"1.23.2",
		"1.23.1",
		"1.23.0",
		"1.22.10",
		"1.22.9",
		"1.22.8",
		"1.22.7",
		"1.22.6",
		"1.22.5",
		"1.22.4",
		"1.22.3",
		"1.22.2",
		"1.22.1",
		"1.22.0",
		"1.21.13",
		"1.21.12",
		"1.21.11",
		"1.21.10",
		"1.21.9",
		"1.21.8",
		"1.21.7",
		"1.21.6",
		"1.21.5",
		"1.21.4",
		"1.21.3",
		"1.21.2",
		"1.21.1",
		"1.21.0",
	}, nil
}

func (p *GoProvider) GetLatestVersion(ctx context.Context) (string, error) {
	return "1.23.5", nil
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
