package http

import (
	"net/http"

	"github.com/english-coach/backend/internal/modules/user/domain"
	usergetprofile "github.com/english-coach/backend/internal/modules/user/usecase/get_profile"
	userlogin "github.com/english-coach/backend/internal/modules/user/usecase/login"
	userregister "github.com/english-coach/backend/internal/modules/user/usecase/register"
	userupdateprofile "github.com/english-coach/backend/internal/modules/user/usecase/update_profile"
	sharederrors "github.com/english-coach/backend/internal/shared/errors"
	"github.com/english-coach/backend/internal/shared/response"
	"github.com/english-coach/backend/internal/transport/http/middleware"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// Handler handles user-related HTTP requests
type Handler struct {
	registerUC      *userregister.Handler
	loginUC         *userlogin.Handler
	getProfileUC    *usergetprofile.Handler
	updateProfileUC *userupdateprofile.Handler
	userRepo        domain.UserRepository
	profileRepo     domain.UserProfileRepository
	logger          *zap.Logger
}

// NewHandler creates a new user handler
func NewHandler(
	registerUC *userregister.Handler,
	loginUC *userlogin.Handler,
	getProfileUC *usergetprofile.Handler,
	updateProfileUC *userupdateprofile.Handler,
	userRepo domain.UserRepository,
	profileRepo domain.UserProfileRepository,
	logger *zap.Logger,
) *Handler {
	return &Handler{
		registerUC:      registerUC,
		loginUC:         loginUC,
		getProfileUC:    getProfileUC,
		updateProfileUC: updateProfileUC,
		userRepo:        userRepo,
		profileRepo:     profileRepo,
		logger:          logger,
	}
}

// Register handles POST /api/v1/auth/register
func (h *Handler) Register(c *gin.Context) {
	ctx := c.Request.Context()

	var req RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		middleware.SetError(c, sharederrors.ErrInvalidRequest.WithDetails(err.Error()))
		return
	}

	// Validate that at least email or username is provided
	if (req.Email == nil || *req.Email == "") && (req.Username == nil || *req.Username == "") {
		middleware.SetError(c, domain.ErrEmailRequired)
		return
	}

	result, err := h.registerUC.Execute(ctx, userregister.Input{
		DisplayName: req.DisplayName,
		Email:       req.Email,
		Username:    req.Username,
		Password:    req.Password,
	})

	if err != nil {
		middleware.SetError(c, err)
		return
	}

	// Create profile with display_name if provided
	if req.DisplayName != nil && *req.DisplayName != "" {
		_, err := h.profileRepo.Create(ctx, result.UserID, req.DisplayName, nil, nil, nil)
		if err != nil {
			h.logger.Warn("failed to create user profile",
				zap.Error(err),
				zap.Int64("user_id", result.UserID),
			)
			// Don't fail registration if profile creation fails
		}
	}

	response.Success(c, http.StatusCreated, result)
}

// Login handles POST /api/v1/auth/login
func (h *Handler) Login(c *gin.Context) {
	ctx := c.Request.Context()

	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		middleware.SetError(c, sharederrors.ErrInvalidRequest.WithDetails(err.Error()))
		return
	}

	// Validate that at least email or username is provided
	if (req.Email == nil || *req.Email == "") && (req.Username == nil || *req.Username == "") {
		middleware.SetError(c, domain.ErrEmailRequired)
		return
	}

	result, err := h.loginUC.Execute(ctx, userlogin.Input{
		Email:    req.Email,
		Username: req.Username,
		Password: req.Password,
	})

	if err != nil {
		middleware.SetError(c, err)
		return
	}

	response.Success(c, http.StatusOK, result)
}

// GetProfile handles GET /api/v1/users/profile
func (h *Handler) GetProfile(c *gin.Context) {
	ctx := c.Request.Context()

	// Get user ID from context (set by auth middleware)
	userID, exists := c.Get("user_id")
	if !exists {
		middleware.SetError(c, sharederrors.ErrUnauthorized)
		return
	}

	userIDInt64, ok := userID.(int64)
	if !ok {
		middleware.SetError(c, sharederrors.ErrInternalError)
		return
	}

	profile, err := h.getProfileUC.Execute(ctx, usergetprofile.Input{UserID: userIDInt64})
	if err != nil {
		middleware.SetError(c, err)
		return
	}

	response.Success(c, http.StatusOK, profile)
}

// UpdateProfile handles PUT /api/v1/users/profile
func (h *Handler) UpdateProfile(c *gin.Context) {
	ctx := c.Request.Context()

	// Get user ID from context (set by auth middleware)
	userID, exists := c.Get("user_id")
	if !exists {
		middleware.SetError(c, sharederrors.ErrUnauthorized)
		return
	}

	userIDInt64, ok := userID.(int64)
	if !ok {
		middleware.SetError(c, sharederrors.ErrInternalError)
		return
	}

	var req UpdateProfileRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		middleware.SetError(c, sharederrors.ErrInvalidRequest.WithDetails(err.Error()))
		return
	}

	result, err := h.updateProfileUC.Execute(ctx, userIDInt64, userupdateprofile.Input{
		DisplayName: req.DisplayName,
		AvatarURL:   req.AvatarURL,
		BirthDay:    req.BirthDay,
		Bio:         req.Bio,
	})

	if err != nil {
		middleware.SetError(c, err)
		return
	}

	response.Success(c, http.StatusOK, result)
}

// CheckEmailAvailability handles GET /api/v1/auth/check-email?email=...
func (h *Handler) CheckEmailAvailability(c *gin.Context) {
	ctx := c.Request.Context()
	email := c.Query("email")

	if email == "" {
		middleware.SetError(c, sharederrors.ErrInvalidParameter.WithDetails("email parameter is required"))
		return
	}

	exists, err := h.userRepo.CheckEmailExists(ctx, email)
	if err != nil {
		middleware.SetError(c, err)
		return
	}

	response.Success(c, http.StatusOK, CheckEmailAvailabilityResponse{
		Available: !exists,
		Exists:    exists,
	})
}

// CheckUsernameAvailability handles GET /api/v1/auth/check-username?username=...
func (h *Handler) CheckUsernameAvailability(c *gin.Context) {
	ctx := c.Request.Context()
	username := c.Query("username")

	if username == "" {
		middleware.SetError(c, sharederrors.ErrInvalidParameter.WithDetails("username parameter is required"))
		return
	}

	exists, err := h.userRepo.CheckUsernameExists(ctx, username)
	if err != nil {
		middleware.SetError(c, err)
		return
	}

	response.Success(c, http.StatusOK, CheckUsernameAvailabilityResponse{
		Available: !exists,
		Exists:    exists,
	})
}
