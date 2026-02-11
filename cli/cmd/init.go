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

		// Check if envm.json already exists in current or parent directories
		currentDir, err := os.Getwd()
		if err != nil {
			ui.PrintError(fmt.Errorf("failed to get current directory: %w", err))
			os.Exit(1)
		}

		// Walk up the directory tree to check for existing envm.json
		dir := currentDir
		for {
			configPath := filepath.Join(dir, "envm.json")
			if _, err := os.Stat(configPath); err == nil {
				ui.PrintError(fmt.Errorf("envm.json already exists at %s", configPath))
				ui.PrintError(fmt.Errorf("cannot initialize in a directory that already has an envm project"))
				os.Exit(1)
			}

			// Move to parent directory
			parent := filepath.Dir(dir)
			if parent == dir {
				// Reached root directory
				break
			}
			dir = parent
		}

		creds, err := auth.LoadCredentials()
		if err != nil || creds == nil || creds.Token == "" {
			ui.PrintError(fmt.Errorf("not authenticated. please log in with 'envm login'"))
			os.Exit(1)
		}

		apiURL := viper.GetString("api-url")
		client := auth.NewClient(apiURL)

		if err := client.Validate(creds.Token); err != nil {
			if err.Error() == "unauthorized" || strings.Contains(err.Error(), "unauthorized") {
				ui.PrintError(fmt.Errorf("Session expired or invalid token. Please run 'envm login' to authenticate."))
			} else {
				ui.PrintError(fmt.Errorf("failed to validate session: %w", err))
			}
			os.Exit(1)
		}

		cwd, _ := os.Getwd()
		defaultProjectName := filepath.Base(cwd)
		projectName := ui.Prompt(fmt.Sprintf("Project Name (default: %s)", defaultProjectName))
		if projectName == "" {
			projectName = defaultProjectName
		}
		projectSlug := strings.ToLower(strings.ReplaceAll(projectName, " ", "-"))

		ui.PrintSuccess("Setting up project...")

		project, err := client.CreateProject(projectName, projectSlug, "Created via CLI", creds.Token)
		if err != nil {
			ui.PrintError(fmt.Errorf("failed to create project: %w", err))
			return
		}
		ui.PrintSuccess(fmt.Sprintf("Project '%s' created.", project.Name))

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

		// Generate envm.json
		config := types.ProjectConfig{
			ProjectID: project.ID,
			// Envs is dynamic now
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
