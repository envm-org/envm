package cmd

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/envm-org/cli/internal/client"
	"github.com/envm-org/cli/internal/resolver"
	"github.com/envm-org/cli/internal/ui"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var projectCmd = &cobra.Command{
	Use:   "project",
	Short: "Manage projects",
	Long:  `Create, read, update, and delete projects.`,
}

var projectListCmd = &cobra.Command{
	Use:   "list",
	Short: "List all projects in an organization",
	Run: func(cmd *cobra.Command, args []string) {
		apiURL := viper.GetString("api-url")
		c := client.New(apiURL)

		orgID, _ := cmd.Flags().GetString("org-id")

		body, err := c.Get("/project/list?organization_id=" + orgID)
		if err != nil {
			ui.PrintError(err)
			os.Exit(1)
		}

		var projects []map[string]interface{}
		if err := json.Unmarshal(body, &projects); err != nil {
			ui.PrintError(fmt.Errorf("failed to parse response: %w", err))
			os.Exit(1)
		}

		var headers = []string{"Name", "Slug"}
		var data [][]string
		for _, p := range projects {
			data = append(data, []string{
				fmt.Sprintf("%v", p["name"]),
				fmt.Sprintf("%v", p["slug"]),
			})
		}

		ui.RenderTable(headers, data)
	},
}

var projectCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a new project",
	Run: func(cmd *cobra.Command, args []string) {
		apiURL := viper.GetString("api-url")
		c := client.New(apiURL)
		r := resolver.New(c)

		orgSlug, _ := cmd.Flags().GetString("org")
		orgID, err := r.ResolveOrg(orgSlug)
		if err != nil {
			ui.PrintError(err)
			os.Exit(1)
		}

		name, _ := cmd.Flags().GetString("name")
		slug, _ := cmd.Flags().GetString("slug")
		desc, _ := cmd.Flags().GetString("desc")

		payload := map[string]interface{}{
			"organization_id": orgID,
			"name":            name,
			"slug":            slug,
		}
		if desc != "" {
			payload["description"] = desc
		}

		body, err := c.Post("/project", payload)
		if err != nil {
			ui.PrintError(err)
			os.Exit(1)
		}

		var project map[string]string
		json.Unmarshal(body, &project)
		ui.PrintSuccess("Project created successfully!")
		ui.RenderKV("Project Details", project)
	},
}

var projectGetCmd = &cobra.Command{
	Use:   "get [slug]",
	Short: "Get a project by slug",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		projectSlug := args[0]
		orgSlug, _ := cmd.Flags().GetString("org")

		apiURL := viper.GetString("api-url")
		c := client.New(apiURL)
		r := resolver.New(c)

		id, err := r.ResolveProject(orgSlug, projectSlug)
		if err != nil {
			ui.PrintError(err)
			os.Exit(1)
		}

		body, err := c.Get("/project?id=" + id)
		if err != nil {
			ui.PrintError(err)
			os.Exit(1)
		}

		var project map[string]string
		json.Unmarshal(body, &project)
		ui.RenderKV("Project Details", project)
	},
}

var projectUpdateCmd = &cobra.Command{
	Use:   "update [slug]",
	Short: "Update a project",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		projectSlug := args[0]
		orgSlug, _ := cmd.Flags().GetString("org")

		apiURL := viper.GetString("api-url")
		c := client.New(apiURL)
		r := resolver.New(c)

		id, err := r.ResolveProject(orgSlug, projectSlug)
		if err != nil {
			ui.PrintError(err)
			os.Exit(1)
		}

		name, _ := cmd.Flags().GetString("name")
		slug, _ := cmd.Flags().GetString("slug")
		desc, _ := cmd.Flags().GetString("desc")

		payload := map[string]interface{}{
			"id": id,
		}
		if name != "" {
			payload["name"] = name
		}
		if slug != "" {
			payload["slug"] = slug
		}
		if desc != "" {
			payload["description"] = desc
		}

		body, err := c.Put("/project?id="+id, payload)
		if err != nil {
			ui.PrintError(err)
			os.Exit(1)
		}

		var project map[string]string
		json.Unmarshal(body, &project)
		ui.PrintSuccess("Project updated successfully!")
		ui.RenderKV("Updated Project Details", project)
	},
}

var projectDeleteCmd = &cobra.Command{
	Use:   "delete [slug]",
	Short: "Delete a project",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		projectSlug := args[0]
		orgSlug, _ := cmd.Flags().GetString("org")

		apiURL := viper.GetString("api-url")
		c := client.New(apiURL)
		r := resolver.New(c)

		id, err := r.ResolveProject(orgSlug, projectSlug)
		if err != nil {
			ui.PrintError(err)
			os.Exit(1)
		}

		_, err = c.Delete("/project?id=" + id)
		if err != nil {
			ui.PrintError(err)
			os.Exit(1)
		}

		ui.PrintSuccess("Project deleted successfully!")
	},
}

func init() {
	rootCmd.AddCommand(projectCmd)
	projectCmd.AddCommand(projectListCmd)
	projectCmd.AddCommand(projectCreateCmd)
	projectCmd.AddCommand(projectGetCmd)
	projectCmd.AddCommand(projectUpdateCmd)
	projectCmd.AddCommand(projectDeleteCmd)

	projectListCmd.Flags().String("org", "", "Organization slug")
	projectListCmd.MarkFlagRequired("org")

	projectCreateCmd.Flags().String("org", "", "Organization slug")
	projectCreateCmd.Flags().String("name", "", "Project name")
	projectCreateCmd.Flags().String("slug", "", "Project slug")
	projectCreateCmd.Flags().String("desc", "", "Project description")
	projectCreateCmd.MarkFlagRequired("org")
	projectCreateCmd.MarkFlagRequired("name")
	projectCreateCmd.MarkFlagRequired("slug")

	projectGetCmd.Flags().String("org", "", "Organization slug")
	projectGetCmd.MarkFlagRequired("org")

	projectUpdateCmd.Flags().String("org", "", "Organization slug")
	projectUpdateCmd.Flags().String("name", "", "Project name")
	projectUpdateCmd.Flags().String("slug", "", "Project slug")
	projectUpdateCmd.Flags().String("desc", "", "Project description")
	projectUpdateCmd.MarkFlagRequired("org")

	projectDeleteCmd.Flags().String("org", "", "Organization slug")
	projectDeleteCmd.MarkFlagRequired("org")
}
