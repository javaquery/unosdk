package system

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestNewPathManager(t *testing.T) {
	pm := NewPathManager()
	if pm == nil {
		t.Fatal("NewPathManager() returned nil")
	}
}

func TestPathManager_IsInPath(t *testing.T) {
	pm := NewPathManager()
	
	// Save original PATH
	originalPath := os.Getenv("PATH")
	defer os.Setenv("PATH", originalPath)
	
	// Set a test PATH
	testPath := "C:\\test\\path1;C:\\test\\path2;C:\\test\\path3"
	os.Setenv("PATH", testPath)
	
	tests := []struct {
		name string
		dir  string
		want bool
	}{
		{
			name: "path exists",
			dir:  "C:\\test\\path1",
			want: true,
		},
		{
			name: "path exists with different case",
			dir:  "c:\\test\\path2",
			want: true,
		},
		{
			name: "path does not exist",
			dir:  "C:\\test\\nonexistent",
			want: false,
		},
		{
			name: "path with trailing slash",
			dir:  "C:\\test\\path3\\",
			want: true,
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := pm.IsInPath(tt.dir)
			if got != tt.want {
				t.Errorf("IsInPath(%v) = %v, want %v", tt.dir, got, tt.want)
			}
		})
	}
}

func TestPathManager_GetSDKInstallPath(t *testing.T) {
	pm := NewPathManager()
	
	path, err := pm.GetSDKInstallPath()
	if err != nil {
		t.Fatalf("GetSDKInstallPath() error = %v", err)
	}
	
	if path == "" {
		t.Error("GetSDKInstallPath() returned empty path")
	}
	
	if !strings.Contains(path, ".unosdk") {
		t.Errorf("GetSDKInstallPath() should contain '.unosdk', got %v", path)
	}
	
	if !strings.Contains(path, "sdks") {
		t.Errorf("GetSDKInstallPath() should contain 'sdks', got %v", path)
	}
	
	// Path should be absolute
	if !filepath.IsAbs(path) {
		t.Errorf("GetSDKInstallPath() should return absolute path, got %v", path)
	}
}

func TestPathManager_GetCachePath(t *testing.T) {
	pm := NewPathManager()
	
	path, err := pm.GetCachePath()
	if err != nil {
		t.Fatalf("GetCachePath() error = %v", err)
	}
	
	if path == "" {
		t.Error("GetCachePath() returned empty path")
	}
	
	if !strings.Contains(path, ".unosdk") {
		t.Errorf("GetCachePath() should contain '.unosdk', got %v", path)
	}
	
	if !strings.Contains(path, "cache") {
		t.Errorf("GetCachePath() should contain 'cache', got %v", path)
	}
	
	// Path should be absolute
	if !filepath.IsAbs(path) {
		t.Errorf("GetCachePath() should return absolute path, got %v", path)
	}
}

func TestPathManager_GetRegistryPath(t *testing.T) {
	pm := NewPathManager()
	
	path, err := pm.GetRegistryPath()
	if err != nil {
		t.Fatalf("GetRegistryPath() error = %v", err)
	}
	
	if path == "" {
		t.Error("GetRegistryPath() returned empty path")
	}
	
	if !strings.Contains(path, ".unosdk") {
		t.Errorf("GetRegistryPath() should contain '.unosdk', got %v", path)
	}
	
	if !strings.Contains(path, "registry.json") {
		t.Errorf("GetRegistryPath() should contain 'registry.json', got %v", path)
	}
	
	// Path should be absolute
	if !filepath.IsAbs(path) {
		t.Errorf("GetRegistryPath() should return absolute path, got %v", path)
	}
}

func TestPathManager_PathsAreConsistent(t *testing.T) {
	pm := NewPathManager()
	
	sdkPath, err := pm.GetSDKInstallPath()
	if err != nil {
		t.Fatalf("GetSDKInstallPath() error = %v", err)
	}
	
	cachePath, err := pm.GetCachePath()
	if err != nil {
		t.Fatalf("GetCachePath() error = %v", err)
	}
	
	registryPath, err := pm.GetRegistryPath()
	if err != nil {
		t.Fatalf("GetRegistryPath() error = %v", err)
	}
	
	// All paths should share the same base .unosdk directory
	baseSDK := filepath.Dir(sdkPath)
	baseCache := filepath.Dir(cachePath)
	baseRegistry := filepath.Dir(registryPath)
	
	if baseSDK != baseCache {
		t.Errorf("SDK and Cache paths have different base directories: %v vs %v", baseSDK, baseCache)
	}
	
	if baseSDK != baseRegistry {
		t.Errorf("SDK and Registry paths have different base directories: %v vs %v", baseSDK, baseRegistry)
	}
}

func TestPathManager_IsInPath_EmptyPath(t *testing.T) {
	pm := NewPathManager()
	
	// Save original PATH
	originalPath := os.Getenv("PATH")
	defer os.Setenv("PATH", originalPath)
	
	// Set empty PATH
	os.Setenv("PATH", "")
	
	result := pm.IsInPath("C:\\test\\path")
	if result {
		t.Error("IsInPath() should return false for empty PATH")
	}
}

func TestPathManager_IsInPath_MultipleOccurrences(t *testing.T) {
	pm := NewPathManager()
	
	// Save original PATH
	originalPath := os.Getenv("PATH")
	defer os.Setenv("PATH", originalPath)
	
	// Set a test PATH with duplicates
	testPath := "C:\\test\\path1;C:\\test\\path2;C:\\test\\path1"
	os.Setenv("PATH", testPath)
	
	// Should still find it
	result := pm.IsInPath("C:\\test\\path1")
	if !result {
		t.Error("IsInPath() should return true even with duplicates")
	}
}
