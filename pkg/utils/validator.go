package utils

import (
	"fmt"
	"regexp"
)

var (
	// Version pattern: 1.2.3, 11.0.21, etc.
	versionPattern = regexp.MustCompile(`^\d+(\.\d+)*$`)
)

// ValidateVersion validates a version string
func ValidateVersion(version string) error {
	if version == "" {
		return fmt.Errorf("version cannot be empty")
	}

	if version == "latest" {
		return nil
	}

	if !versionPattern.MatchString(version) {
		return fmt.Errorf("invalid version format: %s", version)
	}

	return nil
}

// ValidateSDKType validates an SDK type
func ValidateSDKType(sdkType string) error {
	validTypes := map[string]bool{
		"java":   true,
		"node":   true,
		"python": true,
		"go":     true,
	}

	if !validTypes[sdkType] {
		return fmt.Errorf("invalid SDK type: %s", sdkType)
	}

	return nil
}

// ValidateProvider validates a provider name
func ValidateProvider(provider string) error {
	if provider == "" {
		return fmt.Errorf("provider cannot be empty")
	}

	return nil
}
