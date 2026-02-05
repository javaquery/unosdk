package java

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	"github.com/javaquery/unosdk/pkg/models"
)

// GraalVMProvider implements the Provider interface for GraalVM
type GraalVMProvider struct{}

// NewGraalVMProvider creates a new GraalVM provider
func NewGraalVMProvider() *GraalVMProvider {
	return &GraalVMProvider{}
}

func (p *GraalVMProvider) Name() string {
	return "graalvm"
}

func (p *GraalVMProvider) DisplayName() string {
	return "GraalVM"
}

func (p *GraalVMProvider) Type() models.SDKType {
	return models.JavaSDK
}

func (p *GraalVMProvider) GetVersions(ctx context.Context) ([]string, error) {
	// Latest GraalVM versions for JDK 25, 21, 17, and 11
	return []string{
		"25.0.2", // GraalVM 25.0.2 for JDK 25
		"25",     // Simplified: Latest GraalVM for JDK 25
		"21.0.2", // GraalVM 21.0.2 for JDK 21
		"21",     // Simplified: Latest GraalVM for JDK 21
		"17.0.9", // GraalVM 17.0.9 for JDK 17
		"17",     // Simplified: Latest GraalVM for JDK 17
	}, nil
}

func (p *GraalVMProvider) GetLatestVersion(ctx context.Context) (string, error) {
	return "25.0.2", nil
}

func (p *GraalVMProvider) GetDownloadURL(version string, arch string) (string, error) {
	// Map common architecture names to GraalVM's naming convention
	switch arch {
	case "", "amd64", "x86_64":
		arch = "x64"
	case "arm64", "aarch64":
		arch = "aarch64"
	default:
		return "", fmt.Errorf("unsupported architecture: %s", arch)
	}

	// GraalVM download URLs
	var graalVersion string
	
	switch version {
	case "25.0.2", "25":
		graalVersion = "25.0.0"
	case "21.0.2", "21":
		graalVersion = "21.0.2"
	case "17.0.9", "17":
		graalVersion = "17.0.9"
	default:
		return "", fmt.Errorf("unsupported version: %s", version)
	}

	// Format: graalvm-community-jdk-{graalVersion}_windows-{arch}_bin.zip
	// Example: https://github.com/graalvm/graalvm-ce-builds/releases/download/jdk-25.0.2/graalvm-community-jdk-25.0.2_windows-x64_bin.zip
	fileName := fmt.Sprintf("graalvm-community-jdk-%s_windows-%s_bin.zip", graalVersion, arch)
	
	// GraalVM uses a tag format like "jdk-25.0.2"
	releaseTag := fmt.Sprintf("jdk-%s", graalVersion)
	
	return fmt.Sprintf("https://github.com/graalvm/graalvm-ce-builds/releases/download/%s/%s", releaseTag, fileName), nil
}

func (p *GraalVMProvider) GetChecksum(version string, arch string) (string, error) {
	// In production, this would fetch the actual checksum from GraalVM releases
	return "", nil // Checksums would be fetched from provider
}

func (p *GraalVMProvider) GetDefaultInstallPath(version string) string {
	homeDir, _ := os.UserHomeDir()
	return filepath.Join(homeDir, ".unosdk", "java", "graalvm", version)
}

func (p *GraalVMProvider) Validate(version string) error {
	if version == "" {
		return fmt.Errorf("version cannot be empty")
	}
	return nil
}
