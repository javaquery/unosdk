package config

import (
	"os"
	"path/filepath"
)

// Config holds the configuration for the SDK manager
type Config struct {
	ConfigDir   string
	CacheDir    string
	InstallDir  string
	RegistryURL string
}

// New creates a new configuration with default values
func New() (*Config, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return nil, err
	}

	configDir := filepath.Join(homeDir, ".unosdk")
	cacheDir := filepath.Join(configDir, "cache")
	installDir := filepath.Join(configDir, "sdks")

	return &Config{
		ConfigDir:   configDir,
		CacheDir:    cacheDir,
		InstallDir:  installDir,
		RegistryURL: DefaultRegistryURL,
	}, nil
}

// EnsureDirectories creates necessary directories if they don't exist
func (c *Config) EnsureDirectories() error {
	dirs := []string{
		c.ConfigDir,
		c.CacheDir,
		c.InstallDir,
	}

	for _, dir := range dirs {
		if err := os.MkdirAll(dir, 0755); err != nil {
			return err
		}
	}

	return nil
}
