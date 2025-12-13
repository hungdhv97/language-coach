# Data Model: Multilingual Dictionary with Vocabulary Game

**Created**: 2025-01-27  
**Feature**: [spec.md](./spec.md)

This document describes the data model for the multilingual dictionary and vocabulary game feature, based on the database schema defined in `backend/internal/infrastructure/db/migrations/0001_init.sql`.

## Core Dictionary Entities

### Language

Represents a supported language in the system.

**Table**: `languages`

**Fields**:
- `id` (SMALLINT, Primary Key, AUTO_INCREMENT): Unique identifier
- `code` (VARCHAR(10), NOT NULL, UNIQUE): Language code (e.g., 'en', 'vi', 'zh')
- `name` (VARCHAR(100), NOT NULL): Language name in English (e.g., 'English')

**Relationships**:
- One-to-many: Words, Senses, Game Sessions
- Many-to-many: Used as source/target in translations

**Validation Rules**:
- Language code must be unique
- Language code format: ISO 639-1 or similar standard codes
- Name is required

**Indexes**: Primary key on `id`, unique on `code`

---

### Part of Speech

Represents grammatical categories (noun, verb, adjective, etc.).

**Table**: `parts_of_speech`

**Fields**:
- `id` (SMALLINT, Primary Key, AUTO_INCREMENT): Unique identifier
- `code` (VARCHAR(20), NOT NULL, UNIQUE): Part of speech code (e.g., 'noun', 'verb')
- `name` (VARCHAR(100), NOT NULL): Display name

**Relationships**:
- One-to-many: Words (via `part_of_speech_id`)

**Validation Rules**:
- Code must be unique
- Name is required

**Indexes**: Primary key on `id`, unique on `code`

---

### Topic

Represents thematic categories for organizing vocabulary.

**Table**: `topics`

**Fields**:
- `id` (BIGINT, Primary Key, AUTO_INCREMENT): Unique identifier
- `code` (VARCHAR(50), NOT NULL, UNIQUE): Topic code (e.g., 'education', 'travel')
- `name` (VARCHAR(100), NOT NULL): Display name

**Relationships**:
- Many-to-many: Words (via `word_topics` join table)
- One-to-many: User Topic Statistics, Game Sessions

**Validation Rules**:
- Code must be unique
- Name is required

**Indexes**: Primary key on `id`, unique on `code`

---

### Level

Represents difficulty or proficiency levels (e.g., HSK1, A1, N3).

**Table**: `levels`

**Fields**:
- `id` (BIGINT, Primary Key, AUTO_INCREMENT): Unique identifier
- `code` (VARCHAR(50), NOT NULL, UNIQUE): Level code (e.g., 'HSK1', 'A1', 'N3')
- `name` (VARCHAR(100), NOT NULL): Display name
- `description` (TEXT, NULLABLE): Level description
- `language_id` (SMALLINT, NULLABLE, FK → languages.id): Associated language (null if level is general)
- `difficulty_order` (SMALLINT, NULLABLE): Order of difficulty (1 < 2 < 3...)

**Relationships**:
- Many-to-one: Language (optional)
- One-to-many: Senses, Game Sessions

**Validation Rules**:
- Code must be unique
- Name is required
- If `language_id` is set, it must reference a valid language
- `difficulty_order` should be consistent within same `language_id`

**Indexes**: 
- Primary key on `id`
- Unique on `code`
- Composite index on `(language_id, difficulty_order)`

---

## Word Entities

### Word

Represents a vocabulary word in a specific language.

**Table**: `words`

**Fields**:
- `id` (BIGINT, Primary Key, AUTO_INCREMENT): Unique identifier
- `language_id` (SMALLINT, NOT NULL, FK → languages.id): Language of the word
- `lemma` (VARCHAR(255), NOT NULL): Base form (headword)
- `lemma_normalized` (VARCHAR(255), NULLABLE): Normalized form (lowercase, no accents)
- `search_key` (VARCHAR(255), NULLABLE): Search key (pinyin, no accents, etc.)
- `part_of_speech_id` (SMALLINT, NULLABLE, FK → parts_of_speech.id): Part of speech
- `romanization` (VARCHAR(255), NULLABLE): Latin transcription (pinyin, Hán-Việt, etc.)
- `script_code` (VARCHAR(20), NULLABLE): Script code (e.g., 'Latn', 'Hani')
- `frequency_rank` (INT, NULLABLE): Frequency/popularity rank
- `notes` (TEXT, NULLABLE): Additional notes
- `created_at` (DATETIME, DEFAULT CURRENT_TIMESTAMP): Creation timestamp
- `updated_at` (DATETIME, DEFAULT CURRENT_TIMESTAMP ON UPDATE): Last update timestamp

