package cli

import (
	"errors"
	"fmt"

	"github.com/javaquery/unosdk/internal/installer"
	"github.com/javaquery/unosdk/internal/providers"
	"github.com/javaquery/unosdk/internal/registry"
	"github.com/javaquery/unosdk/internal/system"
	"github.com/javaquery/unosdk/pkg/models"
	"github.com/spf13/cobra"
)

var (
	cleanupEnv bool
)

var uninstallCmd = &cobra.Command{
	Use:   "uninstall [sdk-type] [provider] [version]",
	Short: "Uninstall an SDK",
	Long: `Uninstall a previously installed SDK.

When unosdk is run through its shell wrapper (unosdk.cmd for cmd.exe, or the
unosdk function in your PowerShell profile), the current terminal is updated
immediately as well, so no restart is needed.

Examples:
  # Uninstall Amazon Corretto Java 21
  unosdk uninstall java amazoncorretto 21

  # Uninstall Node.js
  unosdk uninstall node nodejs 20.10.0

  # Uninstall and cleanup environment variables
  unosdk uninstall java openjdk 17 --cleanup-env`,
	Args: cobra.ExactArgs(3),
	RunE: runUninstall,
}

func init() {
	uninstallCmd.Flags().BoolVar(&cleanupEnv, "cleanup-env", true, "Remove from PATH and clean environment variables")
}

func runUninstall(cmd *cobra.Command, args []string) error {
	sdkType := models.SDKType(args[0])
	providerName := args[1]
	version := args[2]

	// Initialize registry
	reg, err := registry.NewRegistry()
	if err != nil {
		return fmt.Errorf("failed to initialize registry: %w", err)
	}

	// Get SDK from registry
	sdk, ok := reg.Get(sdkType, providerName, version)
	if !ok {
		return fmt.Errorf("SDK not found: %s %s %s", sdkType, providerName, version)
	}

	fmt.Printf("Uninstalling %s %s %s...\n", sdkType, providerName, version)
	fmt.Printf("  Location: %s\n", sdk.InstallPath)

	// Initialize installer
	providerRegistry := providers.NewRegistry()
	inst := installer.NewInstaller(providerRegistry)

	// Uninstall
	if err := inst.Uninstall(sdk.InstallPath); err != nil {
		return fmt.Errorf("uninstallation failed: %w", err)
	}

	// Remove from registry
	if err := reg.Remove(sdkType, providerName, version); err != nil {
		return fmt.Errorf("failed to remove from registry: %w", err)
	}

	fmt.Printf("✓ Successfully uninstalled %s %s %s\n", sdkType, providerName, version)

	// Cleanup environment variables if requested
	if cleanupEnv {
		wasDefault, err := cleanupEnvironment(sdk)
		if err != nil {
			fmt.Printf("⚠ Warning: Failed to cleanup environment variables: %v\n", err)
		} else {
			fmt.Println("✓ Environment variables cleaned up")

			sessionApplied := false
			newDefaultSet := false

			// If the uninstalled SDK was the default, try to set a new one
			if wasDefault {
				applied, err := setNewDefault(reg, sdk)
				if err != nil {
					fmt.Printf("⚠ No alternative SDK found to set as default\n")
				} else {
					newDefaultSet = true
					sessionApplied = applied
				}
			}

			// No replacement was activated: just drop the removed SDK from the
			// current terminal (and clear JAVA_HOME if it pointed at it)
			if !newDefaultSet {
				clearJavaHome := sdk.Type == models.JavaSDK && wasDefault
				if err := writeSessionRemoval(sdk, clearJavaHome); err == nil {
					sessionApplied = true
				} else if !errors.Is(err, errNoSession) {
					fmt.Printf("⚠ Warning: Failed to prepare current-session update: %v\n", err)
				}
			}

			if sessionApplied {
				fmt.Println("Active in this terminal now. New terminals will pick it up too.")
			} else {
				fmt.Println("Please restart your terminal for changes to take effect.")
				fmt.Println("(Tip: run unosdk through its shell wrapper to avoid restarts.)")
			}
		}
	}

	return nil
}

