package service

import (
	"context"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"web-clean/domain"
	"web-clean/internal/domain/entity"
	"web-clean/internal/domain/repository"
	"web-clean/internal/domain/usecase"
)

var (
	ErrUserNotFound      = errors.New("user not found")
	ErrUserAlreadyExists = errors.New("user already exists")
	ErrInvalidUserData   = errors.New("invalid user data")
)

// UserService implements the UserUseCase interface
// This is the application layer that contains business logic
type UserService struct {
	userRepo repository.UserRepository
	logger   domain.Log
}

// NewUserService creates a new UserService instance
func NewUserService(userRepo repository.UserRepository, logger domain.Log) usecase.UserUseCase {
	return &UserService{
		userRepo: userRepo,
		logger:   logger,
	}
}

// CreateUser creates a new user with business validation
func (s *UserService) CreateUser(ctx context.Context, req usecase.CreateUserRequest) (*entity.User, error) {
	s.logger.Infow("CreateUser", "email", req.Email, "username", req.Username)

	// Business rule: Check if user with email already exists
	existingUser, err := s.userRepo.GetByEmail(ctx, req.Email)
	if err == nil && existingUser != nil {
		s.logger.Warnw("User creation failed - email already exists", "email", req.Email)
		return nil, ErrUserAlreadyExists
	}

	// Business rule: Check if username already exists
	existingUser, err = s.userRepo.GetByUsername(ctx, req.Username)
	if err == nil && existingUser != nil {
		s.logger.Warnw("User creation failed - username already exists", "username", req.Username)
		return nil, ErrUserAlreadyExists
	}

	// Create new user entity
	user := entity.NewUser(req.Email, req.Username, req.Name)

	// Business validation
	if !user.IsValid() {
		s.logger.Errorw("User creation failed - invalid data", "user", user)
		return nil, ErrInvalidUserData
	}

	// Store the user
	if err := s.userRepo.Create(ctx, user); err != nil {
		s.logger.Errorw("Failed to create user", "error", err, "user", user)
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	s.logger.Infow("User created successfully", "userID", user.ID, "email", user.Email)
	return user, nil
}

// GetUserByID retrieves a user by ID
func (s *UserService) GetUserByID(ctx context.Context, id uuid.UUID) (*entity.User, error) {
	s.logger.Infow("GetUserByID", "userID", id)

	user, err := s.userRepo.GetByID(ctx, id)
	if err != nil {
		s.logger.Errorw("Failed to get user by ID", "error", err, "userID", id)
		return nil, ErrUserNotFound
	}

	if user == nil {
		s.logger.Warnw("User not found", "userID", id)
		return nil, ErrUserNotFound
	}

	return user, nil
}

// GetUserByEmail retrieves a user by email
func (s *UserService) GetUserByEmail(ctx context.Context, email string) (*entity.User, error) {
	s.logger.Infow("GetUserByEmail", "email", email)

	user, err := s.userRepo.GetByEmail(ctx, email)
	if err != nil {
		s.logger.Errorw("Failed to get user by email", "error", err, "email", email)
		return nil, ErrUserNotFound
	}

	if user == nil {
		s.logger.Warnw("User not found", "email", email)
		return nil, ErrUserNotFound
	}

	return user, nil
}

// UpdateUserProfile updates user profile information
func (s *UserService) UpdateUserProfile(ctx context.Context, req usecase.UpdateUserProfileRequest) (*entity.User, error) {
	s.logger.Infow("UpdateUserProfile", "userID", req.ID, "name", req.Name)

	// Get existing user
	user, err := s.userRepo.GetByID(ctx, req.ID)
	if err != nil {
		s.logger.Errorw("Failed to get user for update", "error", err, "userID", req.ID)
		return nil, ErrUserNotFound
	}

	if user == nil {
		s.logger.Warnw("User not found for update", "userID", req.ID)
		return nil, ErrUserNotFound
	}

	// Apply business logic for profile update
	user.UpdateProfile(req.Name)

	// Business validation
	if !user.IsValid() {
		s.logger.Errorw("User update failed - invalid data", "user", user)
		return nil, ErrInvalidUserData
	}

	// Update in repository
	if err := s.userRepo.Update(ctx, user); err != nil {
		s.logger.Errorw("Failed to update user", "error", err, "userID", req.ID)
		return nil, fmt.Errorf("failed to update user: %w", err)
	}

	s.logger.Infow("User profile updated successfully", "userID", user.ID)
	return user, nil
}

// DeleteUser removes a user
func (s *UserService) DeleteUser(ctx context.Context, id uuid.UUID) error {
	s.logger.Infow("DeleteUser", "userID", id)

	// Business rule: Check if user exists before deletion
	user, err := s.userRepo.GetByID(ctx, id)
	if err != nil || user == nil {
		s.logger.Warnw("User not found for deletion", "userID", id)
		return ErrUserNotFound
	}

	// Perform deletion
	if err := s.userRepo.Delete(ctx, id); err != nil {
		s.logger.Errorw("Failed to delete user", "error", err, "userID", id)
		return fmt.Errorf("failed to delete user: %w", err)
	}

	s.logger.Infow("User deleted successfully", "userID", id)
	return nil
}

// ListUsers retrieves paginated list of users
func (s *UserService) ListUsers(ctx context.Context, req usecase.ListUsersRequest) (*usecase.ListUsersResponse, error) {
	s.logger.Infow("ListUsers", "offset", req.Offset, "limit", req.Limit)

	// Business rule: Set default limit if not provided
	if req.Limit <= 0 {
		req.Limit = 10
	}

	// Business rule: Maximum limit
	if req.Limit > 100 {
		req.Limit = 100
	}

	// Get total count
	total, err := s.userRepo.Count(ctx)
	if err != nil {
		s.logger.Errorw("Failed to get user count", "error", err)
		return nil, fmt.Errorf("failed to get user count: %w", err)
	}

	// Get users
	users, err := s.userRepo.List(ctx, req.Offset, req.Limit)
	if err != nil {
		s.logger.Errorw("Failed to list users", "error", err, "offset", req.Offset, "limit", req.Limit)
		return nil, fmt.Errorf("failed to list users: %w", err)
	}

	hasMore := int64(req.Offset+req.Limit) < total

	response := &usecase.ListUsersResponse{
		Users:   users,
		Total:   total,
		Offset:  req.Offset,
		Limit:   req.Limit,
		HasMore: hasMore,
	}

	s.logger.Infow("Users listed successfully", "total", total, "returned", len(users))
	return response, nil
}
