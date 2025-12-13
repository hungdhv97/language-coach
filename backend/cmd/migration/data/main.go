package main

import (
	"bufio"
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	appconfig "github.com/english-coach/backend/internal/config"
)

// Seed data models for initial JSON (0001_init_data.json)
type Language struct {
	Code string `json:"code"`
	Name string `json:"name"`
}

type PartOfSpeech struct {
	Code string `json:"code"`
	Name string `json:"name"`
}

type Topic struct {
	Code string `json:"code"`
	Name string `json:"name"`
}

type Level struct {
	Code            string  `json:"code"`
	Name            string  `json:"name"`
	Description     string  `json:"description"`
	Language        *string `json:"language"`         // language code: "en", "vi", "zh" or null
	DifficultyOrder int     `json:"difficulty_order"` // difficulty order: 1 < 2 < 3 ...
}

type SeedData struct {
	Languages     []Language     `json:"languages"`
	PartsOfSpeech []PartOfSpeech `json:"parts_of_speech"`
	Topics        []Topic        `json:"topics"`
	Levels        []Level        `json:"levels"`
}

// Word JSONL models (0002/0003/0004_word_*.jsonl)
// New format: all fields at top level, no entries array
type WordJSON struct {
	Language        string              `json:"language"`
	Lemma           string              `json:"lemma"`
	LemmaNormalized *string             `json:"lemma_normalized,omitempty"`
	SearchKey       *string             `json:"search_key,omitempty"`
	Romanization    *string             `json:"romanization,omitempty"`
	ScriptCode      *string             `json:"script_code,omitempty"`
	FrequencyRank   *int                `json:"frequency_rank,omitempty"`
	Note            *string             `json:"note,omitempty"`
	Topics          []string            `json:"topics,omitempty"`
	Pronunciations  []PronunciationJSON `json:"pronunciations,omitempty"`
	Relations       []WordRelationJSON  `json:"relations,omitempty"`
	Senses          []SenseJSON         `json:"senses,omitempty"`
	Characters      []CharacterJSON     `json:"characters,omitempty"` // used mainly for Chinese words
}

type PronunciationJSON struct {
	Dialect  string  `json:"dialect"`
	IPA      *string `json:"ipa,omitempty"`
	Phonetic *string `json:"phonetic,omitempty"`
	AudioURL *string `json:"audio_url,omitempty"`
}

type WordRelationJSON struct {
	RelationType string          `json:"relation_type"`
	Note         *string         `json:"note,omitempty"`
	TargetWord   RelatedWordJSON `json:"target_word"`
}

type RelatedWordJSON struct {
	Language        string   `json:"language"`
	Lemma           string   `json:"lemma"`
	LemmaNormalized *string  `json:"lemma_normalized,omitempty"`
	SearchKey       *string  `json:"search_key,omitempty"`
	PartOfSpeech    string   `json:"part_of_speech"`
	Romanization    *string  `json:"romanization,omitempty"`
	ScriptCode      *string  `json:"script_code,omitempty"`
	FrequencyRank   *int     `json:"frequency_rank,omitempty"`
	Note            *string  `json:"note,omitempty"`
	Topics          []string `json:"topics,omitempty"`
}

type SenseJSON struct {
	Order              int                    `json:"order"`
	PartOfSpeech       string                 `json:"part_of_speech"`
	DefinitionLanguage string                 `json:"definition_language"`
	Definition         string                 `json:"definition"`
	UsageLabel         *string                `json:"usage_label,omitempty"`
	Level              *string                `json:"level,omitempty"`
	Note               *string                `json:"note,omitempty"`
	Translations       []SenseTranslationJSON `json:"translations,omitempty"`
	Examples           []ExampleJSON          `json:"examples,omitempty"`
}

type SenseTranslationJSON struct {
	Priority   int             `json:"priority"`
	Note       *string         `json:"note,omitempty"`
	TargetWord RelatedWordJSON `json:"target_word"`
}

type ExampleJSON struct {
	Language     string                   `json:"language"`
	Content      string                   `json:"content"`
	AudioURL     *string                  `json:"audio_url,omitempty"`
	Translations []ExampleTranslationJSON `json:"translations,omitempty"`
}

type ExampleTranslationJSON struct {
	Language string `json:"language"`
	Content  string `json:"content"`
}

