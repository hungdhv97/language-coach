package create_session

import (
	"context"
	"math/rand"
	"time"

	dictdomain "github.com/english-coach/backend/internal/modules/dictionary/domain"
	"github.com/english-coach/backend/internal/modules/game/domain"
	"github.com/english-coach/backend/internal/shared/constants"
	sharederrors "github.com/english-coach/backend/internal/shared/errors"
	"github.com/english-coach/backend/internal/shared/logger"
)

// Handler handles game session creation
type Handler struct {
	sessionRepo  domain.GameSessionRepository
	questionRepo domain.GameQuestionRepository
	wordRepo     dictdomain.WordRepository
	logger       logger.ILogger
}

// NewHandler creates a new use case
func NewHandler(
	sessionRepo domain.GameSessionRepository,
	questionRepo domain.GameQuestionRepository,
	wordRepo dictdomain.WordRepository,
	logger logger.ILogger,
) *Handler {
	return &Handler{
		sessionRepo:  sessionRepo,
		questionRepo: questionRepo,
		wordRepo:     wordRepo,
		logger:       logger,
	}
}

// Execute creates a new game session
func (h *Handler) Execute(ctx context.Context, input CreateSessionInput, userID int64) (*CreateSessionOutput, error) {
	// Validate request
	if err := input.Validate(); err != nil {
		return nil, sharederrors.ErrValidationError.WithDetails(err.Error())
	}

	// Create game session model
	// Note: TopicID is kept for backward compatibility with DB schema, but we use TopicIDs array for filtering
	var topicID *int64
	if len(input.TopicIDs) > 0 {
		// Store first topic ID for DB compatibility (schema still has single topic_id)
		topicID = &input.TopicIDs[0]
	}
	levelID := &input.LevelID

	session := &domain.GameSession{
		UserID:           userID,
		Mode:             input.Mode,
		SourceLanguageID: input.SourceLanguageID,
		TargetLanguageID: input.TargetLanguageID,
		TopicID:          topicID,
		LevelID:          levelID,
		TotalQuestions:   0, // Will be set when questions are generated
		CorrectQuestions: 0,
		StartedAt:        time.Now(),
	}

	// Save session to database first (needed for question generation)
	if err := h.sessionRepo.Create(ctx, session); err != nil {
		h.logger.Error("failed to create game session",
			logger.Error(err),
			logger.Int64("user_id", userID),
			logger.String("mode", input.Mode),
		)
		return nil, sharederrors.MapDomainErrorToAppError(err)
	}

	// Generate questions upfront - request up to MaxGameQuestionCount (20)
	questions, options, err := h.generateQuestions(
		ctx,
		session.ID,
		input.SourceLanguageID,
		input.TargetLanguageID,
		input.Mode,
		input.TopicIDs,
		input.LevelID,
		constants.MaxGameQuestionCount,
	)
	if err != nil {
		h.logger.Error("failed to generate questions",
			logger.Error(err),
			logger.Int64("session_id", session.ID),
			logger.String("mode", input.Mode),
			logger.Int("source_language_id", int(input.SourceLanguageID)),
			logger.Int("target_language_id", int(input.TargetLanguageID)),
			logger.Any("topic_ids", input.TopicIDs),
			logger.Any("level_id", input.LevelID),
		)
		return nil, sharederrors.MapDomainErrorToAppError(err)
	}

	// Check if we have at least the minimum required questions (1)
	if len(questions) < constants.MinGameQuestionCount {
		return nil, sharederrors.MapDomainErrorToAppError(domain.ErrInsufficientWords)
	}

	// Save questions and options
	if err := h.questionRepo.CreateBatch(ctx, questions, options); err != nil {
		h.logger.Error("failed to save questions",
			logger.Error(err),
			logger.Int64("session_id", session.ID),
		)
		return nil, sharederrors.MapDomainErrorToAppError(err)
	}

	// Update session with question count
	session.TotalQuestions = int16(len(questions))
	if err := h.sessionRepo.Update(ctx, session); err != nil {
		h.logger.Error("failed to update session with question count",
			logger.Error(err),
			logger.Int64("session_id", session.ID),
		)
		return nil, sharederrors.MapDomainErrorToAppError(err)
	}

	// Log session creation
	h.logger.Info("game session created with questions",
		logger.Int64("session_id", session.ID),
		logger.Int64("user_id", userID),
		logger.String("mode", input.Mode),
		logger.Int("source_language_id", int(input.SourceLanguageID)),
		logger.Int("target_language_id", int(input.TargetLanguageID)),
		logger.Int("question_count", len(questions)),
	)

	return &CreateSessionOutput{
		ID:               session.ID,
		UserID:           session.UserID,
		Mode:             session.Mode,
		SourceLanguageID: session.SourceLanguageID,
		TargetLanguageID: session.TargetLanguageID,
		TopicID:          session.TopicID,
		LevelID:          session.LevelID,
		TotalQuestions:   session.TotalQuestions,
		CorrectQuestions: session.CorrectQuestions,
		StartedAt:        session.StartedAt,
		EndedAt:          session.EndedAt,
	}, nil
}

