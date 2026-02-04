package java

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	"github.com/javaquery/unosdk/pkg/models"
)

// AmazonCorrettoProvider implements the Provider interface for Amazon Corretto
type AmazonCorrettoProvider struct{}

// NewAmazonCorrettoProvider creates a new Amazon Corretto provider
func NewAmazonCorrettoProvider() *AmazonCorrettoProvider {
	return &AmazonCorrettoProvider{}
}

func (p *AmazonCorrettoProvider) Name() string {
	return "amazoncorretto"
}

func (p *AmazonCorrettoProvider) DisplayName() string {
	return "Amazon Corretto"
}

func (p *AmazonCorrettoProvider) Type() models.SDKType {
	return models.JavaSDK
}

func (p *AmazonCorrettoProvider) GetVersions(ctx context.Context) ([]string, error) {
	// In a real implementation, this would fetch from Amazon's API or website
	return []string{
		"25.0.0",
		"21.0.1",
		"17.0.9",
		"11.0.21",
		"8.392.08.1",
	}, nil
}

func (p *AmazonCorrettoProvider) GetLatestVersion(ctx context.Context) (string, error) {
	return "25.0.0", nil
}

func (p *AmazonCorrettoProvider) GetDownloadURL(version string, arch string) (string, error) {
	// Normalize architecture - Amazon Corretto uses x64, not amd64
	if arch == "" || arch == "amd64" {
		arch = "x64"
	} else if arch == "arm64" {
		arch = "aarch64"
	}

	// Map version to download URL
	// Amazon Corretto uses major version in the filename
	baseURL := "https://corretto.aws/downloads/latest"
	
	var majorVersion string
	switch version {
	case "25.0.0", "25":
		majorVersion = "25"
	case "21.0.1", "21":
		majorVersion = "21"
	case "17.0.9", "17":
		majorVersion = "17"
	case "11.0.21", "11":
		majorVersion = "11"
	case "8.392.08.1", "8":
		majorVersion = "8"
	default:
		return "", fmt.Errorf("unsupported version: %s", version)
	}

	// Format: amazon-corretto-{major}-{arch}-windows-jdk.zip
	fileName := fmt.Sprintf("amazon-corretto-%s-%s-windows-jdk.zip", majorVersion, arch)
	return fmt.Sprintf("%s/%s", baseURL, fileName), nil
}

func (p *AmazonCorrettoProvider) GetChecksum(version string, arch string) (string, error) {
	// In production, this would fetch the actual checksum from Amazon
	return "", nil // Checksums would be fetched from provider
}

func (p *AmazonCorrettoProvider) GetDefaultInstallPath(version string) string {
	homeDir, _ := os.UserHomeDir()
	return filepath.Join(homeDir, ".unosdk", "java", "amazoncorretto", version)
}

func (p *AmazonCorrettoProvider) Validate(version string) error {
	// Basic validation
	if version == "" {
		return fmt.Errorf("version cannot be empty")
	}
	return nil
}
