package cmd

import (
	"fmt"

	"github.com/envm-org/cli/internal/auth"
	"github.com/envm-org/cli/internal/ui"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var registerCmd = &cobra.Command{
	Use:   "register",
	Short: "Create a new envm account",
	Long:  `Create a new account with your name, email and password.`,
	Run: func(cmd *cobra.Command, args []string) {
		ui.PrintLogo()

		fullName := ui.Prompt("Full Name")
		if fullName == "" {
			ui.PrintError(fmt.Errorf("name is required"))
			return
		}

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

		ui.PrintSuccess("Registering...")

		apiURL := viper.GetString("api-url")
		client := auth.NewClient(apiURL)

		//  Register logic
		user, err := client.Register(fullName, email, password)
		if err != nil {
			ui.PrintError(err)
			return
		}
		ui.PrintSuccess("Registration successful!")

		// Auto-login logic
		ui.PrintSuccess("Logging you in...")
		loginResp, err := client.Login(email, password)
		if err != nil {
			ui.PrintError(fmt.Errorf("failed to auto-login: %w", err))
			ui.PrintSuccess("Please run 'envm login' to sign in.")
			return
		}

		// Save credentials
		creds := auth.Credentials{
			Token:    loginResp.AccessToken,
			UserID:   user.ID,
			Email:    user.Email,
			FullName: user.FullName,
		}

		if err := auth.SaveCredentials(creds); err != nil {
			ui.PrintError(fmt.Errorf("failed to save credentials: %w", err))
			return
		}

		ui.PrintSuccess(fmt.Sprintf("Welcome to envm, %s!", user.FullName))
	},
}

func init() {
	rootCmd.AddCommand(registerCmd)
}
