-- name: CreateUser :one
INSERT INTO users (email, password_hash, full_name)
VALUES ($1, $2, $3)
RETURNING *;

-- name: GetUser :one
SELECT * FROM users
WHERE id = $1 LIMIT 1;

-- name: GetUserByEmail :one
SELECT * FROM users
WHERE email = $1 LIMIT 1;

-- name: ListUsers :many
SELECT * FROM users
ORDER BY created_at DESC;

-- name: UpdateUser :one
UPDATE users
SET email = $2, password_hash = $3, full_name = $4, updated_at = CURRENT_TIMESTAMP
WHERE id = $1
RETURNING *;

-- name: DeleteUser :exec
DELETE FROM users
WHERE id = $1;

-- name: CreateOrganization :one
INSERT INTO organizations (name, slug)
VALUES ($1, $2)
RETURNING *;

-- name: GetOrganization :one
SELECT * FROM organizations
WHERE id = $1 LIMIT 1;

-- name: ListOrganizations :many
SELECT * FROM organizations
ORDER BY name;

-- name: UpdateOrganization :one
UPDATE organizations
SET name = $2, slug = $3, updated_at = CURRENT_TIMESTAMP
WHERE id = $1
RETURNING *;

-- name: DeleteOrganization :exec
DELETE FROM organizations
WHERE id = $1;

-- name: CreateProject :one
INSERT INTO projects (organization_id, name, slug, description)
VALUES ($1, $2, $3, $4)
RETURNING *;

-- name: GetProject :one
SELECT * FROM projects
WHERE id = $1 LIMIT 1;

-- name: ListProjects :many
SELECT * FROM projects
WHERE organization_id = $1
ORDER BY created_at DESC;

-- name: ListProjectsForMember :many
SELECT p.* FROM projects p
JOIN project_members pm ON p.id = pm.project_id
WHERE p.organization_id = $1 AND pm.user_id = $2
ORDER BY p.created_at DESC;

-- name: UpdateProject :one
UPDATE projects
SET name = $2, slug = $3, description = $4, updated_at = CURRENT_TIMESTAMP
WHERE id = $1
RETURNING *;

-- name: DeleteProject :exec
DELETE FROM projects
WHERE id = $1;

-- name: CreateEnvironment :one
INSERT INTO environments (project_id, name, slug)
VALUES ($1, $2, $3)
RETURNING *;

-- name: GetEnvironment :one
SELECT * FROM environments
WHERE id = $1 LIMIT 1;

-- name: ListEnvironments :many
SELECT * FROM environments
WHERE project_id = $1
ORDER BY name;

-- name: UpdateEnvironment :one
UPDATE environments
SET name = $2, slug = $3, updated_at = CURRENT_TIMESTAMP
WHERE id = $1
RETURNING *;

-- name: DeleteEnvironment :exec
DELETE FROM environments
WHERE id = $1;

-- name: CreateVariable :one
INSERT INTO variables (environment_id, key, value, is_secret)
VALUES ($1, $2, $3, $4)
RETURNING *;

-- name: UpdateVariable :one
UPDATE variables
SET value = $3, is_secret = $4, updated_at = CURRENT_TIMESTAMP
WHERE environment_id = $1 AND key = $2
RETURNING *;

-- name: DeleteVariable :exec
DELETE FROM variables
WHERE environment_id = $1 AND key = $2;

-- name: ListVariables :many
SELECT * FROM variables
WHERE environment_id = $1
ORDER BY key;

-- name: CreateAuditLog :one
INSERT INTO audit_logs (user_id, organization_id, action, resource_type, resource_id, details)
VALUES ($1, $2, $3, $4, $5, $6)
RETURNING *;

-- name: AddOrganizationMember :one
INSERT INTO organization_members (organization_id, user_id, role)
VALUES ($1, $2, $3)
RETURNING *;

-- name: RemoveOrganizationMember :exec
DELETE FROM organization_members
WHERE organization_id = $1 AND user_id = $2;

-- name: GetOrganizationMember :one
SELECT * FROM organization_members
WHERE organization_id = $1 AND user_id = $2 LIMIT 1;

-- name: ListOrganizationMembers :many
SELECT om.*, u.email, u.full_name
FROM organization_members om
JOIN users u ON om.user_id = u.id
WHERE om.organization_id = $1
ORDER BY om.created_at;

-- name: SetPasswordResetToken :exec
UPDATE users
SET password_reset_token = $2, password_reset_expires_at = $3
WHERE email = $1;

-- name: GetUserByResetToken :one
SELECT * FROM users
WHERE password_reset_token = $1 AND password_reset_expires_at > NOW()
LIMIT 1;

-- name: UpdatePassword :exec
UPDATE users
SET password_hash = $2, password_reset_token = NULL, password_reset_expires_at = NULL
WHERE id = $1;

-- name: CreateInvitation :one
INSERT INTO organization_invitations (organization_id, email, role, token, expires_at, invited_by)
VALUES ($1, $2, $3, $4, $5, $6)
RETURNING *;

-- name: GetInvitationByToken :one
SELECT * FROM organization_invitations
WHERE token = $1 AND expires_at > NOW()
LIMIT 1;

-- name: DeleteInvitation :exec
DELETE FROM organization_invitations
WHERE id = $1;

-- name: ListInvitations :many
SELECT * FROM organization_invitations
WHERE organization_id = $1;

-- name: AddProjectMember :one
INSERT INTO project_members (project_id, user_id, role)
VALUES ($1, $2, $3)
RETURNING *;

-- name: RemoveProjectMember :exec
DELETE FROM project_members
WHERE project_id = $1 AND user_id = $2;

-- name: GetProjectMember :one
SELECT * FROM project_members
WHERE project_id = $1 AND user_id = $2;

-- name: ListProjectMembers :many
SELECT pm.*, u.email, u.full_name
FROM project_members pm
JOIN users u ON pm.user_id = u.id
WHERE pm.project_id = $1;
