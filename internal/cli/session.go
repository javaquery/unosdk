package cli

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/javaquery/unosdk/internal/registry"
	"github.com/javaquery/unosdk/pkg/models"
)

var errNoSession = errors.New("not running under unosdk shell wrapper")

// sdkPathEntries returns the PATH directories an SDK contributes.
// (Can also replace the big per-type switches in env_setup.go.)
func sdkPathEntries(sdk *models.SDK) []string {
	switch sdk.Type {
	case models.NodeSDK:
		return []string{sdk.InstallPath}
	case models.PythonSDK:
		return []string{sdk.InstallPath, sdk.InstallPath + `\Scripts`}
	default: // java, maven, gradle, flutter, go, cpp, c
		return []string{sdk.InstallPath + `\bin`}
	}
}

func normalize(p string) string {
	return strings.ToLower(strings.TrimRight(filepath.Clean(strings.TrimSpace(p)), `\`))
}

// rebuildPath removes every entry in `remove` from the current PATH and prepends `add`.
func rebuildPath(current string, remove, add []string) string {
	drop := map[string]bool{}
	for _, r := range remove {
		drop[normalize(r)] = true
	}
	for _, a := range add {
		drop[normalize(a)] = true
	}

	out := append([]string{}, add...)
	for _, p := range strings.Split(current, ";") {
		if strings.TrimSpace(p) == "" || drop[normalize(p)] {
			continue
		}
		out = append(out, p)
	}
	return strings.Join(out, ";")
}

// writeSessionScript writes a script that applies the switch to the *current* shell.
// The syntax is chosen by the file extension: .ps1 => PowerShell, otherwise cmd.
func writeSessionScript(reg *registry.Registry, sdk *models.SDK, setJavaHome bool) error {
	file := os.Getenv("UNOSDK_ENV_FILE")
	if file == "" {
		return errNoSession
	}

	var remove []string
	for _, s := range reg.ListByType(sdk.Type) {
		remove = append(remove, sdkPathEntries(s)...)
	}
	newPath := rebuildPath(os.Getenv("PATH"), remove, sdkPathEntries(sdk))

	vars := [][2]string{{"PATH", newPath}}
	if sdk.Type == models.JavaSDK && setJavaHome {
		vars = append(vars, [2]string{"JAVA_HOME", sdk.InstallPath})
	}

	var sb strings.Builder
	if strings.EqualFold(filepath.Ext(file), ".ps1") {
		for _, v := range vars {
			fmt.Fprintf(&sb, "$env:%s = '%s'\r\n", v[0], strings.ReplaceAll(v[1], "'", "''"))
		}
	} else {
		sb.WriteString("@echo off\r\n")
		for _, v := range vars {
			// % must be doubled inside a batch file
			fmt.Fprintf(&sb, "set \"%s=%s\"\r\n", v[0], strings.ReplaceAll(v[1], "%", "%%"))
		}
	}
	return os.WriteFile(file, []byte(sb.String()), 0600)
}