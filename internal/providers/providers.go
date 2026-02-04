package providers

import (
	"context"

	"github.com/javaquery/unosdk/pkg/models"
)

// Provider defines the interface for SDK providers
type Provider interface {
	// Name returns the provider name (e.g., "amazoncorretto", "openjdk")
	Name() string

	// DisplayName returns the human-readable provider name
	DisplayName() string

	// Type returns the SDK type (java, node, python, etc.)
	Type() models.SDKType

	// GetVersions returns all available versions
	GetVersions(ctx context.Context) ([]string, error)

	// GetLatestVersion returns the latest stable version
	GetLatestVersion(ctx context.Context) (string, error)

	// GetDownloadURL returns the download URL for a specific version
	GetDownloadURL(version string, arch string) (string, error)

	// GetChecksum returns the checksum for verification
	GetChecksum(version string, arch string) (string, error)

	// GetDefaultInstallPath returns the default installation path
	GetDefaultInstallPath(version string) string

	// Validate checks if the provider can handle the given version
	Validate(version string) error
}

// Registry holds all registered providers
type Registry struct {
	providers map[string]Provider
}

// NewRegistry creates a new provider registry
func NewRegistry() *Registry {
	return &Registry{
		providers: make(map[string]Provider),
	}
}

// Register adds a provider to the registry
func (r *Registry) Register(provider Provider) {
	key := string(provider.Type()) + ":" + provider.Name()
	r.providers[key] = provider
}

// Get retrieves a provider by SDK type and name
func (r *Registry) Get(sdkType models.SDKType, providerName string) (Provider, bool) {
	key := string(sdkType) + ":" + providerName
	provider, ok := r.providers[key]
	return provider, ok
}

// List returns all providers for a given SDK type
func (r *Registry) List(sdkType models.SDKType) []Provider {
	var result []Provider
	for _, provider := range r.providers {
		if provider.Type() == sdkType {
			result = append(result, provider)
		}
	}
	return result
}

// ListAll returns all registered providers
func (r *Registry) ListAll() []Provider {
	var result []Provider
	for _, provider := range r.providers {
		result = append(result, provider)
	}
	return result
}
