-- name: FindSensesByWordID :many
SELECT id, word_id, sense_order, definition, definition_language_id,
       usage_label, level_id, note
FROM senses
WHERE word_id = $1
ORDER BY sense_order;

-- name: FindSensesByWordIDs :many
SELECT id, word_id, sense_order, definition, definition_language_id,
       usage_label, level_id, note
FROM senses
WHERE word_id = ANY($1::bigint[])
ORDER BY word_id, sense_order;

