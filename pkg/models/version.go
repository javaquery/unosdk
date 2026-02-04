package models

import (
	"fmt"
	"strconv"
	"strings"
)

// ParseVersion parses a version string into a Version struct
func ParseVersion(v string) (*Version, error) {
	parts := strings.Split(strings.TrimPrefix(v, "v"), ".")
	if len(parts) < 2 {
		return nil, fmt.Errorf("invalid version format: %s", v)
	}

	major, err := strconv.Atoi(parts[0])
	if err != nil {
		return nil, fmt.Errorf("invalid major version: %w", err)
	}

	minor, err := strconv.Atoi(parts[1])
	if err != nil {
		return nil, fmt.Errorf("invalid minor version: %w", err)
	}

	patch := 0
	if len(parts) > 2 {
		patch, _ = strconv.Atoi(parts[2])
	}

	return &Version{
		Major: major,
		Minor: minor,
		Patch: patch,
		Raw:   v,
	}, nil
}

// Compare compares two versions
func (v *Version) Compare(other *Version) int {
	if v.Major != other.Major {
		return v.Major - other.Major
	}
	if v.Minor != other.Minor {
		return v.Minor - other.Minor
	}
	return v.Patch - other.Patch
}

// String returns the string representation
func (v *Version) String() string {
	return v.Raw
}
