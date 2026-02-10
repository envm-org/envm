package resolver

import (
	"encoding/json"
	"fmt"

	"github.com/envm-org/cli/internal/client"
	"github.com/envm-org/cli/internal/types"
)

type Resolver struct {
	Client *client.Client
}

func New(c *client.Client) *Resolver {
	return &Resolver{Client: c}
}

func (r *Resolver) ResolveOrg(slug string) (string, error) {
	body, err := r.Client.Get("/org/list")
	if err != nil {
		return "", fmt.Errorf("failed to list organizations: %w", err)
	}

	var orgs []types.Organization
	if err := json.Unmarshal(body, &orgs); err != nil {
		return "", fmt.Errorf("failed to parse organizations: %w", err)
	}

	for _, org := range orgs {
		if org.Slug == slug {
			return org.ID, nil
		}
	}

	return "", fmt.Errorf("organization with slug '%s' not found", slug)
}

func (r *Resolver) ResolveProject(orgSlug, projectSlug string) (string, error) {
	orgID, err := r.ResolveOrg(orgSlug)
	if err != nil {
		return "", err
	}

	body, err := r.Client.Get("/project/list?organization_id=" + orgID)
	if err != nil {
		return "", fmt.Errorf("failed to list projects for org '%s': %w", orgSlug, err)
	}

	var projects []types.Project
	if err := json.Unmarshal(body, &projects); err != nil {
		return "", fmt.Errorf("failed to parse projects: %w", err)
	}

	for _, p := range projects {
		if p.Slug == projectSlug {
			return p.ID, nil
		}
	}

	return "", fmt.Errorf("project with slug '%s' not found in org '%s'", projectSlug, orgSlug)
}

func (r *Resolver) ResolveEnv(orgSlug, projectSlug, envSlug string) (string, error) {
	projectID, err := r.ResolveProject(orgSlug, projectSlug)
	if err != nil {
		return "", err
	}

	body, err := r.Client.Get("/env/list?project_id=" + projectID)
	if err != nil {
		return "", fmt.Errorf("failed to list environments for project '%s': %w", projectSlug, err)
	}

	var envs []types.Environment
	if err := json.Unmarshal(body, &envs); err != nil {
		return "", fmt.Errorf("failed to parse environments: %w", err)
	}

	for _, e := range envs {
		if e.Slug == envSlug {
			return e.ID, nil
		}
	}

	return "", fmt.Errorf("environment with slug '%s' not found in project '%s'", envSlug, projectSlug)
}

func (r *Resolver) ResolveUser(email string) (string, error) {
	body, err := r.Client.Get("/users/list")
	if err != nil {
		return "", fmt.Errorf("failed to list users: %w", err)
	}

	var users []types.User
	if err := json.Unmarshal(body, &users); err != nil {
		return "", fmt.Errorf("failed to parse users: %w", err)
	}

	for _, u := range users {
		if u.Email == email {
			return u.ID, nil
		}
	}

	return "", fmt.Errorf("user with email '%s' not found", email)
}
