package system

import (
	"fmt"
	"os"
	"strings"
	"syscall"
	"unsafe"

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

// AddToPath adds a directory to the user's PATH (prepends to take precedence)
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
			// Already in PATH, remove it so we can prepend it
			var newPaths []string
			for _, path := range paths {
				if !strings.EqualFold(strings.TrimSpace(path), strings.TrimSpace(dir)) {
					newPaths = append(newPaths, path)
				}
			}
			currentPath = strings.Join(newPaths, ";")
			break
		}
	}

	// Prepend to PATH (add at the beginning to take precedence over Cursor/VS Code extensions)
	newPath := dir
	if currentPath != "" {
		newPath += ";" + currentPath
	}

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
	const (
		HWND_BROADCAST   = 0xFFFF
		WM_SETTINGCHANGE = 0x001A
		SMTO_ABORTIFHUNG = 0x0002
	)

	user32 := syscall.NewLazyDLL("user32.dll")
	sendMessageTimeout := user32.NewProc("SendMessageTimeoutW")

	environment, err := syscall.UTF16PtrFromString("Environment")
	if err != nil {
		return fmt.Errorf("failed to convert string: %w", err)
	}

	// Broadcast WM_SETTINGCHANGE to notify all windows about environment change
	ret, _, _ := sendMessageTimeout.Call(
		uintptr(HWND_BROADCAST),
		uintptr(WM_SETTINGCHANGE),
		0,
		uintptr(unsafe.Pointer(environment)),
		uintptr(SMTO_ABORTIFHUNG),
		2000, // 2 second timeout
		0,
	)

	if ret == 0 {
		// Not necessarily an error - some windows may not respond
		// but changes will still take effect in new processes
	}

	return nil
}

