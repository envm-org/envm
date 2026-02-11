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

var envCmd = &cobra.Command{
	Use:   "env",
	Short: "Manage environments",
	Long:  `Create, read, update, and delete environments.`,
}

var envListCmd = &cobra.Command{
	Use:   "list",
	Short: "List all environments in a project",
	Run: func(cmd *cobra.Command, args []string) {
		projectSlug, _ := cmd.Flags().GetString("project")
		if projectSlug == "" {
			ui.PrintError(fmt.Errorf("project slug is required"))
			os.Exit(1)
		}

		apiURL := viper.GetString("api-url")
		c := client.New(apiURL)
		r := resolver.New(c)

		projectID, err := r.ResolveProject(projectSlug)
		if err != nil {
			ui.PrintError(err)
			os.Exit(1)
		}

		body, err := c.Get("/env/list?project_id=" + projectID)
		if err != nil {
			ui.PrintError(err)
			os.Exit(1)
		}

		var envs []map[string]interface{}
		if err := json.Unmarshal(body, &envs); err != nil {
			ui.PrintError(fmt.Errorf("failed to parse response: %w", err))
			os.Exit(1)
		}

		var headers = []string{"Name", "Slug"}
		var data [][]string
		for _, e := range envs {
			data = append(data, []string{
				fmt.Sprintf("%v", e["name"]),
				fmt.Sprintf("%v", e["slug"]),
			})
		}

		ui.RenderTable(headers, data)
	},
}

var envCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a new environment",
	Run: func(cmd *cobra.Command, args []string) {
		apiURL := viper.GetString("api-url")
		c := client.New(apiURL)
		r := resolver.New(c)

		projectSlug, _ := cmd.Flags().GetString("project")
		projectID, err := r.ResolveProject(projectSlug)
		if err != nil {
			ui.PrintError(err)
			os.Exit(1)
		}

		name, _ := cmd.Flags().GetString("name")
		slug, _ := cmd.Flags().GetString("slug")

		payload := map[string]string{
			"project_id": projectID,
			"name":       name,
			"slug":       slug,
		}

		body, err := c.Post("/env", payload)
		if err != nil {
			ui.PrintError(err)
			os.Exit(1)
		}

		var env map[string]string
		json.Unmarshal(body, &env)
		ui.PrintSuccess("Environment created successfully!")
		ui.RenderKV("Environment Details", env)
	},
}

var envGetCmd = &cobra.Command{
	Use:   "get [slug]",
	Short: "Get an environment by slug",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		envSlug := args[0]
		projectSlug, _ := cmd.Flags().GetString("project")

		apiURL := viper.GetString("api-url")
		c := client.New(apiURL)
		r := resolver.New(c)

		id, err := r.ResolveEnv(projectSlug, envSlug)
		if err != nil {
			ui.PrintError(err)
			os.Exit(1)
		}

		body, err := c.Get("/env?id=" + id)
		if err != nil {
			ui.PrintError(err)
			os.Exit(1)
		}

		var env map[string]string
		json.Unmarshal(body, &env)
		ui.RenderKV("Environment Details", env)
	},
}

var envUpdateCmd = &cobra.Command{
	Use:   "update [slug]",
	Short: "Update an environment",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		envSlug := args[0]
		projectSlug, _ := cmd.Flags().GetString("project")

		apiURL := viper.GetString("api-url")
		c := client.New(apiURL)
		r := resolver.New(c)

		id, err := r.ResolveEnv(projectSlug, envSlug)
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

		body, err := c.Put("/env?id="+id, payload)
		if err != nil {
			ui.PrintError(err)
			os.Exit(1)
		}

		var env map[string]string
		json.Unmarshal(body, &env)
		ui.PrintSuccess("Environment updated successfully!")
		ui.RenderKV("Updated Environment Details", env)
	},
}

var envDeleteCmd = &cobra.Command{
	Use:   "delete [slug]",
	Short: "Delete an environment",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		envSlug := args[0]
		projectSlug, _ := cmd.Flags().GetString("project")

		apiURL := viper.GetString("api-url")
		c := client.New(apiURL)
		r := resolver.New(c)

		id, err := r.ResolveEnv(projectSlug, envSlug)
		if err != nil {
			ui.PrintError(err)
			os.Exit(1)
		}

		_, err = c.Delete("/env?id=" + id)
		if err != nil {
			ui.PrintError(err)
			os.Exit(1)
		}

		ui.PrintSuccess("Environment deleted successfully!")
	},
}

func init() {
	rootCmd.AddCommand(envCmd)
	envCmd.AddCommand(envListCmd)
	envCmd.AddCommand(envCreateCmd)
	envCmd.AddCommand(envGetCmd)
	envCmd.AddCommand(envUpdateCmd)
	envCmd.AddCommand(envDeleteCmd)

	envListCmd.Flags().String("project", "", "Project slug")
	envListCmd.MarkFlagRequired("project")

	envCreateCmd.Flags().String("project", "", "Project slug")
	envCreateCmd.Flags().String("name", "", "Environment name")
	envCreateCmd.Flags().String("slug", "", "Environment slug")
	envCreateCmd.MarkFlagRequired("project")
	envCreateCmd.MarkFlagRequired("name")
	envCreateCmd.MarkFlagRequired("slug")

	envGetCmd.Flags().String("project", "", "Project slug")
	envGetCmd.MarkFlagRequired("project")

	envUpdateCmd.Flags().String("project", "", "Project slug")
	envUpdateCmd.Flags().String("name", "", "Environment name")
	envUpdateCmd.Flags().String("slug", "", "Environment slug")
	envUpdateCmd.MarkFlagRequired("project")

	envDeleteCmd.Flags().String("project", "", "Project slug")
	envDeleteCmd.MarkFlagRequired("project")
}
