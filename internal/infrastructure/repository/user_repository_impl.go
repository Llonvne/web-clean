package repository

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"web-clean/infra/database"
	"web-clean/internal/domain/entity"
	"web-clean/internal/domain/repository"
)

// UserModel represents the database model for users
// This is the infrastructure concern - how we store users in the database
type UserModel struct {
	ID        uuid.UUID `gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	Email     string    `gorm:"type:varchar(255);uniqueIndex;not null"`
	Username  string    `gorm:"type:varchar(50);uniqueIndex;not null"`
	Name      string    `gorm:"type:varchar(100);not null"`
	CreatedAt time.Time `gorm:"autoCreateTime"`
	UpdatedAt time.Time `gorm:"autoUpdateTime"`
}

// TableName specifies the table name for GORM
func (UserModel) TableName() string {
	return "users"
}

// ToEntity converts database model to domain entity
func (m *UserModel) ToEntity() *entity.User {
	return &entity.User{
		ID:        m.ID,
		Email:     m.Email,
		Username:  m.Username,
		Name:      m.Name,
		CreatedAt: m.CreatedAt,
		UpdatedAt: m.UpdatedAt,
	}
}

// FromEntity converts domain entity to database model
func (m *UserModel) FromEntity(user *entity.User) {
	m.ID = user.ID
	m.Email = user.Email
	m.Username = user.Username
	m.Name = user.Name
	m.CreatedAt = user.CreatedAt
	m.UpdatedAt = user.UpdatedAt
}

// UserRepositoryImpl implements the UserRepository interface
// This is the infrastructure layer implementation
type UserRepositoryImpl struct {
	db database.Database
}

// NewUserRepository creates a new user repository implementation
func NewUserRepository(db database.Database) repository.UserRepository {
	return &UserRepositoryImpl{
		db: db,
	}
}

// Register the schema for auto-migration
func init() {
	database.RegisterSchema(UserModel{})
}

// Create stores a new user in the database
func (r *UserRepositoryImpl) Create(ctx context.Context, user *entity.User) error {
	model := &UserModel{}
	model.FromEntity(user)

	return r.db.Transaction(func(tx *gorm.DB) error {
		return tx.WithContext(ctx).Create(model).Error
	})
}

// GetByID retrieves a user by their ID
func (r *UserRepositoryImpl) GetByID(ctx context.Context, id uuid.UUID) (*entity.User, error) {
	var model UserModel

	err := r.db.Transaction(func(tx *gorm.DB) error {
		return tx.WithContext(ctx).Where("id = ?", id).First(&model).Error
	})

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}

	return model.ToEntity(), nil
}

// GetByEmail retrieves a user by their email
func (r *UserRepositoryImpl) GetByEmail(ctx context.Context, email string) (*entity.User, error) {
	var model UserModel

	err := r.db.Transaction(func(tx *gorm.DB) error {
		return tx.WithContext(ctx).Where("email = ?", email).First(&model).Error
	})

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}

	return model.ToEntity(), nil
}

// GetByUsername retrieves a user by their username
func (r *UserRepositoryImpl) GetByUsername(ctx context.Context, username string) (*entity.User, error) {
	var model UserModel

	err := r.db.Transaction(func(tx *gorm.DB) error {
		return tx.WithContext(ctx).Where("username = ?", username).First(&model).Error
	})

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}

	return model.ToEntity(), nil
}

// Update updates an existing user in the database
func (r *UserRepositoryImpl) Update(ctx context.Context, user *entity.User) error {
	model := &UserModel{}
	model.FromEntity(user)

	return r.db.Transaction(func(tx *gorm.DB) error {
		return tx.WithContext(ctx).Model(&UserModel{}).Where("id = ?", user.ID).Updates(model).Error
	})
}

// Delete removes a user from the database
func (r *UserRepositoryImpl) Delete(ctx context.Context, id uuid.UUID) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		return tx.WithContext(ctx).Delete(&UserModel{}, "id = ?", id).Error
	})
}

// List retrieves users with pagination
func (r *UserRepositoryImpl) List(ctx context.Context, offset, limit int) ([]*entity.User, error) {
	var models []UserModel

	err := r.db.Transaction(func(tx *gorm.DB) error {
		return tx.WithContext(ctx).
			Offset(offset).
			Limit(limit).
			Order("created_at DESC").
			Find(&models).Error
	})

	if err != nil {
		return nil, err
	}

	users := make([]*entity.User, len(models))
	for i, model := range models {
		users[i] = model.ToEntity()
	}

	return users, nil
}

// Count returns the total number of users
func (r *UserRepositoryImpl) Count(ctx context.Context) (int64, error) {
	var count int64

	err := r.db.Transaction(func(tx *gorm.DB) error {
		return tx.WithContext(ctx).Model(&UserModel{}).Count(&count).Error
	})

	return count, err
}
