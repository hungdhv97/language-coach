package handler

import (
	"net/http"
	"strconv"

	"github.com/english-coach/backend/internal/domain/game/usecase/query"
	"github.com/english-coach/backend/internal/shared/response"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// StatisticsHandler handles statistics-related HTTP requests
type StatisticsHandler struct {
	getStatisticsUC *query.GetSessionStatisticsUseCase
	logger          *zap.Logger
}

// NewStatisticsHandler creates a new statistics handler
func NewStatisticsHandler(
	getStatisticsUC *query.GetSessionStatisticsUseCase,
	logger *zap.Logger,
) *StatisticsHandler {
	return &StatisticsHandler{
		getStatisticsUC: getStatisticsUC,
		logger:          logger,
	}
}

// GetSessionStatistics handles GET /api/v1/statistics/sessions/{sessionId}
func (h *StatisticsHandler) GetSessionStatistics(c *gin.Context) {
	ctx := c.Request.Context()

	// Get session ID from path
	sessionIDStr := c.Param("sessionId")
	sessionID, err := strconv.ParseInt(sessionIDStr, 10, 64)
	if err != nil {
		response.ErrorResponse(c, http.StatusBadRequest,
			"INVALID_SESSION_ID",
			"Invalid session ID",
			nil,
		)
		return
	}

	// Get user ID from context (set by auth middleware)
	// For now, use a default user ID if not authenticated
	userID, exists := c.Get("user_id")
	if !exists {
		// TODO: In production, this should require authentication
		// For now, use a default user ID for development
		userID = int64(1)
	}

	var userIDInt64 int64
	switch v := userID.(type) {
	case int64:
		userIDInt64 = v
	case int:
		userIDInt64 = int64(v)
	case string:
		parsed, err := strconv.ParseInt(v, 10, 64)
		if err != nil {
			response.ErrorResponse(c, http.StatusBadRequest,
				"INVALID_USER_ID",
				"Invalid user ID",
				nil,
			)
			return
		}
		userIDInt64 = parsed
	default:
		userIDInt64 = 1 // Default for development
	}

	// Execute use case
	stats, err := h.getStatisticsUC.Execute(ctx, sessionID, userIDInt64)
	if err != nil {
		h.logger.Error("failed to get session statistics",
			zap.Error(err),
			zap.Int64("session_id", sessionID),
			zap.Int64("user_id", userIDInt64),
		)

		// Check for specific errors
		if err.Error() == "session does not belong to user" {
			response.ErrorResponse(c, http.StatusForbidden,
				"FORBIDDEN",
				"You do not have permission to access this game session",
				nil,
			)
			return
		}

		if err.Error() == "failed to find session" || err.Error() == "sql: no rows in result set" {
			response.ErrorResponse(c, http.StatusNotFound,
				"SESSION_NOT_FOUND",
				"Game session not found",
				nil,
			)
			return
		}

		response.ErrorResponse(c, http.StatusInternalServerError,
			"INTERNAL_ERROR",
			"Failed to get session statistics",
			nil,
		)
		return
	}

	// Return success response
	response.Success(c, http.StatusOK, stats)
}
