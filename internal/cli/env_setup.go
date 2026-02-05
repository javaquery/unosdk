package cli

import (
	"fmt"

	"github.com/javaquery/unosdk/internal/registry"
	"github.com/javaquery/unosdk/internal/system"
	"github.com/javaquery/unosdk/pkg/models"
)

// cleanupExistingSDKPaths removes all existing PATH entries for the same SDK type
// This includes other unosdk installations and other installations of the same SDK
func cleanupExistingSDKPaths(reg *registry.Registry, sdk *models.SDK) error {
	env := system.NewWindowsEnv()
	isAdmin := env.IsAdmin()

	// Get all installed SDKs of the same type
	installedSDKs := reg.ListByType(sdk.Type)

	for _, installedSDK := range installedSDKs {
		// Skip the current SDK being installed/switched to
		if installedSDK.InstallPath == sdk.InstallPath {
			continue
		}

		// Remove PATH entries based on SDK type
		switch sdk.Type {
		case models.JavaSDK:
			binPath := installedSDK.InstallPath + "\\bin"
			_ = env.RemoveFromPath(binPath)
			if isAdmin {
				_ = env.RemoveFromSystemPathSingle(binPath)
			}

		case models.NodeSDK:
			_ = env.RemoveFromPath(installedSDK.InstallPath)
			if isAdmin {
				_ = env.RemoveFromSystemPathSingle(installedSDK.InstallPath)
			}

		case models.PythonSDK:
			_ = env.RemoveFromPath(installedSDK.InstallPath)
			scriptsPath := installedSDK.InstallPath + "\\Scripts"
			_ = env.RemoveFromPath(scriptsPath)
			if isAdmin {
				_ = env.RemoveFromSystemPathSingle(installedSDK.InstallPath)
				_ = env.RemoveFromSystemPathSingle(scriptsPath)
			}
		}
	}

	return nil
}

// checkSystemPathConflicts detects and removes (if admin) or warns about Java installations in System PATH
func checkSystemPathConflicts(sdk *models.SDK) {
	// Only check for Java SDKs
	if sdk.Type != models.JavaSDK {
		return
	}

	env := system.NewWindowsEnv()
	conflicts := env.DetectJavaConflicts()
	
	if len(conflicts) == 0 {
		return
	}

	fmt.Println("\n⚠ Found Java installation(s) in System PATH that will take precedence:")
	for _, path := range conflicts {
		fmt.Printf("  - %s\n", path)
	}

	// Check if running with admin privileges
	if env.IsAdmin() {
		fmt.Println("\n⚡ Running with administrator privileges - automatically removing conflicts...")
		
		if err := env.RemoveFromSystemPath(conflicts); err != nil {
			fmt.Printf("❌ Failed to remove from System PATH: %v\n", err)
			showManualInstructions()
		} else {
			fmt.Println("✓ Successfully removed conflicting Java paths from System PATH")
			fmt.Println("  Your unosdk-managed Java will now take precedence")
		}
	} else {
		fmt.Println("\n⚠ Not running with administrator privileges")
		fmt.Println("  To automatically fix this, run the command as Administrator")
		fmt.Println("  Or follow these manual steps:")
		showManualInstructions()
	}
}

func showManualInstructions() {
	fmt.Println("\nManual fix steps:")
	fmt.Println("  1. Press Win+R, type 'sysdm.cpl' and press Enter")
	fmt.Println("  2. Go to 'Advanced' tab → 'Environment Variables'")
	fmt.Println("  3. Under 'System variables', select 'Path' → 'Edit'")
	fmt.Println("  4. Remove the Java-related entries listed above")
	fmt.Println("  5. Click 'OK' on all dialogs and restart your terminal")
}

// setupSDKEnvironment configures environment variables for the target SDK
// If setJavaHome is true, JAVA_HOME will be set for Java SDKs
func setupSDKEnvironment(sdk *models.SDK, setJavaHome bool) error {
	env := system.NewWindowsEnv()
	isAdmin := env.IsAdmin()

	switch sdk.Type {
	case models.JavaSDK:
		// Set JAVA_HOME if requested
		if setJavaHome {
			// Set in User environment
			if err := env.SetJavaHome(sdk.InstallPath); err != nil {
				return fmt.Errorf("failed to set User JAVA_HOME: %w", err)
			}
			fmt.Println("  Set User JAVA_HOME=" + sdk.InstallPath)
			
			// Also set in System environment if running as admin
			if isAdmin {
				if err := env.SetSystemJavaHome(sdk.InstallPath); err != nil {
					fmt.Printf("  ⚠ Failed to set System JAVA_HOME: %v\n", err)
				} else {
					fmt.Println("  Set System JAVA_HOME=" + sdk.InstallPath)
				}
			}
		}

		// Add to PATH
		binPath := sdk.InstallPath + "\\bin"
		
		// Add to User PATH
		if err := env.AddToPath(binPath); err != nil {
			return fmt.Errorf("failed to add to User PATH: %w", err)
		}
		fmt.Println("  Added to User PATH: " + binPath)
		
		// Also add to System PATH if running as admin
		if isAdmin {
			if err := env.AddToSystemPath(binPath); err != nil {
				fmt.Printf("  ⚠ Failed to add to System PATH: %v\n", err)
			} else {
				fmt.Println("  Added to System PATH: " + binPath)
			}
		}

	case models.NodeSDK:
		// Add to User PATH
		if err := env.AddToPath(sdk.InstallPath); err != nil {
			return fmt.Errorf("failed to add to User PATH: %w", err)
		}
		fmt.Println("  Added to User PATH: " + sdk.InstallPath)
		
		// Also add to System PATH if running as admin
		if isAdmin {
			if err := env.AddToSystemPath(sdk.InstallPath); err != nil {
				fmt.Printf("  ⚠ Failed to add to System PATH: %v\n", err)
			} else {
				fmt.Println("  Added to System PATH: " + sdk.InstallPath)
			}
		}

	case models.PythonSDK:
		// Add to User PATH
		if err := env.AddToPath(sdk.InstallPath); err != nil {
			return fmt.Errorf("failed to add to User PATH: %w", err)
		}
		fmt.Println("  Added to User PATH: " + sdk.InstallPath)
		
		// Add Scripts directory to User PATH
		scriptsPath := sdk.InstallPath + "\\Scripts"
		if err := env.AddToPath(scriptsPath); err != nil {
			return fmt.Errorf("failed to add Scripts to User PATH: %w", err)
		}
		fmt.Println("  Added to User PATH: " + scriptsPath)
		
		// Also add to System PATH if running as admin
		if isAdmin {
			if err := env.AddToSystemPath(sdk.InstallPath); err != nil {
				fmt.Printf("  ⚠ Failed to add to System PATH: %v\n", err)
			} else {
				fmt.Println("  Added to System PATH: " + sdk.InstallPath)
			}
			
			if err := env.AddToSystemPath(scriptsPath); err != nil {
				fmt.Printf("  ⚠ Failed to add Scripts to System PATH: %v\n", err)
			} else {
				fmt.Println("  Added to System PATH: " + scriptsPath)
			}
		}
	}

	return nil
}
