package http

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"web-clean/domain"
	"web-clean/internal/application/service"
	"web-clean/internal/domain/usecase"
)

// UserHandler handles HTTP requests for user operations
// This is the interface/delivery layer that handles HTTP concerns only
type UserHandler struct {
	userUseCase usecase.UserUseCase
	logger      domain.Log
}

// NewUserHandler creates a new user handler
func NewUserHandler(userUseCase usecase.UserUseCase, logger domain.Log) *UserHandler {
	return &UserHandler{
		userUseCase: userUseCase,
		logger:      logger,
	}
}

// CreateUserRequest represents the HTTP request for creating a user
type CreateUserRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Username string `json:"username" binding:"required,min=3,max=50"`
	Name     string `json:"name" binding:"required,min=1,max=100"`
}

// UserResponse represents the HTTP response for user data
type UserResponse struct {
	ID        string `json:"id"`
	Email     string `json:"email"`
	Username  string `json:"username"`
	Name      string `json:"name"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}

// ListUsersResponse represents the HTTP response for listing users
type ListUsersResponse struct {
	Users   []UserResponse `json:"users"`
	Total   int64          `json:"total"`
	Offset  int            `json:"offset"`
	Limit   int            `json:"limit"`
	HasMore bool           `json:"has_more"`
}

// ErrorResponse represents error response
type ErrorResponse struct {
	Error   string `json:"error"`
	Message string `json:"message,omitempty"`
}

// CreateUser handles POST /users
func (h *UserHandler) CreateUser(c *gin.Context) {
	var req CreateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Warnw("Invalid request for create user", "error", err)
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "invalid_request",
			Message: err.Error(),
		})
		return
	}

	// Convert HTTP request to use case request
	useCaseReq := usecase.CreateUserRequest{
		Email:    req.Email,
		Username: req.Username,
		Name:     req.Name,
	}

	// Call use case
	user, err := h.userUseCase.CreateUser(c.Request.Context(), useCaseReq)
	if err != nil {
		h.handleError(c, err)
		return
	}

	// Convert domain entity to HTTP response
	response := UserResponse{
		ID:        user.ID.String(),
		Email:     user.Email,
		Username:  user.Username,
		Name:      user.Name,
		CreatedAt: user.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
		UpdatedAt: user.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
	}

	c.JSON(http.StatusCreated, response)
}

// GetUserByID handles GET /users/:id
func (h *UserHandler) GetUserByID(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		h.logger.Warnw("Invalid user ID format", "id", idStr, "error", err)
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "invalid_id",
			Message: "Invalid user ID format",
		})
		return
	}

	// Call use case
	user, err := h.userUseCase.GetUserByID(c.Request.Context(), id)
	if err != nil {
		h.handleError(c, err)
		return
	}

	// Convert domain entity to HTTP response
	response := UserResponse{
		ID:        user.ID.String(),
		Email:     user.Email,
		Username:  user.Username,
		Name:      user.Name,
		CreatedAt: user.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
		UpdatedAt: user.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
	}

	c.JSON(http.StatusOK, response)
}

// UpdateUserProfile handles PUT /users/:id
func (h *UserHandler) UpdateUserProfile(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		h.logger.Warnw("Invalid user ID format", "id", idStr, "error", err)
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "invalid_id",
			Message: "Invalid user ID format",
		})
		return
	}

	var req struct {
		Name string `json:"name" binding:"required,min=1,max=100"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Warnw("Invalid request for update user", "error", err)
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "invalid_request",
			Message: err.Error(),
		})
		return
	}

	// Convert HTTP request to use case request
	useCaseReq := usecase.UpdateUserProfileRequest{
		ID:   id,
		Name: req.Name,
	}

	// Call use case
	user, err := h.userUseCase.UpdateUserProfile(c.Request.Context(), useCaseReq)
	if err != nil {
		h.handleError(c, err)
		return
	}

	// Convert domain entity to HTTP response
	response := UserResponse{
		ID:        user.ID.String(),
		Email:     user.Email,
		Username:  user.Username,
		Name:      user.Name,
		CreatedAt: user.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
		UpdatedAt: user.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
	}

	c.JSON(http.StatusOK, response)
}

// DeleteUser handles DELETE /users/:id
func (h *UserHandler) DeleteUser(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		h.logger.Warnw("Invalid user ID format", "id", idStr, "error", err)
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "invalid_id",
			Message: "Invalid user ID format",
		})
		return
	}

	// Call use case
	err = h.userUseCase.DeleteUser(c.Request.Context(), id)
	if err != nil {
		h.handleError(c, err)
		return
	}

	c.JSON(http.StatusNoContent, nil)
}

// ListUsers handles GET /users
func (h *UserHandler) ListUsers(c *gin.Context) {
	// Parse query parameters
	offsetStr := c.DefaultQuery("offset", "0")
	limitStr := c.DefaultQuery("limit", "10")

	offset, err := strconv.Atoi(offsetStr)
	if err != nil || offset < 0 {
		h.logger.Warnw("Invalid offset parameter", "offset", offsetStr)
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "invalid_offset",
			Message: "Offset must be a non-negative integer",
		})
		return
	}

	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit <= 0 || limit > 100 {
		h.logger.Warnw("Invalid limit parameter", "limit", limitStr)
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "invalid_limit",
			Message: "Limit must be a positive integer between 1 and 100",
		})
		return
	}

	// Convert HTTP request to use case request
	useCaseReq := usecase.ListUsersRequest{
		Offset: offset,
		Limit:  limit,
	}

	// Call use case
	result, err := h.userUseCase.ListUsers(c.Request.Context(), useCaseReq)
	if err != nil {
		h.handleError(c, err)
		return
	}

	// Convert domain response to HTTP response
	users := make([]UserResponse, len(result.Users))
	for i, user := range result.Users {
		users[i] = UserResponse{
			ID:        user.ID.String(),
			Email:     user.Email,
			Username:  user.Username,
			Name:      user.Name,
			CreatedAt: user.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
			UpdatedAt: user.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
		}
	}

	response := ListUsersResponse{
		Users:   users,
		Total:   result.Total,
		Offset:  result.Offset,
		Limit:   result.Limit,
		HasMore: result.HasMore,
	}

	c.JSON(http.StatusOK, response)
}

// handleError converts use case errors to appropriate HTTP responses
func (h *UserHandler) handleError(c *gin.Context, err error) {
	switch err {
	case service.ErrUserNotFound:
		c.JSON(http.StatusNotFound, ErrorResponse{
			Error:   "user_not_found",
			Message: "User not found",
		})
	case service.ErrUserAlreadyExists:
		c.JSON(http.StatusConflict, ErrorResponse{
			Error:   "user_already_exists",
			Message: "User with email or username already exists",
		})
	case service.ErrInvalidUserData:
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "invalid_user_data",
			Message: "Invalid user data provided",
		})
	default:
		h.logger.Errorw("Internal server error", "error", err)
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error:   "internal_server_error",
			Message: "An internal error occurred",
		})
	}
}
