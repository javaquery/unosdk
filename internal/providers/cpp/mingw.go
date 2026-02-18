package cpp

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	"github.com/javaquery/unosdk/pkg/models"
)

// MinGWProvider implements the Provider interface for MinGW-w64 C++ toolchain
type MinGWProvider struct{}

// NewMinGWProvider creates a new MinGW-w64 provider
func NewMinGWProvider() *MinGWProvider {
	return &MinGWProvider{}
}

func (p *MinGWProvider) Name() string {
	return "mingw"
}

func (p *MinGWProvider) DisplayName() string {
	return "MinGW-w64"
}

func (p *MinGWProvider) Type() models.SDKType {
	return models.CppSDK
}

func (p *MinGWProvider) GetVersions(ctx context.Context) ([]string, error) {
	// MinGW-w64 releases from GitHub (winlibs builds)
	// These are standalone builds that include GCC, G++, and related tools
	return []string{
		"15.2.0",
		"14.2.0",
		"14.1.0",
		"13.2.0",
		"13.1.0",
		"12.3.0",
		"12.2.0",
		"12.1.0",
		"11.3.0",
		"11.2.0",
		"11.1.0",
	}, nil
}

func (p *MinGWProvider) GetLatestVersion(ctx context.Context) (string, error) {
	return "15.2.0", nil
}