type CharacterJSON struct {
	Literal     string                 `json:"literal"`
	Simplified  *string                `json:"simplified,omitempty"`
	Traditional *string                `json:"traditional,omitempty"`
	ScriptCode  string                 `json:"script_code"`
	Strokes     *int                   `json:"strokes,omitempty"`
	Radical     *string                `json:"radical,omitempty"`
	Level       *string                `json:"level,omitempty"` // level code, will be converted to level_id
	CharOrder   int                    `json:"char_order"`
	Readings    []CharacterReadingJSON `json:"readings,omitempty"`
}

type CharacterReadingJSON struct {
	Language    string  `json:"language"`
	Reading     string  `json:"reading"`
	ReadingType *string `json:"reading_type,omitempty"`
	Note        *string `json:"note,omitempty"`
}

const (
	initDataPath   = "internal/infrastructure/db/migrations/data/0001_init_data.json"
	wordEnDataPath = "internal/infrastructure/db/migrations/data/0002_word_en.jsonl"
	wordViDataPath = "internal/infrastructure/db/migrations/data/0003_word_vi.jsonl"
	wordZhDataPath = "internal/infrastructure/db/migrations/data/0004_word_zh.jsonl"
)

func main() {
	initFlag := flag.Bool("init", false, "Upsert initial dictionary metadata (languages, parts of speech, topics, levels)")
	wordEnFlag := flag.Bool("word-en", false, "Upsert English words from JSONL")
	wordViFlag := flag.Bool("word-vi", false, "Upsert Vietnamese words from JSONL")
	wordZhFlag := flag.Bool("word-zh", false, "Upsert Chinese words from JSONL")
	dsn := flag.String("dsn", "", "PostgreSQL DSN (or use env DATABASE_URL / app config)")
	flag.Parse()

	// If no action flags provided, run full seed: init + all word files
	if !*initFlag && !*wordEnFlag && !*wordViFlag && !*wordZhFlag {
		*initFlag = true
		*wordEnFlag = true
		*wordViFlag = true
		*wordZhFlag = true
	}

	ctx := context.Background()

	pool, err := connectDB(ctx, *dsn)
	if err != nil {
		log.Fatalf("database connection error: %v", err)
	}
	defer pool.Close()

	if *initFlag {
		if err := runInit(ctx, pool, initDataPath); err != nil {
			log.Fatalf("init seed error: %v", err)
		}
		fmt.Println("Initial metadata seed completed successfully.")
	}

	if *wordEnFlag {
		if err := upsertWordsFromJSONL(ctx, pool, "en", wordEnDataPath); err != nil {
			log.Fatalf("word-en upsert error: %v", err)
		}
		fmt.Println("English words upsert completed successfully.")
	}

	if *wordViFlag {
		if err := upsertWordsFromJSONL(ctx, pool, "vi", wordViDataPath); err != nil {
			log.Fatalf("word-vi upsert error: %v", err)
		}
		fmt.Println("Vietnamese words upsert completed successfully.")
	}

	if *wordZhFlag {
		if err := upsertWordsFromJSONL(ctx, pool, "zh", wordZhDataPath); err != nil {
			log.Fatalf("word-zh upsert error: %v", err)
		}
		fmt.Println("Chinese words upsert completed successfully.")
	}
}

func connectDB(ctx context.Context, cliDSN string) (*pgxpool.Pool, error) {
	dsn := cliDSN
	if dsn == "" {
		// First try DATABASE_URL
		dsn = os.Getenv("DATABASE_URL")
	}

	if dsn == "" {
		// Fall back to app config (backend/internal/config/config.go)
		cfg, err := appconfig.Load()
		if err != nil {
			return nil, fmt.Errorf("load app config: %w", err)
		}

		dbCfg := cfg.Database
		dsn = fmt.Sprintf(
			"host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
			dbCfg.Host,
			dbCfg.Port,
			dbCfg.User,
			dbCfg.Password,
			dbCfg.Database,
			dbCfg.SSLMode,
		)
	}

	pool, err := pgxpool.New(ctx, dsn)
	if err != nil {
		return nil, fmt.Errorf("create pgx pool: %w", err)
	}

	if err := pool.Ping(ctx); err != nil {
		pool.Close()
		return nil, fmt.Errorf("ping database: %w", err)
	}

	return pool, nil
}

