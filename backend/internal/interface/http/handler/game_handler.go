package handler

import (
	"net/http"
	"strconv"

	"github.com/english-coach/backend/internal/domain/game/dto"
	"github.com/english-coach/backend/internal/domain/game/model"
	"github.com/english-coach/backend/internal/domain/game/port"
	"github.com/english-coach/backend/internal/domain/game/usecase/command"
	"github.com/english-coach/backend/internal/interface/http/middleware"
	commonerrors "github.com/english-coach/backend/internal/shared/errors"
	"github.com/english-coach/backend/internal/shared/response"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// GameHandler handles game-related HTTP requests
type GameHandler struct {
	createSessionUC *command.CreateGameSessionUseCase
	submitAnswerUC  *command.SubmitAnswerUseCase
	endSessionUC    *command.EndGameSessionUseCase
	questionRepo    port.GameQuestionRepository
	sessionRepo     port.GameSessionRepository
	logger          *zap.Logger
}

// NewGameHandler creates a new game handler
func NewGameHandler(
	createSessionUC *command.CreateGameSessionUseCase,
	submitAnswerUC *command.SubmitAnswerUseCase,
	endSessionUC *command.EndGameSessionUseCase,
	questionRepo port.GameQuestionRepository,
	sessionRepo port.GameSessionRepository,
	logger *zap.Logger,
) *GameHandler {
	return &GameHandler{
		createSessionUC: createSessionUC,
		submitAnswerUC:  submitAnswerUC,
		endSessionUC:    endSessionUC,
		questionRepo:    questionRepo,
		sessionRepo:     sessionRepo,
		logger:          logger,
	}
}

// CreateSession handles POST /api/v1/games/sessions
func (h *GameHandler) CreateSession(c *gin.Context) {
	ctx := c.Request.Context()

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
			middleware.SetError(c, commonerrors.ErrInvalidParameter.WithDetails("invalid user_id"))
			return
		}
		userIDInt64 = parsed
	default:
		userIDInt64 = 1 // Default for development
	}

	// Bind request
	var req dto.CreateGameSessionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		middleware.SetError(c, commonerrors.ErrInvalidRequest.WithDetails(err.Error()))
		return
	}

	// Validate request
	if err := req.Validate(); err != nil {
		middleware.SetError(c, commonerrors.ErrValidationError.WithDetails(err.Error()))
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

	// Log game session creation start
	logger.Info("game session creation started",
		zap.Int64("user_id", userIDInt64),
		zap.String("mode", req.Mode),
		zap.Int16("source_language_id", req.SourceLanguageID),
		zap.Int16("target_language_id", req.TargetLanguageID),
		zap.Int64("level_id", req.LevelID),
		zap.Any("topic_ids", req.TopicIDs),
	)

	// Execute use case
	session, err := h.createSessionUC.Execute(ctx, &req, userIDInt64)
	if err != nil {
		middleware.SetError(c, err)
		return
	}

	// Log successful session creation
	logger.Info("game session created successfully",
		zap.Int64("session_id", session.ID),
		zap.Int64("user_id", userIDInt64),
		zap.Int16("total_questions", session.TotalQuestions),
	)

	response.Success(c, http.StatusCreated, session)
}

// GetSession handles GET /api/v1/games/sessions/{sessionId}
func (h *GameHandler) GetSession(c *gin.Context) {
	ctx := c.Request.Context()

	sessionIDStr := c.Param("sessionId")
	sessionID, err := strconv.ParseInt(sessionIDStr, 10, 64)
	if err != nil {
		response.ErrorResponse(c, http.StatusBadRequest,
			"INVALID_PARAMETER",
			"ID phiên chơi không hợp lệ",
			nil,
		)
		return
	}

	// Get user ID
	userID, exists := c.Get("user_id")
	if !exists {
		userID = int64(1) // Default for development
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
			userIDInt64 = 1
		} else {
			userIDInt64 = parsed
		}
	default:
		userIDInt64 = 1
	}

	// Get session
	session, err := h.sessionRepo.FindByID(ctx, sessionID)
	if err != nil {
		middleware.SetError(c, err)
		return
	}

	// Verify user owns session
	if session.UserID != userIDInt64 {
		middleware.SetError(c, commonerrors.ErrForbidden)
		return
	}

	// Get questions and options
	questions, options, err := h.questionRepo.FindBySessionID(ctx, sessionID)
	if err != nil {
		middleware.SetError(c, err)
		return
	}

	// Group options by question ID
	optionsByQuestion := make(map[int64][]*model.GameQuestionOption)
	for _, opt := range options {
		optionsByQuestion[opt.QuestionID] = append(optionsByQuestion[opt.QuestionID], opt)
	}

	// Build response
	type QuestionWithOptions struct {
		*model.GameQuestion
		Options []*model.GameQuestionOption `json:"options"`
	}

	questionsWithOptions := make([]QuestionWithOptions, 0, len(questions))
	for _, q := range questions {
		questionsWithOptions = append(questionsWithOptions, QuestionWithOptions{
			GameQuestion: q,
			Options:      optionsByQuestion[q.ID],
		})
	}

	responseData := gin.H{
		"session":   session,
		"questions": questionsWithOptions,
	}

	response.Success(c, http.StatusOK, responseData)
}

// SubmitAnswer handles POST /api/v1/games/sessions/{sessionId}/answers
func (h *GameHandler) SubmitAnswer(c *gin.Context) {
	ctx := c.Request.Context()

	sessionIDStr := c.Param("sessionId")
	sessionID, err := strconv.ParseInt(sessionIDStr, 10, 64)
	if err != nil {
		response.ErrorResponse(c, http.StatusBadRequest,
			"INVALID_PARAMETER",
			"ID phiên chơi không hợp lệ",
			nil,
		)
		return
	}

	// Get user ID
	userID, exists := c.Get("user_id")
	if !exists {
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
			userIDInt64 = 1
		} else {
			userIDInt64 = parsed
		}
	default:
		userIDInt64 = 1
	}

	// Bind request
	var req dto.SubmitAnswerRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		middleware.SetError(c, commonerrors.ErrInvalidRequest.WithDetails(err.Error()))
		return
	}

	// Execute use case
	answer, err := h.submitAnswerUC.Execute(ctx, &req, sessionID, userIDInt64)
	if err != nil {
		middleware.SetError(c, err)
		return
	}

	response.Success(c, http.StatusCreated, answer)
}
