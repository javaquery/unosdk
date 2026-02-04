package registry

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/javaquery/unosdk/pkg/models"
)

// Registry manages installed SDK tracking
type Registry struct {
	registryPath string
	sdks         map[string]*models.SDK
}

// NewRegistry creates a new SDK registry
func NewRegistry() (*Registry, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return nil, fmt.Errorf("failed to get user home directory: %w", err)
	}

	registryDir := filepath.Join(homeDir, ".unosdk")
	if err := os.MkdirAll(registryDir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create registry directory: %w", err)
	}

	registryPath := filepath.Join(registryDir, "registry.json")
	
	r := &Registry{
		registryPath: registryPath,
		sdks:         make(map[string]*models.SDK),
	}

	// Load existing registry
	if err := r.load(); err != nil {
		// If file doesn't exist, that's okay
		if !os.IsNotExist(err) {
			return nil, fmt.Errorf("failed to load registry: %w", err)
		}
	}

	return r, nil
}

// Add adds or updates an SDK in the registry
func (r *Registry) Add(sdk *models.SDK) error {
	sdk.InstalledAt = time.Now()
	sdk.Installed = true
	
	key := r.makeKey(sdk)
	r.sdks[key] = sdk

	return r.save()
}

// Remove removes an SDK from the registry
func (r *Registry) Remove(sdkType models.SDKType, provider, version string) error {
	key := fmt.Sprintf("%s:%s:%s", sdkType, provider, version)
	delete(r.sdks, key)

	return r.save()
}

// Get retrieves an SDK from the registry
func (r *Registry) Get(sdkType models.SDKType, provider, version string) (*models.SDK, bool) {
	key := fmt.Sprintf("%s:%s:%s", sdkType, provider, version)
	sdk, ok := r.sdks[key]
	return sdk, ok
}

// List returns all installed SDKs
func (r *Registry) List() []*models.SDK {
	var result []*models.SDK
	for _, sdk := range r.sdks {
		result = append(result, sdk)
	}
	return result
}

// ListByType returns all installed SDKs of a specific type
func (r *Registry) ListByType(sdkType models.SDKType) []*models.SDK {
	var result []*models.SDK
	for _, sdk := range r.sdks {
		if sdk.Type == sdkType {
			result = append(result, sdk)
		}
	}
	return result
}

// makeKey creates a unique key for an SDK
func (r *Registry) makeKey(sdk *models.SDK) string {
	return fmt.Sprintf("%s:%s:%s", sdk.Type, sdk.Provider, sdk.Version)
}

// save persists the registry to disk
func (r *Registry) save() error {
	data, err := json.MarshalIndent(r.sdks, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal registry: %w", err)
	}

	if err := os.WriteFile(r.registryPath, data, 0644); err != nil {
		return fmt.Errorf("failed to write registry: %w", err)
	}

	return nil
}

// load loads the registry from disk
func (r *Registry) load() error {
	data, err := os.ReadFile(r.registryPath)
	if err != nil {
		return err
	}

	if err := json.Unmarshal(data, &r.sdks); err != nil {
		return fmt.Errorf("failed to unmarshal registry: %w", err)
	}

	return nil
}