**Relationships**:
- Many-to-one: Language, Part of Speech
- One-to-many: Senses, Translations, Game Questions
- Many-to-many: Topics (via `word_topics`)
- Many-to-many: Characters (via `word_characters`)
- Many-to-many: Related words (via `word_relations`)

**Validation Rules**:
- `language_id` is required and must reference valid language
- `lemma` is required
- `lemma_normalized` and `search_key` should be generated automatically
- `frequency_rank` should be positive integer if provided

**Indexes**:
- Primary key on `id`
- Composite index on `(language_id, lemma)`
- Composite index on `(language_id, lemma_normalized)`
- Composite index on `(language_id, search_key)`

---

### Sense

Represents a specific meaning or definition of a word.

**Table**: `senses`

**Fields**:
- `id` (BIGINT, Primary Key, AUTO_INCREMENT): Unique identifier
- `word_id` (BIGINT, NOT NULL, FK → words.id): Parent word
- `sense_order` (SMALLINT, NOT NULL): Order of this sense (1, 2, 3...)
- `definition` (TEXT, NOT NULL): Definition text
- `definition_language_id` (SMALLINT, NOT NULL, FK → languages.id): Language of definition
- `usage_label` (VARCHAR(100), NULLABLE): Usage label (e.g., 'figurative', 'slang')
- `level_id` (BIGINT, NULLABLE, FK → levels.id): Associated level
- `note` (TEXT, NULLABLE): Additional notes

**Relationships**:
- Many-to-one: Word, Definition Language, Level
- One-to-many: Translations, Examples, Game Questions

**Validation Rules**:
- `word_id` is required
- `sense_order` must be unique within a word
- `definition` is required
- `definition_language_id` is required

**Indexes**:
- Primary key on `id`
- Composite index on `(word_id, sense_order)`

---

### Sense Translation

Represents translations between senses and words in other languages.

**Table**: `sense_translations`

**Fields**:
- `id` (BIGINT, Primary Key, AUTO_INCREMENT): Unique identifier
- `source_sense_id` (BIGINT, NOT NULL, FK → senses.id): Source sense
- `target_word_id` (BIGINT, NOT NULL, FK → words.id): Target word
- `target_language_id` (SMALLINT, NOT NULL, FK → languages.id): Target language
- `priority` (SMALLINT, DEFAULT 1): Display priority (1 = highest)
- `note` (TEXT, NULLABLE): Translation notes

**Relationships**:
- Many-to-one: Source Sense, Target Word, Target Language

**Validation Rules**:
- `priority` should be positive integer
- Lower priority values indicate higher priority (1 > 2 > 3)

**Indexes**:
- Primary key on `id`
- Index on `source_sense_id`
- Composite index on `(target_language_id, target_word_id)`

---

### Word Relation

Represents relationships between words (synonyms, antonyms, related words).

**Table**: `word_relations`

**Fields**:
- `id` (BIGINT, Primary Key, AUTO_INCREMENT): Unique identifier
- `from_word_id` (BIGINT, NOT NULL, FK → words.id): Source word
- `to_word_id` (BIGINT, NOT NULL, FK → words.id): Target word
- `relation_type` (VARCHAR(20), NOT NULL): Type ('synonym', 'antonym', 'related')
- `note` (TEXT, NULLABLE): Relation notes

**Relationships**:
- Many-to-one: From Word, To Word

**Validation Rules**:
- `relation_type` must be one of: 'synonym', 'antonym', 'related'
- `from_word_id` and `to_word_id` must be different
- Relation should be bidirectional for synonyms/antonyms (separate records)

**Indexes**:
- Primary key on `id`
- Composite index on `(from_word_id, relation_type)`
- Composite index on `(to_word_id, relation_type)`

