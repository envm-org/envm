package auth

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/envm-org/cli/internal/types"
)

const (
	defaultAPIURL = "http://localhost:8080"
)

type Client struct {
	BaseURL string
	HTTP    *http.Client
}

func NewClient(baseURL string) *Client {
	if baseURL == "" {
		baseURL = defaultAPIURL
	}
	return &Client{
		BaseURL: baseURL,
		HTTP: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

func (c *Client) Login(email, password string) (*LoginResponse, error) {
	payload := map[string]string{
		"email":    email,
		"password": password,
	}
	data, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}

	resp, err := c.HTTP.Post(fmt.Sprintf("%s/auth/login", c.BaseURL), "application/json", bytes.NewBuffer(data))
	if err != nil {
		return nil, fmt.Errorf("connection failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, parseError(resp)
	}

	var loginResp LoginResponse
	if err := json.NewDecoder(resp.Body).Decode(&loginResp); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	return &loginResp, nil
}

func (c *Client) Register(fullName, email, password string) (*User, error) {
	payload := map[string]string{
		"full_name": fullName,
		"email":     email,
		"password":  password,
	}
	data, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}

	resp, err := c.HTTP.Post(fmt.Sprintf("%s/auth/register", c.BaseURL), "application/json", bytes.NewBuffer(data))
	if err != nil {
		return nil, fmt.Errorf("connection failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		return nil, parseError(resp)
	}

	var user User
	if err := json.NewDecoder(resp.Body).Decode(&user); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	return &user, nil
}

func parseError(resp *http.Response) error {
	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("api error: %s (failed to read body)", resp.Status)
	}

	var errResp struct {
		Error string `json:"error"`
	}
	// Try to parse JSON error first
	if err := json.Unmarshal(bodyBytes, &errResp); err == nil && errResp.Error != "" {
		return fmt.Errorf("api error: %s", errResp.Error)
	}

	// Fallback to body content (it might be a plain text error from http.Error)
	if len(bodyBytes) > 0 {
		return fmt.Errorf("api error: %s", string(bodyBytes))
	}

	return fmt.Errorf("api error: %s", resp.Status)
}

func (c *Client) CreateOrganization(name, slug, token string) (*types.Organization, error) {
	payload := map[string]string{
		"name": name,
		"slug": slug,
	}
	data, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", fmt.Sprintf("%s/orgs", c.BaseURL), bytes.NewBuffer(data))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)

	resp, err := c.HTTP.Do(req)
	if err != nil {
		return nil, fmt.Errorf("connection failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, parseError(resp)
	}

	var org types.Organization
	if err := json.NewDecoder(resp.Body).Decode(&org); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	return &org, nil
}

func (c *Client) CreateProject(orgID, name, slug, description, token string) (*types.Project, error) {
	payload := map[string]string{
		"organization_id": orgID,
		"name":            name,
		"slug":            slug,
		"description":     description,
	}
	data, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", fmt.Sprintf("%s/projects", c.BaseURL), bytes.NewBuffer(data))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)

	resp, err := c.HTTP.Do(req)
	if err != nil {
		return nil, fmt.Errorf("connection failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, parseError(resp)
	}

	var project types.Project
	if err := json.NewDecoder(resp.Body).Decode(&project); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	return &project, nil
}
