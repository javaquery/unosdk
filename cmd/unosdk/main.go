package main

import (
	"fmt"
	"os"

	"github.com/javaquery/unosdk/internal/cli"
	"github.com/javaquery/unosdk/pkg/utils"
	"github.com/javaquery/unosdk/pkg/version"
)

func main() {
	// Initialize logger
	logger := utils.NewLogger()
	defer logger.Sync()

	// Set version info from version package
	cli.SetVersionInfo(version.Version, version.GitCommit, version.BuildDate)

	// Execute CLI
	if err := cli.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
