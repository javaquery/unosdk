package maven

import (
	"context"
	"strings"
	"testing"
)

func TestMavenProvider_Name(t *testing.T) {
	provider := NewMavenProvider()
	if got := provider.Name(); got != "apache" {
		t.Errorf("Name() = %v, want %v", got, "apache")
	}
}

func TestMavenProvider_DisplayName(t *testing.T) {
	provider := NewMavenProvider()
	if got := provider.DisplayName(); got != "Apache Maven" {
		t.Errorf("DisplayName() = %v, want %v", got, "Apache Maven")
	}
}

func TestMavenProvider_GetVersions(t *testing.T) {
	provider := NewMavenProvider()
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
		if v == "3.9.9" {
			found = true
			break
		}
	}
	if !found {
		t.Error("GetVersions() should contain version 3.9.9")
	}
}

func TestMavenProvider_GetLatestVersion(t *testing.T) {
	provider := NewMavenProvider()
	ctx := context.Background()
	
	version, err := provider.GetLatestVersion(ctx)
	if err != nil {
		t.Fatalf("GetLatestVersion() error = %v", err)
	}
	
	if version != "3.9.9" {
		t.Errorf("GetLatestVersion() = %v, want %v", version, "3.9.9")
	}
}

func TestMavenProvider_GetDownloadURL(t *testing.T) {
	provider := NewMavenProvider()
	
	tests := []struct {
		name    string
		version string
		arch    string
		want    string
		wantErr bool
	}{
		{
			name:    "valid version 3.9.9",
			version: "3.9.9",
			arch:    "x64",
			want:    "https://archive.apache.org/dist/maven/maven-3/3.9.9/binaries/apache-maven-3.9.9-bin.zip",
			wantErr: false,
		},
		{
			name:    "valid version 3.9.8",
			version: "3.9.8",
			arch:    "x64",
			want:    "https://archive.apache.org/dist/maven/maven-3/3.9.8/binaries/apache-maven-3.9.8-bin.zip",
			wantErr: false,
		},
		{
			name:    "valid version 3.8.8",
			version: "3.8.8",
			arch:    "arm64",
			want:    "https://archive.apache.org/dist/maven/maven-3/3.8.8/binaries/apache-maven-3.8.8-bin.zip",
			wantErr: false,
		},
		{
			name:    "unsupported version",
			version: "2.0.0",
			arch:    "x64",
			want:    "",
			wantErr: true,
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

func TestMavenProvider_GetDownloadURL_ContainsBinariesPath(t *testing.T) {
	provider := NewMavenProvider()
	
	url, err := provider.GetDownloadURL("3.9.8", "x64")
	if err != nil {
		t.Fatalf("GetDownloadURL() error = %v", err)
	}
	
	if !strings.Contains(url, "/binaries/") {
		t.Errorf("GetDownloadURL() should contain '/binaries/' in path, got %v", url)
	}
}

func TestMavenProvider_Validate(t *testing.T) {
	provider := NewMavenProvider()
	
	tests := []struct {
		name    string
		version string
		wantErr bool
	}{
		{
			name:    "valid version",
			version: "3.9.9",
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

func TestMavenProvider_GetDefaultInstallPath(t *testing.T) {
	provider := NewMavenProvider()
	version := "3.9.9"
	
	path := provider.GetDefaultInstallPath(version)
	if path == "" {
		t.Error("GetDefaultInstallPath() returned empty path")
	}
	
	if !strings.Contains(path, "maven") {
		t.Errorf("GetDefaultInstallPath() should contain 'maven', got %v", path)
	}
	
	if !strings.Contains(path, version) {
		t.Errorf("GetDefaultInstallPath() should contain version %v, got %v", version, path)
	}
}
