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

// errNoSession is returned when unosdk was not started through its shell
// wrapper (unosdk.cmd or the PowerShell profile function). In that case the
// current terminal cannot be updated and the caller falls back to the
// "restart your terminal" message.
var errNoSession = errors.New("not running under the unosdk shell wrapper")

// sessionVar is one environment change to apply to the current shell.
type sessionVar struct {
	name  string
	value string
	unset bool
}

// sdkPathEntries returns the PATH directories an SDK contributes.
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

// normalizePathEntry makes PATH entries comparable (case, trailing slash, spaces).
func normalizePathEntry(p string) string {
	p = strings.TrimSpace(p)
	if p == "" {
		return ""
	}
	return strings.ToLower(strings.TrimRight(filepath.Clean(p), `\`))
}

// rebuildPath removes every entry in `remove` (and any existing copy of `add`)
// from the current PATH, then prepends `add` so it takes precedence.
func rebuildPath(current string, remove, add []string) string {
	drop := map[string]bool{}
	for _, list := range [][]string{remove, add} {
		for _, p := range list {
			if n := normalizePathEntry(p); n != "" {
				drop[n] = true
			}
		}
	}

	out := append([]string{}, add...)
	for _, p := range strings.Split(current, ";") {
		n := normalizePathEntry(p)
		if n == "" || drop[n] {
			continue
		}
		out = append(out, p)
	}
	return strings.Join(out, ";")
}

// renderSessionScript produces the script text. The syntax is chosen by the
// file extension: .ps1 => PowerShell, anything else => cmd.exe batch.
func renderSessionScript(file string, vars []sessionVar) string {
	var sb strings.Builder

	if strings.EqualFold(filepath.Ext(file), ".ps1") {
		sb.WriteString("\uFEFF") // BOM so Windows PowerShell 5.1 reads UTF-8 correctly
		for _, v := range vars {
			if v.unset {
				fmt.Fprintf(&sb, "Remove-Item Env:%s -ErrorAction SilentlyContinue\r\n", v.name)
				continue
			}
			fmt.Fprintf(&sb, "$env:%s = '%s'\r\n", v.name, strings.ReplaceAll(v.value, "'", "''"))
		}
		return sb.String()
	}

	sb.WriteString("@echo off\r\n")
	for _, v := range vars {
		if v.unset {
			fmt.Fprintf(&sb, "set \"%s=\"\r\n", v.name)
			continue
		}
		// % must be doubled inside a batch file
		fmt.Fprintf(&sb, "set \"%s=%s\"\r\n", v.name, strings.ReplaceAll(v.value, "%", "%%"))
	}
	return sb.String()
}

// writeSession writes the script to the file named by UNOSDK_ENV_FILE.
func writeSession(vars []sessionVar) error {
	file := os.Getenv("UNOSDK_ENV_FILE")
	if file == "" {
		return errNoSession
	}
	return os.WriteFile(file, []byte(renderSessionScript(file, vars)), 0600)
}

// writeSessionScript makes `sdk` the active SDK of its type in the current shell.
// PATH entries of every other installed SDK of the same type are removed, and
// JAVA_HOME is set for Java SDKs when setJavaHome is true.
//
// `removed` lists SDKs that were just uninstalled (and are therefore no longer
// in the registry) whose PATH entries should also be dropped from the session.
func writeSessionScript(reg *registry.Registry, sdk *models.SDK, setJavaHome bool, removed ...*models.SDK) error {
	if os.Getenv("UNOSDK_ENV_FILE") == "" {
		return errNoSession
	}

	var remove []string
	for _, s := range reg.ListByType(sdk.Type) {
		remove = append(remove, sdkPathEntries(s)...)
	}
	for _, s := range removed {
		if s != nil {
			remove = append(remove, sdkPathEntries(s)...)
		}
	}

	vars := []sessionVar{
		{name: "PATH", value: rebuildPath(os.Getenv("PATH"), remove, sdkPathEntries(sdk))},
	}
	if sdk.Type == models.JavaSDK && setJavaHome {
		vars = append(vars, sessionVar{name: "JAVA_HOME", value: sdk.InstallPath})
	}

	return writeSession(vars)
}

// writeSessionRemoval drops an uninstalled SDK from the current shell without
// activating a replacement. JAVA_HOME is cleared when clearJavaHome is true.
func writeSessionRemoval(sdk *models.SDK, clearJavaHome bool) error {
	if os.Getenv("UNOSDK_ENV_FILE") == "" {
		return errNoSession
	}

	vars := []sessionVar{
		{name: "PATH", value: rebuildPath(os.Getenv("PATH"), sdkPathEntries(sdk), nil)},
	}
	if clearJavaHome {
		vars = append(vars, sessionVar{name: "JAVA_HOME", unset: true})
	}

	return writeSession(vars)
}