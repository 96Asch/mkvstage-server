package domain

import (
	"context"
	"time"

	"gorm.io/gorm"
)

const (
	ADMIN = iota + 1
	MEMBER
	GUEST
)

type User struct {
	ID           int64          `json:"id"`
	Email        string         `json:"email" gorm:"unique"`
	Password     string         `json:"-"`
	FirstName    string         `json:"first_name"`
	LastName     string         `json:"last_name"`
	Permission   int8           `json:"permission"`
	ProfileColor string         `json:"profile_color"`
	UpdatedAt    time.Time      `json:"last_modified"`
	DeletedAt    gorm.DeletedAt `json:"-"`
}

func (u User) HasClearance() bool {
	return u.Permission == ADMIN
}

type UserService interface {
	Fetcher[User]
	Store(ctx context.Context, user *User) error
	Update(ctx context.Context, user *User) error
	Remove(ctx context.Context, user *User, id int64) error
}

type UserRepository interface {
	Creator[User]
	Getter[User]
	Update(ctx context.Context, user *User) error
	Delete(ctx context.Context, id int64) error
}
