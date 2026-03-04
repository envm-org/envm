package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/envm-org/cli/internal/ui"
	"github.com/spf13/cobra"
)

var uninstallCmd = &cobra.Command{
	Use:   "uninstall",
	Short: "Uninstall the envm CLI",
	Run: func(cmd *cobra.Command, args []string) {
		home, err := os.UserHomeDir()
		if err != nil {
			ui.PrintError(fmt.Errorf("getting home directory: %w", err))
			return
		}

		//  Remove the binary
		binPath := filepath.Join(home, ".local", "bin", "envm")
		if _, err := os.Stat(binPath); err == nil {
			ui.PrintInfo(fmt.Sprintf("Removing binary at %s...", binPath))
			if err := os.Remove(binPath); err != nil {
				ui.PrintError(fmt.Errorf("removing binary: %w", err))
				ui.PrintWarning("You may need to remove it manually.")
			} else {
				ui.PrintSuccess("Binary removed successfully.")
			}
		} else {
			ui.PrintInfo("envm binary not found. It might have already been removed.")
		}

		// Remove configuration files
		configPath := filepath.Join(home, ".envm-cli")
		if _, err := os.Stat(configPath); err == nil {
			ui.PrintInfo(fmt.Sprintf("Removing configuration directory at %s...", configPath))
			if err := os.RemoveAll(configPath); err != nil {
				ui.PrintError(fmt.Errorf("removing configuration directory: %w", err))
			} else {
				ui.PrintSuccess("Configuration removed successfully.")
			}
		}

		ui.PrintSuccess("envm CLI uninstalled successfully.")
		ui.PrintWarning("Please note: You may need to remove ~/.local/bin from your PATH manually in your shell config if nothing else uses it.")
	},
}

func init() {
	rootCmd.AddCommand(uninstallCmd)
}
