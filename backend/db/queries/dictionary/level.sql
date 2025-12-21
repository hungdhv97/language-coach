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
SELECT l.id, l.code, l.name, l.description, l.language_id, l.difficulty_order
FROM levels AS l
WHERE 
    (
        l.language_id = $1
    )
    OR (
        NOT EXISTS (
            SELECT 1 
            FROM levels AS l2 
            WHERE l2.language_id = $1
        )
        AND l.language_id IS NULL
    )
ORDER BY l.difficulty_order, l.code;

