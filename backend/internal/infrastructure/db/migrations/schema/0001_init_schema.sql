-- PostgreSQL Migration: Initial Schema

CREATE TABLE languages (
    id           SMALLSERIAL PRIMARY KEY, -- language id
    code         VARCHAR(10) NOT NULL UNIQUE, -- language code: 'en', 'vi', 'zh', ...
    name         VARCHAR(100) NOT NULL, -- language name (English)
    native_name  VARCHAR(100) -- native name: 'Tiếng Việt', '中文'
);

CREATE TABLE parts_of_speech (
    id      SMALLSERIAL PRIMARY KEY, -- part-of-speech id
    code    VARCHAR(20) NOT NULL UNIQUE, -- part-of-speech code: 'noun', 'verb', ...
    name    VARCHAR(100) NOT NULL -- part-of-speech name
);

CREATE TABLE topics (
    id      BIGSERIAL PRIMARY KEY, -- topic id
    code    VARCHAR(50) NOT NULL UNIQUE, -- topic code: 'education', 'travel', ...
    name    VARCHAR(100) NOT NULL -- topic name
);

CREATE TABLE levels (
    id                BIGSERIAL PRIMARY KEY, -- level id
    code              VARCHAR(50) NOT NULL, -- level code: 'HSK1', 'A1', 'N3', ...
    name              VARCHAR(100) NOT NULL, -- display name for the level
    description       TEXT, -- level description
    language_id       SMALLINT, -- FK -> languages.id (null if shared level)
    difficulty_order  SMALLINT, -- difficulty order (1 < 2 < 3 ...)
    CONSTRAINT fk_levels_lang
        FOREIGN KEY (language_id) REFERENCES languages(id),
    UNIQUE (language_id, code) -- unique code per language (null language_id = shared level)
);

CREATE INDEX idx_levels_lang ON levels(language_id, difficulty_order);

CREATE TABLE words (
    id                   BIGSERIAL PRIMARY KEY, -- word id
    language_id          SMALLINT NOT NULL, -- FK -> languages.id
    lemma                VARCHAR(255) NOT NULL, -- base form (headword)
    lemma_normalized     VARCHAR(255), -- normalized form (no diacritics, lower-case)
    search_key           VARCHAR(255), -- search key (pinyin, non-diacritics, ...)
    romanization         VARCHAR(255), -- latin transcription (pinyin, Sino-Vietnamese, ...)
    script_code          VARCHAR(20), -- script code: 'Latn', 'Hani', ...
    frequency_rank       INTEGER, -- frequency/popularity rank
    notes                TEXT, -- notes
    created_at           TIMESTAMP DEFAULT CURRENT_TIMESTAMP, -- created at
    updated_at           TIMESTAMP DEFAULT CURRENT_TIMESTAMP, -- updated at
    CONSTRAINT fk_words_language
        FOREIGN KEY (language_id) REFERENCES languages(id)
);

CREATE INDEX idx_words_lang_lemma ON words(language_id, lemma);
CREATE INDEX idx_words_lang_norm ON words(language_id, lemma_normalized);
CREATE INDEX idx_words_lang_search ON words(language_id, search_key);

CREATE TABLE senses (
    id                     BIGSERIAL PRIMARY KEY, -- sense id
    word_id                BIGINT NOT NULL, -- FK -> words.id
    sense_order            SMALLINT NOT NULL, -- order of the sense for a word
    part_of_speech_id      SMALLINT NOT NULL, -- FK -> parts_of_speech.id
    definition             TEXT NOT NULL, -- sense definition
    definition_language_id SMALLINT NOT NULL, -- FK -> languages.id (language of definition)
    usage_label            VARCHAR(100), -- usage label: 'figurative', 'slang', ...
    level_id               BIGINT, -- FK -> levels.id (difficulty level)
    note                   TEXT, -- notes
    CONSTRAINT fk_senses_word
        FOREIGN KEY (word_id) REFERENCES words(id),
    CONSTRAINT fk_senses_pos
        FOREIGN KEY (part_of_speech_id) REFERENCES parts_of_speech(id),
    CONSTRAINT fk_senses_def_lang
        FOREIGN KEY (definition_language_id) REFERENCES languages(id),
    CONSTRAINT fk_senses_level
        FOREIGN KEY (level_id) REFERENCES levels(id),
    UNIQUE (word_id, sense_order)
);

