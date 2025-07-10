package service

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"web-clean/internal/domain/entity"
	"web-clean/internal/domain/usecase"
)

// MockUserRepository is a mock implementation of UserRepository for testing
type MockUserRepository struct {
	mock.Mock
}

func (m *MockUserRepository) Create(ctx context.Context, user *entity.User) error {
	args := m.Called(ctx, user)
	return args.Error(0)
}

func (m *MockUserRepository) GetByID(ctx context.Context, id uuid.UUID) (*entity.User, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.User), args.Error(1)
}

func (m *MockUserRepository) GetByEmail(ctx context.Context, email string) (*entity.User, error) {
	args := m.Called(ctx, email)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.User), args.Error(1)
}

func (m *MockUserRepository) GetByUsername(ctx context.Context, username string) (*entity.User, error) {
	args := m.Called(ctx, username)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.User), args.Error(1)
}

func (m *MockUserRepository) Update(ctx context.Context, user *entity.User) error {
	args := m.Called(ctx, user)
	return args.Error(0)
}

func (m *MockUserRepository) Delete(ctx context.Context, id uuid.UUID) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockUserRepository) List(ctx context.Context, offset, limit int) ([]*entity.User, error) {
	args := m.Called(ctx, offset, limit)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*entity.User), args.Error(1)
}

func (m *MockUserRepository) Count(ctx context.Context) (int64, error) {
	args := m.Called(ctx)
	return args.Get(0).(int64), args.Error(1)
}

// MockLogger is a mock implementation of domain.Log for testing
type MockLogger struct {
	mock.Mock
}

func (m *MockLogger) Debug(args ...interface{})                            {}
func (m *MockLogger) Info(args ...interface{})                             {}
func (m *MockLogger) Warn(args ...interface{})                             {}
func (m *MockLogger) Error(args ...interface{})                            {}
func (m *MockLogger) DPanic(args ...interface{})                           {}
func (m *MockLogger) Panic(args ...interface{})                            {}
func (m *MockLogger) Fatal(args ...interface{})                            {}
func (m *MockLogger) Debugf(template string, args ...interface{})          {}
func (m *MockLogger) Infof(template string, args ...interface{})           {}
func (m *MockLogger) Warnf(template string, args ...interface{})           {}
func (m *MockLogger) Errorf(template string, args ...interface{})          {}
func (m *MockLogger) DPanicf(template string, args ...interface{})         {}
func (m *MockLogger) Panicf(template string, args ...interface{})          {}
func (m *MockLogger) Fatalf(template string, args ...interface{})          {}
func (m *MockLogger) Debugw(msg string, keysAndValues ...interface{})      {}
func (m *MockLogger) Infow(msg string, keysAndValues ...interface{})       {}
func (m *MockLogger) Warnw(msg string, keysAndValues ...interface{})       {}
func (m *MockLogger) Errorw(msg string, keysAndValues ...interface{})      {}
func (m *MockLogger) DPanicw(msg string, keysAndValues ...interface{})     {}
func (m *MockLogger) Panicw(msg string, keysAndValues ...interface{})      {}
func (m *MockLogger) Fatalw(msg string, keysAndValues ...interface{})      {}

func TestUserService_CreateUser_Success(t *testing.T) {
	// Arrange
	mockRepo := new(MockUserRepository)
	mockLogger := new(MockLogger)
	service := NewUserService(mockRepo, mockLogger)

	ctx := context.Background()
	req := usecase.CreateUserRequest{
		Email:    "test@example.com",
		Username: "testuser",
		Name:     "Test User",
	}

	// Mock expectations - user doesn't exist
	mockRepo.On("GetByEmail", ctx, req.Email).Return(nil, errors.New("not found"))
	mockRepo.On("GetByUsername", ctx, req.Username).Return(nil, errors.New("not found"))
	mockRepo.On("Create", ctx, mock.AnythingOfType("*entity.User")).Return(nil)

	// Act
	user, err := service.CreateUser(ctx, req)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, user)
	assert.Equal(t, req.Email, user.Email)
	assert.Equal(t, req.Username, user.Username)
	assert.Equal(t, req.Name, user.Name)
	assert.NotEqual(t, uuid.Nil, user.ID)
	mockRepo.AssertExpectations(t)
}