func runInit(ctx context.Context, pool *pgxpool.Pool, filePath string) error {
	f, err := os.Open(filePath)
	if err != nil {
		return fmt.Errorf("open json file: %w", err)
	}
	defer f.Close()

	var data SeedData
	if err := json.NewDecoder(f).Decode(&data); err != nil {
		return fmt.Errorf("decode json: %w", err)
	}

	tx, err := pool.Begin(ctx)
	if err != nil {
		return fmt.Errorf("begin tx: %w", err)
	}
	defer tx.Rollback(ctx)

	if err := upsertLanguages(ctx, tx, data.Languages); err != nil {
		return err
	}
	if err := upsertPartsOfSpeech(ctx, tx, data.PartsOfSpeech); err != nil {
		return err
	}
	if err := upsertTopics(ctx, tx, data.Topics); err != nil {
		return err
	}
	if err := upsertLevels(ctx, tx, data.Levels); err != nil {
		return err
	}

	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("commit tx: %w", err)
	}

	return nil
}

func upsertLanguages(ctx context.Context, tx pgx.Tx, langs []Language) error {
	const q = `
INSERT INTO languages (code, name)
VALUES ($1, $2)
ON CONFLICT (code) DO UPDATE
SET name = EXCLUDED.name;
`
	for _, l := range langs {
		if _, err := tx.Exec(ctx, q, l.Code, l.Name); err != nil {
			return fmt.Errorf("upsert language %s: %w", l.Code, err)
		}
	}
	return nil
}

func upsertPartsOfSpeech(ctx context.Context, tx pgx.Tx, pos []PartOfSpeech) error {
	const q = `
INSERT INTO parts_of_speech (code, name)
VALUES ($1, $2)
ON CONFLICT (code) DO UPDATE
SET name = EXCLUDED.name;
`
	for _, p := range pos {
		if _, err := tx.Exec(ctx, q, p.Code, p.Name); err != nil {
			return fmt.Errorf("upsert part_of_speech %s: %w", p.Code, err)
		}
	}
	return nil
}

func upsertTopics(ctx context.Context, tx pgx.Tx, topics []Topic) error {
	const q = `
INSERT INTO topics (code, name)
VALUES ($1, $2)
ON CONFLICT (code) DO UPDATE
SET name = EXCLUDED.name;
`
	for _, t := range topics {
		if _, err := tx.Exec(ctx, q, t.Code, t.Name); err != nil {
			return fmt.Errorf("upsert topic %s: %w", t.Code, err)
		}
	}
	return nil
}

func upsertLevels(ctx context.Context, tx pgx.Tx, levels []Level) error {
	const q = `
INSERT INTO levels (code, name, description, language_id, difficulty_order)
VALUES (
    $1,
    $2,
    $3,
    CASE
        WHEN $4::varchar IS NULL THEN NULL::smallint
        ELSE (SELECT id FROM languages WHERE code = $4::varchar)
    END,
    $5
)
ON CONFLICT (language_id, code) DO UPDATE
SET name = EXCLUDED.name,
    description = EXCLUDED.description,
    difficulty_order = EXCLUDED.difficulty_order;
`
	for _, lv := range levels {
		var langCode interface{}
		if lv.Language != nil && *lv.Language != "" {
			langCode = *lv.Language // e.g., "zh"
		} else {
			langCode = nil // NULL -> shared level (CEFR)
		}

		if _, err := tx.Exec(ctx, q, lv.Code, lv.Name, lv.Description, langCode, lv.DifficultyOrder); err != nil {
			return fmt.Errorf("upsert level %s: %w", lv.Code, err)
		}
	}
	return nil
}

