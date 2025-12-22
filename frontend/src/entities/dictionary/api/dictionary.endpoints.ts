/**
 * Dictionary API endpoints
 */

import { httpClient } from '@/shared/api/http-client';
import type { Language, Topic, Level, Word, WordDetail, WordSearchResponse } from '../model/dictionary.types';

export interface ApiResponse<T> {
  success: boolean;
  data: T;
}

export interface PaginatedApiResponse<T> {
  success: boolean;
  data: T;
  pagination: {
    page: number;
    pageSize: number;
    total: number;
    totalPages: number;
    limit: number;
    offset: number;
    hasNext: boolean;
    hasPrev: boolean;
  };
}

export const dictionaryEndpoints = {
  /**
   * Get all languages
   */
  getLanguages: async (): Promise<Language[]> => {
    const response = await httpClient.get<ApiResponse<Language[]>>('/reference/languages');
    return response.data || [];
  },

  /**
   * Get all topics
   */
  getTopics: async (): Promise<Topic[]> => {
    const response = await httpClient.get<ApiResponse<Topic[]>>('/reference/topics');
    return response.data || [];
  },

  /**
   * Get all levels, optionally filtered by language ID
   */
  getLevels: async (languageId?: number): Promise<Level[]> => {
    const url = languageId
      ? `/reference/levels?languageId=${languageId}`
      : '/reference/levels';
    const response = await httpClient.get<ApiResponse<Level[]>>(url);
    return response.data || [];
  },

  /**
   * Search words
   */
  searchWords: async (
    query: string,
    languageId: number,
    limit: number = 20,
    offset: number = 0
  ): Promise<WordSearchResponse> => {
    const params = new URLSearchParams({
      q: query,
      languageId: languageId.toString(),
      limit: limit.toString(),
      offset: offset.toString(),
    });
    const response = await httpClient.get<PaginatedApiResponse<Word[]>>(
      `/dictionary/search?${params.toString()}`
    );
    // Transform response to match WordSearchResponse interface
    return {
      words: response.data || [],
      pagination: response.pagination,
    };
  },

  /**
   * Get word detail by ID
   */
  getWordDetail: async (wordId: number): Promise<WordDetail> => {
    const response = await httpClient.get<ApiResponse<WordDetail>>(
      `/dictionary/words/${wordId}`
    );
    return response.data;
  },
};

