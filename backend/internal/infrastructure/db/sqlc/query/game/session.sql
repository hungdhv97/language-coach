-- name: CreateGameSession :one
INSERT INTO vocab_game_sessions (
    user_id, mode, source_language_id, target_language_id,
    topic_id, level_id, total_questions, correct_questions,
    started_at
) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
RETURNING id, started_at;

-- name: FindGameSessionByID :one
SELECT id, user_id, mode, source_language_id, target_language_id,
       topic_id, level_id, total_questions, correct_questions,
       started_at, ended_at
FROM vocab_game_sessions
WHERE id = $1;

-- name: UpdateGameSession :exec
UPDATE vocab_game_sessions
SET total_questions = $2,
    correct_questions = $3,
    ended_at = $4
WHERE id = $1;

-- name: EndGameSession :exec
UPDATE vocab_game_sessions 
SET ended_at = $2 
WHERE id = $1;

