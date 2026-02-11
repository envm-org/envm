package cmd

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/envm-org/cli/internal/auth"
	"github.com/envm-org/cli/internal/types"
)

func fetchProjectEnvs(c *auth.Client, projectID, token string) (map[string]string, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/env/list?project_id=%s", c.BaseURL, projectID), nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "Bearer "+token)

	resp, err := c.HTTP.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("api error: %s", resp.Status)
	}

	var envs []map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&envs); err != nil {
		return nil, err
	}

	result := make(map[string]string)
	for _, e := range envs {
		name, _ := e["name"].(string)
		id, _ := e["id"].(string)
		if name != "" && id != "" {
			result[name] = id
		}
	}
	return result, nil
}

func createEnvironment(c *auth.Client, projectID, name, slug, token string) (string, error) {
	payload := map[string]string{
		"project_id": projectID,
		"name":       name,
		"slug":       slug,
	}
	data, _ := json.Marshal(payload)

	req, err := http.NewRequest("POST", fmt.Sprintf("%s/env", c.BaseURL), bytes.NewBuffer(data))
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)

	resp, err := c.HTTP.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		return "", fmt.Errorf("api error: %s", resp.Status)
	}

	var env map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&env); err != nil {
		return "", err
	}
	id, _ := env["id"].(string)
	return id, nil
}

func fetchVariables(c *auth.Client, envID, token string) (map[string]bool, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/env/variable/list?environment_id=%s", c.BaseURL, envID), nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "Bearer "+token)

	resp, err := c.HTTP.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("api error: %s", resp.Status)
	}

	var vars []map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&vars); err != nil {
		return nil, err
	}

	result := make(map[string]bool)
	for _, v := range vars {
		key, _ := v["key"].(string)
		if key != "" {
			result[key] = true
		}
	}
	return result, nil
}

func createVariable(c *auth.Client, envID, key, value, path, token string) error {
	payload := map[string]interface{}{
		"environment_id": envID,
		"key":            key,
		"value":          value,
		"is_secret":      false, // Default to false for now
		"path":           path,
	}
	data, _ := json.Marshal(payload)

	req, err := http.NewRequest("POST", fmt.Sprintf("%s/env/variable", c.BaseURL), bytes.NewBuffer(data))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)

	resp, err := c.HTTP.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		return fmt.Errorf("api error: %s", resp.Status)
	}
	return nil
}

func updateVariable(c *auth.Client, envID, key, value, path, token string) error {
	payload := map[string]interface{}{
		"environment_id": envID,
		"key":            key,
		"value":          value,
		"is_secret":      false,
		"path":           path,
	}
	data, _ := json.Marshal(payload)

	req, err := http.NewRequest("PUT", fmt.Sprintf("%s/env/variable", c.BaseURL), bytes.NewBuffer(data))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)

	resp, err := c.HTTP.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("api error: %s", resp.Status)
	}
	return nil
}

func scanEnvFiles() ([]types.EnvEntry, error) {
	var envs []types.EnvEntry
	err := filepath.Walk(".", func(path string, info os.FileInfo, err error) error {
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
	return envs, err
}
