package main

import (
	"fmt"
	"os"

	"github.com/javaquery/unosdk/internal/cli"
	"github.com/javaquery/unosdk/pkg/utils"
)

var (
	version = "dev"
	commit  = "none"
	date    = "unknown"
)

func main() {
	// Initialize logger
	logger := utils.NewLogger()
	defer logger.Sync()

	// Set version info
	cli.SetVersionInfo(version, commit, date)

	// Execute CLI
	if err := cli.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
