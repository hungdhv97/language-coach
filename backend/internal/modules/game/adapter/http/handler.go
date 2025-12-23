package http

import (
	"net/http"
	"strconv"

	"github.com/english-coach/backend/internal/modules/game/domain"
	gamecreatesession "github.com/english-coach/backend/internal/modules/game/usecase/create_session"
	gamesubmitanswer "github.com/english-coach/backend/internal/modules/game/usecase/submit_answer"
	sharederrors "github.com/english-coach/backend/internal/shared/errors"
	"github.com/english-coach/backend/internal/shared/logger"
	"github.com/english-coach/backend/internal/shared/response"
	"github.com/english-coach/backend/internal/transport/http/middleware"
	"github.com/gin-gonic/gin"
)

// Handler handles game-related HTTP requests
type Handler struct {
	createSessionUC *gamecreatesession.Handler
	submitAnswerUC  *gamesubmitanswer.Handler
	questionRepo    domain.GameQuestionRepository
	sessionRepo     domain.GameSessionRepository
	logger          logger.ILogger
}

// NewHandler creates a new game handler
func NewHandler(
	createSessionUC *gamecreatesession.Handler,
	submitAnswerUC *gamesubmitanswer.Handler,
	questionRepo domain.GameQuestionRepository,
	sessionRepo domain.GameSessionRepository,
	logger logger.ILogger,
) *Handler {
	return &Handler{
		createSessionUC: createSessionUC,
		submitAnswerUC:  submitAnswerUC,
		questionRepo:    questionRepo,
		sessionRepo:     sessionRepo,
		logger:          logger,
	}
}

// CreateSession handles POST /api/v1/games/sessions
func (h *Handler) CreateSession(c *gin.Context) {
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
			middleware.SetError(c, sharederrors.ErrInvalidParameter.WithDetails("invalid user_id"))
			return
		}
		userIDInt64 = parsed
	default:
		userIDInt64 = 1 // Default for development
	}

	// Bind request
	var req CreateSessionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		middleware.SetError(c, sharederrors.ErrInvalidRequest.WithDetails(err.Error()))
		return
	}

	// Convert to use case input
	input := gamecreatesession.CreateSessionInput{
		Mode:             req.Mode,
		SourceLanguageID: req.SourceLanguageID,
		TargetLanguageID: req.TargetLanguageID,
		LevelID:          req.LevelID,
		TopicIDs:         req.TopicIDs,
	}

	// Validate request
	if err := input.Validate(); err != nil {
		middleware.SetError(c, sharederrors.ErrValidationError.WithDetails(err.Error()))
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

	// Log game session creation start
	appLogger.Info("game session creation started",
		logger.Int64("user_id", userIDInt64),
		logger.String("mode", input.Mode),
		logger.Int("source_language_id", int(input.SourceLanguageID)),
		logger.Int("target_language_id", int(input.TargetLanguageID)),
		logger.Int64("level_id", input.LevelID),
		logger.Any("topic_ids", input.TopicIDs),
	)

	// Execute use case
	session, err := h.createSessionUC.Execute(ctx, input, userIDInt64)
	if err != nil {
		middleware.SetError(c, err)
		return
	}

	// Log successful session creation
	appLogger.Info("game session created successfully",
		logger.Int64("session_id", session.ID),
		logger.Int64("user_id", userIDInt64),
		logger.Int("total_questions", int(session.TotalQuestions)),
	)

	resp := CreateSessionResponse{
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
	}
	if session.EndedAt != nil {
		resp.EndedAt = session.EndedAt
	}

	response.Success(c, http.StatusCreated, resp)
}

// GetSession handles GET /api/v1/games/sessions/{sessionId}
func (h *Handler) GetSession(c *gin.Context) {
	ctx := c.Request.Context()

	var req GetSessionRequest
	if err := c.ShouldBindUri(&req); err != nil {
		response.ErrorResponse(c, http.StatusBadRequest,
			"INVALID_PARAMETER",
			"ID phiên chơi không hợp lệ",
			nil,
		)
		return
	}
	sessionID := req.SessionID

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
		middleware.SetError(c, sharederrors.ErrForbidden)
		return
	}

	// Get questions and options
	questions, options, err := h.questionRepo.FindBySessionID(ctx, sessionID)
	if err != nil {
		middleware.SetError(c, err)
		return
	}

	// Group options by question ID
	optionsByQuestion := make(map[int64][]*domain.GameQuestionOption)
	for _, opt := range options {
		optionsByQuestion[opt.QuestionID] = append(optionsByQuestion[opt.QuestionID], opt)
	}

	// Map session to response DTO
	sessionResp := GameSessionResponse{
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
	}
	if session.EndedAt != nil {
		sessionResp.EndedAt = session.EndedAt
	}

	// Build response
	questionsWithOptions := make([]QuestionWithOptions, 0, len(questions))
	for _, q := range questions {
		questionsWithOptions = append(questionsWithOptions, QuestionWithOptions{
			GameQuestionResponse: GameQuestionResponse{
				ID:                  q.ID,
				SessionID:           q.SessionID,
				QuestionOrder:       q.QuestionOrder,
				QuestionType:        q.QuestionType,
				SourceWordID:        q.SourceWordID,
				SourceSenseID:       q.SourceSenseID,
				CorrectTargetWordID: q.CorrectTargetWordID,
				SourceLanguageID:    q.SourceLanguageID,
				TargetLanguageID:    q.TargetLanguageID,
				CreatedAt:           q.CreatedAt,
			},
			Options: optionsByQuestion[q.ID],
		})
	}

	response.Success(c, http.StatusOK, GetSessionResponse{
		Session:   sessionResp,
		Questions: questionsWithOptions,
	})
}

// SubmitAnswer handles POST /api/v1/games/sessions/{sessionId}/answers
func (h *Handler) SubmitAnswer(c *gin.Context) {
	ctx := c.Request.Context()

	var pathReq GetSessionRequest
	if err := c.ShouldBindUri(&pathReq); err != nil {
		response.ErrorResponse(c, http.StatusBadRequest,
			"INVALID_PARAMETER",
			"ID phiên chơi không hợp lệ",
			nil,
		)
		return
	}
	sessionID := pathReq.SessionID

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
	var req SubmitAnswerRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		middleware.SetError(c, sharederrors.ErrInvalidRequest.WithDetails(err.Error()))
		return
	}

	// Convert to use case input
	input := gamesubmitanswer.SubmitAnswerInput{
		QuestionID:       req.QuestionID,
		SelectedOptionID: req.SelectedOptionID,
		ResponseTimeMs:   req.ResponseTimeMs,
	}

	// Execute use case
	answer, err := h.submitAnswerUC.Execute(ctx, input, sessionID, userIDInt64)
	if err != nil {
		middleware.SetError(c, err)
		return
	}

	resp := SubmitAnswerResponse{
		ID:               answer.ID,
		QuestionID:       answer.QuestionID,
		SessionID:        answer.SessionID,
		UserID:           answer.UserID,
		SelectedOptionID: answer.SelectedOptionID,
		IsCorrect:        answer.IsCorrect,
		ResponseTimeMs:   answer.ResponseTimeMs,
		AnsweredAt:       answer.AnsweredAt,
	}

	response.Success(c, http.StatusCreated, resp)
}
