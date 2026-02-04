package node

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	"github.com/javaquery/unosdk/pkg/models"
)

// NodeJSProvider implements the Provider interface for Node.js
type NodeJSProvider struct{}

// NewNodeJSProvider creates a new Node.js provider
func NewNodeJSProvider() *NodeJSProvider {
	return &NodeJSProvider{}
}

func (p *NodeJSProvider) Name() string {
	return "nodejs"
}

func (p *NodeJSProvider) DisplayName() string {
	return "Node.js"
}

func (p *NodeJSProvider) Type() models.SDKType {
	return models.NodeSDK
}

func (p *NodeJSProvider) GetVersions(ctx context.Context) ([]string, error) {
	return []string{
		"20.10.0",
		"18.19.0",
		"16.20.2",
		"14.21.3",
	}, nil
}

func (p *NodeJSProvider) GetLatestVersion(ctx context.Context) (string, error) {
	return "20.10.0", nil
}

func (p *NodeJSProvider) GetDownloadURL(version string, arch string) (string, error) {
	// Map common architecture names to Node.js naming convention
	switch arch {
	case "", "amd64", "x86_64":
		arch = "x64"
	case "x86", "i386":
		arch = "x86"
	case "arm64", "aarch64":
		arch = "arm64"
	}

	baseURL := "https://nodejs.org/dist"
	fileName := fmt.Sprintf("node-v%s-win-%s.zip", version, arch)
	
	return fmt.Sprintf("%s/v%s/%s", baseURL, version, fileName), nil
}

func (p *NodeJSProvider) GetChecksum(version string, arch string) (string, error) {
	return "", nil
}

func (p *NodeJSProvider) GetDefaultInstallPath(version string) string {
	homeDir, _ := os.UserHomeDir()
	return filepath.Join(homeDir, ".unosdk", "node", "nodejs", version)
}

func (p *NodeJSProvider) Validate(version string) error {
	if version == "" {
		return fmt.Errorf("version cannot be empty")
	}
	return nil
}
