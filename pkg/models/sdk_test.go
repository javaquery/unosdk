package models

import (
	"testing"
	"time"
)

func TestSDKType_Constants(t *testing.T) {
	tests := []struct {
		name string
		got  SDKType
		want string
	}{
		{"JavaSDK", JavaSDK, "java"},
		{"NodeSDK", NodeSDK, "node"},
		{"PythonSDK", PythonSDK, "python"},
		{"GoSDK", GoSDK, "go"},
		{"FlutterSDK", FlutterSDK, "flutter"},
		{"MavenSDK", MavenSDK, "maven"},
		{"GradleSDK", GradleSDK, "gradle"},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if string(tt.got) != tt.want {
				t.Errorf("SDKType constant %s = %v, want %v", tt.name, tt.got, tt.want)
			}
		})
	}
}

func TestSDK_Creation(t *testing.T) {
	now := time.Now()
	sdk := SDK{
		ID:          "test-id",
		Name:        "Test SDK",
		Type:        JavaSDK,
		Provider:    "test-provider",
		Version:     "11.0.0",
		InstallPath: "/path/to/sdk",
		DownloadURL: "https://example.com/sdk.zip",
		Checksum:    "abc123",
		Installed:   true,
		InstalledAt: now,
	}
	
	if sdk.ID != "test-id" {
		t.Errorf("SDK.ID = %v, want %v", sdk.ID, "test-id")
	}
	if sdk.Name != "Test SDK" {
		t.Errorf("SDK.Name = %v, want %v", sdk.Name, "Test SDK")
	}
	if sdk.Type != JavaSDK {
		t.Errorf("SDK.Type = %v, want %v", sdk.Type, JavaSDK)
	}
	if sdk.Provider != "test-provider" {
		t.Errorf("SDK.Provider = %v, want %v", sdk.Provider, "test-provider")
	}
	if sdk.Version != "11.0.0" {
		t.Errorf("SDK.Version = %v, want %v", sdk.Version, "11.0.0")
	}
	if sdk.InstallPath != "/path/to/sdk" {
		t.Errorf("SDK.InstallPath = %v, want %v", sdk.InstallPath, "/path/to/sdk")
	}
	if sdk.DownloadURL != "https://example.com/sdk.zip" {
		t.Errorf("SDK.DownloadURL = %v, want %v", sdk.DownloadURL, "https://example.com/sdk.zip")
	}
	if sdk.Checksum != "abc123" {
		t.Errorf("SDK.Checksum = %v, want %v", sdk.Checksum, "abc123")
	}
	if !sdk.Installed {
		t.Error("SDK.Installed should be true")
	}
	if sdk.InstalledAt != now {
		t.Errorf("SDK.InstalledAt = %v, want %v", sdk.InstalledAt, now)
	}
}

func TestProviderInfo_Creation(t *testing.T) {
	info := ProviderInfo{
		Name:        "test-provider",
		DisplayName: "Test Provider",
		Type:        NodeSDK,
		Versions:    []string{"1.0.0", "2.0.0", "3.0.0"},
		Website:     "https://example.com",
	}
	
	if info.Name != "test-provider" {
		t.Errorf("ProviderInfo.Name = %v, want %v", info.Name, "test-provider")
	}
	if info.DisplayName != "Test Provider" {
		t.Errorf("ProviderInfo.DisplayName = %v, want %v", info.DisplayName, "Test Provider")
	}
	if info.Type != NodeSDK {
		t.Errorf("ProviderInfo.Type = %v, want %v", info.Type, NodeSDK)
	}
	if len(info.Versions) != 3 {
		t.Errorf("ProviderInfo.Versions length = %v, want 3", len(info.Versions))
	}
	if info.Website != "https://example.com" {
		t.Errorf("ProviderInfo.Website = %v, want %v", info.Website, "https://example.com")
	}
}

func TestVersion_Creation(t *testing.T) {
	v := Version{
		Major: 1,
		Minor: 2,
		Patch: 3,
		Raw:   "1.2.3",
	}
	
	if v.Major != 1 {
		t.Errorf("Version.Major = %v, want 1", v.Major)
	}
	if v.Minor != 2 {
		t.Errorf("Version.Minor = %v, want 2", v.Minor)
	}
	if v.Patch != 3 {
		t.Errorf("Version.Patch = %v, want 3", v.Patch)
	}
	if v.Raw != "1.2.3" {
		t.Errorf("Version.Raw = %v, want %v", v.Raw, "1.2.3")
	}
}

func TestSDK_DefaultValues(t *testing.T) {
	sdk := SDK{
		Type:    PythonSDK,
		Version: "3.11.0",
	}
	
	if sdk.Installed {
		t.Error("SDK.Installed should default to false")
	}
	if !sdk.InstalledAt.IsZero() {
		t.Error("SDK.InstalledAt should be zero value")
	}
}

func TestSDKType_AllTypes(t *testing.T) {
	types := []SDKType{
		JavaSDK,
		NodeSDK,
		PythonSDK,
		GoSDK,
		FlutterSDK,
		MavenSDK,
		GradleSDK,
	}
	
	if len(types) != 7 {
		t.Errorf("Expected 7 SDK types, got %d", len(types))
	}
	
	// Ensure all types are unique
	seen := make(map[SDKType]bool)
	for _, typ := range types {
		if seen[typ] {
			t.Errorf("Duplicate SDKType found: %v", typ)
		}
		seen[typ] = true
	}
}