CREATE INDEX idx_senses_word_order ON senses(word_id, sense_order);

CREATE TABLE sense_translations (
    id                 BIGSERIAL PRIMARY KEY, -- sense translation id
    source_sense_id    BIGINT NOT NULL, -- FK -> senses.id (source sense)
    target_word_id     BIGINT NOT NULL, -- FK -> words.id (target word)
    priority           SMALLINT DEFAULT 1, -- display priority (1 = highest)
    note               TEXT, -- notes
    CONSTRAINT fk_st_source_sense
        FOREIGN KEY (source_sense_id) REFERENCES senses(id),
    CONSTRAINT fk_st_target_word
        FOREIGN KEY (target_word_id) REFERENCES words(id),
    UNIQUE (source_sense_id, target_word_id)
);

CREATE INDEX idx_st_source ON sense_translations(source_sense_id);
CREATE INDEX idx_st_target_word ON sense_translations(target_word_id);

CREATE TABLE word_relations (
    id            BIGSERIAL PRIMARY KEY, -- word relation id
    from_word_id  BIGINT NOT NULL, -- FK -> words.id (source word)
    to_word_id    BIGINT NOT NULL, -- FK -> words.id (target word)
    relation_type VARCHAR(20) NOT NULL, -- relation type: 'synonym', 'antonym', 'related'
    note          TEXT, -- notes
    CONSTRAINT fk_wr_from
        FOREIGN KEY (from_word_id) REFERENCES words(id),
    CONSTRAINT fk_wr_to
        FOREIGN KEY (to_word_id) REFERENCES words(id),
    CONSTRAINT chk_wr_different_words
        CHECK (from_word_id <> to_word_id),
    UNIQUE (from_word_id, to_word_id, relation_type)
);

CREATE INDEX idx_wr_from ON word_relations(from_word_id, relation_type);
CREATE INDEX idx_wr_to ON word_relations(to_word_id, relation_type);

CREATE TABLE word_topics (
    word_id  BIGINT NOT NULL, -- FK -> words.id
    topic_id BIGINT NOT NULL, -- FK -> topics.id
    PRIMARY KEY (word_id, topic_id),
    CONSTRAINT fk_wt_word
        FOREIGN KEY (word_id) REFERENCES words(id),
    CONSTRAINT fk_wt_topic
        FOREIGN KEY (topic_id) REFERENCES topics(id)
);

CREATE TABLE examples (
    id              BIGSERIAL PRIMARY KEY, -- example sentence id
    source_sense_id BIGINT NOT NULL, -- FK -> senses.id (sense being illustrated)
    language_id     SMALLINT NOT NULL, -- FK -> languages.id (language of the sentence)
    content         TEXT NOT NULL, -- content of the example sentence
    audio_url       VARCHAR(500), -- audio URL for the sentence (if any)
    source          VARCHAR(255), -- source of the sentence (book, movie, ...)
    CONSTRAINT fk_examples_sense
        FOREIGN KEY (source_sense_id) REFERENCES senses(id),
    CONSTRAINT fk_examples_lang
        FOREIGN KEY (language_id) REFERENCES languages(id)
);

CREATE INDEX idx_examples_sense ON examples(source_sense_id);
CREATE INDEX idx_examples_lang ON examples(language_id);

CREATE TABLE example_translations (
    id          BIGSERIAL PRIMARY KEY, -- example translation id
    example_id  BIGINT NOT NULL, -- FK -> examples.id (original sentence)
    language_id SMALLINT NOT NULL, -- FK -> languages.id (translation language)
    content     TEXT NOT NULL, -- content of the translation
    CONSTRAINT fk_ext_example
        FOREIGN KEY (example_id) REFERENCES examples(id),
    CONSTRAINT fk_ext_lang
        FOREIGN KEY (language_id) REFERENCES languages(id),
    UNIQUE (example_id, language_id)
);

