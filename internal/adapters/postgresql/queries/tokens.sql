-- name: CreateRefreshToken :one
INSERT INTO refresh_tokens (
    user_id, token, expires_at
) VALUES (
    $1, $2, $3
) RETURNING *;

-- name: GetRefreshToken :one
SELECT * FROM refresh_tokens
WHERE token = $1 LIMIT 1;

-- name: RevokeRefreshToken :exec
UPDATE refresh_tokens
SET revoked = true, updated_at = CURRENT_TIMESTAMP
WHERE token = $1;

-- name: DeleteRefreshToken :exec
DELETE FROM refresh_tokens
WHERE token = $1;

-- name: ListActiveRefreshTokensForUser :many
SELECT * FROM refresh_tokens
WHERE user_id = $1 AND revoked = false AND expires_at > CURRENT_TIMESTAMP;
