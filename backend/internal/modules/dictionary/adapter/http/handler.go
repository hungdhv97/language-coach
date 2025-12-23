package http

import (
	"net/http"
	"strconv"

	"github.com/english-coach/backend/internal/modules/dictionary/domain"
	dictusecase "github.com/english-coach/backend/internal/modules/dictionary/usecase/get_word_detail"
	sharederrors "github.com/english-coach/backend/internal/shared/errors"
	"github.com/english-coach/backend/internal/shared/logger"
	"github.com/english-coach/backend/internal/shared/pagination"
	"github.com/english-coach/backend/internal/shared/response"
	"github.com/english-coach/backend/internal/transport/http/middleware"
	"github.com/gin-gonic/gin"
)

// Handler handles dictionary-related HTTP requests
type Handler struct {
	languageRepo    domain.LanguageRepository
	topicRepo       domain.TopicRepository
	levelRepo       domain.LevelRepository
	wordRepo        domain.WordRepository
	getWordDetailUC *dictusecase.Handler
	logger          logger.ILogger
}

// NewHandler creates a new dictionary handler
func NewHandler(
	languageRepo domain.LanguageRepository,
	topicRepo domain.TopicRepository,
	levelRepo domain.LevelRepository,
	wordRepo domain.WordRepository,
	getWordDetailUC *dictusecase.Handler,
	logger logger.ILogger,
) *Handler {
	return &Handler{
		languageRepo:    languageRepo,
		topicRepo:       topicRepo,
		levelRepo:       levelRepo,
		wordRepo:        wordRepo,
		getWordDetailUC: getWordDetailUC,
		logger:          logger,
	}
}

// GetLanguages handles GET /api/v1/reference/languages
func (h *Handler) GetLanguages(c *gin.Context) {
	ctx := c.Request.Context()

	languages, err := h.languageRepo.FindAll(ctx)
	if err != nil {
		middleware.SetError(c, err)
		return
	}

	response.Success(c, http.StatusOK, languages)
}

// GetTopics handles GET /api/v1/reference/topics
func (h *Handler) GetTopics(c *gin.Context) {
	ctx := c.Request.Context()

	topics, err := h.topicRepo.FindAll(ctx)
	if err != nil {
		middleware.SetError(c, err)
		return
	}

	response.Success(c, http.StatusOK, topics)
}

// GetLevels handles GET /api/v1/reference/levels?languageId=...
func (h *Handler) GetLevels(c *gin.Context) {
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
func (h *Handler) SearchWords(c *gin.Context) {
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
	paginationParams, err := pagination.ParseFromQuery(c)
	if err != nil {
		middleware.SetError(c, err)
		return
	}

	// Get request logger from context (includes request ID)
	requestLogger, _ := c.Get("logger")
	var appLogger logger.ILogger
	if reqLogger, ok := requestLogger.(logger.ILogger); ok {
		appLogger = reqLogger
	} else {
		appLogger = h.logger
	}

	// Log dictionary search start
	appLogger.Info("dictionary search started",
		logger.String("query", query),
		logger.Int("language_id", int(langID)),
		logger.Int("limit", paginationParams.Limit),
		logger.Int("offset", paginationParams.Offset),
		logger.Int("page", paginationParams.Page),
		logger.Int("pageSize", paginationParams.Size),
	)

	// Search words
	words, err := h.wordRepo.SearchWords(ctx, query, langID, paginationParams.Limit, paginationParams.Offset)
	if err != nil {
		middleware.SetError(c, err)
		return
	}

	// Get total count for pagination
	totalCount, err := h.wordRepo.CountSearchWords(ctx, query, langID)
	if err != nil {
		h.logger.Error("failed to count search words",
			logger.Error(err),
			logger.String("query", query),
		)
		// Continue without total count
		totalCount = len(words)
	}

	total := int64(totalCount)

	// Log successful search
	appLogger.Info("dictionary search completed",
		logger.String("query", query),
		logger.Int("results_count", len(words)),
		logger.Int64("total", total),
	)

	// Return paginated response
	response.Paginated(c, http.StatusOK, words, paginationParams, total)
}

// GetWordDetail handles GET /api/v1/dictionary/words/:wordId
func (h *Handler) GetWordDetail(c *gin.Context) {
	ctx := c.Request.Context()

	wordIDStr := c.Param("wordId")
	wordID, err := strconv.ParseInt(wordIDStr, 10, 64)
	if err != nil {
		middleware.SetError(c, sharederrors.ErrInvalidParameter.WithDetails("invalid wordId"))
		return
	}

	// Get request logger from context (includes request ID)
	requestLogger, _ := c.Get("logger")
	var appLogger logger.ILogger
	if reqLogger, ok := requestLogger.(logger.ILogger); ok {
		appLogger = reqLogger
	} else {
		appLogger = h.logger
	}

	// Log word detail lookup start
	appLogger.Info("word detail lookup started",
		logger.Int64("word_id", wordID),
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
	appLogger.Info("word detail lookup completed",
		logger.Int64("word_id", wordID),
		logger.Int("senses_count", len(wordDetail.Senses)),
		logger.Int("pronunciations_count", len(wordDetail.Pronunciations)),
	)

	response.Success(c, http.StatusOK, wordDetail)
}
