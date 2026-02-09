package utils

import (
	"testing"
)

func TestValidator_ValidateVersion(t *testing.T) {
	tests := []struct {
		name    string
		version string
		wantErr bool
	}{
		{
			name:    "valid semantic version",
			version: "1.2.3",
			wantErr: false,
		},
		{
			name:    "valid version with major only",
			version: "11",
			wantErr: false,
		},
		{
			name:    "valid version with major.minor",
			version: "11.0",
			wantErr: false,
		},
		{
			name:    "valid complex version",
			version: "1.23.5",
			wantErr: false,
		},
		{
			name:    "empty version",
			version: "",
			wantErr: true,
		},
		{
			name:    "latest alias",
			version: "latest",
			wantErr: false,
		},
		{
			name:    "invalid version with letters",
			version: "1.2.3abc",
			wantErr: true,
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateVersion(tt.version)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateVersion(%v) error = %v, wantErr %v", tt.version, err, tt.wantErr)
			}
		})
	}
}

func TestValidator_ValidateSDKType(t *testing.T) {
	tests := []struct {
		name    string
		sdkType string
		wantErr bool
	}{
		{
			name:    "valid java",
			sdkType: "java",
			wantErr: false,
		},
		{
			name:    "valid node",
			sdkType: "node",
			wantErr: false,
		},
		{
			name:    "valid python",
			sdkType: "python",
			wantErr: false,
		},
		{
			name:    "valid go",
			sdkType: "go",
			wantErr: false,
		},
		{
			name:    "invalid type",
			sdkType: "invalid",
			wantErr: true,
		},
		{
			name:    "empty type",
			sdkType: "",
			wantErr: true,
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateSDKType(tt.sdkType)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateSDKType(%v) error = %v, wantErr %v", tt.sdkType, err, tt.wantErr)
			}
		})
	}
}

func TestValidator_ValidateProvider(t *testing.T) {
	tests := []struct {
		name     string
		provider string
		wantErr  bool
	}{
		{
			name:     "valid provider",
			provider: "openjdk",
			wantErr:  false,
		},
		{
			name:     "valid provider with dash",
			provider: "amazon-corretto",
			wantErr:  false,
		},
		{
			name:     "empty provider",
			provider: "",
			wantErr:  true,
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateProvider(tt.provider)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateProvider(%v) error = %v, wantErr %v", tt.provider, err, tt.wantErr)
			}
		})
	}
}

