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
	ExtractUser(ctx context.Context, access string) (*User, error)
	CreateAccess(ctx context.Context, currentRefresh string) (*AccessToken, error)
	CreateRefresh(ctx context.Context, uid int64, currentRefresh string) (*RefreshToken, error)
	RemoveRefresh(ctx context.Context, uid int64, refresh string) error
	RemoveAllRefresh(ctx context.Context, uid int64) error
}

type TokenRepository interface {
	GetAll(ctx context.Context, uid int64) (*[]RefreshToken, error)
	Create(ctx context.Context, token *RefreshToken) error
	Delete(ctx context.Context, uid int64, refresh string) error
	DeleteAll(ctx context.Context, uid int64) error
}
