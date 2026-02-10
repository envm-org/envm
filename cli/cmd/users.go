package cmd

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/envm-org/cli/internal/client"
	"github.com/envm-org/cli/internal/resolver"
	"github.com/envm-org/cli/internal/types"
	"github.com/envm-org/cli/internal/ui"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var usersCmd = &cobra.Command{
	Use:   "users",
	Short: "Manage users",
	Long:  `Create, read, update, and delete users.`,
}

var usersListCmd = &cobra.Command{
	Use:   "list",
	Short: "List all users",
	Run: func(cmd *cobra.Command, args []string) {
		apiURL := viper.GetString("api-url")
		c := client.New(apiURL)

		body, err := c.Get("/users/list")
		if err != nil {
			ui.PrintError(err)
			os.Exit(1)
		}

		var users []types.User
		if err := json.Unmarshal(body, &users); err != nil {
			ui.PrintError(fmt.Errorf("failed to parse response: %w", err))
			os.Exit(1)
		}

		var headers = []string{"Name", "Email"}
		var data [][]string
		for _, u := range users {
			data = append(data, []string{
				u.FullName,
				u.Email,
			})
		}

		ui.RenderTable(headers, data)
	},
}

var usersCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a new user",
	Run: func(cmd *cobra.Command, args []string) {
		apiURL := viper.GetString("api-url")
		c := client.New(apiURL)

		email, _ := cmd.Flags().GetString("email")
		name, _ := cmd.Flags().GetString("name")
		password, _ := cmd.Flags().GetString("password")

		payload := map[string]string{
			"email":         email,
			"full_name":     name,
			"password_hash": password,
		}

		body, err := c.Post("/users", payload)
		if err != nil {
			ui.PrintError(err)
			os.Exit(1)
		}

		var user map[string]string
		json.Unmarshal(body, &user)
		ui.PrintSuccess("User created successfully!")
		ui.RenderKV("User Details", user)
	},
}

var usersGetCmd = &cobra.Command{
	Use:   "get [email]",
	Short: "Get a user by email",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		email := args[0]
		apiURL := viper.GetString("api-url")
		c := client.New(apiURL)
		r := resolver.New(c)

		id, err := r.ResolveUser(email)
		if err != nil {
			ui.PrintError(err)
			os.Exit(1)
		}

		body, err := c.Get("/users?id=" + id)
		if err != nil {
			ui.PrintError(err)
			os.Exit(1)
		}

		var user map[string]string
		json.Unmarshal(body, &user)
		ui.RenderKV("User Details", user)
	},
}

var usersUpdateCmd = &cobra.Command{
	Use:   "update [email]",
	Short: "Update a user",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		emailArg := args[0]
		apiURL := viper.GetString("api-url")
		c := client.New(apiURL)
		r := resolver.New(c)

		id, err := r.ResolveUser(emailArg)
		if err != nil {
			ui.PrintError(err)
			os.Exit(1)
		}

		name, _ := cmd.Flags().GetString("name")
		email, _ := cmd.Flags().GetString("email")
		password, _ := cmd.Flags().GetString("password")

		payload := map[string]interface{}{
			"id": id,
		}
		if name != "" {
			payload["full_name"] = name
		}
		if email != "" {
			payload["email"] = email
		}
		if password != "" {
			payload["password_hash"] = password
		}

		body, err := c.Put("/users?id="+id, payload)
		if err != nil {
			ui.PrintError(err)
			os.Exit(1)
		}

		var user map[string]string
		json.Unmarshal(body, &user)
		ui.PrintSuccess("User updated successfully!")
		ui.RenderKV("Updated User Details", user)
	},
}

var usersDeleteCmd = &cobra.Command{
	Use:   "delete [email]",
	Short: "Delete a user",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		email := args[0]
		apiURL := viper.GetString("api-url")
		c := client.New(apiURL)
		r := resolver.New(c)

		id, err := r.ResolveUser(email)
		if err != nil {
			ui.PrintError(err)
			os.Exit(1)
		}

		_, err = c.Delete("/users?id=" + id)
		if err != nil {
			ui.PrintError(err)
			os.Exit(1)
		}

		ui.PrintSuccess("User deleted successfully!")
	},
}

func init() {
	rootCmd.AddCommand(usersCmd)
	usersCmd.AddCommand(usersListCmd)
	usersCmd.AddCommand(usersCreateCmd)
	usersCmd.AddCommand(usersGetCmd)
	usersCmd.AddCommand(usersUpdateCmd)
	usersCmd.AddCommand(usersDeleteCmd)

	usersCreateCmd.Flags().String("email", "", "User email")
	usersCreateCmd.Flags().String("name", "", "User full name")
	usersCreateCmd.Flags().String("password", "", "User password")
	usersCreateCmd.MarkFlagRequired("email")
	usersCreateCmd.MarkFlagRequired("name")
	usersCreateCmd.MarkFlagRequired("password")

	usersUpdateCmd.Flags().String("email", "", "User email")
	usersUpdateCmd.Flags().String("name", "", "User full name")
	usersUpdateCmd.Flags().String("password", "", "User password")

}
