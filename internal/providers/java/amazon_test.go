package java

import (
	"context"
	"strings"
	"testing"
)

func TestAmazonCorrettoProvider_Name(t *testing.T) {
	provider := NewAmazonCorrettoProvider()
	if got := provider.Name(); got != "amazoncorretto" {
		t.Errorf("Name() = %v, want %v", got, "amazoncorretto")
	}
}

func TestAmazonCorrettoProvider_DisplayName(t *testing.T) {
	provider := NewAmazonCorrettoProvider()
	if got := provider.DisplayName(); got != "Amazon Corretto" {
		t.Errorf("DisplayName() = %v, want %v", got, "Amazon Corretto")
	}
}

func TestAmazonCorrettoProvider_GetVersions(t *testing.T) {
	provider := NewAmazonCorrettoProvider()
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

func TestAmazonCorrettoProvider_GetLatestVersion(t *testing.T) {
	provider := NewAmazonCorrettoProvider()
	ctx := context.Background()
	
	version, err := provider.GetLatestVersion(ctx)
	if err != nil {
		t.Fatalf("GetLatestVersion() error = %v", err)
	}
	
	if version != "25.0.0" {
		t.Errorf("GetLatestVersion() = %v, want %v", version, "25.0.0")
	}
}

func TestAmazonCorrettoProvider_GetDownloadURL(t *testing.T) {
	provider := NewAmazonCorrettoProvider()
	
	tests := []struct {
		name    string
		version string
		arch    string
		want    string
		wantErr bool
	}{
		{
			name:    "Corretto 25 x64 architecture",
			version: "25.0.0",
			arch:    "x64",
			want:    "https://corretto.aws/downloads/latest/amazon-corretto-25-x64-windows-jdk.zip",
			wantErr: false,
		},
		{
			name:    "Corretto 21 amd64 architecture",
			version: "21.0.1",
			arch:    "amd64",
			want:    "https://corretto.aws/downloads/latest/amazon-corretto-21-x64-windows-jdk.zip",
			wantErr: false,
		},
		{
			name:    "Corretto 17 x64 architecture",
			version: "17.0.9",
			arch:    "x64",
			want:    "https://corretto.aws/downloads/latest/amazon-corretto-17-x64-windows-jdk.zip",
			wantErr: false,
		},
		{
			name:    "Corretto 11 x64 architecture",
			version: "11.0.21",
			arch:    "x64",
			want:    "https://corretto.aws/downloads/latest/amazon-corretto-11-x64-windows-jdk.zip",
			wantErr: false,
		},
		{
			name:    "Corretto 8 x64 architecture",
			version: "8.392.08.1",
			arch:    "x64",
			want:    "https://corretto.aws/downloads/latest/amazon-corretto-8-x64-windows-jdk.zip",
			wantErr: false,
		},
		{
			name:    "arm64 architecture",
			version: "21.0.1",
			arch:    "arm64",
			want:    "https://corretto.aws/downloads/latest/amazon-corretto-21-aarch64-windows-jdk.zip",
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

func TestAmazonCorrettoProvider_GetDefaultInstallPath(t *testing.T) {
	provider := NewAmazonCorrettoProvider()
	version := "21.0.1"
	
	path := provider.GetDefaultInstallPath(version)
	if path == "" {
		t.Error("GetDefaultInstallPath() returned empty path")
	}
	
	if !strings.Contains(path, "amazoncorretto") {
		t.Errorf("GetDefaultInstallPath() = %v, should contain 'amazoncorretto'", path)
	}
	
	if !strings.Contains(path, version) {
		t.Errorf("GetDefaultInstallPath() = %v, should contain version %v", path, version)
	}
}

func TestAmazonCorrettoProvider_Validate(t *testing.T) {
	provider := NewAmazonCorrettoProvider()
	
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
