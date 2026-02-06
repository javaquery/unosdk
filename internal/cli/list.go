package cli

import (
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/spf13/cobra"
	"github.com/javaquery/unosdk/internal/providers"
	"github.com/javaquery/unosdk/internal/providers/flutter"
	"github.com/javaquery/unosdk/internal/providers/gradle"
	"github.com/javaquery/unosdk/internal/providers/java"
	"github.com/javaquery/unosdk/internal/providers/maven"
	"github.com/javaquery/unosdk/internal/providers/node"
	"github.com/javaquery/unosdk/internal/providers/python"
	"github.com/javaquery/unosdk/internal/registry"
)

var (
	showInstalled bool
	showAvailable bool
)

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List available or installed SDKs",
	Long: `List available providers and their versions, or show installed SDKs.

Examples:
  # List all available providers
  unosdk list

  # List installed SDKs
  unosdk list --installed

  # List available versions (not yet installed)
  unosdk list --available`,
	RunE: runList,
}

func init() {
	listCmd.Flags().BoolVarP(&showInstalled, "installed", "i", false, "Show installed SDKs")
	listCmd.Flags().BoolVarP(&showAvailable, "available", "a", false, "Show available providers and versions")
}

func runList(cmd *cobra.Command, args []string) error {
	if showInstalled {
		return listInstalled()
	}

	if showAvailable {
		return listAvailable()
	}

	// Default: show both
	fmt.Println("=== Available Providers ===")
	if err := listAvailable(); err != nil {
		return err
	}

	fmt.Println("\n=== Installed SDKs ===")
	return listInstalled()
}

func listInstalled() error {
	reg, err := registry.NewRegistry()
	if err != nil {
		return fmt.Errorf("failed to initialize registry: %w", err)
	}

	sdks := reg.List()
	
	if len(sdks) == 0 {
		fmt.Println("No SDKs installed.")
		return nil
	}

	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	fmt.Fprintln(w, "TYPE\tPROVIDER\tVERSION\tINSTALL PATH")
	fmt.Fprintln(w, "----\t--------\t-------\t------------")
	
	for _, sdk := range sdks {
		fmt.Fprintf(w, "%s\t%s\t%s\t%s\n", sdk.Type, sdk.Provider, sdk.Version, sdk.InstallPath)
	}
	
	return w.Flush()
}

func listAvailable() error {
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

	allProviders := providerRegistry.ListAll()

	if len(allProviders) == 0 {
		fmt.Println("No providers available.")
		return nil
	}

	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	fmt.Fprintln(w, "TYPE\tPROVIDER\tDISPLAY NAME")
	fmt.Fprintln(w, "----\t--------\t------------")
	
	for _, provider := range allProviders {
		fmt.Fprintf(w, "%s\t%s\t%s\n", provider.Type(), provider.Name(), provider.DisplayName())
	}
	
	w.Flush()

	fmt.Println("\nExample usage:")
	fmt.Println("  unosdk install java amazoncorretto 21")
	fmt.Println("  unosdk install node nodejs latest")
	fmt.Println("  unosdk install python python 3.12.1")
	fmt.Println("  unosdk install flutter flutter latest")

	return nil
}
