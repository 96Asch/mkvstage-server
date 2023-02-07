package domain

import (
	"context"
	"time"

	"gorm.io/gorm"
)

type User struct {
	ID           int64          `json:"id"`
	Email        string         `json:"email" gorm:"unique"`
	Password     string         `json:"-"`
	FirstName    string         `json:"first_name"`
	LastName     string         `json:"last_name"`
	Permission   string         `json:"permission"`
	ProfileColor string         `json:"profile_color"`
	UpdatedAt    time.Time      `json:"last_modified"`
	DeletedAt    gorm.DeletedAt `json:"-"`
}

type UserService interface {
	Fetcher[User]
	Store(ctx context.Context, user *User) error
	Update(ctx context.Context, user *User) error

	// Remove(ctx context.Context, user *User) error
}

type UserRepository interface {
	Creator[User]
	Getter[User]
	Update(ctx context.Context, user *User) error
}
