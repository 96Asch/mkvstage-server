package domain

import (
	"context"
	"time"

	"gorm.io/gorm"
)

type User struct {
	ID           int64          `json:"id"`
	Email        string         `json:"email" binding:"email" gorm:"unique"`
	Password     string         `json:"password" binding:"required"`
	FirstName    string         `json:"first_name" binding:"required"`
	LastName     string         `json:"last_name" binding:"required"`
	Permission   string         `json:"permission" binding:"required"`
	ProfileColor string         `json:"profile_color" binding:"required"`
	UpdatedAt    time.Time      `json:"last_modified"`
	DeletedAt    gorm.DeletedAt `json:"-"`
}

type PublicUser struct {
	ID           int64     `json:"id"`
	Email        string    `json:"email" binding:"email" gorm:"unique"`
	FirstName    string    `json:"first_name" binding:"required"`
	LastName     string    `json:"last_name" binding:"required"`
	Permission   string    `json:"permission" binding:"required"`
	ProfileColor string    `json:"profile_color" binding:"required"`
	UpdatedAt    time.Time `json:"last_modified"`
}

type UserService interface {
	FetchByID(ctx context.Context, id int64) (*User, error)
	FetchAll(ctx context.Context) (*[]PublicUser, error)
	Store(ctx context.Context, user *User) error
	// Remove(ctx context.Context, user *User) error
}

type UserRepository interface {
	Creator[User]
	Getter[User]
}
