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

		case models.NodeSDK:
			_ = env.RemoveFromPath(installedSDK.InstallPath)

		case models.PythonSDK:
			_ = env.RemoveFromPath(installedSDK.InstallPath)
			scriptsPath := installedSDK.InstallPath + "\\Scripts"
			_ = env.RemoveFromPath(scriptsPath)
		}
	}

	return nil
}

// setupSDKEnvironment configures environment variables for the target SDK
// If setJavaHome is true, JAVA_HOME will be set for Java SDKs
func setupSDKEnvironment(sdk *models.SDK, setJavaHome bool) error {
	env := system.NewWindowsEnv()

	switch sdk.Type {
	case models.JavaSDK:
		// Set JAVA_HOME if requested
		if setJavaHome {
			if err := env.SetJavaHome(sdk.InstallPath); err != nil {
				return fmt.Errorf("failed to set JAVA_HOME: %w", err)
			}
			fmt.Println("  Set JAVA_HOME=" + sdk.InstallPath)
		}

		// Add to PATH
		binPath := sdk.InstallPath + "\\bin"
		if err := env.AddToPath(binPath); err != nil {
			return fmt.Errorf("failed to add to PATH: %w", err)
		}
		fmt.Println("  Added to PATH: " + binPath)

	case models.NodeSDK:
		// Add to PATH
		if err := env.AddToPath(sdk.InstallPath); err != nil {
			return fmt.Errorf("failed to add to PATH: %w", err)
		}
		fmt.Println("  Added to PATH: " + sdk.InstallPath)

	case models.PythonSDK:
		// Add to PATH
		if err := env.AddToPath(sdk.InstallPath); err != nil {
			return fmt.Errorf("failed to add to PATH: %w", err)
		}
		fmt.Println("  Added to PATH: " + sdk.InstallPath)

		// Add Scripts directory
		scriptsPath := sdk.InstallPath + "\\Scripts"
		if err := env.AddToPath(scriptsPath); err != nil {
			return fmt.Errorf("failed to add Scripts to PATH: %w", err)
		}
		fmt.Println("  Added to PATH: " + scriptsPath)
	}

	return nil
}
