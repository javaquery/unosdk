package cli

import (
	"fmt"

	"github.com/spf13/cobra"
)

var (
	version   string
	commit    string
	buildDate string
)

// SetVersionInfo sets the version information
func SetVersionInfo(v, c, d string) {
	version = v
	commit = c
	buildDate = d
}

var rootCmd = &cobra.Command{
	Use:   "unosdk",
	Short: "Multi-SDK installer for Windows",
	Long: `SDK Installer is a CLI tool to install and manage multiple SDKs from different providers.

Supported SDKs:
  - Java (Amazon Corretto, OpenJDK)
  - Node.js
  - Python
  - Maven (Apache Maven)

Examples:
  # Install Amazon Corretto Java 21
  unosdk install java amazoncorretto 21

  # Install latest Node.js
  unosdk install node nodejs latest

  # Install Apache Maven
  unosdk install maven apache 3.9.9

  # Switch to a different Java version
  unosdk switch java openjdk 21

  # List all available providers
  unosdk list

  # Show installed SDKs
  unosdk list --installed`,
}

// Execute runs the root command
func Execute() error {
	return rootCmd.Execute()
}

func init() {
	// Add commands
	rootCmd.AddCommand(installCmd)
	rootCmd.AddCommand(listCmd)
	rootCmd.AddCommand(uninstallCmd)
	rootCmd.AddCommand(switchCmd)
	rootCmd.AddCommand(versionCmd)

	// Global flags
	rootCmd.PersistentFlags().BoolP("verbose", "v", false, "Enable verbose output")
	rootCmd.PersistentFlags().BoolP("quiet", "q", false, "Suppress non-error output")
}

// GetVersion returns the version string
func GetVersion() string {
	if version == "" {
		version = "dev"
	}
	return fmt.Sprintf("unosdk %s (commit: %s, built: %s)", version, commit, buildDate)
}