func (p *MinGWProvider) GetDownloadURL(version string, arch string) (string, error) {
	// Map architecture to MinGW-w64 naming convention
	mingwArch := "x86_64"
	switch arch {
	case "x64", "amd64":
		mingwArch = "x86_64"
	case "x86", "386", "i686":
		mingwArch = "i686"
	default:
		mingwArch = "x86_64"
	}

	// Using winlibs releases from GitHub
	// Format: https://github.com/brechtsanders/winlibs_mingw/releases/download/{release}/winlibs-{arch}-posix-seh-gcc-{version}-mingw-w64ucrt-{mingwversion}.zip
	// Using .zip format for native Go support (no external 7z tool needed)

	downloadURLs := map[string]map[string]string{
		"15.2.0": {
			"x86_64": "https://github.com/brechtsanders/winlibs_mingw/releases/download/15.2.0posix-13.0.0-ucrt-r5/winlibs-x86_64-posix-seh-gcc-15.2.0-mingw-w64ucrt-13.0.0-r5.zip",
			"i686":   "https://github.com/brechtsanders/winlibs_mingw/releases/download/15.2.0posix-13.0.0-ucrt-r5/winlibs-i686-posix-dwarf-gcc-15.2.0-mingw-w64ucrt-13.0.0-r5.zip",
		},
		"14.2.0": {
			"x86_64": "https://github.com/brechtsanders/winlibs_mingw/releases/download/14.2.0posix-19.1.1-12.0.0-ucrt-r2/winlibs-x86_64-posix-seh-gcc-14.2.0-mingw-w64ucrt-12.0.0-r2.zip",
			"i686":   "https://github.com/brechtsanders/winlibs_mingw/releases/download/14.2.0posix-19.1.1-12.0.0-ucrt-r2/winlibs-i686-posix-dwarf-gcc-14.2.0-mingw-w64ucrt-12.0.0-r2.zip",
		},
		"14.1.0": {
			"x86_64": "https://github.com/brechtsanders/winlibs_mingw/releases/download/14.1.0posix-18.1.5-11.0.1-ucrt-r1/winlibs-x86_64-posix-seh-gcc-14.1.0-mingw-w64ucrt-11.0.1-r1.zip",
			"i686":   "https://github.com/brechtsanders/winlibs_mingw/releases/download/14.1.0posix-18.1.5-11.0.1-ucrt-r1/winlibs-i686-posix-dwarf-gcc-14.1.0-mingw-w64ucrt-11.0.1-r1.zip",
		},
		"13.2.0": {
			"x86_64": "https://github.com/brechtsanders/winlibs_mingw/releases/download/13.2.0posix-17.0.6-11.0.1-ucrt-r5/winlibs-x86_64-posix-seh-gcc-13.2.0-mingw-w64ucrt-11.0.1-r5.zip",
			"i686":   "https://github.com/brechtsanders/winlibs_mingw/releases/download/13.2.0posix-17.0.6-11.0.1-ucrt-r5/winlibs-i686-posix-dwarf-gcc-13.2.0-mingw-w64ucrt-11.0.1-r5.zip",
		},
		"13.1.0": {
			"x86_64": "https://github.com/brechtsanders/winlibs_mingw/releases/download/13.1.0-16.0.5-11.0.0-ucrt-r5/winlibs-x86_64-posix-seh-gcc-13.1.0-mingw-w64ucrt-11.0.0-r5.zip",
			"i686":   "https://github.com/brechtsanders/winlibs_mingw/releases/download/13.1.0-16.0.5-11.0.0-ucrt-r5/winlibs-i686-posix-dwarf-gcc-13.1.0-mingw-w64ucrt-11.0.0-r5.zip",
		},
		"12.3.0": {
			"x86_64": "https://github.com/brechtsanders/winlibs_mingw/releases/download/12.3.0-16.0.6-11.0.0-ucrt-r1/winlibs-x86_64-posix-seh-gcc-12.3.0-mingw-w64ucrt-11.0.0-r1.zip",
			"i686":   "https://github.com/brechtsanders/winlibs_mingw/releases/download/12.3.0-16.0.6-11.0.0-ucrt-r1/winlibs-i686-posix-dwarf-gcc-12.3.0-mingw-w64ucrt-11.0.0-r1.zip",
		},
		"12.2.0": {
			"x86_64": "https://github.com/brechtsanders/winlibs_mingw/releases/download/12.2.0-15.0.7-10.0.0-ucrt-r4/winlibs-x86_64-posix-seh-gcc-12.2.0-mingw-w64ucrt-10.0.0-r4.zip",
			"i686":   "https://github.com/brechtsanders/winlibs_mingw/releases/download/12.2.0-15.0.7-10.0.0-ucrt-r4/winlibs-i686-posix-dwarf-gcc-12.2.0-mingw-w64ucrt-10.0.0-r4.zip",
		},
		"12.1.0": {
			"x86_64": "https://github.com/brechtsanders/winlibs_mingw/releases/download/12.1.0-14.0.6-10.0.0-ucrt-r3/winlibs-x86_64-posix-seh-gcc-12.1.0-mingw-w64ucrt-10.0.0-r3.zip",
			"i686":   "https://github.com/brechtsanders/winlibs_mingw/releases/download/12.1.0-14.0.6-10.0.0-ucrt-r3/winlibs-i686-posix-dwarf-gcc-12.1.0-mingw-w64ucrt-10.0.0-r3.zip",
		},
		"11.3.0": {
			"x86_64": "https://github.com/brechtsanders/winlibs_mingw/releases/download/11.3.0-14.0.3-10.0.0-ucrt-r3/winlibs-x86_64-posix-seh-gcc-11.3.0-mingw-w64ucrt-10.0.0-r3.zip",
			"i686":   "https://github.com/brechtsanders/winlibs_mingw/releases/download/11.3.0-14.0.3-10.0.0-ucrt-r3/winlibs-i686-posix-dwarf-gcc-11.3.0-mingw-w64ucrt-10.0.0-r3.zip",
		},
		"11.2.0": {
			"x86_64": "https://github.com/brechtsanders/winlibs_mingw/releases/download/11.2.0-13.0.0-9.0.0-ucrt-r5/winlibs-x86_64-posix-seh-gcc-11.2.0-mingw-w64ucrt-9.0.0-r5.zip",
			"i686":   "https://github.com/brechtsanders/winlibs_mingw/releases/download/11.2.0-13.0.0-9.0.0-ucrt-r5/winlibs-i686-posix-dwarf-gcc-11.2.0-mingw-w64ucrt-9.0.0-r5.zip",
		},
		"11.1.0": {
			"x86_64": "https://github.com/brechtsanders/winlibs_mingw/releases/download/11.1.0-12.0.0-9.0.0-r2/winlibs-x86_64-posix-seh-gcc-11.1.0-mingw-w64-9.0.0-r2.zip",
			"i686":   "https://github.com/brechtsanders/winlibs_mingw/releases/download/11.1.0-12.0.0-9.0.0-r2/winlibs-i686-posix-dwarf-gcc-11.1.0-mingw-w64-9.0.0-r2.zip",
		},
	}

	archURLs, ok := downloadURLs[version]
	if !ok {
		return "", fmt.Errorf("version %s not found", version)
	}

	url, ok := archURLs[mingwArch]
	if !ok {
		return "", fmt.Errorf("architecture %s not supported for version %s", arch, version)
	}

	return url, nil
}

func (p *MinGWProvider) GetChecksum(version string, arch string) (string, error) {
	// Checksums would typically be downloaded from the release page
	// For now, returning empty string (optional verification)
	return "", nil
}

func (p *MinGWProvider) GetDefaultInstallPath(version string) string {
	homeDir, _ := os.UserHomeDir()
	return filepath.Join(homeDir, ".unosdk", "cpp", "mingw", version)
}

func (p *MinGWProvider) Validate(version string) error {
	if version == "" {
		return fmt.Errorf("version cannot be empty")
	}
	return nil
}