CREATE INDEX idx_ext_example ON example_translations(example_id);

CREATE TABLE pronunciations (
    id        BIGSERIAL PRIMARY KEY, -- pronunciation id
    word_id   BIGINT NOT NULL, -- FK -> words.id (corresponding word)
    dialect   VARCHAR(20), -- dialect: 'en-US', 'en-UK', 'vi-North', ...
    ipa       VARCHAR(255), -- IPA transcription: /skuːl/
    phonetic  VARCHAR(255), -- easier-to-read phonetic form: 's-kuul'
    audio_url VARCHAR(500), -- pronunciation audio URL
    CONSTRAINT fk_pron_word
        FOREIGN KEY (word_id) REFERENCES words(id),
    UNIQUE (word_id, dialect)
);

CREATE INDEX idx_pron_word ON pronunciations(word_id);

CREATE TABLE characters (
    id          BIGSERIAL PRIMARY KEY, -- character id
    literal     VARCHAR(2) NOT NULL, -- character: '学'
    simplified  VARCHAR(2), -- simplified form
    traditional VARCHAR(2), -- traditional form
    script_code VARCHAR(10) NOT NULL, -- script code: 'Hani' (Chinese), ...
    strokes     SMALLINT, -- stroke count
    radical     VARCHAR(10), -- radical
    level       VARCHAR(20) -- level: 'HSK1', 'HSK2', ...
);

CREATE TABLE character_readings (
    id           BIGSERIAL PRIMARY KEY, -- character reading id
    character_id BIGINT NOT NULL, -- FK -> characters.id
    language_id  SMALLINT NOT NULL, -- FK -> languages.id
    reading      VARCHAR(100) NOT NULL, -- reading: pinyin, Sino-Vietnamese, ...
    reading_type VARCHAR(50), -- reading type: 'pinyin', 'sino-vietnamese', ...
    note         TEXT, -- notes
    CONSTRAINT fk_cr_char
        FOREIGN KEY (character_id) REFERENCES characters(id),
    CONSTRAINT fk_cr_lang
        FOREIGN KEY (language_id) REFERENCES languages(id)
);

CREATE INDEX idx_cr_char ON character_readings(character_id);
CREATE INDEX idx_cr_lang ON character_readings(language_id);

CREATE TABLE word_characters (
    word_id      BIGINT NOT NULL, -- FK -> words.id (typically Chinese words)
    character_id BIGINT NOT NULL, -- FK -> characters.id
    char_order   SMALLINT NOT NULL, -- position of the character in the word: 1,2,3,...
    PRIMARY KEY (word_id, char_order),
    CONSTRAINT fk_wc_word
        FOREIGN KEY (word_id) REFERENCES words(id),
    CONSTRAINT fk_wc_char
        FOREIGN KEY (character_id) REFERENCES characters(id)
);

CREATE TABLE users (
    id             BIGSERIAL PRIMARY KEY, -- user id
    email          VARCHAR(255) UNIQUE, -- login email (may be null if other login methods are used)
    username       VARCHAR(100) UNIQUE, -- username
    password_hash  VARCHAR(255), -- hashed password
    created_at     TIMESTAMP DEFAULT CURRENT_TIMESTAMP, -- account creation time
    updated_at     TIMESTAMP DEFAULT CURRENT_TIMESTAMP, -- last update time
    is_active      BOOLEAN DEFAULT TRUE -- activation status
);

CREATE TABLE user_profiles (
    user_id       BIGINT PRIMARY KEY, -- FK -> users.id
    display_name  VARCHAR(100), -- display name
    avatar_url    VARCHAR(500), -- avatar URL
    birth_day     DATE, -- birthday (YYYY-MM-DD)
    bio           TEXT, -- user bio
    created_at    TIMESTAMP DEFAULT CURRENT_TIMESTAMP, -- profile creation time
    updated_at    TIMESTAMP DEFAULT CURRENT_TIMESTAMP, -- profile last update time
    CONSTRAINT fk_up_user
        FOREIGN KEY (user_id) REFERENCES users(id)
);

