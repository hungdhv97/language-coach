-- name: FindAllLevels :many
SELECT id, code, name, description, language_id, difficulty_order 
FROM levels 
ORDER BY language_id, difficulty_order NULLS LAST, code;

-- name: FindLevelByID :one
SELECT id, code, name, description, language_id, difficulty_order 
FROM levels 
WHERE id = $1;

-- name: FindLevelByCode :one
SELECT id, code, name, description, language_id, difficulty_order 
FROM levels 
WHERE code = $1;

-- name: FindLevelsByLanguageID :many
SELECT id, code, name, description, language_id, difficulty_order 
FROM levels 
WHERE language_id = $1 
ORDER BY difficulty_order, code;

