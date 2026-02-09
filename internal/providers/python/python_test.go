package python

import (
	"context"
	"strings"
	"testing"
)

func TestPythonProvider_Name(t *testing.T) {
	provider := NewPythonProvider()
	if got := provider.Name(); got != "python" {
		t.Errorf("Name() = %v, want %v", got, "python")
	}
}

func TestPythonProvider_DisplayName(t *testing.T) {
	provider := NewPythonProvider()
	if got := provider.DisplayName(); got != "Python" {
		t.Errorf("DisplayName() = %v, want %v", got, "Python")
	}
}

func TestPythonProvider_GetVersions(t *testing.T) {
	provider := NewPythonProvider()
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
		if v == "3.12.1" {
			found = true
			break
		}
	}
	if !found {
		t.Error("GetVersions() should contain version 3.12.1")
	}
}

func TestPythonProvider_GetLatestVersion(t *testing.T) {
	provider := NewPythonProvider()
	ctx := context.Background()
	
	version, err := provider.GetLatestVersion(ctx)
	if err != nil {
		t.Fatalf("GetLatestVersion() error = %v", err)
	}
	
	if version != "3.12.1" {
		t.Errorf("GetLatestVersion() = %v, want %v", version, "3.12.1")
	}
}

func TestPythonProvider_GetDownloadURL(t *testing.T) {
	provider := NewPythonProvider()
	
	tests := []struct {
		name    string
		version string
		arch    string
		want    string
		wantErr bool
	}{
		{
			name:    "Python 3.12.1 amd64 architecture",
			version: "3.12.1",
			arch:    "amd64",
			want:    "https://www.python.org/ftp/python/3.12.1/python-3.12.1-amd64.exe",
			wantErr: false,
		},
		{
			name:    "Python 3.11.7 amd64 architecture",
			version: "3.11.7",
			arch:    "amd64",
			want:    "https://www.python.org/ftp/python/3.11.7/python-3.11.7-amd64.exe",
			wantErr: false,
		},
		{
			name:    "Python 3.10.13 default architecture",
			version: "3.10.13",
			arch:    "",
			want:    "https://www.python.org/ftp/python/3.10.13/python-3.10.13-amd64.exe",
			wantErr: false,
		},
		{
			name:    "Python 3.9.18 amd64 architecture",
			version: "3.9.18",
			arch:    "amd64",
			want:    "https://www.python.org/ftp/python/3.9.18/python-3.9.18-amd64.exe",
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

func TestPythonProvider_GetDefaultInstallPath(t *testing.T) {
	provider := NewPythonProvider()
	version := "3.12.1"
	
	path := provider.GetDefaultInstallPath(version)
	if path == "" {
		t.Error("GetDefaultInstallPath() returned empty path")
	}
	
	if !strings.Contains(path, "python") {
		t.Errorf("GetDefaultInstallPath() = %v, should contain 'python'", path)
	}
	
	if !strings.Contains(path, version) {
		t.Errorf("GetDefaultInstallPath() = %v, should contain version %v", path, version)
	}
}

func TestPythonProvider_Validate(t *testing.T) {
	provider := NewPythonProvider()
	
	tests := []struct {
		name    string
		version string
		wantErr bool
	}{
		{
			name:    "valid version",
			version: "3.12.1",
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
