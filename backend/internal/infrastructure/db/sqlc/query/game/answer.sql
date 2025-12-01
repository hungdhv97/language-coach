-- name: CreateGameAnswer :one
INSERT INTO vocab_game_question_answers (
    question_id, session_id, user_id,
    selected_option_id, is_correct, response_time_ms, answered_at
) VALUES ($1, $2, $3, $4, $5, $6, $7)
RETURNING id, answered_at;

-- name: FindGameAnswerByQuestionID :one
SELECT id, question_id, session_id, user_id,
       selected_option_id, is_correct, response_time_ms, answered_at
FROM vocab_game_question_answers
WHERE question_id = $1 AND session_id = $2 AND user_id = $3
LIMIT 1;

-- name: FindGameAnswersBySessionID :many
SELECT id, question_id, session_id, user_id,
       selected_option_id, is_correct, response_time_ms, answered_at
FROM vocab_game_question_answers
WHERE session_id = $1 AND user_id = $2
ORDER BY answered_at;

