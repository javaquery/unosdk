package models

// Provider represents an SDK provider
type Provider interface {
	GetName() string
	GetType() SDKType
	ListVersions() ([]string, error)
	GetDownloadURL(version string) (string, error)
	GetChecksum(version string) (string, error)
}