CREATE TABLE user_statistics (
    user_id             BIGINT PRIMARY KEY, -- FK -> users.id
    total_sessions      INTEGER DEFAULT 0, -- total number of game sessions
    total_questions     INTEGER DEFAULT 0, -- total number of answered questions
    total_correct       INTEGER DEFAULT 0, -- total number of correct answers
    total_time_seconds  INTEGER DEFAULT 0, -- total play time (in seconds)
    last_played_at      TIMESTAMP, -- last play time
    CONSTRAINT fk_us_user
        FOREIGN KEY (user_id) REFERENCES users(id)
);

CREATE TABLE user_word_statistics (
    user_id          BIGINT NOT NULL, -- FK -> users.id
    word_id          BIGINT NOT NULL, -- FK -> words.id
    correct_count    INTEGER DEFAULT 0, -- number of times this word was answered correctly
    wrong_count      INTEGER DEFAULT 0, -- number of times this word was answered incorrectly
    last_answered_at TIMESTAMP, -- most recent time this word was answered
    streak           INTEGER DEFAULT 0, -- current correct streak for this word
    PRIMARY KEY (user_id, word_id),
    CONSTRAINT fk_uws_user
        FOREIGN KEY (user_id) REFERENCES users(id),
    CONSTRAINT fk_uws_word
        FOREIGN KEY (word_id) REFERENCES words(id)
);

CREATE TABLE user_topic_statistics (
    user_id         BIGINT NOT NULL, -- FK -> users.id
    topic_id        BIGINT NOT NULL, -- FK -> topics.id
    total_questions INTEGER DEFAULT 0, -- total questions for this topic
    total_correct   INTEGER DEFAULT 0, -- total correct answers for this topic
    last_played_at  TIMESTAMP, -- most recent time this topic was played
    PRIMARY KEY (user_id, topic_id),
    CONSTRAINT fk_uts_user
        FOREIGN KEY (user_id) REFERENCES users(id),
    CONSTRAINT fk_uts_topic
        FOREIGN KEY (topic_id) REFERENCES topics(id)
);

CREATE TABLE vocab_game_sessions (
    id                  BIGSERIAL PRIMARY KEY, -- game session id
    user_id             BIGINT NOT NULL, -- FK -> users.id
    mode                VARCHAR(50) NOT NULL, -- mode: 'level', 'topic', ...
    source_language_id  SMALLINT NOT NULL, -- FK -> languages.id (question language)
    target_language_id  SMALLINT NOT NULL, -- FK -> languages.id (answer language)
    topic_id            BIGINT, -- FK -> topics.id (if playing by topic)
    level_id            BIGINT, -- FK -> levels.id (if playing by level)
    total_questions     SMALLINT DEFAULT 0, -- total number of questions in the session
    correct_questions   SMALLINT DEFAULT 0, -- total number of correct answers
    started_at          TIMESTAMP DEFAULT CURRENT_TIMESTAMP, -- session start time
    ended_at            TIMESTAMP, -- session end time
    CONSTRAINT fk_vgs_user
        FOREIGN KEY (user_id) REFERENCES users(id),
    CONSTRAINT fk_vgs_source_lang
        FOREIGN KEY (source_language_id) REFERENCES languages(id),
    CONSTRAINT fk_vgs_target_lang
        FOREIGN KEY (target_language_id) REFERENCES languages(id),
    CONSTRAINT fk_vgs_topic
        FOREIGN KEY (topic_id) REFERENCES topics(id),
    CONSTRAINT fk_vgs_level
        FOREIGN KEY (level_id) REFERENCES levels(id)
);

CREATE INDEX idx_vgs_user_time ON vocab_game_sessions(user_id, started_at);

