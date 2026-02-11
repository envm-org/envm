package cmd

import (
	"fmt"
	"os"

	"github.com/envm-org/cli/internal/auth"
	"github.com/envm-org/cli/internal/ui"
	"github.com/spf13/cobra"
)

var whoamiCmd = &cobra.Command{
	Use:   "whoami",
	Short: "Show current logged in user",
	Long:  `Display the currently authenticated user's name and email.`,
	Run: func(cmd *cobra.Command, args []string) {
		creds, err := auth.LoadCredentials()
		if err != nil {
			// If error is just "no file", treat as not logged in
			if os.IsNotExist(err) {
				ui.PrintError(fmt.Errorf("not logged in"))
				return
			}
			// If other error (e.g. corruption), show it
			ui.PrintError(fmt.Errorf("error loading credentials: %w", err))
			return
		}

		if creds == nil || creds.Token == "" {
			ui.PrintError(fmt.Errorf("not logged in"))
			return
		}

		// Display user info
		ui.PrintLogo()
		ui.RenderKV("Current User", map[string]string{
			"Name":  creds.FullName,
			"Email": creds.Email,
			"ID":    creds.UserID,
		})
	},
}

func init() {
	rootCmd.AddCommand(whoamiCmd)
}
