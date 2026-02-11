package types

type Organization struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	Slug      string `json:"slug"`
	CreatedAt string `json:"created_at"`
}

type Project struct {
	ID             string `json:"id"`
	OrganizationID string `json:"organization_id"`
	Name           string `json:"name"`
	Slug           string `json:"slug"`
	Description    string `json:"description"`
	CreatedAt      string `json:"created_at"`
}

type Environment struct {
	ID        string `json:"id"`
	ProjectID string `json:"project_id"`
	Name      string `json:"name"`
	Slug      string `json:"slug"`
	CreatedAt string `json:"created_at"`
}

type User struct {
	ID        string `json:"id"`
	Email     string `json:"email"`
	FullName  string `json:"full_name"`
	CreatedAt string `json:"created_at"`
}

type ProjectConfig struct {
	OwnerID     string            `json:"ownerId"`
	ProjectID   string            `json:"projectId"`
	Envs        []EnvEntry        `json:"envs"`
	Credentials []CredentialEntry `json:"credentials"`
}

type EnvEntry struct {
	Name string `json:"name"`
	Path string `json:"path"`
}

type CredentialEntry struct {
	Name string `json:"name"`
	Path string `json:"path"`
}