---

### Word Topic

Join table linking words to topics.

**Table**: `word_topics`

**Fields**:
- `word_id` (BIGINT, NOT NULL, FK → words.id): Word
- `topic_id` (BIGINT, NOT NULL, FK → topics.id): Topic

**Relationships**:
- Many-to-one: Word, Topic

**Validation Rules**:
- Both `word_id` and `topic_id` are required
- Combination must be unique

**Indexes**:
- Composite primary key on `(word_id, topic_id)`

---

## Example and Pronunciation Entities

### Example

Represents example sentences illustrating word usage.

**Table**: `examples`

**Fields**:
- `id` (BIGINT, Primary Key, AUTO_INCREMENT): Unique identifier
- `source_sense_id` (BIGINT, NOT NULL, FK → senses.id): Sense being illustrated
- `language_id` (SMALLINT, NOT NULL, FK → languages.id): Language of example
- `content` (TEXT, NOT NULL): Example sentence text
- `audio_url` (VARCHAR(500), NULLABLE): Audio file URL
- `source` (VARCHAR(255), NULLABLE): Source attribution (book, movie, etc.)

**Relationships**:
- Many-to-one: Sense, Language
- One-to-many: Translations

**Validation Rules**:
- `content` is required
- `audio_url` must be valid URL format if provided

**Indexes**:
- Primary key on `id`
- Index on `source_sense_id`
- Index on `language_id`

---

### Example Translation

Translations of example sentences.

**Table**: `example_translations`

**Fields**:
- `id` (BIGINT, Primary Key, AUTO_INCREMENT): Unique identifier
- `example_id` (BIGINT, NOT NULL, FK → examples.id): Source example
- `language_id` (SMALLINT, NOT NULL, FK → languages.id): Translation language
- `content` (TEXT, NOT NULL): Translated sentence

**Relationships**:
- Many-to-one: Example, Language

**Validation Rules**:
- `content` is required
- Translation language should differ from example language

**Indexes**:
- Primary key on `id`
- Index on `example_id`

---

### Pronunciation

Represents pronunciation information for words.

**Table**: `pronunciations`

**Fields**:
- `id` (BIGINT, Primary Key, AUTO_INCREMENT): Unique identifier
- `word_id` (BIGINT, NOT NULL, FK → words.id): Word
- `dialect` (VARCHAR(20), NULLABLE): Dialect code (e.g., 'en-US', 'en-UK', 'vi-North')
- `ipa` (VARCHAR(255), NULLABLE): IPA transcription (e.g., '/skuːl/')
- `phonetic` (VARCHAR(255), NULLABLE): Simplified phonetic (e.g., 's-kuul')
- `audio_url` (VARCHAR(500), NULLABLE): Pronunciation audio URL

**Relationships**:
- Many-to-one: Word

**Validation Rules**:
- At least one of `ipa`, `phonetic`, or `audio_url` should be provided
- `audio_url` must be valid URL format if provided

**Indexes**:
- Primary key on `id`
- Index on `word_id`

---

## Character Entities (for Chinese/Japanese)

### Character

Represents individual characters (for logographic languages).

**Table**: `characters`

**Fields**:
- `id` (BIGINT, Primary Key, AUTO_INCREMENT): Unique identifier
- `literal` (VARCHAR(2), NOT NULL): Character (e.g., '学')
- `simplified` (VARCHAR(2), NULLABLE): Simplified form
- `traditional` (VARCHAR(2), NULLABLE): Traditional form
- `script_code` (VARCHAR(10), NOT NULL): Script code (e.g., 'Hani')
- `strokes` (SMALLINT, NULLABLE): Stroke count
- `radical` (VARCHAR(10), NULLABLE): Radical component
- `level` (VARCHAR(20), NULLABLE): Level indicator (e.g., 'HSK1')

**Relationships**:
- Many-to-many: Words (via `word_characters`)
- One-to-many: Character Readings

**Validation Rules**:
- `literal` is required
- `script_code` is required

**Indexes**: Primary key on `id`

---

### Character Reading

Represents pronunciation readings for characters.

**Table**: `character_readings`

