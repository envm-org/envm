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

var orgCmd = &cobra.Command{
	Use:   "org",
	Short: "Manage organizations",
	Long:  `Create, read, update, and delete organizations.`,
}

var orgListCmd = &cobra.Command{
	Use:   "list",
	Short: "List all organizations",
	Run: func(cmd *cobra.Command, args []string) {
		apiURL := viper.GetString("api-url")
		c := client.New(apiURL)

		body, err := c.Get("/org/list")
		if err != nil {
			ui.PrintError(err)
			os.Exit(1)
		}

		var orgs []types.Organization
		if err := json.Unmarshal(body, &orgs); err != nil {
			ui.PrintError(fmt.Errorf("failed to parse response: %w", err))
			os.Exit(1)
		}

		var headers = []string{"Name", "Slug"}
		var data [][]string
		for _, org := range orgs {
			data = append(data, []string{
				org.Name,
				org.Slug,
			})
		}

		ui.RenderTable(headers, data)
	},
}

var orgCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a new organization",
	Run: func(cmd *cobra.Command, args []string) {
		apiURL := viper.GetString("api-url")
		c := client.New(apiURL)

		name, _ := cmd.Flags().GetString("name")
		slug, _ := cmd.Flags().GetString("slug")

		payload := map[string]string{
			"name": name,
			"slug": slug,
		}

		body, err := c.Post("/org", payload)
		if err != nil {
			ui.PrintError(err)
			os.Exit(1)
		}

		var org map[string]string
		json.Unmarshal(body, &org)
		ui.PrintSuccess("Organization created successfully!")
		ui.RenderKV("Organization Details", org)
	},
}

var orgGetCmd = &cobra.Command{
	Use:   "get [slug]",
	Short: "Get an organization by slug",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		slug := args[0]
		apiURL := viper.GetString("api-url")
		c := client.New(apiURL)
		r := resolver.New(c)

		id, err := r.ResolveOrg(slug)
		if err != nil {
			ui.PrintError(err)
			os.Exit(1)
		}

		body, err := c.Get("/org?id=" + id)
		if err != nil {
			ui.PrintError(err)
			os.Exit(1)
		}

		var org map[string]string
		json.Unmarshal(body, &org)
		ui.RenderKV("Organization Details", org)
	},
}

var orgUpdateCmd = &cobra.Command{
	Use:   "update [slug]",
	Short: "Update an organization",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		slugArg := args[0]
		apiURL := viper.GetString("api-url")
		c := client.New(apiURL)
		r := resolver.New(c)

		id, err := r.ResolveOrg(slugArg)
		if err != nil {
			ui.PrintError(err)
			os.Exit(1)
		}

		name, _ := cmd.Flags().GetString("name")
		slug, _ := cmd.Flags().GetString("slug")

		payload := map[string]interface{}{
			"id": id,
		}
		if name != "" {
			payload["name"] = name
		}
		if slug != "" {
			payload["slug"] = slug
		}

		body, err := c.Put("/org?id="+id, payload)
		if err != nil {
			ui.PrintError(err)
			os.Exit(1)
		}

		var org map[string]string
		json.Unmarshal(body, &org)
		ui.PrintSuccess("Organization updated successfully!")
		ui.RenderKV("Updated Organization Details", org)
	},
}

var orgDeleteCmd = &cobra.Command{
	Use:   "delete [slug]",
	Short: "Delete an organization",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		slugArg := args[0]
		apiURL := viper.GetString("api-url")
		c := client.New(apiURL)
		r := resolver.New(c)

		id, err := r.ResolveOrg(slugArg)
		if err != nil {
			ui.PrintError(err)
			os.Exit(1)
		}

		_, err = c.Delete("/org?id=" + id)
		if err != nil {
			ui.PrintError(err)
			os.Exit(1)
		}

		ui.PrintSuccess("Organization deleted successfully!")
	},
}

func init() {
	rootCmd.AddCommand(orgCmd)
	orgCmd.AddCommand(orgListCmd)
	orgCmd.AddCommand(orgCreateCmd)
	orgCmd.AddCommand(orgGetCmd)
	orgCmd.AddCommand(orgUpdateCmd)
	orgCmd.AddCommand(orgDeleteCmd)

	orgCreateCmd.Flags().String("name", "", "Organization name")
	orgCreateCmd.Flags().String("slug", "", "Organization slug")
	orgCreateCmd.MarkFlagRequired("name")
	orgCreateCmd.MarkFlagRequired("slug")

	orgUpdateCmd.Flags().String("name", "", "Organization name")
	orgUpdateCmd.Flags().String("slug", "", "Organization slug")
}
