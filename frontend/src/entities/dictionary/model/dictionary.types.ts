/**
 * Dictionary entity types
 */

export interface Language {
  id: number;
  code: string;
  name: string;
  native_name?: string;
}

export interface Topic {
  id: number;
  code: string;
  name: string;
}

export interface Level {
  id: number;
  code: string;
  name: string;
  description?: string;
  language_id?: number;
  difficulty_order?: number;
}

export interface Word {
  id: number;
  language_id: number;
  lemma: string;
  lemma_normalized?: string;
  search_key?: string;
  romanization?: string;
  script_code?: string;
  frequency_rank?: number;
  note?: string;
  topics?: string[];
  created_at: string;
  updated_at: string;
}

export interface WordRelation {
  relation_type: 'synonym' | 'antonym' | 'related';
  note?: string;
  target_word: Word;
}

export interface Sense {
  id: number;
  word_id: number;
  sense_order: number;
  part_of_speech_id: number;
  definition: string;
  definition_language_id: number;
  usage_label?: string;
  level_id?: number;
  note?: string;
}

export interface ExampleTranslation {
  language: string;
  content: string;
}

export interface Example {
  id: number;
  source_sense_id: number;
  language_id: number;
  content: string;
  audio_url?: string;
  source?: string;
  translations?: ExampleTranslation[];
}

export interface Pronunciation {
  id: number;
  word_id: number;
  dialect?: string;
  ipa?: string;
  phonetic?: string;
  audio_url?: string;
}

export interface SenseDetail {
  id: number;
  sense_order: number;
  part_of_speech_id: number;
  part_of_speech_name?: string;
  definition: string;
  definition_language_id: number;
  definition_language_name?: string;
  usage_label?: string;
  level_id?: number;
  level_name?: string;
  note?: string;
  translations: Word[];
  examples: Example[];
}

export interface WordDetail {
  word: Word;
  senses: SenseDetail[];
  pronunciations: Pronunciation[];
  relations?: WordRelation[];
}

export interface WordSearchResponse {
  words: Word[];
  total: number;
  limit: number;
  offset: number;
}

