package cmd

import (
	"fmt"
	"os"

	"github.com/envm-org/cli/internal/ui"
	"github.com/spf13/cobra"
)

var loadCmd = &cobra.Command{
	Use:   "load",
	Short: "Scan for new .env files and update envm.json",
	Run: func(cmd *cobra.Command, args []string) {
		envEntries, err := scanEnvFiles()
		if err != nil {
			ui.PrintError(fmt.Errorf("failed to scan for .env files: %w", err))
			os.Exit(1)
		}

		if len(envEntries) == 0 {
			ui.PrintSuccess("No .env files found.")
			return
		}

		ui.PrintSuccess(fmt.Sprintf("Found %d .env files:", len(envEntries)))
		for _, e := range envEntries {
			fmt.Printf("  + %s (%s)\n", e.Name, e.Path)
		}
		ui.PrintSuccess("Run 'envm push' to sync these to the server.")
	},
}

func init() {
	rootCmd.AddCommand(loadCmd)
}
