package gradle

import (
	"context"
	"strings"
	"testing"
)

func TestGradleProvider_Name(t *testing.T) {
	provider := NewGradleProvider()
	if got := provider.Name(); got != "gradle" {
		t.Errorf("Name() = %v, want %v", got, "gradle")
	}
}

func TestGradleProvider_DisplayName(t *testing.T) {
	provider := NewGradleProvider()
	if got := provider.DisplayName(); got != "Gradle" {
		t.Errorf("DisplayName() = %v, want %v", got, "Gradle")
	}
}

func TestGradleProvider_GetVersions(t *testing.T) {
	provider := NewGradleProvider()
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
		if v == "8.12" {
			found = true
			break
		}
	}
	if !found {
		t.Error("GetVersions() should contain version 8.12")
	}
}

func TestGradleProvider_GetLatestVersion(t *testing.T) {
	provider := NewGradleProvider()
	ctx := context.Background()
	
	version, err := provider.GetLatestVersion(ctx)
	if err != nil {
		t.Fatalf("GetLatestVersion() error = %v", err)
	}
	
	if version != "8.12" {
		t.Errorf("GetLatestVersion() = %v, want %v", version, "8.12")
	}
}

func TestGradleProvider_GetDownloadURL(t *testing.T) {
	provider := NewGradleProvider()
	
	tests := []struct {
		name    string
		version string
		arch    string
		want    string
		wantErr bool
	}{
		{
			name:    "valid version 8.12",
			version: "8.12",
			arch:    "x64",
			want:    "https://services.gradle.org/distributions/gradle-8.12-bin.zip",
			wantErr: false,
		},
		{
			name:    "valid version 8.10",
			version: "8.10",
			arch:    "x64",
			want:    "https://services.gradle.org/distributions/gradle-8.10-bin.zip",
			wantErr: false,
		},
		{
			name:    "valid version 7.6",
			version: "7.6",
			arch:    "arm64",
			want:    "https://services.gradle.org/distributions/gradle-7.6-bin.zip",
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

func TestGradleProvider_GetDownloadURL_Format(t *testing.T) {
	provider := NewGradleProvider()
	
	url, err := provider.GetDownloadURL("8.12", "x64")
	if err != nil {
		t.Fatalf("GetDownloadURL() error = %v", err)
	}
	
	if !strings.HasPrefix(url, "https://services.gradle.org/distributions/") {
		t.Errorf("GetDownloadURL() should start with correct base URL, got %v", url)
	}
	
	if !strings.HasSuffix(url, "-bin.zip") {
		t.Errorf("GetDownloadURL() should end with -bin.zip, got %v", url)
	}
}

func TestGradleProvider_Validate(t *testing.T) {
	provider := NewGradleProvider()
	
	tests := []struct {
		name    string
		version string
		wantErr bool
	}{
		{
			name:    "valid version",
			version: "8.12",
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

func TestGradleProvider_GetDefaultInstallPath(t *testing.T) {
	provider := NewGradleProvider()
	version := "8.12"
	
	path := provider.GetDefaultInstallPath(version)
	if path == "" {
		t.Error("GetDefaultInstallPath() returned empty path")
	}
	
	if !strings.Contains(path, "gradle") {
		t.Errorf("GetDefaultInstallPath() should contain 'gradle', got %v", path)
	}
	
	if !strings.Contains(path, version) {
		t.Errorf("GetDefaultInstallPath() should contain version %v, got %v", version, path)
	}
}
