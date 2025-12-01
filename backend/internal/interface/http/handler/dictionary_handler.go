package handler

import (
	"net/http"
	"strconv"

	"github.com/english-coach/backend/internal/domain/dictionary/port"
	"github.com/english-coach/backend/internal/domain/dictionary/service"
	"github.com/english-coach/backend/internal/shared/response"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// DictionaryHandler handles dictionary-related HTTP requests
type DictionaryHandler struct {
	languageRepo     port.LanguageRepository
	topicRepo        port.TopicRepository
	levelRepo        port.LevelRepository
	wordRepo         port.WordRepository
	dictionaryService *service.DictionaryService
	logger           *zap.Logger
}

// NewDictionaryHandler creates a new dictionary handler
func NewDictionaryHandler(
	languageRepo port.LanguageRepository,
	topicRepo port.TopicRepository,
	levelRepo port.LevelRepository,
	wordRepo port.WordRepository,
	dictionaryService *service.DictionaryService,
	logger *zap.Logger,
) *DictionaryHandler {
	return &DictionaryHandler{
		languageRepo:      languageRepo,
		topicRepo:         topicRepo,
		levelRepo:         levelRepo,
		wordRepo:          wordRepo,
		dictionaryService: dictionaryService,
		logger:            logger,
	}
}

// GetLanguages handles GET /api/v1/reference/languages
func (h *DictionaryHandler) GetLanguages(c *gin.Context) {
	ctx := c.Request.Context()

	languages, err := h.languageRepo.FindAll(ctx)
	if err != nil {
		h.logger.Error("failed to fetch languages",
			zap.Error(err),
			zap.String("path", c.Request.URL.Path),
		)
		response.ErrorResponse(c, http.StatusInternalServerError,
			"INTERNAL_ERROR",
			"Failed to fetch languages",
			nil,
		)
		return
	}

	response.Success(c, http.StatusOK, languages)
}

// GetTopics handles GET /api/v1/reference/topics
func (h *DictionaryHandler) GetTopics(c *gin.Context) {
	ctx := c.Request.Context()

	topics, err := h.topicRepo.FindAll(ctx)
	if err != nil {
		h.logger.Error("failed to fetch topics",
			zap.Error(err),
			zap.String("path", c.Request.URL.Path),
		)
		response.ErrorResponse(c, http.StatusInternalServerError,
			"INTERNAL_ERROR",
			"Failed to fetch topics",
			nil,
		)
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
			response.ErrorResponse(c, http.StatusBadRequest,
				"INVALID_PARAMETER",
				"Invalid languageId parameter",
				nil,
			)
			return
		}

		levels, err := h.levelRepo.FindByLanguageID(ctx, int16(languageID))
		if err != nil {
			h.logger.Error("failed to fetch levels by language",
				zap.Error(err),
				zap.String("path", c.Request.URL.Path),
				zap.Int16("language_id", int16(languageID)),
			)
			response.ErrorResponse(c, http.StatusInternalServerError,
				"INTERNAL_ERROR",
				"Failed to fetch levels",
				nil,
			)
			return
		}

		response.Success(c, http.StatusOK, levels)
		return
	}

	// If no languageId provided, return all levels
	levels, err := h.levelRepo.FindAll(ctx)
	if err != nil {
		h.logger.Error("failed to fetch levels",
			zap.Error(err),
			zap.String("path", c.Request.URL.Path),
		)
		response.ErrorResponse(c, http.StatusInternalServerError,
			"INTERNAL_ERROR",
			"Failed to fetch levels",
			nil,
		)
		return
	}

	response.Success(c, http.StatusOK, levels)
}

// SearchWords handles GET /api/v1/dictionary/search?q=...&languageId=...&limit=...&offset=...
func (h *DictionaryHandler) SearchWords(c *gin.Context) {
	ctx := c.Request.Context()

	query := c.Query("q")
	if query == "" {
		response.ErrorResponse(c, http.StatusBadRequest,
			"INVALID_PARAMETER",
			"Search parameter (q) is required",
			nil,
		)
		return
	}

	// Parse language ID (optional)
	var languageID *int16
	if languageIDStr := c.Query("languageId"); languageIDStr != "" {
		id, err := strconv.ParseInt(languageIDStr, 10, 16)
		if err != nil {
			response.ErrorResponse(c, http.StatusBadRequest,
				"INVALID_PARAMETER",
				"Invalid languageId parameter",
				nil,
			)
			return
		}
		langID := int16(id)
		languageID = &langID
	}

	// Parse pagination parameters
	limit := 20 // default limit
	if limitStr := c.Query("limit"); limitStr != "" {
		parsedLimit, err := strconv.Atoi(limitStr)
		if err != nil || parsedLimit < 1 || parsedLimit > 100 {
			response.ErrorResponse(c, http.StatusBadRequest,
				"INVALID_PARAMETER",
				"Limit parameter must be a number between 1 and 100",
				nil,
			)
			return
		}
		limit = parsedLimit
	}

	offset := 0
	if offsetStr := c.Query("offset"); offsetStr != "" {
		parsedOffset, err := strconv.Atoi(offsetStr)
		if err != nil || parsedOffset < 0 {
			response.ErrorResponse(c, http.StatusBadRequest,
				"INVALID_PARAMETER",
				"Offset parameter must be a non-negative number",
				nil,
			)
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
		zap.Any("language_id", languageID),
		zap.Int("limit", limit),
		zap.Int("offset", offset),
	)

	// Search words
	words, err := h.wordRepo.SearchWords(ctx, query, languageID, limit, offset)
	if err != nil {
		logger.Error("failed to search words",
			zap.Error(err),
			zap.String("query", query),
			zap.String("path", c.Request.URL.Path),
		)
		response.ErrorResponse(c, http.StatusInternalServerError,
			"INTERNAL_ERROR",
			"Failed to search words",
			nil,
		)
		return
	}

	// Get total count for pagination
	total, err := h.wordRepo.CountSearchWords(ctx, query, languageID)
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
			"words": []interface{}{},
			"total": total,
			"limit": limit,
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
		"words": words,
		"total": total,
		"limit": limit,
		"offset": offset,
	})
}

// GetWordDetail handles GET /api/v1/dictionary/words/:wordId
func (h *DictionaryHandler) GetWordDetail(c *gin.Context) {
	ctx := c.Request.Context()

	wordIDStr := c.Param("wordId")
	wordID, err := strconv.ParseInt(wordIDStr, 10, 64)
	if err != nil {
		response.ErrorResponse(c, http.StatusBadRequest,
			"INVALID_PARAMETER",
			"Invalid wordId parameter",
			nil,
		)
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

	wordDetail, err := h.dictionaryService.GetWordDetail(ctx, wordID)
	if err != nil {
		logger.Error("failed to get word detail",
			zap.Error(err),
			zap.Int64("word_id", wordID),
			zap.String("path", c.Request.URL.Path),
		)
		response.ErrorResponse(c, http.StatusInternalServerError,
			"INTERNAL_ERROR",
			"Failed to get word detail",
			nil,
		)
		return
	}

	if wordDetail == nil || wordDetail.Word == nil {
		response.ErrorResponse(c, http.StatusNotFound,
			"NOT_FOUND",
			"Word not found",
			nil,
		)
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