// generateQuestions generates questions for a game session
// This method encapsulates the question generation logic
func (h *Handler) generateQuestions(
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
	if err := h.validateMode(mode); err != nil {
		return nil, nil, err
	}

	// Fetch source words
	sourceWords, err := h.fetchSourceWords(ctx, levelID, topicIDs, sourceLanguageID, targetLanguageID, questionCount)
	if err != nil {
		return nil, nil, err
	}

	// Select and shuffle words
	selectedWords := h.selectAndShuffleWords(sourceWords, questionCount)

	// Build questions and collect target words
	questions, allTargetWords, err := h.buildQuestions(ctx, sessionID, selectedWords, sourceLanguageID, targetLanguageID)
	if err != nil {
		return nil, nil, err
	}
	if len(questions) == 0 {
		return nil, nil, domain.ErrInsufficientWords
	}

	// Generate options for each question
	options, err := h.generateOptions(questions, allTargetWords)
	if err != nil {
		return nil, nil, err
	}

	// Log generation performance
	h.logGenerationPerformance(startTime, sessionID, len(questions))

	return questions, options, nil
}

// validateMode validates the game mode
func (h *Handler) validateMode(mode string) error {
	if mode != "level" {
		return domain.ErrInvalidMode
	}
	return nil
}

// fetchSourceWords fetches source words from the repository
func (h *Handler) fetchSourceWords(
	ctx context.Context,
	levelID int64,
	topicIDs []int64,
	sourceLanguageID, targetLanguageID int16,
	questionCount int,
) ([]*dictdomain.Word, error) {
	// Fetch up to questionCount*3 words to have options for wrong answers
	maxWordsToFetch := questionCount * 3
	if maxWordsToFetch > 60 { // Cap at 60 (20*3) to avoid excessive queries
		maxWordsToFetch = 60
	}

	sourceWords, err := h.wordRepo.FindWordsByLevelAndTopicsAndLanguages(
		ctx, levelID, topicIDs, sourceLanguageID, targetLanguageID, maxWordsToFetch,
	)
	if err != nil {
		h.logger.Error("failed to fetch source words by level and topics",
			logger.Error(err),
			logger.String("mode", "level"),
			logger.Int64("level_id", levelID),
			logger.Any("topic_ids", topicIDs),
			logger.Int("source_language_id", int(sourceLanguageID)),
			logger.Int("target_language_id", int(targetLanguageID)),
			logger.Int("requested_limit", maxWordsToFetch),
		)
		return nil, err
	}

	h.logger.Info("fetched words by level and topics",
		logger.Int64("level_id", levelID),
		logger.Any("topic_ids", topicIDs),
		logger.Int("source_language_id", int(sourceLanguageID)),
		logger.Int("target_language_id", int(targetLanguageID)),
		logger.Int("word_count", len(sourceWords)),
		logger.Int("requested_limit", maxWordsToFetch),
	)

	// Check if we have at least 1 word (minimum required)
	if len(sourceWords) < 1 {
		h.logger.Warn("no words available for question generation",
			logger.String("mode", "level"),
			logger.Int("requested", questionCount),
			logger.Int("available", len(sourceWords)),
			logger.Any("topic_ids", topicIDs),
			logger.Int64("level_id", levelID),
			logger.Int("source_language_id", int(sourceLanguageID)),
			logger.Int("target_language_id", int(targetLanguageID)),
		)
		return nil, domain.ErrInsufficientWords
	}

	return sourceWords, nil
}

// selectAndShuffleWords selects and shuffles words for randomness
func (h *Handler) selectAndShuffleWords(sourceWords []*dictdomain.Word, questionCount int) []*dictdomain.Word {
	// Shuffle words for randomness
	rand.Shuffle(len(sourceWords), func(i, j int) {
		sourceWords[i], sourceWords[j] = sourceWords[j], sourceWords[i]
	})

	// Select up to questionCount words (or all available if fewer)
	wordsToSelect := questionCount
	if len(sourceWords) < questionCount {
		wordsToSelect = len(sourceWords)
		h.logger.Info("using fewer words than requested",
			logger.Int("requested", questionCount),
			logger.Int("available", len(sourceWords)),
			logger.Int("using", wordsToSelect),
		)
	}

	return sourceWords[:wordsToSelect]
}

