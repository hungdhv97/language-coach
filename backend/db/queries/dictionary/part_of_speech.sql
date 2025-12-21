-- name: FindAllPartsOfSpeech :many
SELECT id, code, name 
FROM parts_of_speech 
ORDER BY code;

-- name: FindPartOfSpeechByID :one
SELECT id, code, name 
FROM parts_of_speech 
WHERE id = $1;

-- name: FindPartOfSpeechByCode :one
SELECT id, code, name 
FROM parts_of_speech 
WHERE code = $1;

-- name: FindPartsOfSpeechByIDs :many
SELECT id, code, name 
FROM parts_of_speech 
WHERE id = ANY($1::smallint[])
ORDER BY id;
