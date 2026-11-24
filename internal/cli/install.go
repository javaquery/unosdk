package cli

import (
	"context"
	"errors"
	"fmt"
	"runtime"

	"github.com/javaquery/unosdk/internal/installer"
	"github.com/javaquery/unosdk/internal/providers"
	"github.com/javaquery/unosdk/internal/providers/c"
	"github.com/javaquery/unosdk/internal/providers/cpp"
	"github.com/javaquery/unosdk/internal/providers/flutter"
	"github.com/javaquery/unosdk/internal/providers/go"
	"github.com/javaquery/unosdk/internal/providers/gradle"
	"github.com/javaquery/unosdk/internal/providers/java"
	"github.com/javaquery/unosdk/internal/providers/maven"
	"github.com/javaquery/unosdk/internal/providers/node"
	"github.com/javaquery/unosdk/internal/providers/python"
	"github.com/javaquery/unosdk/internal/registry"
	"github.com/javaquery/unosdk/pkg/models"
	"github.com/spf13/cobra"
)

var (
	installArch  string
	installPath  string
	skipEnvSetup bool
	setAsDefault bool
)

var installCmd = &cobra.Command{
	Use:   "install [sdk-type] [provider] [version]",
	Short: "Install an SDK",
	Long: `Install an SDK from a specific provider.

When unosdk is run through its shell wrapper (unosdk.cmd for cmd.exe, or the
unosdk function in your PowerShell profile), the installed SDK is also made
active in the current terminal immediately, so no restart is needed.

Examples:
  # Install Amazon Corretto Java 21
  unosdk install java amazoncorretto 21

  # Install latest OpenJDK
  unosdk install java openjdk latest

  # Install Node.js 20
  unosdk install node nodejs 20.10.0

  # Install Python 3.12
  unosdk install python python 3.12.1

  # Install Flutter SDK
  unosdk install flutter flutter latest

  # Install Apache Maven
  unosdk install maven apache 3.9.9

  # Install Gradle
  unosdk install gradle gradle 8.12

  # Install Go
  unosdk install go golang 1.23.5

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
	providerRegistry.Register(java.NewGraalVMProvider())
	providerRegistry.Register(node.NewNodeJSProvider())
	providerRegistry.Register(python.NewPythonProvider())
	providerRegistry.Register(flutter.NewFlutterProvider())
	providerRegistry.Register(maven.NewMavenProvider())
	providerRegistry.Register(gradle.NewGradleProvider())
	providerRegistry.Register(golang.NewGoProvider())
	providerRegistry.Register(cpp.NewMinGWProvider())
	providerRegistry.Register(c.NewMinGWProvider())

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

	// Track whether the environment was configured and whether the current
	// terminal session was updated as well
	envConfigured := false
	sessionApplied := false

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
			envConfigured = true
			fmt.Println("✓ Environment variables configured")

			// Check for conflicts with System PATH
			checkSystemPathConflicts(sdk)

			// Apply the change to the current terminal too, when running via the
			// shell wrapper (UNOSDK_ENV_FILE is set by unosdk.cmd / PowerShell function)
			if err := writeSessionScript(reg, sdk, setAsDefault); err == nil {
				sessionApplied = true
			} else if !errors.Is(err, errNoSession) {
				fmt.Printf("⚠ Warning: Failed to prepare current-session update: %v\n", err)
			}
		}
	}

	fmt.Println("\n✓ Installation complete!")
	if sessionApplied {
		fmt.Println("Active in this terminal now. New terminals will pick it up too.")
	} else if envConfigured {
		fmt.Println("Please restart your terminal for changes to take effect.")
		fmt.Println("(Tip: run unosdk through its shell wrapper to avoid restarts.)")
	}

	return nil
}