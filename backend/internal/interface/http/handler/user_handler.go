package handler

import (
	"net/http"

	"github.com/english-coach/backend/internal/domain/user/port"
	"github.com/english-coach/backend/internal/domain/user/usecase/command"
	"github.com/english-coach/backend/internal/domain/user/usecase/query"
	"github.com/english-coach/backend/internal/shared/response"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// UserHandler handles user-related HTTP requests
type UserHandler struct {
	registerUC      *command.RegisterUserUseCase
	loginUC         *command.LoginUseCase
	getProfileUC    *query.GetUserProfileUseCase
	updateProfileUC *command.UpdateUserProfileUseCase
	userRepo        port.UserRepository
	profileRepo     port.UserProfileRepository
	logger          *zap.Logger
}

// NewUserHandler creates a new user handler
func NewUserHandler(
	registerUC *command.RegisterUserUseCase,
	loginUC *command.LoginUseCase,
	getProfileUC *query.GetUserProfileUseCase,
	updateProfileUC *command.UpdateUserProfileUseCase,
	userRepo port.UserRepository,
	profileRepo port.UserProfileRepository,
	logger *zap.Logger,
) *UserHandler {
	return &UserHandler{
		registerUC:      registerUC,
		loginUC:         loginUC,
		getProfileUC:    getProfileUC,
		updateProfileUC: updateProfileUC,
		userRepo:        userRepo,
		profileRepo:     profileRepo,
		logger:          logger,
	}
}

// RegisterRequest represents the request body for user registration
type RegisterRequest struct {
	DisplayName *string `json:"display_name,omitempty" binding:"omitempty,max=100"`
	Email       *string `json:"email,omitempty" binding:"omitempty,email"`
	Username    *string `json:"username,omitempty" binding:"omitempty,min=3,max=100"`
	Password    string  `json:"password" binding:"required,min=6"`
}

// Register handles POST /api/v1/auth/register
func (h *UserHandler) Register(c *gin.Context) {
	ctx := c.Request.Context()

	var req RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ErrorResponse(c, http.StatusBadRequest,
			"INVALID_REQUEST",
			"Invalid request body",
			err.Error(),
		)
		return
	}

	// Validate that at least email or username is provided
	if (req.Email == nil || *req.Email == "") && (req.Username == nil || *req.Username == "") {
		response.ErrorResponse(c, http.StatusBadRequest,
			"INVALID_REQUEST",
			"Either email or username is required",
			nil,
		)
		return
	}

	result, err := h.registerUC.Execute(ctx, command.RegisterUserInput{
		DisplayName: req.DisplayName,
		Email:       req.Email,
		Username:    req.Username,
		Password:    req.Password,
	})

	if err != nil {
		h.logger.Error("failed to register user",
			zap.Error(err),
			zap.String("path", c.Request.URL.Path),
		)

		statusCode := http.StatusInternalServerError
		code := "INTERNAL_ERROR"
		message := "Failed to register user"

		if err == command.ErrEmailRequired {
			statusCode = http.StatusBadRequest
			code = "EMAIL_REQUIRED"
			message = err.Error()
		} else if err == command.ErrEmailExists {
			statusCode = http.StatusConflict
			code = "EMAIL_EXISTS"
			message = err.Error()
		} else if err == command.ErrUsernameExists {
			statusCode = http.StatusConflict
			code = "USERNAME_EXISTS"
			message = err.Error()
		} else if err == command.ErrInvalidPassword {
			statusCode = http.StatusBadRequest
			code = "INVALID_PASSWORD"
			message = err.Error()
		}

		response.ErrorResponse(c, statusCode, code, message, nil)
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

// LoginRequest represents the request body for user login
type LoginRequest struct {
	Email    *string `json:"email,omitempty"`
	Username *string `json:"username,omitempty"`
	Password string  `json:"password" binding:"required"`
}

// Login handles POST /api/v1/auth/login
func (h *UserHandler) Login(c *gin.Context) {
	ctx := c.Request.Context()

	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ErrorResponse(c, http.StatusBadRequest,
			"INVALID_REQUEST",
			"Invalid request body",
			err.Error(),
		)
		return
	}

	// Validate that at least email or username is provided
	if (req.Email == nil || *req.Email == "") && (req.Username == nil || *req.Username == "") {
		response.ErrorResponse(c, http.StatusBadRequest,
			"INVALID_REQUEST",
			"Either email or username is required",
			nil,
		)
		return
	}

	result, err := h.loginUC.Execute(ctx, command.LoginInput{
		Email:    req.Email,
		Username: req.Username,
		Password: req.Password,
	})

	if err != nil {
		h.logger.Error("failed to login user",
			zap.Error(err),
			zap.String("path", c.Request.URL.Path),
		)

		statusCode := http.StatusUnauthorized
		code := "UNAUTHORIZED"
		message := "Invalid credentials"

		if err == command.ErrInvalidCredentials {
			statusCode = http.StatusUnauthorized
			code = "INVALID_CREDENTIALS"
			message = err.Error()
		} else if err == command.ErrUserInactive {
			statusCode = http.StatusForbidden
			code = "USER_INACTIVE"
			message = err.Error()
		}

		response.ErrorResponse(c, statusCode, code, message, nil)
		return
	}

	response.Success(c, http.StatusOK, result)
}

// GetProfile handles GET /api/v1/users/profile
func (h *UserHandler) GetProfile(c *gin.Context) {
	ctx := c.Request.Context()

	// Get user ID from context (set by auth middleware)
	userID, exists := c.Get("user_id")
	if !exists {
		response.ErrorResponse(c, http.StatusUnauthorized,
			"UNAUTHORIZED",
			"User not authenticated",
			nil,
		)
		return
	}

	userIDInt64, ok := userID.(int64)
	if !ok {
		response.ErrorResponse(c, http.StatusInternalServerError,
			"INTERNAL_ERROR",
			"Invalid user ID",
			nil,
		)
		return
	}

	profile, err := h.getProfileUC.Execute(ctx, userIDInt64)
	if err != nil {
		h.logger.Error("failed to get user profile",
			zap.Error(err),
			zap.Int64("user_id", userIDInt64),
			zap.String("path", c.Request.URL.Path),
		)

		statusCode := http.StatusInternalServerError
		code := "INTERNAL_ERROR"
		message := "Failed to get user profile"

		if err == query.ErrProfileNotFound {
			statusCode = http.StatusNotFound
			code = "PROFILE_NOT_FOUND"
			message = err.Error()
		}

		response.ErrorResponse(c, statusCode, code, message, nil)
		return
	}

	response.Success(c, http.StatusOK, profile)
}

// UpdateProfileRequest represents the request body for updating user profile
type UpdateProfileRequest struct {
	DisplayName *string `json:"display_name,omitempty" binding:"omitempty,max=100"`
	AvatarURL   *string `json:"avatar_url,omitempty" binding:"omitempty,url,max=500"`
	BirthDay    *string `json:"birth_day,omitempty" binding:"omitempty,datetime=2006-01-02"`
	Bio         *string `json:"bio,omitempty"`
}

// UpdateProfile handles PUT /api/v1/users/profile
func (h *UserHandler) UpdateProfile(c *gin.Context) {
	ctx := c.Request.Context()

	// Get user ID from context (set by auth middleware)
	userID, exists := c.Get("user_id")
	if !exists {
		response.ErrorResponse(c, http.StatusUnauthorized,
			"UNAUTHORIZED",
			"User not authenticated",
			nil,
		)
		return
	}

	userIDInt64, ok := userID.(int64)
	if !ok {
		response.ErrorResponse(c, http.StatusInternalServerError,
			"INTERNAL_ERROR",
			"Invalid user ID",
			nil,
		)
		return
	}

	var req UpdateProfileRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ErrorResponse(c, http.StatusBadRequest,
			"INVALID_REQUEST",
			"Invalid request body",
			err.Error(),
		)
		return
	}

	result, err := h.updateProfileUC.Execute(ctx, userIDInt64, command.UpdateUserProfileInput{
		DisplayName: req.DisplayName,
		AvatarURL:   req.AvatarURL,
		BirthDay:    req.BirthDay,
		Bio:         req.Bio,
	})

	if err != nil {
		h.logger.Error("failed to update user profile",
			zap.Error(err),
			zap.Int64("user_id", userIDInt64),
			zap.String("path", c.Request.URL.Path),
		)

		response.ErrorResponse(c, http.StatusInternalServerError,
			"INTERNAL_ERROR",
			"Failed to update user profile",
			nil,
		)
		return
	}

	response.Success(c, http.StatusOK, result)
}

