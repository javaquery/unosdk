package gradle

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	"github.com/javaquery/unosdk/pkg/models"
)

// GradleProvider implements the Provider interface for Gradle
type GradleProvider struct{}

// NewGradleProvider creates a new Gradle provider
func NewGradleProvider() *GradleProvider {
	return &GradleProvider{}
}

func (p *GradleProvider) Name() string {
	return "gradle"
}

func (p *GradleProvider) DisplayName() string {
	return "Gradle"
}

func (p *GradleProvider) Type() models.SDKType {
	return models.GradleSDK
}

func (p *GradleProvider) GetVersions(ctx context.Context) ([]string, error) {
	return []string{
		"8.12",
		"8.11.1",
		"8.11",
		"8.10.2",
		"8.10.1",
		"8.10",
		"8.9",
		"8.8",
		"8.7",
		"8.6",
		"8.5",
		"8.4",
		"8.3",
		"8.2.1",
		"8.2",
		"8.1.1",
		"8.1",
		"8.0.2",
		"8.0.1",
		"8.0",
		"7.6.4",
		"7.6.3",
		"7.6.2",
		"7.6.1",
		"7.6",
	}, nil
}

func (p *GradleProvider) GetLatestVersion(ctx context.Context) (string, error) {
	return "8.12", nil
}

func (p *GradleProvider) GetDownloadURL(version string, arch string) (string, error) {
	// Gradle is architecture-independent (pure Java application)
	// Download from official Gradle distribution site
	baseURL := "https://services.gradle.org/distributions"
	fileName := fmt.Sprintf("gradle-%s-bin.zip", version)
	downloadURL := fmt.Sprintf("%s/%s", baseURL, fileName)

	return downloadURL, nil
}

func (p *GradleProvider) GetChecksum(version string, arch string) (string, error) {
	// Checksums would typically be downloaded from Gradle's SHA256 files
	// For now, returning empty string (optional verification)
	return "", nil
}

func (p *GradleProvider) GetDefaultInstallPath(version string) string {
	homeDir, _ := os.UserHomeDir()
	return filepath.Join(homeDir, ".unosdk", "gradle", version)
}

func (p *GradleProvider) Validate(version string) error {
	if version == "" {
		return fmt.Errorf("version cannot be empty")
	}
	return nil
}