**Fields**:
- `id` (BIGINT, Primary Key, AUTO_INCREMENT): Unique identifier
- `character_id` (BIGINT, NOT NULL, FK → characters.id): Character
- `language_id` (SMALLINT, NOT NULL, FK → languages.id): Language of reading
- `reading` (VARCHAR(100), NOT NULL): Reading (pinyin, Hán-Việt, etc.)
- `reading_type` (VARCHAR(50), NULLABLE): Reading type (e.g., 'pinyin', 'sino-vietnamese')
- `note` (TEXT, NULLABLE): Reading notes

**Relationships**:
- Many-to-one: Character, Language

**Validation Rules**:
- `reading` is required

**Indexes**:
- Primary key on `id`
- Index on `character_id`
- Index on `language_id`

---

### Word Character

Join table linking words to characters (for multi-character words).

**Table**: `word_characters`

**Fields**:
- `word_id` (BIGINT, NOT NULL, FK → words.id): Word
- `character_id` (BIGINT, NOT NULL, FK → characters.id): Character
- `char_order` (SMALLINT, NOT NULL): Position in word (1, 2, 3...)

**Relationships**:
- Many-to-one: Word, Character

**Validation Rules**:
- `char_order` must be unique within a word
- Order must be sequential starting from 1

**Indexes**:
- Composite primary key on `(word_id, char_order)`

---

## User Entities

### User

Represents application users.

**Table**: `users`

**Fields**:
- `id` (BIGINT, Primary Key, AUTO_INCREMENT): Unique identifier
- `email` (VARCHAR(255), UNIQUE, NULLABLE): Email address
- `username` (VARCHAR(100), UNIQUE, NULLABLE): Username
- `password_hash` (VARCHAR(255), NULLABLE): Hashed password
- `created_at` (DATETIME, DEFAULT CURRENT_TIMESTAMP): Account creation time
- `updated_at` (DATETIME, DEFAULT CURRENT_TIMESTAMP ON UPDATE): Last update time
- `is_active` (TINYINT(1), DEFAULT 1): Active status (1 = active)

**Relationships**:
- One-to-one: User Profile, User Statistics
- One-to-many: Game Sessions, Game Answers

**Validation Rules**:
- Either `email` or `username` must be provided
- `email` must be valid email format if provided
- `is_active` must be 0 or 1

**Indexes**:
- Primary key on `id`
- Unique on `email`
- Unique on `username`

---

### User Profile

Extended user profile information.

**Table**: `user_profiles`

**Fields**:
- `user_id` (BIGINT, Primary Key, FK → users.id): User reference
- `display_name` (VARCHAR(100), NULLABLE): Display name
- `avatar_url` (VARCHAR(500), NULLABLE): Avatar image URL
- `birth_day` (DATE, NULLABLE): Birth date (YYYY-MM-DD)
- `bio` (TEXT, NULLABLE): User biography
- `created_at` (DATETIME, DEFAULT CURRENT_TIMESTAMP): Profile creation time
- `updated_at` (DATETIME, DEFAULT CURRENT_TIMESTAMP ON UPDATE): Last update time

**Relationships**:
- One-to-one: User

**Validation Rules**:
- `birth_day` must be valid date format if provided
- `avatar_url` must be valid URL format if provided

**Indexes**: Primary key on `user_id`

---

## Statistics Entities

### User Statistics

Aggregated statistics for a user across all game sessions.

**Table**: `user_statistics`

**Fields**:
- `user_id` (BIGINT, Primary Key, FK → users.id): User reference
- `total_sessions` (INT, DEFAULT 0): Total game sessions played
- `total_questions` (INT, DEFAULT 0): Total questions answered
- `total_correct` (INT, DEFAULT 0): Total correct answers
- `total_time_seconds` (INT, DEFAULT 0): Total play time in seconds
- `last_played_at` (DATETIME, NULLABLE): Last game session timestamp

**Relationships**:
- One-to-one: User

**Validation Rules**:
- All count fields must be non-negative
- `total_correct` must be <= `total_questions`
- `last_played_at` should be updated when user plays a game

**Indexes**: Primary key on `user_id`

---

### User Word Statistics

Per-word statistics for a user.

**Table**: `user_word_statistics`

