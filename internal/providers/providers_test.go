package providers

import (
	"context"
	"testing"

	"github.com/javaquery/unosdk/pkg/models"
)

// mockProvider is a mock implementation of the Provider interface for testing
type mockProvider struct {
	name        string
	displayName string
	sdkType     models.SDKType
	versions    []string
	latestVer   string
}

func (m *mockProvider) Name() string {
	return m.name
}

func (m *mockProvider) DisplayName() string {
	return m.displayName
}

func (m *mockProvider) Type() models.SDKType {
	return m.sdkType
}

func (m *mockProvider) GetVersions(ctx context.Context) ([]string, error) {
	return m.versions, nil
}

func (m *mockProvider) GetLatestVersion(ctx context.Context) (string, error) {
	return m.latestVer, nil
}

func (m *mockProvider) GetDownloadURL(version string, arch string) (string, error) {
	return "https://example.com/" + version, nil
}

func (m *mockProvider) GetChecksum(version string, arch string) (string, error) {
	return "checksum123", nil
}

func (m *mockProvider) GetDefaultInstallPath(version string) string {
	return "/path/to/" + version
}

func (m *mockProvider) Validate(version string) error {
	if version == "" {
		return nil
	}
	return nil
}

func TestNewRegistry(t *testing.T) {
	registry := NewRegistry()
	if registry == nil {
		t.Fatal("NewRegistry() returned nil")
	}
}

func TestRegistry_Register(t *testing.T) {
	registry := NewRegistry()
	
	provider := &mockProvider{
		name:        "test",
		displayName: "Test Provider",
		sdkType:     models.JavaSDK,
	}
	
	registry.Register(provider)
	
	// Try to get the registered provider
	got, ok := registry.Get(models.JavaSDK, "test")
	if !ok {
		t.Fatal("Get() failed to retrieve registered provider")
	}
	
	if got.Name() != "test" {
		t.Errorf("Get() returned provider with name %v, want %v", got.Name(), "test")
	}
}

func TestRegistry_Get(t *testing.T) {
	registry := NewRegistry()
	
	provider := &mockProvider{
		name:        "testprovider",
		displayName: "Test Provider",
		sdkType:     models.NodeSDK,
	}
	
	registry.Register(provider)
	
	tests := []struct {
		name         string
		sdkType      models.SDKType
		providerName string
		wantOk       bool
	}{
		{
			name:         "existing provider",
			sdkType:      models.NodeSDK,
			providerName: "testprovider",
			wantOk:       true,
		},
		{
			name:         "non-existing provider",
			sdkType:      models.JavaSDK,
			providerName: "nonexistent",
			wantOk:       false,
		},
		{
			name:         "wrong sdk type",
			sdkType:      models.PythonSDK,
			providerName: "testprovider",
			wantOk:       false,
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, got := registry.Get(tt.sdkType, tt.providerName)
			if got != tt.wantOk {
				t.Errorf("Get() ok = %v, want %v", got, tt.wantOk)
			}
		})
	}
}

func TestRegistry_List(t *testing.T) {
	registry := NewRegistry()
	
	javaProvider1 := &mockProvider{
		name:    "java1",
		sdkType: models.JavaSDK,
	}
	javaProvider2 := &mockProvider{
		name:    "java2",
		sdkType: models.JavaSDK,
	}
	nodeProvider := &mockProvider{
		name:    "node1",
		sdkType: models.NodeSDK,
	}
	
	registry.Register(javaProvider1)
	registry.Register(javaProvider2)
	registry.Register(nodeProvider)
	
	// List Java providers
	javaProviders := registry.List(models.JavaSDK)
	if len(javaProviders) != 2 {
		t.Errorf("List(JavaSDK) returned %d providers, want 2", len(javaProviders))
	}
	
	// List Node providers
	nodeProviders := registry.List(models.NodeSDK)
	if len(nodeProviders) != 1 {
		t.Errorf("List(NodeSDK) returned %d providers, want 1", len(nodeProviders))
	}
	
	// List Python providers (none registered)
	pythonProviders := registry.List(models.PythonSDK)
	if len(pythonProviders) != 0 {
		t.Errorf("List(PythonSDK) returned %d providers, want 0", len(pythonProviders))
	}
}

func TestRegistry_ListAll(t *testing.T) {
	registry := NewRegistry()
	
	provider1 := &mockProvider{
		name:    "provider1",
		sdkType: models.JavaSDK,
	}
	provider2 := &mockProvider{
		name:    "provider2",
		sdkType: models.NodeSDK,
	}
	provider3 := &mockProvider{
		name:    "provider3",
		sdkType: models.PythonSDK,
	}
	
	registry.Register(provider1)
	registry.Register(provider2)
	registry.Register(provider3)
	
	all := registry.ListAll()
	if len(all) != 3 {
		t.Errorf("ListAll() returned %d providers, want 3", len(all))
	}
}

func TestRegistry_MultipleProvidersPerType(t *testing.T) {
	registry := NewRegistry()
	
	// Register multiple providers for the same SDK type
	corretto := &mockProvider{
		name:    "corretto",
		sdkType: models.JavaSDK,
	}
	openjdk := &mockProvider{
		name:    "openjdk",
		sdkType: models.JavaSDK,
	}
	graalvm := &mockProvider{
		name:    "graalvm",
		sdkType: models.JavaSDK,
	}
	
	registry.Register(corretto)
	registry.Register(openjdk)
	registry.Register(graalvm)
	
	// Verify all are retrievable
	if _, ok := registry.Get(models.JavaSDK, "corretto"); !ok {
		t.Error("Failed to retrieve corretto provider")
	}
	if _, ok := registry.Get(models.JavaSDK, "openjdk"); !ok {
		t.Error("Failed to retrieve openjdk provider")
	}
	if _, ok := registry.Get(models.JavaSDK, "graalvm"); !ok {
		t.Error("Failed to retrieve graalvm provider")
	}
	
	// Verify list returns all three
	javaProviders := registry.List(models.JavaSDK)
	if len(javaProviders) != 3 {
		t.Errorf("List(JavaSDK) returned %d providers, want 3", len(javaProviders))
	}
}