func cleanupEnvironment(sdk *models.SDK) (bool, error) {
	env := system.NewWindowsEnv()
	wasDefault := false

	switch sdk.Type {
	case models.JavaSDK:
		// Check if this is the current JAVA_HOME
		currentHome, err := env.GetJavaHome()
		if err == nil && currentHome == sdk.InstallPath {
			wasDefault = true
			if err := env.DeleteUserEnvironmentVariable("JAVA_HOME"); err != nil {
				return wasDefault, err
			}
			fmt.Println("  Removed JAVA_HOME")
		}

		// Remove from PATH
		binPath := sdk.InstallPath + "\\bin"
		if err := env.RemoveFromPath(binPath); err != nil {
			return wasDefault, err
		}
		fmt.Println("  Removed from PATH: " + binPath)

	case models.NodeSDK:
		// Remove from PATH
		if err := env.RemoveFromPath(sdk.InstallPath); err != nil {
			return wasDefault, err
		}
		wasDefault = true // Node/Python are default if they were in PATH
		fmt.Println("  Removed from PATH: " + sdk.InstallPath)

	case models.PythonSDK:
		// Remove from PATH
		if err := env.RemoveFromPath(sdk.InstallPath); err != nil {
			return wasDefault, err
		}
		wasDefault = true // Node/Python are default if they were in PATH
		fmt.Println("  Removed from PATH: " + sdk.InstallPath)

		// Remove Scripts directory
		scriptsPath := sdk.InstallPath + "\\Scripts"
		if err := env.RemoveFromPath(scriptsPath); err != nil {
			return wasDefault, err
		}
		fmt.Println("  Removed from PATH: " + scriptsPath)
	}

	return wasDefault, nil
}

// setNewDefault attempts to set a new default SDK after uninstallation.
// Priority: Same provider versions first, then other providers.
// The returned bool reports whether the current terminal session was updated too.
func setNewDefault(reg *registry.Registry, removed *models.SDK) (bool, error) {
	sdkType := removed.Type
	uninstalledProvider := removed.Provider

	// Get all installed SDKs of the same type
	installedSDKs := reg.ListByType(sdkType)

	if len(installedSDKs) == 0 {
		return false, fmt.Errorf("no alternative SDKs found")
	}

	// Sort: prioritize same provider first
	var sameProvider []*models.SDK
	var otherProviders []*models.SDK

	for _, sdk := range installedSDKs {
		if sdk.Provider == uninstalledProvider {
			sameProvider = append(sameProvider, sdk)
		} else {
			otherProviders = append(otherProviders, sdk)
		}
	}

	// Choose the first available SDK (same provider takes precedence)
	var newDefault *models.SDK
	if len(sameProvider) > 0 {
		newDefault = sameProvider[0]
		fmt.Printf("\n⚡ Setting %s %s %s as new default (same provider)...\n",
			newDefault.Type, newDefault.Provider, newDefault.Version)
	} else if len(otherProviders) > 0 {
		newDefault = otherProviders[0]
		fmt.Printf("\n⚡ Setting %s %s %s as new default (alternative provider)...\n",
			newDefault.Type, newDefault.Provider, newDefault.Version)
	} else {
		return false, fmt.Errorf("no alternative SDKs found")
	}

	// Setup environment for the new default
	setJavaHome := (sdkType == models.JavaSDK)
	if err := setupSDKEnvironment(newDefault, setJavaHome); err != nil {
		return false, fmt.Errorf("failed to set new default: %w", err)
	}

	fmt.Printf("✓ Successfully set %s %s %s as default\n",
		newDefault.Type, newDefault.Provider, newDefault.Version)

	// Also switch the current terminal, dropping the removed SDK's PATH entries
	sessionApplied := false
	if err := writeSessionScript(reg, newDefault, setJavaHome, removed); err == nil {
		sessionApplied = true
	} else if !errors.Is(err, errNoSession) {
		fmt.Printf("⚠ Warning: Failed to prepare current-session update: %v\n", err)
	}

	return sessionApplied, nil
}