// upsertWordsFromJSONL reads a JSONL file and upserts words for a given language.
func upsertWordsFromJSONL(ctx context.Context, pool *pgxpool.Pool, languageCode, filePath string) error {
	f, err := os.Open(filePath)
	if err != nil {
		return fmt.Errorf("open jsonl file: %w", err)
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)

	// Increase the scanner buffer in case of long lines
	const maxLineSize = 1024 * 1024 // 1MB
	buf := make([]byte, 0, 64*1024)
	scanner.Buffer(buf, maxLineSize)

	// Caches to reduce round-trips
	langID, err := getLanguageID(ctx, pool, languageCode)
	if err != nil {
		return err
	}

	languageCache := map[string]int16{
		languageCode: langID,
	}
	posCache := make(map[string]*int16)
	topicCache := make(map[string]int64)
	levelCache := make(map[string]*int64)
	wordCache := make(map[string]int64)      // key: lang|lemma|pos
	characterCache := make(map[string]int64) // key: literal|script

	tx, err := pool.Begin(ctx)
	if err != nil {
		return fmt.Errorf("begin tx: %w", err)
	}
	defer tx.Rollback(ctx)

	lineNumber := 0
	wordCount := 0
	for scanner.Scan() {
		lineNumber++
		line := scanner.Bytes()
		if len(line) == 0 {
			continue
		}

		var w WordJSON
		if err := json.Unmarshal(line, &w); err != nil {
			return fmt.Errorf("decode word json (line %d): %w", lineNumber, err)
		}

		// Upsert the word (single word per line in new format)
		wordID, err := upsertSingleWord(ctx, tx, langID, w.Lemma, w)
		if err != nil {
			return fmt.Errorf("upsert word %s at line %d: %w", w.Lemma, lineNumber, err)
		}
		wordCount++

		if err := upsertWordDetails(
			ctx,
			pool,
			tx,
			wordID,
			w,
			languageCache,
			posCache,
			topicCache,
			levelCache,
			wordCache,
			characterCache,
		); err != nil {
			return fmt.Errorf("upsert word details %s at line %d: %w", w.Lemma, lineNumber, err)
		}
	}

	if err := scanner.Err(); err != nil {
		return fmt.Errorf("scan jsonl: %w", err)
	}

	if wordCount == 0 {
		return fmt.Errorf("no words found in file %s", filePath)
	}

	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("commit tx: %w", err)
	}

	fmt.Printf("  Processed %d words from %s\n", wordCount, filePath)
	return nil
}

func getLanguageID(ctx context.Context, pool *pgxpool.Pool, code string) (int16, error) {
	const q = `SELECT id FROM languages WHERE code = $1`
	var id int16
	if err := pool.QueryRow(ctx, q, code).Scan(&id); err != nil {
		return 0, fmt.Errorf("get language id for %s: %w", code, err)
	}
	return id, nil
}

func getPartOfSpeechID(ctx context.Context, tx pgx.Tx, cache map[string]*int16, code string) (*int16, error) {
	if code == "" {
		return nil, nil
	}

	if v, ok := cache[code]; ok {
		return v, nil
	}

	const q = `SELECT id FROM parts_of_speech WHERE code = $1`
	var id int16
	err := tx.QueryRow(ctx, q, code).Scan(&id)
	if err != nil {
		if err == pgx.ErrNoRows {
			cache[code] = nil
			return nil, nil
		}
		return nil, err
	}
	cache[code] = &id
	return &id, nil
}

// Helper: get language id using a small in-memory cache.
func getLanguageIDWithCache(ctx context.Context, pool *pgxpool.Pool, cache map[string]int16, code string) (int16, error) {
	if id, ok := cache[code]; ok {
		return id, nil
	}
	id, err := getLanguageID(ctx, pool, code)
	if err != nil {
		return 0, err
	}
	cache[code] = id
	return id, nil
}

// Helper: get topic id by code.
func getTopicID(ctx context.Context, pool *pgxpool.Pool, cache map[string]int64, code string) (int64, error) {
	if id, ok := cache[code]; ok {
		return id, nil
	}
	const q = `SELECT id FROM topics WHERE code = $1`
	var id int64
	if err := pool.QueryRow(ctx, q, code).Scan(&id); err != nil {
		return 0, fmt.Errorf("get topic id for %s: %w", code, err)
	}
	cache[code] = id
	return id, nil
}

// Helper: get level id by code (may be nullable).
func getLevelID(ctx context.Context, pool *pgxpool.Pool, cache map[string]*int64, code *string) (*int64, error) {
	if code == nil || *code == "" {
		return nil, nil
	}
	if id, ok := cache[*code]; ok {
		return id, nil
	}
	const q = `SELECT id FROM levels WHERE code = $1`
	var id int64
	if err := pool.QueryRow(ctx, q, *code).Scan(&id); err != nil {
		return nil, fmt.Errorf("get level id for %s: %w", *code, err)
	}
	cache[*code] = &id
	return &id, nil
}

// Cache key for words (for related/translation targets).
func wordCacheKey(lang, lemma string) string {
	return fmt.Sprintf("%s|%s", lang, lemma)
}

