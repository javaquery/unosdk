package java

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	"github.com/javaquery/unosdk/pkg/models"
)

// OpenJDKProvider implements the Provider interface for OpenJDK
type OpenJDKProvider struct{}

// NewOpenJDKProvider creates a new OpenJDK provider
func NewOpenJDKProvider() *OpenJDKProvider {
	return &OpenJDKProvider{}
}

func (p *OpenJDKProvider) Name() string {
	return "openjdk"
}

func (p *OpenJDKProvider) DisplayName() string {
	return "OpenJDK"
}

func (p *OpenJDKProvider) Type() models.SDKType {
	return models.JavaSDK
}

func (p *OpenJDKProvider) GetVersions(ctx context.Context) ([]string, error) {
	return []string{
		"25.0.0",
		"21.0.1",
		"17.0.9",
		"11.0.21",
		"8u392",
	}, nil
}

func (p *OpenJDKProvider) GetLatestVersion(ctx context.Context) (string, error) {
	return "25.0.0", nil
}

func (p *OpenJDKProvider) GetDownloadURL(version string, arch string) (string, error) {
	// Map common architecture names to Adoptium's naming convention
	switch arch {
	case "", "amd64", "x86_64":
		arch = "x64"
	case "x86", "i386":
		arch = "x86-32"
	case "arm64", "aarch64":
		arch = "aarch64"
	}

	// OpenJDK download URLs (using Adoptium/Eclipse Temurin)
	var repoName, releaseTag, fileName string
	
	switch version {
	case "25.0.0", "25":
		repoName = "temurin25-binaries"
		releaseTag = "jdk-25+36"
		fileName = fmt.Sprintf("OpenJDK25U-jdk_%s_windows_hotspot_25_36.zip", arch)
	case "21.0.1", "21":
		repoName = "temurin21-binaries"
		releaseTag = "jdk-21.0.1+12"
		fileName = fmt.Sprintf("OpenJDK21U-jdk_%s_windows_hotspot_21.0.1_12.zip", arch)
	case "17.0.9", "17":
		repoName = "temurin17-binaries"
		releaseTag = "jdk-17.0.9+9"
		fileName = fmt.Sprintf("OpenJDK17U-jdk_%s_windows_hotspot_17.0.9_9.zip", arch)
	case "11.0.21", "11":
		repoName = "temurin11-binaries"
		releaseTag = "jdk-11.0.21+9"
		fileName = fmt.Sprintf("OpenJDK11U-jdk_%s_windows_hotspot_11.0.21_9.zip", arch)
	case "8u392", "8":
		repoName = "temurin8-binaries"
		releaseTag = "jdk8u392-b08"
		fileName = fmt.Sprintf("OpenJDK8U-jdk_%s_windows_hotspot_8u392b08.zip", arch)
	default:
		return "", fmt.Errorf("unsupported version: %s", version)
	}

	return fmt.Sprintf("https://github.com/adoptium/%s/releases/download/%s/%s", repoName, releaseTag, fileName), nil
}

func (p *OpenJDKProvider) GetChecksum(version string, arch string) (string, error) {
	return "", nil
}

func (p *OpenJDKProvider) GetDefaultInstallPath(version string) string {
	homeDir, _ := os.UserHomeDir()
	return filepath.Join(homeDir, ".unosdk", "java", "openjdk", version)
}

func (p *OpenJDKProvider) Validate(version string) error {
	if version == "" {
		return fmt.Errorf("version cannot be empty")
	}
	return nil
}
