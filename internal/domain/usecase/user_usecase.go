package usecase

import (
	"context"
	"github.com/google/uuid"
	"web-clean/internal/domain/entity"
)

// UserUseCase defines the business operations for user management
// This interface belongs to the domain layer and contains business logic
type UserUseCase interface {
	// CreateUser creates a new user with validation
	CreateUser(ctx context.Context, req CreateUserRequest) (*entity.User, error)

	// GetUserByID retrieves a user by ID
	GetUserByID(ctx context.Context, id uuid.UUID) (*entity.User, error)

	// GetUserByEmail retrieves a user by email
	GetUserByEmail(ctx context.Context, email string) (*entity.User, error)

	// UpdateUserProfile updates user profile information
	UpdateUserProfile(ctx context.Context, req UpdateUserProfileRequest) (*entity.User, error)

	// DeleteUser removes a user
	DeleteUser(ctx context.Context, id uuid.UUID) error

	// ListUsers retrieves paginated list of users
	ListUsers(ctx context.Context, req ListUsersRequest) (*ListUsersResponse, error)
}

// CreateUserRequest represents the request to create a new user
type CreateUserRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Username string `json:"username" validate:"required,min=3,max=50"`
	Name     string `json:"name" validate:"required,min=1,max=100"`
}

// UpdateUserProfileRequest represents the request to update user profile
type UpdateUserProfileRequest struct {
	ID   uuid.UUID `json:"id" validate:"required"`
	Name string    `json:"name" validate:"required,min=1,max=100"`
}

// ListUsersRequest represents the request to list users with pagination
type ListUsersRequest struct {
	Offset int `json:"offset" validate:"min=0"`
	Limit  int `json:"limit" validate:"min=1,max=100"`
}

// ListUsersResponse represents the response for listing users
type ListUsersResponse struct {
	Users   []*entity.User `json:"users"`
	Total   int64          `json:"total"`
	Offset  int            `json:"offset"`
	Limit   int            `json:"limit"`
	HasMore bool           `json:"has_more"`
}