func upsertSingleWord(ctx context.Context, tx pgx.Tx, languageID int16, lemma string, w WordJSON) (int64, error) {
	const selectQ = `
SELECT id
FROM words
WHERE language_id = $1
  AND lemma = $2
`

	const insertQ = `
INSERT INTO words (
    language_id,
    lemma,
    lemma_normalized,
    search_key,
    romanization,
    script_code,
    frequency_rank,
    note
)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
RETURNING id
`

	const updateQ = `
UPDATE words
SET
    lemma_normalized = $2,
    search_key       = $3,
    romanization     = $4,
    script_code      = $5,
    frequency_rank   = $6,
    note             = $7,
    updated_at       = CURRENT_TIMESTAMP
WHERE id = $1
`

	var wordID int64
	err := tx.QueryRow(ctx, selectQ, languageID, lemma).Scan(&wordID)
	if err != nil {
		if err != pgx.ErrNoRows {
			return 0, fmt.Errorf("select existing word: %w", err)
		}

		// Insert new word
		if err := tx.QueryRow(
			ctx,
			insertQ,
			languageID,
			lemma,
			w.LemmaNormalized,
			w.SearchKey,
			w.Romanization,
			w.ScriptCode,
			w.FrequencyRank,
			w.Note,
		).Scan(&wordID); err != nil {
			return 0, fmt.Errorf("insert word: %w", err)
		}

		return wordID, nil
	}

	// Update existing word
	if _, err := tx.Exec(
		ctx,
		updateQ,
		wordID,
		w.LemmaNormalized,
		w.SearchKey,
		w.Romanization,
		w.ScriptCode,
		w.FrequencyRank,
		w.Note,
	); err != nil {
		return 0, fmt.Errorf("update word: %w", err)
	}

	return wordID, nil
}

// upsertWordDetails handles all dictionary structures that depend on a word:
// topics, pronunciations, senses, translations, relations, examples, characters, etc.
func upsertWordDetails(
	ctx context.Context,
	pool *pgxpool.Pool,
	tx pgx.Tx,
	wordID int64,
	w WordJSON,
	languageCache map[string]int16,
	posCache map[string]*int16,
	topicCache map[string]int64,
	levelCache map[string]*int64,
	wordCache map[string]int64,
	characterCache map[string]int64,
) error {
	// Topics
	if err := upsertWordTopics(ctx, pool, tx, wordID, w.Topics, topicCache); err != nil {
		return err
	}

	// Pronunciations
	if err := upsertWordPronunciations(ctx, tx, wordID, w.Pronunciations); err != nil {
		return err
	}

	// Senses, translations, examples
	if err := upsertWordSenses(
		ctx,
		pool,
		tx,
		wordID,
		w.Senses,
		languageCache,
		levelCache,
		wordCache,
		posCache,
	); err != nil {
		return err
	}

	// Relations
	if err := upsertWordRelations(
		ctx,
		pool,
		tx,
		wordID,
		w.Relations,
		languageCache,
		wordCache,
	); err != nil {
		return err
	}

	// Characters (mainly for Chinese)
	if len(w.Characters) > 0 {
		if err := upsertWordCharacters(
			ctx,
			pool,
			tx,
			wordID,
			w.Characters,
			languageCache,
			characterCache,
			levelCache,
		); err != nil {
			return err
		}
	}

	return nil
}

// --------- Word topics ----------

func upsertWordTopics(
	ctx context.Context,
	pool *pgxpool.Pool,
	tx pgx.Tx,
	wordID int64,
	topicCodes []string,
	topicCache map[string]int64,
) error {
	if len(topicCodes) == 0 {
		return nil
	}

	const insertQ = `
INSERT INTO word_topics (word_id, topic_id)
VALUES ($1, $2)
ON CONFLICT (word_id, topic_id) DO NOTHING
`

	for _, code := range topicCodes {
		if code == "" {
			continue
		}
		topicID, err := getTopicID(ctx, pool, topicCache, code)
		if err != nil {
			return err
		}
		if _, err := tx.Exec(ctx, insertQ, wordID, topicID); err != nil {
			return fmt.Errorf("insert word_topic (%d, %d): %w", wordID, topicID, err)
		}
	}
	return nil
}

// --------- Pronunciations ----------

