package system

import (
	"fmt"
	"os"
	"strings"

	"golang.org/x/sys/windows/registry"
)

// WindowsEnv handles Windows environment variable management
type WindowsEnv struct{}

// NewWindowsEnv creates a new Windows environment manager
func NewWindowsEnv() *WindowsEnv {
	return &WindowsEnv{}
}

// SetUserEnvironmentVariable sets a user-level environment variable
func (w *WindowsEnv) SetUserEnvironmentVariable(name, value string) error {
	k, err := registry.OpenKey(registry.CURRENT_USER, `Environment`, registry.SET_VALUE)
	if err != nil {
		return fmt.Errorf("failed to open registry key: %w", err)
	}
	defer k.Close()

	if err := k.SetStringValue(name, value); err != nil {
		return fmt.Errorf("failed to set environment variable: %w", err)
	}

	return w.broadcastEnvironmentChange()
}

// GetUserEnvironmentVariable gets a user-level environment variable
func (w *WindowsEnv) GetUserEnvironmentVariable(name string) (string, error) {
	k, err := registry.OpenKey(registry.CURRENT_USER, `Environment`, registry.QUERY_VALUE)
	if err != nil {
		return "", fmt.Errorf("failed to open registry key: %w", err)
	}
	defer k.Close()

	value, _, err := k.GetStringValue(name)
	if err != nil {
		return "", fmt.Errorf("failed to get environment variable: %w", err)
	}

	return value, nil
}

// DeleteUserEnvironmentVariable deletes a user-level environment variable
func (w *WindowsEnv) DeleteUserEnvironmentVariable(name string) error {
	k, err := registry.OpenKey(registry.CURRENT_USER, `Environment`, registry.SET_VALUE)
	if err != nil {
		return fmt.Errorf("failed to open registry key: %w", err)
	}
	defer k.Close()

	if err := k.DeleteValue(name); err != nil {
		return fmt.Errorf("failed to delete environment variable: %w", err)
	}

	return w.broadcastEnvironmentChange()
}

// AddToPath adds a directory to the user's PATH
func (w *WindowsEnv) AddToPath(dir string) error {
	currentPath, err := w.GetUserEnvironmentVariable("Path")
	if err != nil {
		// If Path doesn't exist, create it
		currentPath = ""
	}

	// Check if directory is already in PATH
	paths := strings.Split(currentPath, ";")
	for _, p := range paths {
		if strings.EqualFold(strings.TrimSpace(p), strings.TrimSpace(dir)) {
			// Already in PATH
			return nil
		}
	}

	// Add to PATH
	newPath := currentPath
	if newPath != "" && !strings.HasSuffix(newPath, ";") {
		newPath += ";"
	}
	newPath += dir

	return w.SetUserEnvironmentVariable("Path", newPath)
}

// RemoveFromPath removes a directory from the user's PATH
func (w *WindowsEnv) RemoveFromPath(dir string) error {
	currentPath, err := w.GetUserEnvironmentVariable("Path")
	if err != nil {
		return err
	}

	paths := strings.Split(currentPath, ";")
	var newPaths []string
	
	for _, p := range paths {
		if !strings.EqualFold(strings.TrimSpace(p), strings.TrimSpace(dir)) {
			newPaths = append(newPaths, p)
		}
	}

	newPath := strings.Join(newPaths, ";")
	return w.SetUserEnvironmentVariable("Path", newPath)
}

// SetJavaHome sets the JAVA_HOME environment variable
func (w *WindowsEnv) SetJavaHome(javaPath string) error {
	return w.SetUserEnvironmentVariable("JAVA_HOME", javaPath)
}

// GetJavaHome gets the JAVA_HOME environment variable
func (w *WindowsEnv) GetJavaHome() (string, error) {
	// Try user environment first
	value, err := w.GetUserEnvironmentVariable("JAVA_HOME")
	if err == nil && value != "" {
		return value, nil
	}

	// Fall back to system environment
	return os.Getenv("JAVA_HOME"), nil
}

// broadcastEnvironmentChange notifies Windows that environment variables have changed
func (w *WindowsEnv) broadcastEnvironmentChange() error {
	// In a real implementation, this would use Windows API
	// to broadcast WM_SETTINGCHANGE message
	// For now, we'll just return nil as the registry changes
	// will take effect in new processes
	return nil
}
