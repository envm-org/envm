package cmd

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/envm-org/cli/internal/auth"
	"github.com/envm-org/cli/internal/types"
	"github.com/envm-org/cli/internal/ui"
	"github.com/joho/godotenv"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var pushCmd = &cobra.Command{
	Use:   "push",
	Short: "Push local environment variables to the server",
	Run: func(cmd *cobra.Command, args []string) {
		creds, err := auth.LoadCredentials()
		if err != nil || creds == nil || creds.Token == "" {
			ui.PrintError(fmt.Errorf("not authenticated. please log in with 'envm login'"))
			os.Exit(1)
		}

		apiURL := viper.GetString("api-url")

		configFile, err := os.Open("envm.json")
		if err != nil {
			ui.PrintError(fmt.Errorf("envm.json not found. Run 'envm init' first."))
			os.Exit(1)
		}
		defer configFile.Close()

		var config types.ProjectConfig
		if err := json.NewDecoder(configFile).Decode(&config); err != nil {
			ui.PrintError(fmt.Errorf("failed to parse envm.json: %w", err))
			os.Exit(1)
		}

		c := auth.NewClient(apiURL)

		envMap, err := fetchProjectEnvs(c, config.ProjectID, creds.Token)
		if err != nil {
			ui.PrintError(fmt.Errorf("failed to fetch environments: %w", err))
			os.Exit(1)
		}

		envEntries, err := scanEnvFiles()
		if err != nil {
			ui.PrintError(fmt.Errorf("failed to scan for .env files: %w", err))
			os.Exit(1)
		}

		if len(envEntries) == 0 {
			ui.PrintSuccess("No .env files found to push.")
			return
		}

		for _, envEntry := range envEntries {
			ui.PrintSuccess(fmt.Sprintf("\nSyncing '%s' (%s)...", envEntry.Name, envEntry.Path))

			if _, err := os.Stat(envEntry.Path); os.IsNotExist(err) {
				ui.PrintError(fmt.Errorf("file not found: %s", envEntry.Path))
				continue
			}

			envMapFromFile, err := godotenv.Read(envEntry.Path)
			if err != nil {
				ui.PrintError(fmt.Errorf("failed to read .env file: %w", err))
				continue
			}

			envID, exists := envMap[envEntry.Name]
			if !exists {
				ui.PrintSuccess(fmt.Sprintf("Environment '%s' not found on server. Creating...", envEntry.Name))
				newEnvID, err := createEnvironment(c, config.ProjectID, envEntry.Name, envEntry.Name, creds.Token) // Slug = Name for simplicity?
				if err != nil {
					ui.PrintError(fmt.Errorf("failed to create environment: %w", err))
					continue
				}
				envID = newEnvID
				envMap[envEntry.Name] = envID // Update map
			}

			remoteVars, err := fetchVariables(c, envID, creds.Token)
			if err != nil {
				ui.PrintError(fmt.Errorf("failed to fetch variables: %w", err))
				continue
			}

			for key, value := range envMapFromFile {
				if _, exists := remoteVars[key]; exists {
					if err := updateVariable(c, envID, key, value, envEntry.Path, creds.Token); err != nil {
						ui.PrintError(fmt.Errorf("failed to update variable %s: %w", key, err))
					} else {
						fmt.Printf("  Updated %s\n", key)
					}
				} else {
					if err := createVariable(c, envID, key, value, envEntry.Path, creds.Token); err != nil {
						ui.PrintError(fmt.Errorf("failed to create variable %s: %w", key, err))
					} else {
						fmt.Printf("  Created %s\n", key)
					}
				}
			}
		}
	},
}

func init() {
	rootCmd.AddCommand(pushCmd)
}