func upsertWordPronunciations(
	ctx context.Context,
	tx pgx.Tx,
	wordID int64,
	prons []PronunciationJSON,
) error {
	if len(prons) == 0 {
		return nil
	}

	const insertQ = `
INSERT INTO pronunciations (word_id, dialect, ipa, phonetic, audio_url)
VALUES ($1, $2, $3, $4, $5)
ON CONFLICT (word_id, dialect) DO UPDATE
SET ipa = EXCLUDED.ipa,
    phonetic = EXCLUDED.phonetic,
    audio_url = EXCLUDED.audio_url
`

	for _, p := range prons {
		if _, err := tx.Exec(ctx, insertQ, wordID, p.Dialect, p.IPA, p.Phonetic, p.AudioURL); err != nil {
			return fmt.Errorf("upsert pronunciation: %w", err)
		}
	}

	return nil
}

// --------- Senses, translations, examples ----------

func upsertWordSenses(
	ctx context.Context,
	pool *pgxpool.Pool,
	tx pgx.Tx,
	wordID int64,
	senses []SenseJSON,
	languageCache map[string]int16,
	levelCache map[string]*int64,
	wordCache map[string]int64,
	posCache map[string]*int16,
) error {
	if len(senses) == 0 {
		return nil
	}

	const insertSenseQ = `
INSERT INTO senses (word_id, sense_order, part_of_speech_id, definition, definition_language_id, usage_label, level_id, note)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
ON CONFLICT (word_id, sense_order) DO UPDATE
SET part_of_speech_id = EXCLUDED.part_of_speech_id,
    definition = EXCLUDED.definition,
    usage_label = EXCLUDED.usage_label,
    level_id = EXCLUDED.level_id,
    note = EXCLUDED.note
RETURNING id
`

	for _, s := range senses {
		defLangID, err := getLanguageIDWithCache(ctx, pool, languageCache, s.DefinitionLanguage)
		if err != nil {
			return err
		}

		levelID, err := getLevelID(ctx, pool, levelCache, s.Level)
		if err != nil {
			return err
		}

		posID, err := getPartOfSpeechID(ctx, tx, posCache, s.PartOfSpeech)
		if err != nil {
			return fmt.Errorf("get part_of_speech_id (%s): %w", s.PartOfSpeech, err)
		}
		if posID == nil {
			return fmt.Errorf("part_of_speech is required for sense but was empty or not found: %s", s.PartOfSpeech)
		}

		var senseID int64
		if err := tx.QueryRow(
			ctx,
			insertSenseQ,
			wordID,
			s.Order,
			posID,
			s.Definition,
			defLangID,
			s.UsageLabel,
			levelID,
			s.Note,
		).Scan(&senseID); err != nil {
			return fmt.Errorf("upsert sense: %w", err)
		}

		// Translations for this sense
		if err := upsertSenseTranslations(
			ctx,
			pool,
			tx,
			senseID,
			s.Translations,
			languageCache,
			wordCache,
		); err != nil {
			return err
		}

		// Examples for this sense
		if err := upsertSenseExamples(
			ctx,
			pool,
			tx,
			senseID,
			s.Examples,
			languageCache,
		); err != nil {
			return err
		}
	}

	return nil
}

func upsertSenseTranslations(
	ctx context.Context,
	pool *pgxpool.Pool,
	tx pgx.Tx,
	senseID int64,
	translations []SenseTranslationJSON,
	languageCache map[string]int16,
	wordCache map[string]int64,
) error {
	if len(translations) == 0 {
		return nil
	}

	const insertQ = `
INSERT INTO sense_translations (source_sense_id, target_word_id, priority, note)
VALUES ($1, $2, $3, $4)
ON CONFLICT (source_sense_id, target_word_id) DO UPDATE
SET priority = EXCLUDED.priority,
    note = EXCLUDED.note
RETURNING id
`

	for _, t := range translations {
		targetLangID, err := getLanguageIDWithCache(ctx, pool, languageCache, t.TargetWord.Language)
		if err != nil {
			return err
		}

		targetWordID, err := upsertRelatedWord(ctx, tx, targetLangID, t.TargetWord, wordCache)
		if err != nil {
			return err
		}

		if _, err := tx.Exec(
			ctx,
			insertQ,
			senseID,
			targetWordID,
			t.Priority,
			t.Note,
		); err != nil {
			return fmt.Errorf("upsert sense_translation: %w", err)
		}
	}

	return nil
}

