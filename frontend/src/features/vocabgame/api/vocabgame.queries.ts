/**
 * VocabGame queries using React Query
 */

import { useQuery } from '@tanstack/react-query';
import { vocabGameEndpoints } from '@/entities/vocabgame/api/vocabgame.endpoints';
import type {
  VocabGameSessionWithQuestions,
  SessionStatistics,
} from '@/entities/vocabgame/model/vocabgame.types';

export const vocabGameQueries = {
  /**
   * Query key factory for vocabgame queries
   */
  keys: {
    all: ['vocabgame'] as const,
    sessions: () => [...vocabGameQueries.keys.all, 'sessions'] as const,
    session: (sessionId: number) =>
      [...vocabGameQueries.keys.sessions(), sessionId] as const,
    statistics: () => [...vocabGameQueries.keys.all, 'statistics'] as const,
    sessionStatistics: (sessionId: number) =>
      [...vocabGameQueries.keys.statistics(), sessionId] as const,
  },

  /**
   * Get a vocabgame session with questions
   */
  useSession: (sessionId: number) => {
    return useQuery<VocabGameSessionWithQuestions>({
      queryKey: vocabGameQueries.keys.session(sessionId),
      queryFn: () => vocabGameEndpoints.getSession(sessionId),
      enabled: !!sessionId && sessionId > 0,
      staleTime: 0, // Always fetch fresh data for active vocabgame
    });
  },

  /**
   * Get session statistics
   */
  useSessionStatistics: (sessionId: number) => {
    return useQuery<SessionStatistics>({
      queryKey: vocabGameQueries.keys.sessionStatistics(sessionId),
      queryFn: () => vocabGameEndpoints.getSessionStatistics(sessionId),
      enabled: !!sessionId && sessionId > 0,
    });
  },
};

