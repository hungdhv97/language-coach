-- name: CreateGameQuestion :one
INSERT INTO vocab_game_questions (
    session_id, question_order, question_type,
    source_word_id, source_sense_id, correct_target_word_id,
    source_language_id, target_language_id, created_at
) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
RETURNING id, created_at;

-- name: CreateGameQuestionOption :one
INSERT INTO vocab_game_question_options (
    question_id, option_label, target_word_id, is_correct
) VALUES ($1, $2, $3, $4)
RETURNING id;

-- name: FindGameQuestionsBySessionID :many
SELECT id, session_id, question_order, question_type,
       source_word_id, source_sense_id, correct_target_word_id,
       source_language_id, target_language_id, created_at
FROM vocab_game_questions
WHERE session_id = $1
ORDER BY question_order;

-- name: FindGameQuestionByID :one
SELECT id, session_id, question_order, question_type,
       source_word_id, source_sense_id, correct_target_word_id,
       source_language_id, target_language_id, created_at
FROM vocab_game_questions
WHERE id = $1;

-- name: FindGameQuestionOptionsByQuestionID :many
SELECT id, question_id, option_label, target_word_id, is_correct
FROM vocab_game_question_options
WHERE question_id = $1
ORDER BY option_label;

-- name: FindGameQuestionOptionsByQuestionIDs :many
SELECT id, question_id, option_label, target_word_id, is_correct
FROM vocab_game_question_options
WHERE question_id = ANY($1::bigint[])
ORDER BY question_id, option_label;

