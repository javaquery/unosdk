package java

import (
	"context"
	"strings"
	"testing"
)

func TestOpenJDKProvider_Name(t *testing.T) {
	provider := NewOpenJDKProvider()
	if got := provider.Name(); got != "openjdk" {
		t.Errorf("Name() = %v, want %v", got, "openjdk")
	}
}

func TestOpenJDKProvider_DisplayName(t *testing.T) {
	provider := NewOpenJDKProvider()
	if got := provider.DisplayName(); got != "OpenJDK" {
		t.Errorf("DisplayName() = %v, want %v", got, "OpenJDK")
	}
}

func TestOpenJDKProvider_GetVersions(t *testing.T) {
	provider := NewOpenJDKProvider()
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
		if v == "25.0.0" {
			found = true
			break
		}
	}
	if !found {
		t.Error("GetVersions() should contain version 25.0.0")
	}
}

func TestOpenJDKProvider_GetLatestVersion(t *testing.T) {
	provider := NewOpenJDKProvider()
	ctx := context.Background()
	
	version, err := provider.GetLatestVersion(ctx)
	if err != nil {
		t.Fatalf("GetLatestVersion() error = %v", err)
	}
	
	if version != "25.0.0" {
		t.Errorf("GetLatestVersion() = %v, want %v", version, "25.0.0")
	}
}

func TestOpenJDKProvider_GetDownloadURL(t *testing.T) {
	provider := NewOpenJDKProvider()
	
	tests := []struct {
		name    string
		version string
		arch    string
		want    string
		wantErr bool
	}{
		{
			name:    "JDK 25 x64 architecture",
			version: "25.0.0",
			arch:    "x64",
			want:    "https://github.com/adoptium/temurin25-binaries/releases/download/jdk-25+36/OpenJDK25U-jdk_x64_windows_hotspot_25_36.zip",
			wantErr: false,
		},
		{
			name:    "JDK 21 amd64 architecture",
			version: "21.0.1",
			arch:    "amd64",
			want:    "https://github.com/adoptium/temurin21-binaries/releases/download/jdk-21.0.1+12/OpenJDK21U-jdk_x64_windows_hotspot_21.0.1_12.zip",
			wantErr: false,
		},
		{
			name:    "JDK 17 x64 architecture",
			version: "17.0.9",
			arch:    "x64",
			want:    "https://github.com/adoptium/temurin17-binaries/releases/download/jdk-17.0.9+9/OpenJDK17U-jdk_x64_windows_hotspot_17.0.9_9.zip",
			wantErr: false,
		},
		{
			name:    "JDK 11 x64 architecture",
			version: "11.0.21",
			arch:    "x64",
			want:    "https://github.com/adoptium/temurin11-binaries/releases/download/jdk-11.0.21+9/OpenJDK11U-jdk_x64_windows_hotspot_11.0.21_9.zip",
			wantErr: false,
		},
		{
			name:    "JDK 8 x64 architecture",
			version: "8u392",
			arch:    "x64",
			want:    "https://github.com/adoptium/temurin8-binaries/releases/download/jdk8u392-b08/OpenJDK8U-jdk_x64_windows_hotspot_8u392b08.zip",
			wantErr: false,
		},
		{
			name:    "arm64 architecture",
			version: "21.0.1",
			arch:    "arm64",
			want:    "https://github.com/adoptium/temurin21-binaries/releases/download/jdk-21.0.1+12/OpenJDK21U-jdk_aarch64_windows_hotspot_21.0.1_12.zip",
			wantErr: false,
		},
		{
			name:    "unsupported version",
			version: "99.0.0",
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
			if err == nil && got != tt.want {
				t.Errorf("GetDownloadURL() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestOpenJDKProvider_GetDefaultInstallPath(t *testing.T) {
	provider := NewOpenJDKProvider()
	version := "21.0.1"
	
	path := provider.GetDefaultInstallPath(version)
	if path == "" {
		t.Error("GetDefaultInstallPath() returned empty path")
	}
	
	if !strings.Contains(path, "openjdk") {
		t.Errorf("GetDefaultInstallPath() = %v, should contain 'openjdk'", path)
	}
	
	if !strings.Contains(path, version) {
		t.Errorf("GetDefaultInstallPath() = %v, should contain version %v", path, version)
	}
}

func TestOpenJDKProvider_Validate(t *testing.T) {
	provider := NewOpenJDKProvider()
	
	tests := []struct {
		name    string
		version string
		wantErr bool
	}{
		{
			name:    "valid version",
			version: "21.0.1",
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
