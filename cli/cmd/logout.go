package cmd

import (
	"fmt"

	"github.com/envm-org/cli/internal/auth"
	"github.com/envm-org/cli/internal/ui"
	"github.com/spf13/cobra"
)

var logoutCmd = &cobra.Command{
	Use:   "logout",
	Short: "Log out from envm",
	Long:  `Remove your local session credentials.`,
	Run: func(cmd *cobra.Command, args []string) {
		if err := auth.Logout(); err != nil {
			ui.PrintError(fmt.Errorf("failed to logout: %w", err))
			return
		}

		ui.PrintSuccess("Logged out successfully.")
	},
}

func init() {
	rootCmd.AddCommand(logoutCmd)
}
