package cli

import (
	"context"
	"fmt"
	"runtime"

	"github.com/spf13/cobra"
	"github.com/javaquery/unosdk/internal/installer"
	"github.com/javaquery/unosdk/internal/providers"
	"github.com/javaquery/unosdk/internal/providers/java"
	"github.com/javaquery/unosdk/internal/providers/node"
	"github.com/javaquery/unosdk/internal/providers/python"
	"github.com/javaquery/unosdk/internal/registry"
	"github.com/javaquery/unosdk/pkg/models"
)

var (
	installArch       string
	installPath       string
	skipEnvSetup      bool
	setAsDefault      bool
)

var installCmd = &cobra.Command{
	Use:   "install [sdk-type] [provider] [version]",
	Short: "Install an SDK",
	Long: `Install an SDK from a specific provider.

Examples:
  # Install Amazon Corretto Java 21
  unosdk install java amazoncorretto 21

  # Install latest OpenJDK
  unosdk install java openjdk latest

  # Install Node.js 20
  unosdk install node nodejs 20.10.0

  # Install Python 3.12
  unosdk install python python 3.12.1

  # Install with custom architecture
  unosdk install java openjdk 21 --arch x64`,
	Args: cobra.ExactArgs(3),
	RunE: runInstall,
}

func init() {
	installCmd.Flags().StringVar(&installArch, "arch", runtime.GOARCH, "Architecture (x64, x86, arm64)")
	installCmd.Flags().StringVar(&installPath, "path", "", "Custom installation path")
	installCmd.Flags().BoolVar(&skipEnvSetup, "skip-env", false, "Skip environment variable setup")
	installCmd.Flags().BoolVar(&setAsDefault, "set-default", true, "Set as default SDK for the type")
}

func runInstall(cmd *cobra.Command, args []string) error {
	sdkType := models.SDKType(args[0])
	providerName := args[1]
	version := args[2]

	// Initialize provider registry
	providerRegistry := providers.NewRegistry()
	
	// Register providers
	providerRegistry.Register(java.NewAmazonCorrettoProvider())
	providerRegistry.Register(java.NewOpenJDKProvider())
	providerRegistry.Register(node.NewNodeJSProvider())
	providerRegistry.Register(python.NewPythonProvider())

	// Initialize installer
	inst := installer.NewInstaller(providerRegistry)

	// Initialize registry
	reg, err := registry.NewRegistry()
	if err != nil {
		return fmt.Errorf("failed to initialize registry: %w", err)
	}

	fmt.Printf("Installing %s %s version %s...\n", sdkType, providerName, version)

	// Install SDK
	ctx := context.Background()
	sdk, err := inst.Install(ctx, sdkType, providerName, version, installArch)
	if err != nil {
		return fmt.Errorf("installation failed: %w", err)
	}

	// Add to registry
	if err := reg.Add(sdk); err != nil {
		return fmt.Errorf("failed to register SDK: %w", err)
	}

	fmt.Printf("✓ Successfully installed %s %s %s\n", sdkType, providerName, version)
	fmt.Printf("  Location: %s\n", sdk.InstallPath)

	// Setup environment variables (Windows-specific)
	if !skipEnvSetup && runtime.GOOS == "windows" {
		// Cleanup existing PATH entries first
		if err := cleanupExistingSDKPaths(reg, sdk); err != nil {
			fmt.Printf("⚠ Warning: Failed to cleanup existing PATH entries: %v\n", err)
		}

		if err := setupSDKEnvironment(sdk, setAsDefault); err != nil {
			fmt.Printf("⚠ Warning: Failed to setup environment variables: %v\n", err)
			fmt.Println("  You may need to configure environment variables manually.")
		} else {
			fmt.Println("✓ Environment variables configured")
			
			// Check for conflicts with System PATH
			checkSystemPathConflicts(sdk)
		}
	}

	fmt.Println("\n✓ Installation complete!")

	return nil
}


