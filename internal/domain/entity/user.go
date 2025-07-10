package entity

import (
	"time"
	"github.com/google/uuid"
)

// User represents the core business entity for users
type User struct {
	ID        uuid.UUID `json:"id"`
	Email     string    `json:"email"`
	Username  string    `json:"username"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// NewUser creates a new user entity with generated ID and timestamps
func NewUser(email, username, name string) *User {
	now := time.Now()
	return &User{
		ID:        uuid.New(),
		Email:     email,
		Username:  username,
		Name:      name,
		CreatedAt: now,
		UpdatedAt: now,
	}
}

// UpdateProfile updates user profile information
func (u *User) UpdateProfile(name string) {
	u.Name = name
	u.UpdatedAt = time.Now()
}

// IsValid validates the user entity
func (u *User) IsValid() bool {
	return u.ID != uuid.Nil && 
		   u.Email != "" && 
		   u.Username != "" && 
		   u.Name != ""
}