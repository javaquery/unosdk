package system

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// PathManager manages system PATH operations
type PathManager struct{}

// NewPathManager creates a new path manager
func NewPathManager() *PathManager {
	return &PathManager{}
}

// AddToPath adds a directory to the system PATH
func (pm *PathManager) AddToPath(dir string) error {
	env := NewWindowsEnv()
	return env.AddToPath(dir)
}

// RemoveFromPath removes a directory from the system PATH
func (pm *PathManager) RemoveFromPath(dir string) error {
	env := NewWindowsEnv()
	return env.RemoveFromPath(dir)
}

// IsInPath checks if a directory is in the PATH
func (pm *PathManager) IsInPath(dir string) bool {
	path := os.Getenv("PATH")
	pathSep := string(os.PathListSeparator)
	paths := strings.Split(path, pathSep)

	for _, p := range paths {
		if strings.EqualFold(filepath.Clean(p), filepath.Clean(dir)) {
			return true
		}
	}
	return false
}

// GetSDKInstallPath returns the default SDK installation path
func (pm *PathManager) GetSDKInstallPath() (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("failed to get user home directory: %w", err)
	}
	return filepath.Join(homeDir, ".unosdk", "sdks"), nil
}

// GetCachePath returns the cache directory path
func (pm *PathManager) GetCachePath() (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("failed to get user home directory: %w", err)
	}
	return filepath.Join(homeDir, ".unosdk", "cache"), nil
}

// GetRegistryPath returns the registry file path
func (pm *PathManager) GetRegistryPath() (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("failed to get user home directory: %w", err)
	}
	return filepath.Join(homeDir, ".unosdk", "registry.json"), nil
}

// EnsureSDKDirectories creates necessary directories
func (pm *PathManager) EnsureSDKDirectories() error {
	dirs := []func() (string, error){
		pm.GetSDKInstallPath,
		pm.GetCachePath,
	}

	for _, getDirFunc := range dirs {
		dir, err := getDirFunc()
		if err != nil {
			return err
		}
		if err := os.MkdirAll(dir, 0755); err != nil {
			return fmt.Errorf("failed to create directory %s: %w", dir, err)
		}
	}
	return nil
}
