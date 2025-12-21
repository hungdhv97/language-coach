-- name: FindWordByID :one
SELECT id, language_id, lemma, lemma_normalized, search_key,
       romanization, script_code, frequency_rank,
       note, created_at, updated_at
FROM words
WHERE id = $1;

-- name: FindWordsByIDs :many
SELECT id, language_id, lemma, lemma_normalized, search_key,
       romanization, script_code, frequency_rank,
       note, created_at, updated_at
FROM words
WHERE id = ANY($1::bigint[])
ORDER BY id;

-- name: FindWordsByTopicAndLanguages :many
SELECT DISTINCT w.id, w.language_id, w.lemma, w.lemma_normalized, w.search_key,
       w.romanization, w.script_code, w.frequency_rank,
       w.note, w.created_at, w.updated_at
FROM words w
INNER JOIN word_topics wt ON w.id = wt.word_id
WHERE wt.topic_id = sqlc.arg('topic_id')
  AND w.language_id = sqlc.arg('source_language_id')
  AND EXISTS (
      SELECT 1
      FROM senses s
      INNER JOIN sense_translations st ON s.id = st.source_sense_id
      INNER JOIN words tw ON st.target_word_id = tw.id
      WHERE s.word_id = w.id
        AND tw.language_id = sqlc.arg('target_language_id')
  )
ORDER BY w.frequency_rank NULLS LAST, w.id
LIMIT sqlc.arg('limit');

-- name: FindWordsByLevelAndLanguages :many
SELECT DISTINCT w.id, w.language_id, w.lemma, w.lemma_normalized, w.search_key,
       w.romanization, w.script_code, w.frequency_rank,
       w.note, w.created_at, w.updated_at
FROM words w
INNER JOIN senses s ON w.id = s.word_id
WHERE s.level_id = sqlc.arg('level_id')
  AND w.language_id = sqlc.arg('source_language_id')
  AND EXISTS (
      SELECT 1
      FROM sense_translations st
      INNER JOIN words tw ON st.target_word_id = tw.id
      WHERE st.source_sense_id = s.id
        AND tw.language_id = sqlc.arg('target_language_id')
  )
ORDER BY w.frequency_rank NULLS LAST, w.id
LIMIT sqlc.arg('limit');

-- name: FindWordsByLevelAndTopicsAndLanguages :many
SELECT DISTINCT w.id, w.language_id, w.lemma, w.lemma_normalized, w.search_key,
       w.romanization, w.script_code, w.frequency_rank,
       w.note, w.created_at, w.updated_at
FROM words w
INNER JOIN senses s ON w.id = s.word_id
WHERE s.level_id = sqlc.arg('level_id')
  AND w.language_id = sqlc.arg('source_language_id')
  AND (
    -- If topic_ids array is empty/null, include all words
    -- Otherwise filter by topic_ids using ANY
    sqlc.arg('topic_ids')::bigint[] IS NULL
    OR array_length(sqlc.arg('topic_ids')::bigint[], 1) IS NULL
    OR EXISTS (
      SELECT 1
      FROM word_topics wt
      WHERE wt.word_id = w.id
        AND wt.topic_id = ANY(sqlc.arg('topic_ids')::bigint[])
    )
  )
  AND EXISTS (
      SELECT 1
      FROM sense_translations st
      INNER JOIN words tw ON st.target_word_id = tw.id
      WHERE st.source_sense_id = s.id
        AND tw.language_id = sqlc.arg('target_language_id')
  )
ORDER BY w.frequency_rank NULLS LAST, w.id
LIMIT sqlc.arg('limit');

-- name: FindTranslationsForWord :many
SELECT DISTINCT tw.id, tw.language_id, tw.lemma, tw.lemma_normalized, tw.search_key,
       tw.romanization, tw.script_code, tw.frequency_rank,
       tw.note, tw.created_at, tw.updated_at
FROM words sw
INNER JOIN senses s ON sw.id = s.word_id
INNER JOIN sense_translations st ON s.id = st.source_sense_id
INNER JOIN words tw ON st.target_word_id = tw.id
WHERE sw.id = sqlc.arg('source_word_id')
  AND tw.language_id = sqlc.arg('target_language_id')
ORDER BY st.priority, tw.frequency_rank NULLS LAST, tw.id
LIMIT sqlc.arg('limit');

-- name: SearchWords :many
SELECT w.id, w.language_id, w.lemma, w.lemma_normalized, w.search_key,
       w.romanization, w.script_code, w.frequency_rank,
       w.note, w.created_at, w.updated_at
FROM words w
WHERE w.language_id = sqlc.arg('language_id')
  AND (
    w.lemma ILIKE sqlc.arg('search_pattern')
    OR w.lemma_normalized ILIKE sqlc.arg('search_pattern')
    OR w.search_key ILIKE sqlc.arg('search_pattern')
  )
ORDER BY 
  CASE 
    WHEN w.lemma = sqlc.arg('exact_match') THEN 1
    WHEN w.lemma ILIKE sqlc.arg('search_pattern') THEN 2
    WHEN w.lemma_normalized ILIKE sqlc.arg('search_pattern') THEN 3
    WHEN w.search_key ILIKE sqlc.arg('search_pattern') THEN 4
    ELSE 5
  END,
  w.frequency_rank NULLS LAST,
  w.id
LIMIT sqlc.arg('limit') OFFSET sqlc.arg('offset');

-- name: CountSearchWords :one
SELECT COUNT(DISTINCT w.id)
FROM words w
WHERE w.language_id = sqlc.arg('language_id')
  AND (
    w.lemma ILIKE sqlc.arg('search_pattern')
    OR w.lemma_normalized ILIKE sqlc.arg('search_pattern')
    OR w.search_key ILIKE sqlc.arg('search_pattern')
  );

