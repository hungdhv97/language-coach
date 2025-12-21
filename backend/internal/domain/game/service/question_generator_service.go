package service

import (
	"context"
	"fmt"
	"math/rand"
	"time"

	dictdomain "github.com/english-coach/backend/internal/modules/dictionary/domain"
	"github.com/english-coach/backend/internal/modules/game/domain"
	"github.com/english-coach/backend/internal/shared/errors"
	"go.uber.org/zap"
)

// QuestionGeneratorService implements question generation logic
type QuestionGeneratorService struct {
	wordRepo dictdomain.WordRepository
	logger   *zap.Logger
}

// NewQuestionGeneratorService creates a new question generator service
func NewQuestionGeneratorService(
	wordRepo dictdomain.WordRepository,
	logger *zap.Logger,
) *QuestionGeneratorService {
	return &QuestionGeneratorService{
		wordRepo: wordRepo,
		logger:   logger,
	}
}

// GenerateQuestions generates questions for a game session
func (s *QuestionGeneratorService) GenerateQuestions(
	ctx context.Context,
	sessionID int64,
	sourceLanguageID, targetLanguageID int16,
	mode string,
	topicIDs []int64,
	levelID int64,
	questionCount int,
) ([]*domain.GameQuestion, []*domain.GameQuestionOption, error) {
	startTime := time.Now()

	// Validate mode
	if mode != "level" {
		return nil, nil, domain.ErrInvalidMode.WithDetails(fmt.Sprintf("mode: %s, required: 'level'", mode))
	}

	// Fetch source words by level and optional topics
	// Fetch up to questionCount*3 words to have options for wrong answers
	maxWordsToFetch := questionCount * 3
	if maxWordsToFetch > 60 { // Cap at 60 (20*3) to avoid excessive queries
		maxWordsToFetch = 60
	}

	sourceWords, err := s.wordRepo.FindWordsByLevelAndTopicsAndLanguages(
		ctx, levelID, topicIDs, sourceLanguageID, targetLanguageID, maxWordsToFetch,
	)
	if err != nil {
		s.logger.Error("failed to fetch source words by level and topics",
			zap.Error(err),
			zap.String("mode", mode),
			zap.Int64("level_id", levelID),
			zap.Any("topic_ids", topicIDs),
			zap.Int16("source_language_id", sourceLanguageID),
			zap.Int16("target_language_id", targetLanguageID),
			zap.Int("requested_limit", maxWordsToFetch),
		)
		return nil, nil, errors.WrapError(err, "failed to fetch source words")
	}
	s.logger.Info("fetched words by level and topics",
		zap.Int64("level_id", levelID),
		zap.Any("topic_ids", topicIDs),
		zap.Int16("source_language_id", sourceLanguageID),
		zap.Int16("target_language_id", targetLanguageID),
		zap.Int("word_count", len(sourceWords)),
		zap.Int("requested_limit", maxWordsToFetch),
	)

	// Check if we have at least 1 word (minimum required)
	if len(sourceWords) < 1 {
		s.logger.Warn("no words available for question generation",
			zap.String("mode", mode),
			zap.Int("requested", questionCount),
			zap.Int("available", len(sourceWords)),
			zap.Any("topic_ids", topicIDs),
			zap.Int64("level_id", levelID),
			zap.Int16("source_language_id", sourceLanguageID),
			zap.Int16("target_language_id", targetLanguageID),
		)
		return nil, nil, domain.ErrInsufficientWords.WithDetails(fmt.Sprintf("required: 1, available: %d", len(sourceWords)))
	}

	// Shuffle words for randomness
	rand.Shuffle(len(sourceWords), func(i, j int) {
		sourceWords[i], sourceWords[j] = sourceWords[j], sourceWords[i]
	})

	// Select up to questionCount words (or all available if fewer)
	// This allows games with 1-20 words depending on available data
	wordsToSelect := questionCount
	if len(sourceWords) < questionCount {
		wordsToSelect = len(sourceWords)
		s.logger.Info("using fewer words than requested",
			zap.Int("requested", questionCount),
			zap.Int("available", len(sourceWords)),
			zap.Int("using", wordsToSelect),
		)
	}
	selectedWords := sourceWords[:wordsToSelect]

	// Generate questions
	questions := make([]*domain.GameQuestion, 0, questionCount)
	options := make([]*domain.GameQuestionOption, 0, questionCount*4)

	// Collect all target words for wrong answer options
	allTargetWords := make(map[int64]*dictdomain.Word)

	// Track question order separately to ensure sequential ordering even when words are skipped
	questionOrder := int16(0)

	for _, sourceWord := range selectedWords {
		// Get correct translation
		translations, err := s.wordRepo.FindTranslationsForWord(
			ctx, sourceWord.ID, targetLanguageID, 10,
		)
		if err != nil || len(translations) == 0 {
			s.logger.Warn("no translations found for word",
				zap.Int64("word_id", sourceWord.ID),
				zap.Int16("target_language_id", targetLanguageID),
			)
			continue
		}

		correctWord := translations[0] // Use first translation as correct answer
		allTargetWords[correctWord.ID] = correctWord

		// Collect other translations for wrong answers
		for _, trans := range translations[1:] {
			allTargetWords[trans.ID] = trans
		}

		// Increment question order for each successfully created question
		questionOrder++

		// Create question
		question := &domain.GameQuestion{
			SessionID:           sessionID,
			QuestionOrder:       questionOrder,
			QuestionType:        "word_to_translation",
			SourceWordID:        sourceWord.ID,
			CorrectTargetWordID: correctWord.ID,
			SourceLanguageID:    sourceLanguageID,
			TargetLanguageID:    targetLanguageID,
			CreatedAt:           time.Now(),
		}
		questions = append(questions, question)
	}

	// Generate options for each question
	targetWordList := make([]*dictdomain.Word, 0, len(allTargetWords))
	for _, word := range allTargetWords {
		targetWordList = append(targetWordList, word)
	}

	for i, question := range questions {
		// Get correct word
		correctWord, exists := allTargetWords[question.CorrectTargetWordID]
		if !exists {
			return nil, nil, domain.ErrQuestionNotFound.WithDetails(fmt.Sprintf("question_index: %d", i+1))
		}

		// Get wrong answer candidates (exclude correct answer)
		wrongCandidates := make([]*dictdomain.Word, 0)
		for _, word := range targetWordList {
			if word.ID != correctWord.ID {
				wrongCandidates = append(wrongCandidates, word)
			}
		}

		// Ensure we have at least 3 wrong candidates
		if len(wrongCandidates) < 3 {
			// If not enough, we need to fetch more translations
			// For now, we'll use what we have and pad with duplicates if needed
			for len(wrongCandidates) < 3 {
				// Try to get more translations for other source words
				// This is a simplified approach - in production, you might want to fetch more
				if len(wrongCandidates) > 0 {
					wrongCandidates = append(wrongCandidates, wrongCandidates[0])
				} else {
					return nil, nil, domain.ErrOptionNotFound.WithDetails(fmt.Sprintf("question_index: %d", i+1))
				}
			}
		}

		// Shuffle wrong candidates
		rand.Shuffle(len(wrongCandidates), func(i, j int) {
			wrongCandidates[i], wrongCandidates[j] = wrongCandidates[j], wrongCandidates[i]
		})

		// Select 3 wrong answers
		selectedWrong := wrongCandidates[:3]

		// Combine correct + wrong answers and shuffle
		allAnswers := []*dictdomain.Word{correctWord, selectedWrong[0], selectedWrong[1], selectedWrong[2]}
		rand.Shuffle(len(allAnswers), func(i, j int) {
			allAnswers[i], allAnswers[j] = allAnswers[j], allAnswers[i]
		})

		// Find correct answer position
		correctIndex := -1
		for idx, word := range allAnswers {
			if word.ID == correctWord.ID {
				correctIndex = idx
				break
			}
		}

		if correctIndex == -1 {
			return nil, nil, domain.ErrQuestionNotFound.WithDetails("correct answer not found in shuffled options")
		}

		// Create options (A, B, C, D)
		labels := []string{"A", "B", "C", "D"}
		for j, word := range allAnswers {
			option := &domain.GameQuestionOption{
				QuestionID:   question.ID, // Will be set after question is saved
				OptionLabel:  labels[j],
				TargetWordID: word.ID,
				IsCorrect:    j == correctIndex,
			}
			options = append(options, option)
		}
	}

	duration := time.Since(startTime)
	s.logger.Info("questions generated",
		zap.Int64("session_id", sessionID),
		zap.Int("question_count", len(questions)),
		zap.Duration("duration", duration),
	)

	// Ensure generation completes within 1 second (SC-003)
	if duration > time.Second {
		s.logger.Warn("question generation took longer than 1 second",
			zap.Duration("duration", duration),
		)
	}

	return questions, options, nil
}
