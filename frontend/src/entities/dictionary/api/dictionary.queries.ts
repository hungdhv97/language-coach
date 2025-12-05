/**
 * Dictionary React Query hooks
 */

import { useQuery } from '@tanstack/react-query';
import { dictionaryEndpoints } from './dictionary.endpoints';
import type { Language, Topic, Level, WordDetail, WordSearchResponse } from '../model/dictionary.types';

export const dictionaryQueries = {
  /**
   * Query key factory for dictionary queries
   */
  keys: {
    all: ['dictionary'] as const,
    languages: () => [...dictionaryQueries.keys.all, 'languages'] as const,
    topics: () => [...dictionaryQueries.keys.all, 'topics'] as const,
    levels: (languageId?: number) =>
      [...dictionaryQueries.keys.all, 'levels', languageId] as const,
    search: (query: string, languageId: number, limit?: number, offset?: number) =>
      [...dictionaryQueries.keys.all, 'search', query, languageId, limit, offset] as const,
    wordDetail: (wordId: number) =>
      [...dictionaryQueries.keys.all, 'word', wordId] as const,
  },

  /**
   * Get all languages
   */
  useLanguages: () => {
    return useQuery<Language[]>({
      queryKey: dictionaryQueries.keys.languages(),
      queryFn: dictionaryEndpoints.getLanguages,
    });
  },

  /**
   * Get all topics
   */
  useTopics: () => {
    return useQuery<Topic[]>({
      queryKey: dictionaryQueries.keys.topics(),
      queryFn: dictionaryEndpoints.getTopics,
    });
  },

  /**
   * Get levels, optionally filtered by language ID
   */
  useLevels: (languageId?: number) => {
    return useQuery<Level[]>({
      queryKey: dictionaryQueries.keys.levels(languageId),
      queryFn: () => dictionaryEndpoints.getLevels(languageId),
      enabled: true, // Always enabled, languageId is optional
    });
  },

  /**
   * Search words
   */
  useSearchWords: (
    query: string,
    languageId: number,
    limit: number = 20,
    offset: number = 0,
    enabled: boolean = true
  ) => {
    return useQuery<WordSearchResponse>({
      queryKey: dictionaryQueries.keys.search(query, languageId, limit, offset),
      queryFn: () => dictionaryEndpoints.searchWords(query, languageId, limit, offset),
      enabled: enabled && query.trim().length > 0 && !!languageId,
      staleTime: 30 * 1000, // 30 seconds
    });
  },

  /**
   * Get word detail by ID
   */
  useWordDetail: (wordId: number) => {
    return useQuery<WordDetail>({
      queryKey: dictionaryQueries.keys.wordDetail(wordId),
      queryFn: () => dictionaryEndpoints.getWordDetail(wordId),
      enabled: !!wordId && wordId > 0,
    });
  },
};

