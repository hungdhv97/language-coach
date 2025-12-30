/**
 * Hook for managing vocabgame session state
 */

import { useState, useCallback, useRef } from 'react';
import { vocabGameMutations } from '../api/vocabgame.mutations';
import type { SubmitAnswerRequest, VocabGameAnswer } from '@/entities/vocabgame/model/vocabgame.types';

export interface UseVocabGameSessionOptions {
  sessionId: number;
  onAnswerSubmitted?: (answer: VocabGameAnswer) => void;
  onAllQuestionsAnswered?: () => void;
}

export function useVocabGameSession({ sessionId, onAnswerSubmitted }: UseVocabGameSessionOptions) {
  const [currentQuestionIndex, setCurrentQuestionIndex] = useState(0);
  const [answers, setAnswers] = useState<Map<number, VocabGameAnswer>>(new Map());
  const [responseStartTime, setResponseStartTime] = useState<number | null>(null);
  const questionStartTimeRef = useRef<number | null>(null);

  const submitAnswerMutation = vocabGameMutations.useSubmitAnswer(sessionId);

  // Start tracking response time when question is displayed
  const startQuestion = useCallback(() => {
    questionStartTimeRef.current = Date.now();
    setResponseStartTime(Date.now());
  }, []);

  // Submit answer
  const submitAnswer = useCallback(
    async (request: SubmitAnswerRequest) => {
      if (!questionStartTimeRef.current) {
        return;
      }

      // Calculate response time
      const responseTime = Date.now() - questionStartTimeRef.current;
      const requestWithTime: SubmitAnswerRequest = {
        ...request,
        response_time_ms: responseTime,
      };

      try {
        const answer = await submitAnswerMutation.mutateAsync(requestWithTime);
        setAnswers((prev) => {
          const newMap = new Map(prev);
          newMap.set(request.question_id, answer);
          return newMap;
        });

        onAnswerSubmitted?.(answer);
        return answer;
      } catch (error) {
        // Enhanced error handling for network errors
        if (error instanceof Error) {
          // Network error (connection failed, timeout, etc.)
          if (error.message.includes('fetch') || error.message.includes('network') || error.message.includes('Failed to fetch')) {
            throw new Error('Không thể kết nối đến máy chủ. Vui lòng kiểm tra kết nối mạng và thử lại.');
          }
          // API error (from error interceptor)
          if ('code' in error && 'message' in error) {
            throw error; // Already formatted by error interceptor
          }
        }
        // Generic error
        console.error('Failed to submit answer:', error);
        throw new Error('Không thể gửi câu trả lời. Vui lòng thử lại.');
      }
    },
    [submitAnswerMutation, onAnswerSubmitted]
  );

  // Move to next question
  const nextQuestion = useCallback(() => {
    setCurrentQuestionIndex((prev) => prev + 1);
    questionStartTimeRef.current = Date.now();
    setResponseStartTime(Date.now());
  }, []);

  // Reset to first question
  const reset = useCallback(() => {
    setCurrentQuestionIndex(0);
    setAnswers(new Map());
    questionStartTimeRef.current = null;
    setResponseStartTime(null);
  }, []);

  return {
    currentQuestionIndex,
    answers,
    responseStartTime,
    submitAnswer,
    nextQuestion,
    reset,
    startQuestion,
    isSubmitting: submitAnswerMutation.isPending,
  };
}