**Fields**:
- `user_id` (BIGINT, NOT NULL, FK → users.id): User reference
- `word_id` (BIGINT, NOT NULL, FK → words.id): Word reference
- `correct_count` (INT, DEFAULT 0): Correct answer count
- `wrong_count` (INT, DEFAULT 0): Wrong answer count
- `last_answered_at` (DATETIME, NULLABLE): Last answer timestamp
- `streak` (INT, DEFAULT 0): Consecutive correct answers

**Relationships**:
- Many-to-one: User, Word

**Validation Rules**:
- Count fields must be non-negative
- `streak` resets to 0 when user answers incorrectly

**Indexes**:
- Composite primary key on `(user_id, word_id)`

---

### User Topic Statistics

Per-topic statistics for a user.

**Table**: `user_topic_statistics`

**Fields**:
- `user_id` (BIGINT, NOT NULL, FK → users.id): User reference
- `topic_id` (BIGINT, NOT NULL, FK → topics.id): Topic reference
- `total_questions` (INT, DEFAULT 0): Total questions in this topic
- `total_correct` (INT, DEFAULT 0): Correct answers in this topic
- `last_played_at` (DATETIME, NULLABLE): Last play timestamp for this topic

**Relationships**:
- Many-to-one: User, Topic

**Validation Rules**:
- Count fields must be non-negative
- `total_correct` must be <= `total_questions`

**Indexes**:
- Composite primary key on `(user_id, topic_id)`

---

## Game Entities

### Game Session

Represents a single vocabulary game playthrough.

**Table**: `vocab_game_sessions`

**Fields**:
- `id` (BIGINT, Primary Key, AUTO_INCREMENT): Unique identifier
- `user_id` (BIGINT, NOT NULL, FK → users.id): User playing the game
- `mode` (VARCHAR(50), NOT NULL): Game mode ('level', 'topic')
- `source_language_id` (SMALLINT, NOT NULL, FK → languages.id): Source language
- `target_language_id` (SMALLINT, NOT NULL, FK → languages.id): Target language
- `topic_id` (BIGINT, NULLABLE, FK → topics.id): Topic filter (if mode = 'topic')
- `level_id` (BIGINT, NULLABLE, FK → levels.id): Level filter (if mode = 'level')
- `total_questions` (SMALLINT, DEFAULT 0): Total questions in session
- `correct_questions` (SMALLINT, DEFAULT 0): Correct answers count
- `started_at` (DATETIME, DEFAULT CURRENT_TIMESTAMP): Session start time
- `ended_at` (DATETIME, NULLABLE): Session end time

**Relationships**:
- Many-to-one: User, Source Language, Target Language, Topic (optional), Level (optional)
- One-to-many: Game Questions, Game Answers

**Validation Rules**:
- `source_language_id` and `target_language_id` must be different
- Either `topic_id` or `level_id` must be set (but not both)
- `mode` must match the set filter ('topic' → topic_id, 'level' → level_id)
- `total_questions` and `correct_questions` must be non-negative
- `correct_questions` must be <= `total_questions`
- `ended_at` should be set when session completes

**Indexes**:
- Primary key on `id`
- Composite index on `(user_id, started_at)`

---

### Game Question

Represents a single question within a game session.

**Table**: `vocab_game_questions`

**Fields**:
- `id` (BIGINT, Primary Key, AUTO_INCREMENT): Unique identifier
- `session_id` (BIGINT, NOT NULL, FK → vocab_game_sessions.id): Parent session
- `question_order` (SMALLINT, NOT NULL): Question order in session (1, 2, 3...)
- `question_type` (VARCHAR(30), NOT NULL): Question type (e.g., 'word_to_translation')
- `source_word_id` (BIGINT, NOT NULL, FK → words.id): Source word
- `source_sense_id` (BIGINT, NULLABLE, FK → senses.id): Specific sense (optional)
- `correct_target_word_id` (BIGINT, NOT NULL, FK → words.id): Correct answer word
- `source_language_id` (SMALLINT, NOT NULL, FK → languages.id): Source language
- `target_language_id` (SMALLINT, NOT NULL, FK → languages.id): Target language
- `created_at` (DATETIME, DEFAULT CURRENT_TIMESTAMP): Question creation time

**Relationships**:
- Many-to-one: Session, Source Word, Source Sense (optional), Correct Target Word, Source Language, Target Language
- One-to-many: Question Options, Answers

