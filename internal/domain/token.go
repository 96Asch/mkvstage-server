package domain

import (
	"context"
	"time"

	"github.com/google/uuid"
)

type RefreshToken struct {
	ID                 uuid.UUID     `json:"-"`
	UserID             int64         `json:"-"`
	Refresh            string        `json:"refresh"`
	ExpirationDuration time.Duration `json:"-"`
}

type AccessToken struct {
	Access string `json:"access"`
}

type Tokens struct {
	AccessToken
	RefreshToken
}

type TokenConfig struct {
	IAT         time.Time
	ExpDuration time.Duration
	Secret      string
}

type TokenService interface {
	ExtractUser(ctx context.Context, token *AccessToken) (*User, error)
	CreateAccess(ctx context.Context, user *User) (*AccessToken, error)
	CreateRefresh(ctx context.Context, user *User) (*RefreshToken, error)
}

type TokenRepository interface {
	GetAll(ctx context.Context, uid int64) (*[]RefreshToken, error)
	Create(ctx context.Context, token *RefreshToken) error
	Delete(ctx context.Context, token *RefreshToken) error
	DeleteAll(ctx context.Context, uid int64) error
}
