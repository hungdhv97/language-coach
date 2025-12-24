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
)

// Handler handles user-related HTTP requests
type Handler struct {
	registerUC      *userregister.Handler
	loginUC         *userlogin.Handler
	getProfileUC    *usergetprofile.Handler
	updateProfileUC *userupdateprofile.Handler
	userRepo    domain.UserRepository
	profileRepo domain.UserProfileRepository
}

// NewHandler creates a new user handler
func NewHandler(
	registerUC *userregister.Handler,
	loginUC *userlogin.Handler,
	getProfileUC *usergetprofile.Handler,
	updateProfileUC *userupdateprofile.Handler,
	userRepo domain.UserRepository,
	profileRepo domain.UserProfileRepository,
) *Handler {
	return &Handler{
		registerUC:      registerUC,
		loginUC:         loginUC,
		getProfileUC:    getProfileUC,
		updateProfileUC: updateProfileUC,
		userRepo:        userRepo,
		profileRepo:     profileRepo,
	}
}

// Register handles POST /api/v1/auth/register
func (h *Handler) Register(c *gin.Context) {
	ctx := c.Request.Context()

	var req RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		// Parse error - return as invalid request
		middleware.SetError(c, sharederrors.NewAppError(
			sharederrors.CodeInvalidRequest,
			"Dữ liệu yêu cầu không hợp lệ",
		).WithMetadata("parse_error", err.Error()))
		return
	}

	result, err := h.registerUC.Execute(ctx, userregister.RegisterInput{
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
	// Note: Profile creation failure doesn't fail registration
	// This is a business decision - registration succeeds even if profile creation fails
	if req.DisplayName != nil && *req.DisplayName != "" {
		_, _ = h.profileRepo.Create(ctx, result.UserID, req.DisplayName, nil, nil, nil)
	}

	resp := RegisterResponse{
		UserID:   result.UserID,
		Email:    result.Email,
		Username: result.Username,
	}

	response.Success(c, http.StatusCreated, resp)
}

// Login handles POST /api/v1/auth/login
func (h *Handler) Login(c *gin.Context) {
	ctx := c.Request.Context()

	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		// Parse error - return as invalid request
		middleware.SetError(c, sharederrors.NewAppError(
			sharederrors.CodeInvalidRequest,
			"Dữ liệu yêu cầu không hợp lệ",
		).WithMetadata("parse_error", err.Error()))
		return
	}

	result, err := h.loginUC.Execute(ctx, userlogin.LoginInput{
		Email:    req.Email,
		Username: req.Username,
		Password: req.Password,
	})

	if err != nil {
		middleware.SetError(c, err)
		return
	}

	resp := LoginResponse{
		Token:    result.Token,
		UserID:   result.UserID,
		Email:    result.Email,
		Username: result.Username,
	}

	response.Success(c, http.StatusOK, resp)
}

// GetProfile handles GET /api/v1/users/profile
func (h *Handler) GetProfile(c *gin.Context) {
	ctx := c.Request.Context()

	// Get user ID from context (set by auth middleware)
	userID, exists := c.Get("user_id")
	if !exists {
		middleware.SetError(c, sharederrors.NewAppError(
			sharederrors.CodeUnauthorized,
			"Người dùng chưa được xác thực",
		))
		return
	}

	userIDInt64, ok := userID.(int64)
	if !ok {
		middleware.SetError(c, sharederrors.NewAppError(
			sharederrors.CodeInternalError,
			"Đã xảy ra lỗi hệ thống",
		))
		return
	}

	profile, err := h.getProfileUC.Execute(ctx, usergetprofile.GetProfileInput{UserID: userIDInt64})
	if err != nil {
		middleware.SetError(c, err)
		return
	}

	resp := UserProfileResponse{
		UserID:      profile.UserID,
		DisplayName: profile.DisplayName,
		AvatarURL:   profile.AvatarURL,
		BirthDay:    profile.BirthDay,
		Bio:         profile.Bio,
	}

	response.Success(c, http.StatusOK, resp)
}

// UpdateProfile handles PUT /api/v1/users/profile
func (h *Handler) UpdateProfile(c *gin.Context) {
	ctx := c.Request.Context()

	// Get user ID from context (set by auth middleware)
	userID, exists := c.Get("user_id")
	if !exists {
		middleware.SetError(c, sharederrors.NewAppError(
			sharederrors.CodeUnauthorized,
			"Người dùng chưa được xác thực",
		))
		return
	}

	userIDInt64, ok := userID.(int64)
	if !ok {
		middleware.SetError(c, sharederrors.NewAppError(
			sharederrors.CodeInternalError,
			"Đã xảy ra lỗi hệ thống",
		))
		return
	}

	var req UpdateProfileRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		// Parse error - return as invalid request
		middleware.SetError(c, sharederrors.NewAppError(
			sharederrors.CodeInvalidRequest,
			"Dữ liệu yêu cầu không hợp lệ",
		).WithMetadata("parse_error", err.Error()))
		return
	}

	result, err := h.updateProfileUC.Execute(ctx, userIDInt64, userupdateprofile.UpdateProfileInput{
		DisplayName: req.DisplayName,
		AvatarURL:   req.AvatarURL,
		BirthDay:    req.BirthDay,
		Bio:         req.Bio,
	})

	if err != nil {
		middleware.SetError(c, err)
		return
	}

	resp := UpdateProfileResponse{
		UserID:      result.UserID,
		DisplayName: result.DisplayName,
		AvatarURL:   result.AvatarURL,
		BirthDay:    result.BirthDay,
		Bio:         result.Bio,
	}

	response.Success(c, http.StatusOK, resp)
}

// CheckEmailAvailability handles GET /api/v1/auth/check-email?email=...
func (h *Handler) CheckEmailAvailability(c *gin.Context) {
	ctx := c.Request.Context()
	email := c.Query("email")

	if email == "" {
		middleware.SetError(c, sharederrors.NewAppError(
			sharederrors.CodeInvalidParameter,
			"Tham số không hợp lệ",
		).WithMetadata("field", "email"))
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
		middleware.SetError(c, sharederrors.NewAppError(
			sharederrors.CodeInvalidParameter,
			"Tham số không hợp lệ",
		).WithMetadata("field", "username"))
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
