package models

import "time"

// SDKType represents the type of SDK
type SDKType string

const (
	JavaSDK    SDKType = "java"
	NodeSDK    SDKType = "node"
	PythonSDK  SDKType = "python"
	GoSDK      SDKType = "go"
	FlutterSDK SDKType = "flutter"
	MavenSDK   SDKType = "maven"
)

// SDK represents an SDK installation
type SDK struct {
	ID          string    `json:"id" yaml:"id"`
	Name        string    `json:"name" yaml:"name"`
	Type        SDKType   `json:"type" yaml:"type"`
	Provider    string    `json:"provider" yaml:"provider"`
	Version     string    `json:"version" yaml:"version"`
	InstallPath string    `json:"install_path" yaml:"install_path"`
	DownloadURL string    `json:"download_url" yaml:"download_url"`
	Checksum    string    `json:"checksum" yaml:"checksum"`
	Installed   bool      `json:"installed" yaml:"installed"`
	InstalledAt time.Time `json:"installed_at,omitempty" yaml:"installed_at,omitempty"`
}

// ProviderInfo represents SDK provider information
type ProviderInfo struct {
	Name        string   `json:"name" yaml:"name"`
	DisplayName string   `json:"display_name" yaml:"display_name"`
	Type        SDKType  `json:"type" yaml:"type"`
	Versions    []string `json:"versions" yaml:"versions"`
	Website     string   `json:"website" yaml:"website"`
}

// Version represents a version string with comparison capabilities
type Version struct {
	Major int
	Minor int
	Patch int
	Raw   string
}

// InstallationStatus represents the status of an installation
type InstallationStatus struct {
	SDK      *SDK
	Status   string
	Progress int
	Message  string
	Error    error
}
