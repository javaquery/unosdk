package java

import (
	"context"
	"strings"
	"testing"
)

func TestGraalVMProvider_Name(t *testing.T) {
	provider := NewGraalVMProvider()
	if got := provider.Name(); got != "graalvm" {
		t.Errorf("Name() = %v, want %v", got, "graalvm")
	}
}

func TestGraalVMProvider_DisplayName(t *testing.T) {
	provider := NewGraalVMProvider()
	if got := provider.DisplayName(); got != "GraalVM" {
		t.Errorf("DisplayName() = %v, want %v", got, "GraalVM")
	}
}

func TestGraalVMProvider_GetVersions(t *testing.T) {
	provider := NewGraalVMProvider()
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
		if v == "25.0.2" {
			found = true
			break
		}
	}
	if !found {
		t.Error("GetVersions() should contain version 25.0.2")
	}
}

func TestGraalVMProvider_GetLatestVersion(t *testing.T) {
	provider := NewGraalVMProvider()
	ctx := context.Background()
	
	version, err := provider.GetLatestVersion(ctx)
	if err != nil {
		t.Fatalf("GetLatestVersion() error = %v", err)
	}
	
	if version != "25.0.2" {
		t.Errorf("GetLatestVersion() = %v, want %v", version, "25.0.2")
	}
}

func TestGraalVMProvider_GetDownloadURL(t *testing.T) {
	provider := NewGraalVMProvider()
	
	tests := []struct {
		name    string
		version string
		arch    string
		want    string
		wantErr bool
	}{
		{
			name:    "GraalVM 25 amd64 architecture",
			version: "25.0.2",
			arch:    "amd64",
			want:    "https://github.com/graalvm/graalvm-ce-builds/releases/download/jdk-25.0.0/graalvm-community-jdk-25.0.0_windows-x64_bin.zip",
			wantErr: false,
		},
		{
			name:    "GraalVM 21 amd64 architecture",
			version: "21.0.2",
			arch:    "amd64",
			want:    "https://github.com/graalvm/graalvm-ce-builds/releases/download/jdk-21.0.2/graalvm-community-jdk-21.0.2_windows-x64_bin.zip",
			wantErr: false,
		},
		{
			name:    "GraalVM 17 amd64 architecture",
			version: "17.0.9",
			arch:    "amd64",
			want:    "https://github.com/graalvm/graalvm-ce-builds/releases/download/jdk-17.0.9/graalvm-community-jdk-17.0.9_windows-x64_bin.zip",
			wantErr: false,
		},
		{
			name:    "GraalVM 21 arm64 architecture",
			version: "21",
			arch:    "arm64",
			want:    "https://github.com/graalvm/graalvm-ce-builds/releases/download/jdk-21.0.2/graalvm-community-jdk-21.0.2_windows-aarch64_bin.zip",
			wantErr: false,
		},
		{
			name:    "unsupported version",
			version: "99.0.0",
			arch:    "x64",
			want:    "",
			wantErr: true,
		},
		{
			name:    "unsupported architecture",
			version: "21.0.2",
			arch:    "x86",
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

func TestGraalVMProvider_GetDefaultInstallPath(t *testing.T) {
	provider := NewGraalVMProvider()
	version := "21.0.2"
	
	path := provider.GetDefaultInstallPath(version)
	if path == "" {
		t.Error("GetDefaultInstallPath() returned empty path")
	}
	
	if !strings.Contains(path, "graalvm") {
		t.Errorf("GetDefaultInstallPath() = %v, should contain 'graalvm'", path)
	}
	
	if !strings.Contains(path, version) {
		t.Errorf("GetDefaultInstallPath() = %v, should contain version %v", path, version)
	}
}

func TestGraalVMProvider_Validate(t *testing.T) {
	provider := NewGraalVMProvider()
	
	tests := []struct {
		name    string
		version string
		wantErr bool
	}{
		{
			name:    "valid version",
			version: "21.0.2",
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
