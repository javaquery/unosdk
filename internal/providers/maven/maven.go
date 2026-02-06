package maven

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	"github.com/javaquery/unosdk/pkg/models"
)

// MavenProvider implements the Provider interface for Apache Maven
type MavenProvider struct{}

// NewMavenProvider creates a new Maven provider
func NewMavenProvider() *MavenProvider {
	return &MavenProvider{}
}

func (p *MavenProvider) Name() string {
	return "apache"
}

func (p *MavenProvider) DisplayName() string {
	return "Apache Maven"
}

func (p *MavenProvider) Type() models.SDKType {
	return models.MavenSDK
}

func (p *MavenProvider) GetVersions(ctx context.Context) ([]string, error) {
	return []string{
		"3.9.9",
		"3.9.8",
		"3.9.7",
		"3.9.6",
		"3.8.8",
		"3.8.7",
		"3.8.6",
		"3.6.3",
	}, nil
}

func (p *MavenProvider) GetLatestVersion(ctx context.Context) (string, error) {
	return "3.9.9", nil
}

func (p *MavenProvider) GetDownloadURL(version string, arch string) (string, error) {
	// Maven is architecture-independent (pure Java application)
	// Download from Apache archive
	baseURL := "https://archive.apache.org/dist/maven/maven-3"
	
	// Extract major.minor version for URL path
	var majorMinor string
	switch version {
	case "3.9.9", "3.9.8", "3.9.7", "3.9.6":
		majorMinor = version
	case "3.8.8", "3.8.7", "3.8.6":
		majorMinor = version
	case "3.6.3":
		majorMinor = version
	default:
		return "", fmt.Errorf("unsupported version: %s", version)
	}

	fileName := fmt.Sprintf("apache-maven-%s-bin.zip", majorMinor)
	downloadURL := fmt.Sprintf("%s/%s/binaries/%s", baseURL, majorMinor, fileName)

	return downloadURL, nil
}

func (p *MavenProvider) GetChecksum(version string, arch string) (string, error) {
	// Checksums would typically be downloaded from Apache Maven's SHA512 files
	// For now, returning empty string (optional verification)
	return "", nil
}

func (p *MavenProvider) GetDefaultInstallPath(version string) string {
	homeDir, _ := os.UserHomeDir()
	return filepath.Join(homeDir, ".unosdk", "maven", version)
}

func (p *MavenProvider) Validate(version string) error {
	if version == "" {
		return fmt.Errorf("version cannot be empty")
	}
	return nil
}