// CheckEmailAvailability handles GET /api/v1/auth/check-email?email=...
func (h *UserHandler) CheckEmailAvailability(c *gin.Context) {
	ctx := c.Request.Context()
	email := c.Query("email")

	if email == "" {
		response.ErrorResponse(c, http.StatusBadRequest,
			"INVALID_PARAMETER",
			"Email parameter is required",
			nil,
		)
		return
	}

	exists, err := h.userRepo.CheckEmailExists(ctx, email)
	if err != nil {
		h.logger.Error("failed to check email availability",
			zap.Error(err),
			zap.String("email", email),
		)
		response.ErrorResponse(c, http.StatusInternalServerError,
			"INTERNAL_ERROR",
			"Failed to check email availability",
			nil,
		)
		return
	}

	response.Success(c, http.StatusOK, gin.H{
		"available": !exists,
		"exists":    exists,
	})
}

// CheckUsernameAvailability handles GET /api/v1/auth/check-username?username=...
func (h *UserHandler) CheckUsernameAvailability(c *gin.Context) {
	ctx := c.Request.Context()
	username := c.Query("username")

	if username == "" {
		response.ErrorResponse(c, http.StatusBadRequest,
			"INVALID_PARAMETER",
			"Username parameter is required",
			nil,
		)
		return
	}

	exists, err := h.userRepo.CheckUsernameExists(ctx, username)
	if err != nil {
		h.logger.Error("failed to check username availability",
			zap.Error(err),
			zap.String("username", username),
		)
		response.ErrorResponse(c, http.StatusInternalServerError,
			"INTERNAL_ERROR",
			"Failed to check username availability",
			nil,
		)
		return
	}

	response.Success(c, http.StatusOK, gin.H{
		"available": !exists,
		"exists":    exists,
	})
}