func upsertSenseExamples(
	ctx context.Context,
	pool *pgxpool.Pool,
	tx pgx.Tx,
	senseID int64,
	examples []ExampleJSON,
	languageCache map[string]int16,
) error {
	if len(examples) == 0 {
		return nil
	}

	const selectExampleQ = `
SELECT id
FROM examples
WHERE source_sense_id = $1 AND language_id = $2 AND content = $3
`

	const insertExampleQ = `
INSERT INTO examples (source_sense_id, language_id, content, audio_url, source)
VALUES ($1, $2, $3, $4, NULL)
RETURNING id
`

	const updateExampleQ = `
UPDATE examples
SET audio_url = $2
WHERE id = $1
`

	const insertTransQ = `
INSERT INTO example_translations (example_id, language_id, content)
VALUES ($1, $2, $3)
ON CONFLICT (example_id, language_id) DO UPDATE
SET content = EXCLUDED.content
`

	for _, ex := range examples {
		langID, err := getLanguageIDWithCache(ctx, pool, languageCache, ex.Language)
		if err != nil {
			return err
		}

		var exampleID int64
		err = tx.QueryRow(ctx, selectExampleQ, senseID, langID, ex.Content).Scan(&exampleID)
		if err != nil {
			if err == pgx.ErrNoRows {
				if err := tx.QueryRow(
					ctx,
					insertExampleQ,
					senseID,
					langID,
					ex.Content,
					ex.AudioURL,
				).Scan(&exampleID); err != nil {
					return fmt.Errorf("insert example: %w", err)
				}
			} else {
				return fmt.Errorf("select example: %w", err)
			}
		} else {
			if _, err := tx.Exec(ctx, updateExampleQ, exampleID, ex.AudioURL); err != nil {
				return fmt.Errorf("update example: %w", err)
			}
		}

		for _, tr := range ex.Translations {
			trLangID, err := getLanguageIDWithCache(ctx, pool, languageCache, tr.Language)
			if err != nil {
				return err
			}

			if _, err := tx.Exec(
				ctx,
				insertTransQ,
				exampleID,
				trLangID,
				tr.Content,
			); err != nil {
				return fmt.Errorf("upsert example_translation: %w", err)
			}
		}
	}

	return nil
}

// --------- Word relations ----------

func upsertWordRelations(
	ctx context.Context,
	pool *pgxpool.Pool,
	tx pgx.Tx,
	fromWordID int64,
	relations []WordRelationJSON,
	languageCache map[string]int16,
	wordCache map[string]int64,
) error {
	if len(relations) == 0 {
		return nil
	}

	const insertQ = `
INSERT INTO word_relations (from_word_id, to_word_id, relation_type, note)
VALUES ($1, $2, $3, $4)
ON CONFLICT (from_word_id, to_word_id, relation_type) DO UPDATE
SET note = EXCLUDED.note
`

	for _, r := range relations {
		targetLangID, err := getLanguageIDWithCache(ctx, pool, languageCache, r.TargetWord.Language)
		if err != nil {
			return err
		}

		targetWordID, err := upsertRelatedWord(ctx, tx, targetLangID, r.TargetWord, wordCache)
		if err != nil {
			return err
		}

		// Skip if from_word_id == to_word_id (CHECK constraint will prevent this)
		if fromWordID == targetWordID {
			continue
		}

		if _, err := tx.Exec(
			ctx,
			insertQ,
			fromWordID,
			targetWordID,
			r.RelationType,
			r.Note,
		); err != nil {
			return fmt.Errorf("upsert word_relation: %w", err)
		}
	}

	return nil
}

// upsertRelatedWord ensures target/related words exist in the words table and returns their id.
func upsertRelatedWord(ctx context.Context, tx pgx.Tx, languageID int16, w RelatedWordJSON, wordCache map[string]int64) (int64, error) {
	key := wordCacheKey(w.Language, w.Lemma)
	if id, ok := wordCache[key]; ok {
		return id, nil
	}

	const selectQ = `
SELECT id
FROM words
WHERE language_id = $1 AND lemma = $2
`

	const insertQ = `
INSERT INTO words (
    language_id,
    lemma,
    lemma_normalized,
    search_key,
    romanization,
    script_code,
    frequency_rank,
    note
)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
RETURNING id
`

	var id int64
	err := tx.QueryRow(ctx, selectQ, languageID, w.Lemma).Scan(&id)
	if err != nil {
		if err != pgx.ErrNoRows {
			return 0, fmt.Errorf("select related word: %w", err)
		}
		if err := tx.QueryRow(
			ctx,
			insertQ,
			languageID,
			w.Lemma,
			w.LemmaNormalized,
			w.SearchKey,
			w.Romanization,
			w.ScriptCode,
			w.FrequencyRank,
			w.Note,
		).Scan(&id); err != nil {
			return 0, fmt.Errorf("insert related word: %w", err)
		}
	}

	wordCache[key] = id
	return id, nil
}