// GetSystemEnvironmentVariable gets a system-level environment variable from registry
func (w *WindowsEnv) GetSystemEnvironmentVariable(name string) (string, error) {
	k, err := registry.OpenKey(registry.LOCAL_MACHINE, `SYSTEM\CurrentControlSet\Control\Session Manager\Environment`, registry.QUERY_VALUE)
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

// DetectSDKConflicts checks for SDK installations in System PATH that would take precedence
func (w *WindowsEnv) DetectSDKConflicts(sdkType string) []string {
	systemPath, err := w.GetSystemEnvironmentVariable("Path")
	if err != nil {
		return nil
	}

	var conflicts []string
	paths := strings.Split(systemPath, ";")
	
	for _, p := range paths {
		pTrimmed := strings.TrimSpace(p)
		if pTrimmed == "" {
			continue
		}
		pLower := strings.ToLower(pTrimmed)
		
		switch sdkType {
		case "java":
			// Check for common Java installation paths (excluding unosdk)
			if strings.Contains(pLower, "unosdk") {
				continue // Skip our own installations
			}
			if (strings.Contains(pLower, "java") || strings.Contains(pLower, "jdk") || strings.Contains(pLower, "jre")) && 
			   (strings.Contains(pLower, "bin") || strings.Contains(pLower, "javapath") || 
			    strings.Contains(pLower, "corretto") || strings.Contains(pLower, "openjdk")) {
				conflicts = append(conflicts, pTrimmed)
			}
		case "node":
			// Check for common Node.js installation paths (excluding unosdk)
			if strings.Contains(pLower, "unosdk") {
				continue // Skip our own installations
			}
			if strings.Contains(pLower, "nodejs") || 
			   (strings.Contains(pLower, "node") && (strings.Contains(pLower, "program files") || strings.Contains(pLower, "programfiles"))) {
				conflicts = append(conflicts, pTrimmed)
			}
		case "python":
			// Check for common Python installation paths (excluding unosdk)
			if strings.Contains(pLower, "unosdk") {
				continue // Skip our own installations
			}
			if strings.Contains(pLower, "python") && 
			   (strings.Contains(pLower, "program files") || strings.Contains(pLower, "programfiles") || 
			    strings.Contains(pLower, "appdata") || strings.Contains(pLower, "scripts")) {
				conflicts = append(conflicts, pTrimmed)
			}
		case "maven":
			// Check for common Maven installation paths (excluding unosdk)
			if strings.Contains(pLower, "unosdk") {
				continue // Skip our own installations
			}
			if strings.Contains(pLower, "maven") && strings.Contains(pLower, "bin") {
				conflicts = append(conflicts, pTrimmed)
			}
		case "flutter":
			// Check for common Flutter installation paths (excluding unosdk)
			if strings.Contains(pLower, "unosdk") {
				continue // Skip our own installations
			}
			if strings.Contains(pLower, "flutter") && strings.Contains(pLower, "bin") {
				conflicts = append(conflicts, pTrimmed)
			}
		}
	}
	
	return conflicts
}

// RemoveFromSystemPath removes directories from System PATH (requires admin privileges)
func (w *WindowsEnv) RemoveFromSystemPath(dirsToRemove []string) error {
	k, err := registry.OpenKey(registry.LOCAL_MACHINE, `SYSTEM\CurrentControlSet\Control\Session Manager\Environment`, registry.QUERY_VALUE|registry.SET_VALUE)
	if err != nil {
		return fmt.Errorf("failed to open registry key (admin rights required): %w", err)
	}
	defer k.Close()

	currentPath, _, err := k.GetStringValue("Path")
	if err != nil {
		return fmt.Errorf("failed to get System PATH: %w", err)
	}

	paths := strings.Split(currentPath, ";")
	var newPaths []string
	removed := 0
	
	for _, p := range paths {
		shouldRemove := false
		pTrimmed := strings.TrimSpace(p)
		
		for _, dirToRemove := range dirsToRemove {
			if strings.EqualFold(pTrimmed, strings.TrimSpace(dirToRemove)) {
				shouldRemove = true
				removed++
				break
			}
		}
		
		if !shouldRemove && pTrimmed != "" {
			newPaths = append(newPaths, p)
		}
	}

	if removed == 0 {
		return nil // Nothing to remove
	}

	newPath := strings.Join(newPaths, ";")
	if err := k.SetStringValue("Path", newPath); err != nil {
		return fmt.Errorf("failed to update System PATH: %w", err)
	}

	return w.broadcastEnvironmentChange()
}

// SetSystemEnvironmentVariable sets a system-level environment variable (requires admin)
func (w *WindowsEnv) SetSystemEnvironmentVariable(name, value string) error {
	k, err := registry.OpenKey(registry.LOCAL_MACHINE, `SYSTEM\CurrentControlSet\Control\Session Manager\Environment`, registry.SET_VALUE)
	if err != nil {
		return fmt.Errorf("failed to open registry key (admin rights required): %w", err)
	}
	defer k.Close()

	if err := k.SetStringValue(name, value); err != nil {
		return fmt.Errorf("failed to set system environment variable: %w", err)
	}

	return w.broadcastEnvironmentChange()
}

// AddToSystemPath adds a directory to System PATH, prepending it (requires admin)
func (w *WindowsEnv) AddToSystemPath(dir string) error {
	k, err := registry.OpenKey(registry.LOCAL_MACHINE, `SYSTEM\CurrentControlSet\Control\Session Manager\Environment`, registry.QUERY_VALUE|registry.SET_VALUE)
	if err != nil {
		return fmt.Errorf("failed to open registry key (admin rights required): %w", err)
	}
	defer k.Close()

	currentPath, _, err := k.GetStringValue("Path")
	if err != nil {
		currentPath = ""
	}

	// Check if directory is already in PATH and remove it
	paths := strings.Split(currentPath, ";")
	for _, p := range paths {
		if strings.EqualFold(strings.TrimSpace(p), strings.TrimSpace(dir)) {
			var newPaths []string
			for _, path := range paths {
				if !strings.EqualFold(strings.TrimSpace(path), strings.TrimSpace(dir)) {
					newPaths = append(newPaths, path)
				}
			}
			currentPath = strings.Join(newPaths, ";")
			break
		}
	}

	// Prepend to PATH
	newPath := dir
	if currentPath != "" {
		newPath += ";" + currentPath
	}

	if err := k.SetStringValue("Path", newPath); err != nil {
		return fmt.Errorf("failed to update System PATH: %w", err)
	}

	return w.broadcastEnvironmentChange()
}

// RemoveFromSystemPathSingle removes a single directory from System PATH (requires admin)
func (w *WindowsEnv) RemoveFromSystemPathSingle(dir string) error {
	k, err := registry.OpenKey(registry.LOCAL_MACHINE, `SYSTEM\CurrentControlSet\Control\Session Manager\Environment`, registry.QUERY_VALUE|registry.SET_VALUE)
	if err != nil {
		return fmt.Errorf("failed to open registry key (admin rights required): %w", err)
	}
	defer k.Close()

	currentPath, _, err := k.GetStringValue("Path")
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
	if err := k.SetStringValue("Path", newPath); err != nil {
		return fmt.Errorf("failed to update System PATH: %w", err)
	}

	return w.broadcastEnvironmentChange()
}

// SetSystemJavaHome sets JAVA_HOME in System environment (requires admin)
func (w *WindowsEnv) SetSystemJavaHome(javaPath string) error {
	return w.SetSystemEnvironmentVariable("JAVA_HOME", javaPath)
}

// IsAdmin checks if the current process has administrator privileges
func (w *WindowsEnv) IsAdmin() bool {
	k, err := registry.OpenKey(registry.LOCAL_MACHINE, `SYSTEM\CurrentControlSet\Control\Session Manager\Environment`, registry.SET_VALUE)
	if err != nil {
		return false
	}
	k.Close()
	return true
}