CREATE TABLE vocab_game_questions (
    id                     BIGSERIAL PRIMARY KEY, -- game question id
    session_id             BIGINT NOT NULL, -- FK -> vocab_game_sessions.id
    question_order         SMALLINT NOT NULL, -- question order within the session
    question_type          VARCHAR(30) NOT NULL, -- question type: 'word_to_translation', ...
    source_word_id         BIGINT NOT NULL, -- FK -> words.id (source word)
    source_sense_id        BIGINT, -- FK -> senses.id (specific sense, if used)
    correct_target_word_id BIGINT NOT NULL, -- FK -> words.id (correct answer)
    source_language_id     SMALLINT NOT NULL, -- FK -> languages.id (question language)
    target_language_id     SMALLINT NOT NULL, -- FK -> languages.id (answer language)
    created_at             TIMESTAMP DEFAULT CURRENT_TIMESTAMP, -- question creation time
    CONSTRAINT fk_vgq_session
        FOREIGN KEY (session_id) REFERENCES vocab_game_sessions(id),
    CONSTRAINT fk_vgq_source_word
        FOREIGN KEY (source_word_id) REFERENCES words(id),
    CONSTRAINT fk_vgq_source_sense
        FOREIGN KEY (source_sense_id) REFERENCES senses(id),
    CONSTRAINT fk_vgq_correct_word
        FOREIGN KEY (correct_target_word_id) REFERENCES words(id),
    CONSTRAINT fk_vgq_source_lang
        FOREIGN KEY (source_language_id) REFERENCES languages(id),
    CONSTRAINT fk_vgq_target_lang
        FOREIGN KEY (target_language_id) REFERENCES languages(id)
);

CREATE INDEX idx_vgq_session_order ON vocab_game_questions(session_id, question_order);

CREATE TABLE vocab_game_question_options (
    id             BIGSERIAL PRIMARY KEY, -- option id
    question_id    BIGINT NOT NULL, -- FK -> vocab_game_questions.id
    option_label   CHAR(1) NOT NULL, -- label: 'A', 'B', 'C', 'D'
    target_word_id BIGINT NOT NULL, -- FK -> words.id (word shown as an option)
    is_correct     BOOLEAN NOT NULL DEFAULT FALSE, -- TRUE if this is the correct answer
    CONSTRAINT fk_vgqo_question
        FOREIGN KEY (question_id) REFERENCES vocab_game_questions(id),
    CONSTRAINT fk_vgqo_target_word
        FOREIGN KEY (target_word_id) REFERENCES words(id),
    UNIQUE (question_id, option_label) -- mỗi câu chỉ có 1 A/B/C/D
);

CREATE TABLE vocab_game_question_answers (
    id                 BIGSERIAL PRIMARY KEY, -- answer id
    question_id        BIGINT NOT NULL, -- FK -> vocab_game_questions.id
    session_id         BIGINT NOT NULL, -- FK -> vocab_game_sessions.id
    user_id            BIGINT NOT NULL, -- FK -> users.id
    selected_option_id BIGINT, -- FK -> vocab_game_question_options.id (user's chosen answer)
    is_correct         BOOLEAN NOT NULL DEFAULT FALSE, -- TRUE if the answer is correct
    response_time_ms   INTEGER, -- response time (ms)
    answered_at        TIMESTAMP DEFAULT CURRENT_TIMESTAMP, -- answer time
    CONSTRAINT fk_vgqa_question
        FOREIGN KEY (question_id) REFERENCES vocab_game_questions(id),
    CONSTRAINT fk_vgqa_session
        FOREIGN KEY (session_id) REFERENCES vocab_game_sessions(id),
    CONSTRAINT fk_vgqa_user
        FOREIGN KEY (user_id) REFERENCES users(id),
    CONSTRAINT fk_vgqa_option
        FOREIGN KEY (selected_option_id) REFERENCES vocab_game_question_options(id)
);

CREATE INDEX idx_vgqa_user_time ON vocab_game_question_answers(user_id, answered_at);

-- Create function and trigger for updated_at columns
CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = CURRENT_TIMESTAMP;
    RETURN NEW;
END;
$$ language 'plpgsql';

CREATE TRIGGER update_words_updated_at BEFORE UPDATE ON words
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_users_updated_at BEFORE UPDATE ON users
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_user_profiles_updated_at BEFORE UPDATE ON user_profiles
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();
