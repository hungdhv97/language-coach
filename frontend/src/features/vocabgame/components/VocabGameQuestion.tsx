/**
 * VocabGame Question Component
 * Displays a question with 4 multiple-choice options (A, B, C, D)
 */

import { useState, useMemo } from 'react';
import type { VocabGameQuestionWithOptions } from '@/entities/vocabgame/model/vocabgame.types';
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card';
import { Button } from '@/components/ui/button';
import { cn } from '@/lib/utils';

interface VocabGameQuestionProps {
  question: VocabGameQuestionWithOptions;
  onAnswerSelect: (optionId: number) => void;
  isSubmitting?: boolean;
  selectedOptionId?: number;
  totalQuestions?: number;
}

export default function VocabGameQuestion({
  question,
  onAnswerSelect,
  isSubmitting = false,
  selectedOptionId,
}: VocabGameQuestionProps) {
  // Use a key-based approach: derive localSelected from props when question changes
  // This avoids setState in useEffect
  const questionKey = `${question.id}-${selectedOptionId || 'none'}`;
  const [localSelectedMap, setLocalSelectedMap] = useState<Record<string, number | null>>({});
  
  // Get current local selected for this question
  const localSelected = useMemo(() => {
    return localSelectedMap[questionKey] ?? selectedOptionId ?? null;
  }, [localSelectedMap, questionKey, selectedOptionId]);

  const handleOptionClick = (optionId: number) => {
    if (isSubmitting || localSelected !== null) {
      return; // Prevent multiple selections
    }
    setLocalSelectedMap(prev => ({ ...prev, [questionKey]: optionId }));
    onAnswerSelect(optionId);
  };

  // Sort options by label (A, B, C, D)
  const sortedOptions = [...question.options].sort((a, b) =>
    a.option_label.localeCompare(b.option_label)
  );

  return (
    <Card>
      <CardHeader>
        <CardTitle className="text-center text-2xl">{question.source_word_text}</CardTitle>
        <p className="text-center text-muted-foreground">
          Chọn từ đúng trong ngôn ngữ đích
        </p>
      </CardHeader>
      <CardContent className="space-y-3">
        {sortedOptions.map((option) => {
          const isSelected = localSelected === option.id;
          const isDisabled = isSubmitting || localSelected !== null;

          return (
            <Button
              key={option.id}
              variant={isSelected ? 'default' : 'outline'}
              className={cn(
                'w-full h-auto py-6 justify-start text-left border',
                isSelected && 'bg-primary text-primary-foreground border-primary',
                !isSelected && 'border-input',
                isDisabled && 'cursor-not-allowed opacity-50'
              )}
              onClick={() => handleOptionClick(option.id)}
              disabled={isDisabled}
              aria-label={`Option ${option.option_label}: ${option.word_text}`}
            >
              <span className="font-semibold mr-3 min-w-[2rem] text-center">
                {option.option_label}.
              </span>
              <span>{option.word_text}</span>
            </Button>
          );
        })}
      </CardContent>
    </Card>
  );
}

