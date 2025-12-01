-- name: FindAllTopics :many
SELECT id, code, name 
FROM topics 
ORDER BY code;

-- name: FindTopicByID :one
SELECT id, code, name 
FROM topics 
WHERE id = $1;

-- name: FindTopicByCode :one
SELECT id, code, name 
FROM topics 
WHERE code = $1;

