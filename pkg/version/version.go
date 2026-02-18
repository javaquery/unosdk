package version

// Version is the current version of UnoSDK
// Update this manually for new releases
const Version = "1.2.0"

// Build information (set via ldflags during build)
var (
	GitCommit = "unknown"
	BuildDate = "unknown"
)
