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
ORDER BY name;

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
