package domain

import (
	"context"
	"time"

	"gorm.io/gorm"
)

type Clearance int

const (
	ADMIN Clearance = iota + 1
	EDITOR
	MEMBER
	GUEST
)

type User struct {
	ID           int64          `json:"id"`
	Email        string         `json:"email" gorm:"unique"`
	Password     string         `json:"-"`
	FirstName    string         `json:"first_name"`
	LastName     string         `json:"last_name"`
	Permission   Clearance      `json:"permission"`
	ProfileColor string         `json:"profile_color"`
	UpdatedAt    time.Time      `json:"last_modified"`
	DeletedAt    gorm.DeletedAt `json:"-"`
}

func (u User) HasClearance(clearance Clearance) bool {
	return u.Permission <= clearance
}

type UserService interface {
	Fetcher[User]
	Store(ctx context.Context, user *User) error
	Update(ctx context.Context, user *User) error
	Remove(ctx context.Context, user *User, id int64) (int64, error)
	Authorize(ctx context.Context, email, password string) (*User, error)
	SetPermission(ctx context.Context, permission Clearance, recipient, principal *User) (*User, error)
}

type UserRepository interface {
	Creator[User]
	Getter[User]
	GetByEmail(ctx context.Context, email string) (*User, error)
	Update(ctx context.Context, user *User) error
	Delete(ctx context.Context, id int64) error
}
