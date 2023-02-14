package domain

import (
	"context"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"
)

type RefreshToken struct {
	ID      uuid.UUID `json:"-"`
	UserID  int64     `json:"-"`
	Refresh string    `json:"refresh"`
}

type AccessToken struct {
	Access string `json:"access"`
}

type Tokens struct {
	AccessToken
	RefreshToken
}

type AccessTokenClaims struct {
	User *User `json:"user"`
	jwt.RegisteredClaims
}

type RefreshTokenClaims struct {
	UID int64 `json:"uid"`
	jwt.RegisteredClaims
}

type Claims interface {
	AccessTokenClaims | RefreshTokenClaims
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
}
