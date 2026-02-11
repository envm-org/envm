package cmd

import (
	"fmt"

	"github.com/envm-org/cli/internal/auth"
	"github.com/envm-org/cli/internal/ui"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var loginCmd = &cobra.Command{
	Use:   "login",
	Short: "Log in to envm",
	Long:  `Log in with your email and password to access envm services.`,
	Run: func(cmd *cobra.Command, args []string) {
		ui.PrintLogo()

		email := ui.Prompt("Email")
		if email == "" {
			ui.PrintError(fmt.Errorf("email is required"))
			return
		}

		password := ui.PromptPassword("Password")
		if password == "" {
			ui.PrintError(fmt.Errorf("password is required"))
			return
		}

		ui.PrintSuccess("Authenticating...")

		apiURL := viper.GetString("api-url")
		client := auth.NewClient(apiURL)

		resp, err := client.Login(email, password)
		if err != nil {
			ui.PrintError(err)
			return
		}

		creds := auth.Credentials{
			Token:    resp.AccessToken,
			UserID:   resp.User.ID,
			Email:    resp.User.Email,
			FullName: resp.User.FullName,
		}

		if err := auth.SaveCredentials(creds); err != nil {
			ui.PrintError(fmt.Errorf("failed to save credentials: %w", err))
			return
		}

		ui.PrintSuccess(fmt.Sprintf("Welcome back, %s!", resp.User.FullName))
	},
}

func init() {
	rootCmd.AddCommand(loginCmd)
}
