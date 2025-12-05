-- name: CreateUser :one
INSERT INTO users (email, username, password_hash, is_active)
VALUES ($1, $2, $3, $4)
RETURNING id, email, username, password_hash, created_at, updated_at, is_active;

-- name: FindUserByID :one
SELECT id, email, username, password_hash, created_at, updated_at, is_active
FROM users
WHERE id = $1;

-- name: FindUserByEmail :one
SELECT id, email, username, password_hash, created_at, updated_at, is_active
FROM users
WHERE email = $1;

-- name: FindUserByUsername :one
SELECT id, email, username, password_hash, created_at, updated_at, is_active
FROM users
WHERE username = $1;

-- name: UpdateUserPassword :exec
UPDATE users
SET password_hash = $2, updated_at = CURRENT_TIMESTAMP
WHERE id = $1;

-- name: UpdateUserActiveStatus :exec
UPDATE users
SET is_active = $2, updated_at = CURRENT_TIMESTAMP
WHERE id = $1;

-- name: CheckEmailExists :one
SELECT EXISTS(SELECT 1 FROM users WHERE email = $1) as exists;

-- name: CheckUsernameExists :one
SELECT EXISTS(SELECT 1 FROM users WHERE username = $1) as exists;
