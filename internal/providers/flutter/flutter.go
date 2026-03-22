package flutter

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	"github.com/javaquery/unosdk/pkg/models"
)

// FlutterProvider implements the Provider interface for Flutter SDK
type FlutterProvider struct{}

// NewFlutterProvider creates a new Flutter provider
func NewFlutterProvider() *FlutterProvider {
	return &FlutterProvider{}
}

func (p *FlutterProvider) Name() string {
	return "flutter"
}

func (p *FlutterProvider) DisplayName() string {
	return "Flutter SDK"
}

func (p *FlutterProvider) Type() models.SDKType {
	return models.FlutterSDK
}

func (p *FlutterProvider) GetVersions(ctx context.Context) ([]string, error) {
	// Flutter stable releases
	return []string{
		"3.41.5",
		"3.38.10",
		"3.35.7",
		"3.32.8",
		"3.29.3",
		"3.27.4",
		"3.24.5",
		"latest",
	}, nil
}

func (p *FlutterProvider) GetLatestVersion(ctx context.Context) (string, error) {
	return "3.41.5", nil
}

func (p *FlutterProvider) GetDownloadURL(version string, arch string) (string, error) {
	// Flutter uses the same package for all Windows architectures
	// The Flutter SDK is platform-specific but not architecture-specific on Windows
	
	// Handle "latest" version
	if version == "latest" {
		version = "3.41.5"
	}

	// Flutter download URLs format: https://storage.googleapis.com/flutter_infra_release/releases/stable/windows/flutter_windows_{version}-stable.zip
	baseURL := "https://storage.googleapis.com/flutter_infra_release/releases/stable/windows"
	fileName := fmt.Sprintf("flutter_windows_%s-stable.zip", version)
	
	return fmt.Sprintf("%s/%s", baseURL, fileName), nil
}

func (p *FlutterProvider) GetChecksum(version string, arch string) (string, error) {
	// In production, this would fetch the actual checksum from Flutter releases
	return "", nil
}

func (p *FlutterProvider) GetDefaultInstallPath(version string) string {
	homeDir, _ := os.UserHomeDir()
	return filepath.Join(homeDir, ".unosdk", "flutter", "flutter", version)
}

func (p *FlutterProvider) Validate(version string) error {
	if version == "" {
		return fmt.Errorf("version cannot be empty")
	}
	return nil
}
