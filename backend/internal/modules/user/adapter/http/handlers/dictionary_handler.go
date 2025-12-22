package handler

import (
	"net/http"
	"strconv"

	"github.com/english-coach/backend/internal/modules/dictionary/domain"
	dictusecase "github.com/english-coach/backend/internal/modules/dictionary/usecase/get_word_detail"
	"github.com/english-coach/backend/internal/transport/http/middleware"
	sharederrors "github.com/english-coach/backend/internal/shared/errors"
	"github.com/english-coach/backend/internal/shared/response"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// DictionaryHandler handles dictionary-related HTTP requests
type DictionaryHandler struct {
	languageRepo      domain.LanguageRepository
	topicRepo         domain.TopicRepository
	levelRepo         domain.LevelRepository
	wordRepo        domain.WordRepository
	getWordDetailUC *dictusecase.Handler
	logger          *zap.Logger
}

// NewDictionaryHandler creates a new dictionary handler
func NewDictionaryHandler(
	languageRepo domain.LanguageRepository,
	topicRepo domain.TopicRepository,
	levelRepo domain.LevelRepository,
	wordRepo domain.WordRepository,
	getWordDetailUC *dictusecase.Handler,
	logger *zap.Logger,
) *DictionaryHandler {
	return &DictionaryHandler{
		languageRepo:    languageRepo,
		topicRepo:       topicRepo,
		levelRepo:       levelRepo,
		wordRepo:        wordRepo,
		getWordDetailUC: getWordDetailUC,
		logger:          logger,
	}
}

// GetLanguages handles GET /api/v1/reference/languages
func (h *DictionaryHandler) GetLanguages(c *gin.Context) {
	ctx := c.Request.Context()

	languages, err := h.languageRepo.FindAll(ctx)
	if err != nil {
		middleware.SetError(c, err)
		return
	}

	response.Success(c, http.StatusOK, languages)
}

// GetTopics handles GET /api/v1/reference/topics
func (h *DictionaryHandler) GetTopics(c *gin.Context) {
	ctx := c.Request.Context()

	topics, err := h.topicRepo.FindAll(ctx)
	if err != nil {
		middleware.SetError(c, err)
		return
	}

	response.Success(c, http.StatusOK, topics)
}

// GetLevels handles GET /api/v1/reference/levels?languageId=...
func (h *DictionaryHandler) GetLevels(c *gin.Context) {
	ctx := c.Request.Context()

	languageIDStr := c.Query("languageId")
	if languageIDStr != "" {
		languageID, err := strconv.ParseInt(languageIDStr, 10, 16)
		if err != nil {
			middleware.SetError(c, sharederrors.ErrInvalidParameter.WithDetails("invalid languageId"))
			return
		}

		levels, err := h.levelRepo.FindByLanguageID(ctx, int16(languageID))
		if err != nil {
			middleware.SetError(c, err)
			return
		}

		response.Success(c, http.StatusOK, levels)
		return
	}

	// If no languageId provided, return all levels
	levels, err := h.levelRepo.FindAll(ctx)
	if err != nil {
		middleware.SetError(c, err)
		return
	}

	response.Success(c, http.StatusOK, levels)
}

// SearchWords handles GET /api/v1/dictionary/search?q=...&languageId=...&limit=...&offset=...
func (h *DictionaryHandler) SearchWords(c *gin.Context) {
	ctx := c.Request.Context()

	query := c.Query("q")
	if query == "" {
		middleware.SetError(c, sharederrors.ErrInvalidParameter.WithDetails("query parameter (q) is required"))
		return
	}

	// Parse language ID (required)
	languageIDStr := c.Query("languageId")
	if languageIDStr == "" {
		middleware.SetError(c, sharederrors.ErrInvalidParameter.WithDetails("languageId parameter is required"))
		return
	}

	languageID, err := strconv.ParseInt(languageIDStr, 10, 16)
	if err != nil {
		middleware.SetError(c, sharederrors.ErrInvalidParameter.WithDetails("invalid languageId"))
		return
	}
	langID := int16(languageID)

	// Parse pagination parameters
	limit := 20 // default limit
	if limitStr := c.Query("limit"); limitStr != "" {
		parsedLimit, err := strconv.Atoi(limitStr)
		if err != nil || parsedLimit < 1 || parsedLimit > 100 {
			middleware.SetError(c, sharederrors.ErrInvalidParameter.WithDetails("limit must be between 1 and 100"))
			return
		}
		limit = parsedLimit
	}

	offset := 0
	if offsetStr := c.Query("offset"); offsetStr != "" {
		parsedOffset, err := strconv.Atoi(offsetStr)
		if err != nil || parsedOffset < 0 {
			middleware.SetError(c, sharederrors.ErrInvalidParameter.WithDetails("offset must be non-negative"))
			return
		}
		offset = parsedOffset
	}

	// Get request logger from context (includes request ID)
	requestLogger, _ := c.Get("logger")
	var logger *zap.Logger
	if reqLogger, ok := requestLogger.(*zap.Logger); ok {
		logger = reqLogger
	} else {
		logger = h.logger
	}

	// Log dictionary search start
	logger.Info("dictionary search started",
		zap.String("query", query),
		zap.Int16("language_id", langID),
		zap.Int("limit", limit),
		zap.Int("offset", offset),
	)

	// Search words
	words, err := h.wordRepo.SearchWords(ctx, query, langID, limit, offset)
	if err != nil {
		middleware.SetError(c, err)
		return
	}

	// Get total count for pagination
	total, err := h.wordRepo.CountSearchWords(ctx, query, langID)
	if err != nil {
		h.logger.Error("failed to count search words",
			zap.Error(err),
			zap.String("query", query),
		)
		// Continue without total count
		total = len(words)
	}

	// Handle empty results gracefully
	if len(words) == 0 {
		logger.Info("dictionary search completed - no results",
			zap.String("query", query),
			zap.Int("total", total),
		)
		response.Success(c, http.StatusOK, gin.H{
			"words":  []interface{}{},
			"total":  total,
			"limit":  limit,
			"offset": offset,
		})
		return
	}

	// Log successful search
	logger.Info("dictionary search completed",
		zap.String("query", query),
		zap.Int("results_count", len(words)),
		zap.Int("total", total),
	)

	response.Success(c, http.StatusOK, gin.H{
		"words":  words,
		"total":  total,
		"limit":  limit,
		"offset": offset,
	})
}

// GetWordDetail handles GET /api/v1/dictionary/words/:wordId
func (h *DictionaryHandler) GetWordDetail(c *gin.Context) {
	ctx := c.Request.Context()

	wordIDStr := c.Param("wordId")
	wordID, err := strconv.ParseInt(wordIDStr, 10, 64)
	if err != nil {
		middleware.SetError(c, sharederrors.ErrInvalidParameter.WithDetails("invalid wordId"))
		return
	}

	// Get request logger from context (includes request ID)
	requestLogger, _ := c.Get("logger")
	var logger *zap.Logger
	if reqLogger, ok := requestLogger.(*zap.Logger); ok {
		logger = reqLogger
	} else {
		logger = h.logger
	}

	// Log word detail lookup start
	logger.Info("word detail lookup started",
		zap.Int64("word_id", wordID),
	)

	wordDetail, err := h.getWordDetailUC.Execute(ctx, dictusecase.Input{WordID: wordID})
	if err != nil {
		middleware.SetError(c, err)
		return
	}

	if wordDetail == nil || wordDetail.Word == nil {
		middleware.SetError(c, sharederrors.ErrNotFound)
		return
	}

	// Log successful word detail lookup
	logger.Info("word detail lookup completed",
		zap.Int64("word_id", wordID),
		zap.Int("senses_count", len(wordDetail.Senses)),
		zap.Int("pronunciations_count", len(wordDetail.Pronunciations)),
	)

	response.Success(c, http.StatusOK, wordDetail)
}