func TestUserService_CreateUser_EmailExists(t *testing.T) {
	// Arrange
	mockRepo := new(MockUserRepository)
	mockLogger := new(MockLogger)
	service := NewUserService(mockRepo, mockLogger)

	ctx := context.Background()
	req := usecase.CreateUserRequest{
		Email:    "existing@example.com",
		Username: "newuser",
		Name:     "New User",
	}

	existingUser := &entity.User{
		ID:       uuid.New(),
		Email:    req.Email,
		Username: "olduser",
		Name:     "Existing User",
	}

	// Mock expectations - user with email exists
	mockRepo.On("GetByEmail", ctx, req.Email).Return(existingUser, nil)

	// Act
	user, err := service.CreateUser(ctx, req)

	// Assert
	assert.Error(t, err)
	assert.Equal(t, ErrUserAlreadyExists, err)
	assert.Nil(t, user)
	mockRepo.AssertExpectations(t)
}

func TestUserService_CreateUser_UsernameExists(t *testing.T) {
	// Arrange
	mockRepo := new(MockUserRepository)
	mockLogger := new(MockLogger)
	service := NewUserService(mockRepo, mockLogger)

	ctx := context.Background()
	req := usecase.CreateUserRequest{
		Email:    "new@example.com",
		Username: "existinguser",
		Name:     "New User",
	}

	existingUser := &entity.User{
		ID:       uuid.New(),
		Email:    "existing@example.com",
		Username: req.Username,
		Name:     "Existing User",
	}

	// Mock expectations
	mockRepo.On("GetByEmail", ctx, req.Email).Return(nil, errors.New("not found"))
	mockRepo.On("GetByUsername", ctx, req.Username).Return(existingUser, nil)

	// Act
	user, err := service.CreateUser(ctx, req)

	// Assert
	assert.Error(t, err)
	assert.Equal(t, ErrUserAlreadyExists, err)
	assert.Nil(t, user)
	mockRepo.AssertExpectations(t)
}

func TestUserService_GetUserByID_Success(t *testing.T) {
	// Arrange
	mockRepo := new(MockUserRepository)
	mockLogger := new(MockLogger)
	service := NewUserService(mockRepo, mockLogger)

	ctx := context.Background()
	userID := uuid.New()
	expectedUser := &entity.User{
		ID:        userID,
		Email:     "test@example.com",
		Username:  "testuser",
		Name:      "Test User",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	// Mock expectations
	mockRepo.On("GetByID", ctx, userID).Return(expectedUser, nil)

	// Act
	user, err := service.GetUserByID(ctx, userID)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, user)
	assert.Equal(t, expectedUser.ID, user.ID)
	assert.Equal(t, expectedUser.Email, user.Email)
	assert.Equal(t, expectedUser.Username, user.Username)
	assert.Equal(t, expectedUser.Name, user.Name)
	mockRepo.AssertExpectations(t)
}

func TestUserService_GetUserByID_NotFound(t *testing.T) {
	// Arrange
	mockRepo := new(MockUserRepository)
	mockLogger := new(MockLogger)
	service := NewUserService(mockRepo, mockLogger)

	ctx := context.Background()
	userID := uuid.New()

	// Mock expectations
	mockRepo.On("GetByID", ctx, userID).Return(nil, errors.New("not found"))

	// Act
	user, err := service.GetUserByID(ctx, userID)

	// Assert
	assert.Error(t, err)
	assert.Equal(t, ErrUserNotFound, err)
	assert.Nil(t, user)
	mockRepo.AssertExpectations(t)
}

func TestUserService_UpdateUserProfile_Success(t *testing.T) {
	// Arrange
	mockRepo := new(MockUserRepository)
	mockLogger := new(MockLogger)
	service := NewUserService(mockRepo, mockLogger)

	ctx := context.Background()
	userID := uuid.New()
	req := usecase.UpdateUserProfileRequest{
		ID:   userID,
		Name: "Updated Name",
	}

	originalUpdatedAt := time.Now().Add(-1 * time.Hour)
	existingUser := &entity.User{
		ID:        userID,
		Email:     "test@example.com",
		Username:  "testuser",
		Name:      "Old Name",
		CreatedAt: time.Now().Add(-24 * time.Hour),
		UpdatedAt: originalUpdatedAt,
	}

	// Mock expectations
	mockRepo.On("GetByID", ctx, userID).Return(existingUser, nil)
	mockRepo.On("Update", ctx, mock.AnythingOfType("*entity.User")).Return(nil)

	// Act
	user, err := service.UpdateUserProfile(ctx, req)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, user)
	assert.Equal(t, req.Name, user.Name)
	assert.Equal(t, "test@example.com", user.Email)
	assert.Equal(t, "testuser", user.Username)
	// UpdatedAt should be updated (check that it's after the original time)
	assert.True(t, user.UpdatedAt.After(originalUpdatedAt))
	mockRepo.AssertExpectations(t)
}

