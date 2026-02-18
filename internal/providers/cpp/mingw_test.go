package cpp

import (
	"context"
	"strings"
	"testing"

	"github.com/javaquery/unosdk/pkg/models"
)

func TestMinGWProvider_Name(t *testing.T) {
	p := NewMinGWProvider()
	if p.Name() != "mingw" {
		t.Errorf("expected name 'mingw', got '%s'", p.Name())
	}
}

func TestMinGWProvider_DisplayName(t *testing.T) {
	p := NewMinGWProvider()
	if p.DisplayName() != "MinGW-w64" {
		t.Errorf("expected display name 'MinGW-w64', got '%s'", p.DisplayName())
	}
}

func TestMinGWProvider_Type(t *testing.T) {
	p := NewMinGWProvider()
	if p.Type() != models.CppSDK {
		t.Errorf("expected type 'cpp', got '%s'", p.Type())
	}
}

func TestMinGWProvider_GetVersions(t *testing.T) {
	p := NewMinGWProvider()
	versions, err := p.GetVersions(context.Background())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(versions) == 0 {
		t.Error("expected at least one version")
	}
	// Check that 15.2.0 is in the list
	found := false
	for _, v := range versions {
		if v == "15.2.0" {
			found = true
			break
		}
	}
	if !found {
		t.Error("expected version 15.2.0 to be in the list")
	}
}

func TestMinGWProvider_GetLatestVersion(t *testing.T) {
	p := NewMinGWProvider()
	latest, err := p.GetLatestVersion(context.Background())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if latest != "15.2.0" {
		t.Errorf("expected latest version '15.2.0', got '%s'", latest)
	}
}

func TestMinGWProvider_GetDownloadURL(t *testing.T) {
	p := NewMinGWProvider()

	testCases := []struct {
		version  string
		arch     string
		expected string
		hasError bool
	}{
		{"15.2.0", "x64", "winlibs-x86_64", false},
		{"14.2.0", "x64", "winlibs-x86_64", false},
		{"14.2.0", "amd64", "winlibs-x86_64", false},
		{"14.2.0", "x86", "winlibs-i686", false},
		{"13.2.0", "x64", "winlibs-x86_64", false},
		{"99.99.99", "x64", "", true},
	}

	for _, tc := range testCases {
		url, err := p.GetDownloadURL(tc.version, tc.arch)
		if tc.hasError {
			if err == nil {
				t.Errorf("expected error for version %s, arch %s", tc.version, tc.arch)
			}
		} else {
			if err != nil {
				t.Errorf("unexpected error for version %s, arch %s: %v", tc.version, tc.arch, err)
			}
			if !strings.Contains(url, tc.expected) {
				t.Errorf("expected URL to contain '%s', got '%s'", tc.expected, url)
			}
			// Verify .zip format is used (not .7z)
			if !strings.HasSuffix(url, ".zip") {
				t.Errorf("expected URL to end with .zip, got '%s'", url)
			}
		}
	}
}

func TestMinGWProvider_GetDefaultInstallPath(t *testing.T) {
	p := NewMinGWProvider()
	path := p.GetDefaultInstallPath("14.2.0")
	if !strings.Contains(path, "cpp") || !strings.Contains(path, "mingw") || !strings.Contains(path, "14.2.0") {
		t.Errorf("unexpected install path: %s", path)
	}
}

func TestMinGWProvider_Validate(t *testing.T) {
	p := NewMinGWProvider()

	// Valid version
	if err := p.Validate("14.2.0"); err != nil {
		t.Errorf("unexpected error for valid version: %v", err)
	}

	// Empty version
	if err := p.Validate(""); err == nil {
		t.Error("expected error for empty version")
	}
}
