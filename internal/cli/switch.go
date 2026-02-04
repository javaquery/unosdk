package cli

import (
	"fmt"
	"runtime"

	"github.com/spf13/cobra"
	"github.com/javaquery/unosdk/internal/registry"
	"github.com/javaquery/unosdk/pkg/models"
)

var switchCmd = &cobra.Command{
	Use:   "switch [sdk-type] [provider] [version]",
	Short: "Switch to a different SDK version",
	Long: `Switch the default SDK to a different installed version.

This command updates environment variables (e.g., JAVA_HOME, PATH) to point to
the specified SDK version. The SDK must already be installed.

Examples:
  # Switch to Java OpenJDK 21
  unosdk switch java openjdk 21

  # Switch to Node.js 20.10.0
  unosdk switch node nodejs 20.10.0

  # Switch to Python 3.12.1
  unosdk switch python python 3.12.1`,
	Args: cobra.ExactArgs(3),
	RunE: runSwitch,
}

func init() {
	// No additional flags needed for now
}

func runSwitch(cmd *cobra.Command, args []string) error {
	sdkType := models.SDKType(args[0])
	providerName := args[1]
	version := args[2]

	// Validate SDK type
	if !isValidSDKType(sdkType) {
		return fmt.Errorf("invalid SDK type: %s (valid types: java, node, python)", sdkType)
	}

	// Initialize registry
	reg, err := registry.NewRegistry()
	if err != nil {
		return fmt.Errorf("failed to initialize registry: %w", err)
	}

	// Check if SDK is installed
	sdk, exists := reg.Get(sdkType, providerName, version)
	if !exists {
		return fmt.Errorf("SDK not found: %s %s %s\nPlease install it first using: unosdk install %s %s %s", 
			sdkType, providerName, version, sdkType, providerName, version)
	}

	fmt.Printf("Switching to %s %s %s...\n", sdkType, providerName, version)

	// Setup environment variables (Windows-specific)
	if runtime.GOOS == "windows" {
		// First, cleanup existing PATH entries for this SDK type
		if err := cleanupExistingSDKPaths(reg, sdk); err != nil {
			fmt.Printf("⚠ Warning: Failed to cleanup existing PATH entries: %v\n", err)
		}

		// Setup environment variables for the target SDK (always set JAVA_HOME for switch)
		if err := setupSDKEnvironment(sdk, true); err != nil {
			return fmt.Errorf("failed to setup environment variables: %w", err)
		}

		fmt.Println("✓ Environment variables configured")
		fmt.Printf("  Location: %s\n", sdk.InstallPath)
	} else {
		return fmt.Errorf("switch command is currently only supported on Windows")
	}

	fmt.Printf("\n✓ Successfully switched to %s %s %s\n", sdkType, providerName, version)
	fmt.Println("Please restart your terminal for changes to take effect.")

	return nil
}

// isValidSDKType checks if the SDK type is valid
func isValidSDKType(sdkType models.SDKType) bool {
	switch sdkType {
	case models.JavaSDK, models.NodeSDK, models.PythonSDK:
		return true
	default:
		return false
	}
}
