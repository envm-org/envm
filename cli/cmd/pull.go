package cmd

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"path/filepath"

	"github.com/envm-org/cli/internal/auth"
	"github.com/envm-org/cli/internal/types"
	"github.com/envm-org/cli/internal/ui"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var pullCmd = &cobra.Command{
	Use:   "pull",
	Short: "Pull environment variables from the server",
	Run: func(cmd *cobra.Command, args []string) {
		creds, err := auth.LoadCredentials()
		if err != nil || creds == nil || creds.Token == "" {
			ui.PrintError(fmt.Errorf("not authenticated. please log in with 'envm login'"))
			os.Exit(1)
		}

		apiURL := viper.GetString("api-url")
		c := auth.NewClient(apiURL)

		// Check if envm.json exists
		var config types.ProjectConfig
		configFile, err := os.Open("envm.json")

		if err != nil {
			if os.IsNotExist(err) {
				// Interactive Setup
				ui.PrintSuccess("envm.json not found. Select a project to pull:")
				projects, err := c.ListProjects(creds.Token)
				if err != nil {
					ui.PrintError(fmt.Errorf("failed to list projects: %w", err))
					os.Exit(1)
				}

				if len(projects) == 0 {
					ui.PrintError(fmt.Errorf("no projects found"))
					os.Exit(1)
				}

				for i, p := range projects {
					fmt.Printf("[%d] %s (%s)\n", i+1, p.Name, p.Slug)
				}

				idxStr := ui.Prompt("Select Project (number)")
				var idx int
				fmt.Sscanf(idxStr, "%d", &idx)

				if idx < 1 || idx > len(projects) {
					ui.PrintError(fmt.Errorf("invalid selection"))
					os.Exit(1)
				}
				selectedProject := projects[idx-1]

				config = types.ProjectConfig{
					ProjectID:   selectedProject.ID,
					Credentials: []types.CredentialEntry{},
				}

				// Save minimal envm.json
				file, err := os.Create("envm.json")
				if err != nil {
					ui.PrintError(fmt.Errorf("failed to create envm.json: %w", err))
					os.Exit(1)
				}
				encoder := json.NewEncoder(file)
				encoder.SetIndent("", "    ")
				if err := encoder.Encode(config); err != nil {
					file.Close()
					ui.PrintError(fmt.Errorf("failed to write config: %w", err))
					os.Exit(1)
				}
				file.Close()
				ui.PrintSuccess(fmt.Sprintf("Initialized envm.json for project '%s'", selectedProject.Name))

			} else {
				ui.PrintError(fmt.Errorf("failed to open envm.json: %w", err))
				os.Exit(1)
			}
		} else {
			// Read existing
			if err := json.NewDecoder(configFile).Decode(&config); err != nil {
				configFile.Close()
				ui.PrintError(fmt.Errorf("failed to parse envm.json: %w", err))
				os.Exit(1)
			}
			configFile.Close()
		}

		// Fetch all variables for all environments and group by path
		ui.PrintSuccess("Fetching variables from server...")

		envMap, err := fetchProjectEnvs(c, config.ProjectID, creds.Token)
		if err != nil {
			ui.PrintError(fmt.Errorf("failed to fetch environments: %w", err))
			os.Exit(1)
		}

		// Map to store: path -> []variable
		fileMap := make(map[string][]map[string]interface{})

		for name, id := range envMap {
			ui.PrintSuccess(fmt.Sprintf("Fetching variables for '%s'...", name))

			req, err := http.NewRequest("GET", fmt.Sprintf("%s/env/variable/list?environment_id=%s", c.BaseURL, id), nil)
			if err != nil {
				ui.PrintError(err)
				continue
			}
			req.Header.Set("Authorization", "Bearer "+creds.Token)

			resp, err := c.HTTP.Do(req)
			if err != nil {
				ui.PrintError(err)
				continue
			}
			defer resp.Body.Close()

			if resp.StatusCode != http.StatusOK {
				ui.PrintError(fmt.Errorf("api error: %s", resp.Status))
				continue
			}

			var vars []map[string]interface{}
			if err := json.NewDecoder(resp.Body).Decode(&vars); err != nil {
				ui.PrintError(err)
				continue
			}

			// Group variables by their path
			for _, v := range vars {
				path, ok := v["path"].(string)
				if !ok || path == "" {
					path = ".env" // Default fallback
				}
				fileMap[path] = append(fileMap[path], v)
			}
		}

		// Create files at their specified paths
		for path, vars := range fileMap {
			ui.PrintSuccess(fmt.Sprintf("Creating %s...", path))

			// Create parent directories if they don't exist
			dir := filepath.Dir(path)
			if dir != "." && dir != "" {
				if err := os.MkdirAll(dir, 0755); err != nil {
					ui.PrintError(fmt.Errorf("failed to create directory %s: %w", dir, err))
					continue
				}
			}

			// Write to file
			file, err := os.Create(path)
			if err != nil {
				ui.PrintError(fmt.Errorf("failed to create file %s: %w", path, err))
				continue
			}

			for _, v := range vars {
				key, _ := v["key"].(string)
				val, _ := v["value"].(string)
				if key != "" {
					fmt.Fprintf(file, "%s=%s\n", key, val)
				}
			}
			file.Close()
			ui.PrintSuccess(fmt.Sprintf("✓ Updated %s with %d variables", path, len(vars)))
		}
	},
}

func init() {
	rootCmd.AddCommand(pullCmd)
}
