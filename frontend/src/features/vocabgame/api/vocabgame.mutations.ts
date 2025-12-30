/**
 * VocabGame mutations using React Query
 */

import { useMutation, useQueryClient } from '@tanstack/react-query';
import { vocabGameEndpoints } from '@/entities/vocabgame/api/vocabgame.endpoints';
import type {
  CreateVocabGameSessionRequest,
  VocabGameSession,
  VocabGameAnswer,
  SubmitAnswerRequest,
} from '@/entities/vocabgame/model/vocabgame.types';

export const vocabGameMutations = {
  /**
   * Create a new vocabgame session
   */
  useCreateSession: () => {
    const queryClient = useQueryClient();

    return useMutation<VocabGameSession, Error, CreateVocabGameSessionRequest>({
      mutationFn: vocabGameEndpoints.createSession,
      onSuccess: () => {
        // Invalidate vocabgame session queries if needed
        queryClient.invalidateQueries({ queryKey: ['vocabgame', 'sessions'] });
      },
    });
  },

  /**
   * Submit an answer to a question
   */
  useSubmitAnswer: (sessionId: number) => {
    const queryClient = useQueryClient();

    return useMutation<VocabGameAnswer, Error, SubmitAnswerRequest>({
      mutationFn: (request) => vocabGameEndpoints.submitAnswer(sessionId, request),
      onSuccess: () => {
        // Invalidate session query to refresh data
        queryClient.invalidateQueries({
          queryKey: ['vocabgame', 'sessions', sessionId],
        });
      },
    });
  },
};