// buildQuestions builds questions from selected words and collects target words
func (h *Handler) buildQuestions(
	ctx context.Context,
	sessionID int64,
	selectedWords []*dictdomain.Word,
	sourceLanguageID, targetLanguageID int16,
) ([]*domain.GameQuestion, map[int64]*dictdomain.Word, error) {
	questions := make([]*domain.GameQuestion, 0, len(selectedWords))
	allTargetWords := make(map[int64]*dictdomain.Word)
	questionOrder := int16(0)

	for _, sourceWord := range selectedWords {
		// Get correct translation
		translations, err := h.wordRepo.FindTranslationsForWord(
			ctx, sourceWord.ID, targetLanguageID, 10,
		)
		if err != nil {
			h.logger.Error("failed to find translations for word",
				logger.Error(err),
				logger.Int64("word_id", sourceWord.ID),
				logger.Int("target_language_id", int(targetLanguageID)),
			)
			return nil, nil, err
		}
		if len(translations) == 0 {
			h.logger.Warn("no translations found for word",
				logger.Int64("word_id", sourceWord.ID),
				logger.Int("target_language_id", int(targetLanguageID)),
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

	return questions, allTargetWords, nil
}

// generateOptions generates options (A, B, C, D) for each question
func (h *Handler) generateOptions(
	questions []*domain.GameQuestion,
	allTargetWords map[int64]*dictdomain.Word,
) ([]*domain.GameQuestionOption, error) {
	options := make([]*domain.GameQuestionOption, 0, len(questions)*4)

	// Convert map to slice for easier iteration
	targetWordList := make([]*dictdomain.Word, 0, len(allTargetWords))
	for _, word := range allTargetWords {
		targetWordList = append(targetWordList, word)
	}

	for _, question := range questions {
		// Get correct word
		correctWord, exists := allTargetWords[question.CorrectTargetWordID]
		if !exists {
			return nil, domain.ErrQuestionNotFound
		}

		// Get wrong answer candidates
		wrongCandidates := h.getWrongAnswerCandidates(targetWordList, correctWord.ID)

		// Ensure we have at least 3 wrong candidates
		if len(wrongCandidates) < 3 {
			wrongCandidates = h.padWrongCandidates(wrongCandidates)
			if len(wrongCandidates) < 3 {
				return nil, domain.ErrOptionNotFound
			}
		}

		// Create options for this question
		questionOptions := h.createQuestionOptions(question, correctWord, wrongCandidates)
		options = append(options, questionOptions...)
	}

	return options, nil
}

// getWrongAnswerCandidates gets wrong answer candidates excluding the correct answer
func (h *Handler) getWrongAnswerCandidates(targetWordList []*dictdomain.Word, correctWordID int64) []*dictdomain.Word {
	wrongCandidates := make([]*dictdomain.Word, 0)
	for _, word := range targetWordList {
		if word.ID != correctWordID {
			wrongCandidates = append(wrongCandidates, word)
		}
	}
	return wrongCandidates
}

// padWrongCandidates pads wrong candidates if not enough are available
func (h *Handler) padWrongCandidates(wrongCandidates []*dictdomain.Word) []*dictdomain.Word {
	// If not enough, pad with duplicates if needed
	for len(wrongCandidates) < 3 {
		if len(wrongCandidates) > 0 {
			wrongCandidates = append(wrongCandidates, wrongCandidates[0])
		} else {
			break
		}
	}
	return wrongCandidates
}

// createQuestionOptions creates options (A, B, C, D) for a question
func (h *Handler) createQuestionOptions(
	question *domain.GameQuestion,
	correctWord *dictdomain.Word,
	wrongCandidates []*dictdomain.Word,
) []*domain.GameQuestionOption {
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
	correctIndex := h.findCorrectAnswerIndex(allAnswers, correctWord.ID)
	if correctIndex == -1 {
		// This should not happen, but handle gracefully
		h.logger.Error("correct answer not found in shuffled options",
			logger.Int64("question_id", question.ID),
			logger.Int64("correct_word_id", correctWord.ID),
		)
		correctIndex = 0 // Fallback to first option
	}

	// Create options (A, B, C, D)
	labels := []string{"A", "B", "C", "D"}
	options := make([]*domain.GameQuestionOption, 0, 4)
	for j, word := range allAnswers {
		option := &domain.GameQuestionOption{
			QuestionID:   question.ID, // Will be set after question is saved
			OptionLabel:  labels[j],
			TargetWordID: word.ID,
			IsCorrect:    j == correctIndex,
		}
		options = append(options, option)
	}

	return options
}

// findCorrectAnswerIndex finds the index of the correct answer in the shuffled options
func (h *Handler) findCorrectAnswerIndex(allAnswers []*dictdomain.Word, correctWordID int64) int {
	for idx, word := range allAnswers {
		if word.ID == correctWordID {
			return idx
		}
	}
	return -1
}

// logGenerationPerformance logs the performance of question generation
func (h *Handler) logGenerationPerformance(startTime time.Time, sessionID int64, questionCount int) {
	duration := time.Since(startTime)
	h.logger.Info("questions generated",
		logger.Int64("session_id", sessionID),
		logger.Int("question_count", questionCount),
		logger.Duration("duration", duration),
	)

	// Ensure generation completes within 1 second (SC-003)
	if duration > time.Second {
		h.logger.Warn("question generation took longer than 1 second",
			logger.Duration("duration", duration),
		)
	}
}
