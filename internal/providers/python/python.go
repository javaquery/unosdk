package python

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	"github.com/javaquery/unosdk/pkg/models"
)

// PythonProvider implements the Provider interface for Python
type PythonProvider struct{}

// NewPythonProvider creates a new Python provider
func NewPythonProvider() *PythonProvider {
	return &PythonProvider{}
}

func (p *PythonProvider) Name() string {
	return "python"
}

func (p *PythonProvider) DisplayName() string {
	return "Python"
}

func (p *PythonProvider) Type() models.SDKType {
	return models.PythonSDK
}

func (p *PythonProvider) GetVersions(ctx context.Context) ([]string, error) {
	return []string{
		"3.12.1",
		"3.11.7",
		"3.10.13",
		"3.9.18",
	}, nil
}

func (p *PythonProvider) GetLatestVersion(ctx context.Context) (string, error) {
	return "3.12.1", nil
}

func (p *PythonProvider) GetDownloadURL(version string, arch string) (string, error) {
	if arch == "" {
		arch = "amd64"
	}

	baseURL := "https://www.python.org/ftp/python"
	fileName := fmt.Sprintf("python-%s-%s.exe", version, arch)
	
	return fmt.Sprintf("%s/%s/%s", baseURL, version, fileName), nil
}

func (p *PythonProvider) GetChecksum(version string, arch string) (string, error) {
	return "", nil
}

func (p *PythonProvider) GetDefaultInstallPath(version string) string {
	homeDir, _ := os.UserHomeDir()
	return filepath.Join(homeDir, ".unosdk", "python", "python", version)
}

func (p *PythonProvider) Validate(version string) error {
	if version == "" {
		return fmt.Errorf("version cannot be empty")
	}
	return nil
}
