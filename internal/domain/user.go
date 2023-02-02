package domain

import (
	"context"
	"time"

	"gorm.io/gorm"
)

type User struct {
	ID           uint           `json:"id"`
	Email        string         `json:"email" binding:"email" gorm:"unique"`
	Password     string         `json:"password"`
	FirstName    string         `json:"first_name"`
	LastName     string         `json:"last_name"`
	Permission   string         `json:"permission"`
	ProfileColor string         `binding:"required"`
	UpdatedAt    time.Time      `json:"last_modified"`
	DeletedAt    gorm.DeletedAt `json:"-"`
}

type UserService interface {
	Getter[User]
	Store(ctx context.Context, user *User) error
	Remove(ctx context.Context, user *User) error
}

type UserRepository interface {
	Creator[User]
	Getter[User]
}
