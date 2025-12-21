-- name: CreateUserProfile :one
INSERT INTO user_profiles (user_id, display_name, avatar_url, birth_day, bio)
VALUES ($1, $2, $3, $4, $5)
RETURNING user_id, display_name, avatar_url, birth_day, bio, created_at, updated_at;

-- name: GetUserProfile :one
SELECT user_id, display_name, avatar_url, birth_day, bio, created_at, updated_at
FROM user_profiles
WHERE user_id = $1;

-- name: UpdateUserProfile :one
UPDATE user_profiles
SET display_name = COALESCE($2, display_name),
    avatar_url = COALESCE($3, avatar_url),
    birth_day = COALESCE($4, birth_day),
    bio = COALESCE($5, bio),
    updated_at = CURRENT_TIMESTAMP
WHERE user_id = $1
RETURNING user_id, display_name, avatar_url, birth_day, bio, created_at, updated_at;
