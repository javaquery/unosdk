package golang

import (
	"context"
	"strings"
	"testing"
)

func TestGoProvider_Name(t *testing.T) {
	provider := NewGoProvider()
	if got := provider.Name(); got != "golang" {
		t.Errorf("Name() = %v, want %v", got, "golang")
	}
}

func TestGoProvider_DisplayName(t *testing.T) {
	provider := NewGoProvider()
	if got := provider.DisplayName(); got != "Go" {
		t.Errorf("DisplayName() = %v, want %v", got, "Go")
	}
}

func TestGoProvider_GetVersions(t *testing.T) {
	provider := NewGoProvider()
	ctx := context.Background()
	
	versions, err := provider.GetVersions(ctx)
	if err != nil {
		t.Fatalf("GetVersions() error = %v", err)
	}
	
	if len(versions) == 0 {
		t.Error("GetVersions() returned empty list")
	}
	
	// Check if latest version is present
	found := false
	for _, v := range versions {
		if v == "1.23.5" {
			found = true
			break
		}
	}
	if !found {
		t.Error("GetVersions() should contain version 1.23.5")
	}
}

func TestGoProvider_GetLatestVersion(t *testing.T) {
	provider := NewGoProvider()
	ctx := context.Background()
	
	version, err := provider.GetLatestVersion(ctx)
	if err != nil {
		t.Fatalf("GetLatestVersion() error = %v", err)
	}
	
	if version != "1.23.5" {
		t.Errorf("GetLatestVersion() = %v, want %v", version, "1.23.5")
	}
}

func TestGoProvider_GetDownloadURL(t *testing.T) {
	provider := NewGoProvider()
	
	tests := []struct {
		name    string
		version string
		arch    string
		want    string
		wantErr bool
	}{
		{
			name:    "amd64 architecture",
			version: "1.23.5",
			arch:    "x64",
			want:    "https://go.dev/dl/go1.23.5.windows-amd64.zip",
			wantErr: false,
		},
		{
			name:    "amd64 architecture explicit",
			version: "1.23.5",
			arch:    "amd64",
			want:    "https://go.dev/dl/go1.23.5.windows-amd64.zip",
			wantErr: false,
		},
		{
			name:    "386 architecture",
			version: "1.22.10",
			arch:    "x86",
			want:    "https://go.dev/dl/go1.22.10.windows-386.zip",
			wantErr: false,
		},
		{
			name:    "arm64 architecture",
			version: "1.21.0",
			arch:    "arm64",
			want:    "https://go.dev/dl/go1.21.0.windows-arm64.zip",
			wantErr: false,
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := provider.GetDownloadURL(tt.version, tt.arch)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetDownloadURL() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("GetDownloadURL() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGoProvider_GetDownloadURL_Format(t *testing.T) {
	provider := NewGoProvider()
	
	url, err := provider.GetDownloadURL("1.23.5", "x64")
	if err != nil {
		t.Fatalf("GetDownloadURL() error = %v", err)
	}
	
	if !strings.HasPrefix(url, "https://go.dev/dl/") {
		t.Errorf("GetDownloadURL() should start with correct base URL, got %v", url)
	}
	
	if !strings.Contains(url, ".windows-") {
		t.Errorf("GetDownloadURL() should contain '.windows-', got %v", url)
	}
	
	if !strings.HasSuffix(url, ".zip") {
		t.Errorf("GetDownloadURL() should end with .zip, got %v", url)
	}
}

func TestGoProvider_Validate(t *testing.T) {
	provider := NewGoProvider()
	
	tests := []struct {
		name    string
		version string
		wantErr bool
	}{
		{
			name:    "valid version",
			version: "1.23.5",
			wantErr: false,
		},
		{
			name:    "empty version",
			version: "",
			wantErr: true,
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := provider.Validate(tt.version)
			if (err != nil) != tt.wantErr {
				t.Errorf("Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestGoProvider_GetDefaultInstallPath(t *testing.T) {
	provider := NewGoProvider()
	version := "1.23.5"
	
	path := provider.GetDefaultInstallPath(version)
	if path == "" {
		t.Error("GetDefaultInstallPath() returned empty path")
	}
	
	if !strings.Contains(path, "go") {
		t.Errorf("GetDefaultInstallPath() should contain 'go', got %v", path)
	}
	
	if !strings.Contains(path, version) {
		t.Errorf("GetDefaultInstallPath() should contain version %v, got %v", version, path)
	}
}
