-- name: FindAllLanguages :many
SELECT id, code, name 
FROM languages 
ORDER BY code;

-- name: FindLanguageByID :one
SELECT id, code, name 
FROM languages 
WHERE id = $1;

-- name: FindLanguageByCode :one
SELECT id, code, name 
FROM languages 
WHERE code = $1;

