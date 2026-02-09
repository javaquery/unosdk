package flutter

import (
	"context"
	"strings"
	"testing"
)

func TestFlutterProvider_Name(t *testing.T) {
	provider := NewFlutterProvider()
	if got := provider.Name(); got != "flutter" {
		t.Errorf("Name() = %v, want %v", got, "flutter")
	}
}

func TestFlutterProvider_DisplayName(t *testing.T) {
	provider := NewFlutterProvider()
	if got := provider.DisplayName(); got != "Flutter SDK" {
		t.Errorf("DisplayName() = %v, want %v", got, "Flutter SDK")
	}
}

func TestFlutterProvider_GetVersions(t *testing.T) {
	provider := NewFlutterProvider()
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
		if v == "3.27.2" {
			found = true
			break
		}
	}
	if !found {
		t.Error("GetVersions() should contain version 3.27.2")
	}
	
	// Check if "latest" is present
	foundLatest := false
	for _, v := range versions {
		if v == "latest" {
			foundLatest = true
			break
		}
	}
	if !foundLatest {
		t.Error("GetVersions() should contain 'latest'")
	}
}

func TestFlutterProvider_GetLatestVersion(t *testing.T) {
	provider := NewFlutterProvider()
	ctx := context.Background()
	
	version, err := provider.GetLatestVersion(ctx)
	if err != nil {
		t.Fatalf("GetLatestVersion() error = %v", err)
	}
	
	if version != "3.27.2" {
		t.Errorf("GetLatestVersion() = %v, want %v", version, "3.27.2")
	}
}

func TestFlutterProvider_GetDownloadURL(t *testing.T) {
	provider := NewFlutterProvider()
	
	tests := []struct {
		name    string
		version string
		arch    string
		want    string
		wantErr bool
	}{
		{
			name:    "Flutter 3.27.2 x64 architecture",
			version: "3.27.2",
			arch:    "x64",
			want:    "https://storage.googleapis.com/flutter_infra_release/releases/stable/windows/flutter_windows_3.27.2-stable.zip",
			wantErr: false,
		},
		{
			name:    "Flutter 3.27.1 x64 architecture",
			version: "3.27.1",
			arch:    "x64",
			want:    "https://storage.googleapis.com/flutter_infra_release/releases/stable/windows/flutter_windows_3.27.1-stable.zip",
			wantErr: false,
		},
		{
			name:    "Flutter 3.24.5 arm64 architecture",
			version: "3.24.5",
			arch:    "arm64",
			want:    "https://storage.googleapis.com/flutter_infra_release/releases/stable/windows/flutter_windows_3.24.5-stable.zip",
			wantErr: false,
		},
		{
			name:    "Flutter latest version resolves to 3.27.2",
			version: "latest",
			arch:    "x64",
			want:    "https://storage.googleapis.com/flutter_infra_release/releases/stable/windows/flutter_windows_3.27.2-stable.zip",
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
			if err == nil && got != tt.want {
				t.Errorf("GetDownloadURL() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestFlutterProvider_GetDefaultInstallPath(t *testing.T) {
	provider := NewFlutterProvider()
	version := "3.27.2"
	
	path := provider.GetDefaultInstallPath(version)
	if path == "" {
		t.Error("GetDefaultInstallPath() returned empty path")
	}
	
	if !strings.Contains(path, "flutter") {
		t.Errorf("GetDefaultInstallPath() = %v, should contain 'flutter'", path)
	}
	
	if !strings.Contains(path, version) {
		t.Errorf("GetDefaultInstallPath() = %v, should contain version %v", path, version)
	}
}

func TestFlutterProvider_Validate(t *testing.T) {
	provider := NewFlutterProvider()
	
	tests := []struct {
		name    string
		version string
		wantErr bool
	}{
		{
			name:    "valid version",
			version: "3.27.2",
			wantErr: false,
		},
		{
			name:    "valid latest version",
			version: "latest",
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
