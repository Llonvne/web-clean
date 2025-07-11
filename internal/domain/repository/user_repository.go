package repository

import (
	"context"
	"github.com/google/uuid"
	"web-clean/internal/domain/entity"
)

// UserRepository defines the contract for user data access
// This interface belongs to the domain layer and will be implemented by infrastructure layer
type UserRepository interface {
	// Create stores a new user
	Create(ctx context.Context, user *entity.User) error

	// GetByID retrieves a user by their ID
	GetByID(ctx context.Context, id uuid.UUID) (*entity.User, error)

	// GetByEmail retrieves a user by their email
	GetByEmail(ctx context.Context, email string) (*entity.User, error)

	// GetByUsername retrieves a user by their username
	GetByUsername(ctx context.Context, username string) (*entity.User, error)

	// Update updates an existing user
	Update(ctx context.Context, user *entity.User) error

	// Delete removes a user by ID
	Delete(ctx context.Context, id uuid.UUID) error

	// List retrieves users with pagination
	List(ctx context.Context, offset, limit int) ([]*entity.User, error)

	// Count returns the total number of users
	Count(ctx context.Context) (int64, error)
}