func TestUserService_DeleteUser_Success(t *testing.T) {
	// Arrange
	mockRepo := new(MockUserRepository)
	mockLogger := new(MockLogger)
	service := NewUserService(mockRepo, mockLogger)

	ctx := context.Background()
	userID := uuid.New()

	existingUser := &entity.User{
		ID:       userID,
		Email:    "test@example.com",
		Username: "testuser",
		Name:     "Test User",
	}

	// Mock expectations
	mockRepo.On("GetByID", ctx, userID).Return(existingUser, nil)
	mockRepo.On("Delete", ctx, userID).Return(nil)

	// Act
	err := service.DeleteUser(ctx, userID)

	// Assert
	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestUserService_DeleteUser_NotFound(t *testing.T) {
	// Arrange
	mockRepo := new(MockUserRepository)
	mockLogger := new(MockLogger)
	service := NewUserService(mockRepo, mockLogger)

	ctx := context.Background()
	userID := uuid.New()

	// Mock expectations
	mockRepo.On("GetByID", ctx, userID).Return(nil, errors.New("not found"))

	// Act
	err := service.DeleteUser(ctx, userID)

	// Assert
	assert.Error(t, err)
	assert.Equal(t, ErrUserNotFound, err)
	mockRepo.AssertExpectations(t)
}

func TestUserService_ListUsers_Success(t *testing.T) {
	// Arrange
	mockRepo := new(MockUserRepository)
	mockLogger := new(MockLogger)
	service := NewUserService(mockRepo, mockLogger)

	ctx := context.Background()
	req := usecase.ListUsersRequest{
		Offset: 0,
		Limit:  10,
	}

	expectedUsers := []*entity.User{
		{
			ID:       uuid.New(),
			Email:    "user1@example.com",
			Username: "user1",
			Name:     "User One",
		},
		{
			ID:       uuid.New(),
			Email:    "user2@example.com",
			Username: "user2",
			Name:     "User Two",
		},
	}

	// Mock expectations
	mockRepo.On("Count", ctx).Return(int64(25), nil)
	mockRepo.On("List", ctx, req.Offset, req.Limit).Return(expectedUsers, nil)

	// Act
	response, err := service.ListUsers(ctx, req)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, response)
	assert.Equal(t, int64(25), response.Total)
	assert.Equal(t, len(expectedUsers), len(response.Users))
	assert.Equal(t, req.Offset, response.Offset)
	assert.Equal(t, req.Limit, response.Limit)
	assert.True(t, response.HasMore) // 25 total > 10 limit
	mockRepo.AssertExpectations(t)
}

func TestUserService_ListUsers_WithDefaultLimit(t *testing.T) {
	// Arrange
	mockRepo := new(MockUserRepository)
	mockLogger := new(MockLogger)
	service := NewUserService(mockRepo, mockLogger)

	ctx := context.Background()
	req := usecase.ListUsersRequest{
		Offset: 0,
		Limit:  0, // Should be defaulted to 10
	}

	// Mock expectations
	mockRepo.On("Count", ctx).Return(int64(5), nil)
	mockRepo.On("List", ctx, 0, 10).Return([]*entity.User{}, nil) // Expects limit to be 10

	// Act
	response, err := service.ListUsers(ctx, req)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, 10, response.Limit) // Should be defaulted
	mockRepo.AssertExpectations(t)
}

func TestUserService_ListUsers_MaxLimitEnforced(t *testing.T) {
	// Arrange
	mockRepo := new(MockUserRepository)
	mockLogger := new(MockLogger)
	service := NewUserService(mockRepo, mockLogger)

	ctx := context.Background()
	req := usecase.ListUsersRequest{
		Offset: 0,
		Limit:  200, // Should be capped at 100
	}

	// Mock expectations
	mockRepo.On("Count", ctx).Return(int64(5), nil)
	mockRepo.On("List", ctx, 0, 100).Return([]*entity.User{}, nil) // Expects limit to be 100

	// Act
	response, err := service.ListUsers(ctx, req)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, 100, response.Limit) // Should be capped
	mockRepo.AssertExpectations(t)
}