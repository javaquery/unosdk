package cli

import (
	"fmt"

	"github.com/spf13/cobra"
)

var updateCmd = &cobra.Command{
	Use:   "update [sdk]",
	Short: "Update an installed SDK to the latest version",
	Long:  `Update an installed SDK to the latest version or a specific version.`,
	Args:  cobra.MinimumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		sdkName := args[0]
		version, _ := cmd.Flags().GetString("version")

		fmt.Printf("Updating %s...\n", sdkName)
		if version != "" {
			fmt.Printf("Target version: %s\n", version)
		} else {
			fmt.Println("Target version: latest")
		}

		// TODO: Implement update logic
		fmt.Println("Update functionality coming soon!")
		return nil
	},
}

func init() {
	updateCmd.Flags().StringP("version", "v", "", "Specific version to update to")
	rootCmd.AddCommand(updateCmd)
}