**Validation Rules**:
- `question_order` must be unique within a session
- `question_type` must be 'word_to_translation' or similar valid type
- `source_word_id` and `correct_target_word_id` must be different
- Languages must match session languages

**Indexes**:
- Primary key on `id`
- Composite index on `(session_id, question_order)`

---

### Game Question Option

Represents one of the four multiple-choice answers (A, B, C, D) for a question.

**Table**: `vocab_game_question_options`

**Fields**:
- `id` (BIGINT, Primary Key, AUTO_INCREMENT): Unique identifier
- `question_id` (BIGINT, NOT NULL, FK → vocab_game_questions.id): Parent question
- `option_label` (CHAR(1), NOT NULL): Option label ('A', 'B', 'C', 'D')
- `target_word_id` (BIGINT, NOT NULL, FK → words.id): Word displayed as option
- `is_correct` (TINYINT(1), NOT NULL, DEFAULT 0): Correctness flag (1 = correct)

**Relationships**:
- Many-to-one: Question, Target Word

**Validation Rules**:
- `option_label` must be one of: 'A', 'B', 'C', 'D'
- Each question must have exactly 4 options (A, B, C, D)
- Exactly one option per question must have `is_correct = 1`
- `option_label` must be unique within a question

**Indexes**:
- Primary key on `id`
- Unique constraint on `(question_id, option_label)`

---

### Game Answer

Represents a user's answer to a game question.

**Table**: `vocab_game_question_answers`

**Fields**:
- `id` (BIGINT, Primary Key, AUTO_INCREMENT): Unique identifier
- `question_id` (BIGINT, NOT NULL, FK → vocab_game_questions.id): Question answered
- `session_id` (BIGINT, NOT NULL, FK → vocab_game_sessions.id): Game session
- `user_id` (BIGINT, NOT NULL, FK → users.id): User who answered
- `selected_option_id` (BIGINT, NULLABLE, FK → vocab_game_question_options.id): Selected option
- `is_correct` (TINYINT(1), NOT NULL, DEFAULT 0): Correctness (1 = correct)
- `response_time_ms` (INT, NULLABLE): Response time in milliseconds
- `answered_at` (DATETIME, DEFAULT CURRENT_TIMESTAMP): Answer timestamp

**Relationships**:
- Many-to-one: Question, Session, User, Selected Option

**Validation Rules**:
- `is_correct` must match whether `selected_option_id` is the correct option
- `response_time_ms` must be positive if provided
- One answer per question per user per session

**Indexes**:
- Primary key on `id`
- Composite index on `(user_id, answered_at)`

---

## State Transitions

### Game Session States

1. **Created**: Session created, questions generated
2. **In Progress**: Session started, questions being answered
3. **Completed**: All questions answered, session ended

**Transitions**:
- Created → In Progress: When first question is answered
- In Progress → Completed: When last question is answered

### User Statistics Updates

Statistics are updated incrementally as game sessions complete:
- User Statistics: Aggregated from all sessions
- User Word Statistics: Updated per-word as answers are submitted
- User Topic Statistics: Updated per-topic when sessions complete

---

## Data Integrity Constraints

1. **Foreign Key Constraints**: All foreign keys have referential integrity
2. **Unique Constraints**: Language codes, topic codes, level codes, word-topic pairs
3. **Check Constraints**: 
   - Source and target languages must differ in game sessions
   - Correct answer count cannot exceed total questions
   - Exactly one correct option per question
4. **Cascade Rules**: Consider cascade delete for dependent records (e.g., delete game questions when session deleted)

---

## Query Patterns

### Common Queries

1. **Dictionary Lookup**: Search words by language, lemma, normalized form, or search key
2. **Game Question Generation**: Select words filtered by topic/level and language pair
3. **Statistics Retrieval**: Aggregate statistics by user, word, or topic
4. **Game History**: List game sessions for a user, ordered by time

### Index Usage

Indexes are optimized for:
- Dictionary search: `(language_id, lemma_normalized)`, `(language_id, search_key)`
- Game question generation: Word filtering by topic/level
- Statistics queries: User-based lookups
- Game session history: `(user_id, started_at)`

