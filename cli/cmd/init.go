package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/envm-org/cli/internal/auth"
	"github.com/envm-org/cli/internal/types"
	"github.com/envm-org/cli/internal/ui"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize envm in the current directory",
	Long:  `Scans for .env files and creates an envm.json configuration file.`,
	Run: func(cmd *cobra.Command, args []string) {
		ui.PrintLogo()

		// Check Authentication
		creds, err := auth.LoadCredentials()
		if err != nil || creds == nil || creds.Token == "" {
			ui.PrintError(fmt.Errorf("not authenticated. please log in with 'envm login'"))
			os.Exit(1)
		}

		apiURL := viper.GetString("api-url")
		client := auth.NewClient(apiURL)

		// Prompt for Project Name
		cwd, _ := os.Getwd()
		defaultProjectName := filepath.Base(cwd)
		projectName := ui.Prompt(fmt.Sprintf("Project Name (default: %s)", defaultProjectName))
		if projectName == "" {
			projectName = defaultProjectName
		}
		projectSlug := strings.ToLower(strings.ReplaceAll(projectName, " ", "-"))

		// Prompt for Organization Name
		// TODO: Ideally list existing orgs, but for now prompt to create/use
		defaultOrgName := creds.FullName
		orgName := ui.Prompt(fmt.Sprintf("Organization Name (default: %s)", defaultOrgName))
		if orgName == "" {
			orgName = defaultOrgName
		}
		orgSlug := strings.ToLower(strings.ReplaceAll(orgName, " ", "-"))

		ui.PrintSuccess("Setting up project...")

		// Create/Get Organization
		// Improve: Check if exists first? For now try create, if fails maybe it exists?
		// The current API doesn't easily support "get by slug" for orgs without ID.
		// So we try to create. If 409 conflict, we might need a way to find it or ask user.
		// For this iteration, let's assume valid new org or unique.
		// Actually, to make it robust, we should probably ListOrgs and see if one matches.
		// But let's stick to the plan: Call CreateOrganization.
		org, err := client.CreateOrganization(orgName, orgSlug, creds.Token)
		if err != nil {
			// If error, print but maybe verify if it is "already exists" (not handled in client yet)
			ui.PrintError(fmt.Errorf("failed to create organization: %w", err))
			return
		}
		ui.PrintSuccess(fmt.Sprintf("Organization '%s' ready.", org.Name))

		// 5. Create Project
		project, err := client.CreateProject(org.ID, projectName, projectSlug, "Created via CLI", creds.Token)
		if err != nil {
			ui.PrintError(fmt.Errorf("failed to create project: %w", err))
			return
		}
		ui.PrintSuccess(fmt.Sprintf("Project '%s' created.", project.Name))

		// 6. Scan for .env files
		var envs []types.EnvEntry
		filepath.Walk(".", func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			if info.IsDir() && path != "." {
				if strings.HasPrefix(info.Name(), ".") {
					return filepath.SkipDir
				}
				return nil
			}

			if strings.HasPrefix(info.Name(), ".env") {
				name := strings.TrimPrefix(info.Name(), ".env")
				if name == "" || name == "." {
					name = "default"
				} else if strings.HasPrefix(name, ".") {
					name = strings.TrimPrefix(name, ".")
				}

				envs = append(envs, types.EnvEntry{
					Name: name,
					Path: path,
				})
			}
			return nil
		})

		if len(envs) > 0 {
			ui.PrintSuccess(fmt.Sprintf("Found %d .env files.", len(envs)))
		}

		// 7. Generate envm.json
		config := types.ProjectConfig{
			OwnerID:   org.ID,
			ProjectID: project.ID,
			Envs:      envs,
			// Credentials optional/empty by default now
			Credentials: []types.CredentialEntry{},
		}

		file, err := os.Create("envm.json")
		if err != nil {
			ui.PrintError(fmt.Errorf("failed to create envm.json: %w", err))
			os.Exit(1)
		}
		defer file.Close()

		encoder := json.NewEncoder(file)
		encoder.SetIndent("", "    ")
		if err := encoder.Encode(config); err != nil {
			ui.PrintError(fmt.Errorf("failed to write config: %w", err))
			os.Exit(1)
		}

		ui.PrintSuccess("Initialized envm.json successfully!")
		ui.PrintSuccess("You can now add credentials manually to envm.json if needed.")
	},
}

func init() {
	rootCmd.AddCommand(initCmd)
}