// --------- Characters & readings ----------

func upsertWordCharacters(
	ctx context.Context,
	pool *pgxpool.Pool,
	tx pgx.Tx,
	wordID int64,
	chars []CharacterJSON,
	languageCache map[string]int16,
	characterCache map[string]int64,
	levelCache map[string]*int64,
) error {
	if len(chars) == 0 {
		return nil
	}

	const selectCharQ = `
SELECT id
FROM characters
WHERE literal = $1 AND script_code = $2
`

	const insertCharQ = `
INSERT INTO characters (literal, simplified, traditional, script_code, strokes, radical, level_id)
VALUES ($1, $2, $3, $4, $5, $6, $7)
RETURNING id
`

	const insertWordCharQ = `
INSERT INTO word_characters (word_id, character_id, char_order)
VALUES ($1, $2, $3)
ON CONFLICT (word_id, char_order) DO UPDATE SET character_id = EXCLUDED.character_id
`

	const selectReadingQ = `
SELECT id
FROM character_readings
WHERE character_id = $1 AND language_id = $2 AND reading = $3 AND COALESCE(reading_type, '') = COALESCE($4, '')
`

	const insertReadingQ = `
INSERT INTO character_readings (character_id, language_id, reading, reading_type, note)
VALUES ($1, $2, $3, $4, $5)
RETURNING id
`

	const updateReadingQ = `
UPDATE character_readings
SET note = $2
WHERE id = $1
`

	for _, c := range chars {
		if c.Literal == "" {
			continue
		}

		charKey := fmt.Sprintf("%s|%s", c.Literal, c.ScriptCode)
		var charID int64
		if cachedID, ok := characterCache[charKey]; ok {
			charID = cachedID
		} else {
			err := tx.QueryRow(ctx, selectCharQ, c.Literal, c.ScriptCode).Scan(&charID)
			if err != nil {
				if err == pgx.ErrNoRows {
					// Convert level code to level_id
					var levelID *int64
					if c.Level != nil && *c.Level != "" {
						levelID, err = getLevelID(ctx, pool, levelCache, c.Level)
						if err != nil {
							return fmt.Errorf("get level_id for character level %s: %w", *c.Level, err)
						}
					}

					if err := tx.QueryRow(
						ctx,
						insertCharQ,
						c.Literal,
						c.Simplified,
						c.Traditional,
						c.ScriptCode,
						c.Strokes,
						c.Radical,
						levelID,
					).Scan(&charID); err != nil {
						return fmt.Errorf("insert character: %w", err)
					}
				} else {
					return fmt.Errorf("select character: %w", err)
				}
			}
			characterCache[charKey] = charID
		}

		// Link word to character with order
		if _, err := tx.Exec(ctx, insertWordCharQ, wordID, charID, c.CharOrder); err != nil {
			return fmt.Errorf("insert word_character: %w", err)
		}

		// Readings for this character
		for _, r := range c.Readings {
			langID, err := getLanguageIDWithCache(ctx, pool, languageCache, r.Language)
			if err != nil {
				return err
			}

			var readingID int64
			err = tx.QueryRow(ctx, selectReadingQ, charID, langID, r.Reading, r.ReadingType).Scan(&readingID)
			if err != nil {
				if err == pgx.ErrNoRows {
					if err := tx.QueryRow(
						ctx,
						insertReadingQ,
						charID,
						langID,
						r.Reading,
						r.ReadingType,
						r.Note,
					).Scan(&readingID); err != nil {
						return fmt.Errorf("insert character_reading: %w", err)
					}
				} else {
					return fmt.Errorf("select character_reading: %w", err)
				}
			} else {
				if _, err := tx.Exec(ctx, updateReadingQ, readingID, r.Note); err != nil {
					return fmt.Errorf("update character_reading: %w", err)
				}
			}
		}
	}

	return nil
}
