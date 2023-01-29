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

type UserRepository interface {
	GetAll(ctx context.Context) ([]*User, error)
	GetByID(ctx context.Context, id int64) (*User, error)
}